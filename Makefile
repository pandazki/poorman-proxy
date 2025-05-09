.PHONY: build test debug
build:
	go build -o proxy ./cmd/proxy
	chmod +x proxy

test:
	go test ./...

debug:
	# ./test/openai_chat.sh
	# ./test/claude_chat.sh
	./test/gemini_chat.sh
