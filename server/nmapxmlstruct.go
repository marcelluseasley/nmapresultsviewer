package main

import (
	"encoding/xml"
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
