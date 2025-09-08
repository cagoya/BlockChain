<template>
  <div class="page-container">
    <AssetNav />
    <div class="asset-upload-container">
      <div class="upload-header">
        <h1>NFT 资产上传</h1>
        <p>上传您的数字资产到区块链网络，开启数字收藏之旅</p>
      </div>

      <div class="upload-form">
        <a-form
          :model="formData"
          :rules="rules"
          @finish="handleSubmit"
          layout="vertical"
          class="upload-form-content"
        >
          <a-form-item label="资产名称" name="name">
            <a-input
              v-model:value="formData.name"
              placeholder="请输入 NFT 资产的名称（例如：数字艺术品 #001）"
              size="large"
            />
          </a-form-item>

          <a-form-item label="资产描述" name="description">
            <a-textarea
              v-model:value="formData.description"
              placeholder="请详细描述您的 NFT 资产（可选，但推荐填写）"
              :rows="4"
              size="large"
            />
          </a-form-item>

          <a-form-item label="资产图片" name="image">
            <div class="image-upload-section">
              <a-upload
                :before-upload="beforeUpload"
                :show-upload-list="false"
                accept="image/*"
                class="upload-area"
              >
                <div class="upload-placeholder" v-if="!previewImage && !croppedImage">
                  <a-icon type="cloud-upload" class="upload-icon" />
                  <div class="upload-text">点击或拖拽图片到此处上传</div>
                  <div class="upload-hint">支持 JPG、PNG 格式，建议尺寸 400x400，大小不超过 10MB</div>
                </div>
                <div v-else class="preview-cropper-container">
                  <img :src="previewImage" alt="上传图片" class="full-width-image" v-if="previewImage && !croppedImage" />
                  <img :src="croppedImage" alt="裁剪图片" class="full-width-image" v-if="croppedImage" />

                  <div class="action-buttons">
                    <a-button @click="openCropper" v-if="previewImage && !croppedImage" class="crop-button" type="primary">
                      <a-icon type="scissor" /> 裁剪图片
                    </a-button>
                    <a-button @click="resetUpload" class="reset-button" danger>
                      <a-icon type="delete" /> 重新上传
                    </a-button>
                  </div>
                </div>
              </a-upload>
            </div>
          </a-form-item>

          <a-form-item>
            <a-button
              type="primary"
              html-type="submit"
              size="large"
              :loading="uploading"
              :disabled="!formData.name || !croppedImage"
              class="submit-button"
            >
              {{ uploading ? '上传中...' : '创建 NFT 资产' }}
            </a-button>
          </a-form-item>
        </a-form>
      </div>
    </div>

    <a-modal
      v-model:visible="cropperModalVisible"
      title="裁剪图片"
      @ok="cropImage"
      @cancel="cropperModalVisible = false"
      :width="700"
      :maskClosable="false"
      destroyOnClose
      okText="确认裁剪"
      cancelText="取消"
    >
      <div class="cropper-wrapper">
        <Cropper
          ref="cropperRef"
          :src="previewImage"
          :stencil-props="{
            aspectRatio: 1 / 1,
            resizable: true, // 允许用户调整裁剪框大小
            movable: true,   // 允许用户移动裁剪框
          }"
          :resize-image="{
            adjustHeight: true,
            adjustWidth: true
          }"
          image-restriction="stencil"
        />
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import AssetNav from '../components/AssetNav.vue';
import { ref, reactive } from 'vue';
import { message } from 'ant-design-vue';
import type { Rule } from 'ant-design-vue/es/form';
import axios from '../utils/axios';
import { Cropper } from 'vue-advanced-cropper';
import 'vue-advanced-cropper/dist/style.css';
import 'vue-advanced-cropper/dist/theme.bubble.css'; // 使用 bubble 主题，更现代

interface FormData {
  name: string;
  description: string;
  image: Blob | null; // 修改为Blob，因为裁剪后会生成Blob
}

const formData = reactive<FormData>({
  name: '',
  description: '',
  image: null
});

const rules: Record<string, Rule[]> = {
  name: [
    { required: true, message: '请输入资产名称', trigger: 'blur' },
    { min: 1, max: 50, message: '资产名称长度应在1-50个字符之间', trigger: 'blur' }
  ]
};

