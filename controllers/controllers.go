package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/ryananyangu/gojsrunner/models"
	"github.com/ryananyangu/gojsrunner/services"
	"github.com/ryananyangu/gojsrunner/utils"
)

func RequestTransformation(request *models.Request) {

	jsctx := services.RunCode()
	defer jsctx.Isolate().Dispose()

	scriptfile := fmt.Sprintf("req_%s.js", request.ClientInfo.ServiceCode)

	scriptContent, err := utils.ReadFile("wrapperscripts/" + scriptfile)

	if err != nil {
		utils.Log.Error(err)
		return

	}

	jsctx.RunScript(scriptContent, scriptfile)

	constants, err := json.Marshal(request.ClientInfo.Statics)
	if err != nil {
		utils.Log.Error(err)
		return
	}
	headers, err := json.Marshal(request.ClientInfo.Headers)
	if err != nil {
		utils.Log.Error(err)
		return
	}
	payload, err := json.Marshal(request.Transaction)

	if err != nil {
		utils.Log.Error(err)
		return
	}

	//FIXME: Add Settings Build script call by inserting variable
	funcDataInject := fmt.Sprintf(`main(%s, %s,%s,{})`,
		string(payload[:]),
		string(headers[:]),
		string(constants[:]))

	//FIXME: Capture return of the main js Execute main function
	// Repackage and publish to ask queue
	val, err := jsctx.RunScript(funcDataInject, scriptfile)
	utils.Log.Info(val)
	if err != nil {
		utils.Log.Error(err)
		return
	}
}
