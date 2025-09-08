<template>
  <div class="dashboard-container">
    <header class="user-info-section">
      <div class="user-profile">
        <a-popover v-model:open="popoverVisible" trigger="click" placement="bottomRight">
          <template #content>
            <div class="popover-menu">
              <a-menu @click="handleMenuClick" :selectable="false" class="custom-popover-menu">
                <a-menu-item key="change_avatar">
                  <template #icon><UserOutlined /></template>
                  更改头像
                </a-menu-item>
                <a-menu-item key="edit_profile">
                  <template #icon><SettingOutlined /></template>
                  编辑信息
                </a-menu-item>
                <a-menu-item v-if="isPlatformAdmin" key="update_org">
                  <template #icon><TeamOutlined /></template>
                  更新组织
                </a-menu-item>
                <a-menu-item key="logout">
                  <template #icon><LogoutOutlined /></template>
                  登出
                </a-menu-item>
              </a-menu>
            </div>
          </template>
          <img :src="user.avatarURL" alt="用户头像" class="avatar clickable-avatar" />
        </a-popover>

        <div class="user-details">
          <h1 class="username">{{ user.username }}</h1>
          <p class="greeting">欢迎回来，探索你的数字资产世界！</p>
        </div>
      </div>
    </header>

    <main class="main-content">
      <div class="card-list">
        <router-link to="/asset/upload" class="card card-creation">
          <div class="card-icon">
            <EditOutlined />
          </div>
          <div class="card-title">NFT 创作</div>
          <div class="card-description">将你的创意转化为独一无二的数字艺术品。</div>
        </router-link>

        <router-link to="/market" class="card card-trading">
          <div class="card-icon">
            <TransactionOutlined />
          </div>
          <div class="card-title">NFT 交易</div>
          <div class="card-description">探索、买卖和收藏来自全球的数字资产。</div>
        </router-link>

        <router-link to="/wallet" class="card card-wallet">
          <div class="card-icon">
            <WalletOutlined />
          </div>
          <div class="card-title">我的钱包</div>
          <div class="card-description">管理你的加密货币和数字收藏品。</div>
        </router-link>
      </div>
    </main>


    <!-- 头像上传组件 -->
    <AvatarUploader 
      v-model:visible="changeAvatarModalVisible" 
      @success="handleAvatarUploadSuccess"
    />

    <!-- 个人信息编辑组件 -->
    <ProfileEditor 
      v-model:visible="editProfileModalVisible" 
      :user-info="userInfo"
      @success="handleProfileUpdateSuccess"
    />

    <!-- 组织更新组件 -->
    <OrgUpdater 
      v-model:visible="updateOrgModalVisible" 
      @success="handleOrgUpdateSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from '../utils/axios';
import { useRouter } from 'vue-router';
import {
  EditOutlined,
  TransactionOutlined,
  WalletOutlined,
  UserOutlined,
  SettingOutlined,
  LogoutOutlined,
  TeamOutlined,
} from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import { MenuInfo } from 'ant-design-vue/es/menu/src/interface';
import { orgMap } from '../utils';
import AvatarUploader from '../components/AvatarUploader.vue';
import ProfileEditor from '../components/ProfileEditor.vue';
import OrgUpdater from '../components/OrgUpdater.vue';
interface UserInfo {
  username: string;
  avatarURL: string;
}

const router = useRouter();

const user = ref<UserInfo>({
  username: '游客',
  avatarURL: 'https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png'
});

const popoverVisible = ref<boolean>(false);
const changeAvatarModalVisible = ref<boolean>(false);
const editProfileModalVisible = ref<boolean>(false);
const updateOrgModalVisible = ref<boolean>(false);
const userInfo = ref<any>({
  id: '',
  username: '',
  email: '',
  org: '',
});
const isPlatformAdmin = ref<boolean>(false);

onMounted(() => {
  loadUserInfo();
});

// 加载用户信息
const loadUserInfo = () => {
  const userInfoString = localStorage.getItem('userInfo');
  if (userInfoString) {
    try {
      const parsedUserInfo = JSON.parse(userInfoString);
      user.value.username = parsedUserInfo.username || user.value.username;
      user.value.avatarURL = parsedUserInfo.avatarURL || user.value.avatarURL;
      userInfo.value.id = parsedUserInfo.id || userInfo.value.id;
      userInfo.value.username = parsedUserInfo.username || userInfo.value.username;
      userInfo.value.email = parsedUserInfo.email || userInfo.value.email;
      userInfo.value.org = orgMap[parsedUserInfo.org] || userInfo.value.org;
      
      // 检查是否为平台管理员（组织ID为1）
      isPlatformAdmin.value = parsedUserInfo.org === 1;
    } catch (e) {
      console.error('解析 localStorage 中的 userInfo 失败', e);
    }
  }
};

// 处理 Popover 菜单点击事件
const handleMenuClick = (e: MenuInfo) => {
  popoverVisible.value = false;
  if (e.key === 'change_avatar') {
    changeAvatarModalVisible.value = true;
  } 
  else if (e.key === 'edit_profile') {
    editProfileModalVisible.value = true;
  }
  else if (e.key === 'update_org') {
    updateOrgModalVisible.value = true;
  }
  else if (e.key === 'logout') {
    handleLogout();
  }
};

