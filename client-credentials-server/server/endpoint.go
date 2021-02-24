package server

import (
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	"log"
	"net/http"
	"time"
)

func Register() {
	// handle authentication and return   token
	http.HandleFunc("/client", authEndpoint)
	// perform token introspection for clients
	http.HandleFunc("/oauth2/introspect", introspectionEndpoint)
}

func newSession(user string) *openid.DefaultSession {
	return &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Issuer:      "https://fosite.my-application.com",
			Subject:     user,
			Audience:    []string{"https://my-client.my-application.com"},
			ExpiresAt:   time.Now().Add(time.Hour * 6),
			IssuedAt:    time.Now(),
			RequestedAt: time.Now(),
			AuthTime:    time.Now(),
		},
		Headers: &jwt.Headers{
			Extra: make(map[string]interface{}),
		},
	}
}

//  in auth server we are responsible for client authentication and creation of token
//  as described in 4.4.2   we have to accept basic authentication
// and grant_type  must be client_credentials
func authEndpoint(rw http.ResponseWriter, req *http.Request) {

	//  request has to be POST
	if req.Method != "POST" {
		http.Error(rw, "bad method, only post allowed", http.StatusBadRequest)
	}

	// has to be  authenticated,   in a real we would use soemthing more
	// secure like certificates etc.
	user, _, ok := req.BasicAuth()

	if !ok {
		http.Error(rw, "authentication required", http.StatusForbidden)
	}

	log.Println("basic authentication successful  for ", user)

	//  now we issue token and return it

	// This context will be passed to all methods.
	ctx := req.Context()

	// Create an empty session object which will be passed to the request handlers
	mySessionData := newSession("")

	// This will create an access request object and iterate through the registered TokenEndpointHandlers to validate the request.
	accessRequest, err := oauth2.NewAccessRequest(ctx, req, mySessionData)

	// Catch any errors, e.g.:
	// * unknown client
	// * invalid redirect
	// * ...
	if err != nil {
		log.Printf("Error occurred in NewAccessRequest: %+v", err)
		oauth2.WriteAccessError(rw, accessRequest, err)
		return
	}

	// If this is a client_credentials grant, grant all requested scopes
	// NewAccessRequest validated that all requested scopes the client is allowed to perform
	// based on configured scope matching strategy.
	if accessRequest.GetGrantTypes().ExactOne("client_credentials") {
		for _, scope := range accessRequest.GetRequestedScopes() {
			accessRequest.GrantScope(scope)
		}
	}

	// Next we create a response for the access request. Again, we iterate through the TokenEndpointHandlers
	// and aggregate the result in response.
	response, err := oauth2.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		log.Printf("Error occurred in NewAccessResponse: %+v", err)
		oauth2.WriteAccessError(rw, accessRequest, err)
		return
	}

	// All done, send the response.
	oauth2.WriteAccessResponse(rw, accessRequest, response)

}

// perform token introspection
func introspectionEndpoint(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	mySessionData := newSession("")
	ir, err := oauth2.NewIntrospectionRequest(ctx, req, mySessionData)
	if err != nil {
		log.Printf("Error occurred in NewIntrospectionRequest: %+v", err)
		oauth2.WriteIntrospectionError(rw, err)
		return
	}

	oauth2.WriteIntrospectionResponse(rw, ir)
}
