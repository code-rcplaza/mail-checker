package usecases

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/emailsync/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func assertKeyedError(t *testing.T, err error, sentinel error, want shared.MessageKey) {
	t.Helper()
	if !errors.Is(err, sentinel) {
		t.Fatalf("errors.Is(err, %v) = false (err = %v), want true", sentinel, err)
	}
	var keyed shared.KeyedError
	if !errors.As(err, &keyed) || keyed.Key != want {
		t.Errorf("key = %q, want %q", keyed.Key, want)
	}
}

func TestSyncAll_RejectsMonthsOutOfRangeWithoutReading(t *testing.T) {
	for _, months := range []int{0, 25} {
		fixture := newSyncFixture([]domain.Credential{credential(1, accountOne)}, nil, bankA)

		_, err := fixture.useCase.Execute(context.Background(), months)

		assertKeyedError(t, err, shared.ErrInvalid, domain.KeyMonthsOutOfRange)
		if fixture.directory.credentialsCalls != 0 || len(fixture.reader.queries) != 0 {
			t.Errorf("months = %d: read %d credentials and %d mailboxes, want none", months, fixture.directory.credentialsCalls, len(fixture.reader.queries))
		}
	}
}

func TestSyncAll_SyncsAtRangeLimits(t *testing.T) {
	for _, months := range []int{1, 24} {
		fixture := newSyncFixture([]domain.Credential{credential(1, accountOne)}, nil, bankA)

		results, err := fixture.useCase.Execute(context.Background(), months)

		if err != nil || len(results) != 1 {
			t.Errorf("months = %d: results = %+v, err = %v, want one result and no error", months, results, err)
		}
	}
}

// blockingMailReader holds the first read until the test releases it and
// counts every read, so a second sync that gets through shows up at once.
type blockingMailReader struct {
	entered chan struct{}
	release chan struct{}
	reads   *atomic.Int32
}

func (b blockingMailReader) Read(context.Context, domain.Credential, domain.MailQuery) ([]domain.RawEmail, error) {
	if b.reads.Add(1) == 1 {
		b.entered <- struct{}{}
		<-b.release
	}
	return nil, nil
}

func TestSyncAll_SecondSyncWhileOneRunsIsConflict(t *testing.T) {
	reader := blockingMailReader{entered: make(chan struct{}), release: make(chan struct{}), reads: &atomic.Int32{}}
	directory := &fakeDirectory{credentials: []domain.Credential{credential(1, accountOne)}}
	useCase := NewSyncAll(directory, reader, []ports.BankEmailParser{bankA}, &fakeChargeWriter{}, newFakeUnparsed(), func() time.Time { return testNow })

	firstDone := make(chan error)
	go func() {
		_, err := useCase.Execute(context.Background(), 3)
		firstDone <- err
	}()
	<-reader.entered

	_, err := useCase.Execute(context.Background(), 3)

	close(reader.release)
	if reads := reader.reads.Load(); reads != 1 {
		t.Errorf("mailbox reads = %d, want 1 (the second sync must not read)", reads)
	}
	assertKeyedError(t, err, shared.ErrConflict, domain.KeySyncInProgress)
	if err := <-firstDone; err != nil {
		t.Errorf("first sync error = %v, want nil", err)
	}
}

func TestSyncAll_SyncsAgainAfterAPreviousOneFinished(t *testing.T) {
	fixture := newSyncFixture([]domain.Credential{credential(1, accountOne)}, nil, bankA)
	fixture.directory.credentialsErr = errors.New("database is locked")
	if _, err := fixture.useCase.Execute(context.Background(), 3); err == nil {
		t.Fatalf("first sync error = nil, want the storage error")
	}
	fixture.directory.credentialsErr = nil

	if _, err := fixture.useCase.Execute(context.Background(), 3); err != nil {
		t.Errorf("sync after a failed one: error = %v, want nil", err)
	}
	if _, err := fixture.useCase.Execute(context.Background(), 3); err != nil {
		t.Errorf("sync after a successful one: error = %v, want nil", err)
	}
}
