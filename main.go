package main

import (
	"ecommerce/controller/auth"
	"ecommerce/controller/orderData"
	"ecommerce/controller/products"
	shoppingcart "ecommerce/controller/shopping-cart"
	"ecommerce/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	orders.GET("/", orderData.GetOrder)

	product := public.Group("/products")
	product.GET("/", products.GetAllProducts)
	product.POST("/create", products.PostProduct).Use(auth.Auth())
	product.PUT("/update", products.UpdateProduct).Use(auth.Auth())
	product.DELETE("/delete", products.DeleteProduct).Use(auth.Auth())

	shoppingCart := public.Group("/shoppingcart")
	shoppingCart.Use(auth.Auth())
	shoppingCart.GET("/", shoppingcart.PostShoppingCart)
	shoppingCart.POST("/create", shoppingcart.PostShoppingCart)
	shoppingCart.DELETE("/delete", shoppingcart.DeleteShoppingCart)

	//auth
	authRoute := public.Group("/auth")
	authRoute.GET("/")
	authRoute.POST("/login", auth.Login)
	authRoute.POST("/register", auth.PostRegister)
	r.Run(":8080")
}
