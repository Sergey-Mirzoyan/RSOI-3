package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"lab2/src/database/pgdb"
	"lab2/src/library/handlers"
	"lab2/src/library/repositories/repo_implementation"
	"lab2/src/library/usecases/uc_implementation"
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

	lr := repo_implementation.NewLibraryRepository(manager)
	lc := uc_implementation.NewLibraryUsecase(lr)
	lh := handlers.NewLibraryHandlers(lc)
	br := repo_implementation.NewBooksRepository(manager)
	bc := uc_implementation.NewBooksUsecase(br)
	bh := handlers.NewBooksHandlers(bc)
	lbr := repo_implementation.NewLibraryBooksRepository(manager)
	lbc := uc_implementation.NewLibraryBooksUsecase(lbr)
	lbh := handlers.NewLibraryBooksHandlers(lbc)

	apiRouter.HandleFunc("/libraries", lh.GetByCity).Methods(http.MethodGet)
	apiRouter.HandleFunc("/libraries/{libraryUid:[0-9|a-z|\\-]+}", lh.GetByUid).Methods(http.MethodGet)
	apiRouter.HandleFunc("/libraries/{libraryUid:[0-9|a-z|\\-]+}/books", bh.GetByLibrary).Methods(http.MethodGet)
	apiRouter.HandleFunc("/books/{bookUid:[0-9|a-z|\\-]+}", bh.GetBook).Methods(http.MethodGet)
	apiRouter.HandleFunc("/books/{bookUid:[0-9|a-z|\\-]+}", bh.UpdateBook).Methods(http.MethodPut)
	apiRouter.HandleFunc("/library_books", lbh.UpdateBooksAmount).Methods(http.MethodPut)
	apiRouter.HandleFunc("/library_books", lbh.GetBooksAmount).Methods(http.MethodGet)

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

	fmt.Printf("Library system server is running on %s\n", address)
	return server.ListenAndServe()
}
