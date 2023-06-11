package paginations

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Size       int         `json:"size,omitempty" query:"size"`
	Page       int         `json:"page,omitempty" query:"page"`
	Sort       string      `json:"sort,omitempty" query:"sort"`
	TotalData  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

func GeneratePaginationFromRequest(c *gin.Context) Pagination {
	// Initializing default
	//	var mode string
	size := 10
	page := 1
	sort := "id asc"
	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "size":
			size, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "sort":
			sort = queryValue
			break

		}
	}
	return Pagination{
		Size: size,
		Page: page,
		Sort: sort,
	}

}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetSize()
}
func (p *Pagination) GetSize() int {
	if p.Size == 0 {
		p.Size = 10
	}
	return p.Size
}
func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}
func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}
