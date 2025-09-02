package api

import (
	"application/model"
	"application/pkg/image"
	"application/service"
	"application/utils"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService *service.AccountService
	imageHelper    *image.ImageHelper
}

func NewAccountHandler() *AccountHandler {
	accountService, err := service.NewAccountService()
	if err != nil {
		panic("初始化账号服务失败：" + err.Error())
	}

	return &AccountHandler{
		accountService: accountService,
		imageHelper:    image.NewImageHelper(),
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
		"user":  user,
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
	response := user

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
	avatarURL, err := h.accountService.GetAvatarById(userID.(uint))
	if err != nil {
		utils.ServerError(c, "获取头像失败："+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"avatarURL": avatarURL})
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
	// 保存图片到图床
	avatarURL, err := h.imageHelper.UploadImage(file)
	if err != nil {
		utils.ServerError(c, "保存头像失败："+err.Error())
		return
	}

	// 更新头像
	avatarURL, err = h.accountService.UpdateAvatar(userID.(uint), avatarURL)
	if err != nil {
		utils.ServerError(c, "更新头像失败："+err.Error())
		return
	}
	utils.SuccessWithMessage(c, "更新成功", avatarURL)
}
