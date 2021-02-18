package server

import "net/http"

func RegisterHandlers() {
	// handle authentication and return  authorisation token
	http.HandleFunc("/oauth2/auth", authEndpoint)
	// exchange authorisation token  with real token
	http.HandleFunc("/oauth2/token", tokenEndpoint)
	// introspect is necessary  for authorising rel requests
	http.HandleFunc("/oauth2/introspect", introspectionEndpoint)
}

func introspectionEndpoint(writer http.ResponseWriter, request *http.Request) {

}

func revokeEndpoint(writer http.ResponseWriter, request *http.Request) {

}

func tokenEndpoint(writer http.ResponseWriter, request *http.Request) {

}

func authEndpoint(writer http.ResponseWriter, request *http.Request) {

}
