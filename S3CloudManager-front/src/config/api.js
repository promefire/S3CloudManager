// API 配置文件
const config = {
  // 开发环境配置
  development: {
    baseURL: process.env.VUE_APP_API_BASE_URL || 'http://localhost:8080/api/v1'
  },
  // 生产环境配置
  production: {
    baseURL: process.env.VUE_APP_API_BASE_URL || '/api/v1' // 生产环境通常使用相对路径
  }
};

// 根据当前环境获取配置
const env = process.env.NODE_ENV || 'development';
const apiConfig = config[env];

// 导出 API 基础 URL
export const API_BASE_URL = apiConfig.baseURL;

// 导出常用的 API 端点
export const API_ENDPOINTS = {
  // 健康检查
  health: `${API_BASE_URL.replace('/api/v1', '/api')}/health`,
  
  // 存储桶管理
  buckets: `${API_BASE_URL}/buckets`,
  bucket: (bucketName) => `${API_BASE_URL}/buckets/${bucketName}`,
  
  // 对象管理
  bucketObjects: (bucketName, params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.page) queryParams.append('page', params.page);
    if (params.page_size) queryParams.append('page_size', params.page_size);
    if (params.prefix) queryParams.append('prefix', params.prefix);
    if (params.folder !== undefined) queryParams.append('folder', params.folder);
    if (params.delimiter) queryParams.append('delimiter', params.delimiter);
    
    const queryString = queryParams.toString();
    return `${API_BASE_URL}/buckets/${bucketName}/objects${queryString ? '?' + queryString : ''}`;
  },
  uploadObject: (bucketName) => `${API_BASE_URL}/buckets/${bucketName}/objects`,
  bucketObjectInfo: (bucketName, objectKey) => `${API_BASE_URL}/buckets/${bucketName}/api/objects/${encodeURIComponent(objectKey)}/info`,
  bucketObject: (bucketName, objectKey) => `${API_BASE_URL}/buckets/${bucketName}/api/objects/${encodeURIComponent(objectKey)}`,
  batchDeleteObjects: (bucketName) => `${API_BASE_URL}/buckets/${bucketName}/objects/batch-delete`,
  
  // 文件夹管理
  createFolder: (bucketName) => `${API_BASE_URL}/buckets/${bucketName}/folders`,
  
  // 浏览和下载
  browseFile: (bucketName, filepath) => `${API_BASE_URL}/buckets/${bucketName}/browse/${encodeURIComponent(filepath)}`,
  
  // 浏览文件夹（返回与列出对象相同的格式）
  browseFolder: (bucketName, folderPath, params = {}) => {
    console.log('browseFolder called with:', bucketName, folderPath, params);
    const queryParams = new URLSearchParams();
    if (params.page) queryParams.append('page', params.page);
    if (params.page_size) queryParams.append('page_size', params.page_size);
    if (params.disposition) queryParams.append('disposition', params.disposition);
    
    const queryString = queryParams.toString();
    const path = folderPath.endsWith('/') ? folderPath : folderPath + '/';
    const encodedPath = encodeURIComponent(path);
    const url = `${API_BASE_URL}/buckets/${bucketName}/browse/${encodedPath}${queryString ? '?' + queryString : ''}`;
    console.log('browseFolder generated URL:', url);
    return url;
  }
};

// 图片域名配置
export const IMAGE_DOMAIN = process.env.VUE_APP_IMAGE_DOMAIN || 'https://img.promefire.top';

export default apiConfig; 