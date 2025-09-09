package model

import (
	"time"
)

type Wallet struct {
	ID      int `json:"id"`      // 钱包ID，等于账号ID
	Balance int `json:"balance"` // 钱包余额，这里按照平台代币计数
}

type Transfer struct {
	ID          string    `json:"id"`          // 转账ID
	SenderID    int       `json:"senderId"`    // 转出钱包ID
	RecipientID int       `json:"recipientId"` // 转入钱包ID
	Amount      int       `json:"amount"`      // 转账金额
	TimeStamp   time.Time `json:"timeStamp"`   // 转账时间
}

type TransferRequest struct {
	RecipientID int `json:"recipientId"` // 转入钱包ID
	Amount      int `json:"amount"`      // 转账金额
}

type WithHolding struct {
	ID        string    `json:"id"`        // 预扣款ID
	AccountID int       `json:"accountId"` // 预扣款账号ID
	ListingID string    `json:"listingId"` // 预扣款商品ID
	Amount    int       `json:"amount"`    // 预扣款金额
	TimeStamp time.Time `json:"timeStamp"` // 预扣款时间
}

type WithHoldingRequest struct {
	ListingID string `json:"listingId"` // 预扣款商品ID
	Amount    int    `json:"amount"`    // 预扣款金额
}

type MintTokenRequest struct {
	AccountID int `json:"accountId"` // 铸币账号ID
	Amount    int `json:"amount"`    // 铸币金额
}
