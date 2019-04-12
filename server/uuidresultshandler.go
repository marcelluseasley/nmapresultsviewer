package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func uuidResultsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// if no uuid or uuid not found in database, redirect to index page
	if vars["uuid"] == "" || !uuidScanLookup(vars["uuid"]) { 
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Printf("Error opening index template: %v", err)
		}
		
		err = t.Execute(w, nil)
		if err != nil {
			log.Printf("Error executing index template: %v", err)
		}
	} else {
		log.Println("found it!")
	}

	fmt.Fprintf(w, vars["uuid"])

}

func uuidScanLookup(scanuuid string) bool {
	database, err := sql.Open("sqlite3", "../dbtest1.db")
	if err != nil {
		log.Fatalf("sql.Open error: %v", err)
	}

	scan, err := database.Query(fmt.Sprintf(`
	SELECT DISTINCT uuid
	FROM scandata
	WHERE uuid = '%s' LIMIT 1;`, scanuuid))

	if scan.Next() {
		database.Close()
		return true
	} else {
		database.Close()
		return false
	}

}
