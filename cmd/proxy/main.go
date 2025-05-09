package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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

func createProxy(targetURL string, pathPrefix string, headerRewrite HeaderRewriteFunc, secretKey secret.Secret, outboundProxyURL string) *httputil.ReverseProxy {
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

	// Configure outbound proxy for the transport
	transport := &http.Transport{}
	if outboundProxyURL != "" {
		proxyURL, err := url.Parse(outboundProxyURL)
		if err != nil {
			log.Printf("Error parsing outbound proxy URL '%s': %v. Proceeding without outbound proxy.", outboundProxyURL, err)
		} else {
			transport.Proxy = http.ProxyURL(proxyURL)
			log.Printf("Using outbound proxy: %s", proxyURL.String())
		}
	}
	// If outboundProxyURL is empty or invalid, transport.Proxy will remain nil,
	// and http.DefaultTransport (which uses environment variables like HTTP_PROXY) will be effectively used by default by the ReverseProxy.
	// To ensure our explicit proxy setting (or lack thereof) is used, we always set the transport.
	// If no proxy is set, it will use a new transport with no proxy, overriding env vars.
	// If you want to respect HTTP_PROXY, HTTPS_PROXY, NO_PROXY env vars when outboundProxyURL is not set, use http.DefaultTransport.
	// For this specific feature, we want to explicitly control the proxy via the command-line flag or have no proxy if not specified there.
	if outboundProxyURL == "" {
		log.Println("No outbound proxy URL provided. Using direct connection.")
		// Ensure no proxy is used if not specified, effectively overriding environment variables.
		transport.Proxy = nil // Explicitly set to nil to override environment proxy settings
	} // else, transport.Proxy is already set if outboundProxyURL was valid, or nil if parsing failed (with a log message)

	proxy.Transport = transport

	return proxy
}

func main() {
	outboundProxyURLFlag := flag.String("outbound-proxy-url", "", "Optional. URL of the outbound proxy server (e.g., http://user:pass@host:port or socks5://user:pass@host:port).")
	flag.Parse()

	outboundProxyURL := *outboundProxyURLFlag
	if envURL, ok := os.LookupEnv("OUTBOUND_PROXY_URL"); ok && envURL != "" {
		outboundProxyURL = envURL
		log.Printf("Using outbound proxy URL from environment variable OUTBOUND_PROXY_URL: %s", outboundProxyURL)
	} else if outboundProxyURL != "" {
		log.Printf("Using outbound proxy URL from command-line flag: %s", outboundProxyURL)
	} else {
		log.Println("No outbound proxy URL specified via command-line flag or environment variable.")
	}

	secretKey := secret.Load()
	// Create proxies with their respective header rewrite functions
	openaiProxy := createProxy("https://api.openai.com", "/openai", RewriteOpenAIHeader, secretKey, outboundProxyURL)
	geminiProxy := createProxy("https://generativelanguage.googleapis.com", "/gemini", RewriteGeminiRequest, secretKey, outboundProxyURL)
	claudeProxy := createProxy("https://api.anthropic.com", "/claude", RewriteClaudeHeader, secretKey, outboundProxyURL)

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
