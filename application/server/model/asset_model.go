package model

import (
	"time"
)

// Asset 资产信息
// 先不管稀有度，因为稀有度应该由平台给定，而不是由上传用户给定
type Asset struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ImageName   string    `json:"imageName"`
	AuthorId    int       `json:"authorId"`
	OwnerId     int       `json:"ownerId"`
	Description string    `json:"description"`
	TimeStamp   time.Time `json:"timeStamp"`
}

type TransferAssetRequest struct {
	ID         string `json:"id"`
	NewOwnerId int    `json:"newOwnerId"`
}
