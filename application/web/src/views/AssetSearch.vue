<template>
  <div class="page-container">
    <AssetNav />
    <div class="asset-search-container">
      <div class="search-header">
        <h1>NFT 资产查询</h1>
        <p>搜索和浏览区块链上的数字资产</p>
      </div>

      <div class="search-form">
        <a-form layout="inline">
          <a-form-item label="查询类型">
            <a-select
              v-model:value="searchType"
              style="width: 150px"
              placeholder="选择查询类型"
            >
              <a-select-option value="author">按作者ID</a-select-option>
              <a-select-option value="owner">按拥有者ID</a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item label="查询值">
            <a-input
              v-model:value="searchValue"
              placeholder="请输入查询值"
              style="width: 200px"
            />
          </a-form-item>

          <a-form-item>
            <a-button type="primary" html-type="submit" :loading="searching" @click="handleSearch">
              搜索
            </a-button>
            <a-button @click="handleReset">
              重置
            </a-button>
          </a-form-item>
        </a-form>
      </div>

      <div class="search-results">
        <div v-if="searching" class="loading-container">
          <a-spin size="large" />
          <p>搜索中...</p>
        </div>

        <div v-else-if="assets.length === 0 && hasSearched" class="no-results">
          <a-empty description="未找到相关资产" />
        </div>

        <div v-else-if="assets.length > 0" class="results-grid">
          <div class="results-header">
            <h3>搜索结果 ({{ assets.length }} 个资产)</h3>
          </div>
          <div class="assets-grid">
            <AssetCard
              v-for="asset in assets"
              :key="asset.id"
              :asset="asset"
            />
          </div>
        </div>

        <div v-else class="welcome-message">
          <a-icon type="search" class="welcome-icon" />
          <h3>开始搜索</h3>
          <p>选择查询类型并输入查询值来搜索NFT资产</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AssetNav from '../components/AssetNav.vue';
import AssetCard from '../components/AssetCard.vue';
import { ref } from 'vue';
import { message } from 'ant-design-vue';
import { assetApi } from '../api';

interface Asset {
  id: string;
  name: string;
  description: string;
  imageName: string;
  authorId: number;
  ownerId: number;
  timeStamp: string;
}

const searchType = ref('author');
const searchValue = ref('');
const searching = ref(false);
const hasSearched = ref(false);
const assets = ref<Asset[]>([]);

// 搜索处理
const handleSearch = async () => {
  if (!searchValue.value.trim()) {
    message.warning('请输入查询值');
    return;
  }

  searching.value = true;
  hasSearched.value = true;

  try {
    let response;
    
    switch (searchType.value) {
      case 'author':
        response = await assetApi.getByAuthorId(searchValue.value);
        break;
      case 'owner':
        response = await assetApi.getByOwnerId(searchValue.value);
        break;
      default:
        throw new Error('无效的查询类型');
    }

    const result = await response.data;

    if (result.code === 200) {
      assets.value = result.data || [];
      
      if (assets.value.length === 0) {
        message.info('未找到相关资产');
      } else {
        message.success(`找到 ${assets.value.length} 个资产`);
      }
    } else {
      message.error(result.message || '查询失败');
      assets.value = [];
    }
  } catch (error) {
    console.error('查询失败:', error);
    message.error('查询失败，请重试');
    assets.value = [];
  } finally {
    searching.value = false;
  }
};

// 重置搜索
const handleReset = () => {
  searchType.value = 'author';
  searchValue.value = '';
  assets.value = [];
  hasSearched.value = false;
};


</script>

<style scoped>
/* 页面容器 */
.page-container {
  min-height: 100vh;
  overflow-x: hidden;
  overflow-y: auto;
}

