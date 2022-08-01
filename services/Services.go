package services

import (
	"crypto/sha256"
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
	global.Set("send", CustomFetch(vm), v8go.ReadOnly)
	global.Set("btoa", CustomBtoa(vm), v8go.ReadOnly)
	global.Set("log", CustomLog(vm), v8go.ReadOnly)
	global.Set("SHA256", CustomSHA256(vm), v8go.ReadOnly)

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
		payload := args[0].String()
		headers := map[string][]string{}
		url := args[1].String()
		method := args[3].String()

		if err := json.Unmarshal([]byte(args[2].String()), &headers); err != nil {
			response, err2 := json.Marshal(map[string]string{
				"error": err.Error(),
			})
			val, err1 := v8go.NewValue(vm, []byte(response))
			utils.Log.Error(err)
			utils.Log.Error(err1)
			utils.Log.Error(err2)
			return val
		}

		goResponse, err := utils.Request(payload, headers, url, method)

		if err != nil {
			utils.Log.Error(err)
			return nil
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

		fmt.Println("[JS_SCRIPT]\t" + logdata)
		return nil

	})

	return logFn

}

func CustomSHA256(vm *v8go.Isolate) *v8go.FunctionTemplate {

	sha256Fn := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		str2encode := args[0].String()
		hash := sha256.Sum256([]byte(str2encode))
		val, _ := v8go.NewValue(vm, string(hash[:]))
		return val

	})

	return sha256Fn

}
