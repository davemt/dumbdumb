package main

import (
	"dumbdumb/handler"
	"os"
	"testing"
)

func TestMainForSmoke(t *testing.T) {
	os.Setenv("DUMBDUMB_SMTP_SENDER_DOMAIN_WHITELIST", "test.com")
	os.Setenv("DUMBDUMB_WEATHERUNDERGROUND_API_KEY", "123abc")
	os.Setenv("DUMBDUMB_GOOGLE_API_KEY", "456def")

	server := InitializeServer()

	h, err := server.RouteRequest("weather boston")
	if err != nil {
		t.Error("Routing to weather handler failed, had error:", err)
	}
	if (*h).(handler.WeatherHandler).WuApiKey != "123abc" {
		t.Error("Routing to weather handler failed")
	}

	h, err = server.RouteRequest("Find some place")
	if err != nil {
		t.Error("Routing to PlaceDirectoryHandler handler failed, had error:", err)
	}
	if (*h).(handler.PlaceDirectoryHandler).GoogleAPIKey != "456def" {
		t.Error("Routing to PlaceDirectoryHandler handler failed")
	}

	h, err = server.RouteRequest("Translate homard")
	if err != nil {
		t.Error("Routing to translate handler failed, had error:", err)
	}
	if (*h).(handler.TranslateHandler).GoogleAPIKey != "456def" {
		t.Error("Routing to translate handler failed")
	}
}
