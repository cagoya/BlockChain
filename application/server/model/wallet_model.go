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
	Time        time.Time `json:"time"`        // 转账时间
}

type TransferRequest struct {
	RecipientID int `json:"recipientId"` // 转入钱包ID
	Amount      int `json:"amount"`      // 转账金额
}
