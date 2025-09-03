<template>
    <div>
      <a-upload
        name="file"
        :show-upload-list="false"
        :before-upload="beforeUpload"
      >
        <a-button>
          <upload-outlined />
          选择图片
        </a-button>
      </a-upload>
  
      <a-modal
        v-model:visible="cropperModalVisible"
        title="图片裁剪"
        @ok="handleCropOk"
        @cancel="handleCropCancel"
        :width="600"
      >
        <div class="cropper-container">
          <vue-cropper
            ref="cropper"
            :src="cropperImgSrc"
            :aspectRatio="1"
            preview=".preview"
          />
        </div>
      </a-modal>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  import { message } from 'ant-design-vue';
  import { UploadOutlined } from '@ant-design/icons-vue';
  import VueCropper from 'vue-cropperjs';
  import 'cropperjs/dist/cropper.css';
  
  const cropperModalVisible = ref(false);
  const cropperImgSrc = ref('');
  const cropper = ref(null);
  
  // 1. 拦截 beforeUpload 事件，阻止自动上传并打开裁剪模态框
  const beforeUpload = (file) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      cropperImgSrc.value = e.target.result;
      cropperModalVisible.value = true;
    };
    reader.readAsDataURL(file);
    return false; // 阻止 a-upload 默认行为
  };
  
  // 2. 裁剪模态框的确定按钮事件
  const handleCropOk = () => {
    if (!cropper.value) {
      message.error('裁剪组件未加载完成');
      return;
    }
  
    // 获取裁剪后的 Blob 数据
    cropper.value.getCroppedCanvas().toBlob((blob) => {
      if (blob) {
        // 创建一个新的 File 对象用于上传
        const croppedFile = new File([blob], 'cropped_image.jpeg', { type: 'image/jpeg' });
        
        // 3. 将新文件手动上传
        uploadFile(croppedFile);
  
        // 关闭模态框
        cropperModalVisible.value = false;
      } else {
        message.error('裁剪失败，请重试');
      }
    }, 'image/jpeg');
  };
  
  // 4. 自定义上传逻辑
  const uploadFile = (file) => {
    const formData = new FormData();
    formData.append('file', file);
    message.success('图片上传成功！');
  };
  
  const handleCropCancel = () => {
    cropperModalVisible.value = false;
    cropperImgSrc.value = '';
  };
  </script>
  
  <style scoped>
  .cropper-container {
    max-width: 100%;
    max-height: 400px;
  }
  </style>