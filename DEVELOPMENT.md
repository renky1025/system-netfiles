# 网络文件管理系统开发文档

## 项目概述

本项目是一个基于 Go + Vue3 的企业级网络文件管理系统，采用现代化的技术栈和架构设计，提供完整的文件管理、分享、权限控制等功能。

## 技术栈

### 后端技术栈
- **语言**: Go 1.21+
- **框架**: Gin (HTTP路由)
- **数据库**: PostgreSQL 13+ (主数据库)
- **缓存**: Redis 6+ (缓存和会话)
- **存储**: MinIO (对象存储)
- **ORM**: GORM (数据库操作)
- **认证**: JWT (JSON Web Token)
- **日志**: Zap (结构化日志)
- **配置**: Viper (配置管理)

### 前端技术栈
- **框架**: Vue 3 + TypeScript
- **UI库**: Element Plus
- **构建工具**: Vite
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP客户端**: Axios
- **图表**: ECharts (管理后台)

### 基础设施
- **容器化**: Docker + Docker Compose
- **监控**: Prometheus + Grafana
- **反向代理**: Nginx
- **CI/CD**: GitHub Actions (可选)

## 项目结构

```
netfilessys/
├── backend/                    # 后端服务
│   ├── common/                # 公共模块
│   │   ├── config/           # 配置管理
│   │   ├── database/         # 数据库连接
│   │   ├── middleware/       # 中间件
│   │   └── response/         # 统一响应
│   ├── controller/           # 控制器层
│   ├── service/              # 服务层
│   ├── model/                # 数据模型
│   ├── utils/                # 工具函数
│   ├── scripts/              # 脚本文件
│   ├── main.go              # 主程序入口
│   ├── go.mod               # Go模块定义
│   └── Dockerfile           # Docker构建文件
├── frontend/                  # 前端项目
│   ├── user-portal/          # 用户端
│   └── admin-portal/         # 管理后台
├── database/                  # 数据库相关
│   ├── schemas/              # 数据库表结构
│   └── init/                 # 初始化脚本
├── config/                    # 配置文件
│   ├── config.yaml          # 主配置文件
│   ├── .env.example         # 环境变量模板
│   └── docker-compose.yml   # Docker编排文件
├── docs/                      # 文档目录
├── README.md                 # 项目说明
└── DEVELOPMENT.md            # 开发文档
```

## 核心功能模块

### 1. 用户认证与授权
- **JWT认证**: 无状态token认证
- **RBAC权限**: 基于角色的访问控制
- **ACL权限**: 细粒度对象级权限
- **多因子认证**: 支持邮箱/短信验证
- **密码策略**: 复杂度要求和定期更换

### 2. 文件管理
- **文件上传**: 支持单文件和批量上传
- **秒传技术**: 基于MD5/SHA1的文件去重
- **分片上传**: 大文件断点续传
- **文件预览**: 支持图片、文档预览
- **版本管理**: 文件多版本支持
- **回收站**: 软删除和恢复机制

### 3. 文件夹管理
- **层级结构**: 支持多级文件夹
- **批量操作**: 移动、复制、删除
- **权限继承**: 文件夹权限自动继承
- **路径导航**: 面包屑导航

### 4. 分享功能
- **分享链接**: 生成唯一分享码
- **访问控制**: 密码保护、有效期限制
- **下载限制**: 次数限制、IP白名单
- **访问统计**: 详细的访问日志和统计

### 5. 权限管理
- **角色管理**: 创建、编辑、删除角色
- **权限分配**: 灵活的权限组合
- **组织架构**: 多级组织和部门管理
- **权限继承**: 支持权限继承和覆盖

### 6. 管理后台
- **用户管理**: 用户CRUD、状态管理
- **文件监管**: 全局文件查看和管理
- **分享监管**: 分享记录和控制
- **系统统计**: 各种统计图表和报表
- **审计日志**: 完整的操作记录
- **系统配置**: 系统参数配置

## 数据库设计

### 核心表结构

#### 用户相关表
- `users`: 用户基本信息
- `roles`: 角色定义
- `permissions`: 权限定义
- `user_roles`: 用户角色关联
- `role_permissions`: 角色权限关联
- `organizations`: 组织架构
- `user_organizations`: 用户组织关联

#### 文件相关表
- `files`: 文件元数据
- `folders`: 文件夹信息
- `blobs`: 物理文件存储
- `file_versions`: 文件版本历史
- `file_op_logs`: 文件操作日志
- `recycle_bin`: 回收站记录

#### 分享相关表
- `shares`: 分享记录
- `share_access_logs`: 分享访问日志
- `share_files`: 分享文件关联
- `share_templates`: 分享模板

