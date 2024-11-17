package main

import (
	"log"
	"net/http"
)

// Define the handler for the "/" route
func home(w http.ResponseWriter, r *http.Request) {
	// handle tailing slash issue of go
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Write a byte slice as the response
	w.Write([]byte("Hello from Snippetbox")) // res.send("Hello")
}

// Create a new snippet handler
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Check the HTTP method of the incoming request
	if r.Method != http.MethodPost {
		// response header map
		w.Header().Set("Allow", "POST")
		// Respond with a 405 status code (Method Not Allowed)
		// w.WriteHeader(http.StatusMethodNotAllowed)  instead of these
		// w.Write([]byte("Method Not Allowed"))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// If the request is a POST, handle the creation logic
	w.Write([]byte("Create a new snippet..."))
}

// View a specific snippet handler
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func main() {
	// Create a new servemux (router)
	mux := http.NewServeMux() // app = express() (create a router)

	// Register the home handler for the "/" route
	mux.HandleFunc("/", home) // app.get('/', handler)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Log that the server is starting
	log.Print("Starting server on :4000")

	// Start the HTTP server on port 4000 and use the servemux for routing
	err := http.ListenAndServe(":4000", mux) // app.listen(4000, () => console.log('Server running on 4000'))
	if err != nil {
		// If an error occurs, log it and exit
		log.Fatal(err)
	}
}
