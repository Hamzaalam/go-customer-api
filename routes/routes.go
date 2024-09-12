package routes

import (
	"customer-api/controllers"

	"github.com/gorilla/mux"
)

func RegisterCustomerRoutes(router *mux.Router) {
	router.HandleFunc("/api/customers", controllers.GetAllCustomers).Methods("GET")
	router.HandleFunc("/api/customers/{id}", controllers.GetCustomer).Methods("GET")
	router.HandleFunc("/api/customers", controllers.CreateCustomer).Methods("POST")
	router.HandleFunc("/api/customers/{id}", controllers.UpdateCustomer).Methods("PUT")
	router.HandleFunc("/api/customers/{id}", controllers.DeleteCustomer).Methods("DELETE")
}
