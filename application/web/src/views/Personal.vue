<template>
    <div class="product-list">
      <ProductCard :product="product1" />
      <ProductCard :product="product2" />
    </div>
    <Uploader />
  </template>
  
  <script setup lang="ts">
  import ProductCard from '../components/ProductCard.vue';
  import Uploader from '../components/Uploader.vue';
  
  const product1 = {
    image: 'https://s1.imagehub.cc/images/2025/09/03/c20b2ba72059820571a38d967daeca67.png',
    name: '智能手表 Pro',
    author: '科技公司',
    created_at: '2023-04-20',
    price: 899,
  };
  
  const product2 = {
    image: 'https://s1.imagehub.cc/images/2025/03/02/5ca6e706edd485a936961c5e8c1b6bba.webp',
    name: '创意台灯',
    author: '设计工作室',
    created_at: '2023-03-15',
  };
  </script>
  
  <style>
  .product-list {
    display: flex;
    gap: 20px;
    flex-wrap: wrap;
    justify-content: center;
    padding: 20px;
  }
  </style>