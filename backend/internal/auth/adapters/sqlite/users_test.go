package sqlite

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/shared/sqlite/sqlitetest"
)

var _ ports.UserRepository = Users{}

var santiago = time.FixedZone("America/Santiago", -3*60*60)

// acceptedAt is in Santiago with nanoseconds: it must come back in UTC, in seconds.
var acceptedAt = time.Date(2026, 11, 13, 9, 30, 15, 123456789, santiago)

var createdAt = time.Date(2026, 11, 13, 12, 30, 15, 0, time.UTC)

func newUsers(t *testing.T) Users {
	t.Helper()
	return NewUsers(sqlitetest.OpenMigrated(t), func() time.Time { return createdAt })
}

func ana() domain.User {
	return domain.User{PublicID: "0193f1a2-0000-7000-8000-000000000001", Name: "Ana", Email: "ana@example.com", AcceptedTermsAt: acceptedAt}
}

func createAna(t *testing.T, users Users) domain.User {
	t.Helper()
	created, err := users.Create(context.Background(), ana())
	if err != nil {
		t.Fatalf("Create error = %v, want nil", err)
	}
	return created
}

func TestUsers_CreateAssignsTheIDAndKeepsThePublicID(t *testing.T) {
	users := newUsers(t)

	created := createAna(t, users)

	if created.ID == 0 || created.PublicID != ana().PublicID {
		t.Errorf("created = %+v, want an ID and the public ID %s", created, ana().PublicID)
	}
}

func TestUsers_CreateRejectsATakenEmail(t *testing.T) {
	users := newUsers(t)
	createAna(t, users)
	again := ana()
	again.PublicID = "0193f1a2-0000-7000-8000-000000000002"
	again.Email = "ANA@Example.com"

	_, err := users.Create(context.Background(), again)

	if !errors.Is(err, shared.ErrConflict) {
		t.Errorf("Create error = %v, want ErrConflict", err)
	}
}

func TestUsers_DatesComeBackInUTCSeconds(t *testing.T) {
	users := newUsers(t)
	created := createAna(t, users)

	got, err := users.FindByID(context.Background(), created.ID)

	want := acceptedAt.Truncate(time.Second).UTC()
	if err != nil || !got.AcceptedTermsAt.Equal(want) || got.AcceptedTermsAt.Location() != time.UTC {
		t.Errorf("AcceptedTermsAt = %v, %v; want %v in UTC", got.AcceptedTermsAt, err, want)
	}
}

func TestUsers_FindByEmail(t *testing.T) {
	users := newUsers(t)
	created := createAna(t, users)
	cases := []struct {
		name    string
		email   string
		wantErr error
	}{
		{"same casing", "ana@example.com", nil},
		{"another casing", "Ana@EXAMPLE.com", nil},
		{"an unknown email", "nadie@example.com", shared.ErrNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := users.FindByEmail(context.Background(), tc.email)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("FindByEmail error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil && got != created {
				t.Errorf("FindByEmail = %+v, want %+v", got, created)
			}
		})
	}
}

func TestUsers_FindByID(t *testing.T) {
	users := newUsers(t)
	created := createAna(t, users)

	got, err := users.FindByID(context.Background(), created.ID)
	if err != nil || got != created {
		t.Errorf("FindByID = %+v, %v; want %+v", got, err, created)
	}

	_, err = users.FindByID(context.Background(), created.ID+1)
	if !errors.Is(err, shared.ErrNotFound) {
		t.Errorf("FindByID of an unknown ID error = %v, want ErrNotFound", err)
	}
}
