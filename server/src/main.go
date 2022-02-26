package main

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = 
	dbname   = "users"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Serve the resource.
	case http.MethodPost:
		// Create a new record.
	case http.MethodPut:
		// Update an existing record.
	case http.MethodDelete:
		// Remove the record.
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	fmt.Fprintf(w, "%s %q", r.Method, html.EscapeString(r.URL.Path))
	fmt.Fprintf(w, "HOWDY")
}

func getConn() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}
	log.Println("Connected to database")

	return db
}

func main() {
	port := ":8080"

	con := getConn()

	rows, err := con.Query("SELECT * FROM USERS")
	if err != nil {
		log.Fatalf("Could not fetch users")
	}
	var user_id string
	var name string
	for rows.Next() {
		err = rows.Scan(&user_id, &name)
		if err != nil {
			log.Fatalf("Could not read user with error:", err)
		}
		log.Println(user_id, name)
	}

	//using default mux
	http.HandleFunc("/", handlePing)
	http.HandleFunc("/user", handleUser)
	log.Println("Running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
