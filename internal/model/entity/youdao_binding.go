package entity

// YoudaoBinding 有道云笔记绑定实体
type YoudaoBinding struct {
	BaseEntity
	UserID uint   `gorm:"not null;uniqueIndex:uk_user" json:"user_id"`
	APIKey string `gorm:"type:varchar(512);not null" json:"api_key"`
	Status string `gorm:"type:varchar(20);default:active" json:"status"`
}

func (YoudaoBinding) TableName() string {
	return "youdao_binding"
}
