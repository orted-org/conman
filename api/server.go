package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/orted-org/conman/conman"
)

var config *conman.Config

func configInit() {
	config = conman.NewConfig()
	config.SetFilename("./config.json")
	config.SetFromFile()
}
func ServerInit(addr string) {
	configInit()
	r := chi.NewRouter()
	r.Get("/{key}", getConfig)

	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Fatal(srv.ListenAndServe())
}

func getConfig(rw http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	reqConfig := config.Get(key)
	if reqConfig == nil {
		jsonRes(rw, http.StatusNotFound, "key not found", nil)
		return
	}
	jsonRes(rw, http.StatusOK, "", reqConfig)
}
