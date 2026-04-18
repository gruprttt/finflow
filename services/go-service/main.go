package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go Service 🚀")
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Go service running on port 8080")
	http.ListenAndServe(":8080", nil)
}
