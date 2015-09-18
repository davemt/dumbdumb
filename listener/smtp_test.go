package listener

import (
	"smtpd"
	"testing"
)

// Helper to build a smtpd.Peer object
func makePeer() smtpd.Peer {
	return smtpd.Peer{
		HeloName:   "fakeservername",   // Server name used in HELO/EHLO command
		Username:   "",                 // Username from authentication, if authenticated
		Password:   "",                 // Password from authentication, if authenticated
		Protocol:   "smtp",             // Protocol used, SMTP or ESMTP
		ServerName: "fake.server.name", // A copy of Server.Hostname
		Addr:       nil,                // Network address (net.Addr)
		TLS:        nil,                // *tls.ConnectionState, if on TLS
	}
}

// Helper to build a SMTPListeer object
func makeListener(domainWhitelist []string) SMTPListener {
	if domainWhitelist == nil {
		return SMTPListener{}
	} else {
		return SMTPListener{DomainWhitelist: domainWhitelist}
	}
}

// checkSender should not error when domain whitelist is undefined (nil)
func TestCheckSenderNoWhitelist(t *testing.T) {
	l := makeListener(nil)

	err := l.checkSender(makePeer(), "dave@google.com")
	if err != nil {
		t.Error("Check sender returned error when expected to be OK, error: ", err)
	}
}

// checkSender should not error when domain whitelist is empty
func TestCheckSenderEmptyWhitelist(t *testing.T) {
	l := makeListener([]string{})

	err := l.checkSender(makePeer(), "dave@google.com")
	if err != nil {
		t.Error("Check sender returned error when expected to be OK, error: ", err)
	}
}

// checkSender should not error when sender domain matches domain in whitelist
func TestCheckSenderValidSender(t *testing.T) {
	domainWhitelist := []string{"google.com", "gmail.com"}
	l := makeListener(domainWhitelist)

	err := l.checkSender(makePeer(), "dave@google.com")
	if err != nil {
		t.Error("Check sender returned error when expected to be OK, error: ", err)
	}
}

// checkSender should error when sender domain does not match any whitelist domains
func TestCheckSenderInvalidSender(t *testing.T) {
	domainWhitelist := []string{"google.com", "gmail.com"}
	l := makeListener(domainWhitelist)

	err := l.checkSender(makePeer(), "dave@blah.com")
	if err == nil {
		t.Error("Should have blocked sender with non-whitelisted domain")
	}
}

// checkSender should work w/ full form addr, e.g. "Barry Gibbs <bg@example.com>"
func TestCheckSenderFullFormAddress(t *testing.T) {
	domainWhitelist := []string{"google.com", "gmail.com"}
	l := makeListener(domainWhitelist)

	err := l.checkSender(makePeer(), "Barry Gibbs <bg@gmail.com>")
	if err != nil {
		t.Error("Check sender returned error when expected to be OK, error: ", err)
	}
}
