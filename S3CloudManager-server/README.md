# S3CloudManager Server

一个基于Go语言开发的S3兼容对象存储管理服务器，提供完整的存储桶和对象管理API。

## 🚀 功能特性

- **存储桶管理**: 创建、删除、列出存储桶
- **对象管理**: 上传、下载、删除、更新对象
- **文件夹支持**: 创建文件夹，支持层级目录结构
- **批量操作**: 支持批量删除对象
- **文件浏览**: 支持文件夹浏览和文件下载
- **健康检查**: 提供API健康状态检查
- **CORS支持**: 跨域资源共享支持
- **配置灵活**: 支持环境变量和配置文件

## 📋 系统要求

- Go 1.24.5 或更高版本
- S3兼容的对象存储服务（如AWS S3、Cloudflare R2、MinIO等）

## 🛠️ 安装和配置

### 1. 克隆项目

```bash
git clone <repository-url>
cd S3CloudManager-server
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境变量

创建 `config.env` 文件并配置以下参数：

```env
# S3服务端点
ENDPOINT=your-s3-endpoint.com
# 访问密钥ID
ACCESS_KEY_ID=your-access-key-id
# 访问密钥
SECRET_ACCESS_KEY=your-secret-access-key
# 区域
REGION=auto
# 是否使用SSL
USE_SSL=true
# 签名类型 (V2/V4/V4Streaming/Anonymous)
SIGNATURE_TYPE=V4
# 是否允许删除操作
ALLOW_DELETE=true
# 是否强制下载
FORCE_DOWNLOAD=true
# 服务器端口
PORT=8080
# 是否跳过SSL验证
SKIP_SSL_VERIFICATION=false
```

### 4. 运行服务器

```bash
go run main.go routes.go
```

或者构建后运行：

```bash
go build -o s3cloudmanager-server
./s3cloudmanager-server
```

## 📚 API 文档

### 基础信息

- **基础URL**: `http://localhost:8080`
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

## 🔧 配置说明

### 环境变量

| 变量名 | 描述 | 默认值 | 必需 |
|--------|------|--------|------|
| `ENDPOINT` | S3服务端点 | `s3.amazonaws.com` | 是 |
| `ACCESS_KEY_ID` | 访问密钥ID | - | 是 |
| `SECRET_ACCESS_KEY` | 访问密钥 | - | 是 |
| `REGION` | 区域 | - | 否 |
| `USE_SSL` | 是否使用SSL | `true` | 否 |
| `SIGNATURE_TYPE` | 签名类型 | `V4` | 否 |
| `PORT` | 服务器端口 | `8080` | 否 |
| `SKIP_SSL_VERIFICATION` | 跳过SSL验证 | `false` | 否 |
