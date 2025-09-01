package api

import (
	"application/model"
	"application/service"
	"application/utils"
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"path/filepath"
	"strings"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler() *AccountHandler {
	accountService, err := service.NewAccountService()
	if err != nil {
		panic("初始化账号服务失败：" + err.Error())
	}

	return &AccountHandler{
		accountService: accountService,
	}
}

func (h *AccountHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数格式错误")
		return
	}

	err := h.accountService.Register(&req)
	if err != nil {
		utils.ServerError(c, "注册失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "注册成功", nil)
}

func (h *AccountHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数格式错误")
		return
	}

	token, user, err := h.accountService.Login(&req)
	if err != nil {
		utils.ServerError(c, "登录失败："+err.Error())
		return
	}

	// 返回用户信息（不包含密码）和令牌
	response := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"org":        user.Org,
			"createTime": user.CreateTime,
			"updateTime": user.UpdateTime,
		},
	}

	utils.SuccessWithMessage(c, "登录成功", response)
}

func (h *AccountHandler) Logout(c *gin.Context) {
	// 从请求头获取令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.BadRequest(c, "缺少认证令牌")
		return
	}

	// 提取Bearer令牌
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		utils.BadRequest(c, "无效的认证令牌格式")
		return
	}

	err := h.accountService.Logout(token)
	if err != nil {
		utils.ServerError(c, "登出失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "登出成功", nil)
}

func (h *AccountHandler) GetProfile(c *gin.Context) {
	// 从上下文中获取用户ID（由中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}

	user, err := h.accountService.GetUserByID(userID.(uint))
	if err != nil {
		utils.ServerError(c, "获取用户信息失败："+err.Error())
		return
	}

	// 返回用户信息（不包含密码）
	response := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"org":        user.Org,
		"createTime": user.CreateTime,
		"updateTime": user.UpdateTime,
	}

	utils.Success(c, response)
}

func (h *AccountHandler) UpdateProfile(c *gin.Context) {
	// 从上下文中获取用户ID（由中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数格式错误")
		return
	}

	// 只允许更新邮箱和密码
	updates := make(map[string]interface{})
	if email, ok := req["email"].(string); ok && email != "" {
		updates["email"] = email
	}
	if password, ok := req["password"].(string); ok && password != "" {
		updates["password"] = password
	}

	if len(updates) == 0 {
		utils.BadRequest(c, "没有有效的更新字段")
		return
	}

	err := h.accountService.UpdateUser(userID.(uint), updates)
	if err != nil {
		utils.ServerError(c, "更新用户信息失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (h *AccountHandler) GetAvatar(c *gin.Context) {
	// 从上下文中获取用户ID（由中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	avatar, err := h.accountService.GetAvatarById(userID.(uint))
	if err != nil {
		utils.ServerError(c, "获取头像失败："+err.Error())
		return
	}
	imageType, err := getImageMIMEType(avatar.Bytes())
	if err != nil {
		utils.ServerError(c, "获取图片类型失败："+err.Error())
		return
	}
	c.Header("Content-Type", imageType)
	c.Data(http.StatusOK, imageType, avatar.Bytes())
}

func (h *AccountHandler) UpdateAvatar(c *gin.Context) {
	// 从上下文中获取用户ID（由中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	// 从表单中获取图片
	file, err := c.FormFile("avatar")
	if err != nil {
		utils.ServerError(c, "获取头像失败："+err.Error())
		return
	}
	// 保存图片到本地
	uniqueID := uuid.New().String()
	filePath := filepath.Join(model.AvatarPath, uniqueID+"_"+file.Filename)
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		utils.ServerError(c, "保存头像失败："+err.Error())
		return
	}
	// 更新头像
	err = h.accountService.UpdateAvatar(userID.(uint), uniqueID+"_"+file.Filename)
	if err != nil {
		utils.ServerError(c, "更新头像失败："+err.Error())
		return
	}
	utils.SuccessWithMessage(c, "更新成功", nil)
}

// 辅助函数，判断图片类型
func getImageMIMEType(data []byte) (string, error) {
	// 确保数据足够长以检查文件头
	if len(data) < 8 {
		return "", fmt.Errorf("数据太短，无法识别图片类型")
	}

	// 检查 PNG 文件头 (89 50 4E 47 0D 0A 1A 0A)
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return "image/png", nil
	}

	// 检查 JPEG 文件头 (FF D8 FF)
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return "image/jpeg", nil
	}

	contentType := http.DetectContentType(data)
	if contentType == "application/octet-stream" {
		return "", fmt.Errorf("不支持的图片类型")
	}

	return contentType, nil
}
