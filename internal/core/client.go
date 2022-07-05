package core

import (
	"context"
)

// These interfaces allow the decoupling and ease of unit testing.
// Some developers prefer to break these off into seperate files/packages
// but I prefer to keep these as close to where they are used as possible
// so they remain paid to the client using them
type EmailService interface {
	Send(context.Context, string, string, string) error
}

type SMSService interface {
	Send(context.Context, string, string) error
}

type ClientOptions struct {
	PassedValue string

	Email EmailService
	SMS   SMSService
}

type Client struct {
	// These are lowercased so they are psudo-immutable, that way once you've made
	// an instance of Client, you have an assured single source of truth for configuration
	// There can be methods to change these, but you can hold mutex locks and other practices
	// to make the changes safe
	passedValue string

	email EmailService
	sms   SMSService
}

// Single point of entry to create a new instance of client
// to a known good working order, the user can still just create
// a client with &Client{} however there be dragons.
func New(opts *ClientOptions) (*Client, error) {

	// Any extra initalisation / Setup can exist here
	if opts.PassedValue == "" {
		opts.PassedValue = "Example value"
	}

	return &Client{
		passedValue: opts.PassedValue,

		email: opts.Email,
		sms:   opts.SMS,
	}, nil
}

// Forces a clean completion of New() for initalisation
func Must(client *Client, err error) *Client {
	if err != nil {
		panic(err)
	}

	return client
}

// The actual functions the core controller excutes need to be clear
// easy to read functions that can execute a task at the highest abstraction
// level of the service, what this function does is simple to read and follow
// allowing for a quick knowledge transfer
//
// If you have many of these functions, it is worth seperating them into different files
func (c *Client) Task1(ctx context.Context, in *Task1Input) error {

	if err := c.email.Send(ctx, in.To, in.Subject, in.Body); err != nil {
		return err
	}

	if in.IsNumberSet() {
		if err := c.sms.Send(ctx, in.Number, in.Body); err != nil {
			return err
		}
	}

	return nil
}
