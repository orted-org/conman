package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func getConfig(rw http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	reqConfig := config.Get(key)
	if reqConfig == nil {
		jsonRes(rw, http.StatusNotFound, "key not found", nil)
		return
	}
	jsonRes(rw, http.StatusOK, "", reqConfig)
}
func setWatchFileDuration(rw http.ResponseWriter, r *http.Request) {
	rawDuration := r.URL.Query().Get("duration")
	if len(rawDuration) == 0 {
		jsonRes(rw, http.StatusBadRequest, "duration not found in query params", nil)
		return
	}
	duration, err := strconv.Atoi(rawDuration)
	if err != nil {
		jsonRes(rw, http.StatusBadRequest, "invalid duration", nil)
		return
	}

	if config.GetCurrentWatchInterval() != -1 {
		// the file is being watched already
		config.UnWatchFileChanges()
	}

	// if duration passed is less or equal to 0, stop watching for changes in file
	if duration <= 0 {
		jsonRes(rw, http.StatusOK, "not watching to file changes", nil)
		return
	}
	config.WatchFileChanges(time.Second * time.Duration(duration))
	jsonRes(rw, http.StatusOK, fmt.Sprintf("watching file changes at interval of %d seconds", duration), nil)
}
