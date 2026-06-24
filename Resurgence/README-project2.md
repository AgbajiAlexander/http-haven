# http-fundamentals

A collection of Go HTTP server exercises focused on building deeper understanding of request handling, headers, forms, status codes, routing, and templating.

---

## Project Structure

```
http-fundamentals/
├── exercise1.go   # The Method Inspector
├── exercise2.go   # The Echo Chamber
├── exercise3.go   # Header Detective
├── exercise4.go   # Form Decoder
├── exercise5.go   # Status Code Factory
├── exercise6.go   # The API Subtree
├── exercise7.go   # Template Renderer
└── README.md
```

---

## Running Any Exercise

```bash
go run exerciseN.go
```

The server starts on `http://localhost:8080`.

---

## Exercise 1: The Method Inspector

**File:** `exercise1.go`

**Goal:** Build a `/method-inspector` endpoint that reads the HTTP method of every incoming request and echoes it back in a descriptive sentence. No method is rejected — the handler accepts everything and reports what it sees.

**Concepts learned:**
- Reading `r.Method` directly as a string value
- Using `r.Method` in a response without branching — no `if/else` or `switch` needed
- Any HTTP method works: GET, POST, DELETE, PATCH, etc.

**Test:**
```bash
curl -X GET http://localhost:8080/method-inspector
# You made a GET request.

curl -X POST http://localhost:8080/method-inspector
# You made a POST request.

curl -X DELETE http://localhost:8080/method-inspector
# You made a DELETE request.
```

---

## Exercise 2: The Echo Chamber

**File:** `exercise2.go`

**Goal:** Create an `/echo` endpoint that only accepts POST requests. Read the entire request body and send it straight back — nothing added, nothing removed.

**Concepts learned:**
- Rejecting non-POST requests with `http.StatusMethodNotAllowed` (405)
- Reading the full request body with `io.ReadAll(r.Body)`
- Using `defer r.Body.Close()` immediately after reading to free resources
- Detecting an empty body with `len(body) == 0` and returning 400
- Setting `Content-Type` response header with `w.Header().Set()` before writing
- Why `w.Header().Set()` must be called before `w.Write()` — once the body starts writing, headers are locked

**Test:**
```bash
curl -X POST -d "Hello World" http://localhost:8080/echo
# Hello World

curl -X POST http://localhost:8080/echo
# body cannot be empty (400)

curl -X GET http://localhost:8080/echo
# Method Not Allowed (405)
```

---

## Exercise 3: Header Detective

**File:** `exercise3.go`

**Goal:** Create a `/headers` endpoint that inspects two specific request headers — `X-Custom-Token` and `Content-Type` — reports what it found, and enforces a rule about one of them.

**Concepts learned:**
- Reading custom headers with `r.Header.Get("X-Custom-Token")`
- Returning 400 when a required header is missing or empty
- Reading multiple headers in the same handler
- Building a multi-line response with `\n` inside `fmt.Fprintf`
- `r.Header.Get()` is case-insensitive — `"x-custom-token"` and `"X-Custom-Token"` return the same value
- `r.Header.Get()` returns `""` for any header that was never sent

**Test:**
```bash
curl -H "X-Custom-Token: abc123" -H "Content-Type: application/json" http://localhost:8080/headers
# Token received: abc123
# Content-Type: application/json

curl -H "X-Custom-Token: abc123" http://localhost:8080/headers
# Token received: abc123
# Content-Type not provided

curl http://localhost:8080/headers
# X-Custom-Token header is missing (400)

# Case-insensitivity test
curl -H "x-custom-token: abc123" http://localhost:8080/headers
# Token received: abc123
```

---

## Exercise 4: Form Decoder

**File:** `exercise4.go`

**Goal:** Build a `/form` endpoint that accepts a POST request with a URL-encoded form body containing `username` and `language` fields. Parse, validate, and return a formatted confirmation.

**Concepts learned:**
- Calling `r.ParseForm()` explicitly before reading fields — gives control over parse errors
- Reading form fields with `r.FormValue("field")` — returns `""` if the field is missing
- Difference between `r.ParseForm()` + `r.Form.Get()` vs `r.FormValue()` — `r.FormValue()` calls `ParseForm` internally but swallows the error
- Validating each field independently with a clear error message per missing field
- Checking `Content-Type` with `strings.Contains()` instead of `==` — handles charset suffixes like `application/x-www-form-urlencoded; charset=UTF-8`
- Returning `http.StatusUnsupportedMediaType` (415) when the data format is wrong

**Test:**
```bash
curl -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=Ada&language=Go" \
  http://localhost:8080/form
# Hello Ada, you are coding in Go!

curl -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=Ada" \
  http://localhost:8080/form
# language is required (400)

curl -X POST \
  -H "Content-Type: text/plain" \
  -d "username=Ada" \
  http://localhost:8080/form
# Unsupported Media Type (415)
```

