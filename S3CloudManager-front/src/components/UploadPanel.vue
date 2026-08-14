<template>
  <div v-if="showUploadPanel" class="upload-overlay" @click.self="closeUploadPanel">
    <div class="upload-panel"
         @drop="handleDrop"
         @dragover="handleDragOver"
         @dragenter="handleDragEnter"
         @dragleave="handleDragLeave"
         :class="{ 'drag-over': isDragOver }">
      <!-- 头部 -->
      <div class="upload-header">
        <h5 class="upload-title">
          <i class="material-icons left">cloud_upload</i>文件上传
        </h5>
        <button @click="closeUploadPanel" class="btn-floating btn-small btn-secondary waves-effect waves-light">
          <i class="material-icons" style="color: var(--color-text-secondary);">close</i>
        </button>
      </div>

      <!-- 拖拽区域 -->
      <div class="upload-dropzone" v-if="!uploadFiles.length">
        <i class="material-icons upload-icon">cloud_upload</i>
        <p class="upload-text">拖拽文件到此处</p>
        <p class="upload-divider">或</p>
        <label class="btn btn-primary waves-effect waves-light upload-btn">
          <i class="material-icons left">folder_open</i>选择文件
          <input type="file" multiple @change="handleFileSelect" accept="*/*" style="display: none;">
        </label>
        <p class="upload-divider">或</p>
        <button @click="triggerClipboardPaste" class="btn btn-secondary waves-effect waves-light upload-btn">
          <i class="material-icons left">content_paste</i>从剪贴板粘贴
        </button>
        <p class="upload-hint">支持多文件上传，单个文件最大 100MB</p>
        <p class="upload-hint clipboard-hint">
          <i class="material-icons">keyboard</i>
          也可使用 <kbd>Ctrl</kbd>+<kbd>V</kbd> 快速粘贴剪贴板图片
        </p>
      </div>

      <!-- 已选文件列表 -->
      <div v-if="uploadFiles.length" class="upload-content">
        <div class="upload-files-header">
          <span class="upload-files-count">已选择 {{ uploadFiles.length }} 个文件</span>
          <div class="upload-files-header-actions">
            <button @click="triggerClipboardPaste" class="btn-flat btn-small waves-effect waves-light upload-add-more">
              <i class="material-icons left" style="font-size: 16px;">content_paste</i>粘贴图片
            </button>
            <label class="btn-flat btn-small waves-effect waves-light upload-add-more">
              <i class="material-icons left" style="font-size: 16px;">add</i>添加更多
              <input type="file" multiple @change="handleFileSelect" accept="*/*" style="display: none;">
            </label>
          </div>
        </div>
        <div class="upload-files-list">
          <div v-for="(file, index) in uploadFiles" :key="index" class="upload-file-item">
            <!-- 图片缩略图 -->
            <div v-if="isImageFile(file.name)" class="file-thumbnail">
              <img :src="getFilePreview(file)" alt="" @error="handleFilePreviewError($event)" />
            </div>
            <!-- 文件图标 -->
            <div v-else class="file-icon">
              <i class="material-icons">{{ getFileIcon(file.name) }}</i>
            </div>
            <!-- 文件信息 -->
            <div class="file-info">
              <span class="file-name" :title="file.name">{{ file.name }}</span>
              <span class="file-size">{{ formatFileSize(file.size) }}</span>
            </div>
            <!-- 删除按钮 -->
            <button @click="removeFile(index)" class="btn-floating btn-small btn-remove" title="移除">
              <i class="material-icons" style="font-size: 16px; color: var(--color-text-secondary);">close</i>
            </button>
          </div>
        </div>
      </div>

      <!-- 拖拽覆盖提示 -->
      <div v-if="isDragOver" class="upload-dragover">
        <i class="material-icons">cloud_upload</i>
        <p>释放鼠标上传文件</p>
      </div>

      <!-- 选项区域 -->
      <div v-if="uploadFiles.length" class="upload-options">
        <label class="convert-webp-label">
          <input type="checkbox" v-model="convertToWebp" class="convert-webp-checkbox" />
          <span class="convert-webp-text">图片转 WebP 后上传</span>
        </label>
        <div v-if="convertToWebp" class="webp-quality-slider">
          <span class="webp-quality-label">压缩质量</span>
          <input type="range" min="0.1" max="1" step="0.05" v-model.number="webpQuality" class="webp-quality-range" />
          <span class="webp-quality-value">{{ Math.round(webpQuality * 100) }}%</span>
        </div>
        <label class="convert-webp-label" style="margin-top: 8px;">
          <input type="checkbox" v-model="renameByTime" class="convert-webp-checkbox" />
          <span class="convert-webp-text">使用时间重命名文件</span>
        </label>
        <p v-if="renameByTime" class="rename-hint">格式：年月日_时分秒_随机字符.扩展名</p>
      </div>

      <!-- 操作按钮 -->
      <div v-if="uploadFiles.length" class="upload-actions">
        <button @click="uploadSelectedFiles" :disabled="isUploading" class="btn btn-primary waves-effect waves-light">
          <i class="material-icons left">cloud_upload</i>
          {{ isUploading ? '正在上传...' : `上传 ${uploadFiles.length} 个文件` }}
        </button>
        <button @click="clearUploadFiles" class="btn btn-secondary waves-effect waves-light">
          <i class="material-icons left">clear</i>清空选择
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { API_ENDPOINTS, getHeaders } from '../config/api.js';
import { isImageFile, formatFileSize } from '../utils/file-utils.js';

