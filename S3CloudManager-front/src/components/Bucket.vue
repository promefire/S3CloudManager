<template>
    <div>
      <!-- 页面工具栏 -->
      <div class="page-toolbar">
        <div class="toolbar-container">
          <div class="toolbar-left">
            <a href="#" class="page-title-link" @click.prevent="navigateUp">
              <i class="material-icons">arrow_back</i>
              <span class="page-title">{{ bucketName }}</span>
            </a>
            <div class="toolbar-breadcrumbs" v-if="breadcrumbs.length > 1">
              <a v-for="crumb in breadcrumbs.slice(1)" :key="crumb.path" href="#" class="breadcrumb-link" @click="navigateTo(crumb.path)">{{ crumb.name }}</a>
            </div>
          </div>
          <div class="toolbar-right">
            <button v-if="!objects.length" @click="deleteBucket" class="btn btn-danger waves-effect waves-light">
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

        <!-- 上传文件面板（子组件） -->
        <UploadPanel
          :bucket-name="bucketName"
          :current-path="currentPath"
          :show-upload-panel="showUploadPanel"
          @close="showUploadPanel = false"
          @upload-success="handleUploadSuccess" />

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
          <div v-for="object in objects" :key="`${pagination.current_page}-${object.Key}`" class="thumbnail-item">
            <div class="thumbnail-card" :class="{ 'folder-card': object.IsFolder, 'selected': !object.IsFolder && multiSelectMode && selectedObjects.includes(object.Key) }" @click="object.IsFolder ? handleObjectClick(object) : (multiSelectMode ? toggleObjectSelection(object.Key) : handleObjectClick(object))">
              <!-- 多选复选框 -->
              <div v-if="!object.IsFolder && multiSelectMode" class="thumbnail-checkbox" @click.stop="toggleObjectSelection(object.Key)">
                <span class="thumbnail-checkmark" :class="{ 'checked': selectedObjects.includes(object.Key) }"></span>
              </div>
              <!-- 图片缩略图 -->
              <div v-if="isImageFile(object.Key)" class="thumbnail-image-container">
                <img :data-src="getThumbnailUrl(object.Key)" :alt="object.DisplayName" class="thumbnail-image lazy-image" @error="handleImageError($event)" />
                <div class="lazy-placeholder">
                  <i class="material-icons">image</i>
                </div>
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

      <!-- 复制格式选择菜单（子组件） -->
      <CopyFormatMenu
        :copy-menu="copyMenu"
        @close="hideCopyMenu"
        @copy="hideCopyMenu" />

      <!-- 图片预览 Lightbox（子组件） -->
      <ImageLightbox
        :preview-image="previewImage"
        :image-objects="getImageObjects()"
        :current-image-index="currentImageIndex"
        @close="closePreview"
        @prev="prevImage"
        @next="nextImage"
        @show-copy-menu="onLightboxShowCopyMenu" />
    </div>
  </template>

  <script>
  /* global M */
import { API_ENDPOINTS, getImageDomain, THUMBNAIL_PARAMS, getHeaders } from '../config/api.js';
import { isImageFile, formatFileSize } from '../utils/file-utils.js';
import UploadPanel from './UploadPanel.vue';
import ImageLightbox from './ImageLightbox.vue';
import CopyFormatMenu from './CopyFormatMenu.vue';

