package objects

import (
	"time"

	"gorm.io/gorm"
)

type Response struct {
	Code    string      `json:"code" example:"00"`
	Message string      `json:"message" example:"Succesful"`
	Data    interface{} `json:"data"`
}

type Account struct {
	Name        string `json:"Name" example:"Jemmi"`
	UserType    string `json:"UserType" example:"seller"`
	PhoneNumber string `json:"PhoneNumber"`
	Username    string
	Password    string
	Consumer    Consumer
	Seller      Seller
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type Consumer struct {
	ID        int
	Name      string
	AccountID int
	Order     []Order
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type Seller struct {
	ID        int
	Name      string `json:"Name" example:"Jemmi"`
	AccountID int
	Product   []Product `json:"Products" gorm:"foreignkey:SellerID"`
	Order     []Order
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type Product struct {
	SellerID    int
	Name        string
	Image       string
	Description string `json:"Description" example:"Berenang"`
	Price       int
	Order       []Order `gorm:"many2many:ProductOrder;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type Order struct {
	ConsumerID int
	SellerID   int
	Product    []Product
	Status     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

type Checkout struct {
	Quantity   int
	SellerID   int
	ConsumerID int
	Product    []Product
	Status     bool
	OrderID    int
}

type ShoppingCart struct {
	Quantity  int
	SellerID  int
	ProductID int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type UserData struct {
	UserID      int
	Name        string
	Seller      Seller   `json:"omitempty"`
	Username    string   `json:"omitempty"`
	PhoneNumber string   `json:"omitempty"`
	Address     string   `json:"omitempty"`
	Consumer    Consumer `json:"omitempty"`
	UserType    string   `json:"omitempty"`
}

type JWT struct {
	JWT string
}
