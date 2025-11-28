# 网络文件管理系统 (NetFileSys)

一个基于 Go + Vue3 的企业级网络文件管理系统，支持文件上传、下载、分享、权限管理等功能。

## 🚀 功能特性

### 核心功能
- **文件管理**: 上传、下载、预览、删除、移动、复制
- **文件夹管理**: 创建、删除、重命名、层级管理
- **文件分享**: 生成分享链接、设置密码、过期时间
- **权限控制**: RBAC + ACL 细粒度权限管理
- **用户管理**: 用户注册、登录、角色分配
- **组织管理**: 多级组织架构支持

### 高级功能
- **秒传技术**: 基于MD5/SHA1的文件去重
- **分片上传**: 支持大文件断点续传
- **回收站**: 文件软删除与恢复
- **版本管理**: 文件多版本支持
- **审计日志**: 完整的操作记录
- **管理后台**: 系统监控与管理

### 技术特性
- **高性能**: Go语言高并发处理
- **可扩展**: 微服务架构设计
- **安全性**: JWT认证 + 权限控制
- **存储**: 支持MinIO/S3/本地存储
- **缓存**: Redis缓存加速
- **监控**: Prometheus + Grafana

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        前端层                                │
│  ┌─────────────────┐    ┌─────────────────┐                │
│  │   用户端 (Vue3)  │    │ 管理后台 (Vue3)  │                │
│  └─────────────────┘    └─────────────────┘                │
└─────────────────────────────────────────────────────────────┘
                              │
                         REST/gRPC API
                              │
┌─────────────────────────────────────────────────────────────┐
│                        后端服务层                            │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ 用户服务     │ │ 文件服务     │ │ 分享服务     │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ 权限服务     │ │ 审计服务     │ │ 管理服务     │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      数据存储层                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ PostgreSQL  │ │   Redis     │ │   MinIO     │           │
│  │   (主数据)   │ │   (缓存)     │ │ (文件存储)   │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
```

## 📋 环境要求

### 开发环境
- **Go**: 1.21+
- **Node.js**: 18+
- **PostgreSQL**: 13+
- **Redis**: 6+
- **MinIO**: 最新版本

### 生产环境
- **Docker**: 20.10+
- **Docker Compose**: 2.0+

## 🛠️ 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/your-org/netfilessys.git
cd netfilessys
```

### 2. 配置环境变量
```bash
# 复制环境变量模板
cp config/.env.example config/.env

# 编辑配置文件
vim config/.env
```

### 3. 使用Docker Compose启动（推荐）
```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 4. 手动启动开发环境

#### 启动基础服务
```bash
# 启动PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_DB=netfilessys \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=your_password \
  -p 5432:5432 postgres:15

# 启动Redis
docker run -d --name redis \
  -p 6379:6379 redis:7-alpine

# 启动MinIO
docker run -d --name minio \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  -p 9000:9000 -p 9001:9001 \
  minio/minio server /data --console-address ":9001"
```

#### 启动后端服务
```bash
cd backend

# Windows
scripts\start.bat

# Linux/macOS
chmod +x scripts/start.sh
./scripts/start.sh
```

#### 启动前端服务
```bash
# 用户端
cd frontend/user-portal
npm install
npm run dev

# 管理后台
cd frontend/admin-portal
npm install
npm run dev
```

## 🌐 访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 用户端 | http://localhost:3000 | 普通用户界面 |
| 管理后台 | http://localhost:3001 | 管理员界面 |
| API服务 | http://localhost:8080 | 后端API |
| MinIO控制台 | http://localhost:9001 | 对象存储管理 |
| Grafana | http://localhost:3002 | 监控面板 |

## 📚 API 文档

### 认证接口
```bash
# 用户登录
POST /api/v1/auth/login
{
  "username": "admin",
  "password": "password"
}

# 用户注册
POST /api/v1/auth/register
{
  "username": "user",
  "email": "user@example.com",
  "password": "password"
}
```

### 文件接口
```bash
# 上传文件
POST /api/v1/files/upload
Content-Type: multipart/form-data

# 获取文件列表
GET /api/v1/files/list?page=1&page_size=20

# 下载文件
GET /api/v1/files/{id}/download
```

### 分享接口
```bash
# 创建分享
POST /api/v1/shares
{
  "file_id": 1,
  "share_type": "file",
  "expire_time": "2024-12-31T23:59:59Z"
}

# 访问分享
GET /api/v1/public/share/{code}
```

## 🔧 配置说明

### 主要配置项
```yaml
# 服务器配置
server:
  host: "0.0.0.0"
  port: 8080
  mode: "release"

# 数据库配置
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "your_password"
  dbname: "netfilessys"

# Redis配置
redis:
  host: "localhost"
  port: 6379
  password: ""

# MinIO配置
minio:
  endpoint: "localhost:9000"
  access_key_id: "minioadmin"
  secret_access_key: "minioadmin"
  bucket_name: "netfilessys"

# JWT配置
jwt:
  secret: "your_jwt_secret_key_here_min_32_characters"
  expire_hours: 24
```

## 🚀 部署指南

### Docker部署
```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 扩容服务
docker-compose up -d --scale backend=3
```

### Kubernetes部署
```bash
# 应用配置
kubectl apply -f k8s/

# 查看状态
kubectl get pods -n netfilessys
```

## 🔍 监控与日志

### 系统监控
- **Prometheus**: 指标收集
- **Grafana**: 可视化面板
- **AlertManager**: 告警管理

### 日志管理
- **结构化日志**: JSON格式
- **日志级别**: DEBUG/INFO/WARN/ERROR
- **日志轮转**: 按大小和时间轮转

## 🛡️ 安全特性

### 认证授权
- **JWT Token**: 无状态认证
- **RBAC**: 基于角色的访问控制
- **ACL**: 细粒度权限控制

### 数据安全
- **传输加密**: HTTPS/TLS
- **存储加密**: 敏感数据加密
- **访问审计**: 完整操作日志

### 安全防护
- **限流保护**: API访问限制
- **CORS**: 跨域请求控制
- **XSS防护**: 输入输出过滤

## 🤝 贡献指南

### 开发流程
1. Fork 项目
2. 创建功能分支
3. 提交代码
4. 创建 Pull Request

### 代码规范
- **Go**: 遵循 Go 官方规范
- **Vue**: 遵循 Vue 官方风格指南
- **提交**: 使用 Conventional Commits

### 测试要求
- **单元测试**: 覆盖率 > 80%
- **集成测试**: 关键流程测试
- **E2E测试**: 用户场景测试

## 📄 许可证

本项目采用 [MIT License](LICENSE) 许可证。

## 🆘 支持与帮助

### 问题反馈
- **GitHub Issues**: 提交Bug和功能请求
- **讨论区**: 技术交流和问答

### 文档资源
- **API文档**: `/docs/api`
- **部署文档**: `/docs/deployment`
- **开发文档**: `/docs/development`

---

**NetFileSys** - 让文件管理更简单、更安全、更高效！