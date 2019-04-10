// server

package main

import (
	
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"github.com/mssola/user_agent"
)



func main(){
	r := mux.NewRouter()
	r.HandleFunc("/", userAgentHandler).Methods("GET")
	http.Handle("/", r)
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func userAgentHandler(w http.ResponseWriter, r *http.Request) {
	ua := r.UserAgent()
	uAgent := user_agent.New(ua)
	name, _ := uAgent.Browser()
	log.Printf("User Agent: %s\n", ua)
	log.Printf("%v\n",name)

}