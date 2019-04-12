package main

import (
	"database/sql"
	"fmt"
	"log"

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

func main() {
	var td TemplateData
	var ips []string
	database, err := sql.Open("sqlite3", "../dbtest1.db")
	if err != nil {
		log.Fatalf("sql.Open error: %v", err)
	}

	ipScan, err := database.Query(`
SELECT DISTINCT ip
FROM portdata
WHERE uuid = '1af6effd-ff62-4130-8058-ebe5f218bb27'
ORDER BY ip;
	`)
	if err != nil {
		log.Fatalf("database query error: %v", err)
	}
	var i string
	for ipScan.Next() {
		ipScan.Scan(&i)
		ips = append(ips, i)

	}
	fmt.Println(ips)
	td.Rows = make(map[string][]RowData)
	for _, ii := range ips {
		td.Rows[ii] = []RowData{}
	}

	rowsScan, err := database.Query(`
select *
FROM scandata
WHERE scandata.uuid = '1af6effd-ff62-4130-8058-ebe5f218bb27'
	`)
	if err != nil {
		log.Fatalf("database query error: %v", err)
	}
	var uuidScan, scanargs, scanstart, scantype, scanprotocol, scanservices, scanend, summary string
	for rowsScan.Next() {
		rowsScan.Scan(&uuidScan, &scanargs, &scanstart, &scantype, &scanprotocol, &scanservices, &scanend, &summary)
		s := ScanData{uuidScan, scanargs, scanstart, scantype, scanprotocol, scanservices, scanend, summary}
		td.Scan = s
		fmt.Println(s)

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
		where hostdata.uuid = '1af6effd-ff62-4130-8058-ebe5f218bb27'
		AND hostdata.ip = portdata.ip
		AND hostdata.ip = '%s'
		ORDER BY ip ASC;`, ipr))

		if err != nil {
			log.Fatalf("database query error: %v", err)
		}

		for rowsHost.Next() {
			rowsHost.Scan(&uuidHost, &ip, &hostState, &hReason, &hostname, &port, &pState, &pReason, &service, &method)
			temp = append(temp, RowData{uuidHost, ip, hostState, hReason, hostname, port, pState, pReason, service, method})

			fmt.Println(uuidHost, ip, hostState, hReason, hostname, port, pState, pReason, service, method)

		}
		td.Rows[ipr] = temp
		temp = nil
	}

	fmt.Println(td)

	database.Close()
}
