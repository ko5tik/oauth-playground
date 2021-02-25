package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// simple client requesting token with client credentials grant

const AuthUrl = "http://localhost:3846/client"
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
		"client_id":         {"my-client"},
		"client_credential": {"foobar"},
		"grant_type":        {"client_credentials"},
		"scope":             {"photos openid offline"},
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
	request.SetBasicAuth("my-client", "foobar")

	rd, _ := httputil.DumpRequest(request, true)
	fmt.Println("-----------   request ------------")
	fmt.Printf("%q\n", rd)
	// and call
	res, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("----------- response ------------")
	fmt.Printf("Status:\t%s\n", res.Status)

	dump, _ := httputil.DumpResponse(res, true)

	fmt.Printf("%q", dump)
	fmt.Println("----------- token ------------")

	var token = Token{}

	json.NewDecoder(res.Body).Decode(&token)

	fmt.Printf("token: %s\n", token.AccessToken)
	fmt.Printf("expires: %d\n", token.ExpiresIn)
	fmt.Printf("scope: %s\n", token.Scope)
	fmt.Printf("type: %s\n", token.TokenType)

	//  and now test request for introspection
	//  and also perform introspect in this token,  as if we were a genuine client
	introspectValues := url.Values{
		"token": {token.AccessToken},
	}

	fmt.Printf("-----------  introspect real token ------------\n")
	checkToken(introspectValues, client)

	bogusToken := url.Values{
		"token": {token.AccessToken},
	}
	fmt.Printf("-----------  introspect bogus token ------------\n")
	checkToken(bogusToken, client)
}

func checkToken(introspectValues url.Values, client *http.Client) {
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
