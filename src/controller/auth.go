package auth

import (
	"task-tracker-api/src/database"
	"task-tracker-api/src/model"
	utils "task-tracker-api/src/utils"

	"github.com/gin-gonic/gin"
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

func ResetPassword(ctx *gin.Context) {
	userId := utils.ExtractTokenID(ctx)
	password := ctx.PostForm("password")

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
