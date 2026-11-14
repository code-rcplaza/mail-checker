package sqlite

import (
	"context"
	"errors"
	"testing"
)

func TestIsUniqueViolation(t *testing.T) {
	db := openTempDB(t)
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `CREATE TABLE things (name TEXT NOT NULL UNIQUE)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO things (name) VALUES ('a')`); err != nil {
		t.Fatalf("first insert: %v", err)
	}
	_, duplicate := db.ExecContext(ctx, `INSERT INTO things (name) VALUES ('a')`)
	_, missing := db.ExecContext(ctx, `INSERT INTO things (name) VALUES (NULL)`)
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"a duplicate is a unique violation", duplicate, true},
		{"a wrapped duplicate too", errors.Join(errors.New("create"), duplicate), true},
		{"a NOT NULL violation is not", missing, false},
		{"no error is not", nil, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsUniqueViolation(tc.err); got != tc.want {
				t.Errorf("IsUniqueViolation(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
