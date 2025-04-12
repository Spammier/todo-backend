# Todo列表 API 接口文档

本文档描述了Todo列表应用的API接口，供前端开发人员参考。

## 基础信息

- 基础URL: `http://localhost:8080/api`
- 所有请求和响应均使用JSON格式
- 除了登录和注册外，所有接口都需要认证，并且**待办事项相关接口仅操作当前认证用户的数据**。

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
  // 可能包含用户信息，根据实际实现调整
}
```

- 失败 (400 Bad Request)
```json
{
  "error": "用户名已存在 或 无效的请求数据"
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
  "error": "旧密码不正确 或 用户未认证"
}
```
- 失败 (400 Bad Request)
```json
{
  "error": "无效的请求数据 或 新密码不能为空"
}
```


## Todo接口 (需要认证，仅操作当前用户数据)

### 1. 获取当前用户的所有待办事项

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
    "user_id": 1, // 注意：此字段通常为内部使用，不一定在API响应中返回
    "title": "学习Go语言",
    "description": "完成Todo列表API项目",
    "completed": false,
    "created_at": "2023-04-01T12:00:00Z",
    "updated_at": "2023-04-01T12:00:00Z"
  },
  {
    "id": 3, 
    "user_id": 1,
    "title": "整理笔记",
    "description": "",
    "completed": false,
    "created_at": "2023-04-02T10:00:00Z",
    "updated_at": "2023-04-02T10:00:00Z"
  }
]
```
- 失败 (500 Internal Server Error)
```json
{
  "error": "获取待办事项失败"
}
```

### 2. 获取当前用户的单个待办事项

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
  "user_id": 1,
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
  "error": "待办事项未找到或无权访问"
}
```
- 失败 (400 Bad Request)
```json
{
  "error": "无效的待办事项ID"
}
```

### 3. 为当前用户创建待办事项

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
    // completed 字段可选, 默认为 false
  }
}
```

**响应**

- 成功 (201 Created)
```json
{
  "id": 1,
  "user_id": 1,
  "title": "学习Go语言",
  "description": "完成Todo列表API项目",
  "completed": false,
  "created_at": "2023-04-01T12:00:00Z",
  "updated_at": "2023-04-01T12:00:00Z"
}
```
- 失败 (400 Bad Request)
```json
{
  "error": "无效的请求数据 或 请求体必须包含 todo 或 todos 字段"
}
```
- 失败 (500 Internal Server Error)
```json
{
  "error": "创建待办事项失败"
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
      "user_id": 1,
      "title": "学习Go语言",
      "description": "完成Todo列表API项目",
      "completed": false,
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    },
    {
      "id": 2,
      "user_id": 1,
      "title": "学习GORM",
      "description": "掌握数据库操作",
      "completed": false,
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    }
  ]
}
```
- 失败 (500 Internal Server Error)
```json
{
  "error": "批量创建待办事项失败"
}
```

### 4. 更新当前用户的待办事项

**请求**

```
PUT /todos/{id}
Content-Type: application/json
Authorization: Bearer YOUR_TOKEN_HERE

{
  "title": "学习Go语言进阶", // 可选
  "description": "完成Todo列表API项目并添加新功能", // 可选
  "completed": true // 可选
}
```

**响应**

- 成功 (200 OK)
```json
{
  "id": 1,
  "user_id": 1,
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
  "error": "待办事项未找到或无权更新"
}
```
- 失败 (400 Bad Request)
```json
{
  "error": "无效的请求数据"
}
```
- 失败 (500 Internal Server Error)
```json
{
  "error": "更新待办事项失败"
}
```

### 5. 删除当前用户的待办事项

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
  "error": "待办事项未找到或无权删除"
}
```
- 失败 (500 Internal Server Error)
```json
{
  "error": "删除待办事项失败"
}
```

## 错误码说明

| 状态码 | 说明 | 
|-------|------|
| 200   | 请求成功 |
| 201   | 创建成功 |
| 204   | 删除成功 (无内容返回) |
| 400   | 请求参数错误 (Bad Request) |
| 401   | 未认证或认证失败 (Unauthorized) |
| 403   | 无权限访问 (Forbidden) |
| 404   | 资源未找到 (Not Found) |
| 415   | 不支持的媒体类型 (Unsupported Media Type) |
| 500   | 服务器内部错误 (Internal Server Error) |

## 注意事项

1. Token有效期为24小时，过期后需要重新登录获取新的token。
2. 所有时间字段使用ISO 8601格式（如：`2023-04-01T12:00:00Z`）。
3. 创建和更新待办事项时，`completed`字段如未提供，默认为`false`。
4. 所有待办事项操作（增删改查）都与当前认证用户绑定。
5. API响应中的`user_id`字段仅作示例，实际可能不返回。 