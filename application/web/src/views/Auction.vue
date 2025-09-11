<template>
  <MarketNav />
  <div class="auction-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>拍卖市场</h1>
      <p class="page-description">浏览和参与NFT拍卖</p>
    </div>

    <!-- 筛选器 -->
    <div class="filter-container">
      <div class="filter-section">
        <span class="filter-label">拍卖状态:</span>
        <a-radio-group v-model:value="statusFilter" @change="handleFilterChange" class="filter-radio-group">
          <a-radio-button value="all">全部</a-radio-button>
          <a-radio-button value="0">未开始</a-radio-button>
          <a-radio-button value="1">进行中</a-radio-button>
          <a-radio-button value="2">已结束</a-radio-button>
        </a-radio-group>
      </div>
      <div class="filter-stats">
        <span class="stats-text">共 {{ filteredAuctions.length }} 个拍卖品</span>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <a-spin size="large" />
      <p>正在加载拍卖品...</p>
    </div>

    <!-- 拍卖品列表 -->
    <div v-else-if="filteredAuctions.length > 0" class="auctions-container">
      <div class="auctions-grid">
        <div 
          v-for="auction in filteredAuctions" 
          :key="auction.id" 
          class="auction-card"
          @click="handleAuctionClick(auction)"
        >
          <!-- 资产图片 -->
          <div class="auction-image-container">
            <img 
              :src="getImageURL(auction.asset.imageName)" 
              :alt="auction.asset.name" 
              class="auction-image" 
            />
            <!-- 拍卖状态标签 -->
            <div class="status-overlay">
              <a-tag :color="getAuctionStatusColor(auction)" class="status-tag">
                {{ getAuctionStatusText(auction) }}
              </a-tag>
            </div>
          </div>

          <!-- 拍卖信息 -->
          <div class="auction-info">
            <h3 class="auction-title">{{ auction.title }}</h3>
            <p class="asset-name">{{ auction.asset.name }}</p>
            
            <!-- 拍卖品ID -->
            <div class="lot-id-info">
              <span class="lot-id-label">拍卖品ID:</span>
              <span class="lot-id-value">#{{ auction.id }}</span>
            </div>
            
            <!-- 价格信息 -->
            <div class="price-info">
              <div class="price-item">
                <span class="price-label">起拍价:</span>
                <span class="price-value">{{ auction.reservePrice }} 代币</span>
              </div>
              <div v-if="auction.currentPrice" class="price-item">
                <span class="price-label">当前价:</span>
                <span class="price-value current-price">{{ auction.currentPrice }} 代币</span>
              </div>
            </div>

            <!-- 时间信息 -->
            <div class="time-info">
              <div class="time-item">
                <span class="time-label">开始:</span>
                <span class="time-value">{{ formatDateTime(auction.startTime) }}</span>
              </div>
              <div class="time-item">
                <span class="time-label">结束:</span>
                <span class="time-value">{{ formatDateTime(auction.deadline) }}</span>
              </div>
            </div>

            <!-- 倒计时 -->
            <div class="countdown" v-if="getAuctionStatus(auction) === 1">
              <a-countdown 
                :value="new Date(auction.deadline).getTime()" 
                format="HH:mm:ss"
                @finish="handleCountdownFinish"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <a-empty description="暂无拍卖品" />
    </div>

    <!-- 拍卖详情对话框 -->
    <a-modal
      :visible="dialogVisible"
      @update:visible="dialogVisible = $event"
      :title="selectedAuction?.title || '拍卖详情'"
      :width="700"
      :footer="null"
    >
      <div class="auction-detail-dialog" v-if="selectedAuction">
        <!-- 资产信息展示 -->
        <div class="asset-info-section">
          <div class="asset-preview">
            <img :src="getImageURL(selectedAuction.asset.imageName)" :alt="selectedAuction.asset.name" class="preview-image" />
            <div class="preview-info">
              <h4>{{ selectedAuction.asset.name }}</h4>
              <p class="asset-id">资产ID: {{ selectedAuction.asset.id }}</p>
              <p class="asset-description" v-if="selectedAuction.asset.description">
                {{ selectedAuction.asset.description }}
              </p>
            </div>
          </div>
        </div>

        <!-- 拍卖状态显示 -->
        <div class="status-section">
          <h4>拍卖状态</h4>
          <a-tag :color="getAuctionStatusColor(selectedAuction)" class="status-tag">
            {{ getAuctionStatusText(selectedAuction) }}
          </a-tag>
        </div>

        <!-- 拍卖信息 -->
        <div class="auction-info-section">
          <h4>拍卖信息</h4>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">起拍价:</span>
              <span class="value">{{ selectedAuction.reservePrice }} 代币</span>
            </div>
            <div class="info-item" v-if="dialogCurrentPrice">
              <span class="label">当前价格:</span>
              <span class="value current-price">{{ dialogCurrentPrice }} 代币</span>
            </div>
            <div class="info-item" v-if="dialogMyBid">
              <span class="label">我的出价:</span>
              <span class="value my-bid">{{ dialogMyBid }} 代币</span>
            </div>
          </div>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">开始时间:</span>
              <span class="value">{{ formatDateTime(selectedAuction.startTime) }}</span>
            </div>
            <div class="info-item">
              <span class="label">结束时间:</span>
              <span class="value">{{ formatDateTime(selectedAuction.deadline) }}</span>
            </div>
          </div>
        </div>

        <!-- 拍卖未开始 -->
        <div v-if="getAuctionStatus(selectedAuction) === 0" class="not-started-section">
          <a-alert
            message="拍卖尚未开始"
            description="该拍卖还未开始，请等待开始时间"
            type="info"
            show-icon
          />
        </div>

        <!-- 拍卖进行中 - 出价功能 -->
        <div v-else-if="getAuctionStatus(selectedAuction) === 1" class="bidding-section">
          <h4>参与竞拍</h4>
          
          <!-- 如果是出售者，显示提示信息 -->
          <div v-if="isSeller(selectedAuction)" class="seller-notice">
            <a-alert
              message="您是该拍卖品的出售者"
              description="作为出售者，您不能对自己的拍卖品出价"
              type="info"
              show-icon
            />
          </div>
          
          <!-- 如果不是出售者，显示出价表单 -->
          <div v-else>
            <a-form :model="bidForm" :rules="bidRules" ref="bidFormRef" layout="vertical">
              <a-form-item label="出价金额" name="bidPrice">
                <a-input-number 
                  v-model:value="bidForm.bidPrice" 
                  :min="Math.max(selectedAuction.reservePrice, (dialogCurrentPrice || 0) + 1)" 
                  :precision="0"
                  placeholder="请输入出价金额"
                  style="width: 100%"
                />
                <div class="bid-hint">
                  最低出价: {{ Math.max(selectedAuction.reservePrice, (dialogCurrentPrice || 0) + 1) }} 代币
                </div>
              </a-form-item>
            </a-form>
            
            <div class="bid-actions">
              <a-button @click="handleDialogCancel">取消</a-button>
              <a-button 
                type="primary" 
                @click="handleDialogSubmitBid"
                :loading="bidding"
              >
                确认出价
              </a-button>
            </div>
          </div>
        </div>

        <!-- 拍卖已结束 - 显示结果 -->
        <div v-else class="result-section">
          <h4>拍卖结果</h4>
          <div v-if="auctionResult">
            <div v-if="auctionResult.bidPrice === 0" class="no-sale-result">
              <a-alert
                message="拍卖流拍"
                description="该拍卖未能成功售出"
                type="warning"
                show-icon
              />
            </div>
            <div v-else class="success-result">
              <a-alert
                message="拍卖成功"
                description="该拍卖已成功售出"
                type="success"
                show-icon
              />
              <div class="result-info">
                <div class="info-item">
                  <span class="label">成交价:</span>
                  <span class="value final-price">{{ auctionResult.bidPrice }} 代币</span>
                </div>
                <div class="info-item" v-if="auctionResult.bidderId">
                  <span class="label">中标者:</span>
                  <span class="value">{{ auctionResult.bidderId }}</span>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="loading-result">
            <a-spin />
            <span>正在获取拍卖结果...</span>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { auctionApi, assetApi, getImageURL } from '../api/index';
