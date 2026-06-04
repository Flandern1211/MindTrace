package repository

import (
	"YoudaoNoteLm/internal/model/entity"
	"errors"

	"gorm.io/gorm"
)

// ===== UserSearchConfig =====

type userSearchConfigRepository struct {
	db *gorm.DB
}

func NewUserSearchConfigRepository(db *gorm.DB) UserSearchConfigRepository {
	return &userSearchConfigRepository{db: db}
}

func (r *userSearchConfigRepository) FindByUser(userID uint) ([]*entity.UserSearchConfig, error) {
	var configs []*entity.UserSearchConfig
	err := r.db.Where("user_id = ?", userID).Order("priority ASC").Find(&configs).Error
	return configs, err
}

func (r *userSearchConfigRepository) FindByID(id uint) (*entity.UserSearchConfig, error) {
	var config entity.UserSearchConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

func (r *userSearchConfigRepository) Create(config *entity.UserSearchConfig) error {
	return r.db.Create(config).Error
}

func (r *userSearchConfigRepository) Update(config *entity.UserSearchConfig) error {
	return r.db.Save(config).Error
}

func (r *userSearchConfigRepository) Delete(id uint) error {
	return r.db.Delete(&entity.UserSearchConfig{}, id).Error
}

// ===== UserASRConfig =====

type userASRConfigRepository struct {
	db *gorm.DB
}

func NewUserASRConfigRepository(db *gorm.DB) UserASRConfigRepository {
	return &userASRConfigRepository{db: db}
}

func (r *userASRConfigRepository) FindByUser(userID uint) ([]*entity.UserASRConfig, error) {
	var configs []*entity.UserASRConfig
	err := r.db.Where("user_id = ?", userID).Order("priority ASC").Find(&configs).Error
	return configs, err
}

func (r *userASRConfigRepository) FindByID(id uint) (*entity.UserASRConfig, error) {
	var config entity.UserASRConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

func (r *userASRConfigRepository) Create(config *entity.UserASRConfig) error {
	return r.db.Create(config).Error
}

func (r *userASRConfigRepository) Update(config *entity.UserASRConfig) error {
	return r.db.Save(config).Error
}

func (r *userASRConfigRepository) Delete(id uint) error {
	return r.db.Delete(&entity.UserASRConfig{}, id).Error
}

// ===== UserEmbeddingConfig =====

type userEmbeddingConfigRepository struct {
	db *gorm.DB
}

func NewUserEmbeddingConfigRepository(db *gorm.DB) UserEmbeddingConfigRepository {
	return &userEmbeddingConfigRepository{db: db}
}

func (r *userEmbeddingConfigRepository) FindByUser(userID uint) ([]*entity.UserEmbeddingConfig, error) {
	var configs []*entity.UserEmbeddingConfig
	err := r.db.Where("user_id = ?", userID).Order("priority ASC").Find(&configs).Error
	return configs, err
}

func (r *userEmbeddingConfigRepository) FindByID(id uint) (*entity.UserEmbeddingConfig, error) {
	var config entity.UserEmbeddingConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

func (r *userEmbeddingConfigRepository) Create(config *entity.UserEmbeddingConfig) error {
	return r.db.Create(config).Error
}

func (r *userEmbeddingConfigRepository) Update(config *entity.UserEmbeddingConfig) error {
	return r.db.Save(config).Error
}

func (r *userEmbeddingConfigRepository) Delete(id uint) error {
	return r.db.Delete(&entity.UserEmbeddingConfig{}, id).Error
}
