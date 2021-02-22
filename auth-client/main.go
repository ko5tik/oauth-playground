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

	request, _ := http.NewRequest(http.MethodPost, "http://localhost:3846/oauth2/auth", strings.NewReader(form.Encode()))
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
	fmt.Printf("Location:\t%s\n", location)
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
