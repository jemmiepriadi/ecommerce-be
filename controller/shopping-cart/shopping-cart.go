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
	var ShoppingCarts []objects.ShoppingCart
	var res objects.Response
	var products []objects.Product

	queryBuilder := model.DB.Offset(pagination.GetOffset()).Limit(pagination.GetSize()).Order(pagination.GetSort())

	var result *gorm.DB
	if c.Query("consumerID") != "" {

		consumerID, err := strconv.Atoi(c.Query("sellerID"))
		if err != nil && c.Query("sellerID") != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		pagination.Sort = "created_at asc"
		result = queryBuilder.Model(&model.ShoppingCart{}).Where(&model.ShoppingCart{ConsumerID: consumerID}).Find(&ShoppingCarts)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		result = queryBuilder.Model(&model.Product{}).Where(&model.ShoppingCart{ConsumerID: consumerID}).Find(&ShoppingCarts)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		for _, value := range ShoppingCarts {
			result = queryBuilder.Model(&model.ProductCart{}).Where(model.ProductCart{CartID: value.ID}).Find(&products)
			if result.Error != nil || len(products) == 0 {
				res.Message = "Data not found!"
				c.JSON(http.StatusBadRequest, res)
				return
			}
			value.Product = products
		}
		model.DB.Model(&result).Count(&totalData)

		pagination.TotalData = totalData
		totalPages := int(math.Ceil(float64(totalData) / float64(pagination.Size)))
		pagination.TotalPages = totalPages
	}
	pagination.Data = ShoppingCarts

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
	var consumerID int
	if c.Query("consumerID") != "" {
		consumerID, err = strconv.Atoi(c.Query("consumerID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}
	result = model.DB.Where(&model.ShoppingCart{ConsumerID: consumerID}).Find(&shoppingCartsDB)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	shoppingCartsDB = shoppingcarts
	model.DB.Save(shoppingCartsDB)
	c.JSON(http.StatusOK, gin.H{"message": "Cart Added"})
}

func DeleteShoppingCart(c *gin.Context) {
	res := deletedata.DeleteItem(&model.ShoppingCart{}, c)
	c.JSON(http.StatusOK, res)
}
