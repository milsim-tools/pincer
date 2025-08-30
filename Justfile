git_tag := `git describe --tags --exact-match 2>/dev/null || true`
git_commit := `git rev-parse --short=6 HEAD`
version := if git_tag != "" { git_tag } else { "dev-" + git_commit }

ldflags := "-X main.Version=" + version

clean:
  rm -rf out
  rm -rf pkg/api/gen

gen: clean
  # buf generate buf.build/googleapis/googleapis
  buf generate

proto-lint:
  buf lint

proto-breaking:
  buf breaking --against '.git#branch=main'

build:
  mkdir -p out
  go build \
    -ldflags "{{ ldflags }}" \
    -o out/pincer \
    ./cmd/pincer/main.go

run:
  mkdir -p out
  go build \
    -ldflags "{{ ldflags }}" \
    -o out/pincer \
    ./cmd/pincer/main.go
  ./out/pincer run
