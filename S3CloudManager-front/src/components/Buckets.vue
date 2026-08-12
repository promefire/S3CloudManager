<template>
  <div class="buckets-page">
    <div class="page-header">
      <h1 class="page-title">存储桶</h1>
      <button type="button" class="btn btn-primary waves-effect waves-light modal-trigger" data-target="modal-create-bucket">
        <i class="material-icons left">add</i>创建存储桶
      </button>
    </div>

    <div class="buckets-grid">
      <router-link :to="{ name: 'Bucket', params: { bucketName: bucket.name } }" class="bucket-card" v-for="bucket in buckets" :key="bucket.name">
        <div class="bucket-card-icon">
          <i class="material-icons">folder_open</i>
        </div>
        <div class="bucket-card-info">
          <span class="bucket-card-name">{{ bucket.name }}</span>
          <span class="bucket-card-date">Created on {{ bucket.creationDate }}</span>
        </div>
      </router-link>
    </div>

    <p v-if="!buckets.length" class="empty-state">暂无存储桶，点击右上角按钮创建</p>

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

export default {
  name: 'Buckets',
  data() {
    return {
      buckets: [],
      newBucketName: ''
    }
  },
  mounted() {
    this.fetchBuckets();
    M.Modal.init(document.querySelectorAll('.modal'));
  },
  methods: {
    async fetchBuckets() {
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
        M.toast({ html: 'Failed to load buckets', classes: 'red' });
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
        
        M.toast({ html: 'Bucket created successfully!', classes: 'green' });
        this.newBucketName = '';
        const modal = M.Modal.getInstance(document.getElementById('modal-create-bucket'));
        modal.close();
        this.fetchBuckets();
      } catch (error) {
        console.error('Error creating bucket:', error);
        M.toast({ html: 'Failed to create bucket', classes: 'red' });
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
  text-decoration: none;
  transform: translateY(-1px);
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

.empty-state {
  text-align: center;
  margin-top: var(--spacing-3xl);
  color: var(--color-text-secondary);
  font-size: var(--font-size-md);
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
