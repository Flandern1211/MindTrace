package entity

// UserASRConfig 用户ASR配置
type UserASRConfig struct {
	BaseEntity
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Provider    string `gorm:"type:varchar(50);not null" json:"provider"`
	APIKey      string `gorm:"type:varchar(512)" json:"api_key"`
	APIURL      string `gorm:"type:varchar(512);not null" json:"api_url"`
	ExtraConfig string `gorm:"type:json" json:"extra_config"`
	Enabled     bool   `gorm:"default:true" json:"enabled"`
	Priority    int    `gorm:"default:0" json:"priority"`
}

func (UserASRConfig) TableName() string {
	return "user_asr_config"
}
