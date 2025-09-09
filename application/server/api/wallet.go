package api

import (
	"application/model"
	"application/service"
	"application/utils"
	"fmt"

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
	org, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "组织信息获取失败")
		return
	}
	err := h.walletService.CreateAccount(userID.(int), org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, "钱包开通成功")
}

// 在 WalletHandler 的 GetBalance 方法中增加账户自动创建逻辑
func (h *WalletHandler) GetBalance(c *gin.Context) {
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

	// 先尝试查询余额
	balance, err := h.walletService.GetBalance(userID.(int), org.(int))
	if err != nil {
		// 若查询失败，尝试创建账户
		createErr := h.walletService.CreateAccount(userID.(int), org.(int))
		if createErr != nil {
			utils.ServerError(c, fmt.Sprintf("余额查询失败且账户创建失败：%v", createErr))
			return
		}
		// 账户创建后重新查询余额
		balance, err = h.walletService.GetBalance(userID.(int), org.(int))
		if err != nil {
			utils.ServerError(c, err.Error())
			return
		}
	}
	utils.Success(c, balance)
}

func (h *WalletHandler) Transfer(c *gin.Context) {
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
	var transferRequest model.TransferRequest
	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	recipientId := transferRequest.RecipientID
	amount := transferRequest.Amount
	err := h.walletService.Transfer(userID.(int), recipientId, amount, org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, "转账成功")
}

func (h *WalletHandler) MintToken(c *gin.Context) {
	org, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "组织信息获取失败")
		return
	}
	if org.(int) != 3 {
		utils.ServerError(c, "只有金融组织可以铸币")
		return
	}
	var mintTokenRequest model.MintTokenRequest
	if err := c.ShouldBindJSON(&mintTokenRequest); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	err := h.walletService.MintToken(mintTokenRequest.AccountID, mintTokenRequest.Amount, org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, "铸币成功")
}

func (h *WalletHandler) GetTransferBySenderID(c *gin.Context) {
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
	transfers, err := h.walletService.GetTransferBySenderID(userID.(int), org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, transfers)
}

func (h *WalletHandler) GetTransferByRecipientID(c *gin.Context) {
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
	transfers, err := h.walletService.GetTransferByRecipientID(userID.(int), org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, transfers)
}

func (h *WalletHandler) WithHoldAccount(c *gin.Context) {
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
	var withHoldingRequest model.WithHoldingRequest
	if err := c.ShouldBindJSON(&withHoldingRequest); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	err := h.walletService.WithHoldAccount(userID.(int), withHoldingRequest.ListingID, withHoldingRequest.Amount, org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, "预扣款成功")
}

func (h *WalletHandler) GetWithHoldingByAccountID(c *gin.Context) {
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
	withHoldings, err := h.walletService.GetWithHoldingByAccountID(userID.(int), org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, withHoldings)
}

func (h *WalletHandler) GetWithHoldingByListingID(c *gin.Context) {
	org, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "组织信息获取失败")
		return
	}
	listingID := c.Query("listingID")
	withHoldings, err := h.walletService.GetWithHoldingByListingID(listingID, org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, withHoldings)
}

func (h *WalletHandler) ClearWithHolding(c *gin.Context) {
	org, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "组织信息获取失败")
		return
	}
	listingID := c.Query("listingID")
	err := h.walletService.ClearWithHolding(listingID, org.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, "清除预扣款成功")
}
