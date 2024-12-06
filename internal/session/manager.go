package session

import (
	"os/exec"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string
	Process   *exec.Cmd
	CreatedAt time.Time
	Expiry    time.Time
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func New() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (m *SessionManager) Create() *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.New().String()
	session := &Session{
		ID:        id,
		Process:   exec.Command("python3", "-iq"),
		CreatedAt: time.Now(),
		Expiry:    time.Now().Add(5 * time.Minute),
	}
	m.sessions[id] = session
	return session
}

func (m *SessionManager) Get(id string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[id]
	return session, exists
}
