GOOS = $(shell go env GOOS)
build:
	goreleaser build --id $(GOOS) --single-target --snapshot --rm-dist 
darwin:
	goreleaser build --id darwin --snapshot --rm-dist 
linux:
	goreleaser build --id linux --snapshot --rm-dist 