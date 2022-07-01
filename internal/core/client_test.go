package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/B1scuit/example-pattern-service/internal/core"
)

var errMock = errors.New("mock error")

type MockEmailClient struct {
	Err bool //Whether to throw an error or not
}

// We dont provide param names here as we dont care about the input, only the output
// and making sure it conforms to the interface
func (mec *MockEmailClient) Send(context.Context, string, string, string) error {
	if mec.Err {
		return errMock
	}

	return nil
}

type MockSMSClient struct {
	Err bool
}

func (msc *MockSMSClient) Send(context.Context, string, string) error {
	if msc.Err {
		return errMock
	}

	return nil
}

// Testing for clean init of the client
func TestClient(t *testing.T) {
	client, err := core.New(&core.ClientOptions{
		Email: &MockEmailClient{},
		SMS:   &MockSMSClient{},
	})

	if err != nil {
		t.Error(err)
	}

	// Testing the task 1 function off the client
	// you can make client package level if you want to
	// test these independantly
	t.Run("Task1", func(t *testing.T) {
		if err := client.Task1(context.TODO(), &core.Task1Input{}); err != nil {
			t.Error(err)
		}
	})
}

func TestEmailErr(t *testing.T) {
	client, err := core.New(&core.ClientOptions{
		Email: &MockEmailClient{Err: true},
		SMS:   &MockSMSClient{},
	})

	if err != nil {
		t.Error(err)
	}

	// Testing the task 1 function off the client
	// you can make client package level if you want to
	// test these independantly
	t.Run("Task1", func(t *testing.T) {
		if err := client.Task1(context.TODO(), &core.Task1Input{}); err == nil {
			t.Error("error should have been returned")
		}
	})
}

func TestSMSErr(t *testing.T) {
	client, err := core.New(&core.ClientOptions{
		Email: &MockEmailClient{},
		SMS:   &MockSMSClient{Err: true},
	})

	if err != nil {
		t.Error(err)
	}

	// Testing the task 1 function off the client
	// you can make client package level if you want to
	// test these independantly
	t.Run("Task1", func(t *testing.T) {
		if err := client.Task1(context.TODO(), &core.Task1Input{Number: "0123456789"}); err == nil {
			t.Error("error should have been returned")
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

	core.Must(&core.Client{}, nil)
}

func TestMustPanic(t *testing.T) {
	// This deferal function allows for the testing
	//of panics as it blocks the os.Exit using recover()
	defer func() {
		if r := recover(); r == nil {
			t.Error("panic should have thrown")
		}
	}()

	core.Must(&core.Client{}, errMock)
}
