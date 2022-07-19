package controllers

import (
	"fmt"
	"time"

	v8 "rogchap.com/v8go"

	"github.com/gin-gonic/gin"
)

func HelloWorld(ctx *gin.Context) {

	jsctx := v8.NewContext()
	jsctx.RunScript(`const main = (a, b) => {
  return { "answer": a + b };
};`, "math.js") // executes a script on the global context
	jsctx.RunScript("const response = main(3, 4)", "main.js") // any functions previously added to the context can be called
	val, err := jsctx.RunScript("response", "value.js")       // return a value in JavaScript back to Go
	// fmt.Sprintf("addition result: %s", val)

	if err != nil {
		e := err.(*v8.JSError)    // JavaScript errors will be returned as the JSError struct
		fmt.Println(e.Message)    // the message of the exception thrown
		fmt.Println(e.Location)   // the filename, line number and the column where the error occured
		fmt.Println(e.StackTrace) // the full stack trace of the error, if available

		fmt.Printf("javascript error: %v", e)        // will format the standard error message
		fmt.Printf("javascript stack trace: %+v", e) // will format the full error stack trace
	}

	ctx.JSON(200, gin.H{
		"response": val,
		"time":     time.Now().Unix(),
	})
}
