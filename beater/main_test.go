package beater

import "testing"
import "fmt"
import "log"

func TestProfile(t *testing.T) {
  fmt.Printf("hello")
  err := profilRead("profile.log")

  if err != nil {
    log.Fatal(err)
    t.Fail()
  }
}

