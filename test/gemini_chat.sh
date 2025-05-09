#!/bin/bash

file="$(dirname "$0")/env.sh"
# shellcheck source=./env.sh
. "${file}"

# curl "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:streamGenerateContent?alt=sse&key=${GEMINI_API_KEY}" \
curl "localhost:8080/gemini/v1beta/models/gemini-1.5-pro:streamGenerateContent?key=${GEMINI_PROXY_KEY}" \
        -H 'Content-Type: application/json' \
        --no-buffer \
        -d '{ "contents":[{"parts":[{"text": "Write long a story about a magic backpack."}]}]}' \
        2> /dev/null