export default {
  name: 'Bucket',
  components: {
    UploadPanel,
    ImageLightbox,
    CopyFormatMenu
  },
  props: ['bucketName'],
  data() {
    return {
      objects: [],
      currentPath: '',
      newFolderName: '',
      allowDelete: true,
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
      showUploadPanel: false,
      multiSelectMode: false,
      viewMode: 'thumbnail',
      previewImage: null,
      copyMenu: {
        visible: false,
        objectKey: '',
        filename: ''
      },
      loadedImages: {},
      imageObserver: null,
      currentImageIndex: 0,
      currentDomain: ''
    }
  },
  mounted() {
    this.currentDomain = getImageDomain(this.bucketName);
    this.fetchObjects();
    M.Modal.init(document.querySelectorAll('.modal'));
    document.addEventListener('keydown', this.handleKeydown);
    this.initImageObserver();
  },

  beforeUnmount() {
    document.removeEventListener('keydown', this.handleKeydown);
    document.body.style.overflow = '';
    if (this.imageObserver) {
      this.imageObserver.disconnect();
    }
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
          const params = {
            page: this.pagination.current_page,
            page_size: this.pagination.page_size
          };
          url = API_ENDPOINTS.browseFolder(this.bucketName, this.currentPath, params);
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
          const displayName = isFolder ? object.name.replace(/\/$/, '').split('/').pop() : object.name.split('/').pop();

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

        // 更新当前存储桶的自定义域名
        this.currentDomain = getImageDomain(this.bucketName);

        // 数据更新后设置懒加载观察
        this.setupLazyImages();
      } catch (error) {
        console.error('Error fetching objects:', error);
        this.$root.notify('Failed to load objects', 'error');
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

        this.$root.notify('Bucket deleted successfully!', 'success');
        this.$router.push('/');
      } catch (error) {
        console.error('Error deleting bucket:', error);
        this.$root.notify('Failed to delete bucket', 'error');
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
    navigateUp() {
      // 如果在根目录，返回存储桶列表
      if (!this.currentPath || this.currentPath === '') {
        this.$router.push('/');
        return;
      }
      // 去掉末尾 /
      let path = this.currentPath.replace(/\/$/, '');
      // 找到最后一个 / 的位置
      const lastSlashIndex = path.lastIndexOf('/');
      if (lastSlashIndex <= 0) {
        // 没有更多层级，回到根目录
        this.currentPath = '';
      } else {
        // 回到上一级
        this.currentPath = path.substring(0, lastSlashIndex + 1);
      }
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

        this.$root.notify('Object deleted successfully!', 'success');
        this.fetchObjects();
      } catch (error) {
        console.error('Error deleting object:', error);
        this.$root.notify('Failed to delete object', 'error');
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
      return isImageFile(filename);
    },

    // 获取图片的真实URL
    getImageRealUrl(filename) {
      return `${this.currentDomain}/${filename}`;
    },

    // 获取缩略图URL（带CDN裁剪参数）
    getThumbnailUrl(filename) {
      return `${this.currentDomain}/${filename}?${THUMBNAIL_PARAMS}`;
    },

    // 处理对象点击事件
    handleObjectClick(object) {
      if (object.IsFolder) {
        // 如果是文件夹，导航到文件夹
        this.navigateTo(object.Key);
      } else if (isImageFile(object.Key)) {
        // 如果是图片文件，打开预览
        this.openPreview(object);
      }
      // 其他文件类型不处理点击事件
    },

    // 切换视图模式
    toggleViewMode() {
      this.viewMode = this.viewMode === 'list' ? 'thumbnail' : 'list';
    },

    // 打开图片预览
    openPreview(object) {
      // 找到当前图片在图片列表中的索引（仅图片，排除文件夹）
      this.currentImageIndex = this.getImageIndex(object.Key);
      this.previewImage = {
        name: object.DisplayName,
        url: this.getImageRealUrl(object.Key),
        size: object.Size ? formatFileSize(object.Size) : '',
        date: object.LastModified ? this.formatDateTime(object.LastModified) : ''
      };
      document.body.style.overflow = 'hidden';
      this.$nextTick(() => {
        // focus is now handled by ImageLightbox component
      });
    },

    // 获取当前图片在图片列表中的索引
    getImageIndex(objectKey) {
      const imageObjects = this.objects.filter(obj => !obj.IsFolder && isImageFile(obj.Key));
      return imageObjects.findIndex(obj => obj.Key === objectKey);
    },

    // 获取所有图片对象（排除文件夹）
    getImageObjects() {
      return this.objects.filter(obj => !obj.IsFolder && isImageFile(obj.Key));
    },

    // 切换到上一张图片
    prevImage() {
      const images = this.getImageObjects();
      if (!images.length) return;
      // 如果是第一张，跳到最后一张（循环）
      const newIndex = this.currentImageIndex <= 0 ? images.length - 1 : this.currentImageIndex - 1;
      this.currentImageIndex = newIndex;
      const obj = images[newIndex];
      this.previewImage = {
        name: obj.DisplayName,
        url: this.getImageRealUrl(obj.Key),
        size: obj.Size ? formatFileSize(obj.Size) : '',
        date: obj.LastModified ? this.formatDateTime(obj.LastModified) : ''
      };
    },

    // 切换到下一张图片
    nextImage() {
      const images = this.getImageObjects();
      if (!images.length) return;
      // 如果是最后一张，跳到第一张（循环）
      const newIndex = this.currentImageIndex >= images.length - 1 ? 0 : this.currentImageIndex + 1;
      this.currentImageIndex = newIndex;
      const obj = images[newIndex];
      this.previewImage = {
        name: obj.DisplayName,
        url: this.getImageRealUrl(obj.Key),
        size: obj.Size ? formatFileSize(obj.Size) : '',
        date: obj.LastModified ? this.formatDateTime(obj.LastModified) : ''
      };
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
        filename,
        bucketName: this.bucketName
      };
    },

    // 隐藏复制格式选择菜单
    hideCopyMenu() {
      this.copyMenu.visible = false;
    },

    // Lightbox 触发复制菜单
    onLightboxShowCopyMenu(payload) {
      this.showCopyMenu(payload.objectKey, payload.filename);
    },

    // 上传成功回调
    handleUploadSuccess() {
      this.fetchObjects();
    },

    // 初始化图片懒加载观察器
    initImageObserver() {
      this.imageObserver = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
          if (entry.isIntersecting) {
            const img = entry.target;
            const src = img.dataset.src;
            if (src) {
              // 加载完成：直接操作 DOM 移除占位符 + 添加 class（不依赖 Vue 响应式）
              img.onload = () => {
                img.classList.add('loaded');
                const placeholder = img.parentElement.querySelector('.lazy-placeholder');
                if (placeholder) placeholder.remove();
              };
              // 加载失败：同样移除占位符，显示错误图标
              img.onerror = () => {
                img.classList.add('loaded');
                if (img.parentElement) {
                  const placeholder = img.parentElement.querySelector('.lazy-placeholder');
                  if (placeholder) placeholder.remove();
                }
                this.handleImageError({ target: img });
              };
              img.src = src;
            }
            this.imageObserver.unobserve(img);
          }
        });
      }, {
        rootMargin: '200px 0px',
        threshold: 0.01
      });
    },

    // 观察所有懒加载图片元素
    setupLazyImages() {
      this.$nextTick(() => {
        const lazyImages = document.querySelectorAll('.lazy-image:not([src])');
        lazyImages.forEach(img => {
          this.imageObserver.observe(img);
        });
      });
    },

    // 图片加载失败处理
    handleImageError(event) {
      const img = event.target;
      img.style.display = 'none';
      if (img.parentElement) {
        img.parentElement.innerHTML = `<i class="material-icons large" style="color: #999;">broken_image</i>`;
      }
    },

    // 键盘事件处理
    handleKeydown(event) {
      if (!this.previewImage) return;
      if (event.key === 'Escape') {
        this.closePreview();
      } else if (event.key === 'ArrowLeft') {
        this.prevImage();
      } else if (event.key === 'ArrowRight') {
        this.nextImage();
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
        this.$root.notify('Please enter a folder name', 'error');
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

        this.$root.notify('Folder created successfully!', 'success');
        this.newFolderName = '';
        const modal = M.Modal.getInstance(document.getElementById('modal-create-folder'));
        modal.close();
        this.fetchObjects();
      } catch (error) {
        console.error('Error creating folder:', error);
        this.$root.notify('Failed to create folder', 'error');
      }
    },

    // 复制图片URL
    async copyImageUrl(objectKey) {
      try {
        const imageUrl = this.getImageRealUrl(objectKey);
        await navigator.clipboard.writeText(imageUrl);
        this.$root.notify('图片链接已复制到剪贴板', 'success');
      } catch (error) {
        console.error('复制失败:', error);
        // 降级方案：使用传统方法复制
        const textArea = document.createElement('textarea');
        textArea.value = this.getImageRealUrl(objectKey);
        document.body.appendChild(textArea);
        textArea.select();
        try {
          document.execCommand('copy');
          this.$root.notify('图片链接已复制到剪贴板', 'success');
        } catch (fallbackError) {
          console.error('降级复制也失败:', fallbackError);
          this.$root.notify('复制失败，请手动复制', 'error');
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
        this.$root.notify('文件删除成功！', 'success');
        this.fetchObjects();
      } catch (error) {
        console.error('Error deleting object:', error);
        this.$root.notify('删除文件失败', 'error');
      }
    },

    // 批量删除功能
    async batchDeleteObjects() {
      if (!this.selectedObjects.length) {
        this.$root.notify('Please select objects to delete', 'error');
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
        this.$root.notify(`Deleted ${result.delete_count || result.success_count} objects successfully!`, 'success');
        this.selectedObjects = [];
        this.fetchObjects();
      } catch (error) {
        console.error('Error deleting objects:', error);
        this.$root.notify('Failed to delete objects', 'error');
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
      return formatFileSize(bytes);
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

  /* ========== 图片懒加载 ========== */
  .lazy-image {
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .lazy-image.loaded {
    opacity: 1;
  }

  .lazy-placeholder {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #f5f5f5 25%, #eeeeee 50%, #f5f5f5 75%);
    background-size: 200% 200%;
    animation: shimmer 1.5s ease-in-out infinite;
    z-index: 1;
  }

  .lazy-placeholder i {
    font-size: 32px;
    color: #ccc;
  }

  @keyframes shimmer {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
  }
  </style>
