package main

import (
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO!"))
}

func main() {

	http.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe("localhost:5555", nil))

}