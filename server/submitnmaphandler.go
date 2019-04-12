package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// curl http://myservice --upload-file file.txt
func submitNMAPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}

		log.Print(string(body))
	} else if r.Method == "GET" {
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Printf("Error opening index template: %v", err)
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Printf("Error executing index template: %v", err)
		}
	}

}
