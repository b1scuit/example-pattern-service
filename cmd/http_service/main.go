package main

import (
	"log"
	"os"

	"github.com/B1scuit/example-pattern-service/internal/core"
	"github.com/B1scuit/example-pattern-service/pkg/email"
	"github.com/B1scuit/example-pattern-service/pkg/http"
	"github.com/B1scuit/example-pattern-service/pkg/sms"
)

// This creates a HTTP service init'd from env vars
// as is a common pattern in microservices
// You can see each client as it's loaded
// and what goes into each client
// the main is a reference for how the packages are being run
func main() {
	logger := log.New(os.Stdout, "", 0)

	httpServer := http.Must(http.New(&http.ClientOptions{
		StdLog: logger,
		Core: core.Must(core.New(&core.ClientOptions{
			Email: email.Must(email.New(&email.ClientOptions{
				StdLog:      logger,
				FromAddress: os.Getenv("FROM_EMAIL_ADDRESS"),
			})),
			SMS: sms.Must(sms.New(&sms.ClientOptions{
				StdLog:     logger,
				FromNumber: os.Getenv("FROM_SMS_NUMBER"),
			})),
		})),
	}))

	if err := httpServer.RunServer(); err != nil {
		logger.Fatal(err)
	}
}
