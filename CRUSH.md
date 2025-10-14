# GO-URING Development Guidelines

## Build & Test Commands
- **Test all**: `go test ./...`
- **Test single package**: `go test ./uring` or `go test ./reactor`
- **Test with race detection**: `go test -race ./...`
- **Test with coverage**: `go test -coverprofile=coverage.out ./...`
- **Test single test**: `go test -run TestFunctionName ./package`
- **Lint**: `golangci-lint run -D deadcode -D varcheck -D structcheck -D unused`
- **Build**: `go build ./...`
- **Performance mode**: Use `-tags amd64_atomic` for 1-3% performance gain

## Code Style Guidelines
- **Build tags**: Always include `//go:build linux` for Linux-specific code
- **Package naming**: Use short, lowercase names (uring, reactor, net)
- **Imports**: Group stdlib, third-party, local packages with blank lines
- **Error handling**: Return errors, use `noErr()` helper in examples
- **Types**: Use explicit types for constants and operations (OpCode uint8)
- **Concurrency**: Use channels for async communication, atomic operations for counters
- **Memory**: Use unsafe.Pointer carefully, prefer mmap for ring buffers
- **Naming**: 
  - Public types: PascalCase (Ring, Reactor, CQEvent)
  - Private fields: camelCase (kHead, sqeTail)
  - Constants: camelCase with prefixes (opReadFixed, sqCQOverflow)
- **Comments**: Use godoc format, explain complex io_uring concepts
- **Testing**: Use testify for assertions, separate files for different operation types