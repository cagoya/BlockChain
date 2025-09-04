<template>
    <div class="marketplace-page">
      <header class="header">
        <div class="logo">NFT Market</div>
        <div class="search-bar">
          <input type="text" placeholder="搜索 NFT、创作者或合集" class="search-input" />
          <button class="search-button">搜索</button>
        </div>
      </header>
  
      <div class="main-container">
        <aside class="sidebar">
          <div class="filter-section">
            <h4 class="filter-title">稀有度</h4>
            <ul class="filter-list">
              <li><a href="#">Normal</a></li>
              <li><a href="#">Rare</a></li>
              <li><a href="#">Super Rare</a></li>
              <li><a href="#">Super Super Rare</a></li>
              <li><a href="#">Ultra Rare</a></li>
            </ul>
          </div>
          <div class="filter-section">
            <h4 class="filter-title">价格区间</h4>
            <ul class="filter-list">
              <li><a href="#">0 - 1 ETH</a></li>
              <li><a href="#">1 - 5 ETH</a></li>
              <li><a href="#">5 ETH 以上</a></li>
            </ul>
          </div>
        </aside>
  
        <main class="product-grid">
          <ProductCard 
            v-for="(product, index) in nftList" 
            :key="index"
            :product="product"
            :isSelected="selectedNftId === index"
            @click="selectNft(index)"
          />
        </main>
      </div>
  
      <div class="pagination">
        <button class="page-button active">1</button>
        <button class="page-button">2</button>
        <button class="page-button">3</button>
        <button class="page-button">下一页</button>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue';
  import ProductCard from '../components/ProductCard.vue';
  
  // 模拟的 NFT 数据，硬编码
  const nftList = ref([
    {
      image: 'https://s1.imagehub.cc/images/2025/03/02/5ca6e706edd485a936961c5e8c1b6bba.webp',
      name: 'Cybernetic Voyager',
      author: 'Crypto Artist',
      created_at: '2025-09-03',
      price: 1.5
    },
    {
      image: 'https://s1.imagehub.cc/images/2025/03/02/5ca6e706edd485a936961c5e8c1b6bba.webp',
      name: 'Galactic Wanderer',
      author: 'Pixel Pioneer',
      created_at: '2025-08-29',
      price: 0.8
    },
    {
      image: 'https://s1.imagehub.cc/images/2025/03/02/5ca6e706edd485a936961c5e8c1b6bba.webp',
      name: 'Digital Essence #01',
      author: 'Vivid Dreams',
      created_at: '2025-08-25',
      price: 3.2
    },
    {
      image: 'https://s1.imagehub.cc/images/2025/03/02/5ca6e706edd485a936961c5e8c1b6bba.webp',
      name: 'Future Cityscape',
      author: 'Urban Scape',
      created_at: '2025-08-20',
      price: 2.1
    },
    {
      image: 'https://s1.imagehub.cc/images/2025/03/02/5ca6e706edd485a936961c5e8c1b6bba.webp',
      name: 'Abstract Flow',
      author: 'Neural Art',
      created_at: '2025-08-15',
      price: 0.5
    },
    {
      image: 'https://s1.imagehub.cc/images/2025/03/02/5ca6e706edd485a936961c5e8c1b6bba.webp',
      name: 'Meme Cat',
      author: 'Meme Lords',
      created_at: '2025-08-10',
      price: 0.2
    },
    {
      image: 'https://s1.imagehub.cc/images/2025/03/02/5ca6e706edd485a936961c5e8c1b6bba.webp',
      name: 'Cosmic Traveler',
      author: 'Star Gazer',
      created_at: '2025-08-05',
      price: 1.8
    },
    {
      image: 'https://s1.imagehub.cc/images/2025/03/02/5ca6e706edd485a936961c5e8c1b6bba.webp',
      name: '3D Golden Skull',
      author: '3D Wizard',
      created_at: '2025-08-01',
      price: 4.5
    }
  ]);
  
  const selectedNftId = ref(-1); // 默认没有选中任何卡片
  
  const selectNft = (index: number) => {
    selectedNftId.value = index;
  };
  </script>
  
  <style scoped>
  .marketplace-page {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    background-color: #f0f2f5; /* 浅灰色背景 */
    color: #1a237e;
  }
  
  .header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 15px 50px;
    background-color: #ffffff;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  }
  
  .logo {
    font-size: 28px;
    font-weight: bold;
    color: #2962ff;
    font-style: italic;
  }
  
  .search-bar {
    display: flex;
    width: 500px;
  }
  
  .search-input {
    flex-grow: 1;
    padding: 10px 15px;
    border: 2px solid #e0e0e0;
    border-radius: 20px 0 0 20px;
    outline: none;
    font-size: 16px;
    transition: border-color 0.3s ease;
  }
  
  .search-input:focus {
    border-color: #2962ff;
  }
  
  .search-button {
    padding: 10px 20px;
    background-color: #2962ff;
    color: #ffffff;
    border: none;
    border-radius: 0 20px 20px 0;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.3s ease;
  }
  
  .search-button:hover {
    background-color: #1a237e;
  }
  
  .main-container {
    display: flex;
    padding: 20px 50px;
    gap: 25px;
  }
  
  .sidebar {
    width: 250px;
    background-color: #ffffff;
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
    height: fit-content;
  }
  
  .filter-section {
    margin-bottom: 25px;
  }
  
  .filter-title {
    font-size: 18px;
    color: #1a237e;
    margin-bottom: 10px;
    border-bottom: 2px solid #e0e0e0;
    padding-bottom: 5px;
  }
  
  .filter-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  
  .filter-list li {
    margin-bottom: 8px;
  }
  
  .filter-list a {
    text-decoration: none;
    color: #607d8b;
    font-size: 15px;
    transition: color 0.3s ease;
  }
  
  .filter-list a:hover {
    color: #2962ff;
  }
  
  .product-grid {
    flex-grow: 1;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 25px;
  }
  
  .pagination {
    display: flex;
    justify-content: center;
    gap: 10px;
    padding: 30px 0;
  }
  
  .page-button {
    padding: 10px 18px;
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    background-color: #ffffff;
    color: #1a237e;
    cursor: pointer;
    transition: background-color 0.3s ease, border-color 0.3s ease;
  }
  
  .page-button:hover,
  .page-button.active {
    background-color: #2962ff;
    color: #ffffff;
    border-color: #2962ff;
  }
  </style>