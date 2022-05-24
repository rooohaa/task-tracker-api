package controller

import (
	"net/http"
	"task-tracker-api/src/database"
	"task-tracker-api/src/model"
	utils "task-tracker-api/src/utils"

	"github.com/gin-gonic/gin"
)

func GetTasks(ctx *gin.Context) {
	var tasks []model.Task
	database.GetDB.Find(&tasks)
	ctx.JSON(http.StatusOK, tasks)
}

func GetTaskById(ctx *gin.Context) {
	id := ctx.Param("id")
	var task model.Task

 	err := database.GetDB.Where("id = ?", id).First(&task).Error; 
	
 	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message" : "Task with id: " + id + " not found",
		})
    return
  }

  ctx.JSON(http.StatusOK, task)
}

func CreateTask(ctx *gin.Context) {
	userId := utils.ExtractTokenID(ctx)
	var task model.Task
	task.UserId = uint64(userId)

  err := ctx.BindJSON(&task)
	
	if err != nil {
    ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
      "error": err.Error(),
    })
    return
  }
	
  database.GetDB.Create(&task)
  ctx.JSON(http.StatusOK, &task)
}

func UpdateTask(ctx *gin.Context) {
  id := ctx.Param("id")
  var task model.Task

  err := database.GetDB.Where("id = ?", id).First(&task).Error; 
	
	if err != nil {
    ctx.AbortWithStatus(http.StatusNotFound)
    return
  }

  ctx.BindJSON(&task)
  database.GetDB.Save(&task)
  ctx.JSON(http.StatusOK, &task)
}

func DeleteTask(ctx *gin.Context) {
  id := ctx.Param("id")
  var task model.Task

  if err := database.GetDB.Where("id = ?", id).First(&task).Error; err != nil {
    ctx.AbortWithStatus(http.StatusNotFound)
    return
  }

	database.GetDB.Delete(&task)
	ctx.JSON(http.StatusOK, gin.H{"message": "task with id " + id + " deleted"})
}

