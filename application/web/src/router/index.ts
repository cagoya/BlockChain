import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('../views/Home.vue'),
    },
    {
      path: '/login',
      component: () => import('../views/Login.vue'),
    },
    {
      path: '/register',
      component: () => import('../views/Register.vue'),
    },
    {
      path: '/dashboard',
      component: () => import('../views/DashBoard.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/realty-agency',
      component: () => import('../views/RealtyAgency.vue'),
    },
    {
      path: '/trading-platform',
      component: () => import('../views/TradingPlatform.vue'),
    },
    {
      path: '/bank',
      component: () => import('../views/Bank.vue'),
    },
    {
      path: '/:pathMatch(.*)*', // 匹配所有未匹配的路径
      component: () => import('../views/404.vue'),
    }
  ],
})

// 添加全局前置守卫，用于检查是否登录
router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth) {
    // 检查本地存储中是否存在 token
    const token = localStorage.getItem('userToken');
    
    // 如果没有 token，说明用户未登录
    if (!token) {
      // 重定向到登录页
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      });
    } else {
      // 如果有 token，允许继续跳转
      next();
    }
  } else {
    // 如果目标路由不需要认证，直接允许跳转
    next();
  }
});

export default router 