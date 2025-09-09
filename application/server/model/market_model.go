package model

import "time"

// —— 列表状态常量 ——
const (
	ListingActive = "OPEN"   // 可出价
	ListingSold   = "SOLD"   // 已成交
	ListingClosed = "CLOSED" // 已下架
)

// —— 挂牌（Listing）——
type MarketListing struct {
	ID         int        `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetID    string     `json:"assetId" gorm:"type:varchar(128);not null;index"`
	Title      string     `json:"title" gorm:"type:varchar(200);not null"`
	Price      int64      `json:"price" gorm:"not null"`
	SellerID   int        `json:"sellerId" gorm:"not null"`
	SellerOrg  int32      `json:"sellerOrg" gorm:"not null;default:2"`
	Status     string     `json:"status" gorm:"type:varchar(16);not null;index"` // OPEN/SOLD/CLOSED
	Deadline   *time.Time `json:"deadline" gorm:"index"`
	CreateTime time.Time  `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime time.Time  `json:"updateTime" gorm:"autoUpdateTime"`
	// 销售模式
	IsAuction     bool   `json:"isAuction" gorm:"default:false"`
	ReservePrice  *int64 `json:"reservePrice"` // 低于此价不能成交（可选）
	BuyNowPrice   *int64 `json:"buyNowPrice"`  // 一口价（可选）
	WinnerOfferID *int   `json:"winnerOfferId" gorm:"index"`
}

func (MarketListing) TableName() string { return "market_listings" }

// —— 出价（Offer）——
const (
	OfferPending  = "PENDING"
	OfferAccepted = "ACCEPTED"
	OfferRejected = "REJECTED"
)

type MarketOffer struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ListingID  uint      `json:"listingId" gorm:"not null;index"`
	BidderID   int       `json:"bidderId" gorm:"not null"`
	BidderOrg  int32     `json:"sellerOrg" gorm:"not null;default:2"`
	OfferPrice int64     `json:"offerPrice" gorm:"not null"`
	Status     string    `json:"status" gorm:"type:varchar(16);not null;index"` // PENDING/ACCEPTED/REJECTED
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	// 托管与清算流水（链上 txid）
	IsEscrowed   bool    `json:"isEscrowed" gorm:"default:false"` // 是否已冻结成功
	EscrowHoldID *string `json:"escrowHoldId"`                    // 冻结返回的 holdId（业务关键）
	EscrowTxID   *string `json:"escrowTxId"`                      // WithHold 的 txid
	RefundTxID   *string `json:"refundTxId"`                      // Refund 的 txid（若被拒/撤回/过期）
	PayoutTxID   *string `json:"payoutTxId"`                      // Release 的 txid（若被接受）
}

func (MarketOffer) TableName() string { return "market_offers" }
