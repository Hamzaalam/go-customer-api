package main

import (
	"customer-api/db"
	"customer-api/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()

	router := mux.NewRouter()

	routes.RegisterCustomerRoutes(router)

	// Start the server
	log.Println("Server is running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
