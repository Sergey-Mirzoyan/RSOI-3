package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"lab2/src/database/pgdb"
	"lab2/src/rating/handlers"
	"lab2/src/rating/repositories/repo_implementation"
	"lab2/src/rating/usecases/uc_implementation"
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

	rr := repo_implementation.NewRatingRepository(manager)
	rc := uc_implementation.NewRatingUsecase(rr)
	rh := handlers.NewRatingHandlers(rc)

	apiRouter.HandleFunc("/rating", rh.GetUserRating).Methods(http.MethodGet)
	apiRouter.HandleFunc("/rating", rh.PatchUserRating).Methods(http.MethodPatch)

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

	fmt.Printf("Rating service server is running on %s\n", address)
	return server.ListenAndServe()
}
