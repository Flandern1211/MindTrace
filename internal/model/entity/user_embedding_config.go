package entity

// UserEmbeddingConfig 用户Embedding配置
type UserEmbeddingConfig struct {
	BaseEntity
	UserID     uint   `gorm:"not null;index" json:"user_id"`
	Name       string `gorm:"type:varchar(100);not null" json:"name"`
	Provider   string `gorm:"type:varchar(50);not null" json:"provider"`
	APIKey     string `gorm:"type:varchar(512)" json:"api_key"`
	APIURL     string `gorm:"type:varchar(512);not null" json:"api_url"`
	ModelName  string `gorm:"type:varchar(100)" json:"model_name"`
	Dimensions *int   `json:"dimensions"`
	Enabled    bool   `gorm:"default:true" json:"enabled"`
	Priority   int    `gorm:"default:0" json:"priority"`
}

func (UserEmbeddingConfig) TableName() string {
	return "user_embedding_config"
}
