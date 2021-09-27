package main

import (
	"encoding/json"
	"fmt"
	"time"

	conman "github.com/orted-org/conman/conman"
)

func main() {

	config := conman.NewConfig()
	config.SetFilename("./config.json")
	config.SetFromFile()
	config.WatchFileChanges(time.Second * 5)
	i := 0
	for {
		byteDate, _ := json.Marshal(config.GetAll())
		fmt.Println(string(byteDate))
		time.Sleep(time.Second * 5)
		i++
		if i == 3 {
			config.UnWatchFileChanges()
		}
	}
}
