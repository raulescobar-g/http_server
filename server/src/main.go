package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type userResponse struct {
	Name   string   `json:"name"`
	Colors []string `json:"colors"`
}
type errorResponse struct {
	Error string `json:"error"`
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("%s %s", r.Method, r.URL.Path)) //log path and method
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("%s %s", r.Method, r.URL.Path)) //log path and method
	conn, err := getConn()

	if errorFound(err, "connecting to database", http.StatusInternalServerError, w) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")

		personId := r.URL.Query().Get("id")

		rows, err := conn.Query("SELECT name,color FROM users JOIN users_colors ON users.id = users_colors.user_id AND id=$1;", personId)
		defer rows.Close()

		if errorFound(err, fmt.Sprintf("querying name and color with id=%s", personId), http.StatusBadRequest, w) {
			return
		}

		responseStruct := userResponse{}
		var name, color string

		for rows.Next() {
			rows.Scan(&name, &color)
			responseStruct.Name = name
			responseStruct.Colors = append(responseStruct.Colors, color)
		}

		if rows.Err() != nil {
			errorFound(rows.Err(), "reading rows", http.StatusInternalServerError, w)
			return
		}

		if responseStruct.Name == "" {
			errorFound(sql.ErrNoRows, "returning result", http.StatusNotFound, w)
			return
		}

		json.NewEncoder(w).Encode(responseStruct)

		return

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

func getConn() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("HOST"), os.Getenv("PORT_DB"), os.Getenv("USER_DB"), os.Getenv("PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("pgx", psqlInfo)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func errorFound(err error, expect string, statusCode int, w http.ResponseWriter) bool {
	if err != nil {
		w.WriteHeader(statusCode)
		log.Println(fmt.Sprintf("%s: %s", expect, err.Error()))
		responseStruct := errorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(responseStruct)
		return true
	}
	return false
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	//using default mux
	http.HandleFunc("/", handlePing)
	http.HandleFunc("/user", handleUser)

	log.Println("Running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
