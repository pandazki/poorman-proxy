.PHONY: build test debug
build:
	go build -o proxy ./cmd/proxy
	chmod +x proxy

test:
	go test ./...

debug: build
	./test/openai_chat.sh
