<template>
  <div class="asset-card" :class="{ 'selected': isSelected }">
    <div class="asset-image-container">
      <img :src="getImageUrl(asset.imageName)" :alt="asset.name" class="asset-image" />
    </div>
    <div class="asset-info">
      <h3 class="asset-name">{{ asset.name }}</h3>
      <div class="asset-meta">
        <p class="asset-author">创作者ID: {{ asset.authorId }}</p>
        <p class="asset-owner">拥有者ID: {{ asset.ownerId }}</p>
        <p class="asset-created-at">铸造时间: {{ formatDate(asset.timeStamp) }}</p>
      </div>
      <div v-if="asset.description" class="asset-description">
        <p>{{ asset.description }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue';

interface Asset {
  id: string;
  name: string;
  description: string;
  imageName: string;
  authorId: number;
  ownerId: number;
  timeStamp: string;
}

defineProps<{
  asset: Asset;
  isSelected?: boolean;
}>();

// 获取图片完整URL
const getImageUrl = (imageName: string) => {
  return `http://localhost:8888/public/images/${imageName}`;
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
</script>

<style scoped>
.asset-card {
  width: 280px;
  height: 380px;
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
}

.asset-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
  border: 2px solid #2962ff;
}

.asset-card.selected {
  border: 2px solid #2962ff;
  box-shadow: 0 8px 30px rgba(41, 98, 255, 0.3);
  transform: translateY(-5px);
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
  max-height: 40px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
</style>
