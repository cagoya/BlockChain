<template>
    <div class="user-list-page">
      <h1 class="page-title">管理用户的组织身份</h1>
      <p class="page-subtitle">点击卡片可选中或取消选中用户</p>
      <div class="user-cards-container">
        <UserCard
          v-for="user in users"
          :key="user.id"
          :user="user"
          :isSelected="selectedUserId === user.id"
          @click="toggleSelection(user.id)"
        />
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue';
  import UserCard from '../components/UserCard.vue'; // 确保路径正确
  
  // 硬编码的用户数据
  const users = [
    {
      id: 1,
      avatar: 'http://localhost:8888/public/images/9f7b1307-95ae-4ab7-9368-ed80e0a49d0dda.jpg',
      name: '张三',
      email: 'zhangsan@example.com',
      org: '平台运营方',
    },
    {
      id: 2,
      avatar: 'http://localhost:8888/public/images/9f7b1307-95ae-4ab7-9368-ed80e0a49d0dda.jpg',
      name: '李四',
      email: 'lisi@example.com',
      org: '创作者',
    },
    {
      id: 3,
      avatar: 'http://localhost:8888/public/images/9f7b1307-95ae-4ab7-9368-ed80e0a49d0dda.jpg',
      name: '王五',
      email: 'wangwu@example.com',
      org: '金融机构',
    },
    {
      id: 4,
      avatar: 'http://localhost:8888/public/images/9f7b1307-95ae-4ab7-9368-ed80e0a49d0dda.jpg',
      name: '赵六',
      email: 'zhaoliu@example.com',
      org: '平台运营方 金融机构',
    },
    {
      id: 5,
      avatar: 'http://localhost:8888/public/images/9f7b1307-95ae-4ab7-9368-ed80e0a49d0dda.jpg',
      name: '钱七',
      email: 'qianqi@example.com',
      org: '平台运营方 金融机构',
    },
    {
      id: 6,
      avatar: 'http://localhost:8888/public/images/9f7b1307-95ae-4ab7-9368-ed80e0a49d0dda.jpg',
      name: '孙八',
      email: 'sunba@example.com',
      org: '平台运营方 创作者 金融机构',
    },
  ];
  
  // 使用 ref 存储当前选中的用户 ID
  const selectedUserId = ref<number | null>(null);
  
  // 切换选中的用户
  const toggleSelection = (userId: number) => {
    // 如果点击的是当前已选中的卡片，则取消选中
    if (selectedUserId.value === userId) {
      selectedUserId.value = null;
    } else {
      // 否则，选中新卡片
      selectedUserId.value = userId;
    }
  };
  </script>
  
  <style scoped>
  .user-list-page {
    padding: 40px 20px;
    background-color: #f0f2f5;
    min-height: 100vh;
  }
  
  .page-title {
    text-align: center;
    font-size: 2.5rem;
    color: #1a237e;
    margin-bottom: 5px;
    font-weight: 700;
  }
  
  .page-subtitle {
    text-align: center;
    font-size: 1rem;
    color: #607d8b;
    margin-bottom: 30px;
  }
  
  .user-cards-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 30px;
    justify-content: center;
    max-width: 1300px;
    margin: 0 auto;
    padding: 20px;
  }
  
  @media (max-width: 768px) {
    .user-cards-container {
      grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
      gap: 20px;
    }
  }
  </style>