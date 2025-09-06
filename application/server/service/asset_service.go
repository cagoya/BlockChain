package service

import (
	"application/model"
	"application/pkg/fabric"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type AssetService struct{}

func NewAssetService() *AssetService {
	return &AssetService{}
}

// 创建 nft 资产
func (s *AssetService) CreateAsset(name string, imageName string, authorId int,
	ownerId int, description string, org int) (model.Asset, error) {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return model.Asset{}, fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	uid := uuid.New().String()
	result, err := contract.SubmitTransaction("CreateAsset", uid, imageName, name, fmt.Sprintf("%d", authorId),
		fmt.Sprintf("%d", ownerId), description, time.Now().Format(time.RFC3339))
	if err != nil {
		return model.Asset{}, fmt.Errorf("创建 NFT 失败：%s", fabric.ExtractErrorMessage(err))
	}
	var asset model.Asset
	err = json.Unmarshal(result, &asset)
	if err != nil {
		return model.Asset{}, fmt.Errorf("解析数据失败：%s", err)
	}
	return asset, nil
}

func (s *AssetService) GetAssetByID(id string, org int) (model.Asset, error) {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return model.Asset{}, fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	result, err := contract.EvaluateTransaction("GetAssetByID", id)
	if err != nil {
		return model.Asset{}, fmt.Errorf("获取 NFT 失败：%s", fabric.ExtractErrorMessage(err))
	}
	var asset model.Asset
	err = json.Unmarshal(result, &asset)
	if err != nil {
		return model.Asset{}, fmt.Errorf("解析数据失败：%s", err)
	}
	return asset, nil
}

func (s *AssetService) GetAssetByAuthorID(authorId int, org int) ([]model.Asset, error) {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return nil, fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	results, err := contract.EvaluateTransaction("GetAssetByAuthorID", fmt.Sprintf("%d", authorId))
	if err != nil {
		return nil, fmt.Errorf("获取 NFT 失败：%s", fabric.ExtractErrorMessage(err))
	}
	if len(results) == 0 {
		return nil, nil
	}
	var assets []model.Asset
	err = json.Unmarshal(results, &assets)
	if err != nil {
		return nil, fmt.Errorf("解析数据失败：%s", err)
	}
	return assets, nil
}

func (s *AssetService) GetAssetByOwnerID(ownerId int, org int) ([]model.Asset, error) {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return nil, fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	results, err := contract.EvaluateTransaction("GetAssetByOwnerID", fmt.Sprintf("%d", ownerId))
	if err != nil {
		return nil, fmt.Errorf("获取 NFT 失败：%s", fabric.ExtractErrorMessage(err))
	}
	if len(results) == 0 {
		return nil, nil
	}
	var assets []model.Asset
	err = json.Unmarshal(results, &assets)
	if err != nil {
		return nil, fmt.Errorf("解析数据失败：%s", err)
	}
	return assets, nil
}

func (s *AssetService) TransferAsset(id string, newOwnerId int, userID int, org int) error {
	orgName, err := model.GetOrg(org)
	if err != nil {
		return fmt.Errorf("获取组织失败：%s", err)
	}
	contract := fabric.GetContract(orgName)
	_, err = contract.SubmitTransaction("TransferAsset", id, fmt.Sprintf("%d", newOwnerId), fmt.Sprintf("%d", userID), time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("转移NFT失败：%s", fabric.ExtractErrorMessage(err))
	}
	return nil
}
