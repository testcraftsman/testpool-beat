package beater

import "testing"
import "fmt"
import "log"
import "time"
import "github.com/stretchr/testify/assert"

// TestTimeFormat learning how to use the time.Parse function
func TestTimeFormat(t *testing.T) {

  str := "2017-06-21T04:24:50.452736"
  t1, err := time.Parse("2006-01-02T15:04:05", str)
  assert.Equal(t, err, nil, "parse time")
  assert.Equal(t, t1.Year(), 2017, "year does not match")
  assert.Equal(t, int(t1.Month()), 6, "month does not match")
}

// TestProfile implementing reading a profile.
func TestProfile(t *testing.T) {

  err := profilRead("profile.log")
  assert.Equal(t, err, nil, "read profile")
}
