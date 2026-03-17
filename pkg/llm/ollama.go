package llm

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type Response struct {
	Response string `json:"response"`
}

func CallOllama(model, prompt string) string {
	reqBody := Request{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return "FINAL: LLM unavailable — " + err.Error()
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result Response
	json.Unmarshal(body, &result)

	return strings.TrimSpace(result.Response)

}
