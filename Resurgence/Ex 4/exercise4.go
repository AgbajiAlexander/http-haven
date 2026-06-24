package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// r.Header.Get() returns an empty string "" for a header key that was never sent.
// r.Header.Get() is case-insensitive — "x-custom-token" and "X-Custom-Token" return the same value.

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

func headersHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Custom-Token")
	if token == "" {
		http.Error(w, "X-Custom-Token header is missing", http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "Content-Type not provided"
	} else {
		contentType = "Content-Type: " + contentType
	}

	fmt.Fprintf(w, "Token received: %s\n%s", token, contentType)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/x-www-form-urlencoded") {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	language := r.FormValue("language")

	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	if language == "" {
		http.Error(w, "language is required", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Hello %s, you are coding in %s!", username, language)
}

func main() {
	http.HandleFunc("/method-inspector", methodInspectorHandler)
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/headers", headersHandler)
	http.HandleFunc("/form", formHandler)
	fmt.Println("Server starting on :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}