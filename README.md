# BirdNest 🐦

一个基于 Go + Gin 构建的论坛系统,从用户认证系统起步,逐步演进为完整的社区平台。

## 项目简介

BirdNest 是个人练习项目,目标是通过循序渐进的方式,掌握 Go Web 后端开发的核心技术栈——从基础的用户认证,到论坛核心业务,再到部署运维与 AI 能力集成。

## 技术栈

- **语言**: Go
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **配置管理**: Viper
- **认证**: JWT
- **密码加密**: bcrypt
- **邮件服务**: gomail (基于 QQ 邮箱 SMTP)

## 已实现功能

### 用户系统

- [x] 邮箱验证码注册
    - 验证码生成、Redis 存储(5 分钟有效期)
    - HTML 邮件模板,验证码通过邮件发送
    - 验证码校验与用户创建
- [x] 密码 bcrypt 加密存储
- [x] JWT 鉴权中间件
- [x] 配置文件管理(基于 Viper,支持 YAML 配置)

### 工程规范

- [x] 分层架构:`controller` / `model` / `config` / `utils`
- [x] 敏感配置(数据库密码、邮箱授权码等)通过 `config.yaml` 管理,已加入 `.gitignore`
- [x] Git 提交遵循 Angular Commit 规范

## 项目结构

```
bird-nest/
├── cmd/
│   └── main.go          # 程序入口
├── config/               # 配置加载 + 基础设施初始化(DB、Redis)
│   ├── config.go
│   ├── struct.go
│   ├── pgsql.go
│   ├── redis.go
│   └── config.yaml       # 本地配置(不提交,需自行创建)
├── controller/            # 路由处理层
│   └── user.go
├── model/                 # 数据模型
│   └── user.go
├── utils/                 # 工具函数
│   ├── email.go           # 邮件模板
│   └── random.go          # 验证码生成
└── .gitignore
```

## 快速开始

### 1. 准备依赖服务

确保本地已启动 PostgreSQL 和 Redis。

### 2. 配置

复制配置模板并填入自己的参数:

```bash
cp config/config.example.yaml config/config.yaml
```

编辑 `config/config.yaml`,填入数据库连接信息、Redis 地址、邮箱 SMTP 授权码等。

> ⚠️ 邮箱密码请使用 SMTP **授权码**,而非邮箱登录密码。以 QQ 邮箱为例,需在邮箱设置中开启 SMTP 服务后生成授权码。

### 3. 运行

```bash
go mod tidy
go run cmd/main.go
```

服务默认启动在 `:8080`。

## API 一览

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/v1/verification-codes` | 发送邮箱验证码 |
| POST | `/api/v1/users` | 使用验证码完成注册 |



## License

见 [LICENSE](./LICENSE)