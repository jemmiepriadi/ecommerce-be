package auth

import (
	"ecommerce/model"
	"ecommerce/model/objects"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var res = objects.Response{
	Code:    "00",
	Message: "Success",
}

type JWTClaim struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

func Me(c *gin.Context) {
	var signedToken objects.JWT
	signedToken.JWT = c.Query("JWT")
	token, err := jwt.ParseWithClaims(
		signedToken.JWT,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("secretjemmi"), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "couldn't parse claims"})
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "token expired"})
		return
	}

	user := &model.Account{}
	if err := model.DB.Where("username = ?", claims.Username).First(user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username not found"})
		return
	}
	var seller *objects.Seller
	var consumer *objects.Consumer
	var userData objects.UserData

	if user.UserType == "seller" {
		if err := model.DB.Model(&model.Seller{}).Where("account_id = ?", user.ID).First(&seller).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Username not found"})
			return
		}
		userData.Seller = seller
	} else if user.UserType == "consumer" {
		if err := model.DB.Model(&model.Consumer{}).Where("account_id = ?", user.ID).First(&consumer).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Username not found"})
			return
		}
		userData.Consumer = consumer
	}
	userData.Address = user.Address
	userData.UserID = user.ID
	userData.PhoneNumber = user.PhoneNumber
	userData.Name = claims.Name
	userData.Username = claims.Username
	userData.UserType = user.UserType
	res.Data = userData
	c.JSON(http.StatusOK, res)
}

func PostRegister(c *gin.Context) {
	var req model.Account
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to parse request to profile struct: ", err)
		res.Code = "02"
		res.Message = "Failed parsing request"

		c.JSON(http.StatusBadRequest, res)
		return
	}

	if req.Username == "" || req.Password == "" {
		res.Code = ""
		res.Message = "Password or Username cannot be empty"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		res.Code = "02"
		res.Message = "Password Encryption  failed"
		c.JSON(http.StatusBadRequest, res)
		return
	}
	req.Password = string(pass)

	res.Data = req
	if err := model.DB.Save(&req); err.Error != nil {
		res.Message = "Username already exists"
		res.Data = nil
		c.JSON(http.StatusBadRequest, res)
		return
	}
	AccountID := req.ID

	res.Message = "Account successfully created"
	if req.UserType != "" {
		if req.UserType == "seller" {
			var seller model.Seller
			result := model.DB.Model(&model.Seller{}).Where("account_id = ?", AccountID).Find(&seller)
			if result.Error != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			seller.Name = req.Name
			seller.AccountID = req.ID
			model.DB.Create(&seller)
		} else if req.UserType == "consumer" {
			var consumer model.Consumer
			result := model.DB.Model(&model.Consumer{}).Where("account_id = ?", AccountID).Find(&consumer)
			if result.Error != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			consumer.Name = req.Name
			consumer.AccountID = req.ID
			model.DB.Create(&consumer)
		}
	}
	c.JSON(http.StatusOK, res)
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "request does not contain an access token"})
			ctx.Abort()
			return
		}
		if err := ValidateToken(header); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func Login(c *gin.Context) {
	var user model.Account
	if err := c.ShouldBindJSON(&user); err != nil {
		res.Code = "02"
		res.Message = "Failed parsing request"

		c.JSON(http.StatusBadRequest, res)
		return
	}
	response := FindUser(user.Username, user.Password, c)

	c.JSON(http.StatusOK, response)
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("secretjemmi"), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}

func FindUser(username, password string, c *gin.Context) map[string]interface{} {
	user := &model.Account{}
	if err := model.DB.Where("username = ?", username).First(user).Error; err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Username not found"}
		return resp
	}
	expiresAt := time.Now().Add(1 + time.Hour).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
	}

	claim := &JWTClaim{
		Username: username,
		Name:     user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, error := token.SignedString([]byte("secretjemmi"))
	if error != nil {
		fmt.Println(error)
	}

	var seller *objects.Seller
	var consumer *objects.Consumer
	var userData objects.UserData

	if user.UserType == "seller" {
		if err := model.DB.Model(&model.Seller{}).Where("account_id = ?", user.ID).First(&seller).Error; err != nil {
			var resp = map[string]interface{}{"status": false, "message": "Username not found"}
			return resp
		}
		userData.Seller = seller
	} else if user.UserType == "consumer" {
		if err := model.DB.Model(&model.Consumer{}).Where("account_id = ?", user.ID).First(&consumer).Error; err != nil {
			var resp = map[string]interface{}{"status": false, "message": "Username not found"}
			return resp
		}
		userData.Consumer = consumer
	}
	userData.Address = user.Address
	userData.UserID = user.ID
	userData.PhoneNumber = user.PhoneNumber
	userData.Name = claim.Name
	userData.Username = claim.Username
	userData.UserType = user.UserType
	res.Data = userData

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = res.Data
	c.SetCookie("auth_token", tokenString, 3600, "/", "localhost", false, true)

	return resp
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
