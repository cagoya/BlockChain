<template>
  <div class="wallet-container">
    <!-- 余额信息卡片 -->
    <div class="balance-card">
      <div class="balance-header">
        <h2>我的余额</h2>
        <button @click="refreshBalance" class="refresh-btn">
          <i class="refresh-icon" :class="{ spinning: isRefreshing }"></i>
        </button>
      </div>
      <p class="balance-amount">{{ balance }} ETH</p>
      <p class="last-update">最后更新: {{ lastUpdateTime }}</p>
    </div>

    <!-- 功能按钮区 -->
    <div class="action-buttons">
      <button @click="showTransferDialog = true" class="btn transfer-btn">
        <i class="icon">↗️</i> 转账
      </button>
      <button @click="fetchAndShowTransactionHistory" class="btn history-btn">
        <i class="icon">📜</i> 转账记录
      </button>
      <button
        v-if="isFinanceOrg" 
        @click="showMintDialog = true" 
        class="btn mint-btn"
      >
        <i class="icon">✨</i> 铸币
      </button>
    </div>

    <!-- 我的NFT资产 -->
    <div class="nft-assets">
      <div class="section-header">
        <h2>我的NFT资产</h2>
        <span class="asset-count">{{ myNfts.length }} 件资产</span>
      </div>
      
      <div class="nft-grid">
        <ProductCard 
          v-for="nft in myNfts" 
          :key="nft.id" 
          :product="formatNftToProduct(nft)" 
          class="nft-card"
          @mouseenter="nftHovered = nft.id"
          @mouseleave="nftHovered = null"
        />
      </div>
      
      <p v-if="myNfts.length === 0" class="empty-nft">
        <i class="empty-icon">🖼️</i>
        <span>暂无NFT资产</span>
      </p>
    </div>

    <!-- 转账对话框 -->
    <transition name="modal">
      <div class="modal" v-if="showTransferDialog">
        <div class="modal-content">
          <div class="modal-header">
            <h3>转账</h3>
            <button @click="showTransferDialog = false" class="close-btn">&times;</button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="handleTransfer">
              <div class="form-group">
                <label>接收者ID:</label>
                <input 
                  type="number" 
                  v-model.number="transferForm.recipientId" 
                  required
                  min="1"
                  class="form-input"
                  placeholder="输入接收者ID"
                >
              </div>
              <div class="form-group">
                <label>转账金额:</label>
                <input 
                  type="number" 
                  v-model.number="transferForm.amount" 
                  required
                  min="1"
                  step="0.01"
                  class="form-input"
                  placeholder="输入转账金额"
                >
                <p class="balance-warning" v-if="transferForm.amount > balance">
                  ⚠️ 转账金额超过当前余额
                </p>
              </div>
              <div class="modal-footer">
                <button type="button" @click="showTransferDialog = false" class="cancel-btn">取消</button>
                <button type="submit" class="confirm-btn" :disabled="isTransferring || transferForm.amount > balance">
                  <i class="loading-icon" v-if="isTransferring"></i>
                  确认转账
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </transition>

    <!-- 铸币对话框 -->
    <transition name="modal">
      <div class="modal" v-if="showMintDialog && isFinanceOrg">
        <div class="modal-content">
          <div class="modal-header">
            <h3>铸币</h3>
            <button @click="showMintDialog = false" class="close-btn">&times;</button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="handleMint">
              <div class="form-group">
                <label>目标账户ID:</label>
                <input 
                  type="number" 
                  v-model.number="mintForm.accountId" 
                  required
                  min="1"
                  class="form-input"
                  placeholder="输入目标账户ID"
                >
              </div>
              <div class="form-group">
                <label>铸币金额:</label>
                <input 
                  type="number" 
                  v-model.number="mintForm.amount" 
                  required
                  min="1"
                  step="0.01"
                  class="form-input"
                  placeholder="输入铸币金额"
                >
              </div>
              <div class="modal-footer">
                <button type="button" @click="showMintDialog = false" class="cancel-btn">取消</button>
                <button type="submit" class="confirm-btn" :disabled="isMinting">
                  <i class="loading-icon" v-if="isMinting"></i>
                  确认铸币
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </transition>

    <!-- 转账记录对话框 -->
    <transition name="modal">
      <div class="modal" v-if="showTransactionHistory">
        <div class="modal-content large-modal">
          <div class="modal-header">
            <h3>转账记录</h3>
            <button @click="showTransactionHistory = false" class="close-btn">&times;</button>
          </div>
          <div class="modal-body">
            <div class="tabs">
              <div 
                class="tab-btn" 
                :class="{ active: activeTab === 'sent' }"
                @click="activeTab = 'sent'"
              >
                转出记录
              </div>
              <div 
                class="tab-btn" 
                :class="{ active: activeTab === 'received' }"
                @click="activeTab = 'received'"
              >
                转入记录
              </div>
            </div>
            
            <div class="tab-content" v-if="activeTab === 'sent'">
              <table class="transfer-table">
                <thead>
                  <tr>
                    <th>接收者ID</th>
                    <th>金额(ETH)</th>
                    <th>时间</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="transfer in sentTransfers" :key="transfer.id" class="table-row">
                    <td>{{ transfer.recipientID }}</td>
                    <td class="amount negative">{{ transfer.amount }}</td>
                    <td>{{ formatTime(transfer.timeStamp) }}</td>
                  </tr>
                  <tr v-if="sentTransfers.length === 0">
                    <td colspan="3" class="empty-row">暂无转出记录</td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <div class="tab-content" v-if="activeTab === 'received'">
              <table class="transfer-table">
                <thead>
                  <tr>
                    <th>发送者ID</th>
                    <th>金额(ETH)</th>
                    <th>时间</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="transfer in receivedTransfers" :key="transfer.id" class="table-row">
                    <td>{{ transfer.senderID }}</td>
                    <td class="amount positive">{{ transfer.amount }}</td>
                    <td>{{ formatTime(transfer.timeStamp) }}</td>
                  </tr>
                  <tr v-if="receivedTransfers.length === 0">
                    <td colspan="3" class="empty-row">暂无转入记录</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- 遮罩层 -->
    <transition name="overlay">
      <div class="overlay" v-if="showTransferDialog || showMintDialog || showTransactionHistory" 
        @click="closeAllModals"></div>
    </transition>

    <!-- 通知提示 -->
    <div v-if="notification.show" class="notification" :class="notification.type">
      {{ notification.message }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, watch } from 'vue';
