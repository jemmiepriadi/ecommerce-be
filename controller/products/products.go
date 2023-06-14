package products

import (
	"context"
	"ecommerce/model"
	"ecommerce/model/objects"
	"ecommerce/utils/paginations"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
		id, err := strconv.Atoi(c.Query("sellerID"))
		if err != nil {
			msg := err
			return nil, msg
		}
		result = queryBuilder.Model(&model.Product{}).Where("seller_id = ?", id).Find(&Products)
		if result.Error != nil {
			msg := result.Error
			return nil, msg
		}

		pagination.TotalData = int64(len(Products))
		totalPages := int(math.Ceil(float64(totalData) / float64(pagination.Size)))
		pagination.TotalPages = totalPages
		pagination.Data = Products

		return pagination, nil
	}

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

	pagination.TotalData = totalData
	totalPages := int(math.Ceil(float64(totalData) / float64(pagination.Size)))
	pagination.TotalPages = totalPages
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

// type Form struct {
// 	File *multipart.FileHeader `form:"file" binding:"required"`
// }

func PostProduct(c *gin.Context) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	file, err := c.FormFile("image")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	client := s3.NewFromConfig(cfg)

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to open the file"})
		return
	}

	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("ecommerce-jemmi"),
		Key:    aws.String(file.Filename),
		Body:   f,
		ACL:    "public-read",
	})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	image := result.Location

	var req model.Product
	res := objects.Response{}

	var isBindFail bool
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to parse request to struct: ", err)

		res.Message = "Failed parsing request"
		isBindFail = true
	}
	if isBindFail {
		form, _ := c.MultipartForm()
		req.Name = form.Value["name"][0]
		req.Description = form.Value["description"][0]
		price, _ := strconv.Atoi(form.Value["price"][0])
		req.Price = price
		sellerID, _ := strconv.Atoi(form.Value["sellerID"][0])
		req.SellerID = sellerID
		req.Image = image
	}

	if result := model.DB.Create(&req); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	}
	res.Message = "Create Successful"
	c.JSON(http.StatusOK, res)
}

func UpdateProduct(c *gin.Context) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	file, err := c.FormFile("image")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	client := s3.NewFromConfig(cfg)

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to open the file"})
		return
	}

	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("ecommerce-jemmi"),
		Key:    aws.String(file.Filename),
		Body:   f,
		ACL:    "public-read",
	})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	image := result.Location

	var product model.Product
	var req model.Product
	res := objects.Response{}
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil && c.Query("id") != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var isBindFail bool
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to parse request to struct: ", err)
		res.Message = "Failed parsing request"
		isBindFail = true

	}

	if isBindFail {
		form, _ := c.MultipartForm()

		req.Name = form.Value["name"][0]
		req.Description = form.Value["description"][0]
		price, _ := strconv.Atoi(form.Value["price"][0])
		req.Price = price
		sellerID, _ := strconv.Atoi(form.Value["sellerID"][0])
		req.SellerID = sellerID
		req.Image = image
	}

	if err := model.DB.Table("products").Where("ID = ?", id).First(&product); err.Error != nil {
		res.Message = "Data not found!"
		res.Data = err.Error
		c.JSON(http.StatusBadRequest, res)
		return
	}
	req.ID = product.ID
	if result := model.DB.Save(&req); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	}
	res.Message = "Update Succes"
	res.Data = req
	c.JSON(http.StatusOK, res)
}

func DeleteProduct(c *gin.Context) {
	// res := deletedata.DeleteItem(&model.Product{}, c)
	id, err := strconv.Atoi(c.Query("id"))
	var product model.Product

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	res := &objects.Response{}

	if err := model.DB.Table("products").Where("ID = ?", id).First(&product); err.Error != nil {
		res.Message = "Data not found!"
		res.Data = err.Error
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if result := model.DB.Delete(&product); result.Error != nil {
		res.Message = "Delete Unsucessful"
		res.Data = err.Error
		c.JSON(http.StatusBadRequest, res)
		return
	}

	//also remove people's carts because the product is deleted
	var productCart []model.ProductCart
	if err := model.DB.Model(&model.ProductCart{}).Where("product_id = ?", id).Find(&productCart); err.Error != nil {
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
	res.Message = "Delete successfull"
	c.JSON(http.StatusOK, res)
}
