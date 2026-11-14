package main

import (
	"context"
	"crypto/rand"
	"log/slog"
	"time"
)

// cryptoRandom is the source of codes and session tokens.
var cryptoRandom = rand.Reader

// cleanupInterval: expired codes and sessions are already rejected, so
// cleaning once an hour only keeps the tables small.
const cleanupInterval = time.Hour

// cleanupLoop runs the cleanup every interval until the app stops. A failure
// is logged and the next round tries again: nothing depends on it.
type cleanupLoop struct {
	clean  func(ctx context.Context) error
	every  time.Duration
	logger *slog.Logger
}

func (l cleanupLoop) run(ctx context.Context) {
	ticker := time.NewTicker(l.every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.cleanOnce(ctx)
		}
	}
}

// cleanOnce logs only the error: it carries cutoffs, never emails or tokens.
func (l cleanupLoop) cleanOnce(ctx context.Context) {
	if err := l.clean(ctx); err != nil {
		l.logger.Error("clean expired access failed", "err", err)
	}
}
