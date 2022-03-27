package main

import (
	controllers "Endpoint/controllers"

	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// End Point

	// Login
	router.HandleFunc("/login", controllers.CheckUserLogin).Methods("POST")

	// Logout
	router.HandleFunc("/logout", controllers.Logout).Methods("POST")

	// GET
	router.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	router.HandleFunc("/transactions", controllers.GetAllTransactions).Methods("GET")
	router.HandleFunc("/detail_transactions", controllers.GetDetailTransaction).Methods("GET")

	// POST
	router.HandleFunc("/users", controllers.InsertNewUser).Methods("POST")
	router.HandleFunc("/products", controllers.InsertNewProduct).Methods("POST")
	router.HandleFunc("/transactions", controllers.InsertNewTransaction).Methods("POST")
	// router.HandleFunc("/login", controllers.Login).Methods("POST")

	// PUT
	router.HandleFunc("/users", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/products", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/transactions", controllers.UpdateTransaction).Methods("PUT")

	// DELETE
	router.HandleFunc("/users", controllers.Authenticate(controllers.DeleteUser, 1)).Methods("DELETE")
	router.HandleFunc("/products", controllers.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/transactions", controllers.DeleteTransaction).Methods("DELETE")

	// Connection Notif
	http.Handle("/", router)
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
