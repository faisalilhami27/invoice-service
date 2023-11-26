package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
)

//nolint:revive
type SentryPackage struct {
	Dsn           string  `validate:"required"`
	Env           string  `validate:"required"`
	SampleRate    float64 `validate:"required"`
	EnableTracing bool
	Debug         bool
}

type Option func(*SentryPackage)

func WithDsn(dsn string) Option {
	return func(s *SentryPackage) {
		s.Dsn = dsn
	}
}

func WithDebug(debug bool) Option {
	return func(s *SentryPackage) {
		s.Debug = debug
	}
}

func WithEnv(env string) Option {
	return func(s *SentryPackage) {
		s.Env = env
	}
}

func WithSampleRate(sampleRate float64) Option {
	return func(s *SentryPackage) {
		s.SampleRate = sampleRate
	}
}

func WithEnableTracing(enableTracing bool) Option {
	return func(s *SentryPackage) {
		s.EnableTracing = enableTracing
	}
}

type ISentry interface {
	StartSpan(ctx context.Context, spanName string) *sentry.Span
	Finish(span *sentry.Span)
	CaptureException(exception error) *sentry.EventID
	SpanContext(span *sentry.Span) context.Context
}

func NewSentry(options ...Option) ISentry {
	sentryPkg := &SentryPackage{
		EnableTracing: true,
	}
	for _, option := range options {
		option(sentryPkg)
	}

	err := sentry.Init(sentry.ClientOptions{
		EnableTracing:    sentryPkg.EnableTracing,
		Dsn:              sentryPkg.Dsn,
		Debug:            sentryPkg.Debug,
		Environment:      sentryPkg.Env,
		TracesSampleRate: sentryPkg.SampleRate,
	})
	if err != nil {
		panic(err)
	}

	return sentryPkg
}

func (s *SentryPackage) StartSpan(ctx context.Context, spanName string) *sentry.Span {
	return sentry.StartSpan(ctx, spanName)
}

func (s *SentryPackage) Finish(span *sentry.Span) {
	span.Finish()
}

func (s *SentryPackage) CaptureException(exception error) *sentry.EventID {
	return sentry.CaptureException(exception)
}

func (s *SentryPackage) SpanContext(span *sentry.Span) context.Context {
	return span.Context()
}
