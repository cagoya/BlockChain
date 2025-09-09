<template>
  <div class="avatar-uploader-wrapper">
    <a-modal 
      v-model:visible="uploadModalVisible" 
      title="选择头像图片" 
      :footer="null"
      :width="500"
      centered
      class="avatar-upload-modal"
    >
      <div class="upload-section">
        <div class="upload-header">
          <div class="upload-icon">
            <UserOutlined />
          </div>
          <h3>上传新头像</h3>
          <p>支持 JPG、PNG 格式，文件大小不超过 2MB</p>
        </div>
        
        <a-upload
          name="avatar"
          list-type="picture-card"
          class="avatar-uploader"
          :show-upload-list="false"
          :before-upload="beforeUpload"
          accept="image/jpeg,image/png"
        >
          <div v-if="imageUrl" class="uploaded-image-preview">
            <img :src="imageUrl" alt="avatar preview" />
            <div class="image-overlay">
              <div class="overlay-content">
                <a-button 
                  type="primary" 
                  size="small" 
                  @click.stop="handleUpload"
                  :loading="uploading"
                  class="action-btn"
                >
                  <template #icon><CloudUploadOutlined /></template>
                  确认上传
                </a-button>
                <a-button 
                  type="text" 
                  size="small" 
                  danger
                  @click.stop="handleDeletePreview"
                  class="action-btn delete-btn"
                >
                  <template #icon><DeleteOutlined /></template>
                  删除
                </a-button>
              </div>
            </div>
          </div>
          <div v-else class="upload-placeholder">
            <div class="upload-icon-large">
              <CloudUploadOutlined />
            </div>
            <div class="upload-text">
              点击上传图片
            </div>
          </div>
        </a-upload>
        
        <div class="upload-tips">
          <div class="tip-item">
            <CheckCircleOutlined />
            <span>支持 JPG、PNG 格式</span>
          </div>
          <div class="tip-item">
            <CheckCircleOutlined />
            <span>文件大小不超过 2MB</span>
          </div>
          <div class="tip-item">
            <CheckCircleOutlined />
            <span>建议使用正方形图片</span>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>
  
  <script setup lang="ts">
  import { ref, watch } from 'vue';
  import { message } from 'ant-design-vue';
import { 
  UserOutlined, 
  DeleteOutlined, 
  CloudUploadOutlined, 
  CheckCircleOutlined
} from '@ant-design/icons-vue';
import { accountApi } from '../api';
  
  // Props
  interface Props {
    visible?: boolean;
  }
  
  const props = withDefaults(defineProps<Props>(), {
    visible: false
  });
  
  // Emits
  const emit = defineEmits<{
    'update:visible': [value: boolean];
    'success': [avatarUrl: string];
    'cancel': [];
  }>();
  
// 响应式数据
const uploadModalVisible = ref(false);
const uploading = ref(false);
const imageUrl = ref<string | undefined>(undefined);
const selectedFile = ref<File | null>(null);
  
  // 监听 visible 变化
  watch(() => props.visible, (newVal) => {
    uploadModalVisible.value = newVal;
    if (newVal) {
      // 打开时清空之前的数据
      clearData();
    }
  });
  
  // 监听 uploadModalVisible 变化，同步到父组件
  watch(uploadModalVisible, (newVal) => {
    emit('update:visible', newVal);
  });
  
// 清空数据函数
const clearData = () => {
  imageUrl.value = undefined;
  selectedFile.value = null;
};
  
// 文件上传前的处理
const beforeUpload = (file: File) => {
  // 检查文件类型
  const isImage = ['image/jpeg', 'image/png'].includes(file.type);
  if (!isImage) {
    message.error('只能上传 JPG/PNG 格式的图片!');
    return false;
  }
  
  // 检查文件大小
  const isLt2M = file.size / 1024 / 1024 < 2;
  if (!isLt2M) {
    message.error('图片大小不能超过 2MB!');
    return false;
  }

  // 先清空之前的数据，确保每次上传都是干净的状态
  clearData();

  // 保存选中的文件
  selectedFile.value = file;

  // 使用 FileReader 读取文件并显示预览
  const reader = new FileReader();
  reader.readAsDataURL(file);
  reader.onload = (e) => {
    const result = e.target?.result as string;
    imageUrl.value = result; // 设置预览图片
  };
  reader.onerror = () => {
    message.error('文件读取失败，请重试');
  };

  // 返回 false 阻止 antd 的默认上传行为
  return false;
};
  
// 删除预览图片
const handleDeletePreview = () => {
  clearData();
  message.success('已删除预览图片');
};
  
