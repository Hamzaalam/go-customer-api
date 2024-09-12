package controllers

import (
	"customer-api/db"
	"customer-api/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetCustomer retrieves a customer by ID
func GetCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	err = db.DB.QueryRow("SELECT id, first_name, last_name, email FROM customers WHERE id = $1", id).
		Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email)

	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// GetAllCustomers retrieves all customers
func GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, first_name, last_name, email FROM customers")
	if err != nil {
		http.Error(w, "Error fetching customers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email)
		if err != nil {
			http.Error(w, "Error scanning customer", http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// CreateCustomer adds a new customer
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = db.DB.QueryRow(
		"INSERT INTO customers (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING id",
		customer.FirstName, customer.LastName, customer.Email).Scan(&customer.ID)

	if err != nil {
		http.Error(w, "Error creating customer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// UpdateCustomer updates an existing customer
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	err = json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(
		"UPDATE customers SET first_name = $1, last_name = $2, email = $3 WHERE id = $4",
		customer.FirstName, customer.LastName, customer.Email, id)

	if err != nil {
		http.Error(w, "Error updating customer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

// DeleteCustomer removes a customer by ID
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("DELETE FROM customers WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Error deleting customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
