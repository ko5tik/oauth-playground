package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/storage"
	"github.com/ory/fosite/token/jwt"
	"io/ioutil"
	"log"
	"time"
)

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
	privateKey, _ = readKey() //rsa.GenerateKey(rand.Reader, 2048)
)

// Build a fosite instance with all OAuth2 and OpenID Connect handlers enabled, plugging in our configurations as specified above.
var fosite = compose.Compose(
	config,
	store,
	&oauth2.DefaultJWTStrategy{
		JWTStrategy: &jwt.RS256JWTStrategy{
			PrivateKey: privateKey,
		},
		HMACSHAStrategy: nil,
		Issuer:          "",
		ScopeField:      0,
	},

	nil,

	compose.OAuth2ClientCredentialsGrantFactory,
)

//read key from well known location
func readKey() (*rsa.PrivateKey, error) {
	priv, err := ioutil.ReadFile("server.pem")
	if err != nil {
		log.Fatal("unable to load private key file")
	}

	privPem, _ := pem.Decode(priv)
	var privPemBytes []byte
	if privPem.Type != "RSA PRIVATE KEY" {
		log.Fatal("not RSA key file")
	}

	privPemBytes, err = x509.DecryptPEMBlock(privPem, []byte("foobar"))
	if err != nil {
		log.Fatal("unable to decrypt privatet ekey file")
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			log.Fatal("unable to parse private key")
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		log.Fatal("cast fialed")
	}

	return privateKey, nil
}
