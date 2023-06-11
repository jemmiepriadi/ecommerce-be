package shoppingcart

import (
	"ecommerce/model"
	"ecommerce/model/objects"
	deletedata "ecommerce/utils/deleteData"
	"ecommerce/utils/paginations"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetShoppingCart(c *gin.Context) {
	var pagination = paginations.GeneratePaginationFromRequest(c)
	var totalData int64
	var ShoppingCarts []model.ShoppingCart
	var res objects.Response

	var result *gorm.DB
	if c.Query("consumerID") != "" {
		id, err := strconv.Atoi(c.Query("SellerID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		result = model.DB.Model(&model.ShoppingCart{}).Where("ID = ?", id).Find(&ShoppingCarts)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		model.DB.Model(&result).Count(&totalData)

		pagination.TotalData = totalData
		totalPages := int(math.Ceil(float64(totalData) / float64(pagination.Size)))
		pagination.TotalPages = totalPages
		pagination.Data = ShoppingCarts
	}
	res.Data = pagination
	c.JSON(http.StatusOK, res)
}

func PostShoppingCart(c *gin.Context) {
	var shoppingcarts []model.ShoppingCart
	var shoppingCartsDB []model.ShoppingCart
	// var product model.Product
	var result *gorm.DB
	err := c.ShouldBindJSON(&shoppingcarts)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	consumerID, err := strconv.Atoi(c.Query("consumerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	result = model.DB.Where(&model.ShoppingCart{ConsumerID: consumerID}).Find(&shoppingCartsDB)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	shoppingCartsDB = shoppingcarts
	model.DB.Save(shoppingCartsDB)

}

func DeleteShoppingCart(c *gin.Context) {
	res := deletedata.DeleteItem(&model.ShoppingCart{}, c)
	c.JSON(http.StatusOK, res)
}
