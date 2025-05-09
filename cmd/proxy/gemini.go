package main

import (
	"net/http/httputil"
	"poorman-proxy/secret"

	"golang.org/x/exp/slices"
)

// RewriteGeminiRequest modifies the request header for Gemini API
// here is an exmample
//
//	!curl "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:streamGenerateContent?alt=sse&key=${GEMINI_API_KEY}" \
//			-H 'Content-Type: application/json' \
//			--no-buffer \
//			-d '{ "contents":[{"parts":[{"text": "Write long a story about a magic backpack."}]}]}' \
//			2> /dev/null
func RewriteGeminiRequest(req *httputil.ProxyRequest, key_info secret.Secret) {
	gemini_key := key_info.GeminiKey
	user_query := req.In.URL.Query()
	user_key := user_query.Get("key")

	found := slices.Contains(key_info.GeminiProxyKey, user_key)

	if !found {
		// reject the request by sending empty key
		q := req.Out.URL.Query()
		q.Set("key", "")
		req.Out.URL.RawQuery = q.Encode()
		return
	}
	q := req.Out.URL.Query()
	q.Set("key", gemini_key)
	req.Out.URL.RawQuery = q.Encode()
}
