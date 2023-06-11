package model

import (
	"time"

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

	database.AutoMigrate(&Account{}, &Seller{}, &Consumer{}, &Product{}, &Order{}, &ShoppingCart{})

	DB = database
}

type Account struct {
	ID          int
	Name        string `json:"Name" example:"Jemmi"`
	UserType    string `json:"UserType" example:"seller"`
	PhoneNumber string `json:"PhoneNumber"`
	Username    string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	Consumer    Consumer
	Seller      Seller
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Consumer struct {
	ID           int
	Name         string
	AccountID    int
	ShoppingCart []ShoppingCart
	Order        []Order
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Seller struct {
	ID        int
	Name      string `json:"Name" example:"Jemmi"`
	AccountID int
	Product   []Product `json:"Products" gorm:"foreignkey:SellerID"`
	Order     []Order
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Product struct {
	ID          int
	SellerID    int
	Name        string
	Image       string
	Description string `json:"Description" example:"Berenang"`
	Price       int
	Order       []Order        `gorm:"many2many:ProductOrder;"`
	CreatedAt   time.Time      `json:"CreatedAt"`
	UpdatedAt   time.Time      `json:"UpdatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Order struct {
	ID         int
	ConsumerID int
	SellerID   int
	Product    []Product `gorm:"many2many:ProductOrder;"`
	Status     bool
	CreatedAt  time.Time      `json:"CreatedAt"`
	UpdatedAt  time.Time      `json:"UpdatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type ShoppingCart struct {
	ID         int
	Quantity   int
	ConsumerID int
	ProductID  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type ProductOrder struct {
	ProductID int `gorm:"primaryKey"`
	OrderID   int `gorm:"primaryKey"`
}
