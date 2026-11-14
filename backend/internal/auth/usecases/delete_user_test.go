package usecases

import (
	"context"
	"errors"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestDeleteUser(t *testing.T) {
	cases := []struct {
		name      string
		userID    int64
		deleteErr error
		wantErr   error
		wantLeft  int
	}{
		{"an existing user is deleted", ana.ID, nil, nil, 0},
		{"an unknown user is not found", ana.ID + 1, nil, shared.ErrNotFound, 1},
		{"a repository failure comes back", ana.ID, errDatabase, errDatabase, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			users := usersWith(ana)
			users.deleteErr = tc.deleteErr

			err := NewDeleteUser(users).Execute(context.Background(), tc.userID)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Execute error = %v, want %v", err, tc.wantErr)
			}
			if err != nil && err.Error() == tc.wantErr.Error() {
				t.Errorf("err = %q, want it wrapped with context", err)
			}
			if len(users.users) != tc.wantLeft {
				t.Errorf("users left = %d, want %d", len(users.users), tc.wantLeft)
			}
		})
	}
}
