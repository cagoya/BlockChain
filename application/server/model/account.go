package model

import (
	"github.com/lib/pq"
	"time"
)

// User 用户信息
type User struct {
	ID           int           `json:"id" gorm:"primaryKey;autoIncrement"`                    // 主键，用户 ID
	Username     string        `json:"username" gorm:"uniqueIndex;type:varchar(50);not null"` // 用户名
	Email        string        `json:"email" gorm:"type:varchar(50);not null"`                // 邮箱
	AvatarURL    string        `json:"avatarURL" gorm:"type:varchar(255);not null"`           // 头像URL
	PasswordHash string        `json:"-" gorm:"type:varchar(255);not null"`                   // 密码哈希
	Org          pq.Int32Array `json:"org" gorm:"type:integer[];not null"`                    // 组织，可以有多个
	CreateTime   time.Time     `json:"createTime" gorm:"autoCreateTime"`                      // 创建时间
	UpdateTime   time.Time     `json:"updateTime" gorm:"autoUpdateTime"`                      // 更新时间
}

// Token 令牌信息
type Token struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`          // 主键，令牌 ID
	Token      string    `json:"token" gorm:"uniqueIndex;type:text;not null"` // 令牌
	UserID     int       `json:"userID" gorm:"not null"`                      // 用户 ID
	ExpiresAt  time.Time `json:"expiresAt" gorm:"not null"`                   // 过期时间
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`            // 创建时间
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (Token) TableName() string {
	return "tokens"
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string        `json:"username"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
	Org      pq.Int32Array `json:"org"`
}

// 更新组织请求结构
type UpdateOrgRequest struct {
	UserID int           `json:"userID"`
	Org    pq.Int32Array `json:"org"`
}

// 检查是否属于制定组织
func CheckOrg(org pq.Int32Array, targetOrg int32) bool {
	for _, i := range org {
		if i == targetOrg {
			return true
		}
	}
	return false
}
