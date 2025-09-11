<template>
  <div class="wallet-container">
    <WalletNav />
    <main class="main-content">
      <!-- 余额卡片 -->
      <div class="balance-card">
        <h2>账户余额</h2>
        <p class="balance-amount">{{ balance }} 代币</p>
      </div>

      <!-- 铸造表单 -->
      <div class="transfer-form">
        <h3>代币铸造</h3>
        <a-form @submit.prevent="handleMintToken" class="custom-transfer-form">
          <div class="form-group">
            <label class="form-label" for="accountId">目标账户ID</label>
            <a-input
              v-model:value="mintForm.accountId"
              type="number"
              @blur="validateField('accountId')"
              @input="clearError('accountId')"
              placeholder="请输入要铸造代币的账户ID"
            />
            <p v-if="errors.accountId" class="error-message">{{ errors.accountId }}</p>
          </div>
          
          <div class="form-group">
            <label class="form-label" for="amount">铸造金额</label>
            <a-input
              v-model:value="mintForm.amount"
              type="number"
              @blur="validateField('amount')"
              @input="clearError('amount')"
              placeholder="请输入铸造金额（必须大于0）"
            />
            <p v-if="errors.amount" class="error-message">{{ errors.amount }}</p>
          </div>
              
          <div class="form-group">
            <a-button type="primary" html-type="submit" :disabled="isLoading">
              <span v-if="!isLoading">确认铸造</span>
              <span v-if="isLoading">铸造中...</span>
            </a-button>
          </div>
        </a-form>
        
        <div v-if="message" class="message" :class="{ success: isSuccess, error: !isSuccess }">
          {{ message }}
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import { walletApi } from '../api/index';
import { onMounted, ref } from 'vue';
import { message } from 'ant-design-vue';

export default {
  setup() {
    // 状态定义
    const balance = ref(0);
    const isLoading = ref(false);
    const message = ref('');
    const isSuccess = ref(false);
    
    // 表单数据
    const mintForm = ref({
      accountId: '',
      amount: ''
    });
    
    // 错误信息
    const errors = ref({
      accountId: '',
      amount: ''
    });

    // 加载余额
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

    // 字段验证
    const validateField = (field) => {
      const value = mintForm.value[field];
      const numValue = Number(value);

      if (field === 'accountId') {
        if (!value) {
          errors.value.accountId = '请输入目标账户ID';
        } else if (isNaN(numValue) || numValue < 1) {
          errors.value.accountId = '账户ID必须是大于0的数字';
        } else {
          errors.value.accountId = '';
        }
      }

      if (field === 'amount') {
        if (!value) {
          errors.value.amount = '请输入铸造金额';
        } else if (isNaN(numValue) || numValue < 1) {
          errors.value.amount = '金额必须大于0';
        } else {
          errors.value.amount = '';
        }
      }
    };

    // 清除错误
    const clearError = (field) => {
      errors.value[field] = '';
    };

    // 表单验证
    const validateForm = () => {
      let isValid = true;
      validateField('accountId');
      validateField('amount');
      
      if (errors.value.accountId || errors.value.amount) {
        isValid = false;
      }
      return isValid;
    };

    // 显示消息
    const showMessage = (msg, success) => {
      message.value = msg;
      isSuccess.value = success;
      
      setTimeout(() => {
        message.value = '';
      }, 3000);
    };

    // 处理铸造
    const handleMintToken = async () => {
      if (!validateForm()) {
        return;
      }

      try {
        isLoading.value = true;
        message.value = '';
        
        const accountId = Number(mintForm.value.accountId);
        const amount = Number(mintForm.value.amount);
        
        await walletApi.mintToken(accountId, amount);
        showMessage('代币铸造成功', true);
        
        // 重置表单
        mintForm.value = {
          accountId: '',
          amount: ''
        };
        errors.value = {
          accountId: '',
          amount: ''
        };
        
        // 刷新余额
        loadBalance();
      } catch (error) {
        const errorMsg = error.message || '代币铸造失败，请稍后重试';
        showMessage(errorMsg, false);
        console.error('铸造失败:', error);
      } finally {
        isLoading.value = false;
      }
    };

    // 页面挂载时加载余额
    onMounted(() => {
      loadBalance();
    });

    return {
      balance,
      isLoading,
      message,
      isSuccess,
      mintForm,
      errors,
      validateField,
      clearError,
      handleMintToken
    };
  }
};
</script>

<style scoped>
/* 复用Wallet组件的样式风格 */
.wallet-container {
  min-height: 100vh;
  background-color: #f5f7fa;
  font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
  color: #333;
}

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

/* 表单样式 */
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

/* 消息提示样式 */
.message {
  margin-top: 1rem;
  padding: 0.8rem;
  border-radius: 4px;
  text-align: center;
}

.success {
  background-color: #e6f7ee;
  color: #2a9d54;
  border: 1px solid #b7eb8f;
}

.error {
  background-color: #fff2f0;
  color: #cf1322;
  border: 1px solid #ffccc7;
}
</style>