package main

import (
	"log"

	"net/http"

	"github.com/burakiscoding/go-book-rent/api"
	"github.com/burakiscoding/go-book-rent/database"
	"github.com/burakiscoding/go-book-rent/helpers"
	"github.com/burakiscoding/go-book-rent/store"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewSQL()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	bookStore := store.NewBookStore(db)
	bookHandler := api.NewBookHandler(*bookStore)
	subrouter.HandleFunc("/books", helpers.MakeHandler(bookHandler.HandleGetAll)).Methods(http.MethodGet)
	subrouter.HandleFunc("/books/{id}", helpers.MakeHandler(bookHandler.HandleGetById)).Methods(http.MethodGet)
	subrouter.HandleFunc("/books", helpers.MakeHandler(api.HandleAdminAuth(bookHandler.HandleInsert))).Methods(http.MethodPost)
	subrouter.HandleFunc("/books/{id}", helpers.MakeHandler(api.HandleAdminAuth(bookHandler.HandleUpdate))).Methods(http.MethodPut)
	subrouter.HandleFunc("/books/{id}", helpers.MakeHandler(api.HandleAdminAuth(bookHandler.HandleDelete))).Methods(http.MethodDelete)

	userStore := store.NewUserStore(db)
	userHandler := api.NewUserHandler(*userStore)
	subrouter.HandleFunc("/user/register", helpers.MakeHandler(userHandler.HandleRegister)).Methods(http.MethodPost)
	subrouter.HandleFunc("/user/login", helpers.MakeHandler(userHandler.HandleLogin)).Methods(http.MethodPost)
	subrouter.HandleFunc("/user/details", helpers.MakeHandler(api.HandleAuth(userHandler.HandleGetDetails))).Methods(http.MethodPost)
	subrouter.HandleFunc("/user/admin-register", helpers.MakeHandler(userHandler.HandleAdminRegister)).
		Host("localhost").Methods(http.MethodPost)

	rentStore := store.NewRentStore(db)
	rentHandler := api.NewRentHandler(*rentStore, *bookStore)
	subrouter.HandleFunc("/rent/book", helpers.MakeHandler(api.HandleAuth(rentHandler.HandleRentBook))).Methods(http.MethodPost)
	subrouter.HandleFunc("/rent/history", helpers.MakeHandler(api.HandleAdminAuth(rentHandler.HandleGetAllHistory))).Methods(http.MethodGet)
	subrouter.HandleFunc("/rent/return", helpers.MakeHandler(api.HandleAuth(rentHandler.HandleReturnBook))).Methods(http.MethodPost)
	subrouter.HandleFunc("/rent/user-history", helpers.MakeHandler(api.HandleAuth(rentHandler.HandleGetUserHistory))).Methods(http.MethodGet)

	http.ListenAndServe(":8080", router)
}
