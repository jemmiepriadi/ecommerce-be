package main

import (
	"ecommerce/controller/auth"
	"ecommerce/controller/orderData"
	"ecommerce/controller/products"
	sellerdata "ecommerce/controller/sellerData"
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

	r.Use(auth.CORSMiddleware())
	model.ConnectDataBase()
	public := r.Group("/api")

	seller := public.Group("/sellers")
	seller.GET("/", sellerdata.GetSeller)

	orders := public.Group("/orders")
	orders.GET("/", orderData.GetOrder)
	orders.POST("/", orderData.CreateOrder)
	orders.PUT("/update", auth.Auth(), orderData.UpdateOrder)

	product := public.Group("/products")
	product.GET("/", products.GetAllProducts)
	product.POST("/create", auth.Auth(), products.PostProduct)
	product.PUT("/update", auth.Auth(), products.UpdateProduct)
	product.DELETE("/delete", auth.Auth(), products.DeleteProduct)

	shoppingCart := public.Group("/shoppingcart")
	// shoppingCart.Use(auth.Auth())
	shoppingCart.GET("/", shoppingcart.GetShoppingCart)
	shoppingCart.POST("/", shoppingcart.PostShoppingCart)
	shoppingCart.DELETE("/delete", shoppingcart.DeleteShoppingCart)

	//auth
	authRoute := public.Group("/auth")
	authRoute.GET("/")
	authRoute.GET("/me", auth.Me)
	authRoute.POST("/login", auth.Login)
	authRoute.POST("/register", auth.PostRegister)
	r.Run(":8080")
}
