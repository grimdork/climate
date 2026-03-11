.PHONY: test lint tidy vet fmt ci check

# Run unit tests
test:
	go test ./... -v

# Run go vet across all packages
vet:
	go vet ./...

# Run staticcheck (requires staticcheck in PATH)
lint:
	@command -v staticcheck >/dev/null 2>&1 && staticcheck ./... || echo "staticcheck not found; skipping lint"

# Format sources
fmt:
	gofmt -w .

# Tidy the module and verify
tidy:
	go mod tidy
	go mod verify

# Run the standard CI checks locally
ci: tidy vet lint test fmt_check

# Run a suite of checks useful during development
check: vet lint test fmt_check

# Check formatting without modifying files (fails if formatting needed)
fmt_check:
	if [ -n "$(gofmt -l .)" ]; then echo "gofmt needs to be run"; gofmt -l .; exit 1; fi
