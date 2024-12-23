package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	v1 "github.com/0x41gawor/lupus/api/v1"
	util "github.com/0x41gawor/lupus/internal/util"
)

func sendToDestination(input interface{}, dest v1.Destination) (interface{}, error) {
	switch dest.Type {
	case "http":
		res, err := sendToHTTP(dest.HTTP.Path, dest.HTTP.Method, input)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "opa":
		res, err := sendToOpa(dest.Opa.Path, input)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "gofunc":
		res, err := sendToGoFunc(dest.GoFunc.Name, input)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, fmt.Errorf("no such destination type implemented yet")
	}
}

func sendToOpa(path string, reqBody interface{}) (interface{}, error) {
	wrappedBody := map[string]interface{}{
		"input": reqBody,
	}

	// Call sendToHTTP to get the response
	res, err := sendToHTTP(path, "POST", wrappedBody)
	if err != nil {
		return nil, err
	}
	resMap, err := util.InterfaceToMap(res)
	if err != nil {
		return nil, fmt.Errorf("unexpected response format, not a map")
	}
	// Return only the content of "result"
	if result, ok := resMap["result"]; ok {
		return result, nil
	}
	return nil, fmt.Errorf("no 'result' field in response")
}

func sendToHTTP(path string, method string, body interface{}) (interface{}, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewBuffer(bodyBytes)
	httpReq, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		println("here1\n")
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		println("here2\n")
		return nil, fmt.Errorf("non-ok HTTP Status")
	}

	var res interface{}
	if err := json.Unmarshal(respBody, &res); err != nil {
		println("here3\n")
		return nil, err
	}
	return res, nil
}

func sendToGoFunc(funcName string, body interface{}) (interface{}, error) {
	if fn, exists := FunctionRegistry[funcName]; exists {
		return fn(body)
	} else {
		return nil, fmt.Errorf("no such UserFunction defined")
	}
}
