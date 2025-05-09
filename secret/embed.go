package secret

import (
	_ "embed"
	"encoding/json"
	"os"
	"strings"
)

//go:embed secret.json
var secretJSON []byte

type Secret struct {
	OpenAIKey string   `json:"openai_api_key"`
	GeminiKey string   `json:"gemini_api_key"`
	ClaudeKey string   `json:"claude_api_key"`
	ProxyKey  []string `json:"proxy_key"`
	// GeminiProxyKey []string `json:"gemini_proxy_key"`
	// ClaudeProxyKey []string `json:"claude_proxy_key"`
}

func Load() Secret {
	var secret Secret
	err := json.Unmarshal(secretJSON, &secret)
	if err != nil {
		panic(err)
	}

	if val, ok := os.LookupEnv("OPENAI_API_KEY"); ok {
		secret.OpenAIKey = val
	}
	if val, ok := os.LookupEnv("GEMINI_API_KEY"); ok {
		secret.GeminiKey = val
	}
	if val, ok := os.LookupEnv("CLAUDE_API_KEY"); ok {
		secret.ClaudeKey = val
	}
	if val, ok := os.LookupEnv("PROXY_KEY"); ok {
		if val != "" {
			secret.ProxyKey = strings.Split(val, ",")
		} else {
			secret.ProxyKey = []string{}
		}
	}

	return secret
}
