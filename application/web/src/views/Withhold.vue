<template>
  <div class="withhold-container">
    <WalletNav />

    <!-- 主内容区 -->
    <main class="main-content">
      <!-- 预扣款记录 -->
      <div class="withhold-records">
        <h3>预扣款记录</h3>
        <a-table :data-source="withholdings" :columns="withholdColumns" row-key="id" />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { walletApi } from '../api';

// 类型定义
interface WithHoldingRecord {
  id: string;
  accountId: number;
  listingId: string;
  amount: number;
  timeStamp: string;
}

// 状态定义
const withholdings = ref<WithHoldingRecord[]>([]);

// 预扣款表格列
const withholdColumns = [
  {
    title: '预扣款ID',
    dataIndex: 'id',
    key: 'id',
    ellipsis: true
  },
  {
    title: '账户ID',
    dataIndex: 'accountId',
    key: 'accountId'
  },
  {
    title: '商品ID',
    dataIndex: 'listingId',
    key: 'listingId',
    ellipsis: true
  },
  {
    title: '金额',
    dataIndex: 'amount',
    key: 'amount'
  },
  {
    title: '时间',
    dataIndex: 'timeStamp',
    key: 'timeStamp',
  }
];

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// 加载预扣款记录的方法
const loadWithholdRecords = async () => {
  try {
    const response = await walletApi.getWithholdingsByAccount();
    if (response.data.code === 200) {
      withholdings.value = response.data.data || [];
      // 格式化日期
      withholdings.value.forEach(record => {
        record.timeStamp = formatDate(record.timeStamp);
      });
    }
  } catch (error) {
    message.error('获取预扣款记录失败');
    console.error(error);
  }
};

// 页面挂载时加载数据
onMounted(() => {
  loadWithholdRecords();
});
</script>

<style scoped>
.withhold-container {
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

/* 预扣款记录样式 */
.withhold-records {
  background-color: #fff;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 30px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.withhold-records h3 {
  margin: 0 0 20px 0;
  color: #333;
}
</style>
