package model

import (
	"time"
)

// User 用户信息
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string    `json:"username" gorm:"uniqueIndex;type:varchar(50);not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;type:varchar(100);not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	Org          string    `json:"org" gorm:"type:varchar(20);not null"`
	Role         string    `json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	CreateTime   time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime   time.Time `json:"updateTime" gorm:"autoUpdateTime"`
}

// Token 令牌信息
type Token struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Token      string    `json:"token" gorm:"uniqueIndex;type:text;not null"`
	UserID     uint      `json:"userID" gorm:"not null"`
	ExpiresAt  time.Time `json:"expiresAt" gorm:"not null"`
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`
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
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Org      string `json:"org"`
}
