package server

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	"log"
	"net/http"
	"time"

	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
)

func RegisterHandlers() {
	// handle authentication and return  authorisation token
	http.HandleFunc("/oauth2/auth", authEndpoint)
	// exchange authorisation token  with real token
	http.HandleFunc("/oauth2/token", tokenEndpoint)
	// introspect is necessary  for authorising rel requests
	http.HandleFunc("/oauth2/introspect", introspectionEndpoint)
}

// here we just took everything out of  fosite example.  we will need somethign more substantial for prodiction

// fosite requires four parameters for the server to get up and running:
// 1. config - for any enforcement you may desire, you can do this using `compose.Config`. You like PKCE, enforce it!
// 2. store - no auth service is generally useful unless it can remember clients and users.
//    fosite is incredibly composable, and the store parameter enables you to build and BYODb (Bring Your Own Database)
// 3. secret - required for code, access and refresh token generation.
// 4. privateKey - required for id/jwt token generation.
var (
	// Check the api documentation of `compose.Config` for further configuration options.
	config = &compose.Config{
		AccessTokenLifespan: time.Minute * 30,
		// ...
	}

	// This is the example storage that contains:
	// * an OAuth2 Client with id "my-client" and secret "foobar" capable of all oauth2 and open id connect grant and response types.
	// * a User for the resource owner password credentials grant type with username "peter" and password "secret".
	//
	// You will most likely replace this with your own logic once you set up a real world application.
	store = storage.NewExampleStore()

	// This secret is used to sign authorize codes, access and refresh tokens.
	// It has to be 32-bytes long for HMAC signing. This requirement can be configured via `compose.Config` above.
	// In order to generate secure keys, the best thing to do is use crypto/rand:
	//
	// ```
	// package main
	//
	// import (
	//	"crypto/rand"
	//	"encoding/hex"
	//	"fmt"
	// )
	//
	// func main() {
	//	var secret = make([]byte, 32)
	//	_, err := rand.Read(secret)
	//	if err != nil {
	//		panic(err)
	//	}
	// }
	// ```
	//
	// If you require this to key to be stable, for example, when running multiple fosite servers, you can generate the
	// 32byte random key as above and push it out to a base64 encoded string.
	// This can then be injected and decoded as the `var secret []byte` on server start.
	secret = []byte("some-cool-secret-that-is-32bytes")

	// privateKey is used to sign JWT tokens. The default strategy uses RS256 (RSA Signature with SHA-256)
	privateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
)

// Build a fosite instance with all OAuth2 and OpenID Connect handlers enabled, plugging in our configurations as specified above.
var oauth2 = compose.ComposeAllEnabled(config, store, secret, privateKey)

// A session is passed from the `/auth` to the `/token` endpoint. You probably want to store data like: "Who made the request",
// "What organization does that person belong to" and so on.
// For our use case, the session will meet the requirements imposed by JWT access tokens, HMAC access tokens and OpenID Connect
// ID Tokens plus a custom field

// newSession is a helper function for creating a new session. This may look like a lot of code but since we are
// setting up multiple strategies it is a bit longer.
// Usually, you could do:
//
//  session = new(fosite.DefaultSession)
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

// introspection endpoint is called by  resource provider to  retrieve  information what is client allowed to do.
func introspectionEndpoint(rw http.ResponseWriter, req *http.Request) {
	log.Println("introspectionEndpoint called")
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

func revokeEndpoint(writer http.ResponseWriter, request *http.Request) {

}

func tokenEndpoint(rw http.ResponseWriter, req *http.Request) {
	log.Println("tokenEndpoint called")

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

	// The client now has a valid access token
}

// auth endpoint receives basic authentication data ( in request body, not URL! )
// and returns authorisation code if everything is OK  it is also responsible for returning
// app specific redirect url  as callback into application.   this is configured  and tweaked in example store
func authEndpoint(rw http.ResponseWriter, req *http.Request) {
	log.Println("authEndpoint called")
	// This context will be passed to all methods.
	ctx := req.Context()

	// Let's create an AuthorizeRequest object!
	// It will analyze the request and extract important information like scopes, response type and other
	// important stuff as defined by the standart
	ar, err := oauth2.NewAuthorizeRequest(ctx, req)

	if err != nil {
		log.Printf("Error occurred in NewAuthorizeRequest: %+v", err)
		oauth2.WriteAuthorizeError(rw, ar, err)
		return
	}

	//  retrieve form data
	req.ParseForm()
	// in real world here will be authentication of user  and his consent
	// for demo purposes we have just a vald and invalid user
	if req.PostForm.Get("username") != "valid" {
		http.Error(rw, "you're not allowed to do that", http.StatusForbidden)
		return
	}

	//  now user is authenticated,   we proceed to create respoinse

	// let's see what scopes the user gave consent to, there could be some logic
	// to ensure that only allowed scopes are accepted
	for _, scope := range req.PostForm["scopes"] {
		ar.GrantScope(scope)
	}

	// Now that the user is authorized, we set up a session:
	mySessionData := newSession("peter")

	// When using the HMACSHA strategy you must use something that implements the HMACSessionContainer.
	// It brings you the power of overriding the default values.
	//
	// mySessionData.HMACSession = &strategy.HMACSession{
	//	AccessTokenExpiry: time.Now().Add(time.Day),
	//	AuthorizeCodeExpiry: time.Now().Add(time.Day),
	// }
	//

	// If you're using the JWT strategy, there's currently no distinction between access token and authorize code claims.
	// Therefore, you both access token and authorize code will have the same "exp" claim. If this is something you
	// need let us know on github.
	//
	// mySessionData.JWTClaims.ExpiresAt = time.Now().Add(time.Day)

	// It's also wise to check the requested scopes, e.g.:
	// if ar.GetRequestedScopes().Has("admin") {
	//     http.Error(rw, "you're not allowed to do that", http.StatusForbidden)
	//     return
	// }

	// Now we need to get a response. This is the place where the AuthorizeEndpointHandlers kick in and start processing the request.
	// NewAuthorizeResponse is capable of running multiple response type handlers which in turn enables this library
	// to support open id connect.
	response, err := oauth2.NewAuthorizeResponse(ctx, ar, mySessionData)

	if err != nil {
		log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
		oauth2.WriteAuthorizeError(rw, ar, err)
		return
	}

	// Last but not least, send the response!
	oauth2.WriteAuthorizeResponse(rw, ar, response)

	// response with authorisation token sent, this token will be exchanged for access tokel later
}
