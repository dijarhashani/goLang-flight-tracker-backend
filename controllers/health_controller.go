package controllers

import (
	"github.com/gin-gonic/gin"
)

func HelathController(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
	})
}
