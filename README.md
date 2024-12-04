# Go API Mono

[![Go Version](https://img.shields.io/badge/Go-1.22.4-00ADD8?style=flat&logo=go)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

基于 Go 1.22+ 构建的现代化单体 API 服务，采用清晰架构和行业最佳实践。专注于性能、可维护性和开发体验。

## 特性

### 架构设计
- 基于project-layout的清晰架构
- 模块化和可维护的代码结构
- SOLID 原则实现
- 依赖注入模式

### 核心功能
- 现代 Go 1.22+ 特性和惯用法
- RESTful API 设计
- JWT 认证和授权
- 请求速率限制
- 使用 Zap 的结构化日志
- 基于 YAML 的配置管理
- MySQL 数据库支持（使用 GORM）
- 优雅关机处理

### 中间件支持
- 请求/响应日志记录
- 异常恢复处理
- CORS 跨域支持
- 请求 ID 追踪
- JWT 身份认证
- 令牌桶速率限制
- 请求超时处理

## 快速开始

### 环境要求
- Go 1.22.4 或更高版本
- MySQL 8.0 或更高版本
- Docker 和 Docker Compose（可选）

### 本地开发

1. 克隆仓库
```bash
git clone <repository-url>
cd go-api-mono
```

2. 安装依赖
```bash
go mod download
```

3. 运行数据库迁移
```bash
make migrate-up
```

4. 编译项目
```bash
make build
```

5. 运行服务
```bash
make run
```

### Docker 部署

1. 使用 Docker Compose 启动服务
```bash
docker-compose up -d
```

2. 检查服务状态
```bash
docker-compose ps
```

## API 文档

### 用户管理 API

#### 注册新用户
```
POST /api/v1/auth/register
Content-Type: application/json

{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
}
```

#### 用户登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
    "email": "test@example.com",
    "password": "password123"
}
```

#### 获取用户列表
```
GET /api/v1/users
Authorization: Bearer <token>
```

#### 获取用户详情
```
GET /api/v1/users/{id}
Authorization: Bearer <token>
```

#### 更新用户信息
```
PUT /api/v1/users/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
    "username": "updated_user",
    "email": "updated@example.com"
}
```

#### 删除用户
```
DELETE /api/v1/users/{id}
Authorization: Bearer <token>
```

## 项目结构

```
go-api-mono/
├── api/                    # API 文档
│   └── swagger.yaml        # OpenAPI/Swagger 规范
├── cmd/                    # 主程序入口
│   ├── app/               # API 服务器
│   │   └── main.go        # 主程序入口点
│   └── migrate/           # 数据库迁移工具
│       └── main.go        # 迁移程序入口点
├── configs/               # 配置文件
│   ├── config.yaml        # 基础配置
│   ├── config.development.yaml
│   ├── config.production.yaml
│   ├── config.testing.yaml
│   └── embed.go          # 配置文件嵌入
├── internal/              # 内部代码
│   ├── app/              # 应用核心
│   │   └── user/         # 用户模块
│   └── pkg/              # 内部共享包
├── logs/                  # 日志目录
├── scripts/              # 脚本和工具
│   ├── migrations/       # 数据库迁移脚本
│   └── test_api.sh      # API 测试脚本
└── [其他配置文件]
```

## 开发工具

### 可用的 Make 命令

- \`make build\`: 构建应用
- \`make run\`: 运行应用
- \`make test\`: 运行测试
- \`make check\`: 运行所有检查（格式化、lint、测试）
- \`make migrate-up\`: 运行数据库迁移
- \`make migrate-down\`: 回滚数据库迁移
- \`make docker-build\`: 构建 Docker 镜像
- \`make docker-run\`: 运行 Docker 容器

## 配置说明

配置文件位于 \`configs/\` 目录，支持多环境配置：

- \`config.yaml\`: 基础配置
- \`config.development.yaml\`: 开发环境配置
- \`config.production.yaml\`: 生产环境配置
- \`config.testing.yaml\`: 测试环境配置

### 主要配置项

```yaml
app:
  name: "go-api-mono"
  version: "v0.1.0"
  mode: "development"

server:
  port: 8080
  readTimeout: "10s"
  writeTimeout: "10s"

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "123456"
  database: "go_api_mono"

jwt:
  signingKey: "your-secret-key"
  expirationTime: "24h"
  signingMethod: "HS256"
```

## 测试

### 运行测试
```bash
# 运行所有测试
make test

# 运行集成测试
bash scripts/test_api.sh
```

## 部署

### Docker 部署
```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。
