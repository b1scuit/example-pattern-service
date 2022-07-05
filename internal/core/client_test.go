package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/B1scuit/example-pattern-service/internal/core"
)

// We create a mock client that conformes to the interface we are expecting
// however we hold the function to be executed as a field ans simply call it
// This allow us to really control what the function returns as part of the unit
// tests, this also has the side effect of the unit test being more easily understood.
type MockEmailClient struct {
	SendMock func(context.Context, string, string, string) error
}

func (mec *MockEmailClient) Send(ctx context.Context, s1, s2, s3 string) error {
	return mec.SendMock(ctx, s1, s2, s3)
}

// A similer mock created for the SMS client interface
type MockSMSClient struct {
	SendMock func(context.Context, string, string) error
}

func (mec *MockSMSClient) Send(ctx context.Context, s1, s2 string) error {
	return mec.SendMock(ctx, s1, s2)
}

// We create some "best case" defaults since in most tests that's what
// we'll be using, we create some defaults with the expected ideal behaviour
// and then change away from this behaviour in the unit tests towards what
// we are trying to test, this leaves the unit test more ideally understood
// and consise from a testing pespective
var mockEmailClient core.EmailService = &MockEmailClient{
	SendMock: func(context.Context, string, string, string) error {
		return nil
	},
}

var mockSMSClient core.SMSService = &MockSMSClient{
	SendMock: func(ctx context.Context, s1, s2 string) error {
		return nil
	},
}

// Testing for clean init of the client
func TestClient(t *testing.T) {

	client, err := core.New(&core.ClientOptions{
		Email: mockEmailClient,
		SMS:   mockSMSClient,
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
	mockEmailClient := &MockEmailClient{
		SendMock: func(ctx context.Context, s1, s2, s3 string) error {
			return errors.New("Example Error")
		},
	}

	client, err := core.New(&core.ClientOptions{
		Email: mockEmailClient,
		SMS:   mockSMSClient,
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
	// An example of moving away from the ideal execution path
	// this specifically triggers an error to test error handling
	//
	// Note: using := here created a new variable mockSMSClient scoped to this function
	// it does not override the global best case
	mockSMSClient := &MockSMSClient{
		SendMock: func(ctx context.Context, s1, s2 string) error {
			return errors.New("Example error")
		},
	}

	client, err := core.New(&core.ClientOptions{
		Email: mockEmailClient,
		SMS:   mockSMSClient,
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
	// of panics as it blocks the os.Exit using recover()
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	core.Must(&core.Client{}, nil)
}

func TestMustPanic(t *testing.T) {
	// This deferal function allows for the testing
	// of panics as it blocks the os.Exit using recover()
	defer func() {
		if r := recover(); r == nil {
			t.Error("panic should have thrown")
		}
	}()

	core.Must(&core.Client{}, errors.New("Example"))
}
