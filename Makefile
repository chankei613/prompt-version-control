.PHONY: run ui build test lint tidy smoke clean

run:
	go run .

ui:
	cd frontend && npm run dev

build:
	wails build

test:
	go test ./...

lint:
	golangci-lint run ./internal/...
	cd frontend && npm run lint

smoke:
	go run ./cmd/smoketest

tidy:
	go mod tidy

clean:
	rm -f prompt-version-control *.db *.db-shm *.db-wal
	rm -rf build/bin
