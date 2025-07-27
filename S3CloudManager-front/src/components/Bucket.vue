<template>
    <div>
      <nav class="nav-extended" style="background-color: #E0FFFF;">
        <div class="nav-wrapper container" style="background-color: #E3F2FD;">
          <a href="#" class="brand-logo center" style="color: #E0FFFF;"><i class="material-icons" style="color: #1976D2;">folder_open</i>{{ bucketName }}</a>
          <ul class="right">
            <li v-if="!objects.length" style="margin-right: 20px;">
              <a class="waves-effect waves-light btn red" @click="deleteBucket" style="color: white; font-weight: 500; text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);">
                删除 <i class="material-icons right" style="color: white;">delete</i>
              </a>
            </li>
          </ul>
        </div>
        <div class="nav-wrapper container" style="background-color: #E3F2FD;">
          <router-link to="/" class="breadcrumb" style="color: #1976D2;"><i class="material-icons" style="color: #1976D2;">arrow_back</i> 存储桶</router-link>
          <a href="#" class="breadcrumb" @click="navigateTo('')" style="color: #1976D2;">{{ bucketName }}</a>
          <a v-for="crumb in breadcrumbs.slice(1)" :key="crumb.path" href="#" class="breadcrumb" @click="navigateTo(crumb.path)" style="color: #1976D2;">{{ crumb.name }}</a>
        </div>
      </nav>
              <div class="section" style="margin: 10px; position: relative;">
          
          
          <!-- 文件上传区域 -->
        <div class="row" style="margin-bottom: 20px;">
          <div class="col s12">
            <!-- 小上传按钮 -->
            <div v-if="!showUploadPanel" class="row" style="margin-bottom: 0;">
              <div class="col s6">
                <div class="left">
                  <button @click="showUploadPanel = true" class="btn waves-effect waves-light" style="background-color: #1976D2; margin-right: 10px;">
                    <i class="material-icons left" style="color: white; font-size: 20px; margin-right: 8px;">cloud_upload</i>上传文件
                  </button>
                  <a class="waves-effect waves-light btn modal-trigger" data-target="modal-create-folder" style="background-color: #1976D2; color: white; font-weight: 500; text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3); margin-right: 10px;">
                    <i class="material-icons left" style="color: white; font-size: 20px; margin-right: 8px;">create_new_folder</i>新建文件夹
                  </a>
                  <button @click="toggleMultiSelect" class="btn waves-effect waves-light" :style="multiSelectMode ? 'background-color: #4CAF50;' : 'background-color: #1976D2;'">
                    <i class="material-icons left" style="color: white; font-size: 20px; margin-right: 8px;">check_box</i>{{ multiSelectMode ? '退出多选' : '多选' }}
                  </button>
                </div>
              </div>
              <div class="col s6">
                <div class="right" v-if="multiSelectMode && selectedObjects.length">
                  <button @click="batchDeleteObjects" class="btn red waves-effect waves-light">
                    <i class="material-icons left" style="color: white; font-size: 18px; margin-right: 6px;">delete</i>删除选中 ({{ selectedObjects.length }})
                  </button>
                </div>
              </div>
            </div>
            
            <!-- 完整上传面板 -->
            <div v-if="showUploadPanel" class="card-panel" 
                 style="background-color: #f5f5f5; border: 2px dashed #ccc;"
                 @drop="handleDrop"
                 @dragover="handleDragOver"
                 @dragenter="handleDragEnter"
                 @dragleave="handleDragLeave"
                 :class="{ 'drag-over': isDragOver }">
              <div class="row">
                <div class="col s12">
                  <div class="row" style="margin-bottom: 0;">
                    <div class="col s10">
                      <h6 style="margin-top: 0; color: #666;">
                        <i class="material-icons left">cloud_upload</i>文件上传
                      </h6>
                    </div>
                    <div class="col s2">
                      <button @click="closeUploadPanel" class="btn-floating btn-small grey waves-effect waves-light">
                        <i class="material-icons">close</i>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
              
              <div class="row">
                <div class="col s12">
                  <div class="file-field input-field">
                    <div class="btn waves-effect waves-light">
                      <i class="material-icons left" style="color: white; font-size: 18px; margin-right: 6px;">folder_open</i>
                      <span>选择文件</span>
                      <input type="file" multiple @change="handleFileSelect" accept="*/*">
                    </div>
                    <div class="file-path-wrapper">
                      <input class="file-path validate" type="text" placeholder="选择要上传的文件">
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- 拖拽提示 -->
              <div class="row" style="margin-top: 15px;">
                <div class="col s12 center">
                  <div v-if="isDragOver" class="drag-hint" style="padding: 20px; background-color: #e3f2fd; border: 2px dashed #1976D2; border-radius: 8px;">
                    <i class="material-icons large" style="color: #1976D2;">cloud_upload</i>
                    <p style="margin: 10px 0; color: #1976D2; font-weight: bold;">释放鼠标上传文件</p>
                  </div>
                  <div v-else class="drag-hint" style="padding: 15px; color: #666;">
                    <i class="material-icons" style="vertical-align: middle;">drag_handle</i>
                    <span style="margin-left: 5px;">或者拖拽文件到此处上传</span>
                  </div>
                </div>
              </div>
              
              <!-- 已选择的文件列表 -->
              <div v-if="uploadFiles.length" class="row" style="margin-top: 10px;">
                <div class="col s12">
                  <h6 style="color: #666; margin-bottom: 10px;">已选择的文件 ({{ uploadFiles.length }} 个):</h6>
                  <ul class="collection">
                    <li v-for="(file, index) in uploadFiles" :key="index" class="collection-item">
                      <div class="row" style="margin-bottom: 0;">
                        <div class="col s10">
                          <i class="material-icons left" style="color: #1976D2;">insert_drive_file</i>
                          <span>{{ file.name }}</span>
                          <br>
                          <small class="grey-text">{{ formatFileSize(file.size) }}</small>
                        </div>
                        <div class="col s2">
                          <button @click="removeFile(index)" class="btn-floating btn-small red waves-effect waves-light">
                            <i class="material-icons">close</i>
                          </button>
                        </div>
                      </div>
                    </li>
                  </ul>
                </div>
              </div>
              
              <!-- 上传按钮 -->
              <div v-if="uploadFiles.length" class="row" style="margin-top: 15px;">
                <div class="col s12 center">
                  <button @click="uploadSelectedFiles" :disabled="isUploading" class="btn-large waves-effect waves-light" style="background-color: #1976D2;">
                    <i class="material-icons left" style="color: white; font-size: 20px; margin-right: 8px;">cloud_upload</i>
                    {{ isUploading ? '正在上传...' : `上传 ${uploadFiles.length} 个文件` }}
                  </button>
                  <button @click="clearUploadFiles" class="btn-flat waves-effect waves-light" style="margin-left: 10px;">
                    <i class="material-icons left" style="color: #666; font-size: 18px; margin-right: 6px;">clear</i>清空选择
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <table class="striped" v-if="objects.length">
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
                <button v-if="!object.IsFolder && isImageFile(object.Key)" @click="copyImageUrl(object.Key)" class="btn-floating btn-small waves-effect waves-light" style="background-color: #1976D2; margin-right: 5px;" :title="'复制图片链接'">
                  <i class="material-icons">content_copy</i>
                </button>
                <button v-if="!object.IsFolder" @click="deleteSingleObject(object.Key)" class="btn-floating btn-small waves-effect waves-light red" :title="'删除文件'">
                  <i class="material-icons">delete</i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        
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
            <button type="submit" class="waves-effect waves-green btn">创建</button>
          </div>
        </form>
      </div>
    </div>
  </template>
  
  <script>
  /* global M */
