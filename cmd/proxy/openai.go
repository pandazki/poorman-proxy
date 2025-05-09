package main

import (
	"net/http/httputil"
	"poorman-proxy/secret"
	"slices"
)

// RewriteOpenAIHeader modifies the request header for OpenAI API.
// here is an example:
//
//	curl "https://api.openai.com/v1/chat/completions" \
//		-H "Content-Type: application/json" \
//		-H "Authorization: Bearer $OPENAI_API_KEY" \
//		-d '{
//			"model": "gpt-4.1",
//			"messages": [
//				{
//					"role": "user",
//					"content": "Write a one-sentence bedtime story about a unicorn."
//				}
//			]
//		}'
func RewriteOpenAIHeader(req *httputil.ProxyRequest, key_info secret.Secret) {
	openai_key := key_info.OpenAIKey
	user_key := req.In.Header.Get("Authorization")

	// Strip "Bearer " prefix if present
	if len(user_key) > 7 && user_key[:7] == "Bearer " {
		user_key = user_key[7:]
	}

	found := slices.Contains(key_info.OpenAIProxyKey, user_key)

	if !found {
		// reject the request by sending empty authorization
		req.Out.Header.Set("Authorization", "")
		return
	}

	// Set the proper Authorization header with Bearer prefix
	req.Out.Header.Set("Authorization", "Bearer "+openai_key)
}
