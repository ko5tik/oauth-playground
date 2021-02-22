package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

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

		//  callback url
		"redirect-url": {"callback:"},
	}

	// request goes agains auth server
	res, err := http.Post("http://localhost:3846/oauth2/auth", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	if err != nil {
		log.Fatal(err)
	}

	header := res.Header
	fmt.Printf("----------- header ------------\n")
	for key, value := range header {
		buf := bytes.Buffer{}

		for _, v := range value {
			buf.WriteString(v)
		}
		fmt.Printf("%s : %s\n", key, buf.String())
	}

	fmt.Printf("----------- body ------------\n")
	data, _ := ioutil.ReadAll(res.Body)
	fmt.Printf(string(data))

}