// 确认上传
const handleUpload = async () => {
  if (!selectedFile.value) {
    message.error('请先选择图片');
    return;
  }

  uploading.value = true;

  try {
    const formData = new FormData();
    formData.append('avatar', selectedFile.value);
    
    // 调用后端上传接口
    const response = await accountApi.updateAvatar(formData);
    
    if (response.data && response.data.code === 200) {
      const newAvatarUrl = response.data.data || URL.createObjectURL(selectedFile.value);
      
      message.success('头像上传成功！');
      uploadModalVisible.value = false;
      
      // 清空所有数据
      clearData();
      
      // 通知父组件上传成功
      emit('success', newAvatarUrl);
    } else {
      message.error(response.data?.message || '上传失败，请重试');
    }
  } catch (error) {
    console.error('上传失败:', error);
    message.error('上传失败，请检查网络或稍后重试！');
  } finally {
    uploading.value = false;
  }
};
  </script>
  
  <style scoped>
  /* 上传模态框样式 */
  .avatar-upload-modal :deep(.ant-modal-header) {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 8px 8px 0 0;
  }
  
  .avatar-upload-modal :deep(.ant-modal-title) {
    color: white;
    font-weight: 600;
  }
  
  .avatar-upload-modal :deep(.ant-modal-close) {
    color: white;
  }
  
  .upload-section {
    padding: 20px 0;
  }
  
  .upload-header {
    text-align: center;
    margin-bottom: 32px;
  }
  
  .upload-icon {
    width: 60px;
    height: 60px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 16px;
    font-size: 24px;
    color: white;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }
  
  .upload-header h3 {
    margin: 0 0 8px 0;
    font-size: 20px;
    font-weight: 600;
    color: #333;
  }
  
  .upload-header p {
    margin: 0;
    color: #666;
    font-size: 14px;
  }
  
  /* 上传区域样式 */
  .avatar-uploader {
    display: flex;
    justify-content: center;
    margin-bottom: 24px;
  }
  
  .avatar-uploader > .ant-upload {
    width: 160px;
    height: 160px;
    border-radius: 12px;
    border: 2px dashed #d9d9d9;
    transition: all 0.3s ease;
  }
  
  .avatar-uploader > .ant-upload:hover {
    border-color: #667eea;
    background-color: #f8f9ff;
  }
  
  .uploaded-image-preview {
    position: relative;
    width: 100%;
    height: 100%;
    border-radius: 12px;
    overflow: hidden;
  }
  
  .uploaded-image-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .image-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s ease;
    border-radius: 12px;
  }
  
  .uploaded-image-preview:hover .image-overlay {
    opacity: 1;
  }
  
  .overlay-content {
    display: flex;
    flex-direction: column;
    gap: 12px;
    align-items: center;
  }
  
  .action-btn {
    min-width: 100px;
    border-radius: 6px;
    font-weight: 500;
  }
  
  .delete-btn {
    background: rgba(255, 77, 79, 0.1);
    border: 1px solid rgba(255, 77, 79, 0.3);
  }
  
  .delete-btn:hover {
    background: rgba(255, 77, 79, 0.2);
    border-color: rgba(255, 77, 79, 0.5);
  }
  
  /* 上传占位符样式 */
  .upload-placeholder {
    display: flex; /* 启用 Flexbox */
    flex-direction: column; /* 垂直排列子元素 */
    align-items: center; /* 水平居中 */
    justify-content: center; /* 垂直居中 */
    height: 100%;
    padding: 20px;
  }
  
  .upload-icon-large {
    font-size: 48px;
    color: #d9d9d9;
    margin-bottom: 16px;
    transition: color 0.3s ease;
  }
  
  .avatar-uploader:hover .upload-icon-large {
    color: #667eea;
  }
  
  .upload-text {
    text-align: center;
    font-size: 10px;
    font-weight: 500;
    color: #333;

  }
  
  /* 提示信息样式 */
  .upload-tips {
    background: #f8f9ff;
    border-radius: 8px;
    padding: 16px;
    border: 1px solid #e6f0ff;
  }
  
  .tip-item {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    font-size: 14px;
    color: #666;
  }
  
  .tip-item:last-child {
    margin-bottom: 0;
  }
  
  .tip-item .anticon {
    color: #52c41a;
    margin-right: 8px;
    font-size: 12px;
  }
  
  
/* 响应式设计 */
@media (max-width: 768px) {
  .avatar-uploader > .ant-upload {
    width: 140px;
    height: 140px;
  }
}
  </style>