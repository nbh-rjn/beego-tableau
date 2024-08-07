package controllers

import (
	"context"
	"errors"

	"github.com/cenkalti/backoff/v4"
)

func CallWithRetry(ctx context.Context, call func() error) error {
	wrappedCall := func() error {
		err := call()
		if !canRetryError(err) {
			return backoff.Permanent(err)
		}
		return err
	}

	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
	return backoff.Retry(wrappedCall, backoff.WithContext(b, ctx))
}

func canRetryError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	return true
}
