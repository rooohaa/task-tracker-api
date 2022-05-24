package controller

import (
	"fmt"
	"net/smtp"
	"os"
	"task-tracker-api/src/database"
	"task-tracker-api/src/model"
	utils "task-tracker-api/src/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(ctx *gin.Context) {
	var input RegisterInput

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Email and password are required fields",
		})

		return
	}

	user := model.Users{}

	user.Email = input.Email

	hashedPassword := utils.HashPassword(ctx, input.Password)

	if len(hashedPassword) == 0 {
		return
	}

	user.Password = hashedPassword

	errCreate := database.GetDB.Create(&user).Error

	if errCreate != nil {
		ctx.JSON(400, gin.H{
			"message": "User is already been created",
		})

		return
	}

	token, errToken := utils.GenerateToken(user.ID)

	if errToken != nil {
		ctx.JSON(500, gin.H{
			"message": "Something goes wrong",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "User successfully registered",
		"token":   token,
	})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(ctx *gin.Context) {
	var input LoginInput

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "User is already been created",
		})

		return
	}

	user := model.Users{}

	errCheck := database.GetDB.Model(&user).Where("email = ?", input.Email).Take(&user).Error

	if errCheck != nil {
		ctx.JSON(400, gin.H{
			"message": "User is not registered",
		})

		return
	}

	errVerify := utils.VerifyPassword(input.Password, user.Password)

	if errVerify != nil && errVerify == bcrypt.ErrMismatchedHashAndPassword {
		ctx.JSON(401, gin.H{
			"message": "Password is not correct",
		})

		return
	}

	token, errToken := utils.GenerateToken(user.ID)

	if errToken != nil {
		ctx.JSON(500, gin.H{
			"message": "Something goes wrong",
		})

		return
	}

	ctx.JSON(200, gin.H{
		"message": "User successfully logined",
		"token":   token,
	})
}

type ChangePasswordInput struct {
	Password string `json:"password" binding:"required"`
}

func ChangePassword(ctx *gin.Context) {
	var input ChangePasswordInput

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "User password should be sent",
		})

		return
	}

	userId := utils.ExtractTokenID(ctx)

	user := model.Users{}

	hashedPassword := utils.HashPassword(ctx, input.Password)

	errReset := database.GetDB.Model(&user).Where("id = ?", userId).Update("password", hashedPassword).Error

	if errReset != nil {
		ctx.JSON(400, gin.H{
			"message": "Cant update password",
		})

		return
	}

	ctx.JSON(200, gin.H{
		"message": "User password successfully updated",
	})
}

var verification = make(map[string]string)
var timer *time.Timer

type SendCodeInput struct {
	Email string `json:"email" binding:"required"`
}

func SendCode(ctx *gin.Context) {
	errorENV := godotenv.Load()

	if errorENV != nil {
		fmt.Println("Failed to load env file")
	}

	var input SendCodeInput

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "User email should be sent",
		})

		return
	}

	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	to := []string{
		input.Email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	otpCode, errOtp := utils.GenerateOTP(6)

	if errOtp != nil {
		ctx.JSON(500, gin.H{
			"message": "Something goes wrong",
		})
	}

	verification[input.Email] = otpCode
	StartDeleteOTP(input.Email)

	message := []byte("This is a code for verification of user " + otpCode)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	errMail := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	if errMail != nil {
		ctx.JSON(400, gin.H{
			"message": errMail,
		})

		return
	}

	ctx.JSON(400, gin.H{
		"message": "Email has been sent successfully",
	})
}

type VerifyCodeInput struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func VerifyCode(ctx *gin.Context) {
	var input VerifyCodeInput

	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "User email and code should be sent",
		})

		return
	}

	if verification[input.Email] == input.Code {
		user := model.Users{}
		errVerify := database.GetDB.Model(&user).Where("email = ?", input.Email).Update("verified", true).Error

		if errVerify != nil {
			ctx.JSON(400, gin.H{
				"message": "Can not verify user",
			})

			return
		}

		delete(verification, input.Email)
		defer timer.Stop()

		ctx.JSON(200, gin.H{
			"message": "User successfully verified",
		})

		return
	}

	ctx.JSON(401, gin.H{
		"message": "Code is not correct",
	})
}

func StartDeleteOTP(email string) {
	timer = time.AfterFunc(30*time.Second, func() {
		delete(verification, email)
	})

	defer timer.Stop()
}