/* 基础容器 */
.asset-search-container {
  max-width: 1200px;
  margin: 40px auto;
  padding: 40px;
  background: linear-gradient(135deg, #ffffff, #f0f8ff);
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(41, 98, 255, 0.15);
  border: 1px solid #e0eaff;
}

/* 头部样式 */
.search-header {
  text-align: center;
  margin-bottom: 50px;
}

.search-header h1 {
  color: #1a237e;
  font-size: 36px;
  margin-bottom: 15px;
  font-weight: 700;
  letter-spacing: 1px;
}

.search-header p {
  color: #5a667b;
  font-size: 18px;
  line-height: 1.6;
}

/* 搜索表单样式 */
.search-form {
  background: #f7f9fc;
  padding: 30px;
  border-radius: 12px;
  margin-bottom: 40px;
  border: 1px solid #e0eaff;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.search-form :deep(.ant-form) {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: nowrap;
}

.search-form :deep(.ant-form-item) {
  margin-bottom: 0;
  margin-right: 0;
  flex-shrink: 0;
}

.search-form :deep(.ant-form-item-label > label) {
  font-size: 16px;
  font-weight: 600;
  color: #3f51b5;
  margin-bottom: 0;
  white-space: nowrap;
}

.search-form :deep(.ant-select) {
  border-radius: 8px;
  border: 1px solid #c5cae9;
  transition: all 0.3s ease;
}

.search-form :deep(.ant-select:hover),
.search-form :deep(.ant-select-focused) {
  border-color: #2962ff;
  box-shadow: 0 0 0 2px rgba(41, 98, 255, 0.2);
}

.search-form :deep(.ant-input) {
  border-radius: 8px;
  border: 1px solid #c5cae9;
  padding: 10px 15px;
  font-size: 16px;
  transition: all 0.3s ease;
}

.search-form :deep(.ant-input:hover),
.search-form :deep(.ant-input:focus) {
  border-color: #2962ff;
  box-shadow: 0 0 0 2px rgba(41, 98, 255, 0.2);
}

.search-form :deep(.ant-btn-primary) {
  background-color: #2962ff;
  border-color: #2962ff;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
  margin-right: 8px;
}

.search-form :deep(.ant-btn-primary:hover) {
  background-color: #004acb;
  border-color: #004acb;
  transform: translateY(-1px);
}

.search-form :deep(.ant-btn:not(.ant-btn-primary)) {
  background-color: #ef5350;
  border-color: #ef5350;
  color: #fff;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.search-form :deep(.ant-btn:not(.ant-btn-primary):hover) {
  background-color: #d32f2f;
  border-color: #d32f2f;
  transform: translateY(-1px);
}

/* 搜索结果样式 */
.search-results {
  min-height: 400px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  background: #f7f9fc;
  border-radius: 12px;
  border: 1px solid #e0eaff;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.loading-container p {
  margin-top: 16px;
  color: #3f51b5;
  font-size: 16px;
  font-weight: 500;
}

.no-results {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
  background: #f7f9fc;
  border-radius: 12px;
  border: 1px solid #e0eaff;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.results-header {
  margin-bottom: 20px;
  background: #f7f9fc;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #e0eaff;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.results-header h3 {
  color: #1a237e;
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  text-align: center;
}

.assets-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.welcome-message {
  text-align: center;
  padding: 60px 20px;
  background: #f7f9fc;
  border-radius: 12px;
  border: 1px solid #e0eaff;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.welcome-icon {
  font-size: 48px;
  color: #42a5f5;
  margin-bottom: 16px;
}

.welcome-message h3 {
  color: #1a237e;
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 8px;
}

.welcome-message p {
  color: #5a667b;
  font-size: 16px;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .asset-search-container {
    padding: 20px;
    margin: 20px;
  }
  
  .search-header h1 {
    font-size: 28px;
  }
  
  .search-form {
    padding: 20px;
  }
  
  .search-form :deep(.ant-form) {
    flex-wrap: nowrap;
    gap: 8px;
  }
  
  .search-form :deep(.ant-form-item-label > label) {
    font-size: 14px;
  }
  
  .search-form :deep(.ant-select) {
    width: 120px !important;
  }
  
  .search-form :deep(.ant-input) {
    width: 150px !important;
  }
  
  .search-form :deep(.ant-btn) {
    padding: 4px 8px;
    font-size: 12px;
  }
  
  .assets-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }
}
</style>