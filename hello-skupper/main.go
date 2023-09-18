package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
        hostname, err := os.Hostname()
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        fmt.Fprintf(w, "Welcome from %s\n", hostname)
}

func main() {
        http.HandleFunc("/", handler)
        fmt.Println("Server is listening on :8080...")
        http.ListenAndServe(":8080", nil)
}
