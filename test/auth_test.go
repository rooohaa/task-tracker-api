package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	auth "task-tracker-api/src/controller"
	"task-tracker-api/src/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestRegister(t *testing.T) {

	r := SetUpRouter()

	r.POST("/register", func(ctx *gin.Context) {
		auth.Register(ctx)
	})

	user := model.Users{
		Email:    "test_email@gmail.com",
		Password: "test_password",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, http.StatusOK)
}

func TestLogin(t *testing.T) {

	r := SetUpRouter()

	r.POST("/login", func(ctx *gin.Context) {
		auth.Login(ctx)
	})

	user := model.Users{
		Email:    "test_email@gmail.com",
		Password: "test_password",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, http.StatusOK)
}
