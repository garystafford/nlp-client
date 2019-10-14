// author: Gary A. Stafford
// site: https://programmaticponderings.com
// license: MIT License
// purpose: Simple echo-based microservice: client

package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

const (
	port = ":8080"
	url  = "http://rake:8080/keywords"
)

// Prediction
type prediction struct {
	Prediction int `json:"prediction"`
}

func main() {
	http.HandleFunc("/", makePrediction)
	log.Fatal(http.ListenAndServe(port, nil))
}

func makePrediction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := http.Post(url, "application/json", r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	w.Write(body[:])
}
