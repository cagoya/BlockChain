package api

import (
	"application/service"
	"application/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MarketHandler struct {
	svc *service.MarketService
}

func NewMarketHandler() *MarketHandler {
	return &MarketHandler{svc: service.NewMarketService()}
}

// 1) 公开：查询挂牌
func (h *MarketHandler) ListListings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	items, total, err := h.svc.ListListings(page, size)
	if err != nil {
		utils.ServerError(c, "查询失败："+err.Error())
		return
	}
	utils.Success(c, gin.H{"items": items, "total": total})
}

// 2) 卖家创建挂牌（需要 JWT）
type createListingReq struct {
	AssetID  string  `json:"assetId" binding:"required"`
	Title    string  `json:"title" binding:"required"`
	Price    int64   `json:"price" binding:"required"`
	Deadline *string `json:"deadline"` // RFC3339
}

func (h *MarketHandler) CreateListing(c *gin.Context) {
	uidVal, ok := c.Get("userID")
	if !ok {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		utils.ServerError(c, "用户ID类型错误")
		return
	}

	var req createListingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数格式错误")
		return
	}

	var ddl *time.Time
	if req.Deadline != nil && *req.Deadline != "" {
		t, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			utils.BadRequest(c, "Deadline 格式应为 RFC3339")
			return
		}
		ddl = &t
	}

	l, err := h.svc.CreateListing(userID, req.AssetID, req.Title, req.Price, ddl)
	if err != nil {
		utils.ServerError(c, "创建挂牌失败："+err.Error())
		return
	}
	utils.SuccessWithMessage(c, "创建挂牌成功", l)
}

// 3) 买家出价（需要 JWT）
type createOfferReq struct {
	ListingID  uint  `json:"listingId" binding:"required"`
	OfferPrice int64 `json:"offerPrice" binding:"required"`
}

func (h *MarketHandler) CreateOffer(c *gin.Context) {
	uidVal, ok := c.Get("userID")
	if !ok {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		utils.ServerError(c, "用户ID类型错误")
		return
	}

	var req createOfferReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数格式错误")
		return
	}

	o, err := h.svc.CreateOffer(userID, req.ListingID, req.OfferPrice)
	if err != nil {
		utils.ServerError(c, "提交出价失败："+err.Error())
		return
	}
	utils.SuccessWithMessage(c, "提交出价成功", o)
}

// 4) 卖家接受出价（需要 JWT）
func (h *MarketHandler) AcceptOffer(c *gin.Context) {
	uidVal, ok := c.Get("userID")
	if !ok {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		utils.ServerError(c, "用户ID类型错误")
		return
	}

	offerID, err := strconv.Atoi(c.Param("id"))
	if err != nil || offerID <= 0 {
		utils.BadRequest(c, "出价 ID 非法")
		return
	}

	if err := h.svc.AcceptOffer(userID, uint(offerID)); err != nil {
		utils.ServerError(c, "接受出价失败："+err.Error())
		return
	}
	utils.SuccessWithMessage(c, "接受出价成功", nil)
}

// 5) 我提交的出价列表（需要 JWT）
func (h *MarketHandler) ListMyOffers(c *gin.Context) {
	uidVal, ok := c.Get("userID")
	if !ok {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		utils.ServerError(c, "用户ID类型错误")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	items, total, err := h.svc.ListMyOffers(userID, page, size)
	if err != nil {
		utils.ServerError(c, "查询失败："+err.Error())
		return
	}
	utils.Success(c, gin.H{"items": items, "total": total})
}

func (h *MarketHandler) CancelOffer(c *gin.Context) {
	uidVal, ok := c.Get("userID")
	if !ok {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	userID, ok := uidVal.(int)
	if !ok {
		utils.ServerError(c, "用户ID类型错误")
		return
	}

	oid, err := strconv.Atoi(c.Param("id"))
	if err != nil || oid <= 0 {
		utils.BadRequest(c, "出价 ID 非法")
		return
	}
	if err := h.svc.CancelOffer(userID, uint(oid)); err != nil {
		utils.ServerError(c, "撤回失败："+err.Error())
		return
	}
	utils.SuccessWithMessage(c, "撤回成功", nil)
}

// 一口价直购 BuyNow：如果 listing.BuyNowPrice 不为空
type buyNowReq struct {
	ListingID uint `json:"listingId" binding:"required"`
}

func (h *MarketHandler) BuyNow(c *gin.Context) {
	// 实现思路：对 BuyNowPrice 生成一条“合成的” PENDING 出价（先 WithHold），
	// 然后直接走 AcceptOffer 流程；为保证原子性可在 service 里做一个 BuyNow()
}
