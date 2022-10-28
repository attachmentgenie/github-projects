build:
	goreleaser build --id $(shell go env GOOS) --single-target --snapshot --rm-dist
darwin:
	goreleaser build --id darwin --snapshot --rm-dist
linux:
	goreleaser build --id linux --snapshot --rm-dist
snapshot:
	goreleaser release --snapshot --rm-dist
release:
	git tag $(shell svu next)
	git push --tags
	goreleaser --rm-dist
