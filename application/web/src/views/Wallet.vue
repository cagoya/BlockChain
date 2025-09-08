<template>
  <div class="wallet-container">
    <header class="user-info-section">
      <div class="user-profile">
        <button class="back-button" @click="goToDashboard">
          <i class="back-icon">←</i>
        </button>
        <div class="user-details">
          <h1 class="username">{{ user.username }}</h1>
          <p class="greeting">管理你的的加密货币和数字收藏品。</p>
        </div>
      </div>
    </header>

    <!-- 主内容区 -->
    <main class="main-content">
      <!-- 余额卡片（保持不变） -->
      <div class="balance-card">
        <h2>账户余额</h2>
        <p class="balance-amount">{{ balance}} 代币</p>
      </div>

      <!-- 转账表单（修改为自定义验证） -->
      <div class="transfer-form">
        <h3>转账</h3>
        <a-form @submit.prevent="handleTransfer" class="custom-transfer-form">
          <!-- 接收方ID输入 -->
          <div class="form-group">
            <label class="form-label">接收方ID</label>
            <a-input 
              v-model:value="transferForm.recipientId" 
              type="number"
              @blur="validateField('recipientId')"
              @input="clearError('recipientId')"
            />
            <p v-if="errors.recipientId" class="error-message">{{ errors.recipientId }}</p>
          </div>

          <!-- 转账金额输入 -->
          <div class="form-group">
            <label class="form-label">转账金额</label>
            <a-input 
              v-model:value="transferForm.amount" 
              type="number"
              @blur="validateField('amount')"
              @input="clearError('amount')"
            />
            <p v-if="errors.amount" class="error-message">{{ errors.amount }}</p>
          </div>

          <div class="form-group">
            <a-button type="primary" html-type="submit">确认转账</a-button>
          </div>
        </a-form>
      </div>

      <!-- 转账记录（保持不变） -->
      <div class="transfer-records">
        <h3>转账记录</h3>
        <a-tabs default-active-key="sent">
          <a-tab-pane key="sent" tab="转出记录">
            <a-table :data-source="sentTransfers" :columns="transferColumns" row-key="id" />
          </a-tab-pane>
          <a-tab-pane key="received" tab="转入记录">
            <a-table :data-source="receivedTransfers" :columns="transferColumns" row-key="id" />
          </a-tab-pane>
        </a-tabs>
      </div>

      <!-- 我的NFT资产（保持不变） -->
      <div class="my-nfts">
        <h3>我的NFT资产</h3>
        <div class="nft-list" v-if="nfts.length > 0">
          <AssetCard 
            v-for="nft in nfts" 
            :key="nft.id" 
            :asset="nft" 
          />
        </div>
        <p v-else class="no-nft">暂无NFT资产</p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { message } from 'ant-design-vue';
import AssetCard from '../components/AssetCard.vue';
import { walletApi, assetApi } from '../api';

// 初始化路由
const router = useRouter();

// 添加返回方法
const goToDashboard = () => {
  router.push('/dashboard');
};

// 类型定义（保持不变）
interface UserInfo {
  username: string;
  avatarURL: string;
  id: number;
}

interface TransferForm {
  recipientId: string;
  amount: string;
}

