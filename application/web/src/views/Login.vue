<template>
  <div class="main-container">
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
          <p>用户登录</p>
        </header>
        <section class="main-content">
          <form @submit.prevent="handleLogin">
            <input type="text" placeholder="用户名" v-model="username" autocomplete="username">
            <div class="line"></div>
            <input type="password" placeholder="密码" v-model="password" autocomplete="current-password">
            <button type="submit">登录</button>
          </form>
        </section>
        <footer>
          <p @click="toRegister">还没有账号?</p>
        </footer>
      </section>
    </section>
  </div>
</template>

<style>
@import '../assets/auth.css';
</style>

<script setup lang="ts">
import { ref } from 'vue';
import axios from 'axios';
import { message } from 'ant-design-vue';
import router from '../router';

// 响应式变量，用于存储用户输入
const username = ref('');
const password = ref('');

// 定义登录处理函数
const handleLogin = async () => {
  // 检查用户名和密码是否为空
  if (!username.value || !password.value) {
    message.warning('用户名和密码不能为空');
    return;
  }

  try {
    const response = await axios.post('http://localhost:8888/api/account/login', {
      Username: username.value,
      Password: password.value,
    });
    
    // 检查响应状态码和业务代码
    if (response.status === 200 && response.data.code === 200) {
      localStorage.setItem('userToken', response.data.data.token);
      localStorage.setItem('userInfo', JSON.stringify(response.data.data.user));
      // 设置axios的默认请求头，这样后续所有请求都会自动携带JWT
      // 这样做是无害的，因为即使后端接口不需要JWT，也不会影响请求
      axios.defaults.headers.common['Authorization'] = `Bearer ${response.data.data.token}`;
      message.success('登录成功！');
      
      // 检查路由查询参数中是否有 redirect
      const redirectPath = router.currentRoute.value.query.redirect as string;

      // 如果有 redirect 参数，跳转到该路径，否则跳转到默认页面
      if (redirectPath) {
        router.push(redirectPath);
      } else {
        router.push('/dashboard');
      }
    } else {
      message.error(`${response.data.message}`);
    }
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      message.error(`${error.response.data.message || '用户名或密码错误'}`);
    } else {
      message.error('登录请求失败，请检查网络连接。');
      console.error(error);
    }
  }
};

const toRegister = () => {
  router.push('/register');
};
</script>