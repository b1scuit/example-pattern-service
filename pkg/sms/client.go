// sms
//
// This follows the exact same pattern as core and can be treated as an isolated "mini-core"
// of it's own responsiblilty, this is completly isolated away from any other assumed
// preconception of it's caller
package sms

import (
	"context"
	"log"
	"os"
)

type ClientOptions struct {
	StdLog *log.Logger

	FromNumber string
}

type Client struct {
	stdLog *log.Logger

	fromNumber string
}

func New(opts *ClientOptions) (*Client, error) {

	// If the logger was missed, assume a default
	if opts.StdLog == nil {
		opts.StdLog = log.New(os.Stdout, "email", 0)
	}

	return &Client{
		stdLog: opts.StdLog,

		fromNumber: opts.FromNumber,
	}, nil
}

// Forces a clean completion of New() for initalisation
func Must(client *Client, err error) *Client {
	if err != nil {
		panic(err)
	}

	return client
}

func (c *Client) Send(ctx context.Context, to, body string) error {

	// Complete steps to send sms, for now, we can just log
	c.stdLog.Printf("Sending SMS: To number: %v, from number %v, Body content: %v", to, c.fromNumber, body)

	return nil
}
