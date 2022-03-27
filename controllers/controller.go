package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Print Functions

func PrintSuccess(status int, message string, w http.ResponseWriter) {
	var succResponse SuccessResponse
	succResponse.Status = status
	succResponse.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(succResponse)
}

func PrintError(status int, message string, w http.ResponseWriter) {
	var errResponse ErrorResponse
	errResponse.Status = status
	errResponse.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errResponse)
}

// Login
func CheckUserLogin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	name := r.URL.Query()["Name"]

	row := db.QueryRow("SELECT * FROM users WHERE name = ?", name[0])

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Age, &user.Address)

	if err != nil {
		PrintError(400, "User Not Found", w)
	} else {
		userType := 0
		generateToken(w, user.ID, user.Name, userType)
		PrintSuccess(200, "Logged In", w)
	}
}

// Logout
func Logout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)
	PrintSuccess(200, "Logged Out", w)
}

// ===GET===

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// DB + Query
	db := connect()
	defer db.Close()
	query := "SELECT * FROM users"

	// Get Data
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		PrintError(400, "Rows Are Empty - Users", w)
		return
	}

	// Insert Data To Array
	var user User
	var users []User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Age, &user.Address); err != nil {
			log.Fatal(err.Error())
			PrintError(400, "No User Data Inserted To []User", w)
			return
		} else {
			users = append(users, user)
		}
	}

	// Show Result
	var response UsersResponse
	if len(users) > 0 { // if < 1 -> for else checking
		response.Status = 200
		response.Message = "Success"
		response.Data = users
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		PrintError(400, "No User In []User", w)
		return
	}
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	// DB + Query
	db := connect()
	defer db.Close()
	query := "SELECT * FROM products"

	// Get Data
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		PrintError(400, "Rows Are Empty - Products", w)
		return
	}

	// Insert Data To Array
	var product Product
	var products []Product
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Fatal(err.Error())
			PrintError(400, "No Product Data Inserted To []Product", w)
			return
		} else {
			products = append(products, product)
		}
	}

	// Show Result
	var response ProductsResponse
	if len(products) > 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = products
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		PrintError(400, "No Product In []Product", w)
		return
	}
}

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	// DB + Query
	db := connect()
	defer db.Close()
	query := "SELECT * FROM transactions"

	// Get Data
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		PrintError(400, "Rows Are Empty - Transactions", w)
		return
	}

	// Insert Data To Array
	var transaction Transaction
	var transactions []Transaction
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantitiy); err != nil {
			log.Fatal(err.Error())
			PrintError(400, "No Transaction Data Inserted To []Transaction", w)
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}

	// Show Result
	var response TransactionsResponse
	if len(transactions) > 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = transactions
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		PrintError(400, "No Transaction In []Transaction", w)
		return
	}
}

func GetDetailTransaction(w http.ResponseWriter, r *http.Request) {
	// DB + Query
	db := connect()
	defer db.Close()

	// Get Data From Postman
	err := r.ParseForm()
	if err != nil {
		return
	}
	UserID := r.Form.Get("UserID")

	// Get Data
	query := `SELECT * FROM transactions 
		JOIN users ON transactions.UserID = users.ID 
		JOIN products ON transactions.ProductID = products.ID`

	// If There Is UserID Inserted In Postman
	if len(UserID) > 0 {
		query += ` WHERE transactions.UserID = ` + UserID
	}

	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		PrintError(400, "Rows Are Empty - Transactions", w)
		return
	}

	// Insert Data To Array
	var detailTransaction DetailTransaction
	var detailTransactions DetailTransactions
	for rows.Next() {
		if err := rows.Scan(&detailTransaction.ID,
			&detailTransaction.UserData.ID,
			&detailTransaction.ProductData.ID,
			&detailTransaction.Quantitiy,
			&detailTransaction.UserData.ID,
			&detailTransaction.UserData.Email,
			&detailTransaction.UserData.Password,
			&detailTransaction.UserData.Name,
			&detailTransaction.UserData.Age,
			&detailTransaction.UserData.Address,
			&detailTransaction.ProductData.ID,
			&detailTransaction.ProductData.Name,
			&detailTransaction.ProductData.Price); err != nil {
			log.Fatal(err.Error())
			PrintError(400, "No Product Data Inserted To []Product", w)
		} else {
			detailTransactions.Transactions = append(detailTransactions.Transactions, detailTransaction)
		}
	}

	// Show Result
	var response DetailTransactionsResponse

	if len(detailTransactions.Transactions) > 0 {
		response.Data = detailTransactions
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		PrintError(400, "No Detail Transaction In []detailTransactions", w)
		return
	}
}

