<template>
  <div class="buckets-page">
    <div class="page-header">
      <h1 class="page-title">存储桶</h1>
      <button v-if="hasS3Config" type="button" class="btn btn-primary waves-effect waves-light modal-trigger" data-target="modal-create-bucket">
        <i class="material-icons left">add</i>创建存储桶
      </button>
    </div>

    <div class="buckets-grid">
      <div class="bucket-card" v-for="bucket in buckets" :key="bucket.name">
        <router-link :to="{ name: 'Bucket', params: { bucketName: bucket.name } }" class="bucket-card-link">
          <div class="bucket-card-icon">
            <i class="material-icons">folder_open</i>
          </div>
          <div class="bucket-card-info">
            <span class="bucket-card-name">{{ bucket.name }}</span>
            <span class="bucket-card-date">Created on {{ bucket.creationDate }}</span>
          </div>
        </router-link>
        <div class="bucket-card-actions">
          <button class="btn-more" @click.stop="toggleMenu(bucket.name)" title="更多设置">
            <i class="material-icons">more_vert</i>
          </button>
          <div v-if="activeMenu === bucket.name" class="bucket-dropdown" @click.stop>
            <div class="dropdown-header">自定义域名</div>
            <div class="dropdown-domain-input">
              <input type="text" v-model="bucketDomains[bucket.name]" class="domain-input" placeholder="例如 img.example.com" @keyup.enter="saveBucketDomain(bucket.name); activeMenu = ''" />
              <button @click.stop="saveBucketDomain(bucket.name); activeMenu = ''" class="btn-domain-save" title="保存">
                <i class="material-icons">check</i>
              </button>
            </div>
            <div v-if="bucketDomains[bucket.name]" class="dropdown-hint">已设置: {{ bucketDomains[bucket.name] }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 未配置 S3 引导 -->
    <div v-if="!hasS3Config" class="setup-guide">
      <div class="setup-guide-icon">
        <i class="material-icons">cloud_off</i>
      </div>
      <h3 class="setup-guide-title">尚未配置 S3 连接</h3>
      <p class="setup-guide-desc">
        请先配置您的 S3 兼容存储连接信息，即可开始管理存储桶和文件。
      </p>
      <p class="setup-guide-desc small">
        支持 AWS S3、Cloudflare R2、MinIO、Backblaze B2 等所有 S3 兼容存储服务。
      </p>
      <button type="button" class="btn btn-primary waves-effect waves-light" @click="openSettings">
        <i class="material-icons left">settings</i>立即配置
      </button>
    </div>

    <!-- 已配置但无存储桶 -->
    <p v-else-if="!buckets.length" class="empty-state">暂无存储桶，点击右上角按钮创建</p>

    <!-- 创建存储桶模态框 -->
    <div id="modal-create-bucket" class="modal">
      <form id="create-bucket-form" @submit.prevent="createBucket">
        <div class="modal-content">
          <h4>创建存储桶</h4>
          <div class="row">
            <div class="input-field col m12">
              <input id="name" type="text" name="name" placeholder="My Bucket" v-model="newBucketName">
              <label for="name">名称</label>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="modal-close waves-effect waves-green btn-flat">取消</button>
          <button type="submit" class="waves-effect waves-light btn btn-primary">创建</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
/* global M */
import { API_ENDPOINTS, getHeaders } from '../config/api.js';
import { hasConfig, getBucketDomains, setBucketDomain } from '../config/storage.js';

export default {
  name: 'Buckets',
  data() {
    return {
      buckets: [],
      newBucketName: '',
      hasS3Config: true,
      bucketDomains: {},
      activeMenu: ''
    }
  },
  mounted() {
    this.checkConfig();
    this.fetchBuckets();
    this.loadBucketDomains();
    M.Modal.init(document.querySelectorAll('.modal'));
    document.addEventListener('click', this.closeMenu);
  },
  beforeUnmount() {
    document.removeEventListener('click', this.closeMenu);
  },
  methods: {
    checkConfig() {
      this.hasS3Config = hasConfig();
    },
    openSettings() {
      this.$root.$refs.settingsModal.openModal();
    },
    loadBucketDomains() {
      this.bucketDomains = getBucketDomains();
    },
    toggleMenu(bucketName) {
      this.activeMenu = this.activeMenu === bucketName ? '' : bucketName;
    },
    closeMenu() {
      this.activeMenu = '';
    },
    saveBucketDomain(bucketName) {
      const domain = (this.bucketDomains[bucketName] || '').trim();
      setBucketDomain(bucketName, domain);
      if (domain) {
        M.toast({ html: `${bucketName} 域名已设为: ${domain}`, classes: 'green' });
      } else {
        M.toast({ html: `${bucketName} 域名已清除`, classes: 'green' });
      }
    },
    async fetchBuckets() {
      // 未配置 S3 连接时不请求
      if (!this.hasS3Config) return;

      try {
        const response = await fetch(API_ENDPOINTS.buckets, {
          headers: getHeaders()
        });
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        this.buckets = data.buckets || data;
      } catch (error) {
        console.error('Error fetching buckets:', error);
        this.$root.notify('Failed to load buckets', 'error');
      }
    },
    async createBucket() {
      try {
        const response = await fetch(API_ENDPOINTS.buckets, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            ...getHeaders()
          },
          body: JSON.stringify({ name: this.newBucketName, region: "us-east-1" })
        });
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        this.$root.notify('Bucket created successfully!', 'success');
        this.newBucketName = '';
        const modal = M.Modal.getInstance(document.getElementById('modal-create-bucket'));
        modal.close();
        this.fetchBuckets();
      } catch (error) {
        console.error('Error creating bucket:', error);
        this.$root.notify('Failed to create bucket', 'error');
      }
    }
  }
}
</script>

