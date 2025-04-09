# Todo列表 API 接口文档

本文档描述了Todo列表应用的API接口，供前端开发人员参考。

## 基础信息

- 基础URL: `http://localhost:8080/api`
- 所有请求和响应均使用JSON格式
- 除了登录和注册外，所有接口都需要认证

## 认证方式

API使用JWT（JSON Web Token）进行认证。

- 在登录后获取token
- 在后续请求中，将token添加到请求头：
  ```
  Authorization: Bearer YOUR_TOKEN_HERE
  ```

## 用户接口

### 1. 用户注册

**请求**

```
POST /register
Content-Type: application/json

{
  "username": "用户名",
  "password": "密码"
}
```

**响应**

- 成功 (201 Created)
```json
{
  "message": "用户创建成功"
}
```

- 失败 (400 Bad Request)
```json
{
  "error": "用户名已存在"
}
```

### 2. 用户登录

**请求**

```
POST /login
Content-Type: application/json

{
  "username": "用户名",
  "password": "密码"
}
```

**响应**

- 成功 (200 OK)
```json
{
  "token": "eyJhbGciOiJIUzI1...",
  "user": {
    "id": 1,
    "username": "用户名",
    "created_at": "2023-04-01T12:00:00Z",
    "updated_at": "2023-04-01T12:00:00Z"
  }
}
```

- 失败 (401 Unauthorized)
```json
{
  "error": "用户名或密码错误"
}
```

### 3. 修改密码 (需要认证)

**请求**

```
POST /change-password
Content-Type: application/json
Authorization: Bearer YOUR_TOKEN_HERE

{
  "old_password": "旧密码",
  "new_password": "新密码"
}
```

**响应**

- 成功 (200 OK)
```json
{
  "message": "密码修改成功"
}
```

- 失败 (401 Unauthorized)
```json
{
  "error": "旧密码不正确"
}
```

## Todo接口

### 1. 获取所有待办事项 (需要认证)

**请求**

```
GET /todos
Authorization: Bearer YOUR_TOKEN_HERE
```

**响应**

- 成功 (200 OK)
```json
[
  {
    "id": 1,
    "title": "学习Go语言",
    "description": "完成Todo列表API项目",
    "completed": false,
    "created_at": "2023-04-01T12:00:00Z",
    "updated_at": "2023-04-01T12:00:00Z"
  },
  {
    "id": 2,
    "title": "学习GORM",
    "description": "掌握数据库操作",
    "completed": true,
    "created_at": "2023-04-01T13:00:00Z",
    "updated_at": "2023-04-01T14:00:00Z"
  }
]
```

### 2. 获取单个待办事项 (需要认证)

**请求**

```
GET /todos/{id}
Authorization: Bearer YOUR_TOKEN_HERE
```

**响应**

- 成功 (200 OK)
```json
{
  "id": 1,
  "title": "学习Go语言",
  "description": "完成Todo列表API项目",
  "completed": false,
  "created_at": "2023-04-01T12:00:00Z",
  "updated_at": "2023-04-01T12:00:00Z"
}
```

- 失败 (404 Not Found)
```json
{
  "error": "待办事项未找到"
}
```

### 3. 创建待办事项 (需要认证)

#### 3.1 创建单个待办事项

**请求**

```
POST /todos
Content-Type: application/json
Authorization: Bearer YOUR_TOKEN_HERE

{
  "todo": {
    "title": "学习Go语言",
    "description": "完成Todo列表API项目"
  }
}
```

**响应**

- 成功 (201 Created)
```json
{
  "id": 1,
  "title": "学习Go语言",
  "description": "完成Todo列表API项目",
  "completed": false,
  "created_at": "2023-04-01T12:00:00Z",
  "updated_at": "2023-04-01T12:00:00Z"
}
```

#### 3.2 批量创建待办事项

**请求**

```
POST /todos
Content-Type: application/json
Authorization: Bearer YOUR_TOKEN_HERE

{
  "todos": [
    {
      "title": "学习Go语言",
      "description": "完成Todo列表API项目"
    },
    {
      "title": "学习GORM",
      "description": "掌握数据库操作"
    }
  ]
}
```

**响应**

- 成功 (201 Created)
```json
{
  "message": "批量创建成功",
  "todos": [
    {
      "id": 1,
      "title": "学习Go语言",
      "description": "完成Todo列表API项目",
      "completed": false,
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    },
    {
      "id": 2,
      "title": "学习GORM",
      "description": "掌握数据库操作",
      "completed": false,
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    }
  ]
}
```

### 4. 更新待办事项 (需要认证)

**请求**

```
PUT /todos/{id}
Content-Type: application/json
Authorization: Bearer YOUR_TOKEN_HERE

{
  "title": "学习Go语言进阶",
  "description": "完成Todo列表API项目并添加新功能",
  "completed": true
}
```

**响应**

- 成功 (200 OK)
```json
{
  "id": 1,
  "title": "学习Go语言进阶",
  "description": "完成Todo列表API项目并添加新功能",
  "completed": true,
  "created_at": "2023-04-01T12:00:00Z",
  "updated_at": "2023-04-01T15:00:00Z"
}
```

- 失败 (404 Not Found)
```json
{
  "error": "待办事项未找到"
}
```

### 5. 删除待办事项 (需要认证)

**请求**

```
DELETE /todos/{id}
Authorization: Bearer YOUR_TOKEN_HERE
```

**响应**

- 成功 (204 No Content)

- 失败 (404 Not Found)
```json
{
  "error": "待办事项未找到"
}
```

## 错误码说明

| 状态码 | 说明 |
|-------|------|
| 200   | 请求成功 |
| 201   | 创建成功 |
| 204   | 删除成功 |
| 400   | 请求参数错误 |
| 401   | 未认证或认证失败 |
| 404   | 资源未找到 |
| 415   | 不支持的媒体类型 |
| 500   | 服务器内部错误 |

## 注意事项

1. Token有效期为24小时，过期后需要重新登录获取新的token
2. 所有时间字段使用ISO 8601格式（如：`2023-04-01T12:00:00Z`）
3. 创建和更新待办事项时，`completed`字段如未提供，默认为`false` 