import { message } from 'ant-design-vue';
import MarketNav from '../components/MarketNav.vue';

// 拍卖品接口
interface Auction {
  id: number;
  assetId: string;
  title: string;
  reservePrice: number;
  currentPrice?: number;
  startTime: string;
  deadline: string;
  asset: {
    id: string;
    name: string;
    description: string;
    imageName: string;
    authorId: number;
    ownerId: number;
    timeStamp: string;
  };
}

// 拍卖结果接口
interface AuctionResult {
  lotId: number;
  bidPrice: number;
  bidderId?: number;
}

// 响应式数据
const loading = ref(true);
const auctions = ref<Auction[]>([]);
const dialogVisible = ref(false);
const selectedAuction = ref<Auction | null>(null);
const currentUserId = ref<any>(null);

// 筛选相关数据
const statusFilter = ref('all');
const filteredAuctions = ref<Auction[]>([]);

// 对话框相关数据
const dialogCurrentPrice = ref<number>(0);
const dialogMyBid = ref<number>(0);
const auctionResult = ref<AuctionResult | null>(null);
const bidding = ref(false);

// 出价表单
const bidForm = ref({
  bidPrice: 0
});

const bidFormRef = ref();

// 获取当前用户信息
// 加载用户信息
const loadUserInfo = () => {
  const userInfoString = localStorage.getItem('userInfo');
  if (userInfoString) {
    try {
      const parsedUserInfo = JSON.parse(userInfoString);
      currentUserId.value = parsedUserInfo.id;
    } catch (e) {
      console.error('解析 localStorage 中的 userInfo 失败', e);
    }
  }
};

