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
	if err := model.DB.Create(&req); err.Error != nil {
		res.Message = "Username already exists"
		res.Data = nil
		c.JSON(http.StatusBadRequest, res)
		return
	}
	if req.UserType != "" {
		if req.UserType == "seller" {
			var seller model.Seller
			seller.Name = req.Name
			seller.AccountID = req.ID
			model.DB.Create(&seller)
		} else if req.UserType == "consumer" {
			var consumer model.Consumer
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
	response := FindUser(user.Username, user.Password)
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

func FindUser(username, password string) map[string]interface{} {
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

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = user
	return resp
}
