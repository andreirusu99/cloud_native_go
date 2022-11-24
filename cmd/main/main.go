package main

import (
	"cloud_native_go/pkg/service"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Effector func() error

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func() error {
		for r := 0; ; r++ {
			err := effector()
			if err == nil || r >= retries { //call successful
				return err
			}
			fmt.Printf("Attempt %d failed: %v, retrying in %v\n", r+1, err, delay)
			select {
			case <-time.After(delay):
			}
		}
	}
}

func main() {
	fmt.Println("Server started")

	router := mux.NewRouter()

	service.SetupRoutes(router)
	fmt.Println("Routes initialised")

	retrySetupLogger := Retry(service.SetupLogger, 5, 1*time.Second)

	if err := retrySetupLogger(); err != nil {
		log.Fatal(fmt.Errorf("could not initialize logger: %w", err))
	}
	fmt.Println("Logger initialised")

	fmt.Println("Listening for requests...")
	log.Fatal(http.ListenAndServeTLS(
		"localhost:5555",
		"cert.pem",
		"key.pem",
		router,
	))

}
