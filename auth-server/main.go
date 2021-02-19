package main

import (
	"flag"
	"github.com/ko5tik/oauth-playground/server"
	"log"
	"net/http"
	"os"
)

func main() {

	port := "3846"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	flag.Parse()

	server.RegisterHandlers()

	// fire up webserver on designated port passed via command line
	log.Println("auth server is listening on port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
