package session

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionExpired   = errors.New("session expired")
	ErrExecutionTimeout = errors.New("execution timeout exceeded")
)

type Session struct {
	ID        string
	Process   *exec.Cmd
	CreatedAt time.Time
	Expiry    time.Time
	Stdin     io.WriteCloser
	Stdout    io.ReadCloser
	Stderr    io.ReadCloser
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
	cmd := exec.Command("bash", "-c", "ulimit -v 102400 && python3 -iq")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start Python process: %w", err)
	}
	session := &Session{
		ID:        id,
		Process:   cmd,
		CreatedAt: time.Now(),
		Expiry:    time.Now().Add(5 * time.Minute),
		Stdin:     stdin,
		Stdout:    stdout,
		Stderr:    stderr,
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

func parseStdout(output string) string {
	output = strings.ReplaceAll(output, "__END__\n", "")
	output = strings.ReplaceAll(output, "__END__", "")
	return output
}

func parseStderr(output string) string {
	output = strings.ReplaceAll(output, ">>> ", "")
	output = strings.ReplaceAll(output, "... ", "")
	return strings.TrimSpace(output)
}

func (s *Session) Execute(code string) (stdout string, stderr string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	codeWithMarker := code + "\nprint('__END__')"
	if _, err := s.Stdin.Write([]byte(codeWithMarker + "\n")); err != nil {
		return "", "", fmt.Errorf("failed to write to stdin: %w", err)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutTee := io.TeeReader(s.Stdout, &stdoutBuf)
	stderrTee := io.TeeReader(s.Stderr, &stderrBuf)

	done := make(chan error, 1)
	go func() {
		buf := make([]byte, 1024)
		for {
			_, err := stdoutTee.Read(buf)
			if err != nil {
				done <- err
				return
			}
			if strings.Contains(stdoutBuf.String(), "__END__") {
				done <- nil
				return
			}
		}
	}()

	go func() {
		_, err := io.Copy(io.Discard, stderrTee)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return "", "", fmt.Errorf("error reading output: %w", err)
		}
	case <-ctx.Done():
		return "", "", ErrExecutionTimeout
	}

	stdout = parseStdout(stdoutBuf.String())
	stderr = parseStderr(stderrBuf.String())

	return stdout, stderr, nil
}
