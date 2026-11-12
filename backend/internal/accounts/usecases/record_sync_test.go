package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestRecordSync_MarksTheAccountWithTimeAndErrorKey(t *testing.T) {
	tests := []struct {
		name   string
		errKey shared.MessageKey
	}{
		{"failed sync", "emailsync.mail_auth_failed"},
		{"clean sync", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accounts := &fakeAccounts{}

			if err := NewRecordSync(accounts).Execute(context.Background(), 7, testNow, tt.errKey); err != nil {
				t.Fatalf("Execute error = %v, want nil", err)
			}

			want := []syncMark{{id: 7, at: testNow, errKey: tt.errKey}}
			if !reflect.DeepEqual(accounts.marked, want) {
				t.Errorf("marked %+v, want %+v", accounts.marked, want)
			}
		})
	}
}

func TestRecordSync_AccountRemovedDuringSyncIsIgnored(t *testing.T) {
	accounts := &fakeAccounts{markErr: shared.ErrNotFound}

	err := NewRecordSync(accounts).Execute(context.Background(), 7, testNow, "")

	if err != nil {
		t.Errorf("Execute error = %v, want nil", err)
	}
	if len(accounts.marked) != 1 {
		t.Errorf("MarkSynced called %d times, want 1", len(accounts.marked))
	}
}

func TestRecordSync_RepositoryErrorIsWrapped(t *testing.T) {
	cause := errors.New("boom")
	accounts := &fakeAccounts{markErr: cause}

	err := NewRecordSync(accounts).Execute(context.Background(), 7, testNow, "")

	assertWrapped(t, err, cause, "record account sync")
}
