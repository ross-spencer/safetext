# CLI tools

# Help
help:
    @echo "Command line helpers for this project.\n"
    @just -l

# Run all pre-commit checks
all-checks:
	pre-commit run --all-files

# Fix imports
imports:
  goimports -w ./..

# Go fmt
fmt:
  go fmt ./...

# Go vet
vet:
  go vet ./...

# Snapshot
snapshot:
  goreleaser build --snapshot --single-target --clean -f .goreleaser.yml

# Snapshot (verboise)
snapshot-verbose:
  goreleaser build --verbose --snapshot --single-target --clean -f .goreleaser.yml

# Release
release:
  goreleaser release --skip=publish --clean -f .goreleaser.yml

# Single-target release
target:
  goreleaser build --single-target --clean -f .goreleaser.yml

# Test linux dist.
run-linux *FLAGS:
  ./dist/build-safetext_linux_amd64_v1/safetext {{FLAGS}}

# update safetext.json
update-safetext:
  wget https://github.com/ross-spencer/safetext-json/blob/main/safetext-ucd.json -O ./pkg/safetext/safetext.json

# docs
docs:
  godoc -http=localhost:6060
