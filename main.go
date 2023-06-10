package main

import (
	"ecommerce/controller/auth"
	"ecommerce/model"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

type JWTClaim struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

var res = Response{
	Code:    "00",
	Message: "Success",
}

func main() {
	fmt.Println("hahahah")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	model.ConnectDataBase()
	public := r.Group("/api")
	orders := public.Group("/orders")
	orders.GET("/", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"haha": "haha"}) })
	//auth
	authRoute := public.Group("/auth")
	authRoute.GET("/")
	authRoute.POST("/login", auth.Login)
	authRoute.POST("/register", auth.PostRegister)
	r.Run(":8080")
}
