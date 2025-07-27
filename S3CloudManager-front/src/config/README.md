# API 配置说明

## 配置文件位置
- `src/config/api.js` - API 配置文件

## 配置方式

### 1. 环境变量配置（推荐）
在项目根目录创建 `.env.development` 文件：
```
VUE_APP_API_BASE_URL=http://localhost:8080/api/v1
VUE_APP_IMAGE_DOMAIN=https://img.promefire.top
```

在项目根目录创建 `.env.production` 文件：
```
VUE_APP_API_BASE_URL=/api/v1
VUE_APP_IMAGE_DOMAIN=https://img.promefire.top
```

### 2. 直接修改配置文件
在 `src/config/api.js` 中直接修改 `baseURL` 值。

## 支持的配置项

### API_BASE_URL
- **开发环境**: `http://localhost:8080/api/v1`
- **生产环境**: `/api/v1`

### IMAGE_DOMAIN
- **默认值**: `https://img.promefire.top`
- **用途**: 图片文件的真实访问域名

### API_ENDPOINTS
预定义的 API 端点：

**健康检查**
- `health` - 健康检查

**存储桶管理**
- `buckets` - 获取存储桶列表
- `bucket(bucketName)` - 获取特定存储桶

**对象管理**
- `bucketObjects(bucketName)` - 获取存储桶中的对象
- `bucketObjectInfo(bucketName, objectKey)` - 获取对象详细信息
- `bucketObject(bucketName, objectKey)` - 删除特定对象
- `batchDeleteObjects(bucketName)` - 批量删除对象

**文件夹管理**
- `createFolder(bucketName)` - 创建文件夹

**浏览和下载**
- `browseFile(bucketName, filepath)` - 浏览文件或下载

## 使用示例

```javascript
import { API_ENDPOINTS } from '../config/api.js';

// 健康检查
const response = await fetch(API_ENDPOINTS.health);

// 获取存储桶列表
const response = await fetch(API_ENDPOINTS.buckets);

// 获取特定存储桶的对象
const response = await fetch(API_ENDPOINTS.bucketObjects('my-bucket'));

// 删除特定对象
const response = await fetch(API_ENDPOINTS.bucketObject('my-bucket', 'file.txt'), {
  method: 'DELETE'
});

// 创建文件夹
const response = await fetch(API_ENDPOINTS.createFolder('my-bucket'), {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ name: 'new-folder' })
});

// 下载文件
const downloadUrl = API_ENDPOINTS.browseFile('my-bucket', 'file.txt');
```

## 环境切换
- 开发环境：`npm run serve` 或 `yarn serve`
- 生产环境：`npm run build` 或 `yarn build` 