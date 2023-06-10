package products

import (
	"ecommerce/model"
	"ecommerce/model/objects"
	"ecommerce/utils/paginations"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PaginateProduct(product *model.Product, pagination *paginations.Pagination) (*paginations.Pagination, error) {
	var totalRows int64
	var Products []model.Product
	model.DB.Model(&Products).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Size)))
	pagination.TotalPages = totalPages

	queryBuilder := model.DB.Offset(pagination.GetOffset()).Limit(pagination.GetSize()).Order(pagination.GetSort())
	result := queryBuilder.Model(&model.Product{}).Where(product).Find(&Products)
	pagination.Data = result
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}
	return pagination, nil
}

func GetAllProducts(c *gin.Context) {
	res := &objects.Response{}
	var product model.Product
	paginate := paginations.GeneratePaginationFromRequest(c)
	pagination, err := PaginateProduct(&product, &paginate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	res.Data = pagination
	res.Message = "Success"
	c.JSON(http.StatusOK, paginations.GeneratePaginationFromRequest(c))
}
