package proxy

import (
	"net/http"
	"poorman-proxy/secret"
)

// RewriteGeminiRequest modifies the request header for Gemini API
// here is an exmample
//
//	!curl "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:streamGenerateContent?alt=sse&key=${GEMINI_API_KEY}" \
//			-H 'Content-Type: application/json' \
//			--no-buffer \
//			-d '{ "contents":[{"parts":[{"text": "Write long a story about a magic backpack."}]}]}' \
//			2> /dev/null
func RewriteGeminiRequest(req *http.Request, key_info secret.Secret) {
	gemini_key := key_info.GeminiKey
	user_query := req.URL.Query()
	user_key := user_query.Get("key")

	found := false
	for _, key := range key_info.GeminiProxyKey {
		if user_key == key {
			gemini_key = key
			found = true
			break
		}
	}
	if !found {
		// reject the request by sending empty key
		user_query.Del("key")
		user_query.Set("key", "")
		return
	}
	user_query.Del("key")
	user_query.Set("key", gemini_key)
	return
}
