# Go API Mono

一个基于 Go 1.22+ 的单体 API 服务模板，使用标准库 `net/http` 包构建。

## 特性

- 基于 Go 1.22+ 的新特性
- 使标准库 `net/http` 包，无第三方 Web 框架
- 完整的用户认证功能（注册、登录、JWT）
- 清晰的项目结构和分层架构
- Docker 支持，包含完整的开发和部署配置
- 数据库迁移工具
- 完整的测试覆盖
- 支持多环境配置

## 项目结构

```
.
├── api/                    # API 文档
├── cmd/                    # 命令行工具
│   ├── app/               # 主应用入口
│   └── migrate/           # 数据库迁移工具
├── configs/               # 配置文件
├── internal/              # 内部代码
│   ├── app/              # 应用层
│   │   └── user/         # 用户模块
│   └── pkg/              # 公共包
├── scripts/              # 脚本文件
└── docker-compose.yml    # Docker 编排配置
```

## 快速开始

### 本地开发

1. 启动 MySQL：
```bash
docker run -d --name mysql -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=go_api_mono \
  -e MYSQL_USER=apiuser \
  -e MYSQL_PASSWORD=apipass \
  mysql:8.0
```

2. 执行数据库迁移：
```bash
make migrate-up
```

3. 运行服务：
```bash
make run
```

### Docker 环境

使用 Docker Compose 启动所有服务：

```bash
docker-compose up -d
```

## 测试

运行测试脚本：

```bash
bash scripts/test_api.sh
```

## API 文档

### 用户接口

- POST /api/v1/users/register - 用户注册
- POST /api/v1/users/login - 用户登录
- GET /api/v1/users - 获取用户列表
- GET /api/v1/users/{id} - 获取用户详情
- PUT /api/v1/users/{id} - 更新用户信息
- DELETE /api/v1/users/{id} - 删除用户

## 配置说明

项目支持多环境配置：

- `config.development.yaml` - 本地开发环境
- `config.docker.yaml` - Docker 环境
- `config.production.yaml` - 生产环境
- `config.testing.yaml` - 测试环境

## Makefile 命令

- `make run` - 运行服务
- `make build` - 构建二进制文件
- `make test` - 运行测试
- `make migrate-up` - 执行数据库迁移
- `make migrate-down` - 回滚数据库迁移
- `make docker-build` - 构建 Docker 镜像

## 依赖

- Go 1.22+
- MySQL 8.0
- Docker & Docker Compose（可选）

## 许可证

MIT License
