package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"lab2/src/database/pgdb"
	"lab2/src/reservation/handlers"
	"lab2/src/reservation/repositories/repo_implementation"
	"lab2/src/reservation/usecases/uc_implementation"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	rr := repo_implementation.NewReservationRepository(manager)
	rc := uc_implementation.NewReservationUsecase(rr)
	rh := handlers.NewReservationHandlers(rc)

	apiRouter.HandleFunc("/reservations/count", rh.GetUserReservationsCount).Methods(http.MethodGet)
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
