<template>
  <div v-if="copyMenu.visible" class="copy-format-overlay" @click="hideCopyMenu">
    <div class="copy-format-menu" @click.stop>
      <div class="copy-format-header">
        <span>选择复制格式</span>
      </div>
      <div class="copy-format-item" @click="copyWithFormat('url')">
        <i class="material-icons">link</i>
        <span class="copy-format-label">URL</span>
        <span class="copy-format-preview">{{ getImageRealUrl(copyMenu.objectKey) }}</span>
      </div>
      <div class="copy-format-item" @click="copyWithFormat('markdown')">
        <i class="material-icons">code</i>
        <span class="copy-format-label">Markdown</span>
        <span class="copy-format-preview">![{{ getAltText(copyMenu.filename) }}]({{ getImageRealUrl(copyMenu.objectKey) }})</span>
      </div>
      <div class="copy-format-item" @click="copyWithFormat('html')">
        <i class="material-icons">code</i>
        <span class="copy-format-label">HTML</span>
        <span class="copy-format-preview">&lt;img src=&quot;...&quot; alt=&quot;...&quot; /&gt;</span>
      </div>
      <div class="copy-format-item" @click="copyWithFormat('bbcode')">
        <i class="material-icons">code</i>
        <span class="copy-format-label">BBCode</span>
        <span class="copy-format-preview">[img]...[/img]</span>
      </div>
    </div>
  </div>
</template>

<script>
import { getImageDomain, THUMBNAIL_PARAMS } from '../config/api.js';

export default {
  name: 'CopyFormatMenu',
  props: {
    copyMenu: {
      type: Object,
      default: () => ({
        visible: false,
        objectKey: '',
        filename: ''
      })
    }
  },
  emits: ['close', 'copy'],
  methods: {
    hideCopyMenu() {
      this.$emit('close');
    },

    getAltText(filename) {
      if (!filename) return '';
      return filename.replace(/\.[^.]+$/, '');
    },

    getCopyText(objectKey, filename, format) {
      const url = this.getImageRealUrl(objectKey);
      const alt = this.getAltText(filename);
      switch (format) {
        case 'markdown':
          return `![${alt}](${url})`;
        case 'html':
          return `<img src="${url}" alt="${alt}" />`;
        case 'bbcode':
          return `[img]${url}[/img]`;
        case 'url':
        default:
          return url;
      }
    },

    async copyWithFormat(format) {
      const { objectKey, filename } = this.copyMenu;
      const text = this.getCopyText(objectKey, filename, format);
      try {
        await navigator.clipboard.writeText(text);
        this.$root.notify('已复制到剪贴板', 'success');
      } catch (error) {
        const textArea = document.createElement('textarea');
        textArea.value = text;
        document.body.appendChild(textArea);
        textArea.select();
        try {
          document.execCommand('copy');
          this.$root.notify('已复制到剪贴板', 'success');
        } catch (fallbackError) {
          this.$root.notify('复制失败，请手动复制', 'error');
        }
        document.body.removeChild(textArea);
      }
      this.$emit('copy', format);
    },

    getImageRealUrl(filename) {
      const bucketName = this.copyMenu.bucketName || '';
      return `${getImageDomain(bucketName)}/${filename}`;
    },

    getImageThumbnailUrl(filename) {
      const bucketName = this.copyMenu.bucketName || '';
      return `${getImageDomain(bucketName)}/${filename}?${THUMBNAIL_PARAMS}`;
    }
  }
}
</script>

<style scoped>
/* ========== 复制格式选择菜单 ========== */
.copy-format-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(2px);
  -webkit-backdrop-filter: blur(2px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10001;
  animation: fadeIn 0.15s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.copy-format-menu {
  background: var(--color-surface);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2), 0 0 0 1px rgba(0, 0, 0, 0.05);
  width: 360px;
  max-width: 90vw;
  overflow: hidden;
  animation: slideUp 0.2s ease;
}

.copy-format-header {
  padding: 14px 20px;
  font-size: var(--font-size-md);
  font-weight: 600;
  color: var(--color-text);
  border-bottom: 1px solid var(--color-border);
}

.copy-format-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  cursor: pointer;
  transition: background 0.15s ease;
}

.copy-format-item:hover {
  background: var(--color-primary-light);
}

.copy-format-item:not(:last-child) {
  border-bottom: 1px solid var(--color-border);
}

.copy-format-item i {
  font-size: 20px;
  color: var(--color-primary);
  flex-shrink: 0;
}

.copy-format-label {
  font-size: var(--font-size-sm);
  font-weight: 500;
  color: var(--color-text);
  min-width: 70px;
  flex-shrink: 0;
}

.copy-format-preview {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}
</style>
