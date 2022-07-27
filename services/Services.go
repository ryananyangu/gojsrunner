package services

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ryananyangu/gojsrunner/utils"
	"rogchap.com/v8go"
)

// NOTE: Handles all the custom implementations
func RunCode() *v8go.Context {
	vm := v8go.NewIsolate()

	global := v8go.NewObjectTemplate(vm)

	// FIXME: Request to be an exact simulation of fetch
	global.Set("request", CustomFetch(vm), v8go.ReadOnly)
	global.Set("btoa", CustomBtoa(vm), v8go.ReadOnly)
	global.Set("log", CustomLog(vm), v8go.ReadOnly)

	return v8go.NewContext(vm, global)
}

// NOTE: Custom Btoa function for js
func CustomBtoa(vm *v8go.Isolate) *v8go.FunctionTemplate {

	btoaFn := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		data := args[0].String()
		response := b64.StdEncoding.EncodeToString([]byte(data))
		val, _ := v8go.NewValue(vm, response)
		return val
	})

	return btoaFn
}

// NOTE: Custom fetch
func CustomFetch(vm *v8go.Isolate) *v8go.FunctionTemplate {
	fetchFn := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		url := args[0].String()
		method := args[1].String()
		payload := args[2].String()
		headers := map[string][]string{}
		if err := json.Unmarshal([]byte(args[3].String()), &headers); err != nil {
			response, _ := json.Marshal(map[string]string{
				"error": err.Error(),
			})
			val, _ := v8go.NewValue(vm, response)
			return val
		}

		goResponse, err := utils.Request(payload, headers, url, method)
		if err != nil {
			response, _ := json.Marshal(map[string]string{
				"error": err.Error(),
			})
			val, _ := v8go.NewValue(vm, response)
			return val
		}
		val, _ := v8go.NewValue(vm, goResponse)
		return val
	})

	return fetchFn
}

func CustomLog(vm *v8go.Isolate) *v8go.FunctionTemplate {

	logFn := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {

		args := info.Args()
		logdata := args[0].String()

		fmt.Println(logdata)
		return nil

	})

	return logFn

}
