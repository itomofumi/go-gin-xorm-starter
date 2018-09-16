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
	fruitsService := registry.NewFruits()
	list, err := fruitsService.GetAll()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetFruitByID はフルーツを取得します
func GetFruitByID(c *gin.Context) {
	fruitID := c.MustGet("fruit-id").(uint64)
	registry := c.MustGet(service.RegistryKey).(service.RegistryInterface)
	fruitsService := registry.NewFruits()
	fruit, err := fruitsService.GetByID(fruitID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}
	c.JSON(http.StatusOK, fruit)
}

// PostFruit はフルーツを登録します
func PostFruit(c *gin.Context) {
	registry := c.MustGet(service.RegistryKey).(service.RegistryInterface)
	fruitsService := registry.NewFruits()

	fruitBody := model.FruitBody{}
	if err := c.ShouldBindWith(&fruitBody, binding.JSON); err != nil || !fruitBody.IsValid() {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}

	created, err := fruitsService.Create(&fruitBody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, created)
}

// PutFruit はフルーツを更新します
func PutFruit(c *gin.Context) {
	fruitID := c.MustGet("fruit-id").(uint64)
	registry := c.MustGet(service.RegistryKey).(service.RegistryInterface)

	fruitsService := registry.NewFruits()

	fruitBody := model.FruitBody{}
	c.BindWith(&fruitBody, binding.JSON)

	updated, err := fruitsService.Update(fruitID, &fruitBody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteFruit はフルーツを削除します
func DeleteFruit(c *gin.Context) {
	fruitID := c.MustGet("fruit-id").(uint64)
	registry := c.MustGet(service.RegistryKey).(service.RegistryInterface)

	fruitsService := registry.NewFruits()
	err := fruitsService.Delete(fruitID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
