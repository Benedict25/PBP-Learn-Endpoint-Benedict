package controllers

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
}

type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductsResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Product `json:"data"`
}

type Transaction struct {
	ID        int `json:"id"`
	UserID    int `json:"userid"`
	ProductID int `json:"productid"`
	Quantitiy int `json:"quantity"`
}

type TransactionsResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []Transaction `json:"data"`
}

type DetailTransaction struct {
	ID          int     `json:"ID"`
	UserData    User    `json:"User"`
	ProductData Product `json:"Product"`
	Quantitiy   int     `json:"Quantity"`
}

type DetailTransactions struct {
	Transactions []DetailTransaction `json:"Transactions"`
}

type DetailTransactionsResponse struct {
	Data DetailTransactions `json:"Data"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
