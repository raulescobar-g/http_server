package main

import (
	"database/sql"
	"encoding/json"
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
	password = "Chispis1#"
	dbname   = "users"
)

type userResponse struct {
	Name   string   `json:"name"`
	Colors []string `json:"colors"`
}

func handlePing(w http.ResponseWriter, r *http.Request) {}

func handleUser(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("%s %q", r.Method, html.EscapeString(r.URL.Path))) //basic logging
	conn := getConn()

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")

		personId := r.URL.Query().Get("id")
		rows, err := conn.Query("SELECT name,color FROM users JOIN users_colors ON users.id = users_colors.user_id AND id=$1;", personId)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Error querying user colors with error: %s", err)) // here return error response
		}

		responseStruct := userResponse{}
		var name, color string
		for rows.Next() {
			err = rows.Scan(&name, &color)
			if err != nil {
				log.Fatalf(fmt.Sprintf("Could not read user with error: %s", err))
			}

			responseStruct.Name = name
			responseStruct.Colors = append(responseStruct.Colors, color)

		}
		json.NewEncoder(w).Encode(responseStruct)

	case http.MethodPost:
		// Create a new record.
	case http.MethodPut:
		// Update an existing record.
	case http.MethodDelete:
		// Remove the record.
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

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
	//log.Println("Successfully connected to database") debug statement

	return db
}

func main() {
	port := ":8080"

	//using default mux
	http.HandleFunc("/", handlePing)
	http.HandleFunc("/user", handleUser)

	log.Println("Running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
