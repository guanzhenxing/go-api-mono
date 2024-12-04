package model

import "time"

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:32"`
	Email     string    `json:"email" gorm:"uniqueIndex;size:128"`
	Password  string    `json:"-" gorm:"size:128"` // 密码不返回给前端
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
