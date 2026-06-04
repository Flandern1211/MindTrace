package entity

// Conversation 会话实体
type Conversation struct {
	BaseEntity
	NotebookID int    `gorm:"index;not null;comment:所属笔记本ID"`
	Title      string `gorm:"type:varchar(100);not null;comment:会话标题"`
}

// TableName 指定表名
func (Conversation) TableName() string {
	return "conversations"
}
