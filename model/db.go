package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	errDb error
)

func ConnectDataBase() {
	dsn := "root@tcp(127.0.0.1:3306)/ecommerce?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Account{}, &Seller{}, &Product{})

	DB = database
}

type Account struct {
	ID          int
	Name        string `json:"Name" example:"Jemmi"`
	UserType    string `json:"UserType" example:"seller"`
	PhoneNumber string `json:"PhoneNumber"`
	Username    string `gorm:"unique"`
	Password    string
	Consumer    Consumer
	Seller      Seller
}

type Consumer struct {
	ID        int
	Name      string
	AccountID int
	Order     []Order
}

type Seller struct {
	ID        int
	Name      string `json:"Name" example:"Jemmi"`
	AccountID int
	Product   []Product `json:"Products" gorm:"foreignkey:SellerID"`
	Order     []Order
}

type Product struct {
	ID          int
	SellerID    int
	Name        string
	Image       string
	Description string `json:"Description" example:"Berenang"`
	Price       int
	Order       []Order `gorm:"many2many:ProductOrder;"`
}

type Order struct {
	ID         int
	ConsumerID int
	SellerID   int
	Product    []Product `gorm:"many2many:ProductOrder;"`
	Status     bool
}

type ShoppingCart struct {
	ID        int
	Quantity  int
	SellerID  int
	ProductID int
}
