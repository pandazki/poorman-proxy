#!/bin/sh

file="$(dirname "$0")/env.sh"
# shellcheck source=./env.sh
. "${file}"

# curl https://api.anthropic.com/v1/messages \
curl localhost:8080/claude/v1/messages \
        --header "x-api-key: $CLAUDE_PROXY_KEY" \
        --header "anthropic-version: 2023-06-01" \
        --header "content-type: application/json" \
        --data \
'{
    "model": "claude-3-7-sonnet-20250219",
    "max_tokens": 1024,
    "messages": [
        {"role": "user", "content": "Hello, Claude"}
    ]
}'