package middleware

import (
	"authentication/initializers"
	"authentication/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	//Get the coockie off req
	tokenString, err := c.Cookie("Authorization")

	if err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	//Decode/validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface {}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Method)
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		//Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		//Atach to req
		c.Set("user", user)
		//Continue
		c.Next()

		fmt.Println(claims["foo"], claims["nbf"])
	}else{
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	
}