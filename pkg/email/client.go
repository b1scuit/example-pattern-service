// email
//
// This follows the exact same pattern as core and can be treated as an isolated "mini-core"
// of it's own responsiblilty, this is completly isolated away from any other assumed
// preconception of it's caller
package email

import (
	"context"
	"log"
	"os"
)

type ClientOptions struct {
	StdLog *log.Logger

	FromAddress string
}

type Client struct {
	stdLog *log.Logger

	fromAddress string
}

func New(opts *ClientOptions) (*Client, error) {

	// If the logger was missed, assume a default
	if opts.StdLog == nil {
		opts.StdLog = log.New(os.Stdout, "email", 0)
	}

	return &Client{
		stdLog:      opts.StdLog,
		fromAddress: opts.FromAddress,
	}, nil
}

// Forces a clean completion of New() for initalisation
func Must(client *Client, err error) *Client {
	if err != nil {
		panic(err)
	}

	return client
}

func (c *Client) Send(ctx context.Context, to, subject, body string) error {

	// Complete steps to send message, for now, we can just log
	c.stdLog.Printf("Sending email: To %v, From: %v, Subject: %v, Body: %v", to, c.fromAddress, subject, body)

	return nil
}
