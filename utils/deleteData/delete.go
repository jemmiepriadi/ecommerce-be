package deletedata

import (
	"ecommerce/model"
	"ecommerce/model/objects"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteItem(object interface{}, c *gin.Context) *objects.Response {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	var data []interface{}
	res := &objects.Response{}

	if err := model.DB.Model(object).Where("ID = ?", id).First(&data); err != nil {
		res.Message = "Data not found!"
		return res
	}

	model.DB.Model(object).Where("ID = ", id).Delete(&data)
	res.Message = "Delete successfull"
	return res
}
