<template>
  <div class="org-updater-wrapper">
    <a-modal 
      v-model:visible="modalVisible" 
      title="更新用户组织" 
      :width="600"
      centered
      class="org-updater-modal"
      @ok="handleSubmit"
      @cancel="handleCancel"
      :confirm-loading="submitting"
    >
      <div class="org-updater-content">
        <div class="updater-header">
          <div class="header-icon">
            <TeamOutlined />
          </div>
          <div class="header-text">
            <h3>更新用户组织</h3>
            <p>只有平台运营方可以修改用户组织</p>
          </div>
        </div>

        <a-form 
          ref="formRef"
          :model="formData" 
          :rules="formRules"
          :label-col="{ span: 6 }" 
          :wrapper-col="{ span: 16 }"
          class="org-form"
        >
          <a-form-item label="用户ID" name="userID">
            <a-input 
              v-model:value="formData.userID" 
              placeholder="请输入要更新的用户ID"
              class="editable-input"
            />
          </a-form-item>

          <a-form-item label="新组织" name="org">
            <a-select 
              v-model:value="formData.org" 
              placeholder="请选择新的组织"
              class="editable-input"
            >
              <a-select-option 
                v-for="org in organizationOptions" 
                :key="org.value" 
                :value="org.value"
              >
                {{ org.label }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-form>

        <div class="form-tips">
          <div class="tip-item">
            <InfoCircleOutlined />
            <span>组织变更重新登录后生效</span>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive } from 'vue';
import axios from '../utils/axios';
import { message } from 'ant-design-vue';
import { 
  TeamOutlined, 
  InfoCircleOutlined 
} from '@ant-design/icons-vue';
import type { FormInstance, Rule } from 'ant-design-vue/es/form';

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
  'success': [result: any];
  'cancel': [];
}>();

// 响应式数据
const modalVisible = ref(false);
const submitting = ref(false);
const formRef = ref<FormInstance>();

// 组织选项
const organizationOptions = [
  { label: '平台运营方', value: 1 },
  { label: 'NFT创作者', value: 2 },
  { label: '金融机构', value: 3 }
];

// 表单数据
const formData = reactive({
  userID: '',
  org: undefined as number | undefined
});

// 表单验证规则
const formRules: Record<string, Rule[]> = {
  userID: [
    { 
      required: true, 
      message: '请输入用户ID',
      trigger: 'blur'
    },
    {
      validator: (_rule: any, value: string) => {
        if (value && !/^\d+$/.test(value)) {
          return Promise.reject('用户ID必须是数字');
        }
        return Promise.resolve();
      },
      trigger: 'blur'
    }
  ],
  org: [
    { 
      required: true, 
      message: '请选择新的组织',
      trigger: 'change'
    }
  ]
};

// 监听 visible 变化
watch(() => props.visible, (newVal) => {
  modalVisible.value = newVal;
  if (newVal) {
    // 打开时清空表单数据
    resetForm();
  }
});

// 监听 modalVisible 变化，同步到父组件
watch(modalVisible, (newVal) => {
  emit('update:visible', newVal);
});

// 重置表单
const resetForm = () => {
  formData.userID = '';
  formData.org = undefined;
  formRef.value?.resetFields();
};

// 提交表单
const handleSubmit = async () => {
  try {
    // 验证表单
    await formRef.value?.validate();
    
    submitting.value = true;

    // 准备提交数据
    const submitData = {
      userID: parseInt(formData.userID),
      org: formData.org
    };

    // 调用后端接口
    const response = await axios.put('/api/account/org', submitData, {
      headers: {
        'Content-Type': 'application/json'
      }
    });

    if (response.data && response.data.code === 200) {
      message.success('用户组织更新成功！');
      modalVisible.value = false;
      
      // 通知父组件更新成功
      emit('success', submitData);
    } else {
      message.error(response.data?.message || '更新失败，请重试');
    }
  } catch (error: any) {
    console.error('更新用户组织失败:', error);
    if (error.response?.data?.message) {
      message.error(error.response.data.message);
    } else {
      message.error('更新失败，请检查网络或稍后重试！');
    }
  } finally {
    submitting.value = false;
  }
};

// 取消操作
const handleCancel = () => {
  modalVisible.value = false;
  emit('cancel');
};
</script>

<style scoped>
/* 模态框样式 */
.org-updater-modal :deep(.ant-modal-header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px 8px 0 0;
}

.org-updater-modal :deep(.ant-modal-title) {
  color: white;
  font-weight: 600;
}

.org-updater-modal :deep(.ant-modal-close) {
  color: white;
}

.org-updater-content {
  padding: 20px 0;
}

/* 头部样式 */
.updater-header {
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
.org-form {
  margin-bottom: 24px;
}

.org-form :deep(.ant-form-item-label) {
  font-weight: 500;
  color: #333;
}

.editable-input {
  border-radius: 6px;
  transition: border-color 0.3s ease;
}

.editable-input:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
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
  .org-updater-modal :deep(.ant-modal) {
    margin: 16px;
    max-width: calc(100vw - 32px);
  }
  
  .updater-header {
    flex-direction: column;
    text-align: center;
  }
  
  .header-icon {
    margin-right: 0;
    margin-bottom: 12px;
  }
  
  .org-form :deep(.ant-form-item-label) {
    text-align: left;
  }
}
</style>