interface TransferRecord {
  id: string;
  senderId: number;
  recipientId: number;
  amount: number;
  time: string;
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

// 状态定义
const user = ref<UserInfo>({
  username: '游客',
  avatarURL: 'https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png',
  id: 0
});
const balance = ref<number>(0);
const transferForm = ref<TransferForm>({
  recipientId: '',
  amount: ''
});
// 新增：错误状态管理
const errors = ref({
  recipientId: '',
  amount: ''
});
const sentTransfers = ref<TransferRecord[]>([]);
const receivedTransfers = ref<TransferRecord[]>([]);
const nfts = ref<Asset[]>([]);

// 转账记录表格列定义（保持不变）
const transferColumns = [
  {
    title: '交易ID',
    dataIndex: 'id',
    key: 'id',
    ellipsis: true
  },
  {
    title: '对方ID',
    key: 'counterparty',
    render: (record: TransferRecord) => 
    Number(record.senderId) === Number(user.value.id) 
      ? record.recipientId 
      : record.senderId
  },
  {
    title: '金额',
    dataIndex: 'amount',
    key: 'amount'
  },
  {
    title: '时间',
    dataIndex: 'time',
    key: 'time',
    render: (time: string) => new Date(time).toLocaleString()
  }
];

// 加载数据的方法（保持不变）
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

const loadBalance = async () => {
  try {
    const response = await walletApi.getBalance();
    if (response.data.code === 200) {
      balance.value = response.data.data;
    }
  } catch (error) {
    message.error('获取余额失败');
    console.error(error);
  }
};

const loadTransferRecords = async () => {
  try {
    const sentResponse = await walletApi.getTransfersBySender();
    if (sentResponse.data.code === 200) {
      sentTransfers.value = sentResponse.data.data || [];
    }

    const receivedResponse = await walletApi.getTransfersByRecipient();
    if (receivedResponse.data.code === 200) {
      receivedTransfers.value = receivedResponse.data.data || [];
    }
  } catch (error) {
    message.error('获取转账记录失败');
    console.error(error);
  }
};

const loadMyNFTs = async () => {
  try {
    const response = await assetApi.getByOwnerId(user.value.id.toString());
    if (response.data.code === 200) {
      nfts.value = response.data.data || [];
    }
  } catch (error) {
    message.error('获取NFT资产失败');
    console.error(error);
  }
};

// 自定义验证逻辑
const validateField = (field: keyof TransferForm) => {
  const value = transferForm.value[field];
  const numValue = Number(value);

  if (field === 'recipientId') {
    if (!value) {
      errors.value.recipientId = '请输入接收方ID';
    } else if (isNaN(numValue) || numValue < 1) {
      errors.value.recipientId = '接收方ID必须是大于0的数字';
    } else {
      errors.value.recipientId = '';
    }
  }

  if (field === 'amount') {
    if (!value) {
      errors.value.amount = '请输入转账金额';
    } else if (isNaN(numValue) || numValue < 1) {
      errors.value.amount = '金额必须大于0';
    } else {
      errors.value.amount = '';
    }
  }
};

// 清除单个字段错误
const clearError = (field: keyof TransferForm) => {
  errors.value[field] = '';
};

// 新增：整体表单验证
const validateForm = (): boolean => {
  let isValid = true;
  validateField('recipientId');
  validateField('amount');
  
  if (errors.value.recipientId || errors.value.amount) {
    isValid = false;
  }
  return isValid;
};

// 处理转账逻辑
const handleTransfer = async () => {
  // 使用自定义验证
  if (!validateForm()) {
    return;
  }

  const recipientId = Number(transferForm.value.recipientId);
  const amount = Number(transferForm.value.amount);

  try {
    const response = await walletApi.transfer(recipientId, amount);
    if (response.data.code === 200) {
      message.success('转账成功');
      // 重置表单和错误信息
      transferForm.value = { recipientId: '', amount: '' };
      errors.value = { recipientId: '', amount: '' };
      loadBalance();
      loadTransferRecords();
    }
  } catch (error) {
    message.error('转账失败');
    console.error(error);
  }
};

// 页面挂载时加载数据（保持不变）
onMounted(() => {
  loadUserInfo();
  loadBalance();
  loadTransferRecords();
  loadMyNFTs();
});
</script>

<style scoped>

.back-button {
  position: absolute;
  top: 20px;
  left: 20px;
  background-color: rgba(255, 255, 255, 0.8);
  border: none;
  border-radius: 19px;
  padding: 6px 12px;
  display: flex;
  align-items: center;
  cursor: pointer;
  font-size: 10px;
  color: #4a90e2;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
}

.back-button:hover {
  background-color: white;
  transform: translateX(-2px);
}

.back-icon {
  margin-right: 6px;
  font-style: normal;
}

.custom-transfer-form {
  padding: 10px 0;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
}

.error-message {
  margin: 5px 0 0;
  color: #f5222d;
  font-size: 12px;
  line-height: 1.5;
}

.wallet-container {
  min-height: 100vh;
  background-color: #f5f7fa;
  font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
  color: #333;
}

/* 用户信息头部样式 */
.user-info-section {
  width: 100%;
  padding: 40px 24px;
  background: linear-gradient(135deg, #4a90e2 0%, #76b1f3 100%);
  color: #fff;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.user-profile {
  display: flex;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: 4px solid rgba(255, 255, 255, 0.9);
  object-fit: cover;
  margin-right: 20px;
}

.username {
  font-size: 2rem;
  font-weight: 700;
  margin: 0;
  color: #fff;
  
}

.greeting {
  font-size: 1rem;
  opacity: 0.9;
  margin: 5px 0 0;
}

/* 主内容区样式 */
.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 30px 24px;
}

/* 余额卡片样式 */
.balance-card {
  background-color: #fff;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 30px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.balance-card h2 {
  margin: 0 0 15px 0;
  color: #4a90e2;
}

.balance-amount {
  font-size: 2rem;
  font-weight: bold;
  margin: 0;
  color: #333;
}

/* 转账表单样式 */
.transfer-form {
  background-color: #fff;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 30px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.transfer-form h3 {
  margin: 0 0 20px 0;
  color: #333;
}

/* 转账记录样式 */
.transfer-records {
  background-color: #fff;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 30px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.transfer-records h3 {
  margin: 0 0 20px 0;
  color: #333;
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
}
</style>