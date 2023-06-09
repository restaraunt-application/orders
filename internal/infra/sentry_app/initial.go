package sentryapp

import (
	"orders/internal/infra/config"

	"github.com/getsentry/sentry-go"
)

func InitSentry(sentryConfig *config.SentryConfig) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://examplePublicKey@o0.ingest.sentry.io/0",
		Debug:       sentryConfig.Debug,
		Environment: sentryConfig.Stage,
	})

	return err
}