export default {
  name: 'UploadPanel',
  props: {
    bucketName: {
      type: String,
      required: true
    },
    currentPath: {
      type: String,
      default: ''
    },
    showUploadPanel: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      uploadFiles: [],
      isUploading: false,
      isDragOver: false,
      convertToWebp: false,
      webpQuality: 0.85,
      renameByTime: false
    }
  },
  watch: {
    showUploadPanel(val) {
      if (val) {
        document.addEventListener('paste', this.handlePaste)
      } else {
        document.removeEventListener('paste', this.handlePaste)
      }
    }
  },
  beforeUnmount() {
    document.removeEventListener('paste', this.handlePaste)
  },
  methods: {
    // 点击按钮触发剪贴板粘贴（模拟 Ctrl+V）
    async triggerClipboardPaste() {
      const fakeEvent = new ClipboardEvent('paste', {
        bubbles: true,
        cancelable: true,
        clipboardData: new DataTransfer()
      })

      // 尝试从 Clipboard API 获取图片并设置到 clipboardData
      try {
        if (navigator.clipboard && navigator.clipboard.read) {
          const items = await navigator.clipboard.read()
          for (const item of items) {
            const imageType = item.types.find(t => t.startsWith('image/'))
            if (imageType) {
              const blob = await item.getType(imageType)
              const ext = imageType.split('/')[1] || 'png'
              const fileName = `clipboard_${Date.now()}.${ext}`
              const file = new File([blob], fileName, { type: imageType })
              fakeEvent.clipboardData.items.add(file)
            }
          }
        }
      } catch (e) {
        // Clipboard API 不可用，降级使用 execCommand 触发粘贴
        const pasteTarget = document.createElement('textarea')
        pasteTarget.style.position = 'fixed'
        pasteTarget.style.opacity = '0'
        document.body.appendChild(pasteTarget)
        pasteTarget.focus()
        document.execCommand('paste')
        document.body.removeChild(pasteTarget)
        return
      }

      // 如果有图片数据，触发 handlePaste
      if (fakeEvent.clipboardData.items.length > 0) {
        this.handlePaste(fakeEvent)
      } else {
        this.$root.notify('剪贴板中没有图片', 'info')
      }
    },

    // 处理粘贴事件（剪贴板图片上传）
    handlePaste(event) {
      const items = event.clipboardData?.items
      if (!items) return

      const imageFiles = []
      for (const item of items) {
        if (item.type.startsWith('image/')) {
          const file = item.getAsFile()
          if (file) {
            // 为剪贴板图片生成文件名
            const ext = item.type.split('/')[1] || 'png'
            const fileName = `clipboard_${Date.now()}.${ext}`
            const namedFile = new File([file], fileName, { type: item.type })
            imageFiles.push(namedFile)
          }
        }
      }

      if (imageFiles.length > 0) {
        event.preventDefault()
        this.uploadFiles = [...this.uploadFiles, ...imageFiles]
        this.$root.notify(`已从剪贴板添加 ${imageFiles.length} 张图片`, 'success')
      }
    },
    // 处理文件选择
    handleFileSelect(event) {
      console.log('handleFileSelect called');
      console.log('Files selected:', event.target.files);
      this.uploadFiles = Array.from(event.target.files);
      console.log('uploadFiles array:', this.uploadFiles);
      console.log('uploadFiles length:', this.uploadFiles.length);
      // 强制更新视图
      this.$forceUpdate();
    },

    // 移除单个文件
    removeFile(index) {
      this.uploadFiles.splice(index, 1);
    },

    // 清空所有选择的文件
    clearUploadFiles() {
      this.uploadFiles = [];
    },

    // 关闭上传面板
    closeUploadPanel() {
      this.uploadFiles = [];
      this.isDragOver = false;
      this.$emit('close');
    },

    // 拖拽相关方法
    handleDragOver(event) {
      event.preventDefault();
      event.stopPropagation();
    },

    handleDragEnter(event) {
      event.preventDefault();
      event.stopPropagation();
      this.isDragOver = true;
    },

    handleDragLeave(event) {
      event.preventDefault();
      event.stopPropagation();
      // 只有当离开整个拖拽区域时才设置为false
      if (!event.currentTarget.contains(event.relatedTarget)) {
        this.isDragOver = false;
      }
    },

    handleDrop(event) {
      event.preventDefault();
      event.stopPropagation();
      this.isDragOver = false;

      const files = Array.from(event.dataTransfer.files);
      if (files.length > 0) {
        console.log('拖拽文件:', files);
        // 将拖拽的文件添加到现有文件列表中
        this.uploadFiles = [...this.uploadFiles, ...files];
        this.$root.notify(`已添加 ${files.length} 个文件`, 'success');
      }
    },

    // 文件上传功能
    async uploadSelectedFiles() {
      if (!this.uploadFiles.length) {
        this.$root.notify('Please select files to upload', 'error');
        return;
      }

      this.isUploading = true;

      try {
        // 如果开启了 WebP 转换，先转换图片文件
        let filesToUpload = this.uploadFiles;
        let convertedCount = 0;
        if (this.convertToWebp) {
          this.$root.notify('正在转换图片格式...', 'info');
          const results = await Promise.all(
            this.uploadFiles.map(async (file) => {
              if (isImageFile(file.name)) {
                try {
                  const converted = await this.convertImageToWebp(file, this.webpQuality);
                  if (converted.name !== file.name) convertedCount++;
                  return converted;
                } catch (err) {
                  console.warn(`转换 ${file.name} 失败，使用原文件:`, err);
                  return file;
                }
              }
              return file;
            })
          );
          filesToUpload = results;
        }

        // 如果开启了时间重命名，生成新文件名
        if (this.renameByTime) {
          filesToUpload = filesToUpload.map((file) => {
            const newName = this.getTimestampFileName(file);
            return new File([file], newName, {
              type: file.type,
              lastModified: file.lastModified,
            });
          });
        }

        const uploadPromises = filesToUpload.map(async (file) => {
          const formData = new FormData();
          formData.append('file', file);

          // 如果有当前路径，将文件上传到当前文件夹
          if (this.currentPath) {
            formData.append('object_name', `${this.currentPath}${file.name}`);
          }

          const response = await fetch(API_ENDPOINTS.uploadObject(this.bucketName), {
            method: 'POST',
            headers: getHeaders(),
            body: formData
          });

          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }

          return response.json();
        });

        await Promise.all(uploadPromises);

        const parts = [];
        if (this.convertToWebp && convertedCount > 0) {
          parts.push(`转换 ${convertedCount} 个图片为 WebP`);
        }
        if (this.renameByTime) {
          parts.push(`已使用时间重命名 ${filesToUpload.length} 个文件`);
        }
        const msg = parts.length > 0 ? `上传成功！${parts.join('，')}` : '上传成功！';
        this.$root.notify(msg, 'success');
        this.uploadFiles = [];
        this.$emit('upload-success');
        this.closeUploadPanel();
      } catch (error) {
        console.error('Error uploading files:', error);
        this.$root.notify('Failed to upload files', 'error');
      } finally {
        this.isUploading = false;
      }
    },

    // 获取文件预览 URL
    getFilePreview(file) {
      if (file.preview) return file.preview;
      if (file instanceof File) {
        try {
          file.preview = URL.createObjectURL(file);
          return file.preview;
        } catch (e) {
          return '';
        }
      }
      return '';
    },

    // 获取文件图标
    getFileIcon(filename) {
      const ext = filename.split('.').pop()?.toLowerCase();
      const iconMap = {
        'jpg': 'image', 'jpeg': 'image', 'png': 'image', 'gif': 'image',
        'bmp': 'image', 'svg': 'image', 'webp': 'image',
        'mp4': 'video_library', 'avi': 'video_library', 'mov': 'video_library',
        'wmv': 'video_library', 'flv': 'video_library', 'webm': 'video_library',
        'mp3': 'audiotrack', 'wav': 'audiotrack', 'flac': 'audiotrack',
        'aac': 'audiotrack', 'ogg': 'audiotrack',
        'pdf': 'picture_as_pdf',
        'doc': 'description', 'docx': 'description',
        'xls': 'table_chart', 'xlsx': 'table_chart',
        'ppt': 'slideshow', 'pptx': 'slideshow',
        'txt': 'article',
        'zip': 'archive', 'rar': 'archive', '7z': 'archive',
        'tar': 'archive', 'gz': 'archive'
      };
      return iconMap[ext] || 'insert_drive_file';
    },

    // 文件预览图加载失败
    handleFilePreviewError(event) {
      event.target.style.display = 'none';
      event.target.parentElement.innerHTML = '<i class="material-icons">broken_image</i>';
    },

    // 将图片文件转换为 WebP 格式
    async convertImageToWebp(file, quality) {
      // 如果本身就是 webp，跳过转换
      if (file.name.toLowerCase().endsWith('.webp')) {
        return file;
      }

      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = (e) => {
          const img = new Image();
          img.onload = () => {
            const canvas = document.createElement('canvas');
            canvas.width = img.naturalWidth;
            canvas.height = img.naturalHeight;
            const ctx = canvas.getContext('2d');
            ctx.drawImage(img, 0, 0);

            canvas.toBlob(
              (blob) => {
                if (blob) {
                  // 生成新的文件名（将扩展名改为 .webp）
                  const newName = file.name.replace(/\.[^.]+$/, '.webp');
                  const webpFile = new File([blob], newName, {
                    type: 'image/webp',
                    lastModified: file.lastModified,
                  });
                  resolve(webpFile);
                } else {
                  reject(new Error('WebP 转换失败'));
                }
              },
              'image/webp',
              quality
            );
          };
          img.onerror = () => reject(new Error('图片加载失败'));
          img.src = e.target.result;
        };
        reader.onerror = () => reject(new Error('文件读取失败'));
        reader.readAsDataURL(file);
      });
    },

    // 根据时间戳生成新文件名
    getTimestampFileName(file) {
      const now = new Date();
      const year = now.getFullYear();
      const month = String(now.getMonth() + 1).padStart(2, '0');
      const day = String(now.getDate()).padStart(2, '0');
      const hours = String(now.getHours()).padStart(2, '0');
      const minutes = String(now.getMinutes()).padStart(2, '0');
      const seconds = String(now.getSeconds()).padStart(2, '0');
      const random = Math.random().toString(36).substring(2, 8); // 6位随机字符
      const ext = file.name.split('.').pop();
      return `${year}${month}${day}_${hours}${minutes}${seconds}_${random}.${ext}`;
    },

    isImageFile,
    formatFileSize
  }
}
</script>

