package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/B1scuit/example-pattern-service/internal/core"
)

var errMock = errors.New("mock error")

// We dont provide param names here as we dont care about the input, only the output
// and making sure it conforms to the interface
type MockEmailClient struct {
	SendMock func(context.Context, string, string, string) error
}

func (mec *MockEmailClient) Send(ctx context.Context, s1, s2, s3 string) error {
	return mec.SendMock(ctx, s1, s2, s3)
}

type MockSMSClient struct {
	SendMock func(context.Context, string, string) error
}

func (mec *MockSMSClient) Send(ctx context.Context, s1, s2 string) error {
	return mec.SendMock(ctx, s1, s2)
}

// Setting these globally as defaults
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
			return errMock
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
	mockSMSClient := &MockSMSClient{
		SendMock: func(ctx context.Context, s1, s2 string) error {
			return errMock
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
