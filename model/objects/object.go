package objects

import "time"

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
}

type Consumer struct {
	Name      string
	AccountID int
	Order     []Order
}

type Seller struct {
	Name      string `json:"Name" example:"Jemmi"`
	AccountID int
	Product   []Product `json:"Products" gorm:"foreignkey:SellerID"`
	Order     []Order
}

type Product struct {
	SellerID    int
	Name        string
	Image       string
	Description string `json:"Description" example:"Berenang"`
	Price       int
	Order       []Order `gorm:"many2many:ProductOrder;"`
	Created_at  time.Time
}

type Order struct {
	ConsumerID int
	SellerID   int
	Product    []Product `gorm:"many2many:ProductOrder;"`
	Status     bool
}

type ShoppingCart struct {
	Quantity  int
	SellerID  int
	ProductID int
}
