package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type app struct {
	auth struct {
		username string
		password string
	}
}

func main() {
	app := new(app)

	app.auth.username = os.Getenv("AUTH_USERNAME")
	app.auth.password = os.Getenv("AUTH_PASSWORD")

	validateAuth(app.auth.username, app.auth.password)

	fmt.Printf("hello %v", app.auth.username)

	mux := http.NewServeMux()
	mux.HandleFunc("/public", app.publicHandler)
	mux.HandleFunc("/private", app.basicAuth(app.privateHandler))

	srv := &http.Server{
		Addr:         "localhost:4000",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting server on %s", srv.Addr)
	err := srv.ListenAndServeTLS("./localhost.pem", "./localhost-key.pem")
	log.Fatal(err)
}

func validateAuth(username string, password string) {
	if len(username) == 0 {
		log.Fatal("username must be provided")
	}

	if len(password) == 0 {
		log.Fatal("password must be provided")
	}
}

func (app *app) publicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the public handler")
}

func (app *app) privateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the private handler")
}

func (app *app) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(app.auth.username))
			expectedPasswordHash := sha256.Sum256([]byte(app.auth.password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
