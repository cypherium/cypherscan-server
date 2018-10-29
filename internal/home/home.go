package home

import (
	"github.com/gin-gonic/gin"
)

// controller of "/home"
func GetHome(c *gin.Context) {
	c.JSON(200, gin.H{"message": "good"})
}
