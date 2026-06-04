package repository

import "YoudaoNoteLm/internal/model/entity"

type UserSearchConfigRepository interface {
	FindByUser(userID uint) ([]*entity.UserSearchConfig, error)
	FindByID(id uint) (*entity.UserSearchConfig, error)
	Create(config *entity.UserSearchConfig) error
	Update(config *entity.UserSearchConfig) error
	Delete(id uint) error
}

type UserASRConfigRepository interface {
	FindByUser(userID uint) ([]*entity.UserASRConfig, error)
	FindByID(id uint) (*entity.UserASRConfig, error)
	Create(config *entity.UserASRConfig) error
	Update(config *entity.UserASRConfig) error
	Delete(id uint) error
}

type UserEmbeddingConfigRepository interface {
	FindByUser(userID uint) ([]*entity.UserEmbeddingConfig, error)
	FindByID(id uint) (*entity.UserEmbeddingConfig, error)
	Create(config *entity.UserEmbeddingConfig) error
	Update(config *entity.UserEmbeddingConfig) error
	Delete(id uint) error
}
