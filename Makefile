
test:
	go test -v ./...

generate:
	go generate ./...

build:
	go build -o t8 cmd/t8/main.go

test-release:
	goreleaser --snapshot --skip-publish --rm-dist

release:
	goreleaser --rm-dist