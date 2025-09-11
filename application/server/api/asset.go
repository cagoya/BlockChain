package api

import (
	"application/model"
	"application/service"
	"application/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"path/filepath"
	"strconv"
)

type AssetHandler struct {
	assetService *service.AssetService
}

func NewAssetHandler() *AssetHandler {
	return &AssetHandler{
		assetService: service.NewAssetService(model.GetDB()),
	}
}

func (h *AssetHandler) CreateAsset(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	org, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "组织信息获取失败")
		return
	}
	if org.(int) != 2 {
		utils.ServerError(c, "只有属于NFT创建者组织的用户可以上传NFT")
		return
	}
	// 为了保证原子性，图片必须与资产信息一起提交
	// 所以这里使用表单(form-data)提交
	// 参数就使用 PostForm 获取
	name := c.PostForm("name")
	if name == "" {
		utils.BadRequest(c, "请求参数错误")
		return
	}
	description := c.PostForm("description")
	if description == "" {
		description = "暂无描述"
	}
	image, err := c.FormFile("image")
	if err != nil {
		utils.ServerError(c, "获取请求参数失败")
		return
	}
	imageName := uuid.New().String() + image.Filename
	dst := filepath.Join(model.DefaultImageFolder, imageName)
	if err := c.SaveUploadedFile(image, dst); err != nil {
		utils.ServerError(c, fmt.Sprintf("保存图片失败：%s", err.Error()))
		return
	}
	// 创建时默认所有者是作者本人且是上传者
	asset, err := h.assetService.CreateAsset(name, imageName, userID.(int), userID.(int), description, org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, asset)
}

func (h *AssetHandler) GetAssetByID(c *gin.Context) {
	org, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "组织信息获取失败")
		return
	}
	id := c.Query("id")
	asset, err := h.assetService.GetAssetByID(id, org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, asset)
}

func (h *AssetHandler) GetAssetByAuthorID(c *gin.Context) {
	org, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "组织信息获取失败")
		return
	}
	authorId, err := strconv.Atoi(c.Query("authorId"))
	if err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}
	assets, err := h.assetService.GetAssetByAuthorID(authorId, org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, assets)
}

func (h *AssetHandler) GetAssetByOwnerID(c *gin.Context) {
	org, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "组织信息获取失败")
		return
	}
	ownerId, err := strconv.Atoi(c.Query("ownerId"))
	if err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}
	assets, err := h.assetService.GetAssetByOwnerID(ownerId, org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, assets)
}

func (h *AssetHandler) TransferAsset(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	org, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "组织信息获取失败")
		return
	}
	var req model.TransferAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}
	err := h.assetService.TransferAsset(req.ID, req.NewOwnerId, userID.(int), org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, nil)
}

func (h *AssetHandler) GetAssetStatus(c *gin.Context) {
	id := c.Query("id")
	status, err := h.assetService.GetAssetStatus(id)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, status)
}
