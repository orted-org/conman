package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// handler for getting the current filename and watch duration
func getStats(rw http.ResponseWriter, r *http.Request) {
	stats := struct {
		Filename      string        `json:"file_name"`
		WatchDuration time.Duration `json:"watch_duration"`
	}{Filename: config.GetFileName(), WatchDuration: time.Duration(config.GetCurrentWatchInterval().Seconds())}
	jsonRes(rw, http.StatusOK, "", stats)
}

// handler for getting specified config
func getConfig(rw http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	reqConfig := config.Get(key)
	if reqConfig == nil {
		jsonRes(rw, http.StatusNotFound, "key not found", nil)
		return
	}
	jsonRes(rw, http.StatusOK, "", reqConfig)
}

// handler for getting all the configs
func getAll(rw http.ResponseWriter, r *http.Request) {
	jsonRes(rw, http.StatusOK, "", config.GetAll())
}

// handler for setting duration for file watch
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

// handler for setting the config
func setConfig(rw http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		jsonRes(rw, http.StatusBadRequest, "body could not be read", nil)
		return
	}

	// validate if the json is valid
	var temp interface{}
	if err := json.Unmarshal(b, &temp); err != nil {
		jsonRes(rw, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// setting the config to the config variable
	if err := config.Set(b); err != nil {
		jsonRes(rw, http.StatusInternalServerError, "could not set configuration", nil)
		return
	}

	if hard := r.URL.Query().Get("hard"); hard == "true" {
		// make changes to file also
		if err := config.WriteConfigToFile(b); err != nil {
			jsonRes(rw, http.StatusConflict, "in memory configuration changed but could not change configuration of file", nil)
			return
		} else {
			jsonRes(rw, http.StatusConflict, "changed hard and soft configuration", nil)
			return
		}
	}

	jsonRes(rw, http.StatusOK, "in memory configuration changed", nil)

}

// basic authentication middleware
func auth(secret string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const authBasic = "Basic"
		var h = r.Header.Get("Authorization")
		var key string = ""

		// Basic auth scheme.
		if strings.HasPrefix(h, authBasic) {
			payload, err := base64.StdEncoding.DecodeString(string(strings.Trim(h[len(authBasic):], " ")))
			if err != nil {
				jsonRes(w, http.StatusUnauthorized, "invalid Base64 value in Basic Authorization header",
					nil)
				return
			}

			key = string(payload)
		} else {
			jsonRes(w, http.StatusUnauthorized, "missing Basic Authorization header",
				nil)
			return

		}

		if key != secret {
			jsonRes(w, http.StatusUnauthorized, "invalid API credentials",
				nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
