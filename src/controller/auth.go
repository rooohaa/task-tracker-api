package auth

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

func Register(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	if len(email) == 0 || len(password) == 0 {
		ctx.JSON(400, gin.H{
			"message": "Email and password are required fields",
		})

		return
	}

	user := model.Users{}

	user.Email = email

	hashedPassword := utils.HashPassword(ctx, password)

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

func Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	if len(email) == 0 || len(password) == 0 {
		ctx.JSON(400, gin.H{
			"message": "Email and password are required fields",
		})

		return
	}

	user := model.Users{}

	errCheck := database.GetDB.Model(&user).Where("email = ?", email).Take(&user).Error

	if errCheck != nil {
		ctx.JSON(400, gin.H{
			"message": "User is not registered",
		})

		return
	}

	errVerify := utils.VerifyPassword(password, user.Password)

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

func ChangePassword(ctx *gin.Context) {
	userId := utils.ExtractTokenID(ctx)
	password := ctx.PostForm("password")

	if len(password) == 0 {
		ctx.JSON(400, gin.H{
			"message": "Password are required fields",
		})

		return
	}

	user := model.Users{}

	hashedPassword := utils.HashPassword(ctx, password)

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

func SendCode(ctx *gin.Context) {
	errorENV := godotenv.Load()

	if errorENV != nil {
		fmt.Println("Failed to load env file")
	}

	email := ctx.PostForm("email")

	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	otpCode, errOtp := utils.GenerateOTP(6)

	if errOtp != nil {
		ctx.JSON(500, gin.H{
			"message": "Something goes wrong",
		})
	}

	verification[email] = otpCode
	StartDeleteOTP(email)

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

func VerifyCode(ctx *gin.Context) {
	email := ctx.PostForm("email")
	code := ctx.PostForm("code")

	if len(code) == 0 || len(email) == 0 {
		ctx.JSON(400, gin.H{
			"message": "Email and code fields are required",
		})

		return
	}

	if verification[email] == code {
		user := model.Users{}
		errVerify := database.GetDB.Model(&user).Where("email = ?", email).Update("verified", true).Error

		if errVerify != nil {
			ctx.JSON(400, gin.H{
				"message": "Can not verify user",
			})

			return
		}

		delete(verification, email)
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
