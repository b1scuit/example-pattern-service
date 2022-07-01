package main

import (
	"log"
	"os"

	"github.com/B1scuit/example-pattern-service/internal/core"
	"github.com/B1scuit/example-pattern-service/pkg/email"
	"github.com/B1scuit/example-pattern-service/pkg/sms"
	"github.com/spf13/cobra"
)

// This creates a CLI service init's from CLI flags
// as is a common pattern in cli applications
func main() {
	logger := log.New(os.Stdout, "", 0)

	var (
		to         string
		from       string
		subject    string
		body       string
		number     string
		fromNumber string
	)

	var coreClient *core.Client

	var rootCmd = &cobra.Command{
		// Prerun is a good place to set up client as it acts as a layer before
		// running the task
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			coreClient, err = core.New(&core.ClientOptions{
				Email: email.Must(email.New(&email.ClientOptions{
					StdLog:      logger,
					FromAddress: from,
				})),
				SMS: sms.Must(sms.New(&sms.ClientOptions{
					StdLog:     logger,
					FromNumber: fromNumber,
				})),
			})

			return err
		},
		// Run the task
		RunE: func(cmd *cobra.Command, args []string) error {
			return coreClient.Task1(cmd.Context(), &core.Task1Input{
				To:      to,
				From:    from,
				Number:  number,
				Subject: subject,
				Body:    body,
			})
		},
	}

	// Instead of env vars, this time we are loading config via CLI flags
	rootCmd.Flags().StringVarP(&to, "to", "t", "", "To email address (example@example.comn)")
	rootCmd.Flags().StringVarP(&from, "from", "f", "noreply@company.com", "From email address (example@example.comn)")
	rootCmd.Flags().StringVarP(&subject, "subject", "s", "Default title", "Message subject")
	rootCmd.Flags().StringVarP(&body, "body", "b", "Default content", "Message content")
	rootCmd.Flags().StringVarP(&number, "number", "n", "", "Mobile number for SMS (0123456789)")
	rootCmd.Flags().StringVarP(&fromNumber, "fromnumber", "a", "", "Mobile number to send SMS from (0123456789)")

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
