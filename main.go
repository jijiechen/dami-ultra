package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yuchanns/kong-exercise-microservices/startup"
)

func main() {
	r := gin.Default()

	err := startup.Run(r)
	if err != nil {
		panic(err)
	}
}
