package beater

import "encoding/json"
import "fmt"
import "testing"
import "time"
import "github.com/stretchr/testify/assert"

// TestTimeFormat learning how to use the time.Parse function
func TestTimeFormat(t *testing.T) {

  str := "2017-06-21T04:24:50.452736"
  t1, err := time.Parse("2006-01-02T15:04:05", str)
  assert.Equal(t, err, nil, "parse time")
  assert.Equal(t, t1.Year(), 2017, "year does not match")
  assert.Equal(t, int(t1.Month()), 6,
               "month does not match")
}

// TestUnmarshal test unmarshal
func TestUnmarshal(t *testing.T) {

  var msg Message
  sample := "{\"profile\": \"example\", \"level\": \"info\", \"vm_max\": 3, \"vm_count\": 0, \"timestamp\": 2017-04-15T05:30:58.362953Z06:00}"

  fmt.Printf(string(sample))

  err := json.Unmarshal([]byte(sample), &msg)
  assert.Equal(t, err, nil, err)

}

// TestProfile implementing reading a profile.
func TestProfile(t *testing.T) {

  generator, err := profilRead("profile.log")

  assert.Equal(t, err, nil, "read profile")
  assert.NotEqual(t, generator, nil, "generator failed")

  fmt.Println("MARK: reading")
  for item := range generator {
    fmt.Println("MARK: in")
    fmt.Println(item)
  }
  fmt.Println("MARK: reading done")
}
