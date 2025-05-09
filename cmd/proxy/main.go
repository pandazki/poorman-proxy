package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"poorman-proxy/secret"
)

// HeaderRewriteFunc defines the function type for header modifications
type HeaderRewriteFunc func(*http.Request, secret.Secret)

func createProxy(targetURL string, pathPrefix string, headerRewrite HeaderRewriteFunc, secretKey secret.Secret) *httputil.ReverseProxy {
	target, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Strip the path prefix (e.g., /openai/, /gemini/, /claude/)
		req.URL.Path = strings.TrimPrefix(req.URL.Path, pathPrefix)

		// // Ensure PUT method
		// if req.Method != http.MethodPut {
		// 	req.Method = http.MethodPut
		// }

		// Apply custom header rewrite if provided
		if headerRewrite != nil {
			headerRewrite(req, secretKey)
		}
	}

	return proxy
}

func main() {

	secretKey := secret.Load()
	// Create proxies with their respective header rewrite functions
	openaiProxy := createProxy("https://api.openai.com", "/openai", RewriteOpenAIHeader, secretKey)
	geminiProxy := createProxy("https://generativelanguage.googleapis.com", "/gemini", RewriteGeminiRequest, secretKey)
	claudeProxy := createProxy("https://api.anthropic.com", "/claude", RewriteClaudeHeader, secretKey)

	// Route handlers
	http.HandleFunc("/openai/", func(w http.ResponseWriter, r *http.Request) {
		openaiProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/gemini/", func(w http.ResponseWriter, r *http.Request) {
		geminiProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/claude/", func(w http.ResponseWriter, r *http.Request) {
		claudeProxy.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr: ":8080",
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