const uploading = ref(false);
const previewImage = ref(''); // 用于显示原始上传的图片
const croppedImage = ref(''); // 用于显示裁剪后的图片
const cropperModalVisible = ref(false);
const cropperRef = ref<InstanceType<typeof Cropper> | null>(null);

// 上传前处理
const beforeUpload = (file: File) => {
  const isImage = file.type.startsWith('image/');
  if (!isImage) {
    message.error('只能上传图片文件!');
    return false;
  }

  const isLt10M = file.size / 1024 / 1024 < 10;
  if (!isLt10M) {
    message.error('图片大小不能超过 10MB!');
    return false;
  }

  previewImage.value = URL.createObjectURL(file);
  // 在这里不直接设置 formData.image，而是等到裁剪后再设置
  formData.image = null; // 清空之前的裁剪结果
  croppedImage.value = ''; // 清空之前的裁剪预览

  return false; // 阻止 Ant Design Vue 自动上传
};

// 打开裁剪模态框
const openCropper = () => {
  if (previewImage.value) {
    cropperModalVisible.value = true;
  } else {
    message.warning('请先上传图片');
  }
};

// 执行裁剪
const cropImage = () => {
  if (cropperRef.value) {
    const { canvas } = cropperRef.value.getResult();
    if (canvas) {
      canvas.toBlob((blob) => {
        if (blob) {
          formData.image = blob;
          croppedImage.value = URL.createObjectURL(blob); // 显示裁剪后的图片
          message.success('图片裁剪成功！');
          cropperModalVisible.value = false;
        } else {
          message.error('无法生成裁剪后的图片');
        }
      }, 'image/jpeg', 0.9); // 可以指定图片格式和质量
    }
  }
};

// 重置上传
const resetUpload = () => {
  formData.name = '';
  formData.description = '';
  formData.image = null;
  previewImage.value = '';
  croppedImage.value = '';
  cropperModalVisible.value = false;
  if (cropperRef.value) {
    cropperRef.value.reset();
  }
};

