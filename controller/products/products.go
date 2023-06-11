package products

import (
	"ecommerce/model"
	"ecommerce/model/objects"
	"ecommerce/utils/paginations"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PaginateProduct(product *model.Product, pagination *paginations.Pagination, c *gin.Context) (*paginations.Pagination, error) {
	var totalRows int64
	var Products []model.Product
	model.DB.Model(&Products).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Size)))
	pagination.TotalPages = totalPages
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
		pagination.Data = Products
		return pagination, nil
	}

	//findBySeller
	if c.Query("sellerID") != "" {
		result = model.DB.Model(&model.Product{}).Where("SellerID = ?", c.Query("sellerID")).Find(&Products)
		if result.Error != nil {
			msg := result.Error
			return nil, msg
		}
		pagination.Data = Products
		return pagination, nil
	}

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
