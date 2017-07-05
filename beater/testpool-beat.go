package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/testcraftsman/testpool-beat/config"
)

type TestpoolBeat struct {
	done       chan struct{}
	config     config.Config
	client     publisher.Client
	profileLog string
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
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		timestamp := common.Time(time.Now())
		profiles, err := profileRead(bt.profileLog)
		if err != nil {
			logp.Err(fmt.Sprintf("%v", err))

			for item := range profiles {

				event := common.MapStr{
					"@timestamp": timestamp,
					"type":       b.Name,
					"counter":    counter,
					"profile":    item.Profile,
					"level":      item.Level,
					"vm_max":     item.Vm_max,
					"vm_count":   item.Vm_count,
					"timestamp":  item.Timestamp,
				}
				bt.client.PublishEvent(event)
			}
			logp.Info("Event sent")
			counter++
		}
	}
}

func (bt *TestpoolBeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
