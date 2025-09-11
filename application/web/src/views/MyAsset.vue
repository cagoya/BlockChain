<template>
  <div class="my-asset-container">
    <WalletNav />

    <!-- 主内容区 -->
    <main class="main-content">
      <!-- 我的NFT资产 -->
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
import { message } from 'ant-design-vue';
import AssetCard from '../components/AssetCard.vue';
import { assetApi } from '../api';

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

// 状态定义
const user = ref<UserInfo>({
  username: '游客',
  avatarURL: 'https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png',
  id: 0
});
const nfts = ref<Asset[]>([]);

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
    }
  } catch (error) {
    message.error('获取NFT资产失败');
    console.error(error);
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
