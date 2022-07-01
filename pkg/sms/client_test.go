package sms_test

import (
	"context"
	"errors"
	"testing"

	"github.com/B1scuit/example-pattern-service/pkg/sms"
)

var errMock = errors.New("mock error")

func TestNewClient(t *testing.T) {
	client, err := sms.New(&sms.ClientOptions{})

	if err != nil {
		t.Error(err)
	}

	t.Run("Send", func(t *testing.T) {
		if err := client.Send(context.TODO(), "", ""); err != nil {
			t.Error(err)
		}
	})
}

func TestMustClean(t *testing.T) {
	// This deferal function allows for the testing
	//of panics as it blocks the os.Exit using recover()
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	sms.Must(&sms.Client{}, nil)
}

func TestMustPanic(t *testing.T) {
	// This deferal function allows for the testing
	//of panics as it blocks the os.Exit using recover()
	defer func() {
		if r := recover(); r == nil {
			t.Error("panic should have thrown")
		}
	}()

	sms.Must(&sms.Client{}, errMock)
}
