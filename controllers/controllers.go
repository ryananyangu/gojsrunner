package controllers

import (
	"encoding/json"
	"fmt"

	"net/http"

	"crypto/sha256"
	b64 "encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/ryananyangu/gojsrunner/models"
	"github.com/ryananyangu/gojsrunner/services"
	"github.com/ryananyangu/gojsrunner/utils"
	v8 "rogchap.com/v8go"
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
	funcDataInject := fmt.Sprintf(`main(%s, %s,%s,{})`,
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
	response := map[string]interface{}{}

	if err := ctx.ShouldBindJSON(&response); err != nil {
		utils.Log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsctx := services.RunCode()
	defer jsctx.Isolate().Dispose()
	scriptFile := fmt.Sprintf("res_%s.js", service)
	scriptContent, err := utils.ReadFile("wrapperscripts/" + scriptFile)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	jsctx.RunScript(scriptContent, scriptFile)
	body, err := json.Marshal(response)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

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
	year := "2022"
	month := "08"
	day := "02"
	hrs := "06"
	mins := "35"
	secs := "56"
	noncetime := year + month + day + hrs + mins + secs //`20220802063556`
	str := []byte(noncetime)

	nonce := b64.StdEncoding.EncodeToString(str)
	rawStr := fmt.Sprintf("%s%s-%s-%sT%s:%s:%sZSdpita@20!@", nonce, year, month, day, hrs, mins, secs)
	SHA256 := sha256.Sum256([]byte(rawStr))
	base64 := b64.StdEncoding.EncodeToString(SHA256[:])

	// rawStr := nonce + timespan + AppSecret
	ctx2.JSON(http.StatusOK, map[string]interface{}{
		"nonce":  nonce,
		"rawStr": rawStr,
		"base64": base64,
	})
}
