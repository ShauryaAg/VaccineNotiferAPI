package main

import (
	"log"
	"net/http"
	"time"

	"cov-api/handlers"
	"cov-api/middlewares"
	"cov-api/models/db"
	"cov-api/utils"

	"github.com/gorilla/mux"
)

const (
	CURRENT_HOST = "localhost:8080"
)

func main() {
	utils.AddCronJobs(CURRENT_HOST)

	db.DBCon, _ = db.CreateDatabase() // initialising the database

	r := mux.NewRouter().StrictSlash(true)

	// API Routes
	r.HandleFunc("/login", handlers.Login).Methods("POST")               // POST /login
	r.HandleFunc("/register", handlers.Register).Methods("POST")         // POST /register
	r.HandleFunc("/t/{token}", handlers.VerifyToken).Methods("GET")      // GET /t/<token>; For Email verification
	r.HandleFunc("/u/{token}", handlers.UnsubscribeToken).Methods("GET") // GET /t/<token>; For Email verification

	r.HandleFunc("/", handlers.SendNotification).Methods("GET")

	// Auth routes
	r.Handle("/user", middlewares.AuthMiddleware(
		http.HandlerFunc(handlers.GetUser),
	)).Methods("GET") // GET /user Auth: Bearer <Token>
	r.Handle("/user", middlewares.AuthMiddleware(
		http.HandlerFunc(handlers.UpdateUser),
	)).Methods("PATCH") // GET /user Auth: Bearer <Token>
	r.Handle("/unsubscribe", middlewares.AuthMiddleware(
		http.HandlerFunc(handlers.UnsubscribeUser),
	)).Methods("POST") // POST /unsubscribe Auth: Bearer <Token>

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      r,
	}

	log.Fatal(srv.ListenAndServe())

}
