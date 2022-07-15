package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ryananyangu/gojsrunner/controllers"
)

// var basePath = "/api/v1"

// Routes Function to route mapping
var Routes = map[string]map[string]gin.HandlerFunc{
	"/": {
		"GET": controllers.HelloWorld,
	},
}
