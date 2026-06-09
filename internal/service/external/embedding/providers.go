// internal/service/external/embedding/providers.go
package embedding

import (
	"fmt"

	"YoudaoNoteLm/internal/service/external"
)

const ServiceType = "embedding"

// openaiCompatibleEmbeddingFactory 创建 OpenAI 兼容的 Embedding 客户端
func openaiCompatibleEmbeddingFactory(displayName string) external.FactoryFunc {
	return func(cfg *external.ServiceConfig) (interface{}, error) {
		model := cfg.Model
		if model == "" {
			model = cfg.GetExtraString("model")
		}
		if model == "" {
			return nil, fmt.Errorf("Embedding 模型名称未配置")
		}
		return NewOpenAIEmbedding(cfg.APIURL, cfg.APIKey, model), nil
	}
}

func init() {
	r := external.GetGlobalRegistry()

	// OpenAI Embedding
	r.Register(ServiceType, "openai", "OpenAI Embedding",
		[]string{"api_key", "model"}, []string{"api_url"},
		openaiCompatibleEmbeddingFactory("OpenAI"), map[string]string{
			"api_key": "API Key",
			"model":   "模型名称（如 text-embedding-3-small）",
			"api_url": "API 地址（可选，用于代理）",
		})

	// 智谱 Embedding
	r.Register(ServiceType, "zhipu", "智谱 Embedding",
		[]string{"api_key", "model"}, []string{"api_url"},
		openaiCompatibleEmbeddingFactory("智谱"), map[string]string{
			"api_key": "API Key",
			"model":   "模型名称（如 embedding-3）",
			"api_url": "API 地址（默认 https://open.bigmodel.cn/api/paas/v4）",
		})

	// 火山引擎（豆包）Embedding
	r.Register(ServiceType, "volcengine", "火山引擎（豆包）Embedding",
		[]string{"api_key", "model"}, []string{"api_url"},
		openaiCompatibleEmbeddingFactory("火山引擎"), map[string]string{
			"api_key": "API Key",
			"model":   "模型名称或接入点 ID（如 doubao-embedding）",
			"api_url": "API 地址（默认 https://ark.cn-beijing.volces.com/api/v3）",
		})

	// 通义千问 Embedding
	r.Register(ServiceType, "qwen", "通义千问 Embedding",
		[]string{"api_key", "model"}, []string{"api_url"},
		openaiCompatibleEmbeddingFactory("通义千问"), map[string]string{
			"api_key": "API Key",
			"model":   "模型名称（如 text-embedding-v2）",
			"api_url": "API 地址（默认 https://dashscope.aliyuncs.com/compatible-mode/v1）",
		})
}
