package main

import (
	"github.com/fayca121/cardchecker/src/controllers"
	"github.com/gin-gonic/gin"
)

func main() {

	route := gin.Default()

	route.POST("/check", controllers.CheckCardNumber)

	route.Run()

}
