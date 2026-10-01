.PHONY: setup lint test pinlint image bench ci audit

VERSION ?= dev
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo unknown)

setup:
	go mod download

lint: pinlint
	@files="$$(gofmt -l $$(go list -f '{{$$d := .Dir}}{{range .GoFiles}}{{$$d}}/{{.}} {{end}}{{range .TestGoFiles}}{{$$d}}/{{.}} {{end}}{{range .XTestGoFiles}}{{$$d}}/{{.}} {{end}}' ./...))"; test -z "$$files" || (echo "$$files"; echo "run gofmt -w"; exit 1)
	go vet ./...

# Fail if any workflow uses an action that is not pinned to a full commit SHA.
pinlint:
	go run ./cmd/pinlint

test:
	go test -race -count=1 ./...

# Local image build (needs Docker); CI does the same in .github/workflows/image.yml.
image:
	docker buildx build --load --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) -t secure-supply-chain:local .

bench: pinlint

# Known vulnerabilities in the code paths actually called (needs Go 1.26+, as in CI).
audit:
	go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...

ci: setup lint test
