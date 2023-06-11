package sellerdata

import (
	"ecommerce/model"
	"ecommerce/model/objects"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSeller(c *gin.Context) {
	var seller []model.Seller
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {

		return
	}
	result := model.DB.Model(&model.Seller{}).Where("ID = ?", id).Find(&seller)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	var products model.Product
	result = model.DB.Model(&model.Product{}).Where("seller_id = ?", seller[0].ID).Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	var orders []model.Order
	result = model.DB.Model(&model.Order{}).Where("seller_id = ?", seller[0].ID).Find(&orders)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	// seller[0].Product = products
	seller[0].Order = orders
	res := objects.Response{}
	res.Data = products
	c.JSON(http.StatusOK, res)
}
