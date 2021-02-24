package main

import (
	"client-credentials-server/server"
	"log"
	"net/http"
	"os"
)

//  simple client credentials oauth server using fosite

func main() {

	port := "3846"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	server.Register()

	// fire up webserver on designated port passed via command line
	log.Println("auth server is listening on port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
