package entity

import "time"

// UserSearchConfig 用户搜索配置
type UserSearchConfig struct {
	BaseEntity
	UserID       uint       `gorm:"not null;index" json:"user_id"`
	Name         string     `gorm:"type:varchar(100);not null" json:"name"`
	Provider     string     `gorm:"type:varchar(50);not null" json:"provider"`
	APIKey       string     `gorm:"type:varchar(512);not null" json:"api_key"`
	APIURL       string     `gorm:"type:varchar(512)" json:"api_url"`
	DailyQuota   *int       `json:"daily_quota"`
	QuotaUsed    int        `gorm:"default:0" json:"quota_used"`
	QuotaResetAt *time.Time `json:"quota_reset_at"`
	Enabled      bool       `gorm:"default:true" json:"enabled"`
	Priority     int        `gorm:"default:0" json:"priority"`
}

func (UserSearchConfig) TableName() string {
	return "user_search_config"
}
