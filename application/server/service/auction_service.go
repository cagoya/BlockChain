package service

import (
	"application/model"
	"application/pkg/fabric"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuctionService struct {
	db *gorm.DB
}

func NewAuctionService(db *gorm.DB) *AuctionService {
	return &AuctionService{db: db}
}

func (s *AuctionService) CreateLot(AssetID string, Title string, ReservePrice int, SellerID int,
	SellerOrg int, StartTime time.Time, Deadline time.Time) error {
	// 截止时间必须晚于开始时间
	if Deadline.Before(StartTime) {
		return fmt.Errorf("截止时间必须晚于开始时间")
	}
	// 开始时间必须晚于当前时间
	if StartTime.Before(time.Now()) {
		return fmt.Errorf("开始时间必须晚于当前时间")
	}
	// 检查拍卖品是否已经存在
	// 同一件商品同时只能对应一个有效拍品
	lot := model.Lot{}
	err := s.db.Where("asset_id = ? and deadline > ?", AssetID, time.Now()).First(&lot).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			lot = model.Lot{
				AssetID:      AssetID,
				Title:        Title,
				ReservePrice: ReservePrice,
				CurrentPrice: ReservePrice,
				SellerID:     SellerID,
				SellerOrg:    SellerOrg,
				StartTime:    StartTime,
				Deadline:     Deadline,
			}
			return s.db.Create(&lot).Error
		}
		return fmt.Errorf("查询拍品失败：%v", err)
	}
	return fmt.Errorf("拍卖品已存在")
}

func (s *AuctionService) GetLotBySellerID(SellerID int) ([]model.Lot, error) {
	var lots []model.Lot
	err := s.db.Where("seller_id = ?", SellerID).Find(&lots).Error
	if err != nil {
		return nil, fmt.Errorf("查询拍品失败：%v", err)
	}
	return lots, nil
}

func (s *AuctionService) GetLotByAssetID(AssetID string) (model.Lot, error) {
	lot := model.Lot{}
	err := s.db.Where("asset_id = ? and deadline > ?", AssetID, time.Now()).First(&lot).Error
	if err != nil {
		return model.Lot{}, fmt.Errorf("查询拍品失败：%v", err)
	}
	return lot, nil
}

func (s *AuctionService) GetAllLots() ([]model.Lot, error) {
	var lots []model.Lot
	err := s.db.Find(&lots).Error
	if err != nil {
		return nil, fmt.Errorf("查询拍品失败：%v", err)
	}
	return lots, nil
}

func (s *AuctionService) SubmitBid(LotID int, BidderID int, BidPrice int, BidderOrg int) error {
	// 检查出价是否高于当前价
	lot := model.Lot{}
	err := s.db.Where("id = ?", LotID).First(&lot).Error
	if err != nil {
		return fmt.Errorf("查询拍品失败：%v", err)
	}
	// 拍卖未开始或者已结束
	if lot.Deadline.Before(time.Now()) || lot.StartTime.After(time.Now()) {
		return fmt.Errorf("拍卖品未开始或者已结束，不允许出价")
	}
	if BidPrice <= lot.CurrentPrice {
		return fmt.Errorf("出价必须高于当前价")
	}
	// 更新拍品当前价
	lot.CurrentPrice = BidPrice
	err = s.db.Save(lot).Error
	if err != nil {
		return fmt.Errorf("更新拍品当前价失败：%v", err)
	}
	// 检查该用户是否已有出价
	// 如果不存在则创建
	// 如果存在则更新
	bid := model.Bid{}
	err = s.db.Where("lot_id = ? and bidder_id = ?", LotID, BidderID).First(&bid).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			bid = model.Bid{
				LotID:     LotID,
				BidderID:  BidderID,
				BidderOrg: BidderOrg,
				BidPrice:  BidPrice,
			}
			return s.db.Create(&bid).Error
		}
		return fmt.Errorf("查询出价失败：%v", err)
	}
	bid.BidPrice = BidPrice
	err = s.db.Save(bid).Error
	if err != nil {
		return fmt.Errorf("更新出价失败：%v", err)
	}
	return nil
}

func (s *AuctionService) GetBidPrice(LotID int, BidderID int) (int, error) {
	bid := model.Bid{}
	err := s.db.Where("lot_id = ? and bidder_id = ?", LotID, BidderID).First(&bid).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("出价不存在")
		}
		return 0, fmt.Errorf("查询出价失败：%v", err)
	}
	return bid.BidPrice, nil
}

func (s *AuctionService) FinishAuction(LotID int) error {
	lot := model.Lot{}
	bid := model.Bid{}
	// 查询当前最高出价者
	err := s.db.Where("lot_id = ?", LotID).Order("bid_price desc").First(&bid).Error
	if err != nil {
		// 无人出价，交易失败
		if err == gorm.ErrRecordNotFound {
			auctionResult := model.AuctionResult{
				LotID:    LotID,
				BidPrice: 0,
				BidderID: 0,
			}
			err = s.db.Create(&auctionResult).Error
			if err != nil {
				return fmt.Errorf("记录拍卖结果失败：%v", err)
			}
			return fmt.Errorf("当前无出价，交易失败")
		}
		return fmt.Errorf("查询最高出价失败：%v", err)
	}
	// 查询拍品
	err = s.db.Where("id = ?", LotID).First(&lot).Error
	if err != nil {
		return fmt.Errorf("查询拍品失败：%v", err)
	}
	orgName, err := model.GetOrg(bid.BidderOrg)
	if err != nil {
		return fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	// 转移 NFT 的所有权
	_, err = contract.SubmitTransaction("TransferAsset", lot.AssetID, fmt.Sprintf("%d", bid.BidderID),
		fmt.Sprintf("%d", lot.SellerID), time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("转移 NFT 的所有权失败：%v", err)
	}
	// 转账
	_, err = contract.SubmitTransaction("Transfer", uuid.New().String(), fmt.Sprintf("%d", bid.BidderID), fmt.Sprintf("%d", lot.SellerID),
		fmt.Sprintf("%d", bid.BidPrice), time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("转账失败：%v", err)
	}
	// 记录拍卖结果
	auctionResult := model.AuctionResult{
		LotID:    LotID,
		BidPrice: bid.BidPrice,
		BidderID: bid.BidderID,
	}
	return s.db.Create(&auctionResult).Error
}

func (s *AuctionService) GetAuctionResult(LotID int) (model.AuctionResult, error) {
	auctionResult := model.AuctionResult{}
	// 查询拍卖是否已经截止，未截止不能查询结果
	lot := model.Lot{}
	err := s.db.Where("id = ?", LotID).First(&lot).Error
	if err != nil {
		return model.AuctionResult{}, fmt.Errorf("查询拍卖品失败：%v", err)
	}
	if lot.Deadline.After(time.Now()) {
		return model.AuctionResult{}, fmt.Errorf("拍卖品未截止，不能查询结果")
	}
	err = s.db.Where("lot_id = ?", LotID).First(&auctionResult).Error
	if err != nil {
		return model.AuctionResult{}, fmt.Errorf("查询拍卖结果失败：%v", err)
	}
	return auctionResult, nil
}
