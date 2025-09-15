package service

import (
	"application/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MarketService struct {
	db *gorm.DB
}

func NewMarketService() *MarketService {
	return &MarketService{db: model.GetDB()}
}

// 创建挂牌
func (s *MarketService) CreateListing(userID int, assetId, title string, price int64, deadline *time.Time) (*model.MarketListing, error) {
	// 所有人都在 org2
	const org2 = 2

	// 1) 链上校验：只有 NFT 当前 Owner 才能挂牌
	as := NewAssetService(model.GetDB())
	asset, err := as.GetAssetByID(assetId, org2)
	if err != nil {
		return nil, fmt.Errorf("查询NFT失败：%v", err)
	}
	if asset.OwnerId != userID {
		return nil, errors.New("只有NFT持有人才能挂牌")
	}

	// 2) 落库（显式记录 SellerOrg = 2）
	l := &model.MarketListing{
		AssetID:   assetId,
		Title:     title,
		Price:     price,
		SellerID:  userID,
		SellerOrg: 2,
		Deadline:  deadline,
		Status:    model.ListingActive,
	}
	if err := s.db.Create(l).Error; err != nil {
		return nil, err
	}
	return l, nil
}

// 查询挂牌（保持不变）
func (s *MarketService) ListListings(page, pageSize int) ([]model.MarketListing, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	var items []model.MarketListing
	var total int64
	now := time.Now()

	q := s.db.Model(&model.MarketListing{}).
		Where("status = ?", model.ListingActive).
		Where("deadline IS NULL OR deadline > ?", now) // ← 新增

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *MarketService) CreateOffer(userID int, listingId uint, offerPrice int64) (*model.MarketOffer, error) {
	// 1) 挂牌校验
	var listing model.MarketListing
	if err := s.db.First(&listing, listingId).Error; err != nil {
		return nil, err
	}
	if listing.Status != model.ListingActive {
		return nil, errors.New("挂牌不可出价（已成交或已下架）")
	}
	if listing.Deadline != nil && time.Now().After(*listing.Deadline) {
		return nil, errors.New("已过截止时间，不能出价")
	}
	if offerPrice <= 0 {
		return nil, errors.New("出价必须大于 0")
	}

	// 2) 余额校验（org 固定 2）
	const org2 = 2
	w := NewWalletService()
	bal, err := w.GetBalance(userID, org2)
	if err != nil {
		return nil, fmt.Errorf("查询钱包余额失败：%v", err)
	}
	if int64(bal) < offerPrice {
		return nil, errors.New("余额不足，无法提交出价")
	}

	// 3) 链码级冻结（listingId -> string）
	listingKey := fmt.Sprintf("%d", listingId)
	holdID, holdTx, err := w.WithHoldAccount(userID, listingKey, int(offerPrice), org2)
	if err != nil {
		return nil, fmt.Errorf("冻结失败：%v", err)
	}

	// 4) 落库
	o := &model.MarketOffer{
		ListingID:    listingId,
		BidderID:     userID,
		BidderOrg:    2,
		OfferPrice:   offerPrice,
		Status:       model.OfferPending,
		IsEscrowed:   true,
		EscrowHoldID: &holdID,
		EscrowTxID:   &holdTx,
	}
	if err := s.db.Create(o).Error; err != nil {
		return nil, err
	}
	return o, nil
}

