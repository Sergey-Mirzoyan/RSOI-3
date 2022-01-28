package main

import (
	"errors"
	"fmt"
	"lab1/database/pgdb"
	"lab1/handlers"
	"lab1/repositories"
	"lab1/usecases"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func MakeConStr(host string, port int, dbName string, user string, password string) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?", user, password, host, port, dbName)
}

func RunServer(address string, connectionString string) error {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	manager := pgdb.NewPGDBManager()
	err := manager.Connect(connectionString)
	if err != nil {
		fmt.Print("Failed to connect to database")
	} else {
		fmt.Println("Successfully connected to postgres database")
	}

	pr := repositories.NewPersonsRepository(manager)
	pc := usecases.NewPersonsUsecase(pr)
	ph := handlers.NewPerosnsHandlers(pc)

	apiRouter.HandleFunc("/persons", ph.PostPerson).Methods(http.MethodPost)
	apiRouter.HandleFunc("/persons", ph.GetPersons).Methods(http.MethodGet)
	apiRouter.HandleFunc("/persons/{id:[0-9]+}", ph.GetPersonById).Methods(http.MethodGet)
	apiRouter.HandleFunc("/persons/{id:[0-9]+}", ph.UpdatePerson).Methods(http.MethodPatch)
	apiRouter.HandleFunc("/persons/{id:[0-9]+}", ph.DeletePerson).Methods(http.MethodDelete)

	server := http.Server{
		Addr:    address,
		Handler: apiRouter,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		manager.Disconnect()
		os.Exit(0)
	}()

	fmt.Printf("Server is running on %s\n", address)
	return server.ListenAndServe()
}

func main() {
	var err error = nil
	prod, exists := os.LookupEnv("PS_PROD")
	var address, conString string
	if exists && prod == "true" {
		conString, exists = os.LookupEnv("DATABASE_URL")
		if exists {
			port, exists := os.LookupEnv("PORT")
			if exists {
				address = "0.0.0.0:" + port
			} else {
				fmt.Print("Failed to get port")
				err = errors.New("Failed to get port")
			}
		} else {
			fmt.Print("Failed to get a connection string")
			err = errors.New("Failed to get a connection string")
		}
	} else {
		address = "127.0.0.1:31337"
		conString = MakeConStr("127.0.0.1", 5432, "persons", "program", "test")
	}

	if err == nil {
		fmt.Printf("Starting server on %s\n", address)
		err = RunServer(address, conString)
		if err != nil {
			fmt.Print("Failed to start a server\n")
		}
	}
}
