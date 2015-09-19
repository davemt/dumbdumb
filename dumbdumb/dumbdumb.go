package main

import (
	"dumbdumb"
	"dumbdumb/handler"
	"dumbdumb/listener"
	"os"
)

func main() {
	server := &dumbdumb.Server{}

	// Whitespace-separated list of email address domains
	domainWhitelist = os.Getenv("DUMBDUMB_SMTP_SENDER_DOMAIN_WHITELIST")
	server.AddListener(listener.SMTPListener{
		DomainWhitelist: strings.Fields(domainWhitelist),
	})

	server.AddHandler("[Ww]eather.*", handler.WeatherHandler{
		WuApiKey: os.Getenv("DUMBDUMB_WEATHERUNDERGROUND_API_KEY"),
	})

	server.ListenAndServe()
}
