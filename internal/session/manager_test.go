package session

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	// Initialize the session manager
	manager := New()

	// Create a new session
	session := manager.Create()

	if manager.sessions[session.ID] != session {
		t.Error("Expected session to be stored in sessions map")
	}
	if session.Process == nil {
		t.Error("Expected running Python process")
	}

	if time.Since(session.CreatedAt) > time.Second {
		t.Error("Expected CreatedAt to be recent")
	}
}

func TestGet(t *testing.T) {
	manager := New()

	// Create a new session
	session := manager.Create()

	_, exists := manager.Get(session.ID)

	if !exists {
		t.Error("Expected session to exist")
	}
}
