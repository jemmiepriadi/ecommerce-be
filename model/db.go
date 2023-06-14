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
	if err := database.SetupJoinTable(&ShoppingCart{}, "Product", &ProductCart{}); err != nil {
		println(err.Error())
		panic("Failed to setup join table")
	}

	if err := database.SetupJoinTable(&Order{}, "Product", &ProductOrder{}); err != nil {
		println(err.Error())
		panic("Failed to setup join table")
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
	Address     string
	Consumer    Consumer
	Seller      Seller
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Consumer struct {
	ID           int
	Name         string
	AccountID    int `gorm:"unique;not null"`
	ShoppingCart []ShoppingCart
	Order        []Order
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Seller struct {
	ID        int
	Name      string    `json:"Name" example:"Jemmi"`
	AccountID int       `gorm:"unique;not null"`
	Product   []Product `json:"Products" gorm:"foreignkey:SellerID"`
	Order     []Order
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Product struct {
	ID           int
	SellerID     int
	Name         string
	Image        string
	Description  string `json:"Description" example:"Berenang"`
	Price        int
	Quantity     int
	ShoppingCart []ShoppingCart `gorm:"many2many:ProductCart;"`
	Order        []Order        `gorm:"many2many:ProductOrder;"`
	CreatedAt    time.Time      `json:"CreatedAt"`
	UpdatedAt    time.Time      `json:"UpdatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Order struct {
	ID          int
	ConsumerID  int
	SellerID    int
	Payment     string
	Email       string
	City        string
	TotalPrice  int
	State       string
	ZipCode     int
	PaymentInfo int
	Product     []Product `gorm:"many2many:ProductOrder;"`
	Status      string
	CreatedAt   time.Time      `json:"CreatedAt"`
	UpdatedAt   time.Time      `json:"UpdatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type ShoppingCart struct {
	ID         int
	Quantity   int
	ConsumerID int
	Product    []Product `gorm:"many2many:ProductCart;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type ProductOrder struct {
	ProductID int `gorm:"primaryKey"`
	OrderID   int `gorm:"primaryKey"`
	Quantity  int
	CreatedAt time.Time      `json:"CreatedAt"`
	UpdatedAt time.Time      `json:"UpdatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ProductCart struct {
	ProductID      int `gorm:"primaryKey"`
	ShoppingCartID int `gorm:"primaryKey"`
	Quantity       int
	CreatedAt      time.Time      `json:"CreatedAt"`
	UpdatedAt      time.Time      `json:"UpdatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
