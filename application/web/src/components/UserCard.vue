<template>
    <div class="user-card" :class="{ 'selected': isSelected }">
      <div class="user-image-container">
        <img :src="user.avatar" :alt="user.name" class="user-avatar" />
      </div>
      <div class="user-info">
        <h3 class="user-id">ID: {{ user.id }}</h3>
        <div class="user-meta">
          <p class="user-name">用户名: {{ user.name }}</p>
          <p class="user-email">邮箱: {{ user.email }}</p>
          <p class="user-org">组织: {{ user.org }}</p>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { defineProps } from 'vue';
  
  interface User {
    avatar: string;
    id: number;
    name: string;
    email: string;
    org: string;
  }
  
  defineProps<{
    user: User;
    isSelected?: boolean;
  }>();
  </script>
  
  <style scoped>
  .user-card {
    width: 260px;
    height: 340px;
    background-color: #ffffff;
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
    overflow: hidden;
    text-align: center;
    transition: transform 0.3s ease, box-shadow 0.3s ease, border 0.3s ease;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    padding: 15px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
    border: 2px solid transparent;
  }
  
  .user-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
    border: 2px solid #2962ff;
  }
  
  .user-card.selected {
    border: 2px solid #2962ff;
    box-shadow: 0 8px 30px rgba(41, 98, 255, 0.3);
    transform: translateY(-5px);
  }
  
  .user-image-container {
    width: 100%;
    padding: 10px 0;
  }
  
  .user-avatar {
    width: 100%;
    max-width: 180px;
    height: auto;
    object-fit: contain;
    border-radius: 8px;
  }

  .user-id {
    margin-top: 10px;
    font-weight: bold;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 5px;
  }
  
  .user-info {
    flex-grow: 1;
    width: 100%;
    padding: 10px 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
    text-align: center;
  }
  
  .user-name {
    font-size: 18px;
    font-weight: bold;
    color: #1a237e;
    margin: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .user-meta {
    margin-top: 5px;
    color: #444;
    font-size: 14px;
  }
  
  .user-name,
  .user-email,
  .user-org {
    font-size: 12px;
    color: #607d8b;
    line-height: 1.5;
    margin: 0;
  }
  
  .user-label {
    font-size: 14px;
    color: #607d8b;
  }
  
  .user-value {
    font-size: 18px;
    color: #2962ff;
  }
  </style>