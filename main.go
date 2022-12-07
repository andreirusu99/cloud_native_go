package main

import (
	"cloud_native_go/api"
	"cloud_native_go/core"
	"cloud_native_go/transact"
	"fmt"
	"log"
	"time"
)

type Effector func() (core.TransactionLogger, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func() (core.TransactionLogger, error) {
		for r := 0; ; r++ {
			res, err := effector()
			if err == nil || r >= retries { //call successful
				return res, err
			}
			fmt.Printf("Attempt %d failed: %v, retrying in %v\n", r+1, err, delay)
			select {
			case <-time.After(delay):
			}
		}
	}
}

func main() {
	fmt.Println("-> Server started")

	// create TransactionLogger
	retrySetupLogger := Retry(transact.NewTransactionLogger, 5, time.Second)
	logger, err := retrySetupLogger()

	if err != nil {
		log.Fatal(fmt.Errorf("could not initialize logger: %w", err))
	}
	fmt.Println("-> Logger initialized")

	// create store and pass it the TransactionLogger
	store := core.NewKVStore(logger)
	err = store.Restore()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("-> Store initialized")

	// create frontend
	frontEnd, err := api.NewFrontEnd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("-> FrontEnd initialized")

	fmt.Println("-> Listening for requests...")
	// pass the store to the frontend and start listening for requests
	log.Fatal(frontEnd.Start(store))
}
