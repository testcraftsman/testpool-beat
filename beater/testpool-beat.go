package beater

import (
    "bufio"
	"fmt"
	"time"
    "os"
    "encoding/json"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/testcraftsman/testpool-beat/config"
)

type TestpoolBeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
    lastIndexTime time.Time
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &TestpoolBeat{
		done: make(chan struct{}),
		config: config,
	}
	return bt, nil
}

type Message struct {
    Profile   string
    Level     string
    Vm_max    int
    Vm_count  int
    Timestamp time.Time `json:"RFC3339Nano, string"`
}

func profilRead(profile_path string) error {

  fhndl, err := os.Open(profile_path)

  if err != nil {
    return err
  }

  defer fhndl.Close()

  scanner := bufio.NewScanner(fhndl)
  for scanner.Scan() {
    data := []byte(scanner.Text())
    var msg Message
    if err := json.Unmarshal(data, &msg); err == nil {
        fmt.Println(msg)
    }
  }
   
  if err := scanner.Err(); err != nil {
    return err
  }

  return nil
}

func (bt *TestpoolBeat) Run(b *beat.Beat) error {
	logp.Info("testpool-beat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"counter":    counter,
            "profile":    "mark",
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *TestpoolBeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
