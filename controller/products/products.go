package products

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

func PaginateProduct(product *model.Product, pagination *paginations.Pagination, c *gin.Context) (*paginations.Pagination, error) {
	var totalData int64
	var Products []model.Product

	var result *gorm.DB

	//find By Id
	if c.Query("id") != "" {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			msg := err
			return nil, msg
		}
		result = model.DB.Model(&model.Product{}).Where("ID = ?", id).Find(&Products)
		if result.Error != nil {
			msg := result.Error
			return nil, msg
		}
		model.DB.Model(&result).Count(&totalData)

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
		result = model.DB.Model(&model.Product{}).Where("SellerID = ?", id).Find(&Products)
		if result.Error != nil {
			msg := result.Error
			return nil, msg
		}
		model.DB.Model(&result).Count(&totalData)

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

	queryBuilder := model.DB.Offset(pagination.GetOffset()).Limit(pagination.GetSize()).Order(pagination.GetSort())

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

func UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	if err := model.DB.Model(object).Where("ID = ?", id).Find(&data); err != nil {
		res.Message = "Data not found!"
		return res
	}
}

func DeleteProduct(c *gin.Context) {
	res := deletedata.DeleteItem(&model.Product{}, c)
	c.JSON(http.StatusOK, res)
}
