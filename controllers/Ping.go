package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
)

func HelloWorld(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Hello World!!",
		"time":    time.Now().Unix(),
	})
}