// ===POST===

func InsertNewUser(w http.ResponseWriter, r *http.Request) {
	// DB
	db := connect()
	defer db.Close()

	// Get Data From Postman
	err := r.ParseForm()
	if err != nil {
		return
	}
	Name := r.Form.Get("Name")
	Email := r.Form.Get("Email")
	Password := r.Form.Get("Password")
	Age := r.Form.Get("Age")
	Address := r.Form.Get("Address")

	// Query
	_, errQuery := db.Exec("INSERT INTO users(Name, Email, Password, Age, Address) VALUES(?, ?, ?, ?, ?)",
		Name,
		Email,
		Password,
		Age,
		Address)

	// Show Result
	if errQuery == nil {
		PrintSuccess(200, "User Inserted", w)
	} else {
		PrintError(400, "Insert User Failed", w)
		return
	}
}

func InsertNewProduct(w http.ResponseWriter, r *http.Request) {
	// DB
	db := connect()
	defer db.Close()

	// Get Data From Postman
	err := r.ParseForm()
	if err != nil {
		return
	}
	Name := r.Form.Get("Name")
	Price := r.Form.Get("Price")

	// Query
	_, errQuery := db.Exec("INSERT INTO products(Name, Price) VALUES(?, ?)", Name, Price)

	// Show Result
	if errQuery == nil {
		PrintSuccess(200, "Product Inserted", w)
	} else {
		PrintError(400, "Insert Product Failed", w)
		return
	}
}

func InsertNewTransaction(w http.ResponseWriter, r *http.Request) {
	// DB
	db := connect()
	defer db.Close()

	// Get Data From Postman
	err := r.ParseForm()
	if err != nil {
		return
	}
	UserID := r.Form.Get("UserID")
	ProductID := r.Form.Get("ProductID")
	Quantity := r.Form.Get("Quantity")

	// Check Product Existence
	rows, err := db.Query("SELECT * FROM Products WHERE ID = ?", ProductID)

	if err != nil {
		log.Println(err)
		PrintError(400, "Fetch From Query Failed", w)
		return
	}

	var product Product
	var products []Product
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Fatal(err.Error())
			PrintError(400, "Error In Scanning Data", w)
			return
		} else {
			products = append(products, product)
		}
	}

	if len(products) > 0 {
		fmt.Print("Product Found")
	} else {
		db.Exec("INSERT INTO Products(ID) VALUES(?)", ProductID)
	}

	// Query
	// db.Exec("INSERT INTO Products(ID) VALUES(?)", ProductID)
	_, errQuery := db.Exec("INSERT INTO transactions(UserID, ProductID, Quantity) VALUES(?, ?, ?)", UserID, ProductID, Quantity)

	// Show Result
	if errQuery == nil {
		PrintSuccess(200, "Transaction Inserted", w)
	} else {
		PrintError(400, "Insert Transaction Failed", w)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	// DB
	db := connect()
	defer db.Close()

	// Get Data From Postman
	err := r.ParseForm()
	if err != nil {
		return
	}
	Email := r.Form.Get("Email")
	Password := r.Form.Get("Password")

	// Get Data
	rows, err := db.Query(`
	SELECT * FROM users 
	WHERE Email = ?
	AND Password = ?`,
		Email, Password)

	if err != nil {
		log.Println(err)
		PrintError(400, "Fetch From Query Failed", w)
		return
	}

	// Insert Data To Array
	var user User
	var users []User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Age, &user.Address); err != nil {
			log.Fatal(err.Error())
			PrintError(400, "No User Data Inserted To []User", w)
			return
		} else {
			users = append(users, user)
		}
	}

	// Show Result
	var response UsersResponse
	if len(users) > 0 { // if < 1 -> for else checking
		response.Status = 200
		response.Message = "Logged In As:"
		response.Data = users
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		PrintError(400, "Wrong Email / Password", w)
		return
	}
}

