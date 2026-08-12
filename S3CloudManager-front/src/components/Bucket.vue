<template>
    <div>
      <!-- 页面工具栏 -->
      <div class="page-toolbar">
        <div class="toolbar-container">
          <div class="toolbar-left">
            <router-link to="/" class="page-title-link">
              <i class="material-icons">arrow_back</i>
              <span class="page-title">{{ bucketName }}</span>
            </router-link>
            <div class="toolbar-breadcrumbs">
              <a href="#" class="breadcrumb-link" @click="navigateTo('')">{{ bucketName }}</a>
              <a v-for="crumb in breadcrumbs.slice(1)" :key="crumb.path" href="#" class="breadcrumb-link" @click="navigateTo(crumb.path)">{{ crumb.name }}</a>
            </div>
          </div>
          <div class="toolbar-right" v-if="!objects.length">
            <button @click="deleteBucket" class="btn btn-danger waves-effect waves-light">
              <i class="material-icons left">delete</i>删除
            </button>
          </div>
        </div>
      </div>

      <div class="page-content">
          
          
          <!-- 文件上传区域 -->
        <div class="row" style="margin-bottom: 20px;">
          <div class="col s12">
            <!-- 工具栏按钮 -->
            <div class="row" style="margin-bottom: 0;">
              <div class="col s6">
                <div class="left">
                  <button @click="showUploadPanel = true" class="btn btn-primary waves-effect waves-light">
                    <i class="material-icons left">cloud_upload</i>上传文件
                  </button>
                  <a class="waves-effect waves-light btn btn-primary modal-trigger" data-target="modal-create-folder">
                    <i class="material-icons left">create_new_folder</i>新建文件夹
                  </a>
                  <button @click="toggleMultiSelect" class="btn waves-effect waves-light" :class="multiSelectMode ? 'btn-success' : 'btn-primary'">
                    <i class="material-icons left">check_box</i>{{ multiSelectMode ? '退出多选' : '多选' }}
                  </button>
                  <button @click="toggleViewMode" class="btn btn-secondary waves-effect waves-light">
                    <i class="material-icons left">{{ viewMode === 'list' ? 'grid_on' : 'view_list' }}</i>{{ viewMode === 'list' ? '缩略图' : '列表' }}
                  </button>
                </div>
              </div>
              <div class="col s6">
                <div class="right" v-if="multiSelectMode && selectedObjects.length">
                  <button @click="batchDeleteObjects" class="btn btn-danger waves-effect waves-light">
                    <i class="material-icons left">delete</i>删除选中 ({{ selectedObjects.length }})
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 上传文件浮层 -->
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
              <p class="upload-hint">支持多文件上传，单个文件最大 100MB</p>
            </div>

            <!-- 已选文件列表 -->
            <div v-if="uploadFiles.length" class="upload-content">
              <div class="upload-files-header">
                <span class="upload-files-count">已选择 {{ uploadFiles.length }} 个文件</span>
                <label class="btn-flat btn-small waves-effect waves-light upload-add-more">
                  <i class="material-icons left" style="font-size: 16px;">add</i>添加更多
                  <input type="file" multiple @change="handleFileSelect" accept="*/*" style="display: none;">
                </label>
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

        <!-- 列表视图 -->
        <table class="striped" v-if="objects.length && viewMode === 'list'">
          <thead>
            <tr>
              <th style="width: 50px;">
              </th>
              <th>名称</th>
              <th>大小</th>
              <th>最后修改</th>
              <th style="min-width:100px;">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(object, index) in objects" :key="index">
              <td>
                <label v-if="!object.IsFolder && multiSelectMode">
                  <input type="checkbox" :checked="selectedObjects.includes(object.Key)" @change="toggleObjectSelection(object.Key)" />
                  <span></span>
                </label>
              </td>
              <td :style="{ cursor: object.IsFolder || isImageFile(object.Key) ? 'pointer' : 'default' }" @click="handleObjectClick(object)">
                <i class="material-icons" :title="'Icon: ' + object.Icon" :style="getIconStyle(object.Icon) + ' !important'">{{ object.Icon }}</i> 
                <span v-if="isImageFile(object.Key)" style="color: #1976D2;">{{ object.DisplayName }}</span>
                <span v-else>{{ object.DisplayName }}</span>
              </td>
              <td>{{ object.IsFolder ? '-' : formatFileSize(object.Size) }}</td>
              <td>{{ formatDateTime(object.LastModified) }}</td>
              <td>
                <button v-if="!object.IsFolder && isImageFile(object.Key)" @click="showCopyMenu(object.Key, object.DisplayName)" class="btn-floating btn-small btn-primary waves-effect waves-light" style="margin-right: var(--xs);" title="复制图片链接">
                  <i class="material-icons">content_copy</i>
                </button>
                <button v-if="!object.IsFolder" @click="deleteSingleObject(object.Key)" class="btn-floating btn-small btn-danger waves-effect waves-light" title="删除文件">
                  <i class="material-icons">delete</i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- 缩略图视图 -->
        <div v-if="objects.length && viewMode === 'thumbnail'" class="thumbnail-grid">
          <div v-for="(object, index) in objects" :key="index" class="thumbnail-item">
            <div class="thumbnail-card" :class="{ 'folder-card': object.IsFolder, 'selected': !object.IsFolder && multiSelectMode && selectedObjects.includes(object.Key) }" @click="object.IsFolder ? handleObjectClick(object) : (multiSelectMode ? toggleObjectSelection(object.Key) : handleObjectClick(object))">
              <!-- 多选复选框 -->
              <div v-if="!object.IsFolder && multiSelectMode" class="thumbnail-checkbox" @click.stop="toggleObjectSelection(object.Key)">
                <span class="thumbnail-checkmark" :class="{ 'checked': selectedObjects.includes(object.Key) }"></span>
              </div>
              <!-- 图片缩略图 -->
              <div v-if="isImageFile(object.Key)" class="thumbnail-image-container">
                <img :src="getImageRealUrl(object.Key)" :alt="object.DisplayName" class="thumbnail-image" loading="lazy" @error="handleImageError($event)" />
                <div class="thumbnail-overlay" v-if="!multiSelectMode">
                  <i class="material-icons">zoom_in</i>
                </div>
              </div>
              <!-- 文件夹图标 -->
              <div v-else-if="object.IsFolder" class="thumbnail-icon-container folder-icon-bg">
                <i class="material-icons large" style="color: #2196F3;">folder</i>
              </div>
              <!-- 其他文件图标 -->
              <div v-else class="thumbnail-icon-container">
                <i class="material-icons large" :style="getIconStyle(object.Icon)">{{ object.Icon }}</i>
              </div>
              <!-- 文件名称 -->
              <div class="thumbnail-name" :title="object.DisplayName">{{ object.DisplayName }}</div>
            </div>
            <!-- 操作按钮 -->
            <div class="thumbnail-actions" v-if="!multiSelectMode">
              <button v-if="!object.IsFolder && isImageFile(object.Key)" @click.stop="showCopyMenu(object.Key, object.DisplayName)" class="btn-floating btn-small btn-primary waves-effect waves-light" title="复制图片链接">
                <i class="material-icons">content_copy</i>
              </button>
              <button v-if="!object.IsFolder" @click.stop="deleteSingleObject(object.Key)" class="btn-floating btn-small btn-danger waves-effect waves-light" title="删除文件">
                <i class="material-icons">delete</i>
              </button>
            </div>
          </div>
        </div>
        
        <!-- 分页控件 -->
        <div v-if="pagination.total_pages > 1" class="row" style="margin-top: 20px;">
          <div class="col s12 center">
            <ul class="pagination">
              <li :class="{ disabled: !pagination.has_previous }">
                <a href="#" @click.prevent="goToPage(pagination.current_page - 1)">
                  <i class="material-icons">chevron_left</i>
                </a>
              </li>
              <li v-for="page in getPageNumbers()" :key="page" :class="{ active: page === pagination.current_page }">
                <a href="#" @click.prevent="goToPage(page)">{{ page }}</a>
              </li>
              <li :class="{ disabled: !pagination.has_next }">
                <a href="#" @click.prevent="goToPage(pagination.current_page + 1)">
                  <i class="material-icons">chevron_right</i>
                </a>
              </li>
            </ul>
            <p class="grey-text">第 {{ pagination.current_page }} 页，共 {{ pagination.total_pages }} 页 ({{ pagination.total_count }} 个对象)</p>
          </div>
        </div>
        
        <p v-if="!objects.length" style="text-align:center;margin-top:2em;color:gray;"><strong>{{ bucketName }}/{{ currentPath }}</strong> 中暂无对象</p>
      </div>
      
      <!-- 创建文件夹模态框 -->
      <div id="modal-create-folder" class="modal">
        <form @submit.prevent="createFolder(newFolderName)">
          <div class="modal-content">
            <h4>创建新文件夹</h4>
            <br>
            <div class="row">
              <div class="input-field col m12">
                <input id="folder-name" type="text" v-model="newFolderName" required>
                <label for="folder-name">文件夹名称</label>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="modal-close waves-effect waves-green btn-flat">取消</button>
            <button type="submit" class="waves-effect waves-green btn btn-primary">创建</button>
          </div>
        </form>
      </div>

      <!-- 复制格式选择菜单 -->
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

      <!-- 图片预览 Lightbox -->
      <div v-if="previewImage" class="image-lightbox" @click="closePreview" @keydown.esc="closePreview" tabindex="0" ref="lightbox">
        <div class="lightbox-overlay" @click.stop="closePreview"></div>
        <div class="lightbox-container" @click.stop>
          <div class="lightbox-header">
            <span class="lightbox-title">{{ previewImage.name }}</span>
            <div class="lightbox-header-actions">
              <a :href="previewImage.url" target="_blank" class="btn-floating btn-small btn-primary waves-effect waves-light" title="新窗口打开">
                <i class="material-icons">open_in_new</i>
              </a>
              <button @click.stop="showCopyMenu(previewImage.url.replace(IMAGE_DOMAIN + '/', ''), previewImage.name)" class="btn-floating btn-small btn-primary waves-effect waves-light" style="margin-left: var(--spacing-sm);" title="复制链接">
                <i class="material-icons">content_copy</i>
              </button>
              <button @click.stop="closePreview" class="btn-floating btn-small btn-danger waves-effect waves-light" style="margin-left: var(--spacing-sm);" title="关闭">
                <i class="material-icons">close</i>
              </button>
            </div>
          </div>
          <div class="lightbox-body">
            <img :src="previewImage.url" :alt="previewImage.name" class="lightbox-image" />
          </div>
          <div class="lightbox-footer" v-if="previewImage.size || previewImage.date">
            <span v-if="previewImage.size">{{ previewImage.size }}</span>
            <span v-if="previewImage.date" style="margin-left: 15px;">{{ previewImage.date }}</span>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  /* global M */
