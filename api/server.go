package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/orted-org/conman/conman"
)

var config *conman.Config
var logger = log.Default()

func ServerInit(incomingConfig *conman.Config, secret, addr string) {
	config = incomingConfig
	r := chi.NewRouter()

	// routes
	r.Get("/", auth(secret, getAll))
	r.Post("/", auth(secret, setConfig))
	r.Get("/{key}", auth(secret, getConfig))
	r.Put("/watch", auth(secret, setWatchFileDuration))
	r.Get("/stats", auth(secret, getStats))

	// server config
	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}

	// initial logs
	logger.Println("Current filename:", config.GetFileName())
	if config.GetCurrentWatchInterval() == -1 {
		logger.Println("Not watching for file changes")
	} else {
		logger.Println("Watching file changes at interval of ", config.GetCurrentWatchInterval(), "seconds")
	}

	// listen and serve
	logger.Println("Server running and listening at", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
