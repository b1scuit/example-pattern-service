package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/B1scuit/example-pattern-service/internal/core"
	"github.com/gorilla/mux"
)

type CoreClientInterface interface {
	Task1(context.Context, *core.Task1Input) error
}

type ClientOptions struct {
	StdLog *log.Logger

	HttpServer *http.Server

	Core CoreClientInterface
}

type Client struct {
	stdLog *log.Logger

	httpServer *http.Server

	core CoreClientInterface
}

func New(opts *ClientOptions) (*Client, error) {

	// If the logger was missed, assume a default
	if opts.StdLog == nil {
		opts.StdLog = log.New(os.Stdout, "http", 0)
	}

	if opts.HttpServer == nil {
		opts.HttpServer = &http.Server{
			Addr:         "127.0.0.1:8000",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
	}

	if opts.Core == nil {
		return nil, errors.New("core missing")
	}

	return &Client{
		stdLog:     opts.StdLog,
		httpServer: opts.HttpServer,

		core: opts.Core,
	}, nil
}

func Must(client *Client, err error) *Client {
	if err != nil {
		panic(err)
	}

	return client
}

func (c *Client) RunServer() error {
	router := mux.NewRouter()

	router.HandleFunc("/", c.Task1Handler)

	c.httpServer.Handler = router

	go func() {
		if err := c.httpServer.ListenAndServe(); err != nil {
			c.stdLog.Println(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(signalChan, os.Interrupt)

	// Block until we receive our signal.
	<-signalChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	return c.httpServer.Shutdown(ctx)
}

func (c *Client) Task1Handler(w http.ResponseWriter, r *http.Request) {

	// Decode user input
	var input core.Task1Input
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Run the core function
	if err := c.core.Task1(r.Context(), &input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	// respond all completed
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Done")
}
