package beater

import "encoding/json"
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
	sample := "{\"profile\": \"example\", \"level\": \"info\", \"vm_max\": 3, \"vm_count\": 0, \"timestamp\": \"2017-04-15T05:30:58.362953Z\"}"

	err := json.Unmarshal([]byte(sample), &msg)
	assert.Equal(t, err, nil, err)

}

// TestProfile implementing reading a profile.
func TestProfile(t *testing.T) {

	generator, err := profileRead("profile.log")

	assert.Equal(t, err, nil, "read profile")
	assert.NotEqual(t, generator, nil, "generator failed")

	var counter = 0
	for item := range generator {
		assert.NotEqual(t, item, nil, item)
		counter++
	}
	assert.NotEqual(t, 0, 1)
}

func TestConfigRead(t *testing.T) {

	fpath, err := configRead()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, fpath, "/var/log/testpool/profile.log", err)
}
