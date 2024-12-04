package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:32" binding:"required,min=3,max=32,alphanum,username_valid"`
	Email     string    `json:"email" gorm:"uniqueIndex;size:128" binding:"required,email,max=128,email_domain"`
	Password  string    `json:"password,omitempty" gorm:"size:128" binding:"required,min=6,max=128,password"` // 密码不返回给前端
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// Validate performs custom validation on the User model
func (u *User) Validate() error {
	// 检查用户名是否包含敏感词
	sensitiveWords := []string{"admin", "root", "system", "superuser"}
	usernameLower := strings.ToLower(u.Username)
	for _, word := range sensitiveWords {
		if strings.Contains(usernameLower, word) {
			return fmt.Errorf("username contains prohibited word: %s", word)
		}
	}

	// 检查邮箱域名是否在黑名单中
	emailBlacklist := []string{"tempmail.com", "throwaway.com", "fakemail.com"}
	emailParts := strings.Split(u.Email, "@")
	if len(emailParts) == 2 {
		domain := strings.ToLower(emailParts[1])
		for _, blacklisted := range emailBlacklist {
			if domain == blacklisted {
				return fmt.Errorf("email domain %s is not allowed", domain)
			}
		}
	}

	// 检查用户名格式（不允许连续的特殊字符或数字）
	if matched, _ := regexp.MatchString(`\d{4,}`, u.Username); matched {
		return fmt.Errorf("username cannot contain 4 or more consecutive numbers")
	}

	// 检查密码是否包含用户名
	if strings.Contains(strings.ToLower(u.Password), strings.ToLower(u.Username)) {
		return fmt.Errorf("password cannot contain username")
	}

	return nil
}

// ValidatePasswordStrength checks password strength beyond basic requirements
func (u *User) ValidatePasswordStrength() error {
	// 检查密码长度分数
	lengthScore := len(u.Password) / 8

	// 检查字符种类分数
	var typesScore int
	if regexp.MustCompile(`[A-Z]`).MatchString(u.Password) {
		typesScore++
	}
	if regexp.MustCompile(`[a-z]`).MatchString(u.Password) {
		typesScore++
	}
	if regexp.MustCompile(`[0-9]`).MatchString(u.Password) {
		typesScore++
	}
	if regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(u.Password) {
		typesScore++
	}

	// 检查是否有连续的字符
	for i := 0; i < len(u.Password)-2; i++ {
		if u.Password[i] == u.Password[i+1] && u.Password[i] == u.Password[i+2] {
			return fmt.Errorf("password cannot contain three or more consecutive identical characters")
		}
	}

	// 计算总分
	totalScore := lengthScore + typesScore

	// 密码强度要求
	if totalScore < 5 {
		return fmt.Errorf("password is too weak (score: %d/8), please use a stronger password", totalScore)
	}

	return nil
}
