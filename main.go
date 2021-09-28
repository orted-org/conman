package main

import (
	"os"
	"strconv"
	"time"

	api "github.com/orted-org/conman/api"
	"github.com/orted-org/conman/conman"
)

func main() {
	config := conman.NewConfig()

	// setting file name
	fileName := os.Getenv("CONMAN_FILENAME")
	if len(fileName) != 0 {
		// no file name passed
		config.SetFilename(fileName)
	}

	// check if watch file
	watchDuration := os.Getenv("CONMAN_WATCH_DURATION")
	if len(watchDuration) > 0 {
		dur, err := strconv.Atoi(watchDuration)
		if err != nil {
			panic("error duration set")
		}
		if dur > 0 {
			if len(config.GetFileName()) == 0 {
				// duration is passed to watch file but file name not specified
				panic("file name not specified to watch")
			} else {
				config.SetFromFile()
				config.WatchFileChanges(time.Second * time.Duration(dur))
			}
		}
	}

	api.ServerInit(config, "localhost:4000")

}
