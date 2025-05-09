package main

import (
	"net/http"
	"poorman-proxy/secret"
)

// RewritClaudeHeader modifies the request header for Anthropic API.
// here is an example:
//
// #!/bin/sh
//
//	curl https://api.anthropic.com/v1/messages \
//	     --header "x-api-key: $ANTHROPIC_API_KEY" \
//	     --header "anthropic-version: 2023-06-01" \
//	     --header "content-type: application/json" \
//	     --data \
//
//	'{
//	    "model": "claude-3-7-sonnet-20250219",
//	    "max_tokens": 1024,
//	    "messages": [
//	        {"role": "user", "content": "Hello, Claude"}
//	    ]
//	}'
func RewriteClaudeHeader(req *http.Request, key_info secret.Secret) {

	claude_key := key_info.ClaudeKey
	user_key := req.Header.Get("x-api-key")

	found := false
	for _, key := range key_info.ClaudeProxyKey {
		if user_key == key {
			claude_key = key
			found = true
			break
		}
	}

	if !found {
		// reject the request by sending empty authorization
		req.Header.Del("x-api-key")
		req.Header.Set("x-api-key", "")
		return
	}

	// Set the proper Authorization header with Bearer prefix
	req.Header.Del("x-api-key")
	req.Header.Set("x-api-key", claude_key)
}
