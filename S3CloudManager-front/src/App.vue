<template>
  <div id="app">
    <!-- 顶部 Header -->
    <header class="app-header">
      <div class="header-container">
        <div class="header-left">
          <a href="/" class="header-logo">
            <i class="material-icons">cloud</i>
            <span>S3 Storage</span>
          </a>
        </div>
        <div class="header-right">
          <!-- 设置按钮 -->
          <SettingsModal ref="settingsModal" @config-saved="onConfigSaved" @config-cleared="onConfigCleared" />
        </div>
      </div>
    </header>

    <!-- 主内容区域 -->
    <main class="app-main">
      <router-view/>
    </main>

    <!-- 全局通知容器 -->
    <NotificationContainer ref="notificationContainer" />
  </div>
</template>

<script>
import SettingsModal from './components/SettingsModal.vue'
import NotificationContainer from './components/NotificationContainer.vue'

export default {
  name: 'App',
  components: {
    SettingsModal,
    NotificationContainer
  },
  provide() {
    return {
      notify: this.notify
    };
  },
  methods: {
    openSettings() {
      this.$refs.settingsModal.openModal()
    },
    notify(message, type, duration) {
      this.$refs.notificationContainer.show(message, type, duration);
    },
    onConfigSaved() {
      this.notify('S3 配置已保存，正在刷新...', 'success')
      setTimeout(() => {
        window.location.reload()
      }, 1000)
    },
    onConfigCleared() {
      this.notify('配置已清除，正在刷新...', 'success')
      setTimeout(() => {
        window.location.reload()
      }, 1000)
    }
  }
}
</script>

<style>
/* ========== CSS Variables / Design Tokens ========== */
:root {
  /* 颜色系统 */
  --color-primary: #2563EB;
  --color-primary-hover: #1D4ED8;
  --color-primary-light: #EFF6FF;
  --color-bg: #F8FAFC;
  --color-surface: #FFFFFF;
  --color-text: #0F172A;
  --color-text-secondary: #64748B;
  --color-border: #E2E8F0;
  --color-danger: #EF4444;
  --color-danger-hover: #DC2626;
  --color-success: #16A34A;
  --color-success-hover: #15803D;
  --color-warning: #F59E0B;

  /* 间距系统 */
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 12px;
  --spacing-lg: 16px;
  --spacing-xl: 20px;
  --spacing-2xl: 24px;
  --spacing-3xl: 32px;

  /* 圆角 */
  --radius-sm: 4px;
  --radius-md: 6px;
  --radius-lg: 8px;

  /* 字体 */
  --font-family: Inter, "PingFang SC", "Microsoft YaHei", sans-serif;
  --font-size-xs: 12px;
  --font-size-sm: 13px;
  --font-size-md: 14px;
  --font-size-lg: 16px;
  --font-size-xl: 18px;
  --font-size-2xl: 20px;

  /* 阴影 */
  --shadow-sm: 0 1px 2px rgba(15, 23, 42, 0.04);
  --shadow-md: 0 1px 3px rgba(15, 23, 42, 0.08);
  --shadow-lg: 0 4px 12px rgba(15, 23, 42, 0.08);
}

/* ========== 全局基础样式 ========== */
* {
  box-sizing: border-box;
}

body {
  margin: 0;
  padding: 0;
  font-family: var(--font-family);
  font-size: var(--font-size-md);
  color: var(--color-text);
  background-color: var(--color-bg);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* ========== Header 样式 ========== */
.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  z-index: 1000;
  display: flex;
  align-items: center;
}

