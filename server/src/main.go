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

// body content
type userResponse struct {
	Name   string   `json:"name"`
	Colors []string `json:"colors"`
}

//body content on errors
type errorResponse struct {
	Error string `json:"error"`
}

// ping handler
func handlePing(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
}

// handler for /user path
func handleUser(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
	conn, err := getConn()

	if errorFound(err, "connecting to database", http.StatusInternalServerError, w) {
		return
	}

	switch r.Method {
	//  path: /user?user_id=<int>  ==> body on response {name:<string>, colors:[<string>,...]}
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")

		userId := r.URL.Query().Get("user_id")

		if userId == "" {
			errorFound(sql.ErrNoRows, "getting input parameter", http.StatusBadRequest, w)
			return
		}

		rows, err := conn.Query("SELECT name,color FROM users JOIN users_colors ON users.id = users_colors.user_id AND id=$1;", userId)

		if errorFound(err, fmt.Sprintf("querying name and color with id=%s", userId), http.StatusBadRequest, w) {
			return
		}

		// this must go after error checking the rows, or else will close a nil pointer
		defer rows.Close()

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

		// on nothing found return no rows sql error
		if responseStruct.Name == "" {
			errorFound(sql.ErrNoRows, "returning result", http.StatusNotFound, w)
			return
		}
		// write to response json
		json.NewEncoder(w).Encode(responseStruct)

		return

	case http.MethodPost:
		// Create a new record. not implemented
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	case http.MethodPut:
		// Update an existing record. not implemented
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	case http.MethodDelete:
		// Remove the record. not implemented
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// gets db connection, tests the connection and returns Result pair
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

// on err, will write status code into response, log the expected behavior with error, and write the error to the body
// returns true if err not nill
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

	// if cant load db credentials, log it, and exit
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	//using default mux for paths
	http.HandleFunc("/", handlePing)
	http.HandleFunc("/user", handleUser)

	// log on server listening and server exit
	log.Println("Running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
