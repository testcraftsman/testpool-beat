package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/testcraftsman/testpool-beat/beater"
)

func main() {
	err := beat.Run("testpool-beat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
