package main

import (
	"customer-api/db"
	"customer-api/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	router := mux.NewRouter()

	routes.RegisterCustomerRoutes(router)

	// Start the server
	log.Println("Server is running on port 8000...")

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
