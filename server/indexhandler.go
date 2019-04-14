package nmapserver

import (
	"html/template"
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("server/templates/index.html")
	if err != nil {
		log.Printf("Error opening index template: %v", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing index template: %v", err)
	}
}