// 登出函数
const handleLogout = async () => {
  try{
    const response = await axios.post('/api/account/logout');
    if (response.status === 200 && response.data.code === 200) {
      axios.defaults.headers.common['Authorization'] = '';
      message.success('已成功登出！');
      localStorage.removeItem('userInfo');
      localStorage.removeItem('userToken');
      router.push('/login');
    }
    else{
      message.error('登出失败！');
    }
  } catch (error) {
    message.error('网络错误，请稍后重试！');
  }
};

// 头像上传成功处理
const handleAvatarUploadSuccess = (avatarUrl: string) => {
  user.value.avatarURL = avatarUrl;
  
  // 更新 localStorage
  const userInfoString = localStorage.getItem('userInfo');
  let parsedUserInfo = {};
  if (userInfoString) {
    try {
      parsedUserInfo = JSON.parse(userInfoString);
    } catch (e) {
      console.error('解析 localStorage 中的 userInfo 失败', e);
    }
  }
  localStorage.setItem('userInfo', JSON.stringify({ ...parsedUserInfo, avatarURL: avatarUrl }));
};

// 个人信息更新成功处理
const handleProfileUpdateSuccess = (updatedInfo: any) => {
  // 更新本地用户信息
  userInfo.value = { ...userInfo.value, ...updatedInfo };
  
  // 更新 localStorage
  const userInfoString = localStorage.getItem('userInfo');
  let parsedUserInfo = {};
  if (userInfoString) {
    try {
      parsedUserInfo = JSON.parse(userInfoString);
    } catch (e) {
      console.error('解析 localStorage 中的 userInfo 失败', e);
    }
  }
  localStorage.setItem('userInfo', JSON.stringify({ ...parsedUserInfo, ...updatedInfo }));
};

// 组织更新成功处理
const handleOrgUpdateSuccess = (result: any) => {
  console.log('组织更新成功:', result);
  // 这里可以根据需要添加额外的处理逻辑
  // 比如刷新用户列表、显示通知等
};
</script>

<style scoped>
/* ... 样式保持不变 ... */
.dashboard-container {
  min-height: 100vh;
  background-color: #f5f7fa;
  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
  color: #333;
}

.user-info-section {
  width: 100%;
  padding: 60px 24px;
  background: linear-gradient(135deg, #4a90e2 0%, #76b1f3 100%);
  color: #fff;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: center;
}

.user-profile {
  display: flex;
  align-items: center;
  max-width: 1200px;
  width: 100%;
  position: relative; /* 为 Popover 定位提供上下文 */
}

.avatar {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  border: 5px solid rgba(255, 255, 255, 0.9);
  object-fit: cover;
  margin-right: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  cursor: pointer; /* 增加手型光标表示可点击 */
  transition: transform 0.2s ease-in-out;
}

.avatar:hover {
  transform: scale(1.05); /* 悬停效果 */
}

.user-details {
  display: flex;
  flex-direction: column;
}

.username {
  font-size: 2.8rem;
  font-weight: 700;
  margin: 0;
  letter-spacing: 0.5px;
}

.greeting {
  font-size: 1.2rem;
  font-weight: 400;
  opacity: 0.9;
  margin: 4px 0 0;
}

/* 主内容区 */
.main-content {
  flex-grow: 1;
  width: 100%;
  max-width: 1200px;
  padding: 40px 24px;
  display: flex;
  justify-content: center;
}

.card-list {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
  width: 100%;
}

/* 卡片样式 */
.card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  border-radius: 16px;
  background-color: #fff;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  text-decoration: none;
  color: #555;
  text-align: center;
  border: 1px solid #e0e6ed;
}

.card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.15);
  border-color: #4a90e2;
}

.card-icon {
  font-size: 68px;
  margin-bottom: 24px;
  color: #4a90e2;
  transition: transform 0.4s ease, color 0.4s ease;
}

.card:hover .card-icon {
  transform: scale(1.15) rotate(5deg);
  color: #1a64b3;
}

.card-title {
  font-size: 1.6rem;
  font-weight: 600;
  margin-bottom: 12px;
  color: #333;
}

.card-description {
  font-size: 1rem;
  color: #888;
  line-height: 1.6;
}

/* Popover 菜单样式 */
.custom-popover-menu {
  padding: 0;
  border: none;
  box-shadow: none;
}
.popover-menu .ant-menu-item {
  padding: 10px 16px;
  margin: 0;
  line-height: unset;
  height: auto;
  font-size: 15px;
}
.popover-menu .ant-menu-item:hover {
  background-color: #f0f2f5;
  color: #4a90e2;
}
.popover-menu .ant-menu-item-icon {
  margin-right: 8px;
}

/* 响应式设计 */
@media (max-width: 992px) {
  .card-list {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 600px) {
  .user-profile {
    flex-direction: column;
    text-align: center;
  }
  .avatar {
    margin-right: 0;
    margin-bottom: 16px;
  }
  .username {
    font-size: 2rem;
  }
  .greeting {
    font-size: 1rem;
  }
  .card-list {
    grid-template-columns: 1fr;
    gap: 24px;
  }
}
</style>