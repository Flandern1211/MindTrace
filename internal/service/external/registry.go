// internal/service/external/registry.go
package external

import (
	"fmt"
	"sync"

	"YoudaoNoteLm/pkg/logger"

	"go.uber.org/zap"
)

// FactoryFunc provider 工厂函数签名
// 根据 ServiceConfig 创建具体的服务实例
type FactoryFunc func(cfg *ServiceConfig) (interface{}, error)

// ProviderInfo provider 元信息（用于 API 发现）
type ProviderInfo struct {
	ServiceType  string            `json:"service_type"`
	Provider     string            `json:"provider"`
	DisplayName  string            `json:"display_name"`
	RequiredKeys []string          `json:"required_keys"`
	OptionalKeys []string          `json:"optional_keys"`
	Implemented  bool              `json:"implemented"`             // 是否已实现
	KeyLabels    map[string]string `json:"key_labels,omitempty"`   // 参数中文标签
}

// providerEntry 注册表中的一个 provider
type providerEntry struct {
	info    ProviderInfo
	factory FactoryFunc
}

// Registry provider 注册表
// 线程安全，支持运行时注册和查询
type Registry struct {
	mu      sync.RWMutex
	entries map[string]map[string]providerEntry // [serviceType][providerName] → entry
}

// NewRegistry 创建空的 Registry
func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[string]map[string]providerEntry),
	}
}

// 全局 Registry 实例
var globalRegistry = NewRegistry()

// GetGlobalRegistry 获取全局 Registry
func GetGlobalRegistry() *Registry {
	return globalRegistry
}

// Register 注册一个 provider
func (r *Registry) Register(serviceType, providerName, displayName string,
	requiredKeys, optionalKeys []string, factory FactoryFunc, keyLabels map[string]string) {

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.entries[serviceType] == nil {
		r.entries[serviceType] = make(map[string]providerEntry)
	}

	// 标记是否已实现（factory 不为 nil）
	implemented := factory != nil

	r.entries[serviceType][providerName] = providerEntry{
		info: ProviderInfo{
			ServiceType:  serviceType,
			Provider:     providerName,
			DisplayName:  displayName,
			RequiredKeys: requiredKeys,
			OptionalKeys: optionalKeys,
			Implemented:  implemented,
			KeyLabels:    keyLabels,
		},
		factory: factory,
	}

	// logger 可能在 init() 阶段未初始化，使用安全调用
	if logger.GetLogger() != nil {
		logger.Info("Provider 已注册",
			zap.String("service_type", serviceType),
			zap.String("provider", providerName),
			zap.Bool("implemented", implemented),
		)
	}
}

// Create 根据 serviceType、providerName 和配置创建服务实例
func (r *Registry) Create(serviceType, providerName string, cfg *ServiceConfig) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	providers, ok := r.entries[serviceType]
	if !ok {
		return nil, fmt.Errorf("不支持的服务类型: %s", serviceType)
	}

	entry, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("不支持的 %s provider: %s", serviceType, providerName)
	}

	instance, err := entry.factory(cfg)
	if err != nil {
		return nil, fmt.Errorf("创建 %s/%s 失败: %w", serviceType, providerName, err)
	}

	return instance, nil
}

// ListProviders 列出所有已注册的 provider（或按 serviceType 过滤）
func (r *Registry) ListProviders(serviceType string) []ProviderInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []ProviderInfo

	if serviceType != "" {
		if providers, ok := r.entries[serviceType]; ok {
			for _, entry := range providers {
				result = append(result, entry.info)
			}
		}
		return result
	}

	for _, providers := range r.entries {
		for _, entry := range providers {
			result = append(result, entry.info)
		}
	}
	return result
}

// HasProvider 检查指定 provider 是否已注册
func (r *Registry) HasProvider(serviceType, providerName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if providers, ok := r.entries[serviceType]; ok {
		_, exists := providers[providerName]
		return exists
	}
	return false
}

// GetProviderInfo 获取指定 provider 的信息
func (r *Registry) GetProviderInfo(serviceType, providerName string) *ProviderInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if providers, ok := r.entries[serviceType]; ok {
		if entry, exists := providers[providerName]; exists {
			info := entry.info
			return &info
		}
	}
	return nil
}