// 检查是否是出售者
const isSeller = (auction: Auction) => {
  return currentUserId.value && auction.asset.ownerId === currentUserId.value;
};

// 出价表单验证规则
const bidRules = {
  bidPrice: [
    { required: true, message: '请输入出价金额', trigger: 'blur' },
    { 
      validator: (_rule: any, value: number) => {
        if (!selectedAuction.value) return Promise.resolve();
        const minPrice = Math.max(selectedAuction.value.reservePrice, (dialogCurrentPrice.value || 0) + 1);
        if (value < minPrice) {
          return Promise.reject(`出价不能低于 ${minPrice} 代币`);
        }
        return Promise.resolve();
      }, 
      trigger: 'blur' 
    }
  ]
} as any;

// 获取拍卖品列表
const fetchAuctions = async () => {
  try {
    loading.value = true;
    const response = await auctionApi.list();
    if (response.data.code !== 200) {
      message.error(response.data.message);
      return;
    }
    const auctionList = response.data.data || [];
    
    // 为每个拍卖品获取完整的资产信息
    const auctionsWithAssets = await Promise.all(auctionList.map(async (auction: Auction) => {
      try {
        // 获取完整的资产信息
        const assetResponse = await assetApi.getById(auction.assetId);
        if (assetResponse.data.code !== 200) {
          message.error(assetResponse.data.message);
          return;
        }
        auction.asset = assetResponse.data.data;
        
        return auction;
      } catch (error) {
        console.error(`获取拍卖品 ${auction.id} 的资产信息失败:`, error);
        // 如果获取资产信息失败，返回拍卖品但标记为错误状态
        return auction;
      }
    }));
    
    auctions.value = auctionsWithAssets;
    // 获取拍卖品列表后，应用筛选
    filterAuctions();
  } catch (error) {
    console.error('获取拍卖品列表失败:', error);
    message.error('获取拍卖品列表失败');
  } finally {
    loading.value = false;
  }
};

// 辅助函数：判断拍卖状态
// 0 未开始
// 1 进行中
// 2 已结束
const getAuctionStatus = (auction: Auction) => {
  const now = new Date();
  const startTime = new Date(auction.startTime);
  const deadline = new Date(auction.deadline);
  
  if (now < startTime) {
    return 0; // 未开始
  } else if (now >= startTime && now < deadline) {
    return 1; // 进行中
  } else {
    return 2; // 已结束
  }
};

// 获取拍卖状态颜色
const getAuctionStatusColor = (auction: Auction) => {
  const status = getAuctionStatus(auction);
  const colorMap = {
    0: 'blue',   // 未开始
    1: 'orange', // 进行中
    2: 'gray'     // 已结束
  };
  return colorMap[status as keyof typeof colorMap];
};

// 获取拍卖状态文本
const getAuctionStatusText = (auction: Auction) => {
  const status = getAuctionStatus(auction);
  const textMap = {
    0: '未开始',
    1: '进行中',
    2: '已结束'
  };
  return textMap[status as keyof typeof textMap];
};

// 筛选拍卖品
const filterAuctions = () => {
  if (statusFilter.value === 'all') {
    filteredAuctions.value = auctions.value;
  } else {
    const filterStatus = parseInt(statusFilter.value);
    filteredAuctions.value = auctions.value.filter(auction => {
      return getAuctionStatus(auction) === filterStatus;
    });
  }
};

// 处理筛选器变化
const handleFilterChange = () => {
  filterAuctions();
};

// 格式化日期时间
const formatDateTime = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// 处理拍卖品点击
const handleAuctionClick = (auction: Auction) => {
  selectedAuction.value = auction;
  dialogVisible.value = true;
};

