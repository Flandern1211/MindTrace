package external

import (
	"encoding/json"
	"fmt"

	"YoudaoNoteLm/pkg/config"
	"YoudaoNoteLm/pkg/logger"

	"go.uber.org/zap"
)

// NewASRService 根据配置创建 ASR 服务
// 配置示例：
//
//	asr:
//	  provider: aliyun_nls
//	  params:
//	    access_key_id: "xxx"
//	    access_key_secret: "xxx"
//	    app_key: "xxx"
func NewASRService(cfg config.ASRConfig) ASRService {
	switch cfg.Provider {
	case "aliyun_nls":
		return NewAliyunNLSASRService(
			cfg.GetString("access_key_id"),
			cfg.GetString("access_key_secret"),
			cfg.GetString("app_key"),
		)
	default:
		panic(fmt.Sprintf("不支持的 ASR provider: %s", cfg.Provider))
	}
}

// NewASRServiceFromDB 根据数据库用户配置创建 ASR 服务
// UserConfig 的 ExtraConfig(JSON) 中存放各服务商特有参数
func NewASRServiceFromDB(provider, apiURL, apiKey, extraConfig string) ASRService {
	switch provider {
	case "aliyun_nls":
		var params map[string]interface{}
		if extraConfig != "" {
			if err := json.Unmarshal([]byte(extraConfig), &params); err != nil {
				logger.Error("解析ASR ExtraConfig失败", zap.Error(err))
			}
		}
		getStr := func(key string) string {
			if params != nil {
				if v, ok := params[key].(string); ok {
					return v
				}
			}
			return ""
		}
		accessKeyID := apiKey
		accessKeySecret := getStr("access_key_secret")
		appKey := getStr("app_key")
		return NewAliyunNLSASRService(accessKeyID, accessKeySecret, appKey)
	default:
		logger.Warn("不支持的ASR provider，返回nil", zap.String("provider", provider))
		return nil
	}
}
