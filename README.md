# S3CloudManager

一个基于Vue.js + Go的S3兼容对象存储管理系统，提供完整的存储桶和对象管理功能。

## 🚀 项目结构

```
S3CloudManager/
├── S3CloudManager-front/          # 前端项目 (Vue.js)
│   ├── src/
│   │   ├── components/            # Vue组件
│   │   ├── config/                # API配置
│   │   └── router/                # 路由配置
│   ├── public/                    # 静态资源
│   └── package.json               # 前端依赖
├── S3CloudManager-server/         # 后端项目 (Go)
│   ├── app/                       # 业务逻辑
│   ├── main.go                    # 主程序入口
│   ├── routes.go                  # 路由定义
│   └── config.env                 # 配置文件
└── README.md                      # 项目文档
```

## 🎯 功能特性

### 存储桶管理
- ✅ 列出所有存储桶
- ✅ 创建新存储桶
- ✅ 删除存储桶

### 对象管理
- ✅ 上传文件/对象
- ✅ 下载文件/对象
- ✅ 删除对象
- ✅ 批量删除对象
- ✅ 获取对象信息

### 文件夹支持
- ✅ 创建文件夹
- ✅ 浏览文件夹结构
- ✅ 层级目录导航

### 用户界面
- ✅ 现代化Web界面
- ✅ 响应式设计
- ✅ 文件拖拽上传
- ✅ 实时文件浏览

## 📋 技术栈

### 前端
- **框架**: Vue.js 3
- **路由**: Vue Router 4
- **UI**: Materialize CSS
- **构建工具**: Vue CLI

### 后端
- **语言**: Go 1.24+
- **Web框架**: Gin
- **S3客户端**: MinIO Go Client
- **配置管理**: Viper

### 存储
- **兼容**: AWS S3, Cloudflare R2, MinIO等S3兼容服务

## 🛠️ 快速开始

### 环境要求

- **Node.js**: 16.0+
- **Go**: 1.24+
- **S3兼容存储服务**

### 1. 克隆项目

```bash
git clone <repository-url>
cd S3CloudManager
```

### 2. 配置后端

```bash
# 进入后端目录
cd S3CloudManager-server

# 复制配置文件
cp config.env.example config.env

# 编辑配置文件
nano config.env
```

配置文件示例：
```env
ENDPOINT=your-s3-endpoint.com
ACCESS_KEY_ID=your-access-key-id
SECRET_ACCESS_KEY=your-secret-access-key
REGION=auto
USE_SSL=true
SIGNATURE_TYPE=V4
PORT=9000
```

### 3. 启动后端服务

```bash
# 安装依赖
go mod download

# 启动服务
go run main.go routes.go
```

后端服务将在 `http://localhost:9000` 启动

### 4. 启动前端服务

```bash
# 进入前端目录
cd S3CloudManager-front

# 安装依赖
npm install
# 或
yarn install

# 启动开发服务器
npm run serve
# 或
yarn serve
```

前端服务将在 `http://localhost:9001` 启动

## 📚 API文档

### 基础信息
- **后端地址**: `http://localhost:9000`
- **API版本**: v1

### 主要接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/health` | 健康检查 |
| GET | `/api/v1/buckets` | 列出所有存储桶 |
| POST | `/api/v1/buckets` | 创建存储桶 |
| DELETE | `/api/v1/buckets/:bucket` | 删除存储桶 |
| GET | `/api/v1/buckets/:bucket/objects` | 列出对象 |
| POST | `/api/v1/buckets/:bucket/objects` | 上传对象 |
| POST | `/api/v1/buckets/:bucket/folders` | 创建文件夹 |
| GET | `/api/v1/buckets/:bucket/browse/*filepath` | 浏览/下载文件 |
| DELETE | `/api/v1/buckets/:bucket/api/objects/:object` | 删除对象 |
| POST | `/api/v1/buckets/:bucket/objects/batch-delete` | 批量删除对象 |

### 响应格式
所有API响应均为JSON格式，成功响应包含数据字段，错误响应包含error字段。

## 🚀 部署指南

### 后端部署

#### 1. 构建可执行文件
```bash
cd S3CloudManager-server
GOOS=linux GOARCH=amd64 go build -o s3cloudmanager-server main.go routes.go
```

#### 2. 上传到服务器
```bash
scp s3cloudmanager-server user@server-ip:/opt/s3cloudmanager/
scp config.env user@server-ip:/opt/s3cloudmanager/
```

#### 3. 配置系统服务
```bash
# 创建服务文件
sudo nano /etc/systemd/system/s3cloudmanager.service

# 启动服务
sudo systemctl enable s3cloudmanager
sudo systemctl start s3cloudmanager
```

### 前端部署

#### 1. 构建生产版本
```bash
cd S3CloudManager-front
npm run build
```

#### 2. 部署到Web服务器
将 `dist` 目录内容部署到Nginx或Apache服务器。

#### 3. Nginx配置示例
```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /var/www/s3cloudmanager;
    index index.html;

    location /api/ {
        proxy_pass http://localhost:9000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

## 🔧 配置说明

### 后端配置 (config.env)

| 配置项 | 描述 | 默认值 | 必需 |
|--------|------|--------|------|
| `ENDPOINT` | S3服务端点 | - | 是 |
| `ACCESS_KEY_ID` | 访问密钥ID | - | 是 |
| `SECRET_ACCESS_KEY` | 访问密钥 | - | 是 |
| `REGION` | 区域 | `auto` | 否 |
| `USE_SSL` | 是否使用SSL | `true` | 否 |
| `SIGNATURE_TYPE` | 签名类型 | `V4` | 否 |
| `PORT` | 服务器端口 | `9000` | 否 |

### 前端配置

#### 开发环境 (.env.development)
```
VUE_APP_API_BASE_URL=http://localhost:9000/api/v1
VUE_APP_IMAGE_DOMAIN=https://img.promefire.top
```

#### 生产环境 (.env.production)
```
VUE_APP_API_BASE_URL=/api/v1
VUE_APP_IMAGE_DOMAIN=https://img.promefire.top
```



