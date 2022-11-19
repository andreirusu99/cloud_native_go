package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"cloud_native_go/pkg/service"
)

func main() {
	fmt.Println("Server started")

	router := mux.NewRouter()

	service.SetupRoutes(router)
	fmt.Println("Routes initialised")

	service.SetupLogger("C:\\Users\\andre\\Go\\cloud_native_go\\out\\transaction.log")
	fmt.Println("Logger initialised")

	fmt.Println("Listening for requests...")
	log.Fatal(http.ListenAndServe("localhost:5555", router))

}
