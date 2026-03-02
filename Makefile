build:
	go build -o basalt ./cmd/basalt

test:
	go test ./...

lint:
	golangci-lint run

dev:
	go run ./cmd/basalt dev

release:
	goreleaser release --snapshot --rm-dist
