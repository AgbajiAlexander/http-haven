package main

import (
	"fmt"
	"io"
	"net/http"
)

func methodInspectorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You made a %s request.", r.Method)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "body cannot be empty", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(body)
}

func main() {
	http.HandleFunc("/method-inspector", methodInspectorHandler)
	http.HandleFunc("/echo", echoHandler)
	fmt.Println("Server starting on :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}