package entity

// User 用户实体
type User struct {
	BaseEntity
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"` // 用户名
	Password string `gorm:"type:varchar(255);not null" json:"-"`                   // 密码(不返回给前端)
	Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`            // 邮箱
	Nickname string `gorm:"type:varchar(50)" json:"nickname"`                      // 昵称
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`                       // 头像URL
	Status   int    `gorm:"type:tinyint;default:1" json:"status"`                  // 状态: 1正常 2禁用
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
