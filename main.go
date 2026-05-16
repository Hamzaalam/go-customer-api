package main

import (
	"customer-api/db"
	_ "customer-api/docs"
	"customer-api/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Customer API
// @version 1.0
// @description API documentation for the Customer API.
// @host localhost:8000
// @BasePath /
func main() {
	db.InitDB()
	defer db.DB.Close()

	router := mux.NewRouter()

	routes.RegisterCustomerRoutes(router)
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start the server
	log.Println("Server is running on port 8000...")
	log.Println("Swagger documentation is available at http://localhost:8000/swagger/index.html")

	// Graceful shutdown handler
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(":8000", router))
}
