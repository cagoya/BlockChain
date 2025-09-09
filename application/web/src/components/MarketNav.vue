<template>
    <div class="market-nav">
      <a-menu mode="horizontal" :selected-keys="selectedKeys" @click="handleMenuClick">
        <a-menu-item key="dashboard">
          返回首页
        </a-menu-item>
        <a-menu-item key="sell">
          出售资产
        </a-menu-item>
        <a-menu-item key="buy">
          购买资产
        </a-menu-item>
        <a-menu-item key="auction">
          拍卖资产
        </a-menu-item>
        <a-menu-item key="ask">
          咨询卖家
        </a-menu-item>
      </a-menu>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, watch } from 'vue';
  import { useRouter, useRoute } from 'vue-router';
  import { MenuInfo } from 'ant-design-vue/es/menu/src/interface';
  
  const router = useRouter();
  const route = useRoute();
  const selectedKeys = ref<string[]>([]);
  
  // 根据当前路由设置选中的菜单项
  watch(() => route.path, (newPath) => {
    if (newPath.includes('/market/sell')) {
      selectedKeys.value = ['sell'];
    } else if (newPath.includes('/market/buy')) {
      selectedKeys.value = ['buy'];
    } else if (newPath.includes('/market/auction')) {
      selectedKeys.value = ['auction'];
    } else if (newPath.includes('/market/ask')) {
      selectedKeys.value = ['ask'];
    } else if (newPath.includes('/dashboard')) {
      selectedKeys.value = ['dashboard'];
    } else {
      selectedKeys.value = [];
    }
  }, { immediate: true });
  
  const handleMenuClick = (e: MenuInfo) => {
    switch (e.key) {
      case 'sell':
        router.push('/market/sell');
        break;
      case 'buy':
        router.push('/market/buy');
        break;
      case 'auction':
        router.push('/market/auction');
        break;
      case 'ask':
        router.push('/market/ask');
        break;
      case 'dashboard':
        router.push('/dashboard');
        break;
    }
  };
  </script>
  
  <style scoped>
  .market-nav {
    background: #fff;
    border-bottom: 1px solid #f0f0f0;
    margin-bottom: 20px;
  }
  
  .market-nav .ant-menu {
    border-bottom: none;
  }
  </style>
  