package service

import (
	"application/model"
	"application/pkg/fabric"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type WalletService struct{}

func NewWalletService() *WalletService {
	return &WalletService{}
}

func (s *WalletService) CreateAccount(id int, org int) error {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	_, err = contract.SubmitTransaction("CreateAccount", fmt.Sprintf("%d", id))
	if err != nil {
		return fmt.Errorf("钱包开通失败：%s", fabric.ExtractErrorMessage(err))
	}
	return nil
}

func (s *WalletService) GetBalance(id int, org int) (int, error) {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return 0, fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	result, err := contract.EvaluateTransaction("GetBalance", fmt.Sprintf("%d", id))
	if err != nil {
		return 0, fmt.Errorf("获取余额失败：%s", fabric.ExtractErrorMessage(err))
	}

	var balance int
	if err := json.Unmarshal(result, &balance); err != nil {
		return 0, fmt.Errorf("解析余额失败：%v", err)
	}
	return balance, nil
}

func (s *WalletService) Transfer(senderId int, recipientId int, amount int, org int) error {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	_, err = contract.SubmitTransaction("Transfer", uuid.New().String(), fmt.Sprintf("%d", senderId), fmt.Sprintf("%d", recipientId), fmt.Sprintf("%d", amount), time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("转账失败：%s", fabric.ExtractErrorMessage(err))
	}
	return nil
}

func (s *WalletService) MintToken(accountID int, amount int, org int) error {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	_, err = contract.SubmitTransaction("MintToken", fmt.Sprintf("%d", accountID), fmt.Sprintf("%d", amount))
	if err != nil {
		return fmt.Errorf("铸币失败：%s", fabric.ExtractErrorMessage(err))
	}
	return nil
}

func (s *WalletService) GetTransferBySenderID(senderId int, org int) ([]model.Transfer, error) {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return nil, fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	results, err := contract.EvaluateTransaction("GetTransferBySenderID", fmt.Sprintf("%d", senderId))
	if err != nil {
		return nil, fmt.Errorf("获取转账记录失败：%s", fabric.ExtractErrorMessage(err))
	}
	if len(results) == 0 {
		return nil, nil
	}
	var transfers []model.Transfer
	if err := json.Unmarshal(results, &transfers); err != nil {
		return nil, fmt.Errorf("解析转账记录失败：%v", err)
	}
	return transfers, nil
}

func (s *WalletService) GetTransferByRecipientID(recipientId int, org int) ([]model.Transfer, error) {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return nil, fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	results, err := contract.EvaluateTransaction("GetTransferByRecipientID", fmt.Sprintf("%d", recipientId))
	if err != nil {
		return nil, fmt.Errorf("获取转账记录失败：%s", fabric.ExtractErrorMessage(err))
	}
	if len(results) == 0 {
		return nil, nil
	}
	var transfers []model.Transfer
	if err := json.Unmarshal(results, &transfers); err != nil {
		return nil, fmt.Errorf("解析转账记录失败：%v", err)
	}
	return transfers, nil
}

func (s *WalletService) WithHoldAccount(accountID int, listingID string, amount int, org int) error {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	_, err = contract.SubmitTransaction("WithHoldAccount", uuid.New().String(), fmt.Sprintf("%d", accountID), listingID, fmt.Sprintf("%d", amount), time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("预扣款失败：%s", fabric.ExtractErrorMessage(err))
	}
	return nil
}

func (s *WalletService) GetWithHoldingByAccountID(accountID int, org int) ([]model.WithHolding, error) {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return nil, fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	results, err := contract.EvaluateTransaction("GetWithHoldingByAccountID", fmt.Sprintf("%d", accountID))
	if err != nil {
		return nil, fmt.Errorf("获取预扣款记录失败：%s", fabric.ExtractErrorMessage(err))
	}
	if len(results) == 0 {
		return nil, nil
	}
	var withHoldings []model.WithHolding
	if err := json.Unmarshal(results, &withHoldings); err != nil {
		return nil, fmt.Errorf("解析预扣款记录失败：%v", err)
	}
	return withHoldings, nil
}

func (s *WalletService) GetWithHoldingByListingID(listingID string, org int) ([]model.WithHolding, error) {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return nil, fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	results, err := contract.EvaluateTransaction("GetWithHoldingByListingID", listingID)
	if err != nil {
		return nil, fmt.Errorf("获取预扣款记录失败：%s", fabric.ExtractErrorMessage(err))
	}
	if len(results) == 0 {
		return nil, nil
	}
	var withHoldings []model.WithHolding
	if err := json.Unmarshal(results, &withHoldings); err != nil {
		return nil, fmt.Errorf("解析预扣款记录失败：%v", err)
	}
	return withHoldings, nil
}

func (s *WalletService) ClearWithHolding(listingID string, org int) error {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	_, err = contract.SubmitTransaction("ClearWithHolding", listingID)
	if err != nil {
		return fmt.Errorf("清除预扣款失败：%s", fabric.ExtractErrorMessage(err))
	}
	return nil
}
