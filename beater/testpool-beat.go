package beater


import (
	"fmt"
	"time"
	"os"
	"io"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/testcraftsman/testpool-beat/config"
)

var service string = "service"
var tmp_log string = "/var/tmp/profile.log"

type TestpoolBeat struct {
	done       chan struct{}
	config     config.Config
	client     publisher.Client
	profileLog string
}

func copyfile(src string, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {

	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	////
	// Read configuration. If profile.log is not defined then
	// quit.
	profileLog, err := configRead()
	if err != nil {
		return nil, err
	}
	////

	bt := &TestpoolBeat{
		done:       make(chan struct{}),
		config:     config,
		profileLog: profileLog,
	}
	return bt, nil
}

func (bt *TestpoolBeat) Run(b *beat.Beat) error {
	logp.Info("testpool-beat is running! Hit CTRL-C to stop.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	for {
            select {
	    case <-bt.done:
	        return nil
            case <-ticker.C:
	    }

	    timestamp := common.Time(time.Now())

            if _, err := os.Stat(bt.profileLog); os.IsNotExist(err) {
	        logp.Debug(service, "log does not exist %s\n", bt.profileLog)
                continue
            }
            ////
            // TODO: confirm that moving file worked for 
            // structured log content generated in Python
            // Its possible for the testpool-daemon to maintain
            // a handle to the moved file.
            os.Remove(tmp_log)
            if err := copyfile(bt.profileLog, tmp_log); err != nil {
                logp.Err(err.Error())
                continue
            }
            if err := os.Truncate(bt.profileLog, 0); err != nil {
                logp.Err(err.Error())
            }
            ////

	    profiles, err := profileRead(tmp_log)
	    if err != nil {
	        logp.Err(err.Error())
            } else {
	        for item := range profiles {
                    event := common.MapStr{
		        "@timestamp": timestamp,
			"type":       b.Name,
			"profile":    item.Profile,
			"vm_max":     item.Vm_max,
			"vm_count":   item.Vm_count,
			"timestamp":  item.Timestamp,
		    }
		    bt.client.PublishEvent(event)
	        }
                os.Remove(tmp_log)
	    }
	}
}

func (bt *TestpoolBeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
