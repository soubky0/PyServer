# PyServer

A Python code execution sandbox server written in Go.

## Requirements

### Resource Constraints
- [ ] Process memory limit
- [ ] Session lifetime: 5 minutes
- [ ] Individual code execution timeout: 2 seconds

### Session Management
- [x] UUIDs for session identification
- [x] Multiple concurrent sessions with isolation
- [x] Sessions persist for 5 minutes
- [x] Python standard library only (no external packages)

### API Design
- [ ] Endpoint: /execute
- [ ] No immediate authentication (but designed for future addition)
- [ ] JSON responses with:
  - [ ] session_id (always present)
  - [ ] stdout (successful execution)
  - [ ] stderr (code errors)
  - [ ] error (system errors like timeout, memory limit)

## Current Implementation

### Completed Features
1. Basic session management
   - Session creation with UUID
   - Session retrieval by ID
   - Concurrent access handling with RWMutex
   - 5-minute session expiry time

### Next Steps
1. Process Management
   - Start Python interpreter process
   - Handle stdin/stdout/stderr pipes
   - Process cleanup on session expiry

2. Code Execution
   - Implement code execution in session
   - Capture output
   - Implement timeout mechanism

3. HTTP Server
   - Set up Gin framework
   - Implement /execute endpoint
   - Add error handling

4. Testing
   - Add more unit tests
   - Add integration tests
   - Add API tests
   
## Development Progress
- [x] Basic session management structure
- [x] Session creation and retrieval
- [x] Concurrent access safety
- [ ] Process management
- [ ] Code execution
- [ ] HTTP API
- [ ] Resource constraints
