package shoppingcart

import (
	"ecommerce/model"
	"ecommerce/model/objects"
	"ecommerce/utils/paginations"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetShoppingCart(c *gin.Context) {
	var pagination = paginations.GeneratePaginationFromRequest(c)
	var totalData int64
	var shoppingCarts []objects.ShoppingCart
	var res objects.Response
	var productCarts []model.ProductCart
	var products objects.Product
	var seller model.Seller

	queryBuilder := model.DB.Offset(pagination.GetOffset()).Limit(pagination.GetSize()).Order(pagination.GetSort())

	var result *gorm.DB
	if c.Query("consumerID") != "" {

		consumerID, err := strconv.Atoi(c.Query("consumerID"))
		if err != nil && c.Query("sellerID") != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		pagination.Sort = "created_at asc"
		result = queryBuilder.Model(&model.ShoppingCart{}).Where(&model.ShoppingCart{ConsumerID: consumerID}).Find(&shoppingCarts)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		result = queryBuilder.Model(&model.Product{}).Where(&model.ShoppingCart{ConsumerID: consumerID}).Find(&shoppingCarts)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		result = model.DB.Model(&model.ProductCart{}).Where(&model.ProductCart{ShoppingCartID: shoppingCarts[0].ID}).Find(&productCarts)
		if result.Error != nil {
			res.Message = "Data not ssnnsas!"
			c.JSON(http.StatusBadRequest, res)
			return
		}
		for _, value := range productCarts {

			result = model.DB.Model(&model.Product{}).Where("ID = ?", value.ProductID).Select("id,seller_id,name,image,description,price,quantity,created_at,updated_at,deleted_at").First(&products)
			if result.Error != nil {
				res.Message = "Data not found!"
				c.JSON(http.StatusBadRequest, res)
				return
			}
			fmt.Println(value.ProductID)
			fmt.Println(products.SellerID)
			products.Quantity = value.Quantity

			result = model.DB.Model(&model.Seller{}).Where("ID = ?", products.SellerID).First(&seller)
			if result.Error != nil {
				res.Message = "Data not found!"
				c.JSON(http.StatusBadRequest, res)
				return
			}

			products.SellerName = seller.Name
			shoppingCarts[0].Product = append(shoppingCarts[0].Product, products)
		}

		model.DB.Model(&result).Count(&totalData)

		pagination.TotalData = totalData
		totalPages := int(math.Ceil(float64(totalData) / float64(pagination.Size)))
		pagination.TotalPages = totalPages
	}
	pagination.Data = shoppingCarts

	res.Data = pagination
	c.JSON(http.StatusOK, res)
}

func PostShoppingCart(c *gin.Context) {
	var shoppingcarts model.ShoppingCart
	var shoppingCartsDB []model.ShoppingCart
	// var product model.Product
	var result *gorm.DB
	var productCarts []model.ProductCart
	var res objects.Response
	err := c.ShouldBindJSON(&shoppingcarts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "cannot parse"})
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
	shoppingcarts.ConsumerID = consumerID
	result = model.DB.Where(&model.ShoppingCart{ConsumerID: consumerID}).Find(&shoppingCartsDB)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if len(shoppingCartsDB) > 0 {
		shoppingCartsDB[0].ConsumerID = shoppingcarts.ConsumerID
		shoppingCartsDB[0].Quantity = shoppingcarts.Quantity
		shoppingCartsDB[0].Product = shoppingcarts.Product
		if result := model.DB.Save(shoppingCartsDB); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": result.Error})
			return
		}

		for _, value := range shoppingcarts.Product {

			result = model.DB.Model(&model.ProductCart{}).Where(&model.ProductCart{ShoppingCartID: shoppingCartsDB[0].ID, ProductID: value.ID}).Find(&productCarts)
			if result.Error != nil {
				res.Message = "Data not found!"
				c.JSON(http.StatusBadRequest, res)
				return
			}

			if len(productCarts) > 0 {
				productCarts[0].ProductID = value.ID
				productCarts[0].Quantity = value.Quantity
				productCarts[0].ShoppingCartID = shoppingCartsDB[0].ID
				model.DB.Save(productCarts)
			} else {
				model.DB.Save(&model.ProductCart{ProductID: value.ID, ShoppingCartID: shoppingCartsDB[0].ID, Quantity: value.Quantity})

			}
		}
	} else {

		if result := model.DB.Save(&shoppingcarts); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": result.Error})
			return
		}
		for _, value := range shoppingcarts.Product {
			model.DB.Save(&model.ProductCart{ProductID: value.ID, ShoppingCartID: shoppingcarts.ID, Quantity: value.Quantity})
		}
	}

	res.Data = shoppingcarts
	res.Message = "Cart Added!"
	c.JSON(http.StatusOK, res)
}

// delete productcart
func DeleteProductCart(c *gin.Context) {
	var productCart model.ProductCart
	var res objects.Response
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil && c.Query("id") != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	if err := model.DB.Model(&model.ProductCart{}).Where("product_id = ?", id).First(&productCart); err.Error != nil {
		res.Message = "Data not found!"
		res.Data = err.Error
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if result := model.DB.Delete(&productCart); result.Error != nil {
		res.Message = "Delete Unsucessful"
		res.Data = err.Error
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res.Message = "Delete Successful"
	c.JSON(http.StatusOK, res)
}

// delete entire cart
func DeleteShoppingCart(c *gin.Context) {
	// res := deletedata.DeleteItem(&model.ShoppingCart{}, c) => not needed, will be deleted later
	id, err := strconv.Atoi(c.Query("id"))
	var shoppingCart model.ShoppingCart

	if err != nil && c.Query("id") != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	res := &objects.Response{}

	if err := model.DB.Table("shopping_carts").Where("ID = ?", id).First(&shoppingCart); err.Error != nil {
		res.Message = "Data not found!"
		res.Data = err.Error
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var productCart model.ProductCart

	if err := model.DB.Model(&model.ProductCart{}).Where("shopping_cart = ?", id).First(&productCart); err.Error != nil {
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
