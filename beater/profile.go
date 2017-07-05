package beater

import "io/ioutil"
import "bufio"
import "time"
import "os"
import "encoding/json"
import "github.com/smallfish/simpleyaml"

type Message struct {
	Profile   string    `json: "profle"`
	Level     string    `json: "level"`
	Vm_max    int       `json: "vm_max"`
	Vm_count  int       `json: "vm_count"`
	Timestamp time.Time `json: "RFC3339Nano"`
}

// configRead: Read configuration to find profile log.
func configRead() (string, error) {

	fhndl, err := ioutil.ReadFile("/etc/testpool/testpool.yml")
	if err != nil {
		return "", err
	}

	root, err := simpleyaml.NewYaml(fhndl)
	if err != nil {
		return "", err
	}

	value := root.GetPath("tpldaemon", "profile", "log")
	return value.String()

}

// profileRead: Read profile log content.
// Follows generator pattern by returning a channel.
func profileRead(profile_path string) (<-chan Message, error) {

	fhndl, err := os.Open(profile_path)
	if err != nil {
		return nil, err
	}

	generator := make(chan Message)
	scanner := bufio.NewScanner(fhndl)

	go func() {
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
