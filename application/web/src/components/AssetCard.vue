<template>
  <div class="asset-card" @click="handleClick">
    <div class="asset-image-container">
      <img :src="getImageURL(asset.imageName)" :alt="asset.name" class="asset-image" />
    </div>
    <div class="asset-info">
      <h3 class="asset-name">{{ asset.name }}</h3>
      <div class="asset-meta">
        <p class="asset-id">ID: {{ asset.id }}</p>
        <p class="asset-author">创作者ID: {{ asset.authorId }}</p>
        <p class="asset-owner">拥有者ID: {{ asset.ownerId }}</p>
        <p class="asset-created-at">铸造时间: {{ formatDate(asset.timeStamp) }}</p>
      </div>
      <!-- 状态标签 -->
      <div class="asset-status">
        <a-tag :color="getStatusColor(status)" class="status-tag">
          {{ getStatusText(status) }}
        </a-tag>
      </div>
      <div v-if="asset.description" class="asset-description">
        <p>{{ asset.description }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
import { getImageURL } from '../api/index';

interface Asset {
  id: string;
  name: string;
  description: string;
  imageName: string;
  authorId: number;
  ownerId: number;
  timeStamp: string;
}

type AssetStatus = 'not_trading' | 'listing' | 'auction';

const props = defineProps<{
  asset: Asset;
  status: AssetStatus;
}>();

const emit = defineEmits<{
  click: [asset: Asset];
}>();

// 处理点击事件
const handleClick = () => {
  emit('click', props.asset);
};

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

// 获取状态颜色
const getStatusColor = (status: AssetStatus) => {
  const colorMap = {
    'not_trading': 'green',
    'listing': 'blue',
    'auction': 'orange'
  };
  return colorMap[status];
};

// 获取状态文本
const getStatusText = (status: AssetStatus) => {
  const textMap = {
    'not_trading': '不在交易中',
    'listing': '普通交易中',
    'auction': '拍卖中'
  };
  return textMap[status];
};
</script>

<style scoped>
.asset-card {
  width: 280px;
  height: 420px;
  background-color: #ffffff;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  text-align: center;
  transition: transform 0.3s ease, box-shadow 0.3s ease, border 0.3s ease;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  padding: 15px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  border: 2px solid transparent;
  cursor: pointer;
}

.asset-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
  border: 2px solid #2962ff;
}

.asset-image-container {
  width: 100%;
  padding: 10px 0;
  height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.asset-image {
  width: 100%;
  max-width: 200px;
  height: 160px;
  object-fit: cover;
  border-radius: 8px;
}

.asset-info {
  flex-grow: 1;
  width: 100%;
  padding: 10px 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  text-align: center;
}

.asset-name {
  font-size: 18px;
  font-weight: bold;
  color: #1a237e;
  margin: 0 0 8px 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.asset-meta {
  margin-bottom: 8px;
  color: #444;
  font-size: 12px;
}

.asset-status {
  margin-bottom: 8px;
  display: flex;
  justify-content: center;
}

.status-tag {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
}

.asset-id {
  font-size: 10px;
  color: #3aabe4;
  line-height: 1.4;
  margin: 2px 0;
}

.asset-author,
.asset-owner,
.asset-created-at {
  font-size: 12px;
  color: #607d8b;
  line-height: 1.4;
  margin: 2px 0;
}

.asset-description {
  margin: 8px 0;
  font-size: 12px;
  color: #666;
  max-height: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
}
</style>