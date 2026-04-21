package main

import (
	"log"
	"net/http"
)

func main() {
	h1 := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello from Service 1"))
	}
	http.HandleFunc("/hello", h1)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