import { API_ENDPOINTS, IMAGE_DOMAIN, getHeaders } from '../config/api.js';
  
export default {
  name: 'Bucket',
    props: ['bucketName'],
    data() {
      return {
        objects: [],
        currentPath: '',
        newFolderName: '',
        allowDelete: true, // or false based on your logic
        pagination: {
          current_page: 1,
          page_size: 20,
          total_count: 0,
          total_pages: 0,
          has_next: false,
          has_previous: false
        },
        breadcrumbs: [],
        selectedObjects: [],
        uploadFiles: [],
        isUploading: false,
        showUploadPanel: false,
        isDragOver: false,
        multiSelectMode: false,
        viewMode: 'thumbnail',
        previewImage: null,
        convertToWebp: false,
        webpQuality: 0.85,
        renameByTime: false,
        copyMenu: {
          visible: false,
          objectKey: '',
          filename: ''
        }
      }
    },
    mounted() {
      this.fetchObjects();
      M.Modal.init(document.querySelectorAll('.modal'));
      document.addEventListener('keydown', this.handleKeydown);
    },

    beforeUnmount() {
      document.removeEventListener('keydown', this.handleKeydown);
      document.body.style.overflow = '';
    },
    methods: {
      async fetchObjects() {
        try {
          let url;
          
          console.log('fetchObjects - currentPath:', this.currentPath);
          console.log('fetchObjects - currentPath type:', typeof this.currentPath);
          console.log('fetchObjects - currentPath length:', this.currentPath.length);
          
          // 如果当前路径不为空，使用 browse 接口来浏览文件夹
          if (this.currentPath && this.currentPath.length > 0) {
            console.log('Using browse API for folder:', this.currentPath);
            // 去除 bucket 前缀，因为 API URL 中已经包含 bucket 名
            let pathOnly = this.currentPath;
            if (pathOnly.startsWith(this.bucketName + '/')) {
              pathOnly = pathOnly.substring(this.bucketName.length + 1);
            }
            console.log('browseFolder will be called with:', this.bucketName, pathOnly);
            const params = {
              page: this.pagination.current_page,
              page_size: this.pagination.page_size
            };
            url = API_ENDPOINTS.browseFolder(this.bucketName, pathOnly, params);
          } else {
            console.log('Using objects API for root directory');
            // 如果当前路径为空，使用 objects 接口来列出根目录对象
            const params = {
              page: this.pagination.current_page,
              page_size: this.pagination.page_size,
              folder: true,
              delimiter: '/'
            };
            url = API_ENDPOINTS.bucketObjects(this.bucketName, params);
          }
          console.log('Fetching objects from:', url);
          const response = await fetch(url, { headers: getHeaders() });
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          const data = await response.json();
          console.log('API response:', data);
          console.log('Objects in response:', data.objects);
          console.log('Objects type:', typeof data.objects);
          console.log('Objects is array:', Array.isArray(data.objects));
          
          // 更新分页信息
          if (data.pagination) {
            this.pagination = data.pagination;
          }
          
          // 更新面包屑导航
          if (data.breadcrumbs) {
            this.breadcrumbs = data.breadcrumbs;
          }
          
          // API 返回的是 { objects: [...], bucket: "...", prefix: "...", pagination: {...} } 格式
          // 我们需要提取 objects 数组
          let objects = data.objects;
          console.log('Objects found:', objects ? objects.length : 0);
          
          // 处理 API 返回 null 的情况
          if (objects === null || objects === undefined) {
            console.warn('API returned null/undefined objects, treating as empty array');
            objects = [];
          }
          
          // 确保 objects 是数组
          if (!Array.isArray(objects)) {
            console.warn('API returned non-array objects:', objects);
            this.objects = [];
            return;
          }
          
          // 处理对象数据，为每个对象添加图标和显示名称
          const processedObjects = objects.map(object => {
            // 根据 API 响应，使用 name 字段而不是 Key
            // 检查是否为文件夹：以 / 结尾，或者 type 为 folder，或者包含 //
            const isFolder = object.name.endsWith('/') || object.type === 'folder' || object.name.includes('//');
            const displayName = isFolder ? object.name.replace(/\/$/, '') : object.name.split('/').pop();
            
            // 根据文件扩展名或类型设置图标
            let icon = 'insert_drive_file'; // 默认文件图标
            
            if (isFolder) {
              icon = 'folder';
            } else {
              // 根据文件扩展名设置不同的图标
              const extension = object.name.split('.').pop()?.toLowerCase();
              switch (extension) {
                case 'jpg':
                case 'jpeg':
                case 'png':
                case 'gif':
                case 'bmp':
                case 'svg':
                case 'webp':
                  icon = 'image';
                  break;
                case 'mp4':
                case 'avi':
                case 'mov':
                case 'wmv':
                case 'flv':
                case 'webm':
                  icon = 'video_library';
                  break;
                case 'mp3':
                case 'wav':
                case 'flac':
                case 'aac':
                case 'ogg':
                  icon = 'audiotrack';
                  break;
                case 'pdf':
                  icon = 'picture_as_pdf';
                  break;
                case 'doc':
                case 'docx':
                  icon = 'description';
                  break;
                case 'xls':
                case 'xlsx':
                  icon = 'table_chart';
                  break;
                case 'ppt':
                case 'pptx':
                  icon = 'slideshow';
                  break;
                case 'txt':
                  icon = 'article';
                  break;
                case 'zip':
                case 'rar':
                case '7z':
                case 'tar':
                case 'gz':
                  icon = 'archive';
                  break;
                default:
                  icon = 'insert_drive_file';
              }
            }
            
            return {
              ...object,
              IsFolder: isFolder,
              Icon: icon,
              DisplayName: displayName,
              // 保持与模板兼容的属性名
              Key: object.name,
              Size: object.size,
              LastModified: object.lastModified,
              Owner: 'Unknown' // API 中没有 owner 信息，使用默认值
            };
          });
          
          // 对对象进行排序：文件夹在前，文件按时间排序（越近的越前）
          this.objects = processedObjects.sort((a, b) => {
            // 首先按类型排序：文件夹在前，文件在后
            if (a.IsFolder && !b.IsFolder) {
              return -1; // a 是文件夹，b 是文件，a 排在前面
            }
            if (!a.IsFolder && b.IsFolder) {
              return 1; // a 是文件，b 是文件夹，b 排在前面
            }
            
            // 如果都是文件夹或都是文件，按时间排序
            if (a.IsFolder && b.IsFolder) {
              // 文件夹按名称排序（字母顺序）
              return a.DisplayName.localeCompare(b.DisplayName);
            } else {
              // 文件按最后修改时间排序（越近的越前）
              const timeA = new Date(a.LastModified || 0);
              const timeB = new Date(b.LastModified || 0);
              return timeB - timeA; // 降序排列，最新的在前
            }
          });
        } catch (error) {
          console.error('Error fetching objects:', error);
          M.toast({ html: 'Failed to load objects', classes: 'red' });
        }
      },
      async deleteBucket() {
        if (!confirm('Are you sure you want to delete this bucket?')) {
          return;
        }
        
        try {
          const response = await fetch(API_ENDPOINTS.bucket(this.bucketName), {
            method: 'DELETE',
            headers: getHeaders()
          });
          
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          
          M.toast({ html: 'Bucket deleted successfully!', classes: 'green' });
          this.$router.push('/');
        } catch (error) {
          console.error('Error deleting bucket:', error);
          M.toast({ html: 'Failed to delete bucket', classes: 'red' });
        }
      },
            navigateTo(key) {
        console.log('Navigating to:', key);
        
        // 如果 key 是空字符串，导航到根目录
        if (!key) {
          this.currentPath = '';
        } else {
          // 处理路径，确保格式正确
          let path = key;
          
          // 特殊情况：如果路径就是 bucketName + '/'，这应该被当作文件夹路径
          if (path === this.bucketName + '/') {
            console.log('Case 1: path equals bucketName + /');
            this.currentPath = path;
          } else if (path.startsWith(this.bucketName + '/')) {
            console.log('Case 2: path starts with bucketName + /');
            // 如果路径以 bucketName + '/' 开头，这是一个完整的路径
            // 直接使用这个路径，不要去掉 bucketName 部分
            this.currentPath = path.endsWith('/') ? path : path + '/';
            console.log('Setting currentPath to:', this.currentPath);
          } else {
            console.log('Case 3: path does not start with bucketName + /');
            // 如果路径不以 bucketName 开头，需要处理嵌套文件夹的情况
            if (this.currentPath && this.currentPath.length > 0) {
              // 如果当前已经在某个文件夹内，将新文件夹追加到当前路径
              const folderName = path.endsWith('/') ? path.slice(0, -1) : path;
              this.currentPath = this.currentPath + folderName + '/';
            } else {
              // 如果当前在根目录，直接使用文件夹名称
              this.currentPath = path.endsWith('/') ? path : path + '/';
            }
          }
        }
        
        console.log('Current path set to:', this.currentPath);
        console.log('Current path length:', this.currentPath.length);
        console.log('Will use browse API:', this.currentPath !== '');
        this.fetchObjects();
      },
      async deleteObject(key) {
        if (!confirm('Are you sure you want to delete this object?')) {
          return;
        }
        
        try {
          const response = await fetch(API_ENDPOINTS.bucketObject(this.bucketName, key), {
            method: 'DELETE',
            headers: getHeaders()
          });
          
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          
          M.toast({ html: 'Object deleted successfully!', classes: 'green' });
          this.fetchObjects();
        } catch (error) {
          console.error('Error deleting object:', error);
          M.toast({ html: 'Failed to delete object', classes: 'red' });
        }
      },
      getDownloadUrl(objectKey) {
        return API_ENDPOINTS.browseFile(this.bucketName, objectKey);
      },
      handleOpenDownloadLinkModal(key) {
        // Logic to open download link modal
        console.log('Opening download link modal for:', key);
      },
      
      // 检查是否为图片文件
      isImageFile(filename) {
        const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg', 'webp'];
        const extension = filename.split('.').pop()?.toLowerCase();
        return imageExtensions.includes(extension);
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

      // 获取图片的真实URL
      getImageRealUrl(filename) {
        return `${IMAGE_DOMAIN}/${filename}`;
      },
      
      // 处理对象点击事件
      handleObjectClick(object) {
        if (object.IsFolder) {
          // 如果是文件夹，导航到文件夹
          this.navigateTo(object.Key);
        } else if (this.isImageFile(object.Key)) {
          // 如果是图片文件，打开预览
          this.openPreview(object);
        }
        // 其他文件类型不处理点击事件
      },

      // 切换视图模式
      toggleViewMode() {
        this.viewMode = this.viewMode === 'list' ? 'thumbnail' : 'list';
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

      // 打开图片预览
      openPreview(object) {
        this.previewImage = {
          name: object.DisplayName,
          url: this.getImageRealUrl(object.Key),
          size: object.Size ? this.formatFileSize(object.Size) : '',
          date: object.LastModified ? this.formatDateTime(object.LastModified) : ''
        };
        document.body.style.overflow = 'hidden';
        this.$nextTick(() => {
          if (this.$refs.lightbox) {
            this.$refs.lightbox.focus();
          }
        });
      },

      // 关闭图片预览
      closePreview() {
        this.previewImage = null;
        document.body.style.overflow = '';
      },

      // 显示复制格式选择菜单
      showCopyMenu(objectKey, filename) {
        this.copyMenu = {
          visible: true,
          objectKey,
          filename
        };
      },

      // 隐藏复制格式选择菜单
      hideCopyMenu() {
        this.copyMenu.visible = false;
      },

      // 获取 alt 文本（去掉扩展名的文件名）
      getAltText(filename) {
        if (!filename) return '';
        return filename.replace(/\.[^.]+$/, '');
      },

      // 生成指定格式的复制文本
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

      // 按指定格式复制到剪贴板
      async copyWithFormat(format) {
        const { objectKey, filename } = this.copyMenu;
        const text = this.getCopyText(objectKey, filename, format);
        try {
          await navigator.clipboard.writeText(text);
          M.toast({ html: '已复制到剪贴板', classes: 'green' });
        } catch (error) {
          const textArea = document.createElement('textarea');
          textArea.value = text;
          document.body.appendChild(textArea);
          textArea.select();
          try {
            document.execCommand('copy');
            M.toast({ html: '已复制到剪贴板', classes: 'green' });
          } catch (fallbackError) {
            M.toast({ html: '复制失败，请手动复制', classes: 'red' });
          }
          document.body.removeChild(textArea);
        }
        this.hideCopyMenu();
      },

      // 图片加载失败处理
      handleImageError(event) {
        event.target.style.display = 'none';
        event.target.parentElement.innerHTML = `<i class="material-icons large" style="color: #999;">broken_image</i>`;
      },

      // 键盘事件处理
      handleKeydown(event) {
        if (event.key === 'Escape' && this.previewImage) {
          this.closePreview();
        }
      },
      
      // 获取图标样式
      getIconStyle(iconName) {
        const baseStyle = 'font-size: 20px; vertical-align: middle;';
        
        let color;
        switch (iconName) {
          case 'folder':
            color = '#2196F3'; // 蓝色
            break;
          case 'image':
            color = '#4CAF50'; // 绿色
            break;
          case 'video_library':
            color = '#F44336'; // 红色
            break;
          case 'audiotrack':
            color = '#9C27B0'; // 紫色
            break;
          case 'picture_as_pdf':
            color = '#FF5722'; // 深橙色
            break;
          case 'description':
            color = '#2196F3'; // 蓝色
            break;
          case 'table_chart':
            color = '#4CAF50'; // 绿色
            break;
          case 'slideshow':
            color = '#FF9800'; // 橙色
            break;
          case 'article':
            color = '#607D8B'; // 蓝灰色
            break;
          case 'archive':
            color = '#795548'; // 棕色
            break;
          default:
            color = '#FF9800'; // 橙色（默认文件）
        }
        
        return baseStyle + 'color: ' + color + ';';
      },
      
      async createFolder(folderName) {
        if (!folderName || !folderName.trim()) {
          M.toast({ html: 'Please enter a folder name', classes: 'red' });
          return;
        }
        
        try {
          const response = await fetch(API_ENDPOINTS.createFolder(this.bucketName), {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              ...getHeaders()
            },
            body: JSON.stringify({ folder_name: folderName.trim() })
          });
          
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          
          M.toast({ html: 'Folder created successfully!', classes: 'green' });
          this.newFolderName = '';
          const modal = M.Modal.getInstance(document.getElementById('modal-create-folder'));
          modal.close();
          this.fetchObjects();
        } catch (error) {
          console.error('Error creating folder:', error);
          M.toast({ html: 'Failed to create folder', classes: 'red' });
        }
      },
      
      // 文件上传功能
      async uploadSelectedFiles() {
        if (!this.uploadFiles.length) {
          M.toast({ html: 'Please select files to upload', classes: 'red' });
          return;
        }

        this.isUploading = true;

        try {
          // 如果开启了 WebP 转换，先转换图片文件
          let filesToUpload = this.uploadFiles;
          let convertedCount = 0;
          if (this.convertToWebp) {
            M.toast({ html: '正在转换图片格式...', classes: 'blue' });
            const results = await Promise.all(
              this.uploadFiles.map(async (file) => {
                if (this.isImageFile(file.name)) {
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
          const msg = parts.length > 0 ? `上传成功！${parts.join('，')}` : 'Files uploaded successfully!';
          M.toast({ html: msg, classes: 'green' });
          this.uploadFiles = [];
          this.showUploadPanel = false;
          this.fetchObjects();
        } catch (error) {
          console.error('Error uploading files:', error);
          M.toast({ html: 'Failed to upload files', classes: 'red' });
        } finally {
          this.isUploading = false;
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
        this.showUploadPanel = false;
        this.uploadFiles = [];
        this.isDragOver = false;
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
          M.toast({ html: `已添加 ${files.length} 个文件`, classes: 'green' });
        }
      },
      
      // 复制图片URL
      async copyImageUrl(objectKey) {
        try {
          const imageUrl = this.getImageRealUrl(objectKey);
          await navigator.clipboard.writeText(imageUrl);
          M.toast({ html: '图片链接已复制到剪贴板', classes: 'green' });
        } catch (error) {
          console.error('复制失败:', error);
          // 降级方案：使用传统方法复制
          const textArea = document.createElement('textarea');
          textArea.value = this.getImageRealUrl(objectKey);
          document.body.appendChild(textArea);
          textArea.select();
          try {
            document.execCommand('copy');
            M.toast({ html: '图片链接已复制到剪贴板', classes: 'green' });
          } catch (fallbackError) {
            console.error('降级复制也失败:', fallbackError);
            M.toast({ html: '复制失败，请手动复制', classes: 'red' });
          }
          document.body.removeChild(textArea);
        }
      },
      
      // 删除单个对象
      async deleteSingleObject(objectKey) {
        if (!confirm(`确定要删除文件 "${objectKey}" 吗？`)) {
          return;
        }
        
        try {
          const response = await fetch(API_ENDPOINTS.batchDeleteObjects(this.bucketName), {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              ...getHeaders()
            },
            body: JSON.stringify({ objects: [objectKey] })
          });
          
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          
          await response.json();
          M.toast({ html: '文件删除成功！', classes: 'green' });
          this.fetchObjects();
        } catch (error) {
          console.error('Error deleting object:', error);
          M.toast({ html: '删除文件失败', classes: 'red' });
        }
      },
      
      // 批量删除功能
      async batchDeleteObjects() {
        if (!this.selectedObjects.length) {
          M.toast({ html: 'Please select objects to delete', classes: 'red' });
          return;
        }
        
        if (!confirm(`Are you sure you want to delete ${this.selectedObjects.length} objects?`)) {
          return;
        }
        
        try {
          const response = await fetch(API_ENDPOINTS.batchDeleteObjects(this.bucketName), {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              ...getHeaders()
            },
            body: JSON.stringify({ objects: this.selectedObjects })
          });
          
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          
          const result = await response.json();
          M.toast({ html: `Deleted ${result.delete_count || result.success_count} objects successfully!`, classes: 'green' });
          this.selectedObjects = [];
          this.fetchObjects();
        } catch (error) {
          console.error('Error deleting objects:', error);
          M.toast({ html: 'Failed to delete objects', classes: 'red' });
        }
      },
      
      // 分页功能
      goToPage(page) {
        if (page >= 1 && page <= this.pagination.total_pages) {
          this.pagination.current_page = page;
          this.fetchObjects();
        }
      },
      
      // 切换对象选择
      toggleObjectSelection(objectKey) {
        const index = this.selectedObjects.indexOf(objectKey);
        if (index > -1) {
          this.selectedObjects.splice(index, 1);
        } else {
          this.selectedObjects.push(objectKey);
        }
      },
      
      // 切换多选模式
      toggleMultiSelect() {
        this.multiSelectMode = !this.multiSelectMode;
        if (!this.multiSelectMode) {
          // 退出多选模式时清空选择
          this.selectedObjects = [];
        }
      },
      
      // 获取分页数字数组
      getPageNumbers() {
        const pages = [];
        const totalPages = this.pagination.total_pages;
        const currentPage = this.pagination.current_page;
        
        // 显示最多5个页码
        let start = Math.max(1, currentPage - 2);
        let end = Math.min(totalPages, currentPage + 2);
        
        // 调整起始和结束位置，确保显示5个页码（如果可能）
        if (end - start < 4) {
          if (start === 1) {
            end = Math.min(totalPages, start + 4);
          } else {
            start = Math.max(1, end - 4);
          }
        }
        
        for (let i = start; i <= end; i++) {
          pages.push(i);
        }
        
        return pages;
      },
      
      // 格式化文件大小
      formatFileSize(bytes) {
        if (bytes === 0 || bytes === null || bytes === undefined) {
          return '0 B';
        }
        
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
      },
      
      // 格式化日期时间
      formatDateTime(dateTimeString) {
        if (!dateTimeString) {
          return '';
        }
        
        try {
          // 添加调试信息
          console.log('Original dateTimeString:', dateTimeString);
          
          const date = new Date(dateTimeString);
          
          // 检查日期是否有效
          if (isNaN(date.getTime())) {
            console.log('Invalid date, returning original string');
            return dateTimeString; // 如果解析失败，返回原始字符串
          }
          
          // 添加调试信息
          console.log('Parsed date (UTC):', date.toISOString());
          console.log('Parsed date (local):', date.toString());
          console.log('Local timezone offset:', date.getTimezoneOffset(), 'minutes');
          
          // 获取当前年份
          const currentYear = new Date().getFullYear();
          const year = date.getFullYear();
          
          // 格式化选项 - 明确指定时区处理
          const options = {
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            timeZone: 'Asia/Shanghai' // 明确指定中国时区
          };
          
          // 如果不是当前年份，添加年份显示
          if (year !== currentYear) {
            options.year = 'numeric';
          }
          
          // 使用本地化格式化
          let formattedDate;
          try {
            formattedDate = date.toLocaleDateString('zh-CN', options);
          } catch (error) {
            // 如果时区不可用，使用手动转换
            console.log('Timezone not available, using manual conversion');
            const localDate = new Date(date.getTime() + (8 * 60 * 60 * 1000)); // 手动加8小时
            formattedDate = localDate.toLocaleDateString('zh-CN', {
              month: 'long',
              day: 'numeric',
              hour: '2-digit',
              minute: '2-digit',
              year: year !== currentYear ? 'numeric' : undefined
            });
          }
          
          console.log('Formatted date:', formattedDate);
          
          return formattedDate;
        } catch (error) {
          console.error('Error formatting date:', error);
          return dateTimeString; // 出错时返回原始字符串
        }
      }
    }
  }
  </script>

  <style scoped>
  /* ========== 页面工具栏 ========== */
  .page-toolbar {
    background: var(--color-surface);
    border-bottom: 1px solid var(--color-border);
    padding: var(--spacing-md) 0;
    margin-bottom: var(--spacing-lg);
  }

  .toolbar-container {
    max-width: 1440px;
    margin: 0 auto;
    padding: 0 var(--spacing-2xl);
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .toolbar-left {
    display: flex;
    align-items: center;
    gap: var(--spacing-md);
  }

  .page-title-link {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
    text-decoration: none;
    color: var(--color-text);
    font-size: var(--font-size-xl);
    font-weight: 600;
  }

  .page-title-link i {
    font-size: 22px;
    color: var(--color-text-secondary);
  }

  .page-title {
    color: var(--color-text);
  }

  .toolbar-breadcrumbs {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
    font-size: var(--font-size-md);
  }

  .breadcrumb-link {
    color: var(--color-primary);
    text-decoration: none;
  }

  .breadcrumb-link:hover {
    color: var(--color-primary-hover);
    text-decoration: underline;
  }

  .breadcrumb-link:not(:last-child)::after {
    content: '/';
    margin-left: var(--spacing-xs);
    color: var(--color-text-secondary);
  }

  .toolbar-right {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
  }

  /* ========== 页面内容区域 ========== */
  .page-content {
    max-width: 1440px;
    margin: 0 auto;
    padding: 0 var(--spacing-2xl) var(--spacing-2xl);
  }

  /* ========== 按钮样式覆盖 ========== */
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

  .btn-secondary {
    background-color: var(--color-surface) !important;
    color: var(--color-text) !important;
    border: 1px solid var(--color-border) !important;
  }

  .btn-secondary:hover {
    background-color: var(--color-bg) !important;
  }

  .btn-secondary i {
    color: var(--color-text-secondary) !important;
  }

  .btn-danger {
    background-color: var(--color-danger) !important;
    color: white !important;
  }

  .btn-danger:hover {
    background-color: var(--color-danger-hover) !important;
  }

  .btn-danger i {
    color: white !important;
  }

  .btn-success {
    background-color: var(--color-success) !important;
    color: white !important;
  }

  .btn-success:hover {
    background-color: var(--color-success-hover) !important;
  }

  .btn-success i {
    color: white !important;
  }

  /* ========== 缩略图网格视图 ========== */
  .thumbnail-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 15px;
    padding: 10px 0;
  }

  .thumbnail-item {
    position: relative;
    width: calc(25% - 12px);
    min-width: 180px;
  }

  @media (max-width: 1200px) {
    .thumbnail-item {
      width: calc(33.333% - 10px);
    }
  }

  @media (max-width: 900px) {
    .thumbnail-item {
      width: calc(50% - 8px);
    }
  }

  @media (max-width: 600px) {
    .thumbnail-item {
      width: 100%;
    }
  }

  .thumbnail-card {
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    overflow: hidden;
    cursor: pointer;
    transition: all 0.25s ease;
    border: 1px solid #e8e8e8;
  }

  .thumbnail-card:hover {
    box-shadow: 0 6px 20px rgba(25, 118, 210, 0.2);
    border-color: #1976D2;
    transform: translateY(-2px);
  }

  .thumbnail-image-container {
    position: relative;
    width: 100%;
    height: 160px;
    overflow: hidden;
    background: #f5f5f5;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .thumbnail-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }

  .thumbnail-card:hover .thumbnail-image {
    transform: scale(1.05);
  }

  .thumbnail-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.3);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.25s ease;
  }

  .thumbnail-overlay i {
    color: white;
    font-size: 36px;
  }

  .thumbnail-card:hover .thumbnail-overlay {
    opacity: 1;
  }

  .thumbnail-icon-container {
    width: 100%;
    height: 160px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f5f5f5;
  }

  .folder-icon-bg {
    background: var(--color-primary-light);
  }

  .thumbnail-icon-container i.large {
    font-size: 64px;
  }

  .thumbnail-name {
    padding: 10px 12px;
    font-size: 13px;
    color: #333;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: center;
    background: #fafafa;
    border-top: 1px solid #eee;
  }

  .folder-card .thumbnail-name {
    font-weight: 500;
    color: #1976D2;
  }

  .thumbnail-actions {
    position: absolute;
    top: 8px;
    right: 8px;
    display: flex;
    gap: 4px;
    opacity: 0;
    transition: opacity 0.25s ease;
  }

  .thumbnail-item:hover .thumbnail-actions {
    opacity: 1;
  }

  .thumbnail-actions .btn-small {
    width: 30px;
    height: 30px;
  }

  .thumbnail-actions .btn-small i {
    font-size: 16px;
    line-height: 30px;
  }

  /* 缩略图多选复选框 */
  .thumbnail-checkbox {
    position: absolute;
    top: 8px;
    left: 8px;
    z-index: 5;
    cursor: pointer;
    user-select: none;
  }

  .thumbnail-checkmark {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 22px;
    height: 22px;
    border-radius: 4px;
    background: transparent;
    border: 2px solid rgba(255, 255, 255, 0.7);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
    transition: all 0.15s ease;
  }

  .thumbnail-checkmark.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
    box-shadow: none;
  }

  .thumbnail-checkmark.checked::after {
    content: '';
    display: block;
    width: 6px;
    height: 10px;
    border: solid white;
    border-width: 0 2px 2px 0;
    transform: rotate(45deg) translate(-1px, -1px);
  }

  .thumbnail-checkbox:hover .thumbnail-checkmark {
    border-color: white;
  }

  /* 选中卡片高亮 */
  .thumbnail-card.selected {
    box-shadow: 0 0 0 2px var(--color-primary), 0 4px 16px rgba(37, 99, 235, 0.25);
    border-color: var(--color-primary);
  }

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
  }

  .lightbox-image {
    max-width: 85vw;
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