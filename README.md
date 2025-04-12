# Todo列表后端API

这是一个使用Go语言和Gin框架开发的Todo列表应用后端API。

## 功能特点

- 用户注册和登录
- JWT认证
- **用户隔离**的待办事项CRUD操作 (每个用户只能操作自己的数据)
- 支持批量创建待办事项
- 使用 Redis 缓存优化读取性能

## 技术栈

- Go 1.16+
- Gin Web框架
- GORM ORM框架
- MySQL数据库
- Redis (用于缓存)
- `go-redis/redis/v8` Redis客户端
- `joho/godotenv` (用于加载.env文件)
- JWT认证

## 安装与运行

### 前提条件

- Go 1.16或更高版本
- MySQL数据库 (运行中)
- Redis 服务器 (运行中)

### 安装依赖

```bash
go mod tidy
```

### 配置环境变量

复制`.env.example`文件为`.env`，并根据您的环境修改配置：

```bash
cp .env.example .env
```

编辑`.env`文件，设置您的数据库连接信息、JWT密钥和Redis连接信息：

```dotenv
# 数据库配置
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_NAME=your_db_name

# JWT配置 (请使用强密钥)
JWT_SECRET_KEY=your_strong_jwt_secret_key

# Redis配置
REDIS_ADDR=your_redis_host:your_redis_port
REDIS_PASSWORD=your_redis_password # 如果没有密码则留空
REDIS_DB=0 # Redis数据库编号

# 服务器配置
PORT=8080
```

### 运行应用

1.  确保 MySQL 和 Redis 服务正在运行。
2.  运行命令：
    ```bash
    go run cmd/api/main.go
    ```

服务器将在配置的端口（默认为`8080`）启动。

## API文档

详细的API文档请参考[API接口文档.md](API接口文档.md)。

## 项目结构

```
.
├── cmd
│   └── api
│       └── main.go       # 应用入口, 初始化, 路由
├── handlers
│   ├── todos.go          # 待办事项处理 (包含缓存逻辑)
│   └── users.go          # 用户处理 (注册, 登录, 修改密码)
├── models
│   ├── todo.go           # 待办事项模型, 数据库和Redis初始化
│   └── user.go           # 用户模型
├── .env.example          # 环境变量示例
├── .gitignore            # Git忽略文件
├── go.mod                # Go模块文件
├── go.sum                # Go依赖校验
├── API接口文档.md        # API详细说明
└── README.md             # 项目说明
```

### 环境变量

项目使用`.env`文件或系统环境变量进行配置，主要配置项包括：

- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_HOST`: 数据库主机
- `DB_PORT`: 数据库端口
- `DB_NAME`: 数据库名称
- `JWT_SECRET_KEY`: JWT签名密钥 (请使用强密钥)
- `REDIS_ADDR`: Redis服务器地址 (例如: `localhost:6379`)
- `REDIS_PASSWORD`: Redis密码 (如果需要)
- `REDIS_DB`: Redis数据库编号 (通常是0)
- `PORT`: API服务器监听的端口

## 安全注意事项

- 在生产环境中，请确保使用强密码和安全的JWT密钥。
- **不要**将包含真实凭据的`.env`文件提交到版本控制系统。
- 定期更新依赖包以修复安全漏洞。
- 考虑为 Redis 设置密码保护。

## 许可证

MIT 