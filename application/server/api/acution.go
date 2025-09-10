package api

import (
	"application/model"
	"application/service"
	"application/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AuctionHandler struct {
	service *service.AuctionService
}

func NewAuctionHandler() *AuctionHandler {
	return &AuctionHandler{service: service.NewAuctionService(model.DB)}
}

// 创建拍卖品
func (h *AuctionHandler) CreateLot(c *gin.Context) {
	lotRequest := &model.CreateLotRequest{}
	c.ShouldBindJSON(lotRequest)
	sellerID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	sellerOrg, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "用户组织信息获取失败")
		return
	}
	err := h.service.CreateLot(lotRequest.AssetID, lotRequest.Title, lotRequest.ReservePrice,
		sellerID.(int), sellerOrg.(int), lotRequest.StartTime, lotRequest.Deadline)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	// 用于测试的定时器任务
	//time.AfterFunc(time.Until(lotRequest.Deadline), func() {
	//	fmt.Printf("定时器任务执行，当前时间为%s\n", time.Now().Format("2006-01-02 15:04:05"))
	//})
	//创建一个定时器任务，当拍卖到达结束时间时自动选择最高价
	// 查询自动创建的 LotID
	lot, err := h.service.GetLotByAssetID(lotRequest.AssetID)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	time.AfterFunc(time.Until(lotRequest.Deadline), func() {
		h.service.FinishAuction(lot.ID)
	})
	utils.SuccessWithMessage(c, "创建拍卖品成功", nil)
}

// 更新拍卖品
func (h *AuctionHandler) UpdateLotByID(c *gin.Context) {
	lotRequest := &model.UpdateLotRequest{}
	c.ShouldBindJSON(lotRequest)
	sellerID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	err := h.service.UpdateLotByID(sellerID.(int), lotRequest.AssetID, lotRequest.Title, lotRequest.ReservePrice)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "更新拍卖品成功", nil)
}

// 取消拍卖品
func (h *AuctionHandler) CancelLot(c *gin.Context) {
	assetID := c.Query("assetID")
	sellerID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	err := h.service.CancelLot(sellerID.(int), assetID)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "取消拍卖品成功", nil)
}

// 查询所有拍卖品
func (h *AuctionHandler) GetAllLots(c *gin.Context) {
	lots, err := h.service.GetAllLots()
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "查询拍卖品成功", lots)
}

// 查询个人拍卖品
func (h *AuctionHandler) GetLotBySellerID(c *gin.Context) {
	sellerID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	lots, err := h.service.GetLotBySellerID(sellerID.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "查询拍卖品成功", lots)
}

// 提交出价
func (h *AuctionHandler) SubmitBid(c *gin.Context) {
	bidRequest := &model.BidRequest{}
	c.ShouldBindJSON(bidRequest)
	bidderID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	bidderOrg, exists := c.Get("org")
	if !exists {
		utils.ServerError(c, "用户组织信息获取失败")
		return
	}
	err := h.service.SubmitBid(bidRequest.LotID, bidderID.(int), bidRequest.BidPrice, bidderOrg.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "提交出价成功", nil)
}

// 查询出价
func (h *AuctionHandler) GetBidPrice(c *gin.Context) {
	lotID, err := strconv.Atoi(c.Query("lotID"))
	if err != nil {
		utils.ServerError(c, "拍卖品ID不能为空")
		return
	}
	bidderID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	bidPrice, err := h.service.GetBidPrice(lotID, bidderID.(int))
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "查询出价成功", bidPrice)
}

// 查询最高出价
func (h *AuctionHandler) GetMaxBidPrice(c *gin.Context) {
	lotID, err := strconv.Atoi(c.Query("lotID"))
	if err != nil {
		utils.ServerError(c, "拍卖品ID不能为空")
		return
	}
	bidPrice, err := h.service.GetMaxBidPrice(lotID)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "查询最高出价成功", bidPrice)
}

// 结束拍卖不需要被前端调用，因为我们会在后端创建一个定时器任务
// 当拍卖到达结束时间时自动选择最高价
// 但是这里加这个接口用于调试
func (h *AuctionHandler) FinishAuction(c *gin.Context) {
	lotID, err := strconv.Atoi(c.Query("lotID"))
	if err != nil {
		utils.ServerError(c, "拍卖品ID不能为空")
		return
	}
	err = h.service.FinishAuction(lotID)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "结束拍卖成功", nil)
}

// 查询拍卖结果
func (h *AuctionHandler) GetAuctionResult(c *gin.Context) {
	lotID, err := strconv.Atoi(c.Query("lotID"))
	if err != nil {
		utils.ServerError(c, "拍卖品ID不能为空")
		return
	}
	auctionResult, err := h.service.GetAuctionResult(lotID)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "查询拍卖结果成功", auctionResult)
}
