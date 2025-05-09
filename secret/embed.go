package secret

import (
	_ "embed"
	"encoding/json"
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
	return secret
}
