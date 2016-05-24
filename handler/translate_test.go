package handler

import (
	"fmt"
	"testing"
)

type MockRequest struct {
	Payload string
	Sender  string
}

func (r MockRequest) GetPayload() string             { return r.Payload }
func (r MockRequest) SendOutput(output string) error { return nil }

func makeRequest(payload string) MockRequest {
	return MockRequest{Payload: payload}
}

func makeHandler() TranslateHandler {
	return TranslateHandler{GoogleAPIKey: "123abcd"}
}

func TestBuildAPIParamsFromRequest(t *testing.T) {
	h := makeHandler()
	r := makeRequest("Translate en -> fr lobster")

	params, err := h.BuildAPIParamsFromRequest(r)

	if err != nil {
		t.Error(fmt.Sprintf("Failed to parse a request that was valid: %v", err))
	}
	if params["key"][0] != "123abcd" {
		t.Error("Failed to add google API key to params")
	}
	if params["source"][0] != "en" || params["target"][0] != "fr" {
		t.Error("Failed to parse source/target langs from request")
	}
	if params["q"][0] != "lobster" {
		t.Error("Failed to parse to-translate phrase from request")
	}
}

func TestBuildAPIParamsFromRequestInvalid(t *testing.T) {
	h := makeHandler()
	r := makeRequest("Translate en - lobster")

	params, err := h.BuildAPIParamsFromRequest(r)

	if err == nil {
		t.Error(fmt.Sprintf("Should have returned error with invalid request"))
	}

	if params != nil {
		t.Error(fmt.Sprintf("Should have returned nil params with invalid request"))
	}
}
