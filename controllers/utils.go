package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func jsonDecode(r *http.Request, out interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

func responseJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(data)
}

func responseOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "")
	w.WriteHeader(200)
}

func responseCreated(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "")
	w.WriteHeader(201)
}
