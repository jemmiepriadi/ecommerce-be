package products

import (
	"ecommerce/model"
	"ecommerce/model/objects"
	deletedata "ecommerce/utils/deleteData"
	"ecommerce/utils/paginations"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PaginateProduct(product *model.Product, pagination *paginations.Pagination, c *gin.Context) (*paginations.Pagination, error) {
	var totalData int64
	var Products []model.Product

	var result *gorm.DB

	queryBuilder := model.DB.Offset(pagination.GetOffset()).Limit(pagination.GetSize()).Order(pagination.GetSort())

	//find By Id
	if c.Query("id") != "" {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			msg := err
			return nil, msg
		}
		result = queryBuilder.Model(&model.Product{}).Where("ID = ?", id).Find(&Products)
		if result.Error != nil {
			msg := result.Error
			return nil, msg
		}
		model.DB.Model(&Products).Count(&totalData)

		pagination.TotalData = totalData
		totalPages := int(math.Ceil(float64(totalData) / float64(pagination.Size)))
		pagination.TotalPages = totalPages
		pagination.Data = Products
		return pagination, nil
	}

	//findBySeller
	if c.Query("sellerID") != "" {
		id, err := strconv.Atoi(c.Query("SellerID"))
		if err != nil {
			msg := err
			return nil, msg
		}
		result = queryBuilder.Model(&model.Product{}).Where("SellerID = ?", id).Find(&Products)
		if result.Error != nil {
			msg := result.Error
			return nil, msg
		}
		model.DB.Model(&Products).Count(&totalData)

		pagination.TotalData = totalData
		totalPages := int(math.Ceil(float64(totalData) / float64(pagination.Size)))
		pagination.TotalPages = totalPages
		pagination.Data = Products

		return pagination, nil
	}
	model.DB.Model(&Products).Count(&totalData)

	pagination.TotalData = totalData
	totalPages := int(math.Ceil(float64(totalData) / float64(pagination.Size)))
	pagination.TotalPages = totalPages

	//check if searched by name
	if c.Query("name") != "" {
		result = queryBuilder.Model(&model.Product{}).Where("Name LIKE ?", "%"+c.Query("name")+"%").Find(&Products)
	} else { //else find all
		result = queryBuilder.Model(&model.Product{}).Where(product).Find(&Products)
	}
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}
	pagination.Data = Products
	return pagination, nil
}

func GetAllProducts(c *gin.Context) {
	res := &objects.Response{}
	var product model.Product
	paginate := paginations.GeneratePaginationFromRequest(c)
	pagination, err := PaginateProduct(&product, &paginate, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	res.Data = pagination
	res.Message = "Success"
	c.JSON(http.StatusOK, res)
}

func PostProduct(c *gin.Context) {
	var req model.Product
	res := objects.Response{}
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to parse request to struct: ", err)
		res.Code = "02"
		res.Message = "Failed parsing request"

		c.JSON(http.StatusBadRequest, res)
		return
	}
	if result := model.DB.Create(&req); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	}
	res.Message = "Create Successful"
	c.JSON(http.StatusBadRequest, res)
}

func UpdateProduct(c *gin.Context) {
	var product []model.Product
	var req model.Product
	res := objects.Response{}
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	if err := model.DB.Model(&model.Product{}).Where("ID = ?", id).Find(&product); err != nil || len(product) == 0 {
		res.Message = "Data not found!"
		c.JSON(http.StatusBadRequest, res)
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to parse request to struct: ", err)
		res.Code = "02"
		res.Message = "Failed parsing request"

		c.JSON(http.StatusBadRequest, res)
		return
	}
	req.ID = product[0].ID
	if result := model.DB.Save(&req); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	}
	res.Message = "Create Successful"
	c.JSON(http.StatusBadRequest, res)
}

func DeleteProduct(c *gin.Context) {
	res := deletedata.DeleteItem(&model.Product{}, c)
	c.JSON(http.StatusOK, res)
}
