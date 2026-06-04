package repository

import (
	"YoudaoNoteLm/internal/model/entity"
	"errors"
	"time"

	"gorm.io/gorm"
)

type audioPreviewRepository struct {
	db *gorm.DB
}

func NewAudioPreviewRepository(db *gorm.DB) AudioPreviewRepository {
	return &audioPreviewRepository{db: db}
}

func (r *audioPreviewRepository) Create(preview *entity.AudioPreview) error {
	return r.db.Create(preview).Error
}

func (r *audioPreviewRepository) FindByPreviewID(previewID string) (*entity.AudioPreview, error) {
	var preview entity.AudioPreview
	err := r.db.Where("preview_id = ?", previewID).First(&preview).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &preview, nil
}

func (r *audioPreviewRepository) UpdateStatus(previewID string, status string) error {
	return r.db.Model(&entity.AudioPreview{}).Where("preview_id = ?", previewID).
		Update("status", status).Error
}

func (r *audioPreviewRepository) Delete(previewID string) error {
	return r.db.Where("preview_id = ?", previewID).Delete(&entity.AudioPreview{}).Error
}

func (r *audioPreviewRepository) CleanExpired() error {
	return r.db.Where("expires_at < ? AND status = ?", time.Now(), "pending").
		Delete(&entity.AudioPreview{}).Error
}