import ProductCard from '../components/ProductCard.vue';
import { walletApi } from '../api';

// 状态定义
const balance = ref(0);
const myNfts = ref<any[]>([]);
const isFinanceOrg = ref(false);
const showTransferDialog = ref(false);
const showMintDialog = ref(false);
const showTransactionHistory = ref(false);
const activeTab = ref('sent');
const sentTransfers = ref<any[]>([]);
const receivedTransfers = ref<any[]>([]);
const isRefreshing = ref(false);
const isTransferring = ref(false);
const isMinting = ref(false);
const nftHovered = ref<number | null>(null);
const lastUpdateTime = ref('刚刚');

// 通知提示
const notification = ref({
  show: false,
  message: '',
  type: 'success' // success, error, info
});

// 表单数据
const transferForm = reactive({
  recipientId: 0,
  amount: 0
});

const mintForm = reactive({
  accountId: 0,
  amount: 0
});

// 显示通知
const showNotification = (message: string, type: string = 'success') => {
  notification.value = { show: true, message, type };
  setTimeout(() => {
    notification.value.show = false;
  }, 3000);
};

// 初始化页面
onMounted(async () => {
  try {
    await Promise.all([
      fetchBalance(),
      fetchMyNfts(),
      checkOrgType()
    ]);
    showNotification('数据加载成功');
  } catch (error) {
    showNotification('数据加载失败', 'error');
  }
});

// 格式化NFT数据为ProductCard需要的格式
const formatNftToProduct = (nft: any) => {
  return {
    image: `/images/${nft.imageName}`,
    name: nft.name,
    author: `用户#${nft.authorId}`,
    created_at: formatTime(nft.timeStamp),
  };
};

// 时间格式化
const formatTime = (timeString: string) => {
  const date = new Date(timeString);
  return date.toLocaleString();
};

// API调用函数
const fetchBalance = async () => {
  try {
    const res = await walletApi.getBalance();
    balance.value = res.data;
    lastUpdateTime.value = new Date().toLocaleTimeString();
  } catch (err) {
    showNotification('获取余额失败', 'error');
    console.error(err);
  }
};

const fetchMyNfts = async () => {
  try {
    const res = await walletApi.getAssetsByOwner();
    myNfts.value = res.data;
  } catch (err) {
    showNotification('获取NFT资产失败', 'error');
    console.error(err);
  }
};

const checkOrgType = async () => {
  try {
    const res = await walletApi.getCurrentOrg();
    isFinanceOrg.value = res.data === 3; // 金融组织编号为3
  } catch (err) {
    console.error('获取组织信息失败', err);
  }
};

const handleTransfer = async () => {
  if (!transferForm.recipientId || !transferForm.amount || transferForm.amount <= 0) {
    showNotification('请输入有效的接收者ID和转账金额', 'error');
    return;
  }

  if (transferForm.amount > balance.value) {
    showNotification('转账金额不能超过当前余额', 'error');
    return;
  }

  try {
    isTransferring.value = true;
    await walletApi.transfer(transferForm.recipientId, transferForm.amount);
    showNotification('转账成功');
    showTransferDialog.value = false;
    await fetchBalance();
    // 重置表单
    transferForm.recipientId = 0;
    transferForm.amount = 0;
  } catch (err) {
    showNotification('转账失败: ' + (err as Error).message, 'error');
    console.error(err);
  } finally {
    isTransferring.value = false;
  }
};

