package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	fmt.Fprintf(w, "Hello, %s!", name)
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprint(w, "Send a POST request with text to count words")
		return
	}

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%d", len(body))
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	op := r.URL.Query().Get("op")
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")

	a, err := strconv.Atoi(aStr)
	if err != nil {
		http.Error(w, "Invalid value for 'a'", http.StatusBadRequest)
		return
	}

	b, err := strconv.Atoi(bStr)
	if err != nil {
		http.Error(w, "Invalid value for 'b'", http.StatusBadRequest)
		return
	}

	var result int
	switch op {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	default:
		http.Error(w, "Unknown operation", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Result: %d", result)
}

func main() {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/count", countHandler)
	http.HandleFunc("/calculate", calculateHandler)
	fmt.Println("Server starting on :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
