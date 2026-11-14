package usecases

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/code-rcplaza/mail-checker/backend/internal/auth/domain"
	"github.com/code-rcplaza/mail-checker/backend/internal/auth/ports"
	shared "github.com/code-rcplaza/mail-checker/backend/internal/shared/domain"
)

type CodeDeliveryDependencies struct {
	Codes  ports.LoginCodes
	Sends  ports.CodeSends
	Sender ports.CodeSender
	Random io.Reader
	Now    func() time.Time
}

// CodeDelivery is the part shared by register, login and resend: limit the
// requests of an email, issue a code and send it.
type CodeDelivery struct {
	dependencies CodeDeliveryDependencies
}

func NewCodeDelivery(dependencies CodeDeliveryDependencies) CodeDelivery {
	return CodeDelivery{dependencies: dependencies}
}

// admitRequest counts the request whether or not the email has an account, so
// the limit answers the same for both.
func (d CodeDelivery) admitRequest(ctx context.Context, email string) error {
	now := d.dependencies.Now()
	sends, err := d.dependencies.Sends.List(ctx, email)
	if err != nil {
		return fmt.Errorf("read code sends: %w", err)
	}
	if !domain.CanSendCode(sends, now) {
		return shared.TooMany(domain.KeyTooManySends, map[string]any{"limit": domain.SendLimit})
	}
	if err := d.dependencies.Sends.Record(ctx, domain.CodeSend{Email: email, SentAt: now}); err != nil {
		return fmt.Errorf("record code send: %w", err)
	}
	return nil
}

// sendCode issues a new code for the user and sends it to their email.
func (d CodeDelivery) sendCode(ctx context.Context, user domain.User) error {
	code, err := d.issueCode(ctx, user.ID)
	if err != nil {
		return err
	}
	if err := d.dependencies.Sender.Send(ctx, user.Email, code); err != nil {
		return fmt.Errorf("send login code: %w", err)
	}
	return nil
}

// issueCode stores the hash, replacing the previous code, and returns the code
// in clear only so it can be sent.
func (d CodeDelivery) issueCode(ctx context.Context, userID int64) (string, error) {
	issued, err := domain.IssueLoginCode(d.dependencies.Random, domain.CodeIssue{UserID: userID, Now: d.dependencies.Now()})
	if err != nil {
		return "", err
	}
	if err := d.dependencies.Codes.Save(ctx, issued.LoginCode); err != nil {
		return "", fmt.Errorf("save login code: %w", err)
	}
	return issued.Code, nil
}