// 获取对话框拍卖数据
const fetchDialogAuctionData = async () => {
  if (!selectedAuction.value) return;

  try {
    // 使用拍卖品信息中的当前价格
    dialogCurrentPrice.value = selectedAuction.value.currentPrice || selectedAuction.value.reservePrice;

    // 获取我的出价
    try {
      const myBidResponse = await auctionApi.getBid(selectedAuction.value.id);
      if (myBidResponse.data.code === 200) {
        dialogMyBid.value = myBidResponse.data.data || 0;
      } else {
        dialogMyBid.value = 0;
      }
    } catch (error) {
      // 如果没有出价记录，忽略错误
      dialogMyBid.value = 0;
    }

    // 如果拍卖已结束，获取拍卖结果
    if (getAuctionStatus(selectedAuction.value) === 2) {
      try {
        const resultResponse = await auctionApi.getResult(selectedAuction.value.id);
        if (resultResponse.data.code === 200) {
          auctionResult.value = resultResponse.data.data;
        } else {
          console.warn('获取拍卖结果失败:', resultResponse.data.message);
        }
      } catch (error) {
        console.error('获取拍卖结果失败:', error);
      }
    }
  } catch (error) {
    console.error('获取拍卖数据失败:', error);
  }
};

// 处理对话框取消
const handleDialogCancel = () => {
  dialogVisible.value = false;
};

// 处理对话框出价提交
const handleDialogSubmitBid = async () => {
  if (!selectedAuction.value) return;

  try {
    await bidFormRef.value?.validate();
    
    bidding.value = true;
    const response = await auctionApi.bid({
      id: selectedAuction.value.id,
      bidPrice: bidForm.value.bidPrice
    });
    
    if (response.data.code === 200) {
      message.success('出价成功');
      
      // 刷新拍卖品列表
      await fetchAuctions();
      
      // 关闭对话框
      dialogVisible.value = false;
    } else {
      message.error(response.data.message || '出价失败');
    }
  } catch (error) {
    console.error('出价失败:', error);
    message.error('出价失败，请重试');
  } finally {
    bidding.value = false;
  }
};

// 处理倒计时结束
const handleCountdownFinish = () => {
  // 刷新拍卖品列表
  fetchAuctions();
};

// 监听对话框显示状态
watch(() => dialogVisible.value, (newVisible) => {
  if (newVisible && selectedAuction.value) {
    // 重置数据
    dialogCurrentPrice.value = 0;
    dialogMyBid.value = 0;
    auctionResult.value = null;
    bidForm.value.bidPrice = 0;
    
    // 获取拍卖数据
    fetchDialogAuctionData();
  }
});

// 组件挂载时获取数据
onMounted(async () => {
  await loadUserInfo();
  await fetchAuctions();
});
</script>

<style scoped>
.auction-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
}

.page-header {
  text-align: center;
  margin-bottom: 32px;
  color: white;
}

.page-header h1 {
  font-size: 2.5rem;
  font-weight: bold;
  margin: 0 0 8px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.page-description {
  font-size: 1.1rem;
  opacity: 0.9;
  margin: 0;
}

.filter-container {
  max-width: 1400px;
  margin: 0 auto 32px auto;
  padding: 0 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  padding: 20px 24px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
}

.filter-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.filter-label {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.filter-radio-group {
  display: flex;
  gap: 8px;
}

.filter-radio-group :deep(.ant-radio-button-wrapper) {
  border-radius: 8px;
  border: 2px solid #e9ecef;
  background: white;
  color: #666;
  font-weight: 500;
  transition: all 0.3s ease;
}

.filter-radio-group :deep(.ant-radio-button-wrapper:hover) {
  border-color: #2962ff;
  color: #2962ff;
}

.filter-radio-group :deep(.ant-radio-button-wrapper-checked) {
  background: #2962ff;
  border-color: #2962ff;
  color: white;
}

.filter-radio-group :deep(.ant-radio-button-wrapper-checked:hover) {
  background: #1e53e5;
  border-color: #1e53e5;
  color: white;
}

.filter-stats {
  display: flex;
  align-items: center;
}

.stats-text {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  color: white;
}

.loading-container p {
  margin-top: 16px;
  font-size: 16px;
}

.auctions-container {
  max-width: 1400px;
  margin: 0 auto;
}

.auctions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
  padding: 0 16px;
}

.auction-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
}

.auction-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.15);
}

.auction-image-container {
  position: relative;
  height: 200px;
  overflow: hidden;
}

