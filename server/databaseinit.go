package main


import (
	"database/sql"
	"io/ioutil"
	"log"
	_ "github.com/mattn/go-sqlite3"
)


func createdatabase() {

	database, err := sql.Open("sqlite3", "dbtest1.db")
	if err != nil {
		log.Fatalf("sql.Open error: %v", err)
	}

	sqlDBCreation, err := ioutil.ReadFile("nmapdata.db.sql")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	// if tables haven't been created, then create them 
	
	
	_, err = database.Exec(string(sqlDBCreation))
	if err != nil {
		log.Fatalf("database.Exec table creation error: %v", err)
	}
	
	database.Close()

}