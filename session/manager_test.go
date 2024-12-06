package session

import (
	"strings"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	manager := New()

	session, err := manager.Create()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if session.Process == nil {
		t.Error("Expected running Python process")
	}

	if manager.sessions[session.ID] != session {
		t.Error("Expected session to be stored in sessions map")
	}

	if time.Since(session.CreatedAt) > time.Second {
		t.Error("Expected CreatedAt to be recent")
	}
}

func TestGet(t *testing.T) {
	manager := New()

	session, err := manager.Create()

	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	got, exists := manager.Get(session.ID)
	if !exists {
		t.Error("Expected session to exist")
	}
	if got != session {
		t.Error("Expected to get same session instance")
	}
}

func TestExecuteBasic(t *testing.T) {
	manager := New()
	session, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	stdout, stderr, err := session.Execute("print('hello')")
	if err != nil {
		t.Fatalf("Failed to execute code: %v", err)
	}

	if !strings.Contains(stdout, "hello") {
		t.Errorf("Expected stdout to contain 'hello', got %q", stdout)
	}
	if stderr != "" {
		t.Errorf("Expected empty stderr, got %q", stderr)
	}
}

func TestExecuteError(t *testing.T) {
	manager := New()
	session, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	_, stderr, err := session.Execute("x = 1/0")
	if err != nil {
		t.Fatalf("Expected no error return, got: %v", err)
	}

	if !strings.Contains(stderr, "ZeroDivisionError") {
		t.Errorf("Expected stderr to contain ZeroDivisionError, got %q", stderr)
	}
}

func TestExecuteTimeout(t *testing.T) {
	manager := New()
	session, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	_, _, err = session.Execute("while True: pass")

	if err == nil {
		t.Error("Expected timeout error")
	}
	if err != ErrExecutionTimeout {
		t.Errorf("Expected ErrExecutionTimeout, got: %v", err)
	}
}

func TestExecuteMultiline(t *testing.T) {
	manager := New()
	session, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	code := `
	x = 5
	y = 10
	print(x + y)
	`
	stdout, stderr, err := session.Execute(code)
	if err != nil {
		t.Fatalf("Failed to execute multiline code: %v", err)
	}

	if !strings.Contains(stdout, "15") {
		t.Errorf("Expected stdout to contain '15', got %q", stdout)
	}
	if stderr != "" {
		t.Errorf("Expected empty stderr, got %q", stderr)
	}
}

func TestExecuteStatePreservation(t *testing.T) {
	manager := New()
	session, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	_, _, err = session.Execute("x = 42")
	if err != nil {
		t.Fatalf("Failed to execute first code: %v", err)
	}

	stdout, _, err := session.Execute("print(x)")
	if err != nil {
		t.Fatalf("Failed to execute second code: %v", err)
	}

	if !strings.Contains(stdout, "42") {
		t.Errorf("Expected stdout to contain '42', got %q", stdout)
	}
}
