# S3 Cloud Manager

一个基于 Vue.js 的 S3 兼容对象存储管理前端应用，提供文件上传、下载、删除、文件夹管理等功能。

## 功能特性

- 🔐 用户登录认证
- 📁 存储桶管理
- 📂 多级文件夹导航
- 📤 文件上传（支持拖拽）
- 📋 批量文件操作
- 🔗 图片链接复制
- 🗑️ 文件删除
- 📄 分页显示

## 快速开始

### 环境要求

- Node.js (版本 14 或更高)
- npm 或 yarn

### 安装依赖

```bash
# 使用 npm
npm install

# 或使用 yarn
yarn install
```

### 配置后端API

在 `src/config/api.js` 文件中配置后端API地址：

```javascript
const config = {
  development: {
    baseURL: 'http://localhost:8080/api/v1'  // 修改为你的后端API地址
  },
  production: {
    baseURL: '/api/v1'
  }
};
```

### 启动开发服务器

```bash
# 使用 npm
npm run serve

# 或使用 yarn
yarn serve
```

应用将在 `http://localhost:8080` 启动（如果端口被占用，会自动使用其他端口）。

### 构建生产版本

```bash
# 使用 npm
npm run build

# 或使用 yarn
yarn build
```

构建完成后，`dist` 目录包含可部署的静态文件。

## 使用说明

### 登录

1. 访问应用首页
2. 使用默认账号登录：
   - 用户名：`admin`
   - 密码：`admin`

### 主要功能

1. **存储桶管理**
   - 查看所有存储桶
   - 创建新存储桶
   - 删除空存储桶

2. **文件管理**
   - 浏览文件和文件夹
   - 上传文件（支持多文件选择和拖拽）
   - 创建新文件夹
   - 复制图片链接
   - 删除文件

3. **批量操作**
   - 点击"多选"按钮进入多选模式
   - 选择多个文件
   - 批量删除选中的文件

## 项目结构

```
src/
├── components/          # Vue组件
│   ├── Bucket.vue      # 存储桶详情页
│   ├── Buckets.vue     # 存储桶列表页
│   └── Login.vue       # 登录页
├── config/
│   └── api.js          # API配置
├── router/
│   └── index.js        # 路由配置
└── App.vue             # 根组件
```

## 技术栈

- **前端框架**: Vue.js 2.x
- **UI框架**: Materialize CSS
- **路由**: Vue Router
- **构建工具**: Vue CLI

## 截图展示

### 登录页
![登录页](https://img.promefire.top/blog-img/2025/07/image-20250727173405779.png)

### 存储桶页
![存储桶页](https://img.promefire.top/blog-img/2025/07/image-20250727173509334.png)

### 多级目录
![多级目录](https://img.promefire.top/blog-img/2025/07/image-20250727173541195.png)

