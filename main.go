package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jijiechen/dami-ultra/startup"
)

func main() {
	r := gin.Default()

	err := startup.Run(r)
	if err != nil {
		panic(err)
	}
}
