// server

package main

import (
	//"html/template"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"html/template"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/{uuid}", uuidResultsHandler).Methods("GET")
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/v1/nmap", submitNMAPHandler).Methods("PUT", "GET")

	http.Handle("/", r)
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error opening index template: %v", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing index template: %v", err)
	}

}
