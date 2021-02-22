package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const AuthUrl = "http://localhost:3846/oauth2/auth"
const TokenUrl = "http://localhost:3846/oauth2/token"
const IntrospectUrl = "http://localhost:3846/oauth2/introspect"

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func main() {

	// doctor up request body for post request,
	// we have to place all parameters inside, especially because it is  required by standard
	form := url.Values{
		//  identifies client app
		"client_id": {"my-client"},
		//  how response shall be delivered  values:
		//  code -  explicit grant handler
		// token - implicit grant handler
		"response_type": {"code"},
		//  add some randomnes
		"state": {"some-random-state-foobar"},
		"nonce": {"some-random-nonce-foobar"},
		// credentials supplied by user
		"username": {"peter"},
		"password": {"pan"},
		//  scopes supplied by user
		"scopes": {"photos openid offline"},

		//  callback url -  proper ne has to be configured on server
		// and match the one we are sending.
		//		"redirect_uri": {"callback://callback.url"},
	}

	//  create custom https client, one that does not do automatic redirects
	// we cal also add TLS here - and we shall do it later
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	request, _ := http.NewRequest(http.MethodPost, AuthUrl, strings.NewReader(form.Encode()))
	//important, so values are encoded properly
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// and call
	res, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("----------- response ------------\n")
	fmt.Printf("Status:\t%s\n", res.Status)
	// Location is important, this URL contains "code" parameter
	// which can be exchanged for token later in separate  request to token endpoint
	location, err := res.Location()
	if err != nil {
		log.Fatal("no location was returned")
	}
	fmt.Printf("Location:\t%s\n", location)

	// now we have location,   extract code
	code := location.Query().Get("code")
	fmt.Printf("extracted code: %s\n", code)

	//  exchange code for token
	tokenValues := url.Values{
		"code": {code},
		//
		"grant_type": {"authorization_code"},
		//  identifies client app
		"client_id":     {"my-client"},
		"client_secret": {"foobar"},
	}

	fmt.Printf("----------- retrieve token --------------\n")

	tokenRequest, _ := http.NewRequest(http.MethodPost, TokenUrl, strings.NewReader(tokenValues.Encode()))
	tokenRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	tokenResponse, err := client.Do(tokenRequest)
	if err != nil {
		log.Fatal("token response failed\n", err)
	}

	// token response contained in body as JSON,   parse it
	token := Token{}

	json.NewDecoder(tokenResponse.Body).Decode(&token)

	fmt.Printf("token: %s\n", token.AccessToken)
	fmt.Printf("expires: %d\n", token.ExpiresIn)
	fmt.Printf("scope: %s\n", token.Scope)
	fmt.Printf("type: %s\n", token.TokenType)

	//  and also perform introspect in this token,  as if we were a genuine client
	introspectValues := url.Values{
		"token": {token.AccessToken},
	}
	introspectRequest, _ := http.NewRequest(http.MethodPost, IntrospectUrl, strings.NewReader(introspectValues.Encode()))
	introspectRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	introspectRequest.SetBasicAuth("my-client", "foobar")

	introspectResult, err := client.Do(introspectRequest)
	if err != nil {
		log.Fatal("introspection response failed\n", err)
	}

	fmt.Printf("-----------  introspection response ------------\n")
	fmt.Printf("Status:\t%s\n", introspectResult.Status)

	fmt.Printf("-----------  introspection body ------------\n")
	data, _ := ioutil.ReadAll(introspectResult.Body)
	fmt.Printf(string(data))

}
