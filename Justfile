export CGO_ENABLED := "0"

git_tag := `git describe --tags --exact-match 2>/dev/null || true`
git_commit := `git rev-parse --short=6 HEAD`
version := if git_tag != "" { git_tag } else { "dev-" + git_commit }

ldflags := "-X main.Version=" + version

_clean-gen:
  rm -rf pkg/api/gen

# Generate API bindings for the protobuf definitions
[group("proto")]
gen: _clean-gen
  # buf generate buf.build/googleapis/googleapis
  buf generate

# Lint the protobuf definitions
[group("proto")]
proto-lint:
  buf lint

# Check the protobuf definitions for breaking changes against `main`
[group("proto")]
proto-breaking:
  buf breaking --against '.git#branch=main'

_clean:
  rm -rf out

# Compile the output binary
[group("go")]
build: _clean
  mkdir -p out
  go build \
    -ldflags "{{ ldflags }}" \
    -o out/pincer \
    ./cmd/pincer/main.go

# Run the application (use environment variables to configure the app)
[group("go")]
run: _clean
  mkdir -p out
  go build \
    -ldflags "{{ ldflags }}" \
    -o out/pincer \
    ./cmd/pincer/main.go
  ./out/pincer run

# Format the Go code.
[group("go")]
fmt:
  go fmt ./...

# Lint the Go code.
[group("go")]
lint:
  go vet ./...

# Compile everything without creating artifacts to check the codebase builds.
[group("go")]
check:
  go build ./...

# Run application tests
[group("go")]
test:
	go clean -testcache && go test -v ./...

dotenv:
  if [ ! -f ".env" ]; then cp .env.sample .env; fi
