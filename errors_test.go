package locker_test

import (
	"testing"

	. "github.com/kamilsk/locker"
)

func TestErrors(t *testing.T) {
	if CriticalIssue.Error() != "critical issue" {
		t.Error("unexpected string representation of the error")
		t.FailNow()
	}
	if Interrupted.Error() != "operation interrupted" {
		t.Error("unexpected string representation of the error")
		t.FailNow()
	}
	if InvalidIntent.Error() != "invalid intent" {
		t.Error("unexpected string representation of the error")
		t.FailNow()
	}
}
