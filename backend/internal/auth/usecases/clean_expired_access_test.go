package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
)

var _ ports.ExpiredAccess = (*recordingCleanup)(nil)

type recordingCleanup struct {
	cutoffs []domain.CleanupCutoff
	err     error
}

func (r *recordingCleanup) DeleteExpired(_ context.Context, cutoff domain.CleanupCutoff) error {
	r.cutoffs = append(r.cutoffs, cutoff)
	return r.err
}

func TestCleanExpiredAccess_DeletesWhatExpiredByNow(t *testing.T) {
	cleanup := &recordingCleanup{}
	clock := &testClock{now: testNow}

	err := NewCleanExpiredAccess(cleanup, clock.Now).Execute(context.Background())

	want := domain.CleanupCutoff{ExpiredBy: testNow, SentBy: testNow.Add(-time.Hour)}
	if err != nil || len(cleanup.cutoffs) != 1 || cleanup.cutoffs[0] != want {
		t.Errorf("cutoffs = %+v, err = %v; want one cleanup by %+v", cleanup.cutoffs, err, want)
	}
}

func TestCleanExpiredAccess_WrapsTheStoreError(t *testing.T) {
	cleanup := &recordingCleanup{err: errDatabase}

	err := NewCleanExpiredAccess(cleanup, (&testClock{now: testNow}).Now).Execute(context.Background())

	assertWrapped(t, err, errDatabase)
}
