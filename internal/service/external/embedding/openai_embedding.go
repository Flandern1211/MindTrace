// internal/service/external/embedding/openai_embedding.go
package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OpenAIEmbedding OpenAI 兼容的 Embedding 服务
type OpenAIEmbedding struct {
	apiURL string
	apiKey string
	model  string
	client *http.Client
}

// NewOpenAIEmbedding 创建 OpenAI Embedding 服务
func NewOpenAIEmbedding(apiURL, apiKey, model string) *OpenAIEmbedding {
	if apiURL == "" {
		apiURL = "https://api.openai.com/v1"
	}
	return &OpenAIEmbedding{
		apiURL: apiURL,
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
	}
}

// embeddingRequest Embedding 请求
type embeddingRequest struct {
	Model string      `json:"model"`
	Input interface{} `json:"input"` // string 或 []string
}

// embeddingResponse Embedding 响应
type embeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
}

func (e *OpenAIEmbedding) Name() string {
	return "openai"
}

func (e *OpenAIEmbedding) Embed(text string) ([]float64, error) {
	results, err := e.doRequest([]string{text})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("Embedding 返回空结果")
	}
	return results[0], nil
}

func (e *OpenAIEmbedding) EmbedBatch(texts []string) ([][]float64, error) {
	return e.doRequest(texts)
}

func (e *OpenAIEmbedding) doRequest(input interface{}) ([][]float64, error) {
	reqBody := embeddingRequest{
		Model: e.model,
		Input: input,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	url := fmt.Sprintf("%s/embeddings", e.apiURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.apiKey))

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Embedding API 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Embedding API 返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var result embeddingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 按 index 排序结果
	embeddings := make([][]float64, len(result.Data))
	for _, d := range result.Data {
		if d.Index < len(embeddings) {
			embeddings[d.Index] = d.Embedding
		}
	}

	return embeddings, nil
}
