<template>
  <a-modal
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    title="资产交易状态"
    :width="600"
    :footer="null"
  >
    <div class="asset-status-dialog">
      <!-- 资产信息展示 -->
      <div class="asset-info-section" v-if="asset">
        <div class="asset-preview">
          <img :src="getImageURL(asset.imageName)" :alt="asset.name" class="preview-image" />
          <div class="preview-info">
            <h4>{{ asset.name }}</h4>
            <p class="asset-id">ID: {{ asset.id }}</p>
          </div>
        </div>
        
        <!-- 状态显示 -->
        <div class="status-section">
          <h4>当前状态</h4>
          <a-tag :color="getStatusColor(status)" class="status-tag">
            {{ getStatusText(status) }}
          </a-tag>
        </div>

        <!-- 操作选择 -->
        <div class="action-section" v-if="status === 'not_trading'">
          <h4>选择交易方式</h4>
          <a-radio-group v-model:value="selectedAction" class="action-radio">
            <a-radio value="listing">普通出售</a-radio>
            <a-radio value="auction">拍卖</a-radio>
          </a-radio-group>
        </div>

        <!-- 普通出售表单 -->
        <div class="form-section" v-if="selectedAction === 'listing' && status === 'not_trading'">
          <h4>普通出售信息</h4>
          <a-form :model="listingForm" :rules="listingRules" ref="listingFormRef" layout="vertical">
            <a-form-item label="标题" name="title">
              <a-input v-model:value="listingForm.title" placeholder="请输入出售标题" />
            </a-form-item>
            <a-form-item label="价格" name="price">
              <a-input-number 
                v-model:value="listingForm.price" 
                :min="0" 
                :precision="0"
                placeholder="请输入出售价格"
                style="width: 100%"
              />
            </a-form-item>
            <a-form-item label="截止时间" name="deadline">
              <a-date-picker 
                v-model:value="listingForm.deadline" 
                show-time 
                format="YYYY-MM-DD HH:mm"
                placeholder="选择截止时间"
                style="width: 100%"
              />
            </a-form-item>
          </a-form>
        </div>

        <!-- 拍卖表单 -->
        <div class="form-section" v-if="selectedAction === 'auction' && status === 'not_trading'">
          <h4>拍卖信息</h4>
          <a-form :model="auctionForm" :rules="auctionRules" ref="auctionFormRef" layout="vertical">
            <a-form-item label="标题" name="title">
              <a-input v-model:value="auctionForm.title" placeholder="请输入拍卖标题" />
            </a-form-item>
            <a-form-item label="起拍价" name="reservePrice">
              <a-input-number 
                v-model:value="auctionForm.reservePrice" 
                :min="0" 
                :precision="0"
                placeholder="请输入起拍价"
                style="width: 100%"
              />
            </a-form-item>
            <a-form-item label="开始时间" name="startTime">
              <a-date-picker 
                v-model:value="auctionForm.startTime" 
                show-time 
                format="YYYY-MM-DD HH:mm"
                placeholder="选择开始时间"
                style="width: 100%"
              />
            </a-form-item>
            <a-form-item label="结束时间" name="deadline">
              <a-date-picker 
                v-model:value="auctionForm.deadline" 
                show-time 
                format="YYYY-MM-DD HH:mm"
                placeholder="选择结束时间"
                style="width: 100%"
              />
            </a-form-item>
          </a-form>
        </div>

        <!-- 操作按钮 -->
        <div class="dialog-actions">
          <a-button @click="handleCancel">取消</a-button>
          <a-button 
            type="primary" 
            @click="handleSubmit"
            :loading="submitting"
            :disabled="status !== 'not_trading'"
          >
            确定
          </a-button>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { getImageURL } from '../api/index';
import dayjs from 'dayjs';

interface Asset {
  id: string;
  name: string;
  description: string;
  imageName: string;
  authorId: number;
  ownerId: number;
  timeStamp: string;
}

type AssetStatus = 'not_trading' | 'listing' | 'auction';

const props = defineProps<{
  visible: boolean;
  asset: Asset | null;
  status: AssetStatus;
  submitting?: boolean;
}>();

const emit = defineEmits<{
  'update:visible': [visible: boolean];
  submit: [data: { action: string; formData: any }];
}>();

