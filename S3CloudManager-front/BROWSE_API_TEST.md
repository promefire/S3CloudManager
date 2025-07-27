# Browse API 修复测试

## 问题分析

用户反馈：在存储桶 `blog-img` 下有一个同名的文件夹 `blog-img`，访问这个文件夹的 API 是 `/api/v1/buckets/blog-img/browse/blog-img/`，但前端代码使用的是 `/api/v1/buckets/blog-img/objects` 接口。

## API 接口区别

### 1. Objects API (`/api/v1/buckets/{bucket}/objects`)
- 用于列出存储桶中的对象
- 支持 prefix 参数来过滤对象
- 适用于根目录浏览

### 2. Browse API (`/api/v1/buckets/{bucket}/browse/{filepath}`)
- 用于浏览文件夹或下载文件
- 当 filepath 是文件夹路径时，返回与 objects API 相同的格式
- 适用于文件夹浏览

## 修复方案

### 1. 添加新的 API 端点
```javascript
// 浏览文件夹（返回与列出对象相同的格式）
browseFolder: (bucketName, folderPath, params = {}) => {
  const queryParams = new URLSearchParams();
  if (params.page) queryParams.append('page', params.page);
  if (params.page_size) queryParams.append('page_size', params.page_size);
  if (params.disposition) queryParams.append('disposition', params.disposition);
  
  const queryString = queryParams.toString();
  const path = folderPath.endsWith('/') ? folderPath : folderPath + '/';
  return `${API_BASE_URL}/buckets/${bucketName}/browse/${encodeURIComponent(path)}${queryString ? '?' + queryString : ''}`;
}
```

### 2. 修改 fetchObjects 方法
```javascript
async fetchObjects() {
  try {
    let url;
    
    // 如果当前路径不为空，使用 browse 接口来浏览文件夹
    if (this.currentPath) {
      const params = {
        page: this.pagination.current_page,
        page_size: this.pagination.page_size
      };
      url = API_ENDPOINTS.browseFolder(this.bucketName, this.currentPath, params);
    } else {
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
    // ... 其余代码
  }
}
```

## 测试场景

### 场景 1: 访问存储桶根目录
- **当前路径**: `""` (空字符串)
- **API 调用**: `GET /api/v1/buckets/blog-img/objects?page=1&page_size=20&folder=true&delimiter=/`
- **预期结果**: 显示存储桶根目录下的文件和文件夹

### 场景 2: 访问 blog-img 文件夹
- **当前路径**: `"blog-img/"`
- **API 调用**: `GET /api/v1/buckets/blog-img/browse/blog-img/?page=1&page_size=20`
- **预期结果**: 显示 blog-img 文件夹下的内容（2024/ 文件夹和图片文件）

### 场景 3: 访问 2024 文件夹
- **当前路径**: `"blog-img/2024/"`
- **API 调用**: `GET /api/v1/buckets/blog-img/browse/blog-img/2024/?page=1&page_size=20`
- **预期结果**: 显示 2024 文件夹下的内容

## 验证步骤

1. **进入 blog-img 存储桶**
   - 应该看到根目录内容，包括 blog-img 文件夹

2. **点击 blog-img 文件夹**
   - 应该看到该文件夹下的内容
   - 检查控制台日志，确认使用的是 browse API

3. **使用面包屑导航**
   - 点击面包屑应该能正确导航
   - 确认路径和 API 调用正确

4. **检查 API 响应**
   - 确认不再返回 `objects: null`
   - 确认能正确显示文件列表

## 预期结果

修复后：
- ✅ 访问 blog-img 文件夹时使用正确的 browse API
- ✅ 不再出现 `objects: null` 错误
- ✅ 正确显示文件夹内容
- ✅ 面包屑导航正常工作
- ✅ 分页功能正常工作 