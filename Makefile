build:
	go build ./...

test:
	go test ./...

run-isolated:
	env GOOS=linux GOARCH=amd64 go build
	docker build . -t cbv2
	docker run --rm -it cbv2 bash

release:
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined. Token needs to have full repo permissions)
endif
	curl -sL https://git.io/goreleaser | bash
