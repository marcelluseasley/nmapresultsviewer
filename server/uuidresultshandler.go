package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

func uuidResultsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["uuid"] == "" || !uuidScanLookup(vars["uuid"]) {
		log.Printf("no scan results with: %s", vars["uuid"])
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
