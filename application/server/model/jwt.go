package model

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/lib/pq"
)

// Claims JWT声明
type Claims struct {
	UserID   int           `json:"user_id"`
	Username string        `json:"username"`
	Org      pq.Int32Array `json:"org"`
	jwt.RegisteredClaims
}

// JWT密钥（生产环境应该从配置文件读取）
const JWT_SECRET = "9f2b8c5d1e4a7f0b3d6a2c1e8f5b4a0d9e2c8f6b1a3d5e7c8f9b0a1d2e3c4f5a"
