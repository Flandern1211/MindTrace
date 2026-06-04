package repository

import "YoudaoNoteLm/internal/model/entity"

// AudioPreviewRepository 音频预览仓储接口
type AudioPreviewRepository interface {
	Create(preview *entity.AudioPreview) error
	FindByPreviewID(previewID string) (*entity.AudioPreview, error)
	UpdateStatus(previewID string, status string) error
	Delete(previewID string) error
	CleanExpired() error
}