.auction-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.auction-card:hover .auction-image {
  transform: scale(1.05);
}

.status-overlay {
  position: absolute;
  top: 12px;
  right: 12px;
}

.status-tag {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 12px;
  font-weight: 500;
}

.auction-info {
  padding: 20px;
}

.auction-title {
  font-size: 18px;
  font-weight: bold;
  color: #1a237e;
  margin: 0 0 8px 0;
  line-height: 1.3;
}

.asset-name {
  font-size: 14px;
  color: #666;
  margin: 0 0 12px 0;
  line-height: 1.4;
}

.lot-id-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 8px 12px;
  background-color: #f0f2ff;
  border-radius: 6px;
  border-left: 3px solid #2962ff;
}

.lot-id-label {
  font-size: 12px;
  color: #666;
  font-weight: 500;
}

.lot-id-value {
  font-size: 14px;
  color: #2962ff;
  font-weight: bold;
  font-family: 'Courier New', monospace;
}

.price-info {
  margin-bottom: 16px;
}

.price-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.price-label {
  font-size: 14px;
  color: #666;
}

.price-value {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.current-price {
  color: #4caf50 !important;
  font-size: 16px !important;
}

.time-info {
  margin-bottom: 16px;
  padding: 12px;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.time-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.time-item:last-child {
  margin-bottom: 0;
}

.time-label {
  font-size: 12px;
  color: #666;
}

.time-value {
  font-size: 12px;
  color: #333;
  font-weight: 500;
}

.countdown {
  text-align: center;
  padding: 12px;
  background: linear-gradient(45deg, #ff6b6b, #ffa500);
  border-radius: 8px;
  color: white;
}

.countdown :deep(.ant-statistic-content) {
  color: white;
  font-size: 18px;
  font-weight: bold;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  background: white;
  border-radius: 16px;
  margin: 0 16px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .auction-page {
    padding: 16px;
  }
  
  .page-header h1 {
    font-size: 2rem;
  }
  
  .filter-container {
    flex-direction: column;
    gap: 16px;
    padding: 16px;
    margin-bottom: 24px;
  }
  
  .filter-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
    width: 100%;
  }
  
  .filter-radio-group {
    width: 100%;
    justify-content: space-between;
  }
  
  .filter-radio-group :deep(.ant-radio-button-wrapper) {
    flex: 1;
    text-align: center;
  }
  
  .filter-stats {
    width: 100%;
    justify-content: center;
  }
  
  .auctions-grid {
    grid-template-columns: 1fr;
    gap: 16px;
    padding: 0;
  }
  
  .auction-card {
    margin: 0 8px;
  }
}

@media (max-width: 480px) {
  .auction-info {
    padding: 16px;
  }
  
  .auction-title {
    font-size: 16px;
  }
}

/* 对话框样式 */
.auction-detail-dialog {
  padding: 20px 0;
}

.asset-info-section {
  margin-bottom: 24px;
}

.asset-preview {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 16px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.preview-image {
  width: 120px;
  height: 120px;
  object-fit: cover;
  border-radius: 8px;
}

.preview-info h4 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 18px;
}

.preview-info .asset-id {
  margin: 0 0 8px 0;
  color: #666;
  font-size: 14px;
}

.preview-info .asset-description {
  margin: 0;
  color: #666;
  font-size: 14px;
  line-height: 1.5;
}

.status-section {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.status-section h4 {
  margin: 0;
  color: #333;
  font-size: 16px;
}

.auction-info-section {
  margin-bottom: 24px;
}

.auction-info-section h4 {
  margin: 0 0 16px 0;
  color: #333;
  font-size: 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background-color: #f8f9fa;
  border-radius: 6px;
}

.info-item .label {
  color: #666;
  font-size: 14px;
}

.info-item .value {
  color: #333;
  font-size: 14px;
  font-weight: 500;
}

.current-price {
  color: #4caf50 !important;
}

.my-bid {
  color: #2196f3 !important;
}

.final-price {
  color: #ff6b35 !important;
  font-size: 16px !important;
}

.not-started-section,
.result-section {
  margin-bottom: 24px;
}

.bidding-section {
  margin-bottom: 24px;
}

.bidding-section h4 {
  margin: 0 0 16px 0;
  color: #333;
  font-size: 16px;
}

.seller-notice {
  margin-bottom: 20px;
}

.bid-hint {
  margin-top: 8px;
  color: #666;
  font-size: 12px;
}

.bid-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}

.result-info {
  margin-top: 16px;
}

.loading-result {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px;
  justify-content: center;
}
</style>