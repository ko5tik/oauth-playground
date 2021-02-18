package main

import (
	"flag"
	"github.com/ko5tik/oauth-playground/server"
	"log"
	"net/http"
	"strconv"
)

func main() {

	portPtr := flag.Int("port", 3846, "listener port for server,  default 3846")

	flag.Parse()

	server.RegisterHandlers()

	// fire up webserver on designated port passed via command line
	log.Println("auth server is listening on port ", *portPtr)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil))
}
