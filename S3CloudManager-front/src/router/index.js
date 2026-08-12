import { createRouter, createWebHistory } from 'vue-router'
import Buckets from '../components/Buckets.vue'
import Bucket from '../components/Bucket.vue'

const routes = [
  {
    path: '/',
    name: 'Buckets',
    component: Buckets
  },
  {
    path: '/buckets/:bucketName',
    name: 'Bucket',
    component: Bucket,
    props: true
  },
  // 其他路径重定向到首页
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
