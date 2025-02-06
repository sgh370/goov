package testdata

import "time"

type TestUser struct {
	Username string
	Age      int
	Email    string
}

type Address struct {
	Street  string
	City    string
	Country string
	ZipCode string
}

type Contact struct {
	Type  string
	Value string
}

type Order struct {
	ID          int
	CustomerID  int
	Status      string
	Items       []OrderItem
	BillingAddr Address
	ShippingAddr *Address
	Contacts    map[string]Contact
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderItem struct {
	ProductID   int
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
	Description string
}

type User struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
	Age             int
	Premium         bool
	PremiumDetails  string
}
