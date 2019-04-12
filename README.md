# Bishop Fox Coding Challenge

## REST API Nmap scan result

* Go runtime used: `go1.12.3 linux/amd64`
* Operating System: `Linux parrot 4.19.0-parrot1-20t-amd64 #1 SMP Parrot 4.19.20-2parrot1.20t (2019-03-09) x86_64 GNU/Linux`

### File Structure

All the files needed to build this API server are located in the `server` folder.

* The `nmap.results.xml` file is at the root of the file and isn't necessary for the server to run.
* The `vendor` folder is included, along with the `Gopkg.lock` and `Gopkg.toml` files. I used `dep` for dependency management.

```
.
├── nmap.results.xml    <-- copy of Nmap results>  
├── README.md    <-- this file>
└── server
    ├── databaseinit.go    <-- function to create the database if need be>
    ├── dbtest1.db    <-- database included>
    ├── Gopkg.lock    <-- dep>
    ├── Gopkg.toml    <-- dep>
    ├── indexhandler.go    <-- handler for index route>
    ├── main.go    <-- main go code, router>
    ├── nmapdata.db.sql    <-- sql used to create the tables if need be>
    ├── submitnmaphandler.go    <-- PUT handler that ingests the xml file>
    ├── templates
    │   ├── index.html    <-- page displayed if visiting the index of the server>
    │   └── nmapresults.html    <-- page that displays the Nmap results (scan-id needed)>
    ├── uuidresultshandler.go    <-- handler that takes the scan-id and queries db>
    └── vendor/    <-- dep vendor files (external packages)>
```

## Running the server
To start the server, you can navigate to the `server/` folder and issue:
```
$ go build
```
This will create an executable called `server`, which is hard-coded to run on port 8080. Of course this can be changed.

You can start the server by typing `./server` at the command prompt:
```
$ ./server
2019/04/12 14:56:24 Listening on port 8080
```
### Ingesting the XML file
In order for the `xml` results file to be accepted by the server, it needs to be send via a PUT request to the 
`/v1/nmap` endpoint. 
```
curl http://[SERVER_URL]/v1/nmap --upload-file [scanresults.xml]
```
The response will be a JSON response:
```
{
   "scan-id": "589c048e-cca7-4686-9d70-5a4bd07669f5"
}
```
You can view the results of the scan by taking this `scan-id` (really just a UUID assigned to each different scan import) and vising the following URL:
```
http://[SERVER_URL]/589c048e-cca7-4686-9d70-5a4bd07669f5
```
The `scan-id` will be different, depending on the JSON response you get back from the server.

## Design Decisions

Initially, I was not sure if I wanted to ingest the `xml`, `nmap`, or `gnmap` version of the scan results file. After reading  `https://nmap.org/book/man-output.html`, I decided that XML would be the best option, since it can be easily parsed and ingested into a database.

The challenge requirements were very broad and I could have gone about the design and implementation a number of ways. 