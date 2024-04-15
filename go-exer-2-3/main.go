package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

type app struct {
	auth struct {
		username string
		password string
	}
}

type PasswordType int

const (
	Random PasswordType = iota
	AlphaNumeric
	Pin
)

//export POSTGRESQL_URL='postgres://testuser:testpassword@localhost:5432/passlocker?sslmode=disable'
//PSWLCKRDSN=postgres://testuser:testpassword@localhost/passlocker AUTH_USERNAME=hans AUTH_PASSWORD=password go run .

func main() {
	app := new(app)

	fmt.Printf("hello %v", app.auth.username)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /public", app.publicHandler)
	mux.HandleFunc("GET /private", app.basicAuth(app.privateHandler))

	mux.HandleFunc("GET /generate-password", app.basicAuth(app.generatePasswordHandler))
	// curl trigger : curl -k -u hans:password https://localhost:4000/generate-password

	mux.HandleFunc("POST /save-credentials", app.basicAuth(app.saveCredentialHandle))
	// curl trigger :  curl -k -u hans:password -d
	//'{"url":"www.painhub.com", "username":"hans", "password":"BJ$hjeAI1o"}'  -X POST https://localhost:4000/save-credentials

	mux.HandleFunc("GET /list-credentials", app.basicAuth(app.listCredentialHandle))
	// curl trigger :  curl -k -u hans:password https://localhost:4000/list-credentials

	srv := &http.Server{
		Addr:         "localhost:4000",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting server on %s", srv.Addr)
	error := srv.ListenAndServeTLS("./localhost.pem", "./localhost-key.pem")
	log.Fatal(error)
}

func (app *app) publicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the public handler")
}

func (app *app) privateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the private handler")
}

func (app *app) generatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var password string = generatePassword(10, false, false, false, 0)

	fmt.Fprintf(w, "this is the generated password %v", password)
}

func (app *app) saveCredentialHandle(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON data from the request body
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := data["url"].(string)
	username := data["username"].(string)
	password := data["password"].(string)

	conn, err := pgx.Connect(context.Background(), os.Getenv("PSWLCKRDSN"))
	if err != nil {
		log.Fatal(err)
	}

	Check(err)
	defer conn.Close(context.Background())

	fmt.Fprintf(w, "Successfully saved credentials with the following details: \n URL : %s  \n Credentials %s:%s", url, username, password)
}

func (app *app) listCredentialHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "list all credentials")

	username, password, ok := r.BasicAuth()

	if !ok && len(password) == 0 {
		log.Printf("error")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(context.Background(), os.Getenv("PSWLCKRDSN"))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close(context.Background())

	users, err := dbGetUserByUsername(ctx, conn, username)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, "before credentials db call")
	credentials, err := dbAllCredentialsForUser(ctx, conn, users[0].Id)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, "list all credentials \n", credentials)
}

func generatePassword(passwordLength int, includeNumbersFlag bool, includeSymbolsFlag bool, includeUppercaseFlag bool, passwordType int) string {
	var chars string = "abcdefghijklmnopqrstuvwxyz"

	if includeNumbersFlag || passwordType == int(Random) {
		chars += "0123456789"
	}
	if includeSymbolsFlag || passwordType == int(Random) {
		chars += "!@#$%^&*()_+{}[]:;<>,.?/~`"
	}
	if includeUppercaseFlag || passwordType == int(Random) {
		chars += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	if passwordType == int(Pin) {
		chars = "0123456789"
	}

	password := make([]byte, passwordLength)
	for i := 0; i < passwordLength; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		password[i] = chars[randomIndex.Int64()]
	}

	return string(password)
}

/*
basicAuth is a middleware function that performs basic authentication.

It takes a next http.HandlerFunc as a parameter and returns a new http.HandlerFunc.

The middleware checks if the request contains valid basic authentication credentials.

If the credentials are valid, it calls the next handler function.

If the credentials are not valid, it sets the WWW-Authenticate header and returns an Unauthorized error.

@param next http.HandlerFunc - The next handler function to be called if authentication is successful.

@return http.HandlerFunc - The middleware function that performs basic authentication.

TODO : move this on a separate module
*/
func (app *app) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {

			ctx := context.Background()
			conn, err := pgx.Connect(context.Background(), os.Getenv("PSWLCKRDSN"))
			if err != nil {
				log.Fatal(err)
			}

			defer conn.Close(context.Background())

			users, err := dbGetUserByUsername(ctx, conn, username)
			if err != nil {
				log.Fatal(err)
			}

			var user = users[len(users)-1]

			passwordHash := sha256.Sum256([]byte(user.Password))
			expectedPasswordHash := sha256.Sum256([]byte(password))

			usernameMatch := username == user.Username
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if passwordMatch && usernameMatch {
				next.ServeHTTP(w, r)
				return
			} else {
				app.auth.username = username
				app.auth.password = password
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
