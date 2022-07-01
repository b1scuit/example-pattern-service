package http_test

import (
	"context"
	"errors"
	"fmt"
	h "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/B1scuit/example-pattern-service/internal/core"
	"github.com/B1scuit/example-pattern-service/pkg/http"
)

var errMock = errors.New("mock error")

// Holding the error on the mock logger
// lets me pull it back out again as the HTTP server
// is run in a go routine and the only feedback is through
// the logger where it sends what's happened, this allows me to
// check something was sent to the logger, check TestListenAndServeFail
type MockLogger struct {
	Err string
}

func (ml *MockLogger) Println(in ...any) {
	ml.Err = fmt.Sprint(in[0])
}

func (ml *MockLogger) GetError() error {
	if ml.Err == "" {
		return nil
	}

	return errors.New(ml.Err)
}

type MockCore struct {
	Err bool
}

func (mc *MockCore) Task1(context.Context, *core.Task1Input) error {
	if mc.Err {
		return errMock
	}

	return nil
}

func TestNew(t *testing.T) {

	client, err := http.New(&http.ClientOptions{
		Core: &MockCore{},
	})
	if err != nil {
		t.Error(err)
		return
	}

	if client == nil {
		t.Error("nil client returned")
	}
}

func TestListenAndServeFail(t *testing.T) {

	var log MockLogger

	client, _ := http.New(&http.ClientOptions{
		Core:   &MockCore{},
		StdLog: &log,
		HttpServer: &h.Server{
			Addr: "999.999.999.999:1234567688",
		},
	})

	go func() {
		if err := client.RunServer(); err != nil {
			t.Error(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client.ExposeHttpServer().Shutdown(ctx)

	<-ctx.Done()

	if err := log.GetError(); err == nil {
		t.Error("error should have been triggered")
	}
}

func TestNewMissingCore(t *testing.T) {
	_, err := http.New(&http.ClientOptions{})
	if err == nil {
		t.Error("error should have triggered")
	}
}

func TestMustClean(t *testing.T) {
	// This deferal function allows for the testing
	//of panics as it blocks the os.Exit using recover()
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	http.Must(&http.Client{}, nil)
}

func TestServerRun(t *testing.T) {
	client, err := http.New(&http.ClientOptions{
		Core: &MockCore{},
	})
	if err != nil {
		t.Error(err)
		return
	}

	go func() {
		if err := client.RunServer(); err != nil {
			t.Error(err)
		}
	}()
}

func TestTaskHandler(t *testing.T) {
	httpClient, err := http.New(&http.ClientOptions{
		Core: &MockCore{},
	})
	if err != nil {
		t.Error(err)
		return
	}

	if httpClient == nil {
		t.Error("nil client returned")
		return
	}

	req, _ := h.NewRequest("", "", strings.NewReader("{}"))
	recorder := httptest.NewRecorder()
	handler := h.HandlerFunc(httpClient.Task1Handler)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != h.StatusOK {
		t.Error(recorder.Body.String())
	}

}
func TestTaskHandlerDecodeFail(t *testing.T) {
	httpClient, err := http.New(&http.ClientOptions{
		Core: &MockCore{Err: true},
	})
	if err != nil {
		t.Error(err)
		return
	}

	if httpClient == nil {
		t.Error("nil client returned")
		return
	}

	req, _ := h.NewRequest("", "", strings.NewReader("}"))
	recorder := httptest.NewRecorder()
	handler := h.HandlerFunc(httpClient.Task1Handler)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != h.StatusBadRequest {
		t.Error("error should have returned")
	}
}

func TestTaskHandlerFail(t *testing.T) {
	httpClient, err := http.New(&http.ClientOptions{
		Core: &MockCore{Err: true},
	})
	if err != nil {
		t.Error(err)
		return
	}

	if httpClient == nil {
		t.Error("nil client returned")
		return
	}

	req, _ := h.NewRequest("", "", strings.NewReader("{}"))
	recorder := httptest.NewRecorder()
	handler := h.HandlerFunc(httpClient.Task1Handler)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != h.StatusInternalServerError {
		t.Error("error should have returned")
	}
}

func TestMustPanic(t *testing.T) {
	// This deferal function allows for the testing
	//of panics as it blocks the os.Exit using recover()
	defer func() {
		if r := recover(); r == nil {
			t.Error("panic should have thrown")
		}
	}()

	http.Must(&http.Client{}, errMock)
}