// 表单引用
const listingFormRef = ref();
const auctionFormRef = ref();

// 选中的操作
const selectedAction = ref<'listing' | 'auction'>('listing');

// 普通出售表单
const listingForm = ref({
  title: '',
  price: 0,
  deadline: null as any
});

// 拍卖表单
const auctionForm = ref({
  title: '',
  reservePrice: 0,
  startTime: null as any,
  deadline: null as any
});

// 表单验证规则
const listingRules = {
  title: [{ required: true, message: '请输入出售标题', trigger: 'blur' }],
  price: [{ required: true, message: '请输入出售价格', trigger: 'blur' }],
  deadline: [{ required: true, message: '请选择截止时间', trigger: 'change' }]
} as any;

const auctionRules = {
  title: [{ required: true, message: '请输入拍卖标题', trigger: 'blur' }],
  reservePrice: [{ required: true, message: '请输入起拍价', trigger: 'blur' }],
  startTime: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  deadline: [{ required: true, message: '请选择结束时间', trigger: 'change' }]
} as any;

// 获取状态颜色
const getStatusColor = (status: AssetStatus) => {
  const colorMap = {
    'not_trading': 'green',
    'listing': 'blue',
    'auction': 'orange'
  };
  return colorMap[status];
};

// 获取状态文本
const getStatusText = (status: AssetStatus) => {
  const textMap = {
    'not_trading': '不在交易中',
    'listing': '普通交易中',
    'auction': '拍卖中'
  };
  return textMap[status];
};

// 处理取消
const handleCancel = () => {
  emit('update:visible', false);
};

// 处理提交
const handleSubmit = async () => {
  if (props.status !== 'not_trading') {
    return;
  }

  try {
    if (selectedAction.value === 'listing') {
      // 验证普通出售表单
      await listingFormRef.value?.validate();
      
      const listingData = {
        assetId: props.asset?.id,
        title: listingForm.value.title,
        price: listingForm.value.price,
        deadline: listingForm.value.deadline ? dayjs(listingForm.value.deadline).format() : null
      };

      emit('submit', { action: 'listing', formData: listingData });
    } else if (selectedAction.value === 'auction') {
      // 验证拍卖表单
      await auctionFormRef.value?.validate();
      
      const auctionData = {
        assetId: props.asset?.id,
        title: auctionForm.value.title,
        reservePrice: auctionForm.value.reservePrice,
        startTime: auctionForm.value.startTime ? dayjs(auctionForm.value.startTime).format() : null,
        deadline: auctionForm.value.deadline ? dayjs(auctionForm.value.deadline).format() : null
      };

      emit('submit', { action: 'auction', formData: auctionData });
    }
  } catch (error) {
    console.error('表单验证失败:', error);
  }
};

// 监听对话框关闭，重置表单
watch(() => props.visible, (newVisible) => {
  if (!newVisible) {
    // 重置表单
    listingForm.value = {
      title: '',
      price: 0,
      deadline: null
    };
    auctionForm.value = {
      title: '',
      reservePrice: 0,
      startTime: null,
      deadline: null
    };
    selectedAction.value = 'listing';
  }
});
</script>

<style scoped>
.asset-status-dialog {
  padding: 20px 0;
}

.asset-info-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.asset-preview {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.preview-image {
  width: 80px;
  height: 80px;
  object-fit: cover;
  border-radius: 8px;
}

.preview-info h4 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 16px;
}

.preview-info .asset-id {
  margin: 0;
  color: #666;
  font-size: 12px;
}

.status-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-section h4 {
  margin: 0;
  color: #333;
  font-size: 14px;
}

.status-tag {
  font-size: 12px;
  padding: 4px 8px;
}

.action-section {
  padding: 16px 0;
  border-top: 1px solid #e9ecef;
}

.action-section h4 {
  margin: 0 0 12px 0;
  color: #333;
  font-size: 14px;
}

.action-radio {
  width: 100%;
}

.form-section {
  padding: 16px 0;
  border-top: 1px solid #e9ecef;
}

.form-section h4 {
  margin: 0 0 16px 0;
  color: #333;
  font-size: 14px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid #e9ecef;
}
</style>
