# http-haven

A collection of Go HTTP server exercises built to develop foundational understanding of the `net/http` package.

---

## Project Structure

```
http-haven/
├── exercise1.go   # Basic Ping-Pong Server
├── exercise2.go   # Query Parameters & Path Validation
├── exercise3.go   # Text Counter
├── exercise4.go   # Basic Math API
├── exercise5.go   # User-Agent Echo
├── exercise6.go   # Secure Dashboard
├── exercise7.go   # Simple Redirector
└── README.md
```

---

## Running Any Exercise

```bash
go run exerciseN.go
```

The server starts on `http://localhost:8080`.

---

## Exercise 1: Basic Ping-Pong Server

**File:** `exercise1.go`

**Goal:** Build a minimal web server that listens on port 8080 and responds with `pong` when a user visits `/ping`.

**Concepts learned:**
- Creating a route handler with `http.HandleFunc`
- Writing a plain text response with `fmt.Fprint(w, ...)`
- Starting a server with `http.ListenAndServe`

**Test:**
```bash
curl http://localhost:8080/ping
# pong
```

---

## Exercise 2: Query Parameters & Path Validation

**File:** `exercise2.go`

**Goal:** Create a `/hello` endpoint that reads a `name` query parameter and responds with `Hello, <name>!`. Defaults to `Hello, Guest!` if the parameter is missing. Rejects non-GET requests with a 405.

**Concepts learned:**
- Extracting query parameters with `r.URL.Query().Get("name")`
- Checking the HTTP method with `r.Method`
- Rejecting invalid methods with `http.StatusMethodNotAllowed` (405)

**Test:**
```bash
curl "http://localhost:8080/hello?name=Alice"  # Hello, Alice!
curl "http://localhost:8080/hello"              # Hello, Guest!
curl -X POST "http://localhost:8080/hello"      # 405 Method Not Allowed
```

---

## Exercise 3: Text Counter

**File:** `exercise3.go`

**Goal:** Build a `/count` route. GET requests return an instruction message. POST requests read the body and return the number of characters.

**Concepts learned:**
- Differentiating between GET and POST using `r.Method`
- Reading the full request body with `io.ReadAll(r.Body)`
- Returning the character count with `len(body)`

**Test:**
```bash
curl http://localhost:8080/count                       # GET: instruction message
curl -X POST -d "Golang" http://localhost:8080/count   # POST: 6
curl -X POST -d "Hello World" http://localhost:8080/count  # POST: 11
```

---

## Exercise 4: Basic Math API

**File:** `exercise4.go`

**Goal:** Create a `/calculate` route that accepts `op`, `a`, and `b` as query parameters and returns the result. Supports `add`, `subtract`, and `multiply`. Returns 400 for invalid input or unknown operations.

**Concepts learned:**
- Parsing multiple query parameters
- Converting strings to integers with `strconv.Atoi()`
- Using a `switch` statement to handle operations
- Returning `http.StatusBadRequest` (400) for bad input

**Test:**
```bash
curl "http://localhost:8080/calculate?op=add&a=12&b=8"       # Result: 20
curl "http://localhost:8080/calculate?op=subtract&a=20&b=5"  # Result: 15
curl "http://localhost:8080/calculate?op=multiply&a=4&b=3"   # Result: 12
curl "http://localhost:8080/calculate?op=multiply&a=abc&b=5" # 400 Bad Request
```

---

## Exercise 5: User-Agent Echo

**File:** `exercise5.go`

**Goal:** Create an `/agent` route that reads the `User-Agent` header and echoes it back. Defaults to `Unknown` if the header is missing.

**Concepts learned:**
- Reading request headers with `r.Header.Get("User-Agent")`
- Handling missing or empty headers with a fallback value

**Test:**
```bash
curl -H "User-Agent: CustomTester/1.0" http://localhost:8080/agent
# You are visiting us using: CustomTester/1.0

curl -H "User-Agent: " http://localhost:8080/agent
# You are visiting us using: Unknown
```

---

## Exercise 6: Secure Dashboard

**File:** `exercise6.go`

**Goal:** Create a `/dashboard` route protected by an API key. Requests without the correct `X-API-Key` header are rejected with a 401.

**Concepts learned:**
- Reading custom headers with `r.Header.Get("X-API-Key")`
- Matching against a hardcoded secret value
- Returning `http.StatusUnauthorized` (401) for unauthorized requests

**Test:**
```bash
curl http://localhost:8080/dashboard
# Unauthorized (401)

curl -H "X-API-Key: wrongkey" http://localhost:8080/dashboard
# Unauthorized (401)

curl -H "X-API-Key: secret123" http://localhost:8080/dashboard
# Welcome to the secure dashboard!
```

---

## Exercise 7: Simple Redirector

**File:** `exercise7.go`

**Goal:** Create a `/legacy` route that permanently redirects clients to `/v2`, which returns a friendly welcome message.

**Concepts learned:**
- Redirecting with `http.Redirect(w, r, "/v2", http.StatusMovedPermanently)`
- Using status code 301 for a permanent redirect
- Registering multiple related routes in `main()`

**Test:**
```bash
curl -o /dev/null -w "%{http_code}" http://localhost:8080/legacy
# 301

curl -L http://localhost:8080/legacy
# Welcome to version 2
```

---

## Key Concepts Summary

| Concept | Tool |
|---|---|
| Route registration | `http.HandleFunc(pattern, handler)` |
| Write a response | `fmt.Fprint(w, ...)` / `fmt.Fprintf(w, ...)` |
| Read query params | `r.URL.Query().Get("key")` |
| Check HTTP method | `r.Method` |
| Read request body | `io.ReadAll(r.Body)` |
| Convert string to int | `strconv.Atoi()` |
| Read headers | `r.Header.Get("Header-Name")` |
| Return error + status | `http.Error(w, message, statusCode)` |
| Redirect | `http.Redirect(w, r, url, statusCode)` |
| Start server | `http.ListenAndServe(":8080", nil)` |
