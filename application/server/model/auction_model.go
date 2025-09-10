package model

import (
	"time"
)

type Lot struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`              // 拍品ID
	AssetID      string    `json:"assetId" gorm:"type:varchar(128);not null;index"` // 拍品资产ID
	Title        string    `json:"title" gorm:"type:varchar(200);not null"`         // 拍品标题
	ReservePrice int       `json:"reservePrice" gorm:"not null"`                    // 起拍价
	CurrentPrice int       `json:"currentPrice" gorm:"not null"`                    // 当前价
	SellerID     int       `json:"sellerId" gorm:"not null"`                        // 卖家ID
	SellerOrg    int       `json:"sellerOrg" gorm:"not null;default:2"`             // 卖家组织
	StartTime    time.Time `json:"startTime" gorm:"not null;"`                      // 开始时间
	Deadline     time.Time `json:"deadline" gorm:"not null;"`                       // 结束时间
	Valid        bool      `json:"valid" gorm:"default:true"`                       // 是否有效
	CreateTime   time.Time `json:"createTime" gorm:"autoCreateTime"`                // 创建时间
	UpdateTime   time.Time `json:"updateTime" gorm:"autoUpdateTime"`                // 更新时间
}

func (Lot) TableName() string { return "lots" }

type Bid struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`  // 出价ID
	LotID      int       `json:"lotId" gorm:"not null;index"`         // 拍品ID
	BidderID   int       `json:"bidderId" gorm:"not null"`            // 出价者ID
	BidderOrg  int       `json:"bidderOrg" gorm:"not null;default:2"` // 出价者组织
	BidPrice   int       `json:"bidPrice" gorm:"not null"`            // 出价
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`    // 创建时间
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime"`    // 更新时间
}

func (Bid) TableName() string { return "bids" }

type AuctionResult struct {
	LotID    int `json:"lotId"`    // 拍品ID
	BidPrice int `json:"bidPrice"` // 出价
	BidderID int `json:"bidderId"` // 出价者ID，如果为空则为流拍
}

func (AuctionResult) TableName() string { return "auction_results" }

type CreateLotRequest struct {
	AssetID      string    `json:"assetId"`      // 拍品资产ID
	Title        string    `json:"title"`        // 拍品标题
	ReservePrice int       `json:"reservePrice"` // 起拍价
	StartTime    time.Time `json:"startTime"`    // 开始时间
	Deadline     time.Time `json:"deadline"`     // 结束时间
}

type UpdateLotRequest struct {
	AssetID      string `json:"assetId"`      // 拍品资产ID
	Title        string `json:"title"`        // 拍品标题
	ReservePrice int    `json:"reservePrice"` // 起拍价
}

type BidRequest struct {
	LotID    int `json:"lotId"`    // 拍品ID
	BidPrice int `json:"bidPrice"` // 出价
}
