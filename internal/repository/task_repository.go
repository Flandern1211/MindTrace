package repository

import (
	"YoudaoNoteLm/internal/model/entity"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type importTaskRepository struct {
	db *gorm.DB
}

func NewImportTaskRepository(db *gorm.DB) ImportTaskRepository {
	return &importTaskRepository{db: db}
}

func (r *importTaskRepository) FindByID(id uint) (*entity.ImportTask, error) {
	var task entity.ImportTask
	err := r.db.First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *importTaskRepository) Create(task *entity.ImportTask) error {
	return r.db.Create(task).Error
}

func (r *importTaskRepository) Update(task *entity.ImportTask) error {
	return r.db.Save(task).Error
}

func (r *importTaskRepository) IncrementSuccess(id uint) error {
	return r.db.Model(&entity.ImportTask{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"success_count": gorm.Expr("success_count + 1"),
		}).Error
}

func (r *importTaskRepository) IncrementFail(id uint, errMsg string) error {
	task, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}

	var errorDetail map[string]string
	if task.ErrorDetail != "" {
		_ = json.Unmarshal([]byte(task.ErrorDetail), &errorDetail)
	} else {
		errorDetail = make(map[string]string)
	}
	errorDetail[errMsg] = errMsg
	detailJSON, _ := json.Marshal(errorDetail)

	return r.db.Model(&entity.ImportTask{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"fail_count":   gorm.Expr("fail_count + 1"),
			"error_detail": string(detailJSON),
		}).Error
}
