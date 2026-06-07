package service

import (
	"YoudaoNoteLm/internal/model/dto/response"
	"YoudaoNoteLm/internal/model/entity"
	"YoudaoNoteLm/internal/repository"
	bizerrors "YoudaoNoteLm/pkg/errors"
	"encoding/json"
)

type adminService struct {
	userRepo   repository.UserRepository
	configRepo repository.SysConfigRepository
}

func NewAdminService(userRepo repository.UserRepository, configRepo repository.SysConfigRepository) AdminService {
	return &adminService{userRepo: userRepo, configRepo: configRepo}
}

func (s *adminService) ListUsers(page, size int, keyword string) ([]*response.AdminUserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	users, total, err := s.userRepo.List(offset, size)
	if err != nil {
		return nil, 0, err
	}

	list := make([]*response.AdminUserResponse, 0, len(users))
	for _, u := range users {
		list = append(list, &response.AdminUserResponse{
			ID: u.ID, Username: u.Username, Email: u.Email,
			Nickname: u.Nickname, Status: u.Status, CreatedAt: u.CreatedAt,
		})
	}

	return list, total, nil
}

func (s *adminService) UpdateUserStatus(userID uint, enabled bool) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return bizerrors.ErrUserNotFound
	}

	if enabled {
		user.Status = 1
	} else {
		user.Status = 2
	}
	return s.userRepo.Update(user)
}

func (s *adminService) GetConfigs(group string) ([]*entity.SysConfig, error) {
	return s.configRepo.FindByGroup(group)
}

func (s *adminService) UpdateConfig(group, key string, value json.RawMessage, enabled bool) error {
	config, err := s.configRepo.FindByGroupAndKey(group, key)
	if err != nil {
		return err
	}
	if config == nil {
		return bizerrors.ErrNotFound
	}

	config.ConfigValue = string(value)
	config.Enabled = enabled
	return s.configRepo.Update(config)
}

func (s *adminService) AddConfig(group, key string, value json.RawMessage, description string) error {
	existing, _ := s.configRepo.FindByGroupAndKey(group, key)
	if existing != nil {
		return bizerrors.New(bizerrors.CodeResourceAlreadyExists, "配置已存在")
	}

	config := &entity.SysConfig{
		ConfigGroup: group, ConfigKey: key,
		ConfigValue: string(value), Enabled: true, Description: description,
	}
	return s.configRepo.Create(config)
}

func (s *adminService) GetConfigStatus() ([]response.ConfigStatusGroupResponse, error) {
	summaries, err := s.configRepo.GetConfigStatusSummary()
	if err != nil {
		return nil, err
	}

	result := make([]response.ConfigStatusGroupResponse, 0, len(summaries))
	for _, sm := range summaries {
		result = append(result, response.ConfigStatusGroupResponse{
			Group: sm["group"].(string), Total: sm["total"].(int64),
			Enabled: sm["enabled"].(int64), Description: sm["description"].(string),
		})
	}
	return result, nil
}