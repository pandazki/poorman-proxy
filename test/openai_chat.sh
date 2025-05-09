#!/bin/bash

file="$(dirname "$0")/env.sh"
# shellcheck source=../.env
source "${file}"

curl "localhost:8080/openai/v1/chat/completions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $OPENAI_PROXY_KEY" \
    -d '{
        "model": "gpt-4.1",
        "messages": [
            {
                "role": "user",
                "content": "Write a one-sentence bedtime story about a unicorn."
            }
        ]
    }'