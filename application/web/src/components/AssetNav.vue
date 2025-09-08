<template>
  <div class="asset-nav">
    <a-menu mode="horizontal" :selected-keys="selectedKeys" @click="handleMenuClick">
      <a-menu-item key="dashboard">
        <template #icon>
          <a-icon type="home" />
        </template>
        返回首页
      </a-menu-item>
      <a-menu-item key="upload">
        <template #icon>
          <a-icon type="upload" />
        </template>
        上传NFT
      </a-menu-item>
      <a-menu-item key="search">
        <template #icon>
          <a-icon type="search" />
        </template>
        查询资产
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
  if (newPath.includes('/asset/upload')) {
    selectedKeys.value = ['upload'];
  } else if (newPath.includes('/asset/search')) {
    selectedKeys.value = ['search'];
  } else if (newPath.includes('/dashboard')) {
    selectedKeys.value = ['dashboard'];
  } else {
    selectedKeys.value = [];
  }
}, { immediate: true });

const handleMenuClick = (e: MenuInfo) => {
  switch (e.key) {
    case 'upload':
      router.push('/asset/upload');
      break;
    case 'search':
      router.push('/asset/search');
      break;
    case 'dashboard':
      router.push('/dashboard');
      break;
  }
};
</script>

<style scoped>
.asset-nav {
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 20px;
}

.asset-nav .ant-menu {
  border-bottom: none;
}
</style>
