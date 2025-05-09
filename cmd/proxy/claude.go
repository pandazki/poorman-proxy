package main

import (
	"net/http/httputil"
	"poorman-proxy/secret"
	"slices"
)

// RewritClaudeHeader modifies the request header for Anthropic API.
// here is an example:
//
// #!/bin/sh

// 	curl https://api.anthropic.com/v1/messages \
// 	     --header "x-api-key: $ANTHROPIC_API_KEY" \
// 	     --header "anthropic-version: 2023-06-01" \
// 	     --header "content-type: application/json" \
// 	     --data \

//	'{
//	    "model": "claude-3-7-sonnet-20250219",
//	    "max_tokens": 1024,
//	    "messages": [
//	        {"role": "user", "content": "Hello, Claude"}
//	    ]
//	}'
func RewriteClaudeHeader(req *httputil.ProxyRequest, key_info secret.Secret) {
	claude_key := key_info.ClaudeKey
	user_key := req.In.Header.Get("x-api-key")

	found := slices.Contains(key_info.ProxyKey, user_key)

	if !found {
		// reject the request by sending empty authorization
		req.Out.Header.Set("x-api-key", "")
		return
	}

	// Set the proper Authorization header with Bearer prefix
	req.Out.Header.Set("x-api-key", claude_key)
}
