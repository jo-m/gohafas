.PHONY: lint check format

lint:
	gofmt -l .; test -z "$$(gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000 ./...
	go run github.com/securego/gosec/v2/cmd/gosec@latest ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

check: lint

format:
	gofmt -w -s .