package main

import (
	"fmt"
	"net/http"
)

func methodInspectorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You made a %s request.", r.Method)
}

func main() {
	http.HandleFunc("/method-inspector", methodInspectorHandler)
	fmt.Println("Server starting on :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}