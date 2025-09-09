<template>
  <div class="profile-editor-wrapper">
    <a-modal 
      v-model:visible="modalVisible" 
      title="编辑个人信息" 
      :width="600"
      centered
      class="profile-editor-modal"
      @ok="handleSubmit"
      :confirm-loading="submitting"
    >
      <div class="profile-editor-content">
        <div class="editor-header">
          <div class="header-icon">
            <UserOutlined />
          </div>
          <div class="header-text">
            <h3>修改个人信息</h3>
            <p>更新你的邮箱和密码信息</p>
          </div>
        </div>

        <a-form 
          ref="formRef"
          :model="formData" 
          :rules="formRules"
          :label-col="{ span: 6 }" 
          :wrapper-col="{ span: 16 }"
          class="profile-form"
        >
          <a-form-item label="用户ID">
            <a-input v-model:value="formData.id" disabled class="disabled-input" />
          </a-form-item>

          <a-form-item label="用户名">
            <a-input v-model:value="formData.username" disabled class="disabled-input" />
          </a-form-item>

          <a-form-item label="组织">
            <a-input v-model:value="formData.org" disabled class="disabled-input" />
          </a-form-item>

          <a-form-item label="邮箱" name="email">
            <a-input 
              v-model:value="formData.email" 
              placeholder="请输入新的邮箱地址"
              class="editable-input"
            />
            <div class="validation-text" :class="{ 'error': !isEmailValid && formData.email.length > 0 }" v-if="formData.email.length > 0">
              {{ !isEmailValid ? '请输入有效的邮箱地址' : '邮箱格式正确' }}
            </div>
          </a-form-item>

          <a-form-item label="新密码" name="password">
            <a-input-password 
              v-model:value="formData.password" 
              placeholder="请输入新密码（留空则不修改）"
              class="editable-input"
            />
            <div class="validation-text" :class="{ 'weak': passwordStrength === '弱', 'medium': passwordStrength === '中等', 'strong': passwordStrength === '强' }" v-if="formData.password.length > 0">
              密码强度: {{ passwordStrength }}
            </div>
          </a-form-item>

          <a-form-item label="确认密码" name="confirmPassword">
            <a-input-password 
              v-model:value="formData.confirmPassword" 
              placeholder="请再次输入新密码"
              class="editable-input"
            />
            <div class="validation-text" :class="{ 'error': !isPasswordMatch && formData.confirmPassword.length > 0 }" v-if="formData.confirmPassword.length > 0">
              {{ !isPasswordMatch ? '两次输入的密码不一致' : '密码匹配' }}
            </div>
          </a-form-item>
        </a-form>

        <div class="form-tips">
          <div class="tip-item">
            <InfoCircleOutlined />
            <span>密码留空表示不修改当前密码</span>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive, computed } from 'vue';
import { message } from 'ant-design-vue';
import { 
  UserOutlined, 
  InfoCircleOutlined 
} from '@ant-design/icons-vue';
import type { FormInstance, Rule } from 'ant-design-vue/es/form';
import { accountApi } from '../api';

// Props
interface Props {
  visible?: boolean;
  userInfo?: any;
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  userInfo: () => ({})
});

// Emits
const emit = defineEmits<{
  'update:visible': [value: boolean];
  'success': [updatedInfo: any];
  'cancel': [];
}>();

// 响应式数据
const modalVisible = ref(false);
const submitting = ref(false);
const formRef = ref<FormInstance>();

// 表单数据
const formData = reactive({
  id: '',
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  org: ''
});

// 邮箱格式验证
const isEmailValid = computed(() => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(formData.email) || formData.email.length === 0;
});

// 密码强度验证
const passwordStrength = computed(() => {
  let score = 0;
  if (!formData.password) return '输入密码...';

  const hasLower = /[a-z]/.test(formData.password);
  const hasUpper = /[A-Z]/.test(formData.password);
  const hasNumber = /[0-9]/.test(formData.password);
  const hasSpecial = /[^a-zA-Z0-9]/.test(formData.password);

  if (formData.password.length >= 8) score++;
  if (hasLower) score++;
  if (hasUpper) score++;
  if (hasNumber) score++;
  if (hasSpecial) score++;

  if (score >= 4) return '强';
  if (score >= 2) return '中等';
  return '弱';
});

// 密码匹配验证
const isPasswordMatch = computed(() => {
  return formData.password === formData.confirmPassword || formData.confirmPassword.length === 0;
});

// 表单验证规则
const formRules: Record<string, Rule[]> = {
  email: [
    { 
      type: 'email', 
      message: '请输入有效的邮箱地址',
      trigger: 'blur'
    }
  ],
  password: [
    {
      validator: (_rule: any, value: string) => {
        if (value && value.length < 6) {
          return Promise.reject('密码长度至少6位');
        }
        if (value && passwordStrength.value === '弱') {
          return Promise.reject('密码强度太弱，请使用更复杂的密码');
        }
        return Promise.resolve();
      },
      trigger: 'blur'
    }
  ],
  confirmPassword: [
    {
      validator: (_rule: any, value: string) => {
        if (formData.password && value !== formData.password) {
          return Promise.reject('两次输入的密码不一致');
        }
        return Promise.resolve();
      },
      trigger: 'blur'
    }
  ]
};

