SHELL = /bin/bash -o pipefail

BUMP_VERSION := $(GOPATH)/bin/bump_version
MEGACHECK := $(GOPATH)/bin/megacheck

deps:
	go get -u ./...

test: lint
	go test ./...

lint: | $(MEGACHECK)
	go vet ./...
	$(MEGACHECK) ./...

$(MEGACHECK):
	go get honnef.co/go/tools/cmd/megacheck

race-test: lint
	go test -race ./...

$(BUMP_VERSION):
	go get github.com/kevinburke/bump_version

release: race-test | $(BUMP_VERSION)
	$(BUMP_VERSION) minor client.go
