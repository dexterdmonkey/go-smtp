package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
)

// Interface defines the methods that any SMTP client must implement.
type Interface interface {
	GetSenderAddress() string
	GetPassword() string
	GetHost() string
	GetPort() int
	ParseBody(body string, parameters map[string]interface{}) string
	SendMail(email Email) error
}

// Email struct represents the email structure with recipients, subject, and body.
type Email struct {
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
}

// SMTP struct represents the SMTP client with necessary credentials and configurations.
type SMTP struct {
	senderAddress string
	password      string
	host          string
	port          string
	auth          smtp.Auth
}

// New initializes and returns a new SMTP client.
func New(senderAddress, password, host string, port int) (*SMTP, error) {
	auth := smtp.PlainAuth("", senderAddress, password, host)
	if auth == nil {
		return nil, fmt.Errorf("auth error, empty auth")
	}

	c := &SMTP{
		senderAddress: senderAddress,
		password:      password,
		host:          host,
		port:          strconv.Itoa(port),
		auth:          auth,
	}

	return c, nil
}

// GetSenderAddress returns the sender's email address.
func (c *SMTP) GetSenderAddress() string {
	return c.senderAddress
}

// GetPassword returns the password for the SMTP client.
func (c *SMTP) GetPassword() string {
	return c.password
}

// GetHost returns the host for the SMTP client.
func (c *SMTP) GetHost() string {
	return c.host
}

// GetPort returns the port for the SMTP client as an integer.
func (c *SMTP) GetPort() int {
	port, _ := strconv.Atoi(c.port)
	return port
}

// GetClient initializes and returns an SMTP client.
func (c *SMTP) GetClient() (*smtp.Client, error) {
	client, err := smtp.Dial(c.host + ":" + c.port)
	if err != nil {
		return nil, fmt.Errorf("client error, failed to dial; %s", err.Error())
	}

	if err = client.StartTLS(&tls.Config{InsecureSkipVerify: true, ServerName: c.host}); err != nil {
		return nil, fmt.Errorf("client error, failed to start tls; %s", err.Error())
	}

	if err = client.Auth(c.auth); err != nil {
		return nil, fmt.Errorf("client error, failed to apply auth; %s", err.Error())
	}

	if err = client.Mail(c.senderAddress); err != nil {
		return nil, fmt.Errorf("client error, failed to create mail; %s", err.Error())
	}

	return client, nil
}

// SendMail sends an email with the specified content and recipients.
func (c *SMTP) SendMail(email Email) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	defer client.Close()

	// Send mail to recipients
	for _, addr := range email.To {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("send error, failed to add recipients; %s", err.Error())
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("send error, failed to create data; %s", err.Error())
	}
	defer func() {
		err = w.Close()
		if err != nil {
			fmt.Printf("send error, failed to close email writer; %s\n", err.Error())
		}
	}()

	ccStmt := ""
	if len(email.Cc) != 0 {
		ccStmt = "Cc: " + strings.Join(email.Cc, ",") + "\r\n"
	}

	bccStmt := ""
	if len(email.Bcc) != 0 {
		bccStmt = "Bcc: " + strings.Join(email.Bcc, ",") + "\r\n"
	}

	message := []byte(
		"Subject: " + email.Subject + "\r\n" +
			"To: " + strings.Join(email.To, ",") + "\r\n" +
			ccStmt +
			bccStmt +
			"\r\n" +
			email.Body + "\r\n",
	)

	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("send error, failed to send email from %s [%s:%s], %s", c.senderAddress, c.host, c.port, err.Error())
	}

	return nil
}

// ParseBody replaces placeholders in the email body with actual values from the parameters map.
func (c *SMTP) ParseBody(body string, parameters map[string]interface{}) string {
	for key, value := range parameters {
		placeholder := "{{" + key + "}}"
		body = strings.Replace(body, placeholder, fmt.Sprintf("%v", value), -1)
	}

	return body
}
