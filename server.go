package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO AGAIN!"))
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe("localhost:5555", router))

}