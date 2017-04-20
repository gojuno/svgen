PKG=github.com/gojuno/svgen

all: test

generate:
	go run svgen.go -i ${PKG}/tests

test: generate
	go test ${PKG}/tests
