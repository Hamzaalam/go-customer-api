package routes

import (
	"customer-api/controllers"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

// ValidateIDMiddleware ensures ID parameter is a valid number
func ValidateIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		if id, ok := params["id"]; ok {
			if !regexp.MustCompile(`^\d+$`).MatchString(id) {
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func RegisterCustomerRoutes(router *mux.Router) {
	router.Use(ValidateIDMiddleware)

	router.HandleFunc("/api/customers", controllers.GetAllCustomers).Methods("GET")
	router.HandleFunc("/api/customers/{id}", controllers.GetCustomer).Methods("GET")
	router.HandleFunc("/api/customers", controllers.CreateCustomer).Methods("POST")
	router.HandleFunc("/api/customers/{id}", controllers.UpdateCustomer).Methods("PUT")
	router.HandleFunc("/api/customers/{id}", controllers.DeleteCustomer).Methods("DELETE")
}
