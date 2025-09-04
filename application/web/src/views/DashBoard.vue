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
        <router-link to="/creation" class="card card-creation">
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

    <template>
      <a-modal v-model:visible="editProfileModalVisible" title="编辑信息" :footer="null">
        <a-form :model="editProfileForm" :label-col="{ span: 4 }" :wrapper-col="{ span: 14 }">
          <a-form-item label="ID">
            <a-input v-model:value="editProfileForm.id" disabled />
          </a-form-item>
          <a-form-item label="用户名">
            <a-input v-model:value="editProfileForm.username" disabled />
          </a-form-item>
          <a-form-item label="密码">
            <a-input v-model:value="editProfileForm.password" placeholder="请输入新的密码" />
          </a-form-item>
          <a-form-item label="邮箱">
            <a-input v-model:value="editProfileForm.email" />
          </a-form-item>
          <a-form-item label="组织">
            <a-input v-model:value="editProfileForm.org" disabled />
          </a-form-item>
        </a-form>
      </a-modal>
    </template>

    <a-modal v-model:visible="changeAvatarModalVisible" title="选择头像图片" :footer="null">
      <a-upload
        name="avatar"
        list-type="picture-card"
        class="avatar-uploader"
        :show-upload-list="false"
        :before-upload="beforeUpload"
        accept="image/jpeg,image/png,image/gif"
      >
        <div v-if="imageUrl" class="uploaded-image-preview">
          <img :src="imageUrl" alt="avatar preview" style="width: 100%; height: 100%; object-fit: cover;" />
        </div>
        <div v-else>
          <plus-outlined />
          <div class="ant-upload-text">上传</div>
        </div>
      </a-upload>
      <div class="upload-tip">点击上方区域选择图片，将自动进入裁剪。</div>
    </a-modal>

    <a-modal
      v-model:visible="cropperModalVisible"
      title="裁剪头像"
      @ok="handleCropOk"
      @cancel="handleCropCancel"
      :width="600"
      :confirm-loading="uploading"
    >
      <div class="cropper-container">
        <vue-cropper
          ref="cropper"
          :src="cropperImgSrc"
          :aspectRatio="1"
          :autoCropArea="0.8"
          :viewMode="1"
          :guides="true"
          :cropBoxMovable="true"
          :cropBoxResizable="true"
          dragMode="move"
          />
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';
import {
  EditOutlined,
  TransactionOutlined,
  WalletOutlined,
  UserOutlined,
  SettingOutlined,
  LogoutOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import { MenuInfo } from 'ant-design-vue/es/menu/src/interface';
import VueCropper from 'vue-cropperjs';
import 'cropperjs/dist/cropper.css';
import { transformOrg } from '../utils';
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
const changeAvatarModalVisible = ref<boolean>(false); // 选择文件模态框
const cropperModalVisible = ref<boolean>(false); // 裁剪模态框
const uploading = ref<boolean>(false); 
const imageUrl = ref<string | undefined>(undefined); // 用于`a-upload`的预览
const cropperImgSrc = ref(''); // 用于`vue-cropper`的图片源
const cropper = ref<any>(null); // `vue-cropper`实例的引用
const originalFileType = ref(''); // 存储原始文件类型
const editProfileModalVisible = ref<boolean>(false); // 编辑信息模态框
const editProfileForm = ref<any>({
  id: '',
  username: '',
  password: '',
  email: '',
  org: '',
});


onMounted(() => {
  loadUserInfo();
});

// 加载用户信息
const loadUserInfo = () => {
  const userInfoString = localStorage.getItem('userInfo');
  if (userInfoString) {
    try {
      const userInfo = JSON.parse(userInfoString);
      user.value.username = userInfo.username || user.value.username;
      user.value.avatarURL = userInfo.avatarURL || user.value.avatarURL;
      editProfileForm.value.id = userInfo.id || editProfileForm.value.id;
      editProfileForm.value.username = userInfo.username || editProfileForm.value.username;
      editProfileForm.value.email = userInfo.email || editProfileForm.value.email;
      editProfileForm.value.org = transformOrg(userInfo.org) || editProfileForm.value.org;
    } catch (e) {
      console.error('解析 localStorage 中的 userInfo 失败', e);
    }
  }
  // 这里不再需要将 user.value.avatarURL 赋给 imageUrl.value，因为 imageUrl 只用于裁剪前的临时预览
};

// 处理 Popover 菜单点击事件
const handleMenuClick = (e: MenuInfo) => {
  popoverVisible.value = false;
  if (e.key === 'change_avatar') {
    changeAvatarModalVisible.value = true; // 打开文件选择弹窗
    imageUrl.value = undefined; // 清空上次选择的预览图
  } 
  else if (e.key === 'edit_profile') {
    editProfileModalVisible.value = true;
  }
  else if (e.key === 'logout') {
    handleLogout();
  }
};

// 登出函数
const handleLogout = async () => {
  localStorage.removeItem('userInfo');
  localStorage.removeItem('userToken');
  try{
    const response = await axios.post('http://localhost:8888/api/account/logout');
    if (response.status === 200 && response.data.code === 200) {
      axios.defaults.headers.common['Authorization'] = '';
      message.success('已成功登出！');
      router.push('/login');
    }
    else{
      message.error('登出失败！');
    }
  } catch (error) {
    message.error('网络错误，请稍后重试！');
  }
};

// a-upload 的 beforeUpload 钩子
const beforeUpload = (file: File) => {
  // 检查文件类型
  const isImage = ['image/jpeg', 'image/png', 'image/gif'].includes(file.type);
  if (!isImage) {
    message.error('只能上传 JPG/PNG/GIF 格式的图片!');
    return false;
  }
  // 检查文件大小
  const isLt2M = file.size / 1024 / 1024 < 2;
  if (!isLt2M) {
    message.error('图片大小不能超过 2MB!');
    return false;
  }

  // 使用 FileReader 读取文件并将其作为裁剪组件的源
  const reader = new FileReader();
  reader.readAsDataURL(file);
  reader.onload = (e) => {
    cropperImgSrc.value = e.target?.result as string; // 设置裁剪组件的图片源
    changeAvatarModalVisible.value = false; // 关闭文件选择模态框
    cropperModalVisible.value = true; // 打开裁剪模态框
  };
  originalFileType.value = file.type; // 保存原始文件类型

  // 返回 false 阻止 antd 的默认上传行为
  return false;
};


// 裁剪模态框的确定按钮事件
const handleCropOk = () => {
  if (!cropper.value) {
    message.error('裁剪组件未加载完成');
    return;
  }

  uploading.value = true; // 开始上传，显示加载状态

  // 获取裁剪后的 Blob 数据，并使用动态类型
  cropper.value.getCroppedCanvas().toBlob(async (blob: Blob | null) => {
    if (blob) {
      // 创建一个新的 File 对象用于上传
      const fileExtension = originalFileType.value.split('/')[1] || 'jpeg'; // 默认jpeg
      const croppedFile = new File([blob], `avatar.${fileExtension}`, { type: originalFileType.value });
      
      try {
        const formData = new FormData();
        formData.append('avatar', croppedFile);
        
        // 调用后端上传接口
        const response = await axios.put('http://localhost:8888/api/account/avatar', formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        });
        
        if (response.data && response.data.code === 200) {
          const newAvatarUrl = response.data.data.avatarURL || URL.createObjectURL(croppedFile);
          user.value.avatarURL = newAvatarUrl;
          
          // 更新 localStorage
          const userInfoString = localStorage.getItem('userInfo');
          let userInfo = {};
          if (userInfoString) {
            try {
              userInfo = JSON.parse(userInfoString);
            } catch (e) {
              console.error('解析 localStorage 中的 userInfo 失败', e);
            }
          }
          localStorage.setItem('userInfo', JSON.stringify({ ...userInfo, avatarURL: newAvatarUrl }));
          
          message.success('头像上传成功！');
          cropperModalVisible.value = false; // 上传成功后关闭裁剪弹窗
        } else {
          message.error(response.data?.message || '上传失败，请重试');
        }
      } catch (error) {
        console.error('上传失败:', error);
        message.error('上传失败，请检查网络或稍后重试！');
      } finally {
        uploading.value = false;
      }
      
    } else {
      message.error('裁剪失败，请重试');
      uploading.value = false;
    }
  }, originalFileType.value, 0.9); // 第二个参数是 MIME type，第三个参数是图片质量 (0-1)
};

// 裁剪模态框取消事件
const handleCropCancel = () => {
  cropperModalVisible.value = false;
  cropperImgSrc.value = ''; // 清空图片源
  originalFileType.value = ''; // 清空文件类型
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

.avatar-uploader {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 24px;
}

.avatar-uploader > .ant-upload {
  width: 128px;
  height: 128px;
}

.ant-upload-select-picture-card i {
  font-size: 32px;
  color: #999;
}
.ant-upload-select-picture-card .ant-upload-text {
  margin-top: 8px;
  color: #666;
}

.upload-tip {
  text-align: center;
  color: #999;
  font-size: 0.9em;
  margin-top: 10px;
}

/* 裁剪容器样式 */
.cropper-container {
  max-width: 100%;
  height: 300px; /* 给裁剪容器一个固定高度，防止图片过大撑开 */
  display: flex;
  justify-content: center;
  align-items: center;
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