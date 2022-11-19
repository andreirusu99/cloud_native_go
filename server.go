package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO AGAIN!"))
}

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	// get value from body
	value, err := io.ReadAll(r.Body)

	// handle possible errors
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	// store value and handle possible errors
	if err = Put(key, string(value)); err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	// no errors so far => Success
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Stored value \"%s\" to key \"%s\"", value, key)))
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value, err := Get(key)

	// key not found
	if errors.Is(err, ErrorKeyNotFound) {
		http.Error(w,
			err.Error() + ": " + key,
			http.StatusNotFound,
		)
		return
	}

	// some other unexpected error
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Write([]byte(fmt.Sprintf("Key \"%s\" contains value \"%s\"", key, value)))
}

func keyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	if err := Delete(key); err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	Delete(key)
	w.Write([]byte(fmt.Sprintf("Key \"%s\" was deleted if present", key)))

}

func setupRoutes(router *mux.Router) {
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	router.HandleFunc("/v1/{key}", keyDeleteHandler).Methods("DELETE")
}

func main() {
	fmt.Println("Server started")

	router := mux.NewRouter()

	setupRoutes(router)

	fmt.Println("Routes initialised")
	fmt.Println("Listening...")

	log.Fatal(http.ListenAndServe("localhost:5555", router))

}