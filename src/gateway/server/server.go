package server

import (
	"fmt"
	"lab2/src/gateway/connector"
	"lab2/src/gateway/connector/connector_implementation"
	"lab2/src/gateway/handler"
	"lab2/src/gateway/usecase/uc_implementation"
	"net/http"

	"github.com/gorilla/mux"
)

func RunServer(address string, config *connector.Config) error {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	con := connector_implementation.NewGatewayConnector(config)
	gc := uc_implementation.NewGatewayUsecase(con)
	gh := handler.NewGatewayHandlers(gc)

	apiRouter.HandleFunc("/libraries", gh.GetCityLibraries).Methods(http.MethodGet)
	apiRouter.HandleFunc("/libraries/{libraryUid:[0-9|a-z|\\-]+}/books", gh.GetLibraryBooks).Methods(http.MethodGet)
	apiRouter.HandleFunc("/reservations", gh.GetUserReservations).Methods(http.MethodGet)
	apiRouter.HandleFunc("/reservations", gh.RentBook).Methods(http.MethodPost)
	apiRouter.HandleFunc("/reservations/{reservationUid:[0-9|a-z|\\-]+}/return", gh.ReturnBook).Methods(http.MethodPost)
	apiRouter.HandleFunc("/rating", gh.GetUserRating).Methods(http.MethodGet)

	server := http.Server{
		Addr:    address,
		Handler: apiRouter,
	}

	fmt.Printf("Gateway service server is running on %s\n", address)
	return server.ListenAndServe()
}
