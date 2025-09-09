package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/v2/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract 提供NFT交易的功能
type SmartContract struct {
	contractapi.Contract
}

// 常量，用于构建复合键
// 为了保证不重复，每个不同的键值对都要有自己的前缀
const (
	ACCOUNT_KEY       = "account"
	TRANSFER_KEY      = "transfer"
	SENDER_KEY        = "sender"
	RECIPIENT_KEY     = "recipient"
	WITH_HOLDING_KEY1 = "withHolding1"
	WITH_HOLDING_KEY2 = "withHolding2"
	ASSET_KEY1        = "asset1"
	ASSET_KEY2        = "asset2"
	ASSET_KEY3        = "asset3"
)

// Account 账户信息
type Account struct {
	ID      int `json:"id"`
	Balance int `json:"balance"`
}

// 转账记录
type Transfer struct {
	ID          string    `json:"id"`
	SenderID    int       `json:"senderId"`
	RecipientID int       `json:"recipientId"`
	Amount      int       `json:"amount"`
	TimeStamp   time.Time `json:"timeStamp"`
}

// 预扣款
type WithHolding struct {
	ID        string    `json:"id"`
	AccountID int       `json:"accountID"`
	ListingID string    `json:"listingID"`
	Amount    int       `json:"amount"`
	TimeStamp time.Time `json:"timeStamp"`
}

// asset
type Asset struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ImageName   string    `json:"imageName"`
	AuthorId    int       `json:"authorId"`
	OwnerId     int       `json:"ownerId"`
	Description string    `json:"description"`
	Rarity      string    `json:"rarity"`
	TimeStamp   time.Time `json:"timeStamp"`
}

// QueryResult 分页查询结果
type QueryResult struct {
	Records             []interface{} `json:"records"`             // 记录列表
	RecordsCount        int32         `json:"recordsCount"`        // 本次返回的记录数
	Bookmark            string        `json:"bookmark"`            // 书签，用于下一页查询
	FetchedRecordsCount int32         `json:"fetchedRecordsCount"` // 总共获取的记录数
}

// 组织 MSP ID 常量
const (
	PLATFORM_ORG_MSPID = "Org1MSP" // 平台组织 MSP ID
	CREATOR_ORG_MSPID  = "Org2MSP" // NFT 创建者组织 MSP ID
	FINANCE_ORG_MSPID  = "Org3MSP" // 金融组织 MSP ID
)

// 通用方法: 获取客户端身份信息
func (s *SmartContract) getClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientID, err := cid.New(ctx.GetStub())
	if err != nil {
		return "", fmt.Errorf("获取客户端身份信息失败：%v", err)
	}
	return clientID.GetMSPID()
}

// 通用方法：创建和获取复合键
func (s *SmartContract) getCompositeKey(ctx contractapi.TransactionContextInterface, objectType string, attributes []string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(objectType, attributes)
	if err != nil {
		return "", fmt.Errorf("创建复合键失败：%v", err)
	}
	return key, nil
}

// 通用方法：获取状态
func (s *SmartContract) getState(ctx contractapi.TransactionContextInterface, key string, value interface{}) error {
	bytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("读取状态失败：%v", err)
	}
	if bytes == nil {
		return fmt.Errorf("键 %s 不存在", key)
	}

	err = json.Unmarshal(bytes, value)
	if err != nil {
		return fmt.Errorf("解析数据失败：%v", err)
	}
	return nil
}

// 通用方法：保存状态
func (s *SmartContract) putState(ctx contractapi.TransactionContextInterface, key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化数据失败：%v", err)
	}

	err = ctx.GetStub().PutState(key, bytes)
	if err != nil {
		return fmt.Errorf("保存状态失败：%v", err)
	}
	return nil
}

// Hello 用于验证
func (s *SmartContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "hello", nil
}

