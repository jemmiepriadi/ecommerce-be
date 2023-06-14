package orderData

import (
	"ecommerce/model"
	"ecommerce/model/objects"
	"ecommerce/utils/paginations"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOrder(c *gin.Context) {
	var order []model.Order
	var pagination = paginations.GeneratePaginationFromRequest(c)
	var totalData int64
	var products []model.Product

	var res objects.Response

	pagination.Sort = "created_at asc"
	consumerID, err := strconv.Atoi(c.Query("consumerID"))
	if err != nil && c.Query("consumerID") != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	sellerID, err := strconv.Atoi(c.Query("sellerID"))
	if err != nil && c.Query("sellerID") != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	queryBuilder := model.DB.Offset(pagination.GetOffset()).Limit(pagination.GetSize()).Order(pagination.GetSort())

	result := queryBuilder.Model(&model.Order{}).Where(model.Order{ConsumerID: consumerID, SellerID: sellerID}).Find(&order)
	if result.Error != nil || ((c.Query("consumerID") != "" || c.Query("sellerID") != "") && len(order) == 0) {
		res.Message = "Data not found!"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	for _, value := range order {
		result = queryBuilder.Model(&model.ProductOrder{}).Where(model.ProductOrder{OrderID: value.ID}).Find(&products)
		if result.Error != nil || len(products) == 0 {
			res.Message = "Data not found!"
			c.JSON(http.StatusBadRequest, res)
			return
		}
		value.Product = products
	}

	model.DB.Model(&order).Count(&totalData)

	pagination.TotalData = totalData
	totalPages := int(math.Ceil(float64(totalData) / float64(pagination.Size)))
	pagination.TotalPages = totalPages
	pagination.Data = order
	res.Data = pagination
	c.JSON(http.StatusOK, res)

}

func CreateOrder(c *gin.Context) {
	//when checkout
	var res objects.Response
	var req []model.Order

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to parse request to struct: ", err)
		res.Code = "02"
		res.Message = "Failed parsing request"
		res.Data = err

		c.JSON(http.StatusOK, res)
		return
	}

	var array []interface{}
	for _, val := range req {
		if result := model.DB.Create(&val); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": result.Error,
			})
			return
		}
		for _, value := range val.Product {
			if result := model.DB.Save(&model.ProductOrder{ProductID: value.ID, OrderID: val.ID}); result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": result.Error,
				})
				return
			}
		}
		array = append(array, val)
	}
	res.Data = array

	res.Message = "succesfull create order"

	//remove related carts
	var shoppingCart model.ShoppingCart
	if err := model.DB.Table("shopping_carts").Where("consumer_id = ?", req[0].ConsumerID).First(&shoppingCart); err.Error != nil {
		res.Message = "Data not found!"
		res.Data = err.Error
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var productCart model.ProductCart

	if err := model.DB.Model(&model.ProductCart{}).Where("shopping_cart_id = ?", shoppingCart.ID).First(&productCart); err.Error != nil {
		res.Message = "Data not found!"
		res.Data = err.Error
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if result := model.DB.Delete(&shoppingCart); result.Error != nil {
		res.Message = "Delete Unsucessful"
		res.Data = result.Error
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if result := model.DB.Delete(&productCart); result.Error != nil {
		res.Message = "Delete Unsucessful"
		res.Data = result.Error
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, res)
}

func UpdateOrder(c *gin.Context) {
	var res objects.Response
	var order objects.Order
	var req model.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to parse request to struct: ", err)
		res.Code = "02"
		res.Message = "Failed parsing request"
		res.Data = err
		c.JSON(http.StatusBadRequest, res)
		return
	}
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil && c.Query("id") != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	if result := model.DB.Model(&req).Where("id =?", id).Find(&order); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	}
	order.Status = req.Status
	if result := model.DB.Save(&order); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	}
	res.Data = order
	res.Message = "succesfull create order"
	c.JSON(http.StatusOK, res)
}
