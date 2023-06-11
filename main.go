package main

import (
	"ecommerce/controller/auth"
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
	orders.GET("/", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"haha": "haha"}) })

	product := public.Group("/products")
	product.GET("/", products.GetAllProducts)

	shoppingCart := public.Group("/shoppingcart")
	shoppingCart.Use(auth.Auth())
	shoppingCart.GET("/", shoppingcart.PostShoppingCart)
	shoppingCart.POST("/create", shoppingcart.PostShoppingCart)

	//auth
	authRoute := public.Group("/auth")
	authRoute.GET("/")
	authRoute.POST("/login", auth.Login)
	authRoute.POST("/register", auth.PostRegister)
	r.Run(":8080")
}
