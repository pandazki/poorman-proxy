package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"poorman-proxy/secret"
)

// HeaderRewriteFunc defines the function type for header modifications
type HeaderRewriteFunc func(*httputil.ProxyRequest, secret.Secret)

func debugRequest(req *httputil.ProxyRequest) {
	// Debug logging for outbound request
	log.Printf("Outbound Request:\n"+
		"  Method: %s\n"+
		"  URL: %s\n"+
		"  Host: %s\n"+
		"  Headers: %v\n",
		req.Out.Method,
		req.Out.URL.String(),
		req.Out.Host,
		req.Out.Header)
	// Debug logging for inbound request
	log.Printf("Inbound Request:\n"+
		"  Method: %s\n"+
		"  URL: %s\n"+
		"  Host: %s\n"+
		"  Headers: %v\n",
		req.In.Method,
		req.In.URL.String(),
		req.In.Host,
		req.In.Header)

}

func createProxy(targetURL string, pathPrefix string, headerRewrite HeaderRewriteFunc, secretKey secret.Secret) *httputil.ReverseProxy {
	target, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}
	target.Scheme = "https"

	proxy := &httputil.ReverseProxy{
		Rewrite: func(req *httputil.ProxyRequest) {
			req.SetURL(target)
			// Strip the path prefix (e.g., /openai/, /gemini/, /claude/)
			req.Out.URL.Path = strings.TrimPrefix(req.In.URL.Path, pathPrefix)
			req.Out.URL.Scheme = "https"

			// Apply custom header rewrite if provided
			if headerRewrite != nil {
				headerRewrite(req, secretKey)
			}
			// TODO: print the req.Out for debugging
			debugRequest(req)
		},
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

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr: ":8080",
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
