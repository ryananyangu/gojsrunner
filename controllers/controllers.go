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
	defer jsctx.Isolate().Dispose()

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
	scriptFile := fmt.Sprintf("res_%s.js", service)
	scriptContent, err := utils.ReadFile(scriptFile)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	jsctx.RunScript(scriptContent, scriptFile)

	funcDataInject := fmt.Sprintf(`main(%s)`,
		string(body[:]))

	// Execute main function
	val, err := jsctx.RunScript(funcDataInject, scriptFile)
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
	val, err1 := ctx.RunScript(`send('','{"Content-Type" : ["application/json"]}','https://dummy.restapiexample.com/api/v1/employee/1','GET')`, "print.js")

	if err1 != nil {
		ctx2.JSON(http.StatusBadRequest, err1)
		return
	}
	res, _ := val.Object().MarshalJSON()
	ctx2.JSON(http.StatusOK, string(res))
}
