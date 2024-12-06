package session

import (
	"fmt"
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
	sessions   map[string]*Session
	mu         sync.RWMutex
	pythonPath string
}

func New() *SessionManager {
	return &SessionManager{
		sessions:   make(map[string]*Session),
		pythonPath: "python3",
	}
}

func (m *SessionManager) Create() (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.New().String()
	cmd := exec.Command(m.pythonPath, "-iq")
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start Python process: %w", err)
	}
	session := &Session{
		ID:        id,
		Process:   cmd,
		CreatedAt: time.Now(),
		Expiry:    time.Now().Add(5 * time.Minute),
	}
	m.sessions[id] = session
	return session, nil
}

func (m *SessionManager) Get(id string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[id]
	return session, exists
}
