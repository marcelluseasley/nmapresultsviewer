package main

import (
	"database/sql"
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"log"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// curl http://myservice --upload-file file.txt
func submitNMAPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var scanID uuid.UUID
		scanID = uuid.New()

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}

		var nmapContent Nmaprun
		err = xml.Unmarshal(content, &nmapContent)
		if err != nil {
			log.Fatal(err)
		}

		database, err := sql.Open("sqlite3", "dbtest1.db")
		if err != nil {
			log.Fatalf("sql.Open error: %v", err)
		}

		// SCANDATA
		statement, err := database.Prepare("INSERT INTO scandata (uuid, scanargs, scanstart, scantype, scanprotocol, scanservices, scanend, summary) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatalf("scandata table - database.Prepare error: %v", err)
		}

		_, err = statement.Exec(
			scanID.String(),
			nmapContent.Args,
			nmapContent.Startstr,
			nmapContent.Scaninfo.Type,
			nmapContent.Scaninfo.Protocol,
			nmapContent.Scaninfo.Services,
			nmapContent.Runstats.Finished.Timestr,
			nmapContent.Runstats.Finished.Summary)

		if err != nil {
			log.Fatal(err)
		}

		// HOSTDATA
		statement, err = database.Prepare("INSERT INTO hostdata (uuid, ip, host_state, reason, hostname) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatalf("hostdata - database.Prepare error: %v", err)
		}

		for _, host := range nmapContent.Host {
			_, err = statement.Exec(
				scanID.String(),
				host.Address.Addr,
				host.Status.State,
				host.Status.Reason,
				host.Hostnames.Hostname.Name)

			if err != nil {
				log.Fatal(err)
			}
			// PORTDATA - nested in Host
			portStatement, err := database.Prepare("INSERT INTO portdata (uuid, ip, port, state, reason, service, method) VALUES (?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				log.Fatalf("portdata - database.Prepare error: %v", err)
			}
			for _, port := range host.Ports.Port {
				_, err = portStatement.Exec(
					scanID.String(),
					host.Address.Addr,
					port.Portid,
					port.State.State,
					port.State.Reason,
					port.Service.Name,
					port.Service.Method)
				if err != nil {
					log.Fatalf("portStatement.Exec error: %v", err)
				}
			}
		}

		statement.Close()
		database.Close()

		uuidJSONResponse := fmt.Sprintf(`
		{
			"scan-id": "%s"
		}
		`,scanID.String())
		fmt.Fprintf(w,uuidJSONResponse)
		

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
