<template>
  <div id="app">
    <nav v-if="isLoggedIn">
      <div class="nav-wrapper container">
        <a href="/" class="brand-logo">S3 云存储管理平台</a>
        <ul class="right">
          <li>
            <span class="user-info">
              <i class="material-icons left">account_circle</i>
              {{ username }}
            </span>
          </li>
          <li>
            <a href="#" @click="handleLogout" class="waves-effect waves-light btn" style="background-color: #1976D2; color: white; font-weight: 500; text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);">
              <i class="material-icons left" style="color: white;">logout</i>登出
            </a>
          </li>
        </ul>
      </div>
    </nav>
    <router-view/>
  </div>
</template>

<script>
/* global M */
export default {
  name: 'App',
  data() {
    return {
      isLoggedIn: false,
      username: ''
    }
  },
  methods: {
    checkAuthStatus() {
      this.isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
      this.username = localStorage.getItem('username') || ''
    },
    handleLogout() {
      localStorage.removeItem('isLoggedIn')
      localStorage.removeItem('username')
      this.isLoggedIn = false
      this.username = ''
      M.toast({ html: '已登出', classes: 'green' })
      this.$router.push('/login')
    }
  },
  mounted() {
    this.checkAuthStatus()
    // 监听存储变化
    window.addEventListener('storage', this.checkAuthStatus)
  },
  beforeUnmount() {
    window.removeEventListener('storage', this.checkAuthStatus)
  }
}
</script>

<style>
.breadcrumb:before {
    content: '/';
}
#notifications {
    top: 20px;
    right: 30px;
    position: fixed;
    z-index: 2
}

.user-info {
    color: white !important;
    font-size: 14px;
    display: flex;
    align-items: center;
    font-weight: 500;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
    background-color: rgba(255, 255, 255, 0.1);
    padding: 8px 12px;
    border-radius: 4px;
    border: 1px solid rgba(255, 255, 255, 0.2);
}

.user-info i {
    margin-right: 8px;
    color: white !important;
    font-size: 20px;
    filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
}

nav ul a {
    color: white !important;
    font-weight: 500;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

nav ul a:hover {
    background-color: rgba(255, 255, 255, 0.2) !important;
    color: white !important;
}

nav ul a i {
    color: white !important;
    margin-right: 5px;
}

/* 导航栏按钮样式 */
nav .btn {
    background-color: #1976D2 !important;
    color: white !important;
    font-weight: 500 !important;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3) !important;
    border: none !important;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2) !important;
}

nav .btn:hover {
    background-color: #1565C0 !important;
    color: white !important;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3) !important;
}

nav .btn.red {
    background-color: #f44336 !important;
}

nav .btn.red:hover {
    background-color: #d32f2f !important;
}

nav .btn i {
    color: white !important;
    font-size: 18px !important;
    margin-right: 6px !important;
}

/* 确保所有按钮图标清晰可见 */
.btn i.material-icons {
    color: white !important;
    font-size: 18px !important;
    margin-right: 6px !important;
    filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3)) !important;
}

.btn-large i.material-icons {
    font-size: 20px !important;
    margin-right: 8px !important;
}

.btn-floating i.material-icons {
    color: white !important;
    font-size: 18px !important;
}

/* 导航栏右侧按钮间距 */
nav ul.right li {
    margin-left: 10px;
}

nav ul.right li:last-child {
    margin-right: 20px;
}

/* 确保品牌logo清晰 */
.brand-logo {
    color: white !important;
    font-weight: 600 !important;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3) !important;
}

.brand-logo i {
    color: white !important;
    margin-right: 8px;
}

/* 蓝色主题样式 */
nav {
    background-color: #1976D2 !important;
}

.brand-logo {
    color: white !important;
}

.btn {
    background-color: #1976D2 !important;
    color: white !important;
    font-weight: 500;
}

.btn:hover {
    background-color: #1565C0 !important;
    color: white !important;
}

.btn:disabled {
    background-color: #ccc !important;
    color: #666 !important;
}

.btn-large {
    background-color: #1976D2 !important;
    color: white !important;
    font-weight: 500;
}

.btn-large:hover {
    background-color: #1565C0 !important;
    color: white !important;
}

.btn-floating {
    background-color: #1976D2 !important;
    color: white !important;
}

.btn-floating:hover {
    background-color: #1565C0 !important;
    color: white !important;
}

.btn-floating i {
    color: white !important;
}

.pagination li.active {
    background-color: #1976D2 !important;
}

.pagination li a:hover {
    background-color: #E3F2FD !important;
    color: #1976D2 !important;
}

.collection-item:hover {
    background-color: #E3F2FD !important;
}

.card-panel {
    border-left: 4px solid #1976D2 !important;
}

/* 文件上传区域蓝色主题 */
.file-field .btn {
    background-color: #1976D2 !important;
}

.file-field .btn:hover {
    background-color: #1565C0 !important;
}

/* 表格头部蓝色 */
table thead th {
    background-color: #E3F2FD !important;
    color: #1976D2 !important;
    font-weight: 600;
}

/* 表格内容确保可读性 */
table tbody tr {
    border-bottom: 1px solid #e0e0e0;
}

table tbody tr:hover {
    background-color: #f5f5f5;
}

/* 确保文字颜色对比度 */
table tbody td {
    color: #333 !important;
}

/* 确保所有文字都有足够对比度 */
p, span, div, label {
    color: #333 !important;
}

.grey-text {
    color: #666 !important;
}

/* 确保链接文字清晰 */
a {
    color: #1976D2 !important;
    text-decoration: none;
}

a:hover {
    color: #1565C0 !important;
    text-decoration: underline;
}

/* 确保图标文字清晰 */
.material-icons {
    color: #1976D2 !important;
}

/* 确保卡片内容清晰 */
.card-content {
    color: #333 !important;
}

.card-title {
    color: #1976D2 !important;
    font-weight: 600;
}

/* 分页文字颜色 */
.pagination li a {
    color: #1976D2 !important;
}

.pagination li.active a {
    color: white !important;
}

.pagination li.disabled a {
    color: #ccc !important;
}

/* 链接颜色 */
a {
    color: #1976D2 !important;
}

a:hover {
    color: #1565C0 !important;
}

/* 图标颜色 */
.material-icons {
    color: #1976D2 !important;
}

/* 文件夹图标保持蓝色 */
.material-icons[style*="folder"] {
    color: #1976D2 !important;
}

/* 图片文件保持绿色 */
.material-icons[style*="image"] {
    color: #4CAF50 !important;
}

/* 面包屑导航样式 */
.breadcrumb {
    color: #1976D2 !important;
    font-weight: 500;
}

.breadcrumb:hover {
    color: #1565C0 !important;
    text-decoration: underline;
}

.breadcrumb:before {
    color: #666 !important;
}

/* 确保文字对比度 */
.nav-wrapper {
    background-color: #1976D2 !important;
}

.nav-wrapper .breadcrumb {
    color: white !important;
}

.nav-wrapper .breadcrumb:hover {
    color: #E3F2FD !important;
}

.nav-wrapper .breadcrumb:before {
    color: rgba(255, 255, 255, 0.7) !important;
}



/* 卡片面板背景 */
.card-panel {
    border-left: 4px solid #1976D2 !important;
}
</style>