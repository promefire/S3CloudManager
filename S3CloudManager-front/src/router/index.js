import { createRouter, createWebHistory } from 'vue-router'
import Buckets from '../components/Buckets.vue'
import Bucket from '../components/Bucket.vue'
import Login from '../components/Login.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'Buckets',
    component: Buckets,
    meta: { requiresAuth: true }
  },
  {
    path: '/buckets/:bucketName',
    name: 'Bucket',
    component: Bucket,
    props: true,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
  
  if (to.meta.requiresAuth && !isLoggedIn) {
    // 需要登录但未登录，重定向到登录页
    next('/login')
  } else if (to.path === '/login' && isLoggedIn) {
    // 已登录但访问登录页，重定向到首页
    next('/')
  } else {
    // 其他情况正常跳转
    next()
  }
})

export default router