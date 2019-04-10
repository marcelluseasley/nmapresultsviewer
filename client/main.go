// client

package main

import (
	"log"
	//"net/http"
	"github.com/google/uuid"
)

func main() {
	//resp, err := http.Get("http://localhost:8080")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// resp.Body.Close()
	var uid uuid.UUID
	uid = uuid.New()

	log.Print(uid)
}
