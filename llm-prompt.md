You are helping develop a Python code execution sandbox server called PyServer using Go. Here's the current context and progress:

CURRENT IMPLEMENTATION:
```go
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
```

CURRENT TESTS:
```go
package session

import (
    "testing"
    "time"
)

func TestCreate(t *testing.T) {
    manager := New()
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
    session := manager.Create()
    _, exists := manager.Get(session.ID)
    if !exists {
        t.Error("Expected session to exist")
    }
}
```

PROJECT REQUIREMENTS:
1. Resource Constraints:
   - Process memory limit
   - Session lifetime: 5 minutes (implemented)
   - Individual code execution timeout: 2 seconds

2. Session Management:
   - UUIDs for session identification (implemented)
   - Complete variable state within sessions
   - Multiple concurrent sessions with isolation (implemented)
   - Sessions persist for 5 minutes (implemented)
   - Python standard library only

3. API Design:
   - Endpoint: /execute
   - No immediate authentication
   - JSON responses with:
     - session_id (always present)
     - stdout (successful execution)
     - stderr (code errors)
     - error (system errors like timeout, memory limit)

CURRENT PROGRESS:
- Implemented basic session management with UUIDs
- Implemented concurrent access safety with RWMutex
- Basic session creation and retrieval
- Session expiry time set to 5 minutes
- Basic test coverage

NEXT STEPS NEEDED:
1. Implement Python process management:
   - Start the interpreter process
   - Handle stdin/stdout/stderr pipes
   - Implement process cleanup
   - Add error handling for process operations

2. Implement code execution:
   - Method to execute code in a session
   - Capture output
   - Implement 2-second timeout
   - Handle memory limits

3. Implement HTTP server using Gin:
   - Set up routes
   - Implement /execute endpoint
   - Add error handling
   - Structure JSON responses

We are following Test-Driven Development (TDD) practices. The immediate next step is to implement proper Python process management and add error handling to the Create method.

Please help continue the development from this point, focusing on proper error handling and Python process management. Please explain your reasoning and provide tests first, following TDD principles.