<style scoped>
/* ========== 上传文件浮层 ========== */
.upload-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(15, 23, 42, 0.5);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.upload-panel {
  background: var(--color-surface);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15), 0 0 0 1px rgba(0, 0, 0, 0.05);
  width: 560px;
  max-width: 90vw;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: slideUp 0.25s ease;
  transition: border-color 0.2s;
}

.upload-panel.drag-over {
  border-color: var(--color-primary);
}

/* 头部 */
.upload-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--color-border);
  flex-shrink: 0;
}

.upload-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text);
}

.upload-title i {
  font-size: 22px;
  color: var(--color-primary);
}

/* 拖拽区域 */
.upload-dropzone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
  flex: 1;
  min-height: 240px;
}

.upload-dropzone .upload-icon {
  font-size: 56px;
  color: var(--color-primary);
  opacity: 0.6;
  margin-bottom: 16px;
}

.upload-dropzone .upload-text {
  font-size: var(--font-size-md);
  color: var(--color-text);
  margin: 0 0 8px;
}

.upload-dropzone .upload-divider {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0 0 16px;
}

.upload-dropzone .upload-hint {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  margin: 16px 0 0;
}

.clipboard-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 12px !important;
}

.clipboard-hint i {
  font-size: 16px;
}

.clipboard-hint kbd {
  display: inline-block;
  padding: 2px 6px;
  font-size: 11px;
  font-family: inherit;
  color: var(--color-text);
  background: #f1f5f9;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  box-shadow: 0 1px 0 rgba(0, 0, 0, 0.1);
}

