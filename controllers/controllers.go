package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/ryananyangu/gojsrunner/models"
	"github.com/ryananyangu/gojsrunner/services"
	"github.com/ryananyangu/gojsrunner/utils"
)

func RequestTransformation(request *models.Request) {

	// Initiate js vm
	jsctx := services.RunCode()
	defer jsctx.Isolate().Dispose()

	// Get payload generation script from storage based on the service code
	scriptfile := fmt.Sprintf("req_%s.js", request.ClientInfo.ServiceCode)

	// read the content of the script
	scriptContent, err := utils.ReadFile("wrapperscripts/" + scriptfile)
	if err != nil {
		utils.Log.Error(err)
		return

	}

	// initiate the content within the js vm
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
	settings, err := json.Marshal(request.ClientInfo.Settings)

	if err != nil {
		utils.Log.Error(err)
		return
	}

	// Build string to js function call
	funcDataInject := fmt.Sprintf(`main(%s, %s,%s,%s)`,
		string(payload[:]),
		string(headers[:]),
		string(constants[:]), string(settings[:]))

	// Run the js function call from golang
	val, err := jsctx.RunScript(funcDataInject, scriptfile)
	if err != nil {
		utils.Log.Error(err)
		return
	}

	// json encode the js script response
	reqScriptRes, err := val.MarshalJSON()
	if err != nil {
		utils.Log.Error(err)
		return
	}

	// cast the response to a valid struct
	builtRequest := models.RequestBuilt{}
	err = json.Unmarshal(reqScriptRes, &builtRequest)
	if err != nil {
		utils.Log.Error(err)
		return
	}

	// Send the main service request
	serviceResponse, err := utils.Request(builtRequest.Payload,
		builtRequest.Headers,
		request.ClientInfo.ServiceURL,
		request.ClientInfo.HTTPMethod)
	if err != nil {
		utils.Log.Error(err)
		return
	}

	// FIXME: Handle xml based transformation to map here
	// For json transform to map
	response := map[string]interface{}{}
	if err := json.Unmarshal([]byte(serviceResponse), &response); err != nil {
		utils.Log.Error(err)
		return
	}

	//

	// Prepare response processing script
	ResScriptFile := fmt.Sprintf("res_%s.js", request.ClientInfo.ServiceCode)
	resScriptContent, err := utils.ReadFile("wrapperscripts/" + ResScriptFile)
	if err != nil {
		utils.Log.Error(err)
		return

	}

	jsctx.RunScript(resScriptContent, ResScriptFile)

	// Pass the services response to function FIXME: Only works for json for now
	resInjectScript := fmt.Sprintf(`main(%s)`, serviceResponse)

	// Execute main function
	resVal, err := jsctx.RunScript(resInjectScript, ResScriptFile)
	if err != nil {
		utils.Log.Error(err)
		return
	}

	// get json of script response
	finalres, err := resVal.MarshalJSON()
	if err != nil {
		utils.Log.Error(err)
		return
	}

	// Cast response to struct to make sure to malformation of the response [Validation]
	requestresp := models.Response{}
	err = json.Unmarshal(finalres, &requestresp)
	if err != nil {
		utils.Log.Error(err)
		return
	}

	// Publish to ack Queue
	services.PublishPaymentAck(finalres)

}
