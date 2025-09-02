package service

import (
	"application/model"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

// 创建账号服务实例
func NewAccountService() (*AccountService, error) {
	db := model.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return &AccountService{db: db}, nil
}

// 用户注册
func (s *AccountService) Register(req *model.RegisterRequest) error {
	// 检查用户名是否已存在
	var existingUser model.User
	err := s.db.Where("username = ?", req.Username).First(&existingUser).Error
	if err == nil {
		return fmt.Errorf("用户名已存在")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("检查用户名失败：%v", err)
	}

	// 加密密码
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败：%v", err)
	}

	// 创建用户，使用默认头像
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		AvatarURL:    model.DefaultAvatarURL,
		PasswordHash: string(passwordHash),
		Org:          req.Org,
	}

	fmt.Printf("%T", user.Org)
	// 保存用户到数据库
	err = s.db.Create(user).Error
	if err != nil {
		return fmt.Errorf("保存用户失败：%v", err)
	}

	return nil
}

// 用户登录
func (s *AccountService) Login(req *model.LoginRequest) (string, *model.User, error) {
	// 查找用户
	var user model.User
	err := s.db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, fmt.Errorf("用户不存在")
		}
		return "", nil, fmt.Errorf("查询用户失败：%v", err)
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", nil, fmt.Errorf("密码错误")
	}

	// 生成JWT令牌
	token, err := s.generateJWT(&user)
	if err != nil {
		return "", nil, fmt.Errorf("生成JWT令牌失败：%v", err)
	}

	// 保存令牌到数据库
	err = s.saveToken(user.ID, token)
	if err != nil {
		return "", nil, fmt.Errorf("保存令牌失败：%v", err)
	}

	return token, &user, nil
}

// 用户登出
func (s *AccountService) Logout(tokenString string) error {
	// 从数据库中删除令牌
	err := s.db.Where("token = ?", tokenString).Delete(&model.Token{}).Error
	if err != nil {
		return fmt.Errorf("删除令牌失败：%v", err)
	}
	return nil
}

// 根据用户ID获取用户信息
func (s *AccountService) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败：%v", err)
	}
	return &user, nil
}

// 更新用户信息
func (s *AccountService) UpdateUser(userID uint, updates map[string]interface{}) error {
	// 检查用户是否存在
	var user model.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("查询用户失败：%v", err)
	}

	// 更新字段
	if email, ok := updates["email"].(string); ok && email != "" {
		user.Email = email
	}

	if password, ok := updates["password"].(string); ok && password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("密码加密失败：%v", err)
		}
		user.PasswordHash = string(passwordHash)
	}

	// 保存更新
	err = s.db.Save(&user).Error
	if err != nil {
		return fmt.Errorf("更新用户失败：%v", err)
	}

	return nil
}

// 根据用户ID获取头像URL
func (s *AccountService) GetAvatarById(userID uint) (string, error) {
	var user model.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("用户不存在")
		}
		return "", fmt.Errorf("查询用户失败：%v", err)
	}
	return user.AvatarURL, nil
}

// 更新头像
func (s *AccountService) UpdateAvatar(userID uint, newAvatarURL string) (string, error) {
	var user model.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("用户不存在")
		}
		return "", fmt.Errorf("查询用户失败：%v", err)
	}
	user.AvatarURL = newAvatarURL
	err = s.db.Save(&user).Error
	if err != nil {
		return "", fmt.Errorf("更新头像失败：%v", err)
	}
	return user.AvatarURL, nil
}

// 辅助方法

// 生成JWT令牌
func (s *AccountService) generateJWT(user *model.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24小时过期

	claims := &model.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Org:      user.Org,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(model.JWT_SECRET))
}

// saveToken 保存令牌
func (s *AccountService) saveToken(userID uint, token string) error {
	tokenRecord := &model.Token{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err := s.db.Create(tokenRecord).Error
	if err != nil {
		return fmt.Errorf("保存令牌失败：%v", err)
	}

	return nil
}
