package session

import (
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
