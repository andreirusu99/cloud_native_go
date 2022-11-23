package service

import (
	"cloud_native_go/pkg/config"
	store "cloud_native_go/pkg/db"
	"cloud_native_go/pkg/log"
	"cloud_native_go/pkg/misc"
	"errors"
	"fmt"
	"io"
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
	if err = store.Put(key, string(value)); err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	// no errors so far => Success
	logger.LogPut(key, string(value))
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("Stored value \"%s\" to key \"%s\"", value, key)))
	if err != nil {
		return
	}

	fmt.Printf("Handled PUT request: %s @ %s\n", key, value)
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value, err := store.Get(key)

	// key not found
	if errors.Is(err, store.ErrorKeyNotFound) {
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

func keyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	if err := store.Delete(key); err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	logger.LogDelete(key)
	_, err := w.Write([]byte(fmt.Sprintf("Key \"%s\" was deleted if present", key)))
	if err != nil {
		return
	}
	fmt.Printf("Handled DELETE request: %s\n", key)
}

func SetupRoutes(router *mux.Router) {
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	router.HandleFunc("/v1/{key}", keyDeleteHandler).Methods("DELETE")
}

var logger log.TransactionLogger

func SetupLogger() error {
	var err error

	fmt.Println("-> Creating Transaction Logger...")
	logger, err = log.NewPostgresTransactionLogger(config.PostgresConfig)

	if err != nil {
		return fmt.Errorf("failed to create transaction logger: %w", err)
	}
	fmt.Println("-> Done creating Transaction Logger")

	fmt.Println("-> Replaying Logs...")
	events, replayErrors := logger.ReplayEvents()

	event, ok := misc.Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <-replayErrors: // got an error
			return fmt.Errorf("error replaying logs: %w", err)

		case event, ok = <-events: // got an event

			switch event.Type {
			case misc.EventPut:
				err = store.Put(event.Key, event.Value)

			case misc.EventDelete:
				err = store.Delete(event.Key)
			}
		}
	}

	fmt.Printf("-> Recreated DataStore: %v\n", store.GetAll())
	fmt.Println("-> Done replaying Logs")

	fmt.Println("-> Starting the Logger...")
	logger.Run()
	fmt.Println("-> Done starting the Logger")

	go func() {
		for err = range logger.Err() {
			fmt.Printf("logger error: %v", err)
		}
	}()

	return err
}