---

## Exercise 5: Status Code Factory

**File:** `exercise5.go`

**Goal:** Build a `/status` endpoint that accepts a `code` query parameter and responds using that exact HTTP status code.

**Concepts learned:**
- Setting the response status code explicitly with `w.WriteHeader(code)`
- Critical ordering rule — `w.WriteHeader()` must come before `w.Write()` or `fmt.Fprintf()`; writing the body first locks in a 200 and makes `w.WriteHeader()` a silent no-op
- Getting the official status name from a code with `http.StatusText(code)`
- Validating a range with `code < 100 || code > 599`
- Handling codes with no official name by falling back to `"Unknown"`

**Test:**
```bash
curl -v "http://localhost:8080/status?code=404"
# HTTP/1.1 404 Not Found
# Responding with status 404 Not Found

curl -v "http://localhost:8080/status?code=201"
# HTTP/1.1 201 Created
# Responding with status 201 Created

curl "http://localhost:8080/status?code=abc"
# code must be a valid integer (400)

curl "http://localhost:8080/status?code=999"
# code must be a valid HTTP status code (100-599) (400)
```

---

## Exercise 6: The API Subtree

**File:** `exercise6.go`

**Goal:** Build a mini API under the `/api/v1/` path prefix using a separate `http.ServeMux`. Mount it onto the main server at `/api/`.

**Concepts learned:**
- Creating a custom mux with `http.NewServeMux()` — independent from the global default mux
- Registering routes on a custom mux with `apiMux.HandleFunc()`
- Mounting a submux with `http.StripPrefix("/api", apiMux)` — strips the prefix before the submux sees the path
- Trailing slash pattern `/api/` matches any path that starts with `/api/`
- Passing a custom mux to `http.ListenAndServe(":8080", mainMux)` instead of `nil`
- How request routing flows: `mainMux` → `StripPrefix` → `apiMux` → handler

**Test:**
```bash
curl http://localhost:8080/api/v1/ping
# pong

curl "http://localhost:8080/api/v1/greet?name=Zion"
# Greetings, Zion!

curl http://localhost:8080/api/v1/greet
# Greetings, Stranger!
```

---

## Exercise 7: Template Renderer

**File:** `exercise7.go`

**Goal:** Build a `/render` endpoint that accepts `title`, `body`, and `style` query parameters and renders them into an inline HTML template defined as a string constant in the Go file.

**Concepts learned:**
- Defining an HTML template as a raw string constant with backticks
- Parsing a template once at package level with `template.Must(template.New("page").Parse(tmplStr))`
- `template.Must()` — panics if the template has syntax errors; correct for startup-time templates because a broken template means the program should not run at all
- Passing a struct as template data — fields must be exported (capitalised) to be accessible in the template
- `{{.FieldName}}` syntax — the dot refers to the data passed into `Execute`
- `{{if eq .Style "bold"}}<strong>{{.Body}}</strong>{{else}}{{.Body}}{{end}}` — template conditionals
- Setting `Content-Type: text/html` before calling `tmpl.Execute(w, data)` — Execute writes directly to `w`, so headers must be set first
- Handling template execution errors with a 500 Internal Server Error

**Test:**
```bash
curl "http://localhost:8080/render?title=Hello&body=World"
# renders plain HTML page

curl "http://localhost:8080/render?title=Hello&body=World&style=bold"
# renders with <strong>World</strong>

curl "http://localhost:8080/render?title=Hello"
# title and body are required (400)
```

---

## Key Concepts Summary

| Concept | Tool |
|---|---|
| Read HTTP method | `r.Method` |
| Read request body | `io.ReadAll(r.Body)` |
| Close body stream | `defer r.Body.Close()` |
| Read request headers | `r.Header.Get("Header-Name")` |
| Set response headers | `w.Header().Set("key", "value")` |
| Parse a form body | `r.ParseForm()` |
| Read a form field | `r.FormValue("field")` |
| Set status code explicitly | `w.WriteHeader(code)` |
| Get status code name | `http.StatusText(code)` |
| Custom router | `http.NewServeMux()` |
| Mount a submux | `http.StripPrefix("/prefix", mux)` |
| Define HTML template | `template.Must(template.New("name").Parse(str))` |
| Execute a template | `tmpl.Execute(w, data)` |

---

## Critical Rules to Remember

- **Headers before body** — always call `w.Header().Set()` and `w.WriteHeader()` before any `w.Write()` or `fmt.Fprintf()` call. Once the body starts writing, headers are permanently locked.
- **Close the body** — always `defer r.Body.Close()` after reading `r.Body` to free the connection resource.
- **Parse before reading forms** — call `r.ParseForm()` before `r.FormValue()` so you can handle parse errors explicitly.
- **Templates at package level** — parse templates once at startup, not inside the handler on every request.
