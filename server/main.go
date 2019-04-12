// server

package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	createdatabase()

	r := mux.NewRouter()

	r.HandleFunc("/{uuid}", uuidResultsHandler).Methods("GET")
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/v1/nmap", submitNMAPHandler).Methods("PUT", "GET")

	http.Handle("/", r)
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
