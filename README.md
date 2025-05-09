# Poorman Proxy

Poorman Proxy is a simple HTTP reverse proxy server written in Go. It allows you to route requests to various AI service providers (OpenAI, Google Gemini, Anthropic Claude) through a single endpoint, managing API key authentication for each service. It also supports routing its own outgoing connections through an external HTTP/SOCKS5 proxy (configurable via command-line argument) if direct access to AI services is restricted.

## Features

*   Proxies requests to:
    *   OpenAI API (`https://api.openai.com`) via `/openai/`
    *   Google Gemini API (`https://generativelanguage.googleapis.com`) via `/gemini/`
    *   Anthropic Claude API (`https://api.anthropic.com`) via `/claude/`
*   Manages API keys and proxy authorization keys securely using an embedded `secret.json` file.
*   Supports an outbound proxy (HTTP/SOCKS5) for its own connections to AI services, **configurable via command-line argument**.
*   Listens on port `8080` by default.
*   Provides a `/health` endpoint for health checks.

## Prerequisites

*   Go version 1.23.4 or higher.

## Setup

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd poorman-proxy
    ```

2.  **Configure Secrets:**
    Create a `secret.json` file in the `secret/` directory with your API keys:
    ```json
    {
      "openai_api_key": "sk-your-openai-api-key",
      "gemini_api_key": "your-gemini-api-key",
      "claude_api_key": "sk-ant-your-claude-api-key",
      "openai_proxy_key": ["optional-proxy-auth-key-for-openai"],
      "gemini_proxy_key": ["optional-proxy-auth-key-for-gemini"],
      "claude_proxy_key": ["optional-proxy-auth-key-for-claude"]
      // Note: Outbound proxy is now configured via command-line argument, not in this file.
    }
    ```
    The `*_proxy_key` fields are optional and can be used if you want to add an extra layer of authorization to access the proxy itself. If configured, requests to the proxy for a specific service must include one of the corresponding `*_proxy_key` values in the `Authorization` header (e.g., `Authorization: Bearer your-proxy-key`).

    **Important:** The `secret/secret.json` file is gitignored by default to prevent accidental commitment of sensitive information.

## Usage

### Build

To build the proxy executable:

```bash
make build
```

This will create an executable file named `proxy` in the project root.

### Run

To run the proxy server:

```bash
./proxy
```

To run with an outbound proxy for all AI service requests:

```bash
./proxy -outbound-proxy-url="http://user:password@your-proxy-server.com:port"
```
Or for a SOCKS5 proxy:
```bash
./proxy -outbound-proxy-url="socks5://user:password@your-socks-proxy.com:port"
```
If the proxy URL is invalid or parsing fails, Poorman Proxy will log an error and proceed without an outbound proxy (direct connection).
If no `-outbound-proxy-url` is provided, it will also use a direct connection.

The server will start listening on `http://localhost:8080`.

### Configuring an Outbound Proxy

If the machine running Poorman Proxy cannot directly access the AI service APIs (e.g., due to network restrictions), you can configure Poorman Proxy to use an intermediary proxy for its outgoing connections by providing the `-outbound-proxy-url` command-line argument when starting the server.

Examples:
*   HTTP Proxy: `./proxy -outbound-proxy-url="http://username:password@proxy.example.com:8080"`
*   SOCKS5 Proxy: `./proxy -outbound-proxy-url="socks5://username:password@proxy.example.com:1080"`

If your proxy does not require authentication, you can omit `username:password@` from the URL.
If the argument is not provided, Poorman Proxy will attempt direct connections to the AI services.

### Example Requests

Once the proxy is running, you can send requests to the AI services through it:

*   **OpenAI:**
    ```bash
    curl -X POST http://localhost:8080/openai/v1/chat/completions \
      -H "Content-Type: application/json" \
      # -H "Authorization: Bearer your-openai_proxy_key" # If openai_proxy_key is set
      -d '{
            "model": "gpt-3.5-turbo",
            "messages": [{"role": "user", "content": "Hello!"}]
          }'
    ```

*   **Gemini:**
    ```bash
    curl -X POST http://localhost:8080/gemini/v1beta/models/gemini-pro:generateContent \
      -H "Content-Type: application/json" \
      # -H "Authorization: Bearer your-gemini_proxy_key" # If gemini_proxy_key is set
      -d '{
            "contents": [{"parts":[{"text": "Write a story about a magic backpack."}]}]
          }'
    ```

*   **Claude:**
    ```bash
    curl -X POST http://localhost:8080/claude/v1/messages \
      -H "Content-Type: application/json" \
      -H "x-api-key: dummy" # Claude API requires an x-api-key header, it will be replaced by the proxy.
      -H "anthropic-version: 2023-06-01" \
      # -H "Authorization: Bearer your-claude_proxy_key" # If claude_proxy_key is set
      -d '{
            "model": "claude-3-opus-20240229",
            "max_tokens": 1024,
            "messages": [
              {"role": "user", "content": "Hello, Claude"}
            ]
          }'
    ```
    *Note for Claude:* The proxy handles the actual `x-api-key` and `anthropic-version` headers. You might need to send placeholder values if your client requires them, but the proxy will overwrite them with the correct ones from `secret.json` and the required version.

    *   (Optional) If `*_proxy_key` is configured in `secret.json`, it checks the incoming request's `Authorization` header for a matching proxy key.
    *   (New) If the `-outbound-proxy-url` command-line argument is provided, all outgoing requests from Poorman Proxy to the AI services are routed through this specified proxy.

## Development

### Run Tests

To run the project tests:

```bash
make test
```

### Debugging

The `Makefile` includes a `debug` target that runs test scripts located in the `test/` directory (e.g., `test/gemini_chat.sh`). You can uncomment and modify these scripts for debugging specific AI service integrations.

```bash
make debug
```

## How it Works

The proxy uses `net/http/httputil.ReverseProxy` to forward requests. For each supported AI service:
1.  A specific path prefix (`/openai/`, `/gemini/`, `/claude/`) routes the request to the corresponding upstream API.
2.  The path prefix is stripped before forwarding.
3.  A dedicated header rewrite function is invoked to:
    *   Set the correct `Host` header for the upstream service.
    *   Inject the appropriate API key (e.g., `Authorization: Bearer <OPENAI_API_KEY>`, `x-goog-api-key: <GEMINI_API_KEY>`, `x-api-key: <CLAUDE_API_KEY>`) from the `secret.json` file.
    *   (Optional) If `*_proxy_key` is configured in `secret.json`, it checks the incoming request's `Authorization` header for a matching proxy key.

## TODO

*   Improve error handling and logging.
*   Add more comprehensive tests.
*   Consider support for more AI services.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue. 