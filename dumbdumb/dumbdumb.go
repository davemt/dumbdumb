package main

import (
	"dumbdumb"
	"dumbdumb/handler"
	"dumbdumb/listener"
	"os"
	"strings"
)

func InitializeServer() *dumbdumb.Server {
	server := dumbdumb.NewServer()

	// Whitespace-separated list of email address domains
	domainWhitelist := os.Getenv("DUMBDUMB_SMTP_SENDER_DOMAIN_WHITELIST")
	server.AddListener(listener.SMTPListener{
		DomainWhitelist: strings.Fields(domainWhitelist),
	})

	server.AddHandler("[Ww]eather.*", handler.WeatherHandler{
		WuApiKey: os.Getenv("DUMBDUMB_WEATHERUNDERGROUND_API_KEY"),
	})

	server.AddHandler("411.*", handler.PlaceDirectoryHandler{
		GoogleAPIKey: os.Getenv("DUMBDUMB_GOOGLE_API_KEY"),
	})

	return server
}

func main() {
	server := InitializeServer()
	server.ListenAndServe()
}