<style scoped>
.buckets-page {
  max-width: 1440px;
  margin: 0 auto;
  padding: var(--spacing-2xl);
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-2xl);
}

.page-title {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.buckets-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-lg);
}

.bucket-card {
  position: relative;
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  padding: var(--spacing-lg);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  text-decoration: none;
  transition: all 0.15s ease;
  box-shadow: var(--shadow-sm);
}

.bucket-card:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.bucket-card-link {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  text-decoration: none;
  color: inherit;
}

.bucket-card-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  background: var(--color-primary-light);
  display: flex;
  align-items: center;
  justify-content: center;
}

.bucket-card-icon i {
  font-size: 28px;
  color: var(--color-primary);
}

.bucket-card-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.bucket-card-name {
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text);
}

.bucket-card-date {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

/* 卡片右下角更多按钮 */
.bucket-card-actions {
  position: absolute;
  bottom: 8px;
  right: 8px;
}

.btn-more {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 50%;
  background: transparent;
  cursor: pointer;
  transition: all 0.15s ease;
  opacity: 0.4;
}

.btn-more:hover {
  opacity: 1;
  background: rgba(0, 0, 0, 0.05);
}

.btn-more i {
  font-size: 18px;
  color: var(--color-text-secondary);
}

.bucket-card:hover .btn-more {
  opacity: 0.7;
}

/* 下拉菜单 */
.bucket-dropdown {
  position: absolute;
  bottom: 36px;
  right: 0;
  width: 240px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  padding: var(--spacing-md);
  z-index: 100;
  animation: dropdownIn 0.15s ease;
}

@keyframes dropdownIn {
  from {
    opacity: 0;
    transform: translateY(4px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.dropdown-header {
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: var(--spacing-sm);
}

.dropdown-domain-input {
  display: flex;
  align-items: center;
  gap: 4px;
}

.domain-input {
  flex: 1;
  height: 30px;
  padding: 0 8px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  background: #fff;
  outline: none;
  transition: border-color 0.2s;
}

.domain-input:focus {
  border-color: var(--color-primary);
}

.domain-input::placeholder {
  color: var(--color-text-secondary);
}

.btn-domain-save {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border: none;
  border-radius: var(--radius-sm);
  background: var(--color-primary);
  cursor: pointer;
  transition: all 0.15s ease;
  flex-shrink: 0;
}

.btn-domain-save:hover {
  background: var(--color-primary-hover);
}

.btn-domain-save i {
  font-size: 14px;
  color: white;
}

.dropdown-hint {
  margin-top: var(--spacing-sm);
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  word-break: break-all;
}

.empty-state {
  text-align: center;
  margin-top: var(--spacing-3xl);
  color: var(--color-text-secondary);
  font-size: var(--font-size-md);
}

/* S3 未配置引导 */
.setup-guide {
  text-align: center;
  margin-top: var(--spacing-3xl);
  padding: var(--spacing-3xl);
  background: var(--color-surface);
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-lg);
  max-width: 480px;
  margin-left: auto;
  margin-right: auto;
}

.setup-guide-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--color-primary-light);
  margin: 0 auto var(--spacing-lg);
}

.setup-guide-icon i {
  font-size: 32px;
  color: var(--color-primary);
}

.setup-guide-title {
  font-size: var(--font-size-xl);
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 var(--spacing-sm);
}

.setup-guide-desc {
  font-size: var(--font-size-md);
  color: var(--color-text-secondary);
  margin: 0 0 var(--spacing-xs);
  line-height: 1.6;
}

.setup-guide-desc.small {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  opacity: 0.8;
  margin-bottom: var(--spacing-xl);
}

/* Primary button override */
.btn-primary {
  background-color: var(--color-primary) !important;
  color: white !important;
}

.btn-primary:hover {
  background-color: var(--color-primary-hover) !important;
}

.btn-primary i {
  color: white !important;
}
</style>
