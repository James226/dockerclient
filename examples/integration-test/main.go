package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	// Start up a simple web server
	r := mux.NewRouter()
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	http.ListenAndServe(":10000", r)
}
