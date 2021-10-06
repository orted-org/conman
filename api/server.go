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

	// routes
	r.Get("/", getAll)
	r.Post("/", setConfig)
	r.Get("/{key}", getConfig)
	r.Put("/watch", setWatchFileDuration)
	r.Get("/stats", getStats)

	// server config
	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}

	// initial logs
	logger.Println("Server running on address", addr)
	logger.Println("Current filename:", config.GetFileName())
	if config.GetCurrentWatchInterval() == -1 {
		logger.Println("Not watching for file changes")
	} else {
		logger.Println("Watching file changes at interval of ", config.GetCurrentWatchInterval(), "seconds")
	}

	// listen and serve
	log.Fatal(srv.ListenAndServe())
}
