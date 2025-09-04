package api

import (
	"application/model"
	"application/service"
	"application/utils"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService *service.WalletService
}

func NewWalletHandler() *WalletHandler {
	walletService := service.NewWalletService()
	return &WalletHandler{walletService: walletService}
}

func (h *WalletHandler) CreateAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	err := h.walletService.CreateAccount(userID.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, "钱包开通成功")
}

func (h *WalletHandler) GetBlance(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	blance, err := h.walletService.GetBlance(userID.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, blance)
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	var transferRequest model.TransferRequest
	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	recipientID := transferRequest.RecipientID
	amount := transferRequest.Amount
	err := h.walletService.Transfer(userID.(int), recipientID, amount)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, "转账成功")
}

func (h *WalletHandler) GetTransfer(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	transfers, err := h.walletService.GetTransfer(userID.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, transfers)
}
