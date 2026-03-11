.PHONY: test lint tidy vet fmt

test:
	go test ./... -v

vet:
	go vet ./...

lint:
	staticcheck ./...

fmt:
	gofmt -w .

# Tidy the module and verify
tidy:
	go mod tidy
	go mod verify
