# Todo列表后端API

这是一个使用Go语言和Gin框架开发的Todo列表应用后端API。

## 功能特点

- 用户注册和登录
- JWT认证
- 待办事项的CRUD操作
- 支持批量创建待办事项

## 技术栈

- Go 1.16+
- Gin Web框架
- GORM ORM框架
- MySQL数据库
- JWT认证

## 安装与运行

### 前提条件

- Go 1.16或更高版本
- MySQL数据库

### 安装依赖

```bash
go mod tidy
```

### 配置环境变量

复制`.env.example`文件为`.env`，并根据您的环境修改配置：

```bash
cp .env.example .env
```

编辑`.env`文件，设置您的数据库连接信息和JWT密钥：

```
# 数据库配置
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_NAME=your_db_name

# JWT配置
JWT_SECRET_KEY=your_jwt_secret_key

# 服务器配置
PORT=8080
```

### 运行应用

```bash
go run cmd/api/main.go
```

服务器将在`http://localhost:8080`启动。

## API文档

详细的API文档请参考[API接口文档.md](API接口文档.md)。

## 开发说明

### 项目结构

```
.
├── cmd
│   └── api
│       └── main.go       # 应用入口
├── handlers
│   ├── todos.go          # 待办事项处理
│   └── users.go          # 用户处理
├── models
│   ├── todo.go           # 待办事项模型
│   └── user.go           # 用户模型
├── .env.example          # 环境变量示例
├── .gitignore            # Git忽略文件
├── go.mod                # Go模块文件
├── go.sum                # Go依赖校验
└── README.md             # 项目说明
```

### 环境变量

项目使用环境变量进行配置，主要配置项包括：

- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_HOST`: 数据库主机
- `DB_PORT`: 数据库端口
- `DB_NAME`: 数据库名称
- `JWT_SECRET_KEY`: JWT签名密钥
- `PORT`: 服务器端口

## 安全注意事项

- 在生产环境中，请确保使用强密码和安全的JWT密钥
- 不要将包含真实凭据的`.env`文件提交到版本控制系统
- 定期更新依赖包以修复安全漏洞

## 许可证

MIT 