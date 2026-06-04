package entity

// SysConfig 系统配置实体
type SysConfig struct {
	BaseEntity
	ConfigGroup string `gorm:"type:varchar(50);not null;uniqueIndex:uk_group_key" json:"config_group"`
	ConfigKey   string `gorm:"type:varchar(100);not null;uniqueIndex:uk_group_key" json:"config_key"`
	ConfigValue string `gorm:"type:json;not null" json:"config_value"`
	Enabled     bool   `gorm:"default:true" json:"enabled"`
	Description string `gorm:"type:varchar(255)" json:"description"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}