#### 权限相关表
- `acl_entries`: ACL权限条目
- `user_permission_cache`: 用户权限缓存
- `permission_audit_logs`: 权限审计日志

#### 系统相关表
- `system_configs`: 系统配置
- `login_logs`: 登录日志
- `admin_logs`: 管理员操作日志

## API设计

### RESTful API规范
- 使用标准HTTP方法 (GET, POST, PUT, DELETE)
- 统一的响应格式
- 合理的HTTP状态码
- 分页和排序支持
- 错误处理和消息

### 主要接口分组

#### 认证接口 (/api/v1/auth)
```
POST /login          # 用户登录
POST /logout         # 用户登出
POST /register       # 用户注册
POST /refresh        # 刷新token
POST /forgot-password # 忘记密码
POST /reset-password  # 重置密码
```

#### 文件接口 (/api/v1/files)
```
POST /upload         # 上传文件
GET  /list          # 文件列表
GET  /:id           # 文件详情
PUT  /:id           # 更新文件
DELETE /:id         # 删除文件
GET  /:id/download  # 下载文件
POST /move          # 移动文件
```

#### 分享接口 (/api/v1/shares)
```
POST /              # 创建分享
GET  /              # 分享列表
GET  /:id           # 分享详情
PUT  /:id           # 更新分享
DELETE /:id         # 删除分享
```

#### 管理接口 (/api/v1/admin)
```
GET  /users         # 用户管理
GET  /files         # 文件管理
GET  /shares        # 分享管理
GET  /stats         # 统计信息
GET  /logs          # 审计日志
```

## 安全设计

### 认证安全
- JWT token有效期控制
- Refresh token机制
- 登录失败次数限制
- 异地登录提醒

### 权限安全
- 最小权限原则
- 权限缓存机制
- 权限变更审计
- 敏感操作二次确认

### 数据安全
- 密码加密存储 (bcrypt)
- 敏感数据脱敏
- SQL注入防护
- XSS攻击防护

### 传输安全
- HTTPS强制加密
- CORS跨域控制
- 请求签名验证
- 限流防护

## 性能优化

### 数据库优化
- 合理的索引设计
- 查询语句优化
- 连接池配置
- 读写分离 (可选)

### 缓存策略
- Redis缓存热点数据
- 用户权限缓存
- 文件元数据缓存
- 分享信息缓存

### 文件存储优化
- 对象存储 (MinIO/S3)
- CDN加速 (可选)
- 文件压缩
- 缩略图生成

### 前端优化
- 组件懒加载
- 图片懒加载
- 虚拟滚动
- 请求防抖

## 部署方案

### Docker部署
```bash
# 一键启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f [service-name]
```

### 生产环境部署
1. **负载均衡**: Nginx反向代理
2. **数据库**: PostgreSQL主从复制
3. **缓存**: Redis集群
4. **存储**: MinIO集群或云存储
5. **监控**: Prometheus + Grafana
6. **日志**: ELK Stack (可选)

## 开发规范

### 代码规范
- Go: 遵循官方Go代码规范
- Vue: 遵循Vue官方风格指南
- 统一的代码格式化工具
- 完善的注释和文档

### Git规范
- 使用Conventional Commits
- 功能分支开发
- Code Review流程
- 自动化测试

### 测试规范
- 单元测试覆盖率 > 80%
- 集成测试
- API测试
- 前端E2E测试

## 监控与运维

### 系统监控
- 服务健康检查
- 性能指标监控
- 错误率监控
- 资源使用监控

### 业务监控
- 用户活跃度
- 文件上传下载量
- 分享访问统计
- 存储空间使用

### 告警机制
- 服务异常告警
- 性能阈值告警
- 安全事件告警
- 业务异常告警

## 扩展计划

### 功能扩展
- [ ] 在线编辑功能
- [ ] 协作编辑
- [ ] 文件同步客户端
- [ ] 移动端APP
- [ ] 第三方集成 (钉钉、企微等)

### 技术扩展
- [ ] 微服务架构改造
- [ ] 消息队列集成
- [ ] 搜索引擎集成
- [ ] AI功能集成
- [ ] 区块链存储

## 常见问题

### Q: 如何修改默认管理员密码？
A: 登录后在个人设置中修改，或通过数据库直接更新。

### Q: 如何配置邮件服务？
A: 在config.yaml中配置SMTP相关参数。

### Q: 如何扩展存储容量？
A: 可以配置多个MinIO节点或使用云存储服务。

### Q: 如何备份数据？
A: 定期备份PostgreSQL数据库和MinIO存储数据。

### Q: 如何监控系统状态？
A: 通过Grafana面板查看各项监控指标。

---

更多详细信息请参考项目README和相关文档。