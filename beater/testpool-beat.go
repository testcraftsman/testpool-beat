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
    Profile   string    `json: "profle"`
    Level     string    `json: "level"`
    Vm_max    int       `json: "vm_max"`
    Vm_count  int       `json: "vm_count"`
    Timestamp time.Time `json: "RFC3339Nano"`
}


// profileRead: Read profile log content.
// Follows generator pattern by returning a channel.
func profilRead(profile_path string) (<-chan Message, error) {

  fhndl, err := os.Open(profile_path)
  if err != nil {
    return nil, err
  }

  generator := make (chan Message)
  scanner := bufio.NewScanner(fhndl)

  go func () {
    // Defer must be in here because entering go routine will
    // cause defer.
    defer fhndl.Close()

    for scanner.Scan() {
      data := []byte(scanner.Text())
      var msg Message

      err := json.Unmarshal(data, &msg)
  
      if err == nil {
        generator <- msg
      }
    }
    close(generator)
  }()
   
  if err := scanner.Err(); err != nil {
    return nil, err
  }

  return generator, nil
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
