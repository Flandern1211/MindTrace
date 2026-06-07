package service

import (
	"YoudaoNoteLm/internal/model/entity"
	"YoudaoNoteLm/internal/repository"
	bizerrors "YoudaoNoteLm/pkg/errors"
)

type userConfigService struct {
	configRepo repository.UserConfigRepository
}

func NewUserConfigService(configRepo repository.UserConfigRepository) UserConfigService {
	return &userConfigService{configRepo: configRepo}
}

// ===== Search Config =====

func (s *userConfigService) ListSearchConfigs(userID uint) ([]*entity.UserConfig, error) {
	config, err := s.configRepo.FindByUserAndType(userID, "search")
	if err != nil {
		return nil, err
	}
	if config == nil {
		return []*entity.UserConfig{}, nil
	}
	return []*entity.UserConfig{config}, nil
}

func (s *userConfigService) CreateSearchConfig(userID uint, config *entity.UserConfig) error {
	config.UserID = userID
	config.ConfigType = "search"
	return s.configRepo.Create(config)
}

func (s *userConfigService) UpdateSearchConfig(id uint, config *entity.UserConfig) error {
	existing, err := s.configRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return bizerrors.ErrNotFound
	}
	config.ID = id
	config.UserID = existing.UserID
	config.ConfigType = "search"
	return s.configRepo.Update(config)
}

func (s *userConfigService) DeleteSearchConfig(id uint) error {
	return s.configRepo.Delete(id)
}

// ===== ASR Config =====

func (s *userConfigService) ListASRConfigs(userID uint) ([]*entity.UserConfig, error) {
	config, err := s.configRepo.FindByUserAndType(userID, "asr")
	if err != nil {
		return nil, err
	}
	if config == nil {
		return []*entity.UserConfig{}, nil
	}
	return []*entity.UserConfig{config}, nil
}

func (s *userConfigService) CreateASRConfig(userID uint, config *entity.UserConfig) error {
	config.UserID = userID
	config.ConfigType = "asr"
	return s.configRepo.Create(config)
}

func (s *userConfigService) UpdateASRConfig(id uint, config *entity.UserConfig) error {
	existing, err := s.configRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return bizerrors.ErrNotFound
	}
	config.ID = id
	config.UserID = existing.UserID
	config.ConfigType = "asr"
	return s.configRepo.Update(config)
}

func (s *userConfigService) DeleteASRConfig(id uint) error {
	return s.configRepo.Delete(id)
}

// ===== Embedding Config =====

func (s *userConfigService) ListEmbeddingConfigs(userID uint) ([]*entity.UserConfig, error) {
	config, err := s.configRepo.FindByUserAndType(userID, "embedding")
	if err != nil {
		return nil, err
	}
	if config == nil {
		return []*entity.UserConfig{}, nil
	}
	return []*entity.UserConfig{config}, nil
}

func (s *userConfigService) CreateEmbeddingConfig(userID uint, config *entity.UserConfig) error {
	config.UserID = userID
	config.ConfigType = "embedding"
	return s.configRepo.Create(config)
}

func (s *userConfigService) UpdateEmbeddingConfig(id uint, config *entity.UserConfig) error {
	existing, err := s.configRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return bizerrors.ErrNotFound
	}
	config.ID = id
	config.UserID = existing.UserID
	config.ConfigType = "embedding"
	return s.configRepo.Update(config)
}

func (s *userConfigService) DeleteEmbeddingConfig(id uint) error {
	return s.configRepo.Delete(id)
}