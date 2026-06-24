package main

import (
	"fmt"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func main() {
	http.HandleFunc("/ping", pingHandler)
	fmt.Println("Server starting on :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
