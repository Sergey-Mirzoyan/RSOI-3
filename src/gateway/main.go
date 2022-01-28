package main

import (
	"errors"
	"fmt"
	"lab2/src/gateway/connector"
	"lab2/src/gateway/server"
	"os"
)

func main() {
	var err error = nil
	var config connector.Config
	prod, exists := os.LookupEnv("PROD")
	var address string
	if exists && prod == "true" {
		config = connector.Config{
			LibraryAddress:     "https://rsoi-2-library-mirzoyan.herokuapp.com/",
			ReservationAddress: "https://rsoi-2-reserve-mirzoyan.herokuapp.com/",
			RatingAddress:      "https://rsoi-2-rating-mirzoyan.herokuapp.com/",
			ApiPath:            "api/v1/",
		}
		port, exists := os.LookupEnv("PORT")
		if exists {
			address = "0.0.0.0:" + port
		} else {
			fmt.Print("Failed to get port")
			err = errors.New("Failed to get port")
		}
	} else {
		address = "127.0.0.1:31337"
		config = connector.Config{
			LibraryAddress:     "http://127.0.0.1:31338/",
			ReservationAddress: "http://127.0.0.1:31339/",
			RatingAddress:      "http://127.0.0.1:31340/",
			ApiPath:            "api/v1/",
		}
	}

	if err == nil {
		fmt.Printf("Starting server on %s\n", address)
		err = server.RunServer(address, &config)
		if err != nil {
			fmt.Print("Failed to start a server\n")
		}
	}
}
