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

func TestInvalidSessionID(t *testing.T) {
	manager := New()

	_, exists := manager.Get("non-existent-id")
	if exists {
		t.Error("Expected non-existent session to return exists=false")
	}
}

func TestMemoryLimit(t *testing.T) {
	manager := New()
	session, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	code := `
x = bytearray(120 * 1024 * 1024)  # 120MB
print('allocated')
`
	_, _, err = session.Execute(code)
	if err == nil {
		t.Error("Expected memory limit error")
	}
	if err != ErrExecutionTimeout {
		t.Errorf("Expected ErrExecutionTimeout, got: %v", err)
	}
}

func TestFileSystemAccess(t *testing.T) {
	manager := New()
	session, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	testCases := []struct {
		name string
		code string
	}{
		{
			name: "file read attempt",
			code: "with open('/etc/passwd', 'r') as f: print(f.read())",
		},
		{
			name: "file write attempt",
			code: "with open('test.txt', 'w') as f: f.write('hello')",
		},
		{
			name: "directory listing attempt",
			code: "import os; print(os.listdir('/'))",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, stderr, err := session.Execute(tc.code)
			if err == nil {
				t.Error("Expected permission error")
			}
			if !strings.Contains(stderr, "PermissionError") {
				t.Errorf("Expected PermissionError in stderr, got %q", stderr)
			}
		})
	}
}

func TestNetworkAccess(t *testing.T) {
	manager := New()
	session, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	testCases := []struct {
		name string
		code string
	}{
		{
			name: "http request attempt",
			code: "import urllib.request; urllib.request.urlopen('http://example.com')",
		},
		{
			name: "socket creation attempt",
			code: "import socket; socket.socket()",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, stderr, err := session.Execute(tc.code)
			if err == nil {
				t.Error("Expected network access to be blocked")
			}
			if !strings.Contains(stderr, "PermissionError") {
				t.Errorf("Expected PermissionError in stderr, got %q", stderr)
			}
		})
	}
}

func TestConcurrentExecution(t *testing.T) {
	manager := New()
	session, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Test concurrent access to the same session
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, _, err := session.Execute("print('concurrent')")
			if err != nil {
				t.Errorf("Concurrent execution failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestSessionStateIsolation(t *testing.T) {
	manager := New()

	session1, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session1: %v", err)
	}

	session2, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session2: %v", err)
	}

	// Set variable in session1
	_, _, err = session1.Execute("x = 42")
	if err != nil {
		t.Fatalf("Failed to execute in session1: %v", err)
	}

	// Try to access variable in session2
	stdout, _, err := session2.Execute("try:\n    print(x)\nexcept NameError:\n    print('variable not found')")
	if err != nil {
		t.Fatalf("Failed to execute in session2: %v", err)
	}

	if !strings.Contains(stdout, "variable not found") {
		t.Error("Expected sessions to be isolated")
	}
}

func TestLongRunningSession(t *testing.T) {
	manager := New()
	session, err := manager.Create()
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Execute multiple commands in sequence
	commands := []string{
		"x = 1",
		"y = 2",
		"z = x + y",
		"print(z)",
	}

	for _, cmd := range commands {
		stdout, stderr, err := session.Execute(cmd)
		if err != nil {
			t.Errorf("Failed to execute command %q: %v", cmd, err)
		}
		if stderr != "" {
			t.Errorf("Got unexpected stderr for command %q: %q", cmd, stderr)
		}
		if cmd == "print(z)" && !strings.Contains(stdout, "3") {
			t.Errorf("Expected stdout to contain '3', got %q", stdout)
		}
	}
}
