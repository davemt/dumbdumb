package listener

import (
	"bytes"
	"dumbdumb"
	"fmt"
	"github.com/mattbaird/gochimp"
	"log"
	"net/mail"
	"os"
	"smtpd"
	"strings"
)

// Listener for requests sent via SMTP.  The output generated when the request
// is handled will be sent back via SMTP as well.
type SMTPListener struct {
	incoming chan dumbdumb.Request
	// whitelist of sender domains to accept, nil means no restrictions
	DomainWhitelist []string
}

// Called after MAIL FROM, validates sender address (check whitelist, etc.)
func (l *SMTPListener) checkSender(peer smtpd.Peer, addr string) error {
	if l.DomainWhitelist == nil || len(l.DomainWhitelist) == 0 {
		return nil
	}
	address, err := mail.ParseAddress(addr)
	if err != nil {
		log.Printf("Failed to parse sender address: %v", addr)
		return smtpd.Error{Code: 501, Message: "Bad sender address"}
	}
	parts := strings.SplitN(address.Address, "@", 2)
	_, domain := parts[0], parts[1]

	// check that domain is on whitelist
	domainOk := false
	for _, allowedDomain := range l.DomainWhitelist {
		if domain == allowedDomain {
			domainOk = true
			break
		}
	}
	if !domainOk {
		log.Printf("Rejected a sender that was not on domain whitelist: %v", addr)
		return smtpd.Error{Code: 554, Message: "Bad sender address"}
	}
	return nil
}

func (l *SMTPListener) handleMail(peer smtpd.Peer, env smtpd.Envelope) error {
	sender := env.Sender
	msgReader, err := mail.ReadMessage(bytes.NewReader(env.Data))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(msgReader.Body)
	body := buf.String()

	log.Printf("Got SMTP message from %v: %v", sender, body)
	l.incoming <- SMTPRequest{Payload: strings.TrimSpace(body), Sender: sender}
	return nil
}

// TODO: try changing to pointer receiver and see what happens
func (l SMTPListener) Listen(incoming chan dumbdumb.Request) {
	l.incoming = incoming

	server := &smtpd.Server{
		WelcomeMessage: "SMTP Listener ready.",
		Handler:        l.handleMail,
		SenderChecker:  l.checkSender,
	}

	err := server.ListenAndServe("0.0.0.0:25")

	if err != nil {
		log.Fatal(err)
	}
}

type SMTPRequest struct {
	Payload string
	Sender  string
}

func (r SMTPRequest) GetPayload() string {
	return r.Payload
}

func (r SMTPRequest) SendOutput(output string) error {
	response := SMTPResponse{Payload: output, Sender: r.Sender}
	err := response.Send()
	return err
}

type SMTPResponse struct {
	Payload string
	Sender  string
}

func (r SMTPResponse) Send() error {
	log.Printf("Responding via SMTP with: %v", r.Payload)

	apiKey := os.Getenv("DUMBDUMB_SMTP_MANDRILL_API_KEY")
	mandrillApi, err := gochimp.NewMandrill(apiKey)

	recipients := []gochimp.Recipient{gochimp.Recipient{Email: r.Sender}}

	message := gochimp.Message{
		Text:      r.Payload,
		FromEmail: os.Getenv("DUMBDUMB_SMTP_FROM_EMAIL"),
		To:        recipients,
	}

	_, err = mandrillApi.MessageSend(message, false)

	if err != nil {
		fmt.Println("Error sending SMTP message")
	}
	return err
}
