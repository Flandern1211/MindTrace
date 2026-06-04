package entity

// Notebook 笔记本实体
type Notebook struct {
	BaseEntity
	UserID int    `gorm:"index;not null;comment:所属用户ID"`
	Name   string `gorm:"type:varchar(100);not null;comment:笔记本名称"`
}

// TableName 指定表名
func (Notebook) TableName() string {
	return "notebooks"
}