func (s *MarketService) AcceptOffer(userID int, offerId uint) error {
	// A. 读取 & 校验
	var offer model.MarketOffer
	if err := s.db.First(&offer, offerId).Error; err != nil {
		return err
	}
	var listing model.MarketListing
	if err := s.db.First(&listing, offer.ListingID).Error; err != nil {
		return err
	}

	if listing.SellerID != userID {
		return errors.New("无权接受该出价")
	}
	if listing.Status != model.ListingActive {
		return errors.New("挂牌已非可售状态")
	}
	if offer.Status != model.OfferPending || !offer.IsEscrowed || offer.EscrowHoldID == nil {
		return errors.New("该出价不可被接受")
	}
	if listing.Deadline != nil && time.Now().After(*listing.Deadline) {
		return errors.New("已过截止时间，无法接受出价")
	}
	if listing.IsAuction && listing.ReservePrice != nil && offer.OfferPrice < *listing.ReservePrice {
		return errors.New("未达保留价，无法成交")
	}

	// B. 先链上清算（org 固定 2）
	const org2 = 2
	w := NewWalletService()

	// 1) 赢家释放到卖家
	listingKey := fmt.Sprintf("%d", offer.ListingID)
	payoutTx, err := w.ReleaseHolding(listingKey, listing.SellerID, int(offer.OfferPrice), org2)
	if err != nil {
		return fmt.Errorf("释放失败（向卖家结算）：%v", err)
	}

	// 2) 其他人退款
	var others []model.MarketOffer
	if err := s.db.
		Where("listing_id = ? AND id <> ? AND status = ?", offer.ListingID, offer.ID, model.OfferPending).
		Find(&others).Error; err != nil {
		return err
	}
	refundTxMap := make(map[int]string, len(others))
	for _, o := range others {
		okey := fmt.Sprintf("%d", o.ListingID)
		rtx, e := w.RefundHolding(okey, o.BidderID, int(o.OfferPrice), org2)
		if e != nil {
			return fmt.Errorf("退款失败（offer %d）：%v", o.ID, e)
		}
		refundTxMap[o.ID] = rtx
	}

	// 2.5) 清算成功后 → 转移 NFT（卖家 → 赢家）
	as := NewAssetService(model.GetDB())
	if err := as.TransferAsset(listing.AssetID, offer.BidderID, listing.SellerID, org2); err != nil {
		return fmt.Errorf("NFT转移失败：%v", err)
	}

	// C. 事务落库
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&listing, listing.ID).Error; err != nil {
			return err
		}
		if listing.Status != model.ListingActive {
			return errors.New("挂牌状态已改变，无法成交")
		}
		now := time.Now()

		// 赢家
		if err := tx.Model(&model.MarketOffer{}).
			Where("id = ? AND status = ?", offer.ID, model.OfferPending).
			Updates(map[string]interface{}{
				"status":       model.OfferAccepted,
				"payout_tx_id": payoutTx,
				"update_time":  now,
			}).Error; err != nil {
			return err
		}
		// 其他人
		for _, o := range others {
			if err := tx.Model(&model.MarketOffer{}).
				Where("id = ? AND status = ?", o.ID, model.OfferPending).
				Updates(map[string]interface{}{
					"status":       model.OfferRejected,
					"is_escrowed":  false,
					"refund_tx_id": refundTxMap[o.ID],
					"update_time":  now,
				}).Error; err != nil {
				return err
			}
		}
		// 挂牌 SOLD
		winnerID := int(offer.ID)
		if err := tx.Model(&model.MarketListing{}).
			Where("id = ? AND status = ?", listing.ID, model.ListingActive).
			Updates(map[string]interface{}{
				"status":          model.ListingSold,
				"winner_offer_id": &winnerID,
				"update_time":     now,
			}).Error; err != nil {
			return err
		}
		return nil
	})
}

// 撤回出价（只能撤回 PENDING 状态的出价）
func (s *MarketService) CancelOffer(userID int, offerId uint) error {
	var o model.MarketOffer
	if err := s.db.First(&o, offerId).Error; err != nil {
		return err
	}
	if o.BidderID != userID {
		return errors.New("无权撤回该出价")
	}
	var listing model.MarketListing
	if err := s.db.First(&listing, o.ListingID).Error; err != nil {
		return err
	}
	if listing.Status != model.ListingActive {
		return errors.New("挂牌不可撤回")
	}
	if o.Status != model.OfferPending || !o.IsEscrowed || o.EscrowHoldID == nil {
		return errors.New("该出价不可撤回")
	}
	if listing.Deadline != nil && time.Now().After(*listing.Deadline) {
		return errors.New("已过截止时间，不能撤回")
	}

	// 链上退款
	w := NewWalletService()
	listingKey := fmt.Sprintf("%d", o.ListingID)
	rtx, err := w.RefundHolding(listingKey, o.BidderID, int(o.OfferPrice), int(o.BidderOrg))
	if err != nil {
		return fmt.Errorf("退款失败：%v", err)
	}

	// 落库
	return s.db.Model(&model.MarketOffer{}).
		Where("id = ? AND status = ?", o.ID, model.OfferPending).
		Updates(map[string]any{
			"status":       model.OfferRejected,
			"is_escrowed":  false,
			"refund_tx_id": rtx,
			"update_time":  time.Now(),
		}).Error
}