.header-container {
  width: 100%;
  max-width: 1440px;
  margin: 0 auto;
  padding: 0 var(--spacing-2xl);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-logo {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  text-decoration: none;
  color: var(--color-text);
  font-size: var(--font-size-xl);
  font-weight: 600;
}

.header-logo i {
  font-size: 24px;
  color: var(--color-primary);
}

.header-right {
  position: relative;
  display: flex;
  align-items: center;
}

/* ========== 主内容区域 ========== */
.app-main {
  padding-top: 64px; /* Header height */
  min-height: 100vh;
}

/* ========== 全局按钮样式优化 ========== */
.btn,
.btn-large,
.btn-small,
.btn-floating {
  border-radius: var(--radius-md) !important;
  font-family: var(--font-family);
  text-transform: none;
  box-shadow: none !important;
}

.btn {
  height: 36px;
  line-height: 36px;
  padding: 0 var(--spacing-lg);
  font-size: var(--font-size-md);
}

.btn-large {
  height: 44px;
  line-height: 44px;
  padding: 0 var(--spacing-xl);
  font-size: var(--font-size-lg);
}

.btn-small {
  height: 30px;
  line-height: 30px;
  padding: 0 var(--spacing-md);
  font-size: var(--font-size-sm);
}

.btn-floating {
  width: 40px;
  height: 40px;
}

.btn-floating.btn-large {
  width: 56px;
  height: 56px;
}

.btn-floating.btn-small {
  width: 30px;
  height: 30px;
}

/* Primary button */
.btn[style*="background-color: #1976D2"],
.btn[style*="background-color: rgb(25, 118, 210)"] {
  background-color: var(--color-primary) !important;
}

/* ========== 全局卡片样式优化 ========== */
.card {
  border-radius: var(--radius-lg) !important;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm) !important;
}

.card:hover {
  box-shadow: var(--shadow-md) !important;
}

.card-content {
  padding: var(--spacing-lg);
}

.card-title {
  font-size: var(--font-size-lg) !important;
  font-weight: 600 !important;
}

/* ========== 全局表格样式优化 ========== */
table.striped {
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--color-border);
}

table.striped thead th {
  background-color: var(--color-bg) !important;
  color: var(--color-text-secondary) !important;
  font-weight: 600;
  font-size: var(--font-size-sm);
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--color-border);
}

table.striped tbody td {
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--color-border);
}

table.striped tbody tr:last-child td {
  border-bottom: none;
}

table.striped tbody tr:hover {
  background-color: var(--color-primary-light) !important;
}

/* ========== 全局分页样式优化 ========== */
.pagination li {
  border-radius: var(--radius-sm) !important;
}

.pagination li a {
  color: var(--color-text);
  font-size: var(--font-size-md);
}

.pagination li.active {
  background-color: var(--color-primary) !important;
}

.pagination li.active a {
  color: white !important;
}

/* ========== 全局表单样式优化 ========== */
.input-field input:focus,
.input-field textarea:focus {
  border-bottom: 1px solid var(--color-primary) !important;
  box-shadow: 0 1px 0 0 var(--color-primary) !important;
}

.input-field input:focus + label,
.input-field textarea:focus + label {
  color: var(--color-primary) !important;
}

.input-field .prefix.active {
  color: var(--color-primary) !important;
}

/* ========== 全局 Modal 优化 ========== */
.modal {
  border-radius: var(--radius-lg) !important;
  box-shadow: var(--shadow-lg) !important;
}

.modal .modal-content {
  padding: var(--spacing-2xl);
}

.modal .modal-footer {
  padding: var(--spacing-md) var(--spacing-2xl);
  border-top: 1px solid var(--color-border);
}

.modal .modal-footer .btn,
.modal .modal-footer .btn-flat {
  margin-left: var(--spacing-sm);
}

/* ========== 全局 Collection 优化 ========== */
.collection {
  border-radius: var(--radius-lg) !important;
  border: 1px solid var(--color-border);
}

.collection .collection-item {
  padding: var(--spacing-md) var(--spacing-lg);
}

/* ========== 面包屑优化 ========== */
.breadcrumb {
  font-size: var(--font-size-md);
  color: var(--color-primary);
}

.breadcrumb:hover {
  color: var(--color-primary-hover);
  text-decoration: underline;
}

.breadcrumb:before {
  color: var(--color-text-secondary) !important;
  content: '/';
}

/* ========== 通用链接颜色 ========== */
a {
  color: var(--color-primary);
}

a:hover {
  color: var(--color-primary-hover);
}

/* ========== 滚动条优化 ========== */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: var(--color-border);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--color-text-secondary);
}
</style>
