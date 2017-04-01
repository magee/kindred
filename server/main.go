package main

import (
	"net/http"
	"github.com/codegangsta/negroni"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	// "github.com/jinzhu/gorm"
// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type user struct {
	Username string
	Name string
	Email string
	Password string
}

type userBcrypt struct {
	Username string
	Name string
	Email string
	Password []byte
}

func main() {
	r := http.NewServeMux()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r.Handle("/api/protected", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(protected))))

	http.Handle("/", http.FileServer(http.Dir("../public/")))
	http.Handle("/bundles/", http.StripPrefix("/bundles/", http.FileServer(http.Dir("../bundles/"))))
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/signup", signup)
	http.Handle("/api/protected", r)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}


//redirects
//http.Redirect(w, req, "/", http.MovedPermanently, http.StatusSeeOther, or http.StatusTemporaryRedirect)
//lower level
//w.Header().Set("Location", "/")
//w.WriteHeader(http.StatusSeeOther)

// 1) On client, ask for userName and password
// 2) Exchange the userName and password for a time-limited access token via HTTPS. Use jwt-go on the server
//     to create the token. Use bcrypt to encrypt and compare passwords.
// 3) Add the recieved access token to the request header for any RESTful API requiring authorization
// {
// 	"iss": "http://kindrechat.io",
// 	"user": "xxxxxx"
// }
// 4) On the server, add an access token checker middleware for those routes. JWT tokens have an expire (exp)
//     and not before (nbf) timestamp. JWT validates those when it parses the token from the header.
// 5) On client, periodically refresh the token. Our tokens expire in 5 minutes. I refresh them every 4 minutes.