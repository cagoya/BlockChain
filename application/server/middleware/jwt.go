package middleware

import (
	"application/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// JWT中间件
type JWTMiddleware struct {
	db *gorm.DB
}

// 创建JWT中间件实例
func NewJWTMiddleware() (*JWTMiddleware, error) {
	db := model.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return &JWTMiddleware{db: db}, nil
}

// JWT认证中间件
func (m *JWTMiddleware) Auth(targetOrgs int32) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "缺少认证令牌",
			})
			c.Abort()
			return
		}

		// 提取Bearer令牌
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "无效的认证令牌格式",
			})
			c.Abort()
			return
		}

		// 验证令牌
		claims, err := m.validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "认证失败：" + err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("org", claims.Org)

		// 如果targetOrgs为-1，则不进行组织权限认证
		if targetOrgs == -1 {
			c.Next()
			return
		}

		// 组织权限认证
		org, exists := c.Get("org")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "获取组织信息失败",
			})
			c.Abort()
			return
		}
		for _, org := range org.(pq.Int32Array) {
			if org == targetOrgs {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "权限不足，不属于指定组织",
		})
		c.Abort()
	}
}

// 验证JWT令牌
func (m *JWTMiddleware) validateToken(tokenString string) (*model.Claims, error) {
	// 解析JWT令牌
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(model.JWT_SECRET), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析JWT令牌失败：%v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("无效的JWT令牌")
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		return nil, fmt.Errorf("无效的JWT声明")
	}

	// 检查令牌是否在数据库中
	isInDatabase, err := m.isTokenInDatabase(tokenString)
	if err != nil {
		return nil, fmt.Errorf("检查令牌数据库失败：%v", err)
	}
	if !isInDatabase {
		return nil, fmt.Errorf("令牌不在数据库中")
	}

	return claims, nil
}

// 检查令牌是否在数据库中
func (m *JWTMiddleware) isTokenInDatabase(token string) (bool, error) {
	var count int64
	err := m.db.Model(&model.Token{}).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
