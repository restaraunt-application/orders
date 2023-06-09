package main

import (
	"context"
	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/otel"
	appshutdown "orders/internal/infra/app_shutdown"
	"orders/internal/infra/config"
	"orders/internal/infra/infraerrors"
	"orders/internal/infra/logger"
	sentryapp "orders/internal/infra/sentry_app"
	"orders/internal/infra/tracing"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := logger.GetLog()

	appConfig, err := config.NewAppConfig()
	if err != nil {
		log.Fatal(infraerrors.AppConfigError.WrapWithError(err))
	}

	log, flushFn, err := logger.NewZapLog(appConfig.LogConfig)
	if err != nil {
		log.Fatal(infraerrors.LogInitError.WrapWithError(err))
	}

	logger.ReplaceGlobal(log)

	err = sentryapp.InitSentry(appConfig.SentryConfig)
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	tp, err := tracing.InitTracerProvider(appConfig.TracingConfig)
	if err != nil {
		log.Fatalf("tracing.Init: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	otel.SetTracerProvider(tp)

	tr := tp.Tracer("component-main")

	ctx, span := tr.Start(ctx, "main")
	defer span.End()

	applicationTimeout := time.Duration(appConfig.ShutDownApplicationTimeout) * time.Second
	waitCh := appshutdown.ShutDownApplication(ctx, applicationTimeout, map[string]appshutdown.Handler{
		"sentry": func(ctx context.Context) error {
			trs := tp.Tracer("component-sentry-shutdown")

			ctx, spanSentry := trs.Start(ctx, "sentryShutdown")
			defer spanSentry.End()

			sentry.Flush(time.Duration(appConfig.SentryConfig.FlushTimeout) * time.Second)
			return nil
		},
		"tracing": func(ctx context.Context) error {
			trs := tp.Tracer("component-tracing-shutdown")

			ctx, spanSentry := trs.Start(ctx, "tracingShutdown")
			defer spanSentry.End()

			return tp.Shutdown(ctx)
		},
		"logger": func(ctx context.Context) error {
			trs := tp.Tracer("component-logger-shutdown")

			ctx, spanSentry := trs.Start(ctx, "loggerShutdown")
			defer spanSentry.End()

			return flushFn()
		},
	})

	<-waitCh
}
