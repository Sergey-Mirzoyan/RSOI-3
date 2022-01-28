package main

import (
	"errors"
	"fmt"
	"lab2/src/library/server"
	"os"
)

func MakeConStr(host string, port int, dbName string, user string, password string) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?", user, password, host, port, dbName)
}

func main() {
	var err error = nil
	prod, exists := os.LookupEnv("PROD")
	var address, conString string
	if exists && prod == "true" {
		conString, exists = os.LookupEnv("DATABASE_URL")
		if exists {
			port, exists := os.LookupEnv("PORT")
			if exists {
				address = "0.0.0.0:"+port
			} else {
				fmt.Print("Failed to get port")
				err = errors.New("Failed to get port")
			}
		} else {
			fmt.Print("Failed to get a connection string")
			err = errors.New("Failed to get a connection string")
		}
	} else {
		address = "127.0.0.1:31338"
		conString = MakeConStr("127.0.0.1", 5432, "lab2", "program", "test")
	}

	if err == nil {
		fmt.Printf("Starting server on %s\n", address)
		err = server.RunServer(address, conString)
		if err != nil {
			fmt.Print("Failed to start a server\n")
		}
	}
}