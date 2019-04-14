package nmapserver

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// Nmaprun - represents the results of an nmap run's XML output
type Nmaprun struct {
	XMLName xml.Name `xml:"nmaprun"`

	Scanner          string `xml:"scanner,attr"`
	Args             string `xml:"args,attr"`
	Start            string `xml:"start,attr"`
	Startstr         string `xml:"startstr,attr"`
	Version          string `xml:"version,attr"`
	Xmloutputversion string `xml:"xmloutputversion,attr"`
	Scaninfo         struct {
		Type        string `xml:"type,attr"`
		Protocol    string `xml:"protocol,attr"`
		Numservices string `xml:"numservices,attr"`
		Services    string `xml:"services,attr"`
	} `xml:"scaninfo"`
	Verbose struct {
		Level string `xml:"level,attr"`
	} `xml:"verbose"`
	Debugging struct {
		Level string `xml:"level,attr"`
	} `xml:"debugging"`
	Taskbegin []struct {
		Task string `xml:"task,attr"`
		Time string `xml:"time,attr"`
	} `xml:"taskbegin"`
	Taskend []struct {
		Task      string `xml:"task,attr"`
		Time      string `xml:"time,attr"`
		Extrainfo string `xml:"extrainfo,attr"`
	} `xml:"taskend"`
	Host []struct {
		Starttime string `xml:"starttime,attr"`
		Endtime   string `xml:"endtime,attr"`
		Status    struct {
			State     string `xml:"state,attr"`
			Reason    string `xml:"reason,attr"`
			ReasonTTL string `xml:"reason_ttl,attr"`
		} `xml:"status"`
		Address struct {
			Addr     string `xml:"addr,attr"`
			Addrtype string `xml:"addrtype,attr"`
		} `xml:"address"`
		Hostnames struct {
			Hostname struct {
				Name string `xml:"name,attr"`
				Type string `xml:"type,attr"`
			} `xml:"hostname"`
		} `xml:"hostnames"`
		Ports struct {
			Port []struct {
				Protocol string `xml:"protocol,attr"`
				Portid   string `xml:"portid,attr"`
				State    struct {
					State     string `xml:"state,attr"`
					Reason    string `xml:"reason,attr"`
					ReasonTTL string `xml:"reason_ttl,attr"`
				} `xml:"state"`
				Service struct {
					Name   string `xml:"name,attr"`
					Method string `xml:"method,attr"`
					Conf   string `xml:"conf,attr"`
				} `xml:"service"`
			} `xml:"port"`
		} `xml:"ports"`
		Times struct {
			Srtt   string `xml:"srtt,attr"`
			Rttvar string `xml:"rttvar,attr"`
			To     string `xml:"to,attr"`
		} `xml:"times"`
	} `xml:"host"`
	Runstats struct {
		Finished struct {
			Time    string `xml:"time,attr"`
			Timestr string `xml:"timestr,attr"`
			Elapsed string `xml:"elapsed,attr"`
			Summary string `xml:"summary,attr"`
			Exit    string `xml:"exit,attr"`
		} `xml:"finished"`
		Hosts struct {
			Up    string `xml:"up,attr"`
			Down  string `xml:"down,attr"`
			Total string `xml:"total,attr"`
		} `xml:"hosts"`
	} `xml:"runstats"`
}

// curl http://myservice --upload-file file.txt
func SubmitNMAPHandler(w http.ResponseWriter, r *http.Request) {
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
		`, scanID.String())
		fmt.Fprintf(w, uuidJSONResponse)

	} else if r.Method == "GET" {
		t, err := template.ParseFiles("server/templates/index.html")
		if err != nil {
			log.Printf("Error opening index template: %v", err)
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Printf("Error executing index template: %v", err)
		}
	}

}
