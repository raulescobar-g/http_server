package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello_world(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HOWDY")
}
func main() {
	http.HandleFunc("/", hello_world)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