const handleMint = async () => {
  if (!mintForm.accountId || !mintForm.amount || mintForm.amount <= 0) {
    showNotification('请输入有效的账户ID和铸币金额', 'error');
    return;
  }

  try {
    isMinting.value = true;
    await walletApi.mintToken(mintForm.accountId, mintForm.amount);
    showNotification('铸币成功');
    showMintDialog.value = false;
    // 重置表单
    mintForm.accountId = 0;
    mintForm.amount = 0;
  } catch (err) {
    showNotification('铸币失败: ' + (err as Error).message, 'error');
    console.error(err);
  } finally {
    isMinting.value = false;
  }
};

const refreshBalance = async () => {
  isRefreshing.value = true;
  try {
    await fetchBalance();
    showNotification('余额已更新');
  } catch (error) {
    showNotification('余额更新失败', 'error');
  } finally {
    isRefreshing.value = false;
  }
};

// 查看转账记录
const fetchAndShowTransactionHistory = async () => {
  try {
    showNotification('正在加载转账记录...', 'info');
    const [sentRes, receivedRes] = await Promise.all([
      walletApi.getTransfersBySender(),
      walletApi.getTransfersByRecipient()
    ]);
    sentTransfers.value = sentRes.data;
    receivedTransfers.value = receivedRes.data;
    showTransactionHistory.value = true;
  } catch (err) {
    showNotification('获取转账记录失败', 'error');
    console.error(err);
  }
};

// 关闭所有弹窗
const closeAllModals = () => {
  showTransferDialog.value = false;
  showMintDialog.value = false;
  showTransactionHistory.value = false;
};

// 监听转账表单金额变化
watch(() => transferForm.amount, () => {
  // 可以在这里添加金额验证逻辑
});
</script>

<style scoped>
.wallet-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background-color: #f8f9fa;
  min-height: 100vh;
}

/* 通知提示样式 */
.notification {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 12px 20px;
  border-radius: 8px;
  color: white;
  z-index: 1001;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  animation: slideIn 0.3s ease-out forwards, fadeOut 0.3s ease-in forwards 2.7s;
}

.notification.success {
  background-color: #4caf50;
}

.notification.error {
  background-color: #f44336;
}

.notification.info {
  background-color: #2196f3;
}

@keyframes slideIn {
  from { transform: translateX(100%); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
}

@keyframes fadeOut {
  from { opacity: 1; }
  to { opacity: 0; visibility: hidden; }
}

/* 余额卡片样式 */
.balance-card {
  background: linear-gradient(135deg, #2962ff 0%, #42a5f5 100%);
  border-radius: 16px;
  padding: 30px 20px;
  box-shadow: 0 8px 30px rgba(41, 98, 255, 0.25);
  margin-bottom: 30px;
  color: white;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  position: relative;
  overflow: hidden;
}

.balance-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  opacity: 0.5;
}

.balance-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 40px rgba(41, 98, 255, 0.3);
}

.balance-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.balance-card h2 {
  margin: 0;
  font-weight: 500;
  opacity: 0.9;
  position: relative;
}

.balance-amount {
  font-size: 42px;
  font-weight: 700;
  margin: 10px 0 10px;
  letter-spacing: 0.5px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  position: relative;
}

.last-update {
  font-size: 14px;
  opacity: 0.8;
  margin: 0;
  position: relative;
}

.refresh-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  cursor: pointer;
  color: white;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.3s ease, transform 0.3s ease;
  position: relative;
}

.refresh-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: rotate(90deg);
}

.refresh-icon {
  display: inline-block;
  transition: transform 0.5s ease;
}

.refresh-icon.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 功能按钮区 */
.action-buttons {
  display: flex;
  gap: 15px;
  margin-bottom: 30px;
  flex-wrap: wrap;
}

.btn {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  color: white;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  position: relative;
  overflow: hidden;
}

.btn::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
  transition: all 0.5s ease;
}

.btn:hover::after {
  left: 100%;
}

.btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

.btn:active {
  transform: translateY(0);
}

.transfer-btn {
  background: #2962ff;
}

.history-btn {
  background: #4caf50;
}

.mint-btn {
  background: #ff9800;
}

.icon {
  font-size: 18px;
}

/* NFT资产区域 */
.nft-assets {
  margin-top: 30px;
  background: white;
  border-radius: 16px;
  padding: 25px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
}

.section-header h2 {
  margin: 0;
  color: #2d3748;
  font-weight: 600;
  position: relative;
  padding-bottom: 8px;
}

