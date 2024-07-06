# go-smtp

`go-smtp` is a Go library for sending emails via SMTP. It provides an easy-to-use interface for sending emails with support for custom authentication, TLS, and more.

## Installation

To install the library, use `go get`:

```sh
go get github.com/dexterdmonkey/go-smtp
```

## Usage

Here's a simple example of how to use the `go-smtp` library:

```go
package main

import (
	"fmt"

	"github.com/dexterdmonkey/go-smtp"
)

func main() {
	var mail smtp.Interface

	SMTPUser := "your@email.com"
	SMTPPassword := "yourpassword"
	SMTPHost := "smtp.email.com"
	SMTPPort := 587

	mail, err := smtp.New(
		SMTPUser,
		SMTPPassword,
		SMTPHost,
		SMTPPort,
	)

	if err != nil {
		fmt.Println("error; ", err.Error())
		return
	}

	err = mail.SendMail(smtp.Email{
		To:      []string{"recipient1@email.com"},
		Subject: "subject",
		Body:    "body",
	})

	if err != nil {
		fmt.Println("error; ", err.Error())
		return
	}

	err = mail.SendMail(smtp.Email{
		To:      []string{"recipient2@email.com"},
		Subject: "subject",
		Body:    "body",
	})

	if err != nil {
		fmt.Println("error; ", err.Error())
		return
	}

	fmt.Println("Success")
}
```

## API

### Types

#### Interface

The `Interface` defines the methods available for the SMTP client:

```go
type Interface interface {
	GetSenderAddress() string
	GetPassword() string
	GetHost() string
	GetPort() int
	ParseBody(body string, parameters map[string]interface{}) string
	SendMail(email Email) error
}
```

#### Email

The `Email` struct represents an email to be sent:

```go
type Email struct {
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
}
```

### Functions

#### New

Creates a new SMTP client:

```go
func New(senderAddress, password, host string, port int) (*SMTP, error)
```

#### GetSenderAddress

Returns the sender address:

```go
func (c *SMTP) GetSenderAddress() string
```

#### GetPassword

Returns the password:

```go
func (c *SMTP) GetPassword() string
```

#### GetHost

Returns the host:

```go
func (c *SMTP) GetHost() string
```

#### GetPort

Returns the port:

```go
func (c *SMTP) GetPort() int
```

#### GetClient

Returns an authenticated SMTP client:

```go
func (c *SMTP) GetClient() (*smtp.Client, error)
```

#### SendMail

Sends an email:

```go
func (c *SMTP) SendMail(email Email) error
```

#### ParseBody

Parses the body of the email with the provided parameters:

```go
func (c *SMTP) ParseBody(body string, parameters map[string]interface{}) string
```

## License

This library is licensed under the Apache License, Version 2.0. See the LICENSE file for details.