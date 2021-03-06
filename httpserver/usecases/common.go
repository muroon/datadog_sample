package usecases

import (
	"encoding/json"
	"net/http"

	"github.com/muroon/datadog_sample/config"

	"github.com/muroon/datadog_sample/httpserver/jsonmodel"
)

var conf config.Config

func Init() error {
	var err error
	conf, err = config.GetConfig()
	if err != nil {
		return err
	}

	err = initDB()
	if err != nil {
		return err
	}

	err = openDB()
	if err != nil {
		return err
	}

	err = openGrpc()
	return err
}

func End() {
	_ = closeDB()
	_ = closeGrpc()
}

func renderJSON(w http.ResponseWriter, b []byte) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func renderErrorJSON(w http.ResponseWriter, err error) {
	res := &jsonmodel.ErrorResult{Message: err.Error()}
	b, err := json.Marshal(res)
	if err != nil {
		return
	}

	renderJSON(w, b)
}
