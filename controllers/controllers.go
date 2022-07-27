package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"

	v8 "rogchap.com/v8go"

	"github.com/gin-gonic/gin"
	"github.com/ryananyangu/gojsrunner/models"
	"github.com/ryananyangu/gojsrunner/services"
	"github.com/ryananyangu/gojsrunner/utils"
)

func RequestTransformation(ctx *gin.Context) {

	service := ctx.Param("service")

	request := models.Request{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.Log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsctx := services.RunCode()

	scriptfile := fmt.Sprintf("req_%s.js", service)

	scriptContent, err := utils.ReadFile("wrapperscripts/" + scriptfile)

	if err != nil {
		utils.Log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	jsctx.RunScript(scriptContent, scriptfile)

	constants, err := json.Marshal(request.Constants)
	if err != nil {
		utils.Log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	headers, err := json.Marshal(request.Headers)
	if err != nil {
		utils.Log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload, err := json.Marshal(request.Payload)

	if err != nil {
		utils.Log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build script call by inserting variable
	funcDataInject := fmt.Sprintf(`main(%s, %s,%s)`,
		string(payload[:]),
		string(headers[:]),
		string(constants[:]))

	// Execute main function
	val, err := jsctx.RunScript(funcDataInject, scriptfile)
	if err != nil {
		utils.Log.Error(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	// FIXME: validate in the recieving application
	ctx.JSON(http.StatusOK, val)
}

func ResponseTransformation(ctx *gin.Context) {

	service := ctx.Param("service")

	ioRead, err := ctx.Request.GetBody()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, err := ioutil.ReadAll(ioRead)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsctx := v8.NewContext()

	scriptContent, err := utils.ReadFile(fmt.Sprintf("res_%s.js", service))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	jsctx.RunScript(scriptContent, "main.js")

	funcDataInject := fmt.Sprintf(`const response = main(%s)`,
		string(body[:]))

	// Execute main function
	_, err = jsctx.RunScript(funcDataInject, "main.js")
	if err != nil {
		e := err.(*v8.JSError)
		ctx.JSON(http.StatusBadRequest, models.RequestBuilt{
			Error: models.Error{
				Message:    e.Message,
				Location:   e.Location,
				StackTrace: e.StackTrace,
			},
		})
		return
	}

	// Capture result from the function ran
	val, err := jsctx.RunScript("response", "value.js")
	if err != nil {
		e := err.(*v8.JSError)
		ctx.JSON(http.StatusBadRequest, models.RequestBuilt{
			Error: models.Error{
				Message:    e.Message,
				Location:   e.Location,
				StackTrace: e.StackTrace,
			},
		})
		return
	}

	// FIXME: validate in the recieving application
	ctx.JSON(http.StatusOK, val)
}

func TestApi(ctx2 *gin.Context) {
	ctx := services.RunCode()
	defer ctx.Isolate().Dispose()
	_, err1 := ctx.RunScript("log('data')", "print.js")

	// _, err1 := ctx.RunScript("const main = (data) => { log(data)}", "print.js")

	if err1 != nil {
		ctx2.JSON(http.StatusBadRequest, err1)
		return
	}

	// val, err := ctx.RunScript(fmt.Sprintf(`main('%s')`, "My data is this"), "response.js")
	// if err != nil {
	// 	ctx2.JSON(http.StatusBadRequest, err)
	// 	return
	// }
	// utils.Log.Info(val)
	ctx2.JSON(http.StatusOK, "val")
}

// {
//     "Message": "ReferenceError: log is not defined",
//     "Location": "print.js:1:26",
//     "StackTrace": "ReferenceError: log is not defined\n    at main (print.js:1:26)\n    at response.js:1:1"
// }
