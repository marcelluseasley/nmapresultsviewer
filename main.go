package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var scanID uuid.UUID
	scanID = uuid.New()
	log.Print(scanID)

	var nmapContent Nmaprun
	content, err := ioutil.ReadFile("nmap.results.xml")
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(content, &nmapContent)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", scanID.String())
	fmt.Printf("%s\n", nmapContent.Args)
	fmt.Printf("%s\n", nmapContent.Startstr)
	fmt.Printf("%s\n", nmapContent.Scaninfo.Type)
	fmt.Printf("%s\n", nmapContent.Scaninfo.Protocol)
	fmt.Printf("%s\n", nmapContent.Scaninfo.Services)
	fmt.Printf("%s\n", nmapContent.Runstats.Finished.Timestr)
	fmt.Printf("%s\n", nmapContent.Runstats.Finished.Summary)

	database, err := sql.Open("sqlite3", "./sqlite_idea/dbtest1.db")
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
				log.Fatalf("portStatement.Exec error: %v",err)
			}
		}
	}

	statement.Close()
	database.Close()
}
