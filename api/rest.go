package api

import (
	"cloud_native_go/core"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type RestFrontEnd struct {
	store *core.KVStore
}

func NewRestFrontEnd() (*RestFrontEnd, error) {
	return &RestFrontEnd{}, nil
}

func (restFrontEnd *RestFrontEnd) keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
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
	if err = restFrontEnd.store.Put(key, string(value), false); err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	// no errors so far => Success
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("Stored value \"%s\" to key \"%s\"", value, key)))
	if err != nil {
		return
	}

	fmt.Printf("Handled PUT request: %s @ %s\n", key, value)
}

func (restFrontEnd *RestFrontEnd) keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value, err := restFrontEnd.store.Get(key)

	// key not found
	if errors.Is(err, core.ErrorKeyNotFound) {
		http.Error(w,
			err.Error()+": "+key,
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

	_, err = w.Write([]byte(fmt.Sprintf("Key \"%s\" contains value \"%s\"", key, value)))
	if err != nil {
		return
	}

	fmt.Printf("Handled GET request: %s\n", key)
}

func (restFrontEnd *RestFrontEnd) keyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	if err := restFrontEnd.store.Delete(key, false); err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	_, err := w.Write([]byte(fmt.Sprintf("Key \"%s\" was deleted if present", key)))
	if err != nil {
		return
	}
	fmt.Printf("Handled DELETE request: %s\n", key)
}

func (restFrontEnd *RestFrontEnd) Start(store *core.KVStore) error {
	restFrontEnd.store = store
	router := mux.NewRouter()

	setupRoutes(router, restFrontEnd)

	return http.ListenAndServeTLS(
		"0.0.0.0:5555",
		"cert.pem",
		"key.pem",
		router,
	)
}

func setupRoutes(router *mux.Router, restFrontEnd *RestFrontEnd) {
	router.HandleFunc("/v1/{key}", restFrontEnd.keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/v1/{key}", restFrontEnd.keyValueGetHandler).Methods("GET")
	router.HandleFunc("/v1/{key}", restFrontEnd.keyDeleteHandler).Methods("DELETE")
}
