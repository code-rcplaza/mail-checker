package usecases

import (
	"context"
	"errors"
	"testing"

	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

func TestGetMe(t *testing.T) {
	got, err := NewGetMe(usersWith(ana)).Execute(context.Background(), ana.ID)
	if err != nil || got != ana {
		t.Errorf("GetMe = %+v, %v; want %+v", got, err, ana)
	}

	_, err = NewGetMe(usersWith(ana)).Execute(context.Background(), ana.ID+1)
	if !errors.Is(err, shared.ErrNotFound) {
		t.Errorf("GetMe of an unknown user error = %v, want ErrNotFound", err)
	}
}
