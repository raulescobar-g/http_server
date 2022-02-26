package main

import (
	"fmt"
	"log"
	"net/http"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HOWDY")
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HOWDY")
}

func main() {
	http.HandleFunc("/", handlePing)
	http.HandleFunc("/user", handleUser)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
