package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ryananyangu/gojsrunner/utils"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(utils.LogrusLogger())
	r.Use(gin.Recovery())

	for path, handlers := range Routes {
		for method, handler := range handlers {
			switch method {
			case "GET":
				r.GET(path, handler)

			case "POST":
				r.POST(path, handler)

			case "PUT":
				r.PUT(path, handler)

			case "PATCH":
				r.PATCH(path, handler)

			case "DELETE":
				r.DELETE(path, handler)
			}
		}
	}
	return r
}

func main() {

	gin_instance := SetupRouter()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	gin_instance.Run(fmt.Sprintf(":%s", port))

}
