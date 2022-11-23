package main

import (
	"cloud_native_go/pkg/service"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Server started")

	router := mux.NewRouter()

	service.SetupRoutes(router)
	fmt.Println("Routes initialised")

	if err := service.SetupLogger(); err != nil {
		log.Fatal(fmt.Errorf("could not initialize logger: %w", err))
	}
	fmt.Println("Logger initialised")

	fmt.Println("Listening for requests...")
	log.Fatal(http.ListenAndServeTLS(
		"localhost:5555",
		"C:\\Users\\andre\\Go\\cloud_native_go\\key\\cert.pem",
		"C:\\Users\\andre\\Go\\cloud_native_go\\key\\key.pem",
		router,
	))

}
