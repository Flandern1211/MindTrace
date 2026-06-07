package service

import "YoudaoNoteLm/internal/model/entity"

// UserConfigService 用户配置服务接口
type UserConfigService interface {
	// 搜索配置
	ListSearchConfigs(userID uint) ([]*entity.UserConfig, error)
	CreateSearchConfig(userID uint, config *entity.UserConfig) error
	UpdateSearchConfig(id uint, config *entity.UserConfig) error
	DeleteSearchConfig(id uint) error

	// ASR 配置
	ListASRConfigs(userID uint) ([]*entity.UserConfig, error)
	CreateASRConfig(userID uint, config *entity.UserConfig) error
	UpdateASRConfig(id uint, config *entity.UserConfig) error
	DeleteASRConfig(id uint) error

	// Embedding 配置
	ListEmbeddingConfigs(userID uint) ([]*entity.UserConfig, error)
	CreateEmbeddingConfig(userID uint, config *entity.UserConfig) error
	UpdateEmbeddingConfig(id uint, config *entity.UserConfig) error
	DeleteEmbeddingConfig(id uint) error
}