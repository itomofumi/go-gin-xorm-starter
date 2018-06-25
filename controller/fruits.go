package controller

import (
	"net/http"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// GetFruits はフルーツ一覧取得
func GetFruits(c *gin.Context) {
	registry := c.MustGet(service.RegistryKey).(service.RegistryInterface)
	Fruits := registry.NewFruits()
	list, err := Fruits.GetAll()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetFruitByID はフルーツを取得します
func GetFruitByID(c *gin.Context) {
	noticeID := c.MustGet("fruit-id").(int64)
	registry := c.MustGet(service.RegistryKey).(service.RegistryInterface)
	Fruitservice := registry.NewFruits()
	notice, err := Fruitservice.GetByID(noticeID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}
	c.JSON(http.StatusOK, notice)
}

// PostFruit はフルーツを登録します
func PostFruit(c *gin.Context) {
	registry := c.MustGet(service.RegistryKey).(service.RegistryInterface)
	Fruitservice := registry.NewFruits()

	noticeBody := model.FruitBody{}
	if err := c.ShouldBindWith(&noticeBody, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}

	created, err := Fruitservice.Create(&noticeBody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, created)
}

// PutFruit はフルーツを更新します
func PutFruit(c *gin.Context) {
	noticeID := c.MustGet("fruit-id").(int64)
	registry := c.MustGet(service.RegistryKey).(service.RegistryInterface)

	Fruitservice := registry.NewFruits()

	noticeBody := model.FruitBody{}
	if err := c.ShouldBindWith(&noticeBody, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}

	updated, err := Fruitservice.Update(noticeID, &noticeBody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteFruit はフルーツを削除します
func DeleteFruit(c *gin.Context) {
	noticeID := c.MustGet("fruit-id").(int64)
	registry := c.MustGet(service.RegistryKey).(service.RegistryInterface)

	Fruitservice := registry.NewFruits()
	err := Fruitservice.Delete(noticeID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
