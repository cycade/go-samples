package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	service := initApp()

	r.GET("/restart", func(c *gin.Context) {
		id := c.Query("id")
		resp := service.Restart(id)
		c.JSON(200, resp)
	})
	r.Run()
}
