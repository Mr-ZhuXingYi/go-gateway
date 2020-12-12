package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()
	r.Handle("GET", "/user", func(context *gin.Context) {
		context.JSON(200, "hello world")
	})
	r.Run(":8080")
}