import { API_ENDPOINTS, IMAGE_DOMAIN } from '../config/api.js';
  
  // eslint-disable-next-line vue/multi-word-component-names
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
        multiSelectMode: false
      }
    },
    mounted() {
      this.fetchObjects();
      M.Modal.init(document.querySelectorAll('.modal'));
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
            console.log('browseFolder will be called with:', this.bucketName, this.currentPath);
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
          const response = await fetch(url);
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
            method: 'DELETE'
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
            method: 'DELETE'
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
          // 如果是图片文件，打开真实URL
          const imageUrl = this.getImageRealUrl(object.Key);
          window.open(imageUrl, '_blank');
        }
        // 其他文件类型不处理点击事件
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
          const uploadPromises = this.uploadFiles.map(async (file) => {
            const formData = new FormData();
            formData.append('file', file);
            
            // 如果有当前路径，将文件上传到当前文件夹
            if (this.currentPath) {
              formData.append('object_name', `${this.currentPath}${file.name}`);
            }
            
            const response = await fetch(API_ENDPOINTS.uploadObject(this.bucketName), {
              method: 'POST',
              body: formData
            });
            
            if (!response.ok) {
              throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            return response.json();
          });
          
          await Promise.all(uploadPromises);
          
          M.toast({ html: 'Files uploaded successfully!', classes: 'green' });
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