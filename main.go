// server

package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	nmapserver "github.com/marcelluseasley/nmapresultsviewer/server"
)

func NewNMAPRouter () *mux.Router{
	r := mux.NewRouter()

	r.HandleFunc("/{uuid}", nmapserver.UUIDResultsHandler).Methods("GET")
	r.HandleFunc("/", nmapserver.IndexHandler).Methods("GET")
	r.HandleFunc("/v1/nmap", nmapserver.SubmitNMAPHandler).Methods("PUT", "GET")

	return r
}

func main() {

	nmapserver.Createdatabase()

	r := NewNMAPRouter()
	
	http.Handle("/", r)
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
