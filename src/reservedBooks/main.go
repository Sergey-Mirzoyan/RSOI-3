package main

import (
	"errors"
	"fmt"
	"lab2/src/database/pgdb"
	"lab2/src/reservedBooks/handlers"
	"lab2/src/reservedBooks/repositories"
	"lab2/src/reservedBooks/usecases"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

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

	rr := repositories.NewReservationRepository(manager)
	rc := usecases.NewReservationUsecase(rr)
	rh := handlers.NewReservationHandlers(rc)

	apiRouter.HandleFunc("/reservations/count", rh.GetCountReservations).Methods(http.MethodGet)
	apiRouter.HandleFunc("/reservations", rh.GetUserReservations).Methods(http.MethodGet)
	apiRouter.HandleFunc("/reservations/{reservationUid:[0-9|a-z|\\-]+}", rh.GetReservation).Methods(http.MethodGet)
	apiRouter.HandleFunc("/reservations", rh.PostReservation).Methods(http.MethodPost)

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

	fmt.Printf("Reservations service server is running on %s\n", address)
	return server.ListenAndServe()
}

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
		address = "127.0.0.1:31339"
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
