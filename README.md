# PyServer

A Python code execution sandbox server written in Go.

## Requirements

### Resource Constraints
- [ ] Process memory limit
- [x] Session lifetime: 5 minutes
- [x] Individual code execution timeout: 2 seconds

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
1. Session Management
   - Session creation with UUID
   - Session retrieval by ID
   - Concurrent access handling with RWMutex
   - 5-minute session expiry time

2. Process Management
   - Python interpreter process management
   - Stdin/stdout/stderr pipe handling
   - 2-second execution timeout

3. Code Execution
   - Execute code in interactive Python session
   - Separate stdout/stderr capture
   - State preservation between executions
   - Output cleaning (removal of Python prompts)

### Next Steps
1. HTTP Server
   - Set up Gin framework
   - Implement /execute endpoint
   - Add error handling
   - Structure JSON responses

2. Resource Management
   - Implement process memory limits
   - Add process cleanup on session expiry

3. Testing
   - Add integration tests
   - Add API tests

## Development Progress
- [x] Basic session management structure
- [x] Session creation and retrieval
- [x] Concurrent access safety
- [x] Process management
- [x] Code execution
- [x] Execution timeout
- [ ] HTTP API
- [ ] Resource constraints