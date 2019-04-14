package nmapserver

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type ScanData struct {
	UUIDScan     string
	Scanargs     string
	Scanstart    string
	Scantype     string
	Scanprotocol string
	Scanservices string
	Scanend      string
	Summary      string
}
type RowData struct {
	UUIDHost  string
	IP        string
	HostState string
	HReason   string
	Hostname  string
	Port      string
	PState    string
	PReason   string
	Service   string
	Method    string
}

type TemplateData struct {
	Scan ScanData
	Rows map[string][]RowData
}

func UUIDResultsHandler(w http.ResponseWriter, r *http.Request) {

	uuidR := strings.TrimLeft(r.RequestURI, "/")
	log.Printf("VALUE: %v", uuidR)

	// if no uuid or uuid not found in database, redirect to index page
	if !uuidScanLookup(uuidR) {

		t, err := template.ParseFiles("server/templates/index.html")
		if err != nil {
			log.Printf("Error opening index template: %v", err)
		}

		err = t.Execute(w, nil)
		if err != nil {
			log.Printf("Error executing index template: %v", err)
		}
	} else {
		log.Println("found uuid! generate UI page")

		var TD TemplateData
		var ips []string
		database, err := sql.Open("sqlite3", "dbtest1.db")
		if err != nil {
			log.Fatalf("sql.Open error: %v", err)
		}

		ipScan, err := database.Query(fmt.Sprintf(`
SELECT DISTINCT ip
FROM portdata
WHERE uuid = '%s'
ORDER BY ip;
	`, uuidR))
		if err != nil {
			log.Fatalf("database query error: %v", err)
		}
		var i string
		for ipScan.Next() {
			ipScan.Scan(&i)
			ips = append(ips, i)

		}

		TD.Rows = make(map[string][]RowData)
		for _, ii := range ips {
			TD.Rows[ii] = []RowData{}
		}

		rowsScan, err := database.Query(fmt.Sprintf(`
select *
FROM scandata
WHERE scandata.uuid = '%s'
	`, uuidR))
		if err != nil {
			log.Fatalf("database query error: %v", err)
		}
		var uuidScan, scanargs, scanstart, scantype, scanprotocol, scanservices, scanend, summary string
		for rowsScan.Next() {
			rowsScan.Scan(&uuidScan, &scanargs, &scanstart, &scantype, &scanprotocol, &scanservices, &scanend, &summary)
			s := ScanData{uuidScan, scanargs, scanstart, scantype, scanprotocol, scanservices, scanend, summary}
			TD.Scan = s

		}

		var uuidHost, ip, hostState, hReason, hostname, port, pState, pReason, service, method string
		var temp []RowData
		for _, ipr := range ips {
			rowsHost, err := database.Query(fmt.Sprintf(`
		SELECT hostdata.uuid,
		hostdata.ip as ip,
		hostdata.host_state as host_state,
		hostdata.reason as h_reason,
		hostdata.hostname as hostname,
		portdata.port as port,
		portdata.state as p_state,
		portdata.reason as p_reason,
		portdata.service as service,
		portdata.method as method
		FROM hostdata
		INNER JOIN portdata 
		ON hostdata.uuid = portdata.uuid 
		where hostdata.uuid = '%s'
		AND hostdata.ip = portdata.ip
		AND hostdata.ip = '%s'
		ORDER BY ip ASC;`, uuidR, ipr))

			if err != nil {
				log.Fatalf("database query error: %v", err)
			}

			for rowsHost.Next() {
				rowsHost.Scan(&uuidHost, &ip, &hostState, &hReason, &hostname, &port, &pState, &pReason, &service, &method)
				temp = append(temp, RowData{uuidHost, ip, hostState, hReason, hostname, port, pState, pReason, service, method})

			}
			TD.Rows[ipr] = temp
			temp = nil
		}

		t, err := template.ParseFiles("server/templates/nmapresults.html")
		if err != nil {
			log.Printf("Error opening nmapresults template: %v", err)
		}
		err = t.Execute(w, TD)
		if err != nil {
			log.Printf("Error executing nmapresults template: %v", err)
		}

	}

}

func uuidScanLookup(scanuuid string) bool {
	database, err := sql.Open("sqlite3", "dbtest1.db")
	if err != nil {
		log.Fatalf("sql.Open error: %v", err)
	}

	scan, err := database.Query(fmt.Sprintf(`
	SELECT DISTINCT uuid
	FROM scandata
	WHERE uuid = '%s' LIMIT 1;`, scanuuid))
	if err != nil {
		log.Printf("Database Query failed: %v", err)
	}

	if scan.Next() {
		database.Close()
		return true
	} else {
		database.Close()
		return false
	}

}