// 创建账户信息
func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, id int) error {
	// 创建复合键
	key, err := s.getCompositeKey(ctx, ACCOUNT_KEY, []string{fmt.Sprintf("%d", id)})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	// 检查是否已经存在
	var account Account
	err = s.getState(ctx, key, &account)
	if err == nil {
		return fmt.Errorf("账户已存在")
	}
	// 初始赠送 100 代币
	err = s.putState(ctx, key, Account{ID: id, Balance: 100})
	if err != nil {
		return err
	}
	return nil
}

// 获取余额
func (s *SmartContract) GetBalance(ctx contractapi.TransactionContextInterface, id int) (int, error) {
	var account Account
	key, err := s.getCompositeKey(ctx, ACCOUNT_KEY, []string{fmt.Sprintf("%d", id)})
	if err != nil {
		return 0, fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.getState(ctx, key, &account)
	if err != nil {
		return 0, fmt.Errorf("查询余额失败：%v", err)
	}
	return account.Balance, nil
}

// 转账
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, id string, senderId int, recipientId int, amount int, timeStamp time.Time) error {
	// 转账金额检查
	if amount <= 0 {
		return fmt.Errorf("转账金额必须大于 0")
	}
	// 检查发送方和接收方是否是同一个账户
	if senderId == recipientId {
		return fmt.Errorf("发送方和接收方不能是同一个账户")
	}
	var senderAccount Account
	key1, err := s.getCompositeKey(ctx, ACCOUNT_KEY, []string{fmt.Sprintf("%d", senderId)})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.getState(ctx, key1, &senderAccount)
	if err != nil {
		return fmt.Errorf("查询发送方账户失败：%v", err)
	}
	var recipientAccount Account
	key2, err := s.getCompositeKey(ctx, ACCOUNT_KEY, []string{fmt.Sprintf("%d", recipientId)})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.getState(ctx, key2, &recipientAccount)
	if err != nil {
		return fmt.Errorf("查询接收方账户失败：%v", err)
	}
	// 发送方余额检查
	if senderAccount.Balance < amount {
		return fmt.Errorf("发送方账户 %d 余额不足", senderId)
	}
	senderAccount.Balance -= amount
	recipientAccount.Balance += amount
	// 移除旧的记录
	err = ctx.GetStub().DelState(key1)
	if err != nil {
		return fmt.Errorf("移除发送方账户状态失败：%v", err)
	}
	err = ctx.GetStub().DelState(key2)
	if err != nil {
		return fmt.Errorf("移除接收方账户状态失败：%v", err)
	}
	// 添加新记录
	err = s.putState(ctx, key1, senderAccount)
	if err != nil {
		return fmt.Errorf("更新发送方账户状态失败：%v", err)
	}
	err = s.putState(ctx, key2, recipientAccount)
	if err != nil {
		return fmt.Errorf("更新接收方账户状态失败：%v", err)
	}
	// 添加转账记录
	transfer := Transfer{
		ID:          id,
		SenderID:    senderId,
		RecipientID: recipientId,
		Amount:      amount,
		TimeStamp:   timeStamp,
	}
	// 转账记录需要存两份，一份主键是发送方，一份主键是接收方
	// 创建复合键(SenderID, ID)
	key3, err := s.getCompositeKey(ctx, SENDER_KEY, []string{fmt.Sprintf("%d", transfer.SenderID), transfer.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.putState(ctx, key3, transfer)
	if err != nil {
		return fmt.Errorf("保存转账记录失败：%v", err)
	}
	// 创建复合键(RecipientID, ID)
	key4, err := s.getCompositeKey(ctx, RECIPIENT_KEY, []string{fmt.Sprintf("%d", transfer.RecipientID), transfer.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.putState(ctx, key4, transfer)
	if err != nil {
		return fmt.Errorf("保存转账记录失败：%v", err)
	}
	return nil
}

// 铸币，暂时不存记录
func (s *SmartContract) MintToken(ctx contractapi.TransactionContextInterface, accountID int, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("铸币金额必须大于 0")
	}
	var account Account
	key, err := s.getCompositeKey(ctx, ACCOUNT_KEY, []string{fmt.Sprintf("%d", accountID)})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.getState(ctx, key, &account)
	if err != nil {
		return fmt.Errorf("查询账户失败：%v", err)
	}
	account.Balance += amount
	err = ctx.GetStub().DelState(key)
	if err != nil {
		return fmt.Errorf("移除旧的记录失败：%v", err)
	}
	err = s.putState(ctx, key, account)
	if err != nil {
		return fmt.Errorf("更新账户余额失败：%v", err)
	}
	return nil
}

// 查询某个账户的转账转出记录
func (s *SmartContract) GetTransferBySenderID(ctx contractapi.TransactionContextInterface, senderId int) ([]Transfer, error) {
	var transfers []Transfer
	results, err := ctx.GetStub().GetStateByPartialCompositeKey(SENDER_KEY, []string{fmt.Sprintf("%d", senderId)})
	if err != nil {
		return nil, fmt.Errorf("查询转账记录失败：%v", err)
	}
	for results.HasNext() {
		var transfer Transfer
		result, err := results.Next()
		if err != nil {
			return nil, fmt.Errorf("查询转账记录失败：%v", err)
		}
		err = json.Unmarshal(result.Value, &transfer)
		if err != nil {
			return nil, fmt.Errorf("解析数据失败：%v", err)
		}
		transfers = append(transfers, transfer)
	}
	return transfers, nil
}

// 查询某个账户的转账转入记录
func (s *SmartContract) GetTransferByRecipientID(ctx contractapi.TransactionContextInterface, recipientId int) ([]Transfer, error) {
	var transfers []Transfer
	results, err := ctx.GetStub().GetStateByPartialCompositeKey(RECIPIENT_KEY, []string{fmt.Sprintf("%d", recipientId)})
	if err != nil {
		return nil, fmt.Errorf("查询转账记录失败：%v", err)
	}
	for results.HasNext() {
		var transfer Transfer
		result, err := results.Next()
		if err != nil {
			return nil, fmt.Errorf("查询转账记录失败：%v", err)
		}
		err = json.Unmarshal(result.Value, &transfer)
		if err != nil {
			return nil, fmt.Errorf("解析数据失败：%v", err)
		}
		transfers = append(transfers, transfer)
	}
	return transfers, nil
}

// 预扣款一定金额
func (s *SmartContract) WithHoldAccount(ctx contractapi.TransactionContextInterface, id string, accountId int, listingID string, amount int, timeStamp time.Time) error {
	// 检查 ammount 是否大于 0
	if amount <= 0 {
		return fmt.Errorf("预扣款金额必须大于 0")
	}
	var account Account
	key1, err := s.getCompositeKey(ctx, ACCOUNT_KEY, []string{fmt.Sprintf("%d", accountId)})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.getState(ctx, key1, &account)
	if err != nil {
		return fmt.Errorf("查询账户失败：%v", err)
	}
	// 检查余额是否足够
	if account.Balance < amount {
		return fmt.Errorf("账户余额不足")
	}
	account.Balance -= amount
	err = s.putState(ctx, key1, account)
	if err != nil {
		return fmt.Errorf("更新账户余额失败：%v", err)
	}
	// 添加预扣款记录
	withHolding := WithHolding{
		ID:        id,
		AccountID: accountId,
		ListingID: listingID,
		Amount:    amount,
		TimeStamp: timeStamp,
	}
	// 这个也需要存两份，一份主键是 AccountID，一份主键是ListingID
	key2, err := s.getCompositeKey(ctx, WITH_HOLDING_KEY1, []string{fmt.Sprintf("%d", withHolding.AccountID), withHolding.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.putState(ctx, key2, withHolding)
	if err != nil {
		return fmt.Errorf("保存预扣款记录失败：%v", err)
	}
	key3, err := s.getCompositeKey(ctx, WITH_HOLDING_KEY2, []string{withHolding.ListingID, withHolding.ID})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.putState(ctx, key3, withHolding)
	if err != nil {
		return fmt.Errorf("保存预扣款记录失败：%v", err)
	}
	return nil
}

// 查询某个账户的预扣款记录
func (s *SmartContract) GetWithHoldingByAccountID(ctx contractapi.TransactionContextInterface, accountID int) ([]WithHolding, error) {
	var withHoldings []WithHolding
	results, err := ctx.GetStub().GetStateByPartialCompositeKey(WITH_HOLDING_KEY1, []string{fmt.Sprintf("%d", accountID)})
	if err != nil {
		return nil, fmt.Errorf("查询预扣款记录失败：%v", err)
	}
	for results.HasNext() {
		var withHolding WithHolding
		result, err := results.Next()
		if err != nil {
			return nil, fmt.Errorf("查询预扣款记录失败：%v", err)
		}
		err = json.Unmarshal(result.Value, &withHolding)
		if err != nil {
			return nil, fmt.Errorf("解析数据失败：%v", err)
		}
		withHoldings = append(withHoldings, withHolding)
	}
	return withHoldings, nil
}

// 查询某个商品的预扣款记录
func (s *SmartContract) GetWithHoldingByListingID(ctx contractapi.TransactionContextInterface, listingID string) ([]WithHolding, error) {
	var withHoldings []WithHolding
	results, err := ctx.GetStub().GetStateByPartialCompositeKey(WITH_HOLDING_KEY2, []string{listingID})
	if err != nil {
		return nil, fmt.Errorf("查询预扣款记录失败：%v", err)
	}
	for results.HasNext() {
		var withHolding WithHolding
		result, err := results.Next()
		if err != nil {
			return nil, fmt.Errorf("查询预扣款记录失败：%v", err)
		}
		err = json.Unmarshal(result.Value, &withHolding)
		if err != nil {
			return nil, fmt.Errorf("解析数据失败：%v", err)
		}
		withHoldings = append(withHoldings, withHolding)
	}
	return withHoldings, nil
}

// 清除所有预扣款
func (s *SmartContract) ClearWithHolding(ctx contractapi.TransactionContextInterface, listingID string) error {
	// 查询该商品的扣款记录
	withHoldings, err := s.GetWithHoldingByListingID(ctx, listingID)
	if err != nil {
		return fmt.Errorf("查询扣款记录失败：%v", err)
	}
	if len(withHoldings) == 0 {
		return fmt.Errorf("没有相关商品的扣款记录")
	}
	for _, withHolding := range withHoldings {
		var account Account
		key1, err := s.getCompositeKey(ctx, ACCOUNT_KEY, []string{fmt.Sprintf("%d", withHolding.AccountID)})
		if err != nil {
			return fmt.Errorf("创建复合键失败：%v", err)
		}
		err = s.getState(ctx, key1, &account)
		if err != nil {
			return fmt.Errorf("查询余额失败：%v", err)
		}
		account.Balance += withHolding.Amount
		// 移除旧的记录
		err = ctx.GetStub().DelState(key1)
		if err != nil {
			return fmt.Errorf("移除余额记录失败：%v", err)
		}
		// 添加新的记录
		err = s.putState(ctx, key1, account)
		if err != nil {
			return fmt.Errorf("更新余额失败：%v", err)
		}
		key2, err := s.getCompositeKey(ctx, WITH_HOLDING_KEY1, []string{fmt.Sprintf("%d", withHolding.AccountID), withHolding.ID})
		if err != nil {
			return fmt.Errorf("创建复合键失败：%v", err)
		}
		err = ctx.GetStub().DelState(key2)
		if err != nil {
			return fmt.Errorf("删除扣款记录失败：%v", err)
		}
		key3, err := s.getCompositeKey(ctx, WITH_HOLDING_KEY2, []string{withHolding.ListingID, withHolding.ID})
		if err != nil {
			return fmt.Errorf("创建复合键失败：%v", err)
		}
		err = ctx.GetStub().DelState(key3)
		if err != nil {
			return fmt.Errorf("删除扣款记录失败：%v", err)
		}
	}
	return nil
}

// 创建 NFT
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, imageName string,
	name string, authorId int, ownerId int, description string, timeStamp time.Time) (Asset, error) {
	asset := Asset{
		ID:          id,
		ImageName:   imageName,
		Name:        name,
		AuthorId:    authorId,
		OwnerId:     ownerId,
		Description: description,
		TimeStamp:   timeStamp,
	}
	// 这里存三份，一份主键是 ID，一份主键是 AuthorId，一份主键是 OwnerId
	key1, err := s.getCompositeKey(ctx, ASSET_KEY1, []string{id})
	if err != nil {
		return Asset{}, fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.putState(ctx, key1, asset)
	if err != nil {
		return Asset{}, fmt.Errorf("保存 NFT 失败：%v", err)
	}
	key2, err := s.getCompositeKey(ctx, ASSET_KEY2, []string{fmt.Sprintf("%d", authorId), id})
	if err != nil {
		return Asset{}, fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.putState(ctx, key2, asset)
	if err != nil {
		return Asset{}, fmt.Errorf("保存 NFT 失败：%v", err)
	}
	key3, err := s.getCompositeKey(ctx, ASSET_KEY3, []string{fmt.Sprintf("%d", ownerId), id})
	if err != nil {
		return Asset{}, fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.putState(ctx, key3, asset)
	if err != nil {
		return Asset{}, fmt.Errorf("保存 NFT 失败：%v", err)
	}
	return asset, nil
}

// 根据ID查询某个NFT
func (s *SmartContract) GetAssetByID(ctx contractapi.TransactionContextInterface, id string) (Asset, error) {
	var asset Asset
	key, err := s.getCompositeKey(ctx, ASSET_KEY1, []string{id})
	if err != nil {
		return Asset{}, fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.getState(ctx, key, &asset)
	if err != nil {
		return Asset{}, fmt.Errorf("查询 NFT 失败：%v", err)
	}
	return asset, nil
}

// 根据AuthorId查询某个NFT
func (s *SmartContract) GetAssetByAuthorID(ctx contractapi.TransactionContextInterface, authorId int) ([]Asset, error) {
	var assets []Asset
	results, err := ctx.GetStub().GetStateByPartialCompositeKey(ASSET_KEY2, []string{fmt.Sprintf("%d", authorId)})
	if err != nil {
		return nil, fmt.Errorf("查询 NFT 失败：%v", err)
	}
	for results.HasNext() {
		var asset Asset
		result, err := results.Next()
		if err != nil {
			return nil, fmt.Errorf("查询 NFT 失败：%v", err)
		}
		err = json.Unmarshal(result.Value, &asset)
		if err != nil {
			return nil, fmt.Errorf("解析数据失败：%v", err)
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

// 根据OwnerId查询某个NFT
func (s *SmartContract) GetAssetByOwnerID(ctx contractapi.TransactionContextInterface, ownerId int) ([]Asset, error) {
	var assets []Asset
	results, err := ctx.GetStub().GetStateByPartialCompositeKey(ASSET_KEY3, []string{fmt.Sprintf("%d", ownerId)})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败：%v", err)
	}
	for results.HasNext() {
		var asset Asset
		result, err := results.Next()
		if err != nil {
			return nil, fmt.Errorf("查询 NFT 失败：%v", err)
		}
		err = json.Unmarshal(result.Value, &asset)
		if err != nil {
			return nil, fmt.Errorf("解析数据失败：%v", err)
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

// 转移 NFT 的所有权
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwnerId int, userId int, timeStamp time.Time) error {
	var asset Asset
	//三份记录都需要修改
	key1, err := s.getCompositeKey(ctx, ASSET_KEY1, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.getState(ctx, key1, &asset)
	if err != nil {
		return fmt.Errorf("查询 NFT 失败：%v", err)
	}
	// 确保转移请求是所有者发起的
	if asset.AuthorId != userId {
		return fmt.Errorf("只有 NFT 的所有者可以转移所有权")
	}
	if asset.OwnerId == newOwnerId {
		return fmt.Errorf("新旧主人不能相同")
	}
	oldOwnerId := asset.OwnerId
	asset.OwnerId = newOwnerId
	// 删除旧的记录
	err = ctx.GetStub().DelState(key1)
	if err != nil {
		return fmt.Errorf("删除旧的所有权记录失败：%v", err)
	}
	err = s.putState(ctx, key1, asset)
	if err != nil {
		return fmt.Errorf("保存 NFT 失败：%v", err)
	}
	key2, err := s.getCompositeKey(ctx, ASSET_KEY2, []string{fmt.Sprintf("%d", asset.AuthorId), id})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = ctx.GetStub().DelState(key2)
	if err != nil {
		return fmt.Errorf("删除旧的所有权记录失败：%v", err)
	}
	err = s.putState(ctx, key2, asset)
	if err != nil {
		return fmt.Errorf("保存 NFT 失败：%v", err)
	}
	key3, err := s.getCompositeKey(ctx, ASSET_KEY3, []string{fmt.Sprintf("%d", oldOwnerId), id})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = ctx.GetStub().DelState(key3)
	if err != nil {
		return fmt.Errorf("删除旧的所有权记录失败：%v", err)
	}
	// 重新构造键，因为所有者变了
	key3, err = s.getCompositeKey(ctx, ASSET_KEY3, []string{fmt.Sprintf("%d", newOwnerId), id})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}
	err = s.putState(ctx, key3, asset)
	if err != nil {
		return fmt.Errorf("保存 NFT 失败：%v", err)
	}
	return nil
}

// InitLedger 初始化账本
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("InitLedger")
	return nil
}

// 卖家接受出价 -> 释放冻结资金到卖家
func (s *SmartContract) ReleaseHolding(ctx contractapi.TransactionContextInterface, listingID string, sellerID int, amount int, timeStamp time.Time) error {
	// 给卖家加钱
	var seller Account
	sellerKey, _ := s.getCompositeKey(ctx, ACCOUNT_KEY, []string{fmt.Sprintf("%d", sellerID)})
	_ = s.getState(ctx, sellerKey, &seller) // 如果卖家没钱包，应该先开通
	seller.Balance += amount
	if err := s.putState(ctx, sellerKey, seller); err != nil {
		return fmt.Errorf("更新卖家账户失败：%v", err)
	}

	// 删除 listing 下的冻结记录
	withHoldings, err := s.GetWithHoldingByListingID(ctx, listingID)
	if err != nil {
		return err
	}
	for _, w := range withHoldings {
		key1, _ := s.getCompositeKey(ctx, WITH_HOLDING_KEY1, []string{fmt.Sprintf("%d", w.AccountID), w.ID})
		_ = ctx.GetStub().DelState(key1)
		key2, _ := s.getCompositeKey(ctx, WITH_HOLDING_KEY2, []string{w.ListingID, w.ID})
		_ = ctx.GetStub().DelState(key2)
	}
	return nil
}

// 买家退款 -> 把冻结金额退回买家
func (s *SmartContract) RefundHolding(ctx contractapi.TransactionContextInterface, listingID string, bidderID int, amount int, timeStamp time.Time) error {
	// 给买家退钱
	var buyer Account
	buyerKey, _ := s.getCompositeKey(ctx, ACCOUNT_KEY, []string{fmt.Sprintf("%d", bidderID)})
	_ = s.getState(ctx, buyerKey, &buyer)
	buyer.Balance += amount
	if err := s.putState(ctx, buyerKey, buyer); err != nil {
		return fmt.Errorf("更新买家账户失败：%v", err)
	}

	// 删除冻结记录
	withHoldings, err := s.GetWithHoldingByListingID(ctx, listingID)
	if err != nil {
		return err
	}
	for _, w := range withHoldings {
		if w.AccountID == bidderID {
			key1, _ := s.getCompositeKey(ctx, WITH_HOLDING_KEY1, []string{fmt.Sprintf("%d", w.AccountID), w.ID})
			_ = ctx.GetStub().DelState(key1)
			key2, _ := s.getCompositeKey(ctx, WITH_HOLDING_KEY2, []string{w.ListingID, w.ID})
			_ = ctx.GetStub().DelState(key2)
		}
	}
	return nil
}
func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("创建智能合约失败：%v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("启动智能合约失败：%v", err)
	}
}