// 定时任务：关闭过期挂牌（若有），并退款所有 PENDING 出价
func (s *MarketService) CloseExpired() error {
	now := time.Now()
	var listings []model.MarketListing
	if err := s.db.
		Where("status = ? AND deadline IS NOT NULL AND deadline < ?", model.ListingActive, now).
		Find(&listings).Error; err != nil {
		return err
	}

	w := NewWalletService()

	for _, l := range listings {
		// 查该 listing 下仍 PENDING 的出价
		var offs []model.MarketOffer
		if err := s.db.Where("listing_id = ? AND status = ?", l.ID, model.OfferPending).
			Find(&offs).Error; err != nil {
			return err
		}

		// 全部退款（链上）
		refundMap := map[int]string{}
		for _, o := range offs {
			okey := fmt.Sprintf("%d", o.ListingID)
			rtx, e := w.RefundHolding(okey, o.BidderID, int(o.OfferPrice), int(o.BidderOrg))
			if e != nil {
				return fmt.Errorf("listing %d 退款失败（offer %d）：%v", l.ID, o.ID, e)
			}
			refundMap[o.ID] = rtx
		}

		// 事务：更新 offers & listing
		if err := s.db.Transaction(func(tx *gorm.DB) error {
			// offers
			for _, o := range offs {
				if err := tx.Model(&model.MarketOffer{}).
					Where("id = ? AND status = ?", o.ID, model.OfferPending).
					Updates(map[string]any{
						"status":       model.OfferRejected,
						"is_escrowed":  false,
						"refund_tx_id": refundMap[o.ID],
						"update_time":  now,
					}).Error; err != nil {
					return err
				}
			}
			// listing CLOSED
			return tx.Model(&model.MarketListing{}).
				Where("id = ? AND status = ?", l.ID, model.ListingActive).
				Updates(map[string]any{
					"status":      model.ListingClosed,
					"update_time": now,
				}).Error
		}); err != nil {
			return err
		}
	}
	return nil
}

// 一口价直购
func (s *MarketService) BuyNow(buyerID int, listingId uint) error {
	const org2 = 2
	w := NewWalletService()

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1) 锁定 & 校验
		var l model.MarketListing
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&l, listingId).Error; err != nil {
			return err
		}
		if l.Status != model.ListingActive {
			return errors.New("挂牌已不可售")
		}
		if l.Deadline != nil && time.Now().After(*l.Deadline) {
			return errors.New("已过截止时间，无法购买")
		}
		if l.SellerID == buyerID {
			return errors.New("不能购买自己的挂牌")
		}
		if l.Price <= 0 {
			return errors.New("价格非法")
		}

		// 2) 余额校验
		bal, err := w.GetBalance(buyerID, org2)
		if err != nil {
			return fmt.Errorf("查询余额失败：%v", err)
		}
		if int64(bal) < l.Price {
			return errors.New("余额不足")
		}

		// 3) 直接链上转账给卖家
		txid, err := w.Transfer(buyerID, l.SellerID, int(l.Price), org2)
		if err != nil {
			return fmt.Errorf("直接转账失败：%v", err)
		}

		// 4) NFT 过户（卖家 -> 买家）
		as := NewAssetService(model.GetDB())
		if err := as.TransferAsset(l.AssetID, buyerID, l.SellerID, org2); err != nil {
			// 如果资产转移失败，理论上要退款给买家（这里可以再调 w.Transfer(l.SellerID,buyerID,...))
			return fmt.Errorf("NFT转移失败：%v", err)
		}

		// 5) 插入成交记录
		off := &model.MarketOffer{
			ListingID:  uint(l.ID),
			BidderID:   buyerID,
			BidderOrg:  org2,
			OfferPrice: l.Price,
			Status:     model.OfferAccepted,
			// 不再有 escrow 字段
			PayoutTxID: &txid,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		if err := tx.Create(off).Error; err != nil {
			return fmt.Errorf("保存成交记录失败：%v", err)
		}

		// 6) 标记挂牌 SOLD
		winnerID := uint(off.ID)
		if err := tx.Model(&model.MarketListing{}).
			Where("id = ? AND status = ?", l.ID, model.ListingActive).
			Updates(map[string]any{
				"status":          model.ListingSold,
				"winner_offer_id": &winnerID,
				"update_time":     time.Now(),
			}).Error; err != nil {
			return fmt.Errorf("更新挂牌状态失败：%v", err)
		}

		return nil
	})
}

// 我提交的出价（保持不变）
func (s *MarketService) ListMyOffers(userID int, page, pageSize int) ([]model.MarketOffer, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	var (
		items []model.MarketOffer
		total int64
	)
	q := s.db.Model(&model.MarketOffer{}).Where("bidder_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
