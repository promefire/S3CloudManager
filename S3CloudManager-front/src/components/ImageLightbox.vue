<template>
  <div v-if="previewImage" class="image-lightbox" @click="closePreview" @keydown.esc="closePreview" tabindex="0" ref="lightbox">
    <div class="lightbox-overlay" @click.stop="closePreview"></div>
    <div class="lightbox-container" @click.stop>
      <div class="lightbox-header">
        <span class="lightbox-title">{{ previewImage.name }}</span>
        <div class="lightbox-header-actions">
          <a :href="previewImage.url" target="_blank" class="btn-floating btn-small btn-primary waves-effect waves-light" title="新窗口打开">
            <i class="material-icons">open_in_new</i>
          </a>
          <button @click.stop="$emit('show-copy-menu', getCopyMenuData())" class="btn-floating btn-small btn-primary waves-effect waves-light" style="margin-left: var(--spacing-sm);" title="复制链接">
            <i class="material-icons">content_copy</i>
          </button>
          <button @click.stop="closePreview" class="btn-floating btn-small btn-danger waves-effect waves-light" style="margin-left: var(--spacing-sm);" title="关闭">
            <i class="material-icons">close</i>
          </button>
        </div>
      </div>
      <div class="lightbox-body">
        <button class="lightbox-nav lightbox-nav-prev" @click.stop="$emit('prev')" title="上一张 (←)">
          <i class="material-icons">chevron_left</i>
        </button>
        <img :src="previewImage.url" :alt="previewImage.name" class="lightbox-image" />
        <button class="lightbox-nav lightbox-nav-next" @click.stop="$emit('next')" title="下一张 (→)">
          <i class="material-icons">chevron_right</i>
        </button>
      </div>
      <div class="lightbox-footer" v-if="previewImage.size || previewImage.date">
        <span v-if="previewImage.size">{{ previewImage.size }}</span>
        <span v-if="previewImage.date" style="margin-left: 15px;">{{ previewImage.date }}</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ImageLightbox',
  props: {
    previewImage: {
      type: Object,
      default: null
    },
    imageObjects: {
      type: Array,
      default: () => []
    },
    currentImageIndex: {
      type: Number,
      default: 0
    }
  },
  emits: ['close', 'prev', 'next', 'show-copy-menu'],
  data() {
    return {
      imageDomain: ''
    }
  },
  methods: {
    closePreview() {
      this.$emit('close');
    },
    handleKeydown(event) {
      if (event.key === 'Escape') {
        this.closePreview();
      } else if (event.key === 'ArrowLeft') {
        this.$emit('prev');
      } else if (event.key === 'ArrowRight') {
        this.$emit('next');
      }
    },
    getCopyMenuData() {
      // 从 previewImage.url 中提取域名部分（协议+主机名）
      const urlObj = new URL(this.previewImage.url);
      const baseUrl = `${urlObj.protocol}//${urlObj.host}`;
      const objectKey = this.previewImage.url.replace(baseUrl + '/', '');
      return { objectKey, filename: this.previewImage.name };
    }
  }
}
</script>

<style scoped>
/* ========== 图片预览 Lightbox ========== */
.image-lightbox {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  outline: none;
}

.lightbox-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.85);
}

.lightbox-container {
  position: relative;
  z-index: 1;
  max-width: 90vw;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.lightbox-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 15px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 8px 8px 0 0;
}

.lightbox-title {
  color: white;
  font-size: 14px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  margin-right: 15px;
}

.lightbox-header-actions {
  display: flex;
  align-items: center;
}

.lightbox-body {
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.3);
  min-height: 200px;
  position: relative;
}

.lightbox-image {
  max-width: 75vw;
  max-height: 75vh;
  object-fit: contain;
  border-radius: 4px;
}

.lightbox-footer {
  padding: 8px 15px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 0 0 8px 8px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 12px;
  text-align: center;
}

/* ========== Lightbox 导航箭头 ========== */
.lightbox-nav {
  position: fixed;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.4);
  border: none;
  cursor: pointer;
  transition: background 0.2s ease;
  z-index: 10002;
}

.lightbox-nav:hover {
  background: rgba(0, 0, 0, 0.65);
}

.lightbox-nav i {
  color: white;
  font-size: 36px;
}

.lightbox-nav-prev {
  left: 24px;
}

.lightbox-nav-next {
  right: 24px;
}
</style>
