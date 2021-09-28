package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/orted-org/conman/conman"
)

var config *conman.Config
var logger = log.Default()

func ServerInit(incomingConfig *conman.Config, addr string) {
	config = incomingConfig
	r := chi.NewRouter()
	r.Get("/{key}", getConfig)
	r.Put("/watch", setWatchFileDuration)
	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}

	logger.Println("Server running on address", addr)
	PrintInitialLogs()
	log.Fatal(srv.ListenAndServe())
}
func PrintInitialLogs() {

	if config.GetFileName() == "" {
		logger.Println("Running in", "Memory", "mode")
	} else {
		logger.Println("Running in", "File", "mode")
		logger.Println("Filename:", config.GetFileName())
	}

	if config.GetCurrentWatchInterval() == -1 {
		logger.Println("Not watching for file changes")
	} else {
		logger.Println("Watching file changes at interval of ", config.GetCurrentWatchInterval(), "seconds")
	}
}
