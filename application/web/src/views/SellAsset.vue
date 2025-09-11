<template>
    <div class="my-asset-container">
      <MarketNav />
  
      <!-- 主内容区 -->
      <main class="main-content">
        <!-- 我的NFT资产 -->
        <div class="my-nfts">
          <div class="header-section">
            <h3>我的NFT资产</h3>
            <!-- 状态筛选 -->
            <div class="filter-section">
              <a-select
                v-model:value="statusFilter"
                placeholder="筛选状态"
                style="width: 150px"
                allow-clear
              >
                <a-select-option value="">全部状态</a-select-option>
                <a-select-option value="not_trading">不在交易中</a-select-option>
                <a-select-option value="listing">普通交易中</a-select-option>
                <a-select-option value="auction">拍卖中</a-select-option>
              </a-select>
            </div>
          </div>
          <div class="nft-list" v-if="filteredNfts.length > 0">
            <AssetCard 
              v-for="nft in filteredNfts" 
              :key="nft.id" 
              :asset="nft"
              :status="assetStatusMap[nft.id] || 'not_trading'"
              @click="handleAssetClick"
            />
          </div>
          <p v-else class="no-nft">
            {{ statusFilter ? '没有符合条件的NFT资产' : '暂无NFT资产' }}
          </p>
        </div>
      </main>

      <!-- 资产状态对话框 -->
      <AssetStatusDialog
        v-model:visible="dialogVisible"
        :asset="selectedAsset"
        :status="assetStatus"
        :submitting="submitting"
        @submit="handleDialogSubmit"
      />
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted, computed } from 'vue';
  import { message } from 'ant-design-vue';
  import AssetCard from '../components/AssetCard.vue';
  import AssetStatusDialog from '../components/AssetStatusDialog.vue';
  import { assetApi, auctionApi, marketApi } from '../api';
  
  // 类型定义
  interface UserInfo {
    username: string;
    avatarURL: string;
    id: number;
  }
  
  interface Asset {
    id: string;
    name: string;
    description: string;
    imageName: string;
    authorId: number;
    ownerId: number;
    timeStamp: string;
  }

  // 资产状态类型
  type AssetStatus = 'not_trading' | 'listing' | 'auction';
  
  // 状态映射
  const statusMap: { [key: number]: AssetStatus } = {
    0: 'not_trading',
    1: 'listing', 
    2: 'auction'
  };
  
  // 状态定义
  const user = ref<UserInfo>({
    username: '',
    avatarURL: '',
    id: 0
  });
  const nfts = ref<Asset[]>([]);
  
  // 资产状态映射
  const assetStatusMap = ref<{ [key: string]: AssetStatus }>({});
  
  // 对话框相关状态
  const dialogVisible = ref(false);
  const selectedAsset = ref<Asset | null>(null);
  const assetStatus = ref<AssetStatus>('not_trading');
  const submitting = ref(false);

  // 筛选状态
  const statusFilter = ref<string>('');

  // 计算属性：根据筛选条件过滤资产
  const filteredNfts = computed(() => {
    if (!statusFilter.value) {
      return nfts.value;
    }
    return nfts.value.filter(nft => {
      const nftStatus = assetStatusMap.value[nft.id] || 'not_trading';
      return nftStatus === statusFilter.value;
    });
  });
  
  // 加载用户信息
  const loadUserInfo = () => {
    const userInfoString = localStorage.getItem('userInfo');
    if (userInfoString) {
      try {
        const parsedUserInfo = JSON.parse(userInfoString);
        user.value = {
          username: parsedUserInfo.username || user.value.username,
          avatarURL: parsedUserInfo.avatarURL || user.value.avatarURL,
          id: parsedUserInfo.id || 0
        };
      } catch (e) {
        console.error('解析用户信息失败', e);
      }
    }
  };
  
  // 加载我的NFT资产
  const loadMyNFTs = async () => {
    try {
      const response = await assetApi.getByOwnerId(user.value.id.toString());
      if (response.data.code === 200) {
        nfts.value = response.data.data || [];
        
        // 为每个资产获取状态
        for (const nft of nfts.value) {
          await loadAssetStatus(nft.id);
        }
      }
    } catch (error) {
      message.error('获取NFT资产失败');
      console.error(error);
    }
  };

  // 加载单个资产状态
  const loadAssetStatus = async (assetId: string) => {
    try {
      const response = await assetApi.getStatus(assetId);
      if (response.data.code === 200) {
        const statusCode = response.data.data;
        assetStatusMap.value[assetId] = statusMap[statusCode] || 'not_trading';
      } else {
        assetStatusMap.value[assetId] = 'not_trading';
      }
    } catch (error) {
      console.error(`获取资产 ${assetId} 状态失败:`, error);
      assetStatusMap.value[assetId] = 'not_trading';
    }
  };

  // 处理资产点击事件
  const handleAssetClick = async (asset: Asset) => {
    selectedAsset.value = asset;
    dialogVisible.value = true;
    
    // 使用已存储的状态
    assetStatus.value = assetStatusMap.value[asset.id] || 'not_trading';
  };

  // 处理对话框提交
  const handleDialogSubmit = async (data: { action: string; formData: any }) => {
    submitting.value = true;

    try {
      if (data.action === 'listing') {
        await marketApi.createListing(data.formData);
        message.success('普通出售创建成功');
      } else if (data.action === 'auction') {
        await auctionApi.create(data.formData);
        message.success('拍卖创建成功');
      }

      dialogVisible.value = false;
      // 重新加载资产列表
      await loadMyNFTs();
    } catch (error) {
      console.error('提交失败:', error);
      message.error('操作失败，请重试');
    } finally {
      submitting.value = false;
    }
  };
  
  // 页面挂载时加载数据
  onMounted(() => {
    loadUserInfo();
    loadMyNFTs();
  });
  </script>
  
  <style scoped>
  .my-asset-container {
    min-height: 100vh;
    background-color: #f5f7fa;
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: #333;
  }
  
  /* 主内容区样式 */
  .main-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 30px 24px;
  }
  
  /* NFT列表样式 */
  .my-nfts {
    background-color: #fff;
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  }
  
  .my-nfts h3 {
    margin: 0 0 20px 0;
    color: #333;
  }

  .header-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .header-section h3 {
    margin: 0;
    color: #333;
  }

  .filter-section {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  
  .nft-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }
  
  .no-nft {
    text-align: center;
    padding: 40px 0;
    color: #888;
  }
  
  /* 响应式调整 */
  @media (max-width: 768px) {
    .nft-list {
      grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    }

    .header-section {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .filter-section {
      width: 100%;
    }
  }

  /* 对话框样式 */
  .asset-status-dialog {
    padding: 20px 0;
  }

  .asset-info-section {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .asset-preview {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    background-color: #f8f9fa;
    border-radius: 8px;
    border: 1px solid #e9ecef;
  }

  .preview-image {
    width: 80px;
    height: 80px;
    object-fit: cover;
    border-radius: 8px;
  }

  .preview-info h4 {
    margin: 0 0 8px 0;
    color: #333;
    font-size: 16px;
  }

  .preview-info .asset-id {
    margin: 0;
    color: #666;
    font-size: 12px;
  }

  .status-section {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .status-section h4 {
    margin: 0;
    color: #333;
    font-size: 14px;
  }

  .status-tag {
    font-size: 12px;
    padding: 4px 8px;
  }

  .action-section {
    padding: 16px 0;
    border-top: 1px solid #e9ecef;
  }

  .action-section h4 {
    margin: 0 0 12px 0;
    color: #333;
    font-size: 14px;
  }

  .action-radio {
    width: 100%;
  }

  .form-section {
    padding: 16px 0;
    border-top: 1px solid #e9ecef;
  }

  .form-section h4 {
    margin: 0 0 16px 0;
    color: #333;
    font-size: 14px;
  }

  .dialog-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding-top: 20px;
    border-top: 1px solid #e9ecef;
  }
  </style>
  