.section-header h2::after {
  content: '';
  position: absolute;
  left: 0;
  bottom: 0;
  width: 40px;
  height: 3px;
  background-color: #2962ff;
  border-radius: 3px;
}

.asset-count {
  color: #718096;
  background: #f7fafc;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 14px;
}

.nft-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 25px;
  margin-top: 20px;
}

.nft-card {
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  border-radius: 12px;
  overflow: hidden;
}

.nft-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 12px 20px rgba(0, 0, 0, 0.1);
}

.empty-nft {
  text-align: center;
  color: #666;
  padding: 60px 0;
  background: #f7fafc;
  border-radius: 12px;
  margin: 20px 0 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
}

.empty-icon {
  font-size: 48px;
  opacity: 0.3;
}

/* 弹窗样式 */
.modal-enter-active, .modal-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.modal-enter-from, .modal-leave-to {
  opacity: 0;
  transform: translate(-50%, -55%);
}

.modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 1000;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  width: 90%;
  max-width: 500px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.modal:focus-within {
  box-shadow: 0 15px 45px rgba(41, 98, 255, 0.25);
}

.large-modal {
  max-width: 800px;
}

.modal-header {
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #2d3748;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #718096;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease, transform 0.2s ease;
}

.close-btn:hover {
  background-color: #f7fafc;
  transform: rotate(90deg);
}

.modal-body {
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #4a5568;
  font-weight: 500;
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-sizing: border-box;
  font-size: 16px;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: #2962ff;
  box-shadow: 0 0 0 3px rgba(41, 98, 255, 0.1);
}

.form-input::placeholder {
  color: #a0aec0;
}

.balance-warning {
  color: #e53e3e;
  font-size: 14px;
  margin: 5px 0 0 0;
  display: flex;
  align-items: center;
  gap: 5px;
  animation: shake 0.5s ease;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-5px); }
  40%, 80% { transform: translateX(5px); }
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.cancel-btn {
  padding: 10px 20px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
}

.cancel-btn:hover {
  background-color: #f7fafc;
  border-color: #cbd5e0;
}

.confirm-btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  background: #2962ff;
  color: white;
  cursor: pointer;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: background-color 0.2s ease, transform 0.2s ease;
}

.confirm-btn:hover:not(:disabled) {
  background: #1e40af;
  transform: translateY(-2px);
}

.confirm-btn:active:not(:disabled) {
  transform: translateY(0);
}

.confirm-btn:disabled {
  background: #94a3b8;
  cursor: not-allowed;
}

/* 遮罩层动画 */
.overlay-enter-active, .overlay-leave-active {
  transition: opacity 0.3s ease;
}

.overlay-enter-from, .overlay-leave-to {
  opacity: 0;
}

.overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
  backdrop-filter: blur(3px);
  transition: backdrop-filter 0.3s ease;
}

.overlay:hover {
  backdrop-filter: blur(5px);
}

/* 标签页样式 */
.tabs {
  display: flex;
  border-bottom: 1px solid #e2e8f0;
  margin-bottom: 24px;
  position: relative;
}

.tab-btn {
  padding: 12px 20px;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  font-weight: 500;
  color: #718096;
  transition: all 0.2s ease;
  position: relative;
  z-index: 1;
}

.tab-btn.active {
  border-bottom-color: #2962ff;
  color: #2962ff;
}

.tab-btn:hover:not(.active) {
  color: #4a5568;
  background-color: #f7fafc;
}

/* 表格样式 */
.transfer-table {
  width: 100%;
  border-collapse: collapse;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
}

.transfer-table th,
.transfer-table td {
  padding: 14px 16px;
  text-align: left;
  border-bottom: 1px solid #f0f0f0;
}

.transfer-table th {
  background-color: #f8fafc;
  font-weight: 600;
  color: #4a5568;
  position: sticky;
  top: 0;
}

.table-row {
  transition: background-color 0.2s ease, transform 0.2s ease;
}

.table-row:hover {
  background-color: #f7fafc;
  transform: translateX(5px);
}

.empty-row {
  text-align: center;
  color: #718096;
  padding: 60px 0;
}

.amount {
  font-weight: 600;
}

.amount.positive {
  color: #4caf50;
}

.amount.negative {
  color: #f44336;
}

/* 加载动画 */
.loading-icon {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.5);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s linear infinite;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .balance-amount {
    font-size: 32px;
  }
  
  .action-buttons {
    justify-content: center;
  }
  
  .nft-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  }
  
  .modal {
    width: 95%;
  }
  
  .large-modal {
    max-height: 80vh;
    overflow-y: auto;
  }
  
  .transfer-table th,
  .transfer-table td {
    padding: 10px 8px;
    font-size: 14px;
  }
}
</style>