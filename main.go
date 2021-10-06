package main

import (
	"log"
	"os"
	"strconv"
	"time"

	api "github.com/orted-org/conman/api"
	"github.com/orted-org/conman/conman"
)

func main() {
	logger := log.Default()
	config := conman.NewConfig()

	// setting file name
	fileName := os.Getenv("CONMAN_FILENAME")
	if len(fileName) != 0 {
		// file name is specified
		config.SetFilename(fileName)
	} else {
		logger.Println("File not specified")
		// no file name is specified
		_, err := os.Stat("./temp.json")
		if err != nil {
			// temp.json not present
			logger.Println("Creating temp.json in current directory to store configurations")
			_, err := os.Create("./temp.json")
			if err != nil {
				panic("no file specified and could not create a temp file to store configurations")
			}
		} else {
			// temp.json present
			logger.Println("Using temp.json of the current directory for storing configurations")
		}
		config.SetFilename("./temp.json")
	}
	config.SetFromFile()

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
	secret := os.Getenv("CONMAN_API_SECRET")
	if secret == "" {
		// using a default api secret
		secret = "secret@api"
	}
	port := os.Getenv("CONMAN_PORT")
	if port == "" {
		// using the default port as 4000
		port = "4000"
	}
	api.ServerInit(config, secret, "0.0.0.0:"+port)
}
