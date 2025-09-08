<template>
    <div class="product-card" :class="{ 'selected': isSelected }">
      <div class="product-image-container">
        <img :src="product.image" :alt="product.name" class="product-image" />
      </div>
      <div class="product-info">
        <h3 class="product-name">{{ product.name }}</h3>
        <div class="product-meta">
          <p class="product-author">创作者: {{ product.author }}</p>
          <p class="product-created-at">铸造时间: {{ product.created_at }}</p>
        </div>
        <div v-if="product.price" class="product-price">
          <span class="price-label">价格:</span>
          <span class="price-value">{{ product.price }} ETH</span>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { defineProps } from 'vue';
  
  interface Product {
    image: string;
    name: string;
    author: string;
    created_at: string;
    price?: number;
  }
  
  defineProps<{
    product: Product;
    isSelected?: boolean;
  }>();
  </script>
  
  <style scoped>
  .product-card {
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
  
  .product-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
    border: 2px solid #2962ff;
  }
  
  .product-card.selected {
    border: 2px solid #2962ff;
    box-shadow: 0 8px 30px rgba(41, 98, 255, 0.3);
    transform: translateY(-5px);
  }
  
  .product-image-container {
    width: 100%;
    padding: 10px 0;
  }
  
  .product-image {
    width: 100%;
    max-width: 180px;
    height: auto;
    object-fit: contain;
    border-radius: 8px;
  }
  
  .product-info {
    flex-grow: 1;
    width: 100%;
    padding: 10px 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
    text-align: center;
  }
  
  .product-name {
    font-size: 18px;
    font-weight: bold;
    color: #1a237e;
    margin: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .product-meta {
    margin-top: 5px;
    color: #444;
    font-size: 14px;
  }
  
  .product-author,
  .product-created-at {
    font-size: 12px;
    color: #607d8b;
    line-height: 1.5;
    margin: 0;
  }
  
  .product-price {
    margin-top: 10px;
    font-weight: bold;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 5px;
  }
  
  .price-label {
    font-size: 14px;
    color: #607d8b;
  }
  
  .price-value {
    font-size: 18px;
    color: #2962ff;
  }
  </style>