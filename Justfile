clean:
  rm -rf gen

gen: clean
  buf generate buf.build/googleapis/googleapis
  buf generate

proto-lint:
  buf lint

proto-breaking:
  buf breaking --against '.git#branch=main'