.upload-dropzone .upload-btn {
  cursor: pointer;
  padding: 0 20px;
  height: 36px;
  line-height: 36px;
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
}

/* 已选文件内容区 */
.upload-content {
  flex: 1;
  overflow-y: auto;
  min-height: 200px;
  max-height: 50vh;
}

.upload-files-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 0;
  background: var(--color-surface);
  z-index: 1;
}

.upload-files-header-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.upload-files-count {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.upload-add-more {
  font-size: var(--font-size-sm);
  color: var(--color-primary);
  cursor: pointer;
  padding: 4px 8px;
  height: auto;
  line-height: normal;
}

.upload-add-more:hover {
  background: var(--color-primary-light);
  border-radius: var(--radius-sm);
}

.upload-files-list {
  padding: 8px 20px;
}

/* 文件项 */
.upload-file-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid #f1f5f9;
}

.upload-file-item:last-child {
  border-bottom: none;
}

.file-thumbnail {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md);
  overflow: hidden;
  flex-shrink: 0;
  background: #f1f5f9;
}

.file-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.file-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md);
  background: var(--color-primary-light);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.file-icon i {
  font-size: 22px;
  color: var(--color-primary);
}

.file-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.file-name {
  font-size: var(--font-size-sm);
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-size {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}

.btn-remove {
  background: transparent;
  box-shadow: none;
  width: 28px;
  height: 28px;
  flex-shrink: 0;
}

.btn-remove:hover {
  background: #fee2e2;
}

.btn-remove:hover i {
  color: var(--color-danger) !important;
}

/* 拖拽覆盖层 */
.upload-dragover {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(37, 99, 235, 0.08);
  border: 2px dashed var(--color-primary);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 10;
  pointer-events: none;
}

.upload-dragover i {
  font-size: 48px;
  color: var(--color-primary);
  margin-bottom: 12px;
}

.upload-dragover p {
  font-size: var(--font-size-md);
  color: var(--color-primary);
  font-weight: 500;
  margin: 0;
}

/* 操作按钮 */
.upload-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 20px;
  border-top: 1px solid var(--color-border);
  flex-shrink: 0;
}

