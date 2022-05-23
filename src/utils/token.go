package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func GenerateToken(userId uint64) (string, error) {
	errorENV := godotenv.Load()

	if errorENV != nil {
		fmt.Println("Failed to load env file")
	}

	tokenLife, errTokenLife := strconv.Atoi(os.Getenv("TOKEN_LIFE"))

	if errTokenLife != nil {
		fmt.Println("Can not generate token")
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenLife)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))

}

func TokenValid(ctx *gin.Context) error {
	tokenString := ExtractToken(ctx)

	_, errParse := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if errParse != nil {
		return errParse
	}

	return nil
}

func ExtractToken(ctx *gin.Context) string {
	token := ctx.Query("token")

	if token != "" {
		return token
	}

	bearerToken := ctx.Request.Header.Get("Authorization")

	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func ExtractTokenID(ctx *gin.Context) uint {

	tokenString := ExtractToken(ctx)

	token, errToken := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if errToken != nil {
		ctx.JSON(500, gin.H{
			"message": "Something goes wrong",
		})

		return 0
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		userId, errUserId := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)

		if errUserId != nil {
			ctx.JSON(500, gin.H{
				"message": "Something goes wrong",
			})

			return 0
		}

		return uint(userId)
	}

	return 0
}
