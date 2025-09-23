<template>
  <div class="main-container register-container">
    <div class="svg-top">
      <svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" height="1337" width="1337">
        <defs>
          <path id="path-1" opacity="1" fill-rule="evenodd" d="M1337,668.5 C1337,1037.455193874239 1037.455193874239,1337 668.5,1337 C523.6725684305388,1337 337,1236 370.50000000000006,1094 C434.03835568300906,824.6732385973953 6.906089672974592e-14,892.6277623047779 0,668.5000000000001 C0,299.5448061257611 299.5448061257609,1.1368683772161603e-13 668.4999999999999,0 C1037.455193874239,0 1337,299.544806125761 1337,668.5Z"/>
          <linearGradient id="linearGradient-2" x1="0.79" y1="0.62" x2="0.21" y2="0.86">
            <stop offset="0" stop-color="rgb(88,62,213)" stop-opacity="1"/>
            <stop offset="1" stop-color="rgb(23,215,250)" stop-opacity="1"/>
          </linearGradient>
        </defs>
        <g opacity="1">
          <use xlink:href="#path-1" fill="url(#linearGradient-2)" fill-opacity="1"/>
        </g>
      </svg>
    </div>
    <div class="svg-bottom">
      <svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" height="896" width="967.8852157128662">
        <defs>
          <path id="path-2" opacity="1" fill-rule="evenodd" d="M896,448 C1142.6325445712241,465.5747656464056 695.2579309733121,896 448,896 C200.74206902668806,896 5.684341886080802e-14,695.2579309733121 0,448.0000000000001 C0,200.74206902668806 200.74206902668791,5.684341886080802e-14 447.99999999999994,0 C695.2579309733121,0 475,418 896,448Z"/>
          <linearGradient id="linearGradient-3" x1="0.5" y1="0" x2="0.5" y2="1">
            <stop offset="0" stop-color="rgb(40,175,240)" stop-opacity="1"/>
            <stop offset="1" stop-color="rgb(18,15,196)" stop-opacity="1"/>
          </linearGradient>
        </defs>
        <g opacity="1">
          <use xlink:href="#path-2" fill="url(#linearGradient-3)" fill-opacity="1"/>
        </g>
      </svg>
    </div>
    <section class="container">
      <section class="wrapper">
        <header>
          <h1>NFT 交易系统</h1>
          <p>用户注册</p>
        </header>
        <section class="main-content">
          <form @submit.prevent="handleRegister">
            <input type="text" placeholder="用户名" v-model="username">
            <div class="line"></div>
            <input type="email" placeholder="邮箱" v-model="email">
            <div class="line"></div>
            <div v-if="!isEmailValid && email.length > 0" class="validation-text error">
              请输入有效的邮箱地址
            </div>
            <input type="password" placeholder="密码" v-model="password">
            <div class="line"></div>
            <div class="validation-text" :class="{ 'weak': passwordStrength === '弱', 'medium': passwordStrength === '中等', 'strong': passwordStrength === '强' }" v-if="password.length > 0">
              密码强度: {{ passwordStrength }}
            </div>
            <div class="line"></div>
            <input type="password" placeholder="确认密码" v-model="confirmPassword">
            <div class="validation-text" :class="{ 'error': !isPasswordMatch && confirmPassword.length > 0 }" v-if="confirmPassword.length > 0">
              {{ !isPasswordMatch ? '两次输入的密码不一致' : '密码匹配' }}
            </div>
            <div class="line"></div>
            <div class="organization-group">
              <label>组织:</label>
              <a-radio-group v-model:value="selectedOrg" @change="handleOrgChange">
                <a-radio 
                  v-for="org in organizations" 
                  :key="org.value" 
                  :value="org.value"
                >
                  {{ org.label }}
                </a-radio>
              </a-radio-group>
              <div class="validation-text" :class="{ 'error': !isOrgValid }">
                {{ !isOrgValid ? '请选择一个组织' : '' }}
              </div>
              <div v-if="orgWarning" class="warning-text">
                {{ orgWarning }}
              </div>
            </div>
            <button type="submit" :disabled="!isFormValid">注册</button>
          </form>
        </section>
        <footer>
          <p @click="toLogin">已经有账号了，去登录</p>
        </footer>
      </section>
    </section>
  </div>
</template>

<style>
@import '../assets/auth.css';

.organization-group {
  margin-bottom: 20px;
}

.organization-group label {
  display: block;
  margin-bottom: 10px;
  font-weight: 500;
  color: #333;
}

.ant-radio-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ant-radio-wrapper {
  margin-right: 0;
  margin-bottom: 8px;
}

.ant-radio-wrapper-disabled {
  opacity: 0.6;
}

.warning-text {
  color: #ff4d4f;
  font-size: 12px;
  margin-top: 5px;
  padding: 5px 10px;
  background-color: #fff2f0;
  border: 1px solid #ffccc7;
  border-radius: 4px;
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
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

.validation-text {
  font-size: 12px;
  margin-top: 5px;
  transition: color 0.3s ease;
}
</style>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { message } from 'ant-design-vue';
import router from '../router';
import { accountApi } from '../api';

// 响应式数据
const username = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const selectedOrg = ref<number>(2);
const orgWarning = ref('');

// 组织选项配置
const organizations = [
  { label: '平台运营方', value: 1 },
  { label: 'NFT创作者', value: 2 },
  { label: '金融机构', value: 3 }
];

// 邮箱格式验证
const isEmailValid = computed(() => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email.value) || email.value.length === 0;
});

// 密码强度验证
const passwordStrength = computed(() => {
  let score = 0;
  if (!password.value) return '输入密码...';

  const hasLower = /[a-z]/.test(password.value);
  const hasUpper = /[A-Z]/.test(password.value);
  const hasNumber = /[0-9]/.test(password.value);
  const hasSpecial = /[^a-zA-Z0-9]/.test(password.value);

  if (password.value.length >= 8) score++;
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
  return password.value === confirmPassword.value || confirmPassword.value.length === 0;
});

// 组织选择验证
const isOrgValid = computed(() => {
  return selectedOrg.value !== null && selectedOrg.value === 2; // 只允许选择NFT创作者
});

// 组织选择处理函数
const handleOrgChange = (e: any) => {
  const value = e.target.value;
  if (value === 1 || value === 3) {
    orgWarning.value = '不支持选择该组织，请联系管理员';
    selectedOrg.value = 2; // 清空选择
  } else {
    orgWarning.value = '';
  }
};

// 表单整体验证
const isFormValid = computed(() => {
  return username.value.length > 0 && 
         isEmailValid.value && 
         passwordStrength.value !== '弱' && 
         passwordStrength.value !== '输入密码...' &&
         isPasswordMatch.value &&
         confirmPassword.value.length > 0 &&
         isOrgValid.value;
});

// 注册处理函数
const handleRegister = async () => {
  if (isFormValid.value) {
    try {
      const response = await accountApi.register(
        username.value,
        email.value,
        password.value,
        selectedOrg.value
      );
      if (response.status === 200 && response.data.code === 200) {
        message.success('注册成功！');
        router.push('/login');
      } else {
        message.error(`${response.data.message}`);
      }
    } catch (error) {
      message.error('注册请求失败，请检查网络连接。');
      console.error(error);
    }
  } else {
    message.error('表单验证失败，无法注册。');
  }
};

// 跳转登录
const toLogin = () => {
  router.push('/login');
};
</script>
