<template>
    <div class="wallet-nav">
      <a-menu mode="horizontal" :selected-keys="selectedKeys" @click="handleMenuClick">
        <a-menu-item key="dashboard">
          返回首页
        </a-menu-item>
        <a-menu-item key="balance">
          余额与转账
        </a-menu-item>
        <a-menu-item key="withhold">
          预扣款记录
        </a-menu-item>
        <a-menu-item key="assets">
          我的资产
        </a-menu-item>
        <a-menu-item key="mintToken" v-if="isFinancialOrg">
          铸币
        </a-menu-item>
      </a-menu>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, watch, onMounted } from 'vue';
  import { useRouter, useRoute } from 'vue-router';
  import { MenuInfo } from 'ant-design-vue/es/menu/src/interface';
  
  const router = useRouter();
  const route = useRoute();
  const selectedKeys = ref<string[]>([]);
  const isFinancialOrg = ref(false);

// 加载用户信息并判断是否为金融组织
const loadUserInfo = () => {
  const userInfoString = localStorage.getItem('userInfo');
  if (userInfoString) {
    try {
      const parsedUserInfo = JSON.parse(userInfoString);
      // 检查是否为金融组织（组织ID为3）
      isFinancialOrg.value = parsedUserInfo.org === 3;
    } catch (e) {
      console.error('解析 localStorage 中的 userInfo 失败', e);
      isFinancialOrg.value = false;
    }
  }
};

onMounted(() => {
  loadUserInfo();
});
  
  // 根据当前路由设置选中的菜单项
  watch(() => route.path, (newPath) => {
    if (newPath.includes('/wallet/balance')) {
      selectedKeys.value = ['balance'];
    } else if (newPath.includes('/wallet/withhold')) {
      selectedKeys.value = ['withhold'];
    } else if (newPath.includes('/wallet/assets')) {
      selectedKeys.value = ['assets'];
    }else if (newPath.includes('/wallet/mintToken')) {
      selectedKeys.value = ['mintToken'];
    } else if (newPath.includes('/dashboard')) {
      selectedKeys.value = ['dashboard'];
    } else {
      selectedKeys.value = [];
    }
  }, { immediate: true });
  
  const handleMenuClick = (e: MenuInfo) => {
    switch (e.key) {
      case 'balance':
        router.push('/wallet/balance');
        break;
      case 'withhold':
        router.push('/wallet/withhold');
        break;
      case 'assets':
        router.push('/wallet/assets');
        break;
      case 'mintToken':
        router.push('/wallet/mintToken'); // 导航到铸币页面
        break;
      case 'dashboard':
        router.push('/dashboard');
        break;
    }
  };
  </script>
  
  <style scoped>
  .wallet-nav {
    background: #fff;
    border-bottom: 1px solid #f0f0f0;
    margin-bottom: 20px;
  }
  
  .wallet-nav .ant-menu {
    border-bottom: none;
  }
  </style>
  