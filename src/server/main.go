package main

import (
	"fmt"
	"github.com/webdevwilson/serverless-multicloud/src/boundary"
	"log"
	"net/http"
	"os"
)

func main() {

	var port string
	var ok bool
	if port, ok = os.LookupEnv("PORT"); !ok {
		port = "8080"
	}

	http.Handle("/", boundary.NewRequestHandler(""))

	log.Printf("[INFO] Starting HTTP server on port %s", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}