// 监听 visible 变化
watch(() => props.visible, (newVal) => {
  modalVisible.value = newVal;
  if (newVal) {
    // 打开时填充表单数据
    fillFormData();
  }
});

// 监听 modalVisible 变化，同步到父组件
watch(modalVisible, (newVal) => {
  emit('update:visible', newVal);
});

// 填充表单数据
const fillFormData = () => {
  if (props.userInfo) {
    formData.id = props.userInfo.id || '';
    formData.username = props.userInfo.username || '';
    formData.email = props.userInfo.email || '';
    formData.org = props.userInfo.org || '';
    formData.password = '';
    formData.confirmPassword = '';
  }
};

// 提交表单
const handleSubmit = async () => {
  try {
    // 验证表单
    await formRef.value?.validate();
    
    submitting.value = true;

    // 准备提交数据
    const submitData: any = {};
    
    // 只提交有值的字段
    if (formData.email) {
      submitData.email = formData.email;
    }
    if (formData.password) {
      submitData.password = formData.password;
    }

    // 如果没有要更新的字段
    if (Object.keys(submitData).length === 0) {
      message.warning('请至少修改一个字段');
      submitting.value = false;
      return;
    }

    // 如果修改密码，检查密码强度
    if (submitData.password && passwordStrength.value === '弱') {
      message.error('密码强度太弱，请使用更复杂的密码');
      submitting.value = false;
      return;
    }

    // 如果修改密码，检查密码匹配
    if (submitData.password && !isPasswordMatch.value) {
      message.error('两次输入的密码不一致');
      submitting.value = false;
      return;
    }

    // 调用后端接口
    const response = await accountApi.updateProfile(submitData);

    if (response.data && response.data.code === 200) {
      message.success('个人信息更新成功！');
      modalVisible.value = false;
      
      // 通知父组件更新成功
      emit('success', {
        ...props.userInfo,
        ...submitData
      });
    } else {
      message.error(response.data?.message || '更新失败，请重试');
    }
  } catch (error: any) {
    console.error('更新个人信息失败:', error);
    if (error.response?.data?.message) {
      message.error(error.response.data.message);
    } else {
      message.error('更新失败，请检查网络或稍后重试！');
    }
  } finally {
    submitting.value = false;
  }
};
</script>

<style scoped>
/* 模态框样式 */
.profile-editor-modal :deep(.ant-modal-header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px 8px 0 0;
}

.profile-editor-modal :deep(.ant-modal-title) {
  color: white;
  font-weight: 600;
}

.profile-editor-modal :deep(.ant-modal-close) {
  color: white;
}

.profile-editor-content {
  padding: 20px 0;
}

/* 头部样式 */
.editor-header {
  display: flex;
  align-items: center;
  margin-bottom: 32px;
  padding: 20px;
  background: #f8f9ff;
  border-radius: 12px;
  border: 1px solid #e6f0ff;
}

.header-icon {
  width: 50px;
  height: 50px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  font-size: 20px;
  color: white;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.header-text h3 {
  margin: 0 0 4px 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.header-text p {
  margin: 0;
  font-size: 14px;
  color: #666;
}

/* 表单样式 */
.profile-form {
  margin-bottom: 24px;
}

.profile-form :deep(.ant-form-item-label) {
  font-weight: 500;
  color: #333;
}

.disabled-input {
  background-color: #f5f5f5 !important;
  color: #999 !important;
  cursor: not-allowed !important;
}

.editable-input {
  border-radius: 6px;
  transition: border-color 0.3s ease;
}

.editable-input:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
}

/* 验证文本样式 */
.validation-text {
  font-size: 12px;
  margin-top: 4px;
  transition: color 0.3s ease;
}

.validation-text.error {
  color: #ff4d4f;
}

.validation-text.weak {
  color: #ff4d4f;
}

.validation-text.medium {
  color: #faad14;
}

.validation-text.strong {
  color: #52c41a;
}

/* 提示信息样式 */
.form-tips {
  background: #f0f9ff;
  border-radius: 8px;
  padding: 16px;
  border: 1px solid #bae6fd;
}

.tip-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 13px;
  color: #0369a1;
}

.tip-item:last-child {
  margin-bottom: 0;
}

.tip-item .anticon {
  color: #0ea5e9;
  margin-right: 8px;
  font-size: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .profile-editor-modal :deep(.ant-modal) {
    margin: 16px;
    max-width: calc(100vw - 32px);
  }
  
  .editor-header {
    flex-direction: column;
    text-align: center;
  }
  
  .header-icon {
    margin-right: 0;
    margin-bottom: 12px;
  }
  
  .profile-form :deep(.ant-form-item-label) {
    text-align: left;
  }
}
</style>
