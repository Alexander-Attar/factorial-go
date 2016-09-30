package main

import (
	"log"
	"net/http"

	"github.com/reservemedia/factorial-go/handler"
)

func main() {
	http.HandleFunc("/", handler.Factorial)
	log.Fatal(http.ListenAndServe("localhost:12345", nil))
}
