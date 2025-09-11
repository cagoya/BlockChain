package api

import (
	"application/model"
	"application/service"
	"application/utils"
	"fmt"
	"path/filepath"
	"strconv"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	user, err := h.accountService.GetUserByID(userID.(int))
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

	err := h.accountService.UpdateUser(userID.(int), updates)
	if err != nil {
		utils.ServerError(c, "更新用户信息失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (h *AccountHandler) GetAvatar(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	avatarURL, err := h.accountService.GetAvatarById(userID)
	if err != nil {
		utils.ServerError(c, "获取头像失败："+err.Error())
		return
	}
	utils.Success(c, avatarURL)
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

	// 为图片添加 uid，避免名称冲突
	newFileName := uuid.New().String() + file.Filename
	dst := filepath.Join("public", "images", newFileName)

	// 保存图片到 pulic/images
	if err := c.SaveUploadedFile(file, dst); err != nil {
		utils.ServerError(c, fmt.Sprintf("保存图片失败：%s", err.Error()))
		return
	}

	// 更新头像
	avatarURL, err := h.accountService.UpdateAvatar(userID.(int), newFileName)
	if err != nil {
		utils.ServerError(c, "更新头像失败："+err.Error())
		return
	}
	utils.SuccessWithMessage(c, "更新成功", avatarURL)
}

func (h *AccountHandler) UpdateOrg(c *gin.Context) {
	// 从上下文获取用户组织
	org, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "用户组织获取失败")
		return
	}
	if org.(int) != 1 {
		utils.ServerError(c, "只有平台管理员可以更新用户组织")
		return
	}
	var req model.UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数格式错误")
		return
	}

	err := h.accountService.UpdateOrg(req.UserID, req.Org)
	if err != nil {
		utils.ServerError(c, "更新组织失败："+err.Error())
		return
	}
	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (h *AccountHandler) GetUserNameById(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		utils.BadRequest(c, "请求参数格式错误")
		return
	}
	userName, err := h.accountService.GetUserNameById(userID)
	if err != nil {
		utils.ServerError(c, "获取用户名失败："+err.Error())
		return
	}
	utils.Success(c, userName)
}