// 提交表单
const handleSubmit = async () => {
  if (!formData.image) {
    message.error('请上传并裁剪图片');
    return;
  }

  uploading.value = true;

  try {
    const formDataToSend = new FormData();
    formDataToSend.append('name', formData.name);
    formDataToSend.append('description', formData.description || '暂无描述');
    formDataToSend.append('image', formData.image, 'cropped_image.jpeg'); // Blob需要指定文件名

    const response = await axios.post('/api/asset/create', formDataToSend, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
    });

    const result = response.data;

    if (result.code === 200) {
      message.success('NFT 资产创建成功！');
      resetUpload(); // 重置表单和图片
    } else {
      message.error(result.message || '创建失败');
    }
  } catch (error) {
    console.error('上传失败:', error);
    message.error('上传失败，请重试');
  } finally {
    uploading.value = false;
  }
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
.asset-upload-container {
  max-width: 900px; /* 增加最大宽度 */
  margin: 40px auto; /* 增加上下边距 */
  padding: 40px;
  background: linear-gradient(135deg, #ffffff, #f0f8ff); /* 渐变背景 */
  border-radius: 16px; /* 更大的圆角 */
  box-shadow: 0 10px 40px rgba(41, 98, 255, 0.15); /* 蓝色系更深的阴影 */
  border: 1px solid #e0eaff; /* 浅蓝色边框 */
}

/* 头部样式 */
.upload-header {
  text-align: center;
  margin-bottom: 50px; /* 增加间距 */
}

.upload-header h1 {
  color: #1a237e; /* 深蓝色 */
  font-size: 36px; /* 字体加大 */
  margin-bottom: 15px;
  font-weight: 700; /* 加粗 */
  letter-spacing: 1px;
}

.upload-header p {
  color: #5a667b; /* 柔和的灰色 */
  font-size: 18px;
  line-height: 1.6;
}

/* 表单内容 */
.upload-form-content {
  max-width: 680px; /* 调整表单最大宽度 */
  margin: 0 auto;
}

/* Ant Design Form Item 标签样式 */
:deep(.ant-form-item-label > label) {
  font-size: 16px;
  font-weight: 600;
  color: #3f51b5; /* 中等蓝色 */
  margin-bottom: 8px;
}

/* Ant Design Input 和 Textarea 样式 */
:deep(.ant-input), :deep(.ant-textarea) {
  border-radius: 8px;
  border: 1px solid #c5cae9; /* 浅蓝色边框 */
  padding: 10px 15px;
  font-size: 16px;
  transition: all 0.3s ease;
}

:deep(.ant-input:focus), :deep(.ant-textarea:focus) {
  border-color: #2962ff; /* 聚焦时深蓝色 */
  box-shadow: 0 0 0 2px rgba(41, 98, 255, 0.2);
}

/* 图片上传区域 */
.image-upload-section {
  margin-bottom: 30px;
}

.upload-area {
  width: 100%;
}

.upload-placeholder {
  width: 100%;
  min-height: 350px; /* 增加高度 */
  border: 3px dashed #90caf9; /* 蓝色虚线边框 */
  border-radius: 12px; /* 大圆角 */
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background-color: #e3f2fd; /* 浅蓝色背景 */
}

.upload-placeholder:hover {
  border-color: #2196f3; /* 悬停时深蓝色 */
  background-color: #bbdefb; /* 悬停时背景更深 */
  transform: translateY(-2px); /* 轻微上浮效果 */
  box-shadow: 0 5px 15px rgba(33, 150, 243, 0.2);
}

.upload-icon {
  font-size: 64px; /* 图标加大 */
  color: #42a5f5; /* 蓝色 */
  margin-bottom: 20px;
}

.upload-text {
  font-size: 18px;
  color: #3f51b5; /* 中等蓝色 */
  margin-bottom: 10px;
  font-weight: 500;
}

.upload-hint {
  font-size: 14px;
  color: #7986cb; /* 柔和的蓝色 */
  text-align: center;
  padding: 0 20px;
}

.preview-cropper-container {
  width: 100%;
  min-height: 350px;
  border-radius: 12px;
  overflow: hidden;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 1px solid #bbdefb;
  padding: 20px; /* 增加内边距 */
  gap: 15px; /* 增加元素间距 */
  background-color: #f7f9fc;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.full-width-image {
  max-width: 100%;
  max-height: 280px; /* 调整预览图高度 */
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  background-color: #fff; /* 防止透明背景 */
}

.action-buttons {
  display: flex;
  gap: 15px; /* 按钮间距 */
  margin-top: 15px;
}

.crop-button, .reset-button {
  width: 160px; /* 按钮宽度统一 */
  height: 45px; /* 按钮高度 */
  font-size: 16px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.crop-button {
  background-color: #2962ff; /* 深蓝色 */
  border-color: #2962ff;
}

.crop-button:hover {
  background-color: #004acb;
  border-color: #004acb;
  transform: translateY(-1px);
}

.reset-button {
  background-color: #ef5350; /* 红色，表示删除/重置操作 */
  border-color: #ef5350;
  color: #fff;
}

.reset-button:hover {
  background-color: #d32f2f;
  border-color: #d32f2f;
  transform: translateY(-1px);
}


/* 提交按钮 */
.submit-button {
  width: 100%;
  height: 55px; /* 按钮高度增加 */
  font-size: 18px;
  font-weight: 600;
  border-radius: 10px; /* 大圆角 */
  background-color: #1a237e; /* 深蓝色 */
  border-color: #1a237e;
  box-shadow: 0 5px 20px rgba(26, 35, 126, 0.2);
  transition: all 0.3s ease;
}

.submit-button:hover:not([disabled]) {
  background-color: #0d125f; /* 悬停时更深 */
  border-color: #0d125f;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(26, 35, 126, 0.3);
}

.submit-button[disabled] {
  opacity: 0.7;
  cursor: not-allowed;
  background-color: #9fa8da !important;
  border-color: #9fa8da !important;
  box-shadow: none;
}

/* 裁剪器样式 */
.cropper-wrapper {
  width: 100%;
  height: 450px; /* 裁剪区域的高度增加 */
  background: #f0f4f8; /* 浅灰色背景 */
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  overflow: hidden; /* 确保裁剪器不会溢出 */
}

/* 覆盖 Ant Design Vue 的 modal footer 样式，让按钮居中 */
:deep(.ant-modal-footer) {
  text-align: center;
  padding: 15px 16px;
  border-top: 1px solid #e0eaff; /* 浅蓝色边框 */
}

:deep(.ant-modal-title) {
  color: #1a237e;
  font-weight: 600;
}

:deep(.ant-modal-header) {
  border-bottom: 1px solid #e0eaff;
}
</style>