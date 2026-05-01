package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfirmDeployment_Yes(t *testing.T) {
	var out bytes.Buffer
	ok, err := confirmDeployment(strings.NewReader("yes\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected confirmation to pass")
	}
	if !strings.Contains(out.String(), "Continue with deployment") {
		t.Errorf("expected confirmation prompt, got %q", out.String())
	}
}

func TestConfirmDeployment_CaseInsensitive(t *testing.T) {
	ok, err := confirmDeployment(strings.NewReader("YES\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected uppercase YES to pass")
	}
}

func TestConfirmDeployment_No(t *testing.T) {
	ok, err := confirmDeployment(strings.NewReader("no\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected non-yes response to cancel")
	}
}

func TestConfirmDeployment_EOF(t *testing.T) {
	ok, err := confirmDeployment(strings.NewReader("yes"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected yes without trailing newline to pass")
	}
}