// ===PUT===

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// DB
	db := connect()
	defer db.Close()

	// Get Data From Postman
	err := r.ParseForm()
	if err != nil {
		return
	}
	ID := r.Form.Get("ID")
	Name := r.Form.Get("Name")
	Age := r.Form.Get("Age")
	Address := r.Form.Get("Address")

	// Query
	result, errQuery := db.Exec(
		"UPDATE users SET Name=?, Age=?, Address=? WHERE ID=?", Name, Age, Address, ID)

	num, _ := result.RowsAffected() //num, err

	// Show Result
	if errQuery == nil {
		if num == 0 {
			PrintError(400, "User Update Failed", w)
			return
		} else {
			PrintSuccess(200, "User Updated", w)
		}
	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// DB
	db := connect()
	defer db.Close()

	// Get Data From Postman
	err := r.ParseForm()
	if err != nil {
		return
	}
	ID := r.Form.Get("ID")
	Name := r.Form.Get("Name")
	Price := r.Form.Get("Price")

	// Query
	result, errQuery := db.Exec(
		"UPDATE products SET Name=?, Price=? WHERE ID=?", Name, Price, ID)

	num, _ := result.RowsAffected() //num, err

	// Show Result
	if errQuery == nil {
		if num == 0 {
			PrintError(400, "Product Update Failed", w)
			return
		} else {
			PrintSuccess(200, "Product Updated", w)
		}
	}
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	// DB
	db := connect()
	defer db.Close()

	// Get Data From Postman
	err := r.ParseForm()
	if err != nil {
		return
	}
	ID := r.Form.Get("ID")
	UserID := r.Form.Get("UserID")
	ProductID := r.Form.Get("ProductID")
	Quantitiy := r.Form.Get("Quantitiy")

	// Query
	result, errQuery := db.Exec(
		"UPDATE transactions SET UserID=?, ProductID=?, Quantity=? WHERE ID=?", UserID, ProductID, Quantitiy, ID)

	num, _ := result.RowsAffected() //num, err

	// Show Result
	if errQuery == nil {
		if num == 0 {
			PrintError(400, "Transaction Update Failed", w)
			return
		} else {
			PrintSuccess(200, "Transaction Updated", w)
		}
	}
}

// DELETE

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// DB
	db := connect()
	defer db.Close()

	// Get Data From Postman
	err := r.ParseForm()
	if err != nil {
		return
	}
	ID := r.Form.Get("ID")
	Name := r.Form.Get("Name")

	// Query
	result, errQuery := db.Exec(
		"DELETE FROM users WHERE ID=? OR Name=?", ID, Name)

	num, _ := result.RowsAffected() //num, err

	// Show Result
	if errQuery == nil {
		if num == 0 {
			PrintError(400, "Delete User Failed", w)
			return
		} else {
			PrintSuccess(200, "User Deleted", w)
		}
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// DB
	db := connect()
	defer db.Close()

	// Get Data From Postman
	err := r.ParseForm()
	if err != nil {
		return
	}
	ID := r.Form.Get("ID")
	Name := r.Form.Get("Name")

	// Query
	db.Exec("DELETE FROM transactions WHERE ProductID=?", ID)
	result, errQuery := db.Exec("DELETE FROM products WHERE ID=? OR Name=?", ID, Name)

	num, _ := result.RowsAffected() //num, err

	// Show Result
	if errQuery == nil {
		if num == 0 {
			PrintError(400, "Delete Product Failed", w)
			return
		} else {
			PrintSuccess(200, "Product & Related Transaction/s Deleted", w)
		}
	}
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	// DB
	db := connect()
	defer db.Close()

	// Get Data From Postman
	err := r.ParseForm()
	if err != nil {
		return
	}
	ID := r.Form.Get("ID")

	// Query
	result, errQuery := db.Exec(
		"DELETE FROM transactions WHERE ID=?", ID)

	num, _ := result.RowsAffected() //num, err

	// Show Result
	if errQuery == nil {
		if num == 0 {
			PrintError(400, "Delete Transaction Failed", w)
			return
		} else {
			PrintSuccess(200, "Transaction Deleted", w)
		}
	}
}
