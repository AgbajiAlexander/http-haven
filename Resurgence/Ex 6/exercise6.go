package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// w.WriteHeader() must be called BEFORE w.Write() or fmt.Fprintf().
// If you write the body first, Go automatically sends a 200 header and
// any subsequent w.WriteHeader() call is silently ignored.
// Observation: curl -v shows "< HTTP/1.1 200 OK" even when w.WriteHeader(404)
// is called after w.Write() — the 404 never makes it into the response.

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

func statusHandler(w http.ResponseWriter, r *http.Request) {
	codeStr := r.URL.Query().Get("code")
	if codeStr == "" {
		http.Error(w, "code parameter is required", http.StatusBadRequest)
		return
	}

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		http.Error(w, "code must be a valid integer", http.StatusBadRequest)
		return
	}

	if code < 100 || code > 599 {
		http.Error(w, "code must be a valid HTTP status code (100–599)", http.StatusBadRequest)
		return
	}

	statusText := http.StatusText(code)
	if statusText == "" {
		statusText = "Unknown"
	}

	w.WriteHeader(code)
	fmt.Fprintf(w, "Responding with status %d %s", code, statusText)
}

func apiPingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func apiGreetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Stranger"
	}
	fmt.Fprintf(w, "Greetings, %s!", name)
}

func main() {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/v1/ping", apiPingHandler)
	apiMux.HandleFunc("/v1/greet", apiGreetHandler)

	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/method-inspector", methodInspectorHandler)
	mainMux.HandleFunc("/echo", echoHandler)
	mainMux.HandleFunc("/headers", headersHandler)
	mainMux.HandleFunc("/form", formHandler)
	mainMux.HandleFunc("/status", statusHandler)
	mainMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	fmt.Println("Server starting on :8080")
	fmt.Println(http.ListenAndServe(":8080", mainMux))
}