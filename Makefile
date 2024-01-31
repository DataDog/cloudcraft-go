.POSIX:
.SUFFIXES:

GO = go
GIT = git
RM = rm

-include .env
export

all: pre-commit

pre-commit: tidy fmt lint vulnerabilities test clean # Runs all pre-commit checks.

commit: pre-commit # Commits the changes to the repository.
	$(GIT) commit -s

doc: # Serve the documentation locally.
	$(GO) run golang.org/x/tools/cmd/godoc@latest -http=localhost:1967

tidy: # Updates the go.mod file to use the latest versions of all direct and indirect dependencies.
	$(GO) mod tidy

fmt: # Formats Go source files in this repository.
	$(GO) run mvdan.cc/gofumpt@latest -e -extra -w .

lint: # Runs golangci-lint using the config at the root of the repository.
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...

vulnerabilities: # Analyzes the codebase and looks for vulnerabilities affecting it.
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

test: # Runs unit tests.
	$(GO) test -short -cover -race -vet all -mod readonly ./...

test/integration: # Runs integration tests.
	$(GO) test -cover -race -vet all -mod readonly ./tests/integration

test/coverage: # Generates a coverage profile and open it in a browser.
	$(GO) test -short -coverprofile cover.out ./...
	$(GO) tool cover -html=cover.out

clean: # Cleans cache files from tests and deletes any build output.
	$(RM) -f cover.out

.PHONY: all pre-commit commit doc tidy fmt lint vulnerabilities test test/integration test/coverage clean
