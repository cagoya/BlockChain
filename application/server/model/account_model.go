package model

import (
	"github.com/lib/pq"
	"time"
)

// User 用户信息
type User struct {
	ID           uint          `json:"id" gorm:"primaryKey;autoIncrement"`                    // 主键，用户 ID
	Username     string        `json:"username" gorm:"uniqueIndex;type:varchar(50);not null"` // 用户名
	Email        string        `json:"email" gorm:"type:varchar(50);not null"`                // 邮箱
	AvatarName   string        `json:"-" gorm:"type:varchar(255);not null"`                   // 头像图片名称
	PasswordHash string        `json:"-" gorm:"type:varchar(255);not null"`                   // 密码哈希
	Org          pq.Int32Array `json:"org" gorm:"type:integer[];not null"`                    // 组织，可以有多个
	CreateTime   time.Time     `json:"createTime" gorm:"autoCreateTime"`                      // 创建时间
	UpdateTime   time.Time     `json:"updateTime" gorm:"autoUpdateTime"`                      // 更新时间
}

// Token 令牌信息
type Token struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`          // 主键，令牌 ID
	Token      string    `json:"token" gorm:"uniqueIndex;type:text;not null"` // 令牌
	UserID     uint      `json:"userID" gorm:"not null"`                      // 用户 ID
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

// 头像文件夹路径
const AvatarPath = "file/avatar"
