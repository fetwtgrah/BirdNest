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
- [x] 密码登录,签发 JWT,并将 token 存入 Redis 做会话校验
- [x] JWT 鉴权中间件(`CheckToken`)
  - 校验 token 合法性
  - 校验 token 是否与 Redis 中记录一致(支持强制登出/单点登录)
  - 解析用户身份并注入请求上下文
- [x] 配置文件管理(基于 Viper,支持 YAML 配置)

### 文章系统

- [x] 发布文章:内容通过 DTO 校验,作者身份从登录态获取,不接受客户端传入
- [x] 查询单篇文章
- [x] 查询当前用户的全部文章
- [x] 更新文章:仅作者本人可修改
- [x] 删除文章:仅作者本人可删除,校验所有权防止越权操作
- [x] User - Passage 一对多关联(GORM 外键约束)

### 工程规范

- [x] 分层架构:`controller` / `model` / `config` / `midware` / `utils`
- [x] 请求参数与数据库模型分离(DTO 模式),避免敏感字段被客户端篡改
- [x] 敏感配置(数据库密码、邮箱授权码、JWT 密钥等)通过 `config.yaml` 管理,已加入 `.gitignore`
- [x] Git 提交遵循 Angular Commit 规范

## 项目结构

```
bird-nest/
├── backend/
│   ├── cmd/
│   │   └── main.go         # 程序入口
│   ├── configs/            # 配置加载 + 基础设施初始化(DB、Redis)
│   │   ├── config.go
│   │   ├── struct.go
│   │   ├── pgsql.go
│   │   ├── redis.go
│   │   └── config.yaml
│   ├── controller/         # 路由处理层
│   ├── midware/            # 中间件
│   ├── model/              # 数据模型
│   ├── utils/              # 工具函数
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── components/     # 通用界面组件
│   │   ├── router/         # 页面路由与登录守卫
│   │   ├── services/       # 后端 API 调用
│   │   └── views/          # 登录、注册、文章发布页面
│   ├── package.json
│   └── vite.config.ts
├── .gitignore
└── README.md
```

## 快速开始

### 1. 准备依赖服务

确保本地已启动 PostgreSQL 和 Redis。

### 2. 配置

编辑 `backend/configs/config.yaml`,填入数据库连接信息、Redis 地址、邮箱 SMTP 授权码、JWT 密钥等。

> ⚠️ 邮箱密码请使用 SMTP **授权码**,而非邮箱登录密码。以 QQ 邮箱为例,需在邮箱设置中开启 SMTP 服务后生成授权码。
> ⚠️ JWT 密钥请使用足够长度的随机字符串,不要使用示例中的默认值。

### 3. 运行

```bash
cd backend
go mod tidy
go run cmd/main.go
```

服务默认启动在 `:8080`。

### 4. 启动前端

另开一个终端:

```bash
cd frontend
npm install
npm run dev
```

前端开发服务器默认运行在 `http://localhost:5173`,并将 `/api` 请求代理到 `http://localhost:8080`。

## API 一览

### 用户相关

| 方法 | 路径 | 说明 | 是否需要鉴权 |
|---|---|---|---|
| POST | `/api/v1/code` | 发送邮箱验证码 | 否 |
| POST | `/api/v1/user` | 使用验证码完成注册 | 否 |
| POST | `/api/v1/user/login` | 密码登录,返回 JWT | 否 |

### 文章相关

| 方法 | 路径 | 说明 | 是否需要鉴权 |
|---|---|---|---|
| POST | `/api/v1/passage` | 发布文章 | 是 |
| GET | `/api/v1/passage/:id` | 获取指定文章 | 是 |
| GET | `/api/v1/passage` | 获取当前用户的全部文章 | 是 |
| PUT | `/api/v1/passage/:id` | 更新文章(仅作者本人) | 是 |
| DELETE | `/api/v1/passage/:id` | 删除文章(仅作者本人) | 是 |

> 需要鉴权的接口,请在请求头中携带登录时获得的 token。

## License

见 [LICENSE](./LICENSE)
