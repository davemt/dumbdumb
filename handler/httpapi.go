package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type buildUrlFromRequest func(request string) string

type getOutputFromResponseData func(data interface{}) string

type HTTPAPIHandler struct {
	// This function must be specified when an instance of HTTPAPIHandler is
	// initialized. It is responsible for taking a request string and building
	// a URL.
	BuildUrlFromRequest buildUrlFromRequest
	// This function must be specified when an instance of HTTPAPIHandler is
	// initialized. It is responsible for taking the data returned in the HTTP
	// response and translating it to an output string.
	GetOutputFromResponseData getOutputFromResponseData
}

func (h HTTPAPIHandler) HandleRequest(request string) (string, error) {
	log.Printf("HTTP API handler got request: %v", request)
	resp, err := http.Get(h.BuildUrlFromRequest(request))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var f interface{}
	err2 := json.Unmarshal(body, &f)
	output := h.GetOutputFromResponseData(f)

	return output, err2
}
