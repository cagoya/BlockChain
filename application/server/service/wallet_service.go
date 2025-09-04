package service

import (
	"application/model"
	"application/pkg/fabric"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	PLATFORM_ORG  = "org1" // 平台组织
	CREATOR_ORG   = "org2" // 创作者组织
	FINANCIAL_ORG = "org3" // 金融机构组织
)

type WalletService struct{}

func NewWalletService() *WalletService {
	return &WalletService{}
}

func (s *WalletService) CreateAccount(id int) error {
	contract := fabric.GetContract(PLATFORM_ORG)
	_, err := contract.SubmitTransaction("CreateAccount", fmt.Sprintf("%d", id))
	if err != nil {
		return fmt.Errorf("钱包开通失败：%s", fabric.ExtractErrorMessage(err))
	}
	return nil
}

func (s *WalletService) GetBlance(id int) (int, error) {
	contract := fabric.GetContract(CREATOR_ORG)
	result, err := contract.EvaluateTransaction("GetBlance", fmt.Sprintf("%d", id))
	if err != nil {
		return 0, fmt.Errorf("获取余额失败：%s", fabric.ExtractErrorMessage(err))
	}

	var balance int
	if err := json.Unmarshal(result, &balance); err != nil {
		return 0, fmt.Errorf("解析余额失败：%v", err)
	}
	return balance, nil
}

func (s *WalletService) Transfer(senderID int, recipientID int, amount int) error {
	contract := fabric.GetContract(CREATOR_ORG)
	_, err := contract.SubmitTransaction("Transfer", uuid.New().String(), fmt.Sprintf("%d", senderID), fmt.Sprintf("%d", recipientID), fmt.Sprintf("%d", amount), time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("转账失败：%s", fabric.ExtractErrorMessage(err))
	}
	return nil
}

func (s *WalletService) GetTransfer(senderID int) ([]model.Transfer, error) {
	contract := fabric.GetContract(CREATOR_ORG)
	results, err := contract.EvaluateTransaction("GetTransfer", fmt.Sprintf("%d", senderID))
	if err != nil {
		return nil, fmt.Errorf("获取转账记录失败：%s", fabric.ExtractErrorMessage(err))
	}
	var transfers []model.Transfer
	if err := json.Unmarshal(results, &transfers); err != nil {
		return nil, fmt.Errorf("解析转账记录失败：%v", err)
	}
	return transfers, nil
}
