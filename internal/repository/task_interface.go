package repository

import "YoudaoNoteLm/internal/model/entity"

// ImportTaskRepository 导入任务仓储接口
type ImportTaskRepository interface {
	FindByID(id uint) (*entity.ImportTask, error)
	Create(task *entity.ImportTask) error
	Update(task *entity.ImportTask) error
	IncrementSuccess(id uint) error
	IncrementFail(id uint, errMsg string) error
}