.upload-actions .btn {
  font-size: var(--font-size-sm);
  padding: 0 16px;
  height: 36px;
  line-height: 36px;
  border-radius: var(--radius-md);
}

.upload-actions .btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 上传选项区域 */
.upload-options {
  padding: 12px 20px;
  border-top: 1px solid var(--color-border);
  background: #f8fafc;
  flex-shrink: 0;
}

.convert-webp-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.convert-webp-checkbox {
  width: 16px;
  height: 16px;
  accent-color: var(--color-primary);
  cursor: pointer;
  flex-shrink: 0;
}

.convert-webp-text {
  font-size: var(--font-size-sm);
  color: var(--color-text);
}

.webp-quality-slider {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
  padding-left: 24px;
}

.webp-quality-label {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.webp-quality-range {
  flex: 1;
  height: 4px;
  -webkit-appearance: none;
  appearance: none;
  background: var(--color-border);
  border-radius: 2px;
  outline: none;
  cursor: pointer;
}

.webp-quality-range::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: var(--color-primary);
  cursor: pointer;
}

.webp-quality-range::-moz-range-thumb {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: var(--color-primary);
  border: none;
  cursor: pointer;
}

.webp-quality-value {
  font-size: var(--font-size-xs);
  color: var(--color-primary);
  font-weight: 500;
  min-width: 32px;
  text-align: right;
  flex-shrink: 0;
}

.rename-hint {
  margin: 4px 0 0 24px;
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}
</style>
