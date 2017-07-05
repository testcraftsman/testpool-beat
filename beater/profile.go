package beater

import (
    "bufio"
    "time"
    "os"
    "encoding/json"
)


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
