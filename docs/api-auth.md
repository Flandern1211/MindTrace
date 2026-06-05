# YoudaoNoteLM 认证模块接口文档

> 基础路径：`/api/v1`
>
> Content-Type：`application/json`

---

## 通用响应格式

所有接口统一返回以下结构：

```json
{
  "code": 0,
  "message": "成功",
  "data": {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码，`0` 表示成功，非 0 为错误码 |
| message | string | 提示信息 |
| data | object/null | 业务数据，部分接口可能为空 |

---

## 错误码一览

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 500 | 服务器内部错误（如 Redis 服务异常） |
| 1001 | 用户不存在 |
| 1002 | 用户已存在（邮箱已被注册） |
| 1003 | 邮箱或密码错误 |
| 1004 | 用户已被禁用 |
| 1005 | 无效的令牌（含已吊销的令牌） |
| 1006 | 令牌已过期 |
| 1007 | 账户已被锁定，请15分钟后重试 |
| 1101 | 验证码已过期，请重新获取 |
| 1102 | 验证码错误 |
| 1103 | 验证码输入错误次数过多，请重新获取 |
| 1104 | 验证码发送过于频繁，请60秒后重试 |
| 2001 | 参数错误（如滑块验证失败） |

---

## 1. 获取滑块验证码

登录前必须先获取滑块验证码。

**请求**

```
GET /api/v1/auth/captcha
```

**响应**

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "captcha_id": "aBcDeFgHiJkLmNoP",
    "background": "data:image/png;base64,iVBORw0KGgo...",
    "slider": "data:image/png;base64,iVBORw0KGgo...",
    "slider_size": 40,
    "bg_width": 300,
    "slider_start_x": 0
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| captcha_id | string | 验证码唯一标识，登录时需回传 |
| background | string | 背景图的 base64（含凹槽阴影），直接作为 `<img>` 的 src |
| slider | string | 滑块图的 base64，需叠加在背景图上供用户拖动 |
| slider_size | int | 滑块边长（像素），宽高相等 |
| bg_width | int | 背景图总宽度（像素） |
| slider_start_x | int | 滑块初始 X 坐标（相对背景图左边缘），目前固定为 0 |

**⚠️ captcha_x 坐标说明（重要）**

`captcha_x` 是**滑块左边缘**在背景图上的**绝对 X 坐标**（像素），不是拖拽距离、不是百分比。

```
背景图 (宽 300px)
┌──────────────────────────────────────────────┐
│  滑块起始位置(0)         凹槽位置(correctX)    │
│     ↓                       ↓                 │
│    ┌───┐                                     │
│    │   │ ← slider                             │
│    └───┘                                      │
│     |---------- drag_distance ------------>|  │
│                                              │
│  captcha_x = slider_start_x + drag_distance  │
│  captcha_x = 0 + drag_distance               │
│  captcha_x = drag_distance                   │
└──────────────────────────────────────────────┘
```

- 滑块从 `slider_start_x`（目前固定为 0）开始拖动
- 用户拖动了 `drag_distance` 像素
- **`captcha_x = slider_start_x + drag_distance`**，即滑块左边缘在背景图上的绝对位置
- 后端容差 ±5 像素

**前端交互说明**

1. 渲染背景图（宽 `bg_width`），将滑块图叠放在 `slider_start_x` 位置
2. 用户拖动滑块到凹槽位置后松手
3. 计算 `captcha_x = slider_start_x + 拖拽偏移量`
4. 将 `captcha_id` 和 `captcha_x` 提交到登录接口

**前端示例代码**

```javascript
// 假设使用 sliderLeft 记录滑块左边缘的 CSS left 值
const captcha_x = startX + dragDistance;  // 滑块左边缘在背景图上的绝对像素位置

// ❌ 错误：传百分比
// captcha_x = 52  (52%)

// ❌ 错误：传滑块中心位置
// captcha_x = dragDistance + sliderSize / 2

// ✅ 正确：传滑块左边缘的绝对像素位置
// captcha_x = dragDistance  (因为 slider_start_x = 0)
```

---

## 2. 用户登录

**请求**

```
POST /api/v1/auth/login
```

```json
{
  "email": "user@example.com",
  "password": "myPassword123",
  "captcha_id": "aBcDeFgHiJkLmNoP",
  "captcha_x": 156
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| email | string | ✅ | 邮箱地址 |
| password | string | ✅ | 密码 |
| captcha_id | string | ✅ | 从 `/captcha` 接口获取的验证码 ID |
| captcha_x | int | ✅ | 滑块左边缘在背景图上的**绝对 X 坐标**（像素），计算方式：`slider_start_x + 拖拽偏移量` |

**成功响应**

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "user",
      "email": "user@example.com",
      "nickname": "",
      "avatar": "",
      "status": 1,
      "created_at": "2026-06-04T10:00:00Z",
      "updated_at": "2026-06-04T10:00:00Z"
    }
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| access_token | string | 访问令牌，有效期 **15 分钟**，请求业务接口时放入 `Authorization: Bearer <token>` |
| refresh_token | string | 刷新令牌，有效期 **1 天**，用于无感刷新 access_token |
| user | object | 用户信息 |

**错误场景**

| code | message | 原因 |
|------|---------|------|
| 2001 | 验证码已过期，请重新获取 | captcha_id 无效或已过期 |
| 2001 | 滑块验证失败，请重试 | 用户拖动位置偏差过大 |
| 1003 | 邮箱或密码错误 | 账号不存在或密码错误 |
| 1004 | 用户已被禁用 | 账号被管理员禁用 |
| 1007 | 账户已被锁定，请15分钟后重试 | 连续输错密码 3 次 |

---

## 3. 刷新 Token

access_token 过期后，使用 refresh_token 换取新的 token 对。

**请求**

```
POST /api/v1/auth/refresh
```

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| refresh_token | string | ✅ | 登录时返回的刷新令牌 |

**成功响应**（结构与登录响应一致）

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...(新的)",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...(新的)",
    "user": {
      "id": 1,
      "username": "user",
      "email": "user@example.com",
      "nickname": "",
      "avatar": "",
      "status": 1,
      "created_at": "2026-06-04T10:00:00Z",
      "updated_at": "2026-06-04T10:00:00Z"
    }
  }
}
```

**错误场景**

| code | message | 原因 |
|------|---------|------|
| 1005 | 无效的令牌 | refresh_token 格式错误或被篡改 |
| 1006 | 令牌已过期 | refresh_token 已超过 1 天有效期 |
| 1005 | 请使用 refresh_token 进行刷新 | 误传了 access_token |
| 1005 | refresh token 已失效，请重新登录 | refresh_token 已被吊销（已登出或已刷新过） |

---

## 4. 用户登出

将当前 access_token 和 refresh_token 加入服务端黑名单，使其立即失效。登出后，这两个 token 将无法再使用。

**请求**

```
POST /api/v1/auth/logout
```

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| access_token | string | ❌ | 当前的访问令牌，传入后立即失效 |
| refresh_token | string | ❌ | 当前的刷新令牌，传入后立即失效 |

> 两个字段至少传一个，建议同时传入以确保完全登出。

**成功响应**

```json
{
  "code": 0,
  "message": "成功",
  "data": null
}
```

**错误场景**

| code | message | 原因 |
|------|---------|------|
| 500 | 吊销 access token 失败 | Redis 服务异常 |
| 500 | 吊销 refresh token 失败 | Redis 服务异常 |

**前端登出推荐流程**

```
1. 用户点击「退出登录」
2. POST /auth/logout  { access_token, refresh_token }
3. 无论成功与否，清除本地存储的 token
4. 跳转登录页
```

---

## 5. 发送邮箱验证码

用于注册和找回密码两个场景。

**请求**

```
POST /api/v1/auth/send-code
```

```json
{
  "email": "user@example.com",
  "type": "register"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| email | string | ✅ | 邮箱地址 |
| type | string | ✅ | 场景：`register`（注册）或 `reset`（找回密码） |

**成功响应**

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "retry_after": 60
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| retry_after | int | 距离下次可发送的冷却秒数（60 秒） |

**错误场景**

| code | message | 原因 |
|------|---------|------|
| 1002 | 邮箱已被注册 | type=register 时，该邮箱已注册过 |
| 1001 | 用户不存在 | type=reset 时，该邮箱未注册 |
| 1104 | 验证码发送过于频繁，请60秒后重试 | 60 秒内重复请求 |

---

## 6. 用户注册

**请求**

```
POST /api/v1/auth/register
```

```json
{
  "email": "user@example.com",
  "password": "myPassword123",
  "confirm_password": "myPassword123",
  "code": "382916"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| email | string | ✅ | 邮箱地址，不可重复 |
| password | string | ✅ | 密码，8-20 位，至少包含字母和数字 |
| confirm_password | string | ✅ | 确认密码，必须与 password 一致 |
| code | string | ✅ | 6 位邮箱验证码 |

**成功响应**

```json
{
  "code": 0,
  "message": "成功",
  "data": null
}
```

**错误场景**

| code | message | 原因 |
|------|---------|------|
| 400 | 两次密码输入不一致 | password 与 confirm_password 不一致 |
| 1101 | 验证码已过期，请重新获取 | 验证码超过 5 分钟有效期 |
| 1102 | 验证码错误 | 验证码输入错误 |
| 1103 | 验证码输入错误次数过多，请重新获取 | 验证码累计输错 5 次 |
| 1002 | 邮箱已被注册 | 该邮箱已注册过 |

---

## 7. 找回密码（重置密码）

**请求**

```
POST /api/v1/auth/reset-password
```

```json
{
  "email": "user@example.com",
  "code": "591738",
  "new_password": "newPass123456",
  "confirm_password": "newPass123456"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| email | string | ✅ | 注册时使用的邮箱 |
| code | string | ✅ | 6 位邮箱验证码 |
| new_password | string | ✅ | 新密码，8-20 位 |
| confirm_password | string | ✅ | 确认新密码，必须与 new_password 一致 |

**成功响应**

```json
{
  "code": 0,
  "message": "成功",
  "data": null
}
```

**错误场景**

| code | message | 原因 |
|------|---------|------|
| 400 | 两次密码输入不一致 | new_password 与 confirm_password 不一致 |
| 1101 | 验证码已过期，请重新获取 | 验证码超过 5 分钟有效期 |
| 1102 | 验证码错误 | 验证码输入错误 |
| 1103 | 验证码输入错误次数过多，请重新获取 | 验证码累计输错 5 次 |
| 1001 | 用户不存在 | 该邮箱未注册 |

---

## 8. 上传头像

上传用户头像，支持 jpg/jpeg/png 格式，最大 2MB。同一用户重复上传会自动覆盖旧头像文件。

**请求**

```
POST /api/v1/user/avatar
Content-Type: multipart/form-data
Authorization: Bearer <access_token>
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| avatar | file | ✅ | 头像文件（jpg/jpeg/png，≤ 2MB） |

**成功响应**

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "avatar": "/uploads/avatars/1.jpg"
  }
}
```

> 头像 URL 格式：`/uploads/avatars/{user_id}.{ext}`，可通过 `GET /uploads/avatars/{user_id}.{ext}` 访问。

**错误场景**

| code | message | 原因 |
|------|---------|------|
| 400 | 请上传头像文件 | 未选择文件 |
| 400 | 头像文件大小不能超过 2MB | 文件过大 |
| 400 | 仅支持 jpg/jpeg/png 格式 | 文件格式不支持 |

---

## 9. 修改用户名

修改当前登录用户的用户名。

**请求**

```
PUT /api/v1/user/username
Authorization: Bearer <access_token>
```

```json
{
  "username": "new_username"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| username | string | ✅ | 新用户名，3-50 位 |

**成功响应**

```json
{
  "code": 0,
  "message": "成功",
  "data": null
}
```

**错误场景**

| code | message | 原因 |
|------|---------|------|
| 1002 | 用户名已被使用 | 该用户名已存在 |

---

## 10. 修改密码

修改当前登录用户的密码。修改成功后，当前 token 会失效，需要重新登录。

**请求**

```
POST /api/v1/user/password
Authorization: Bearer <access_token>
```

```json
{
  "old_password": "oldPass123",
  "new_password": "newPass456"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| old_password | string | ✅ | 当前密码 |
| new_password | string | ✅ | 新密码，8-20 位 |

**成功响应**

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "message": "密码修改成功，请重新登录"
  }
}
```

**错误场景**

| code | message | 原因 |
|------|---------|------|
| 1003 | 邮箱或密码错误 | 旧密码不正确 |

---

## Token 使用说明

### 双 Token 机制

| Token | 有效期 | 用途 | 存储建议 |
|-------|--------|------|----------|
| access_token | 15 分钟 | 请求业务接口时的身份凭证 | 内存 / sessionStorage |
| refresh_token | 1 天 | 静默刷新 access_token | localStorage |

### 请求业务接口

所有需要登录的接口（如用户信息、笔记管理等），在请求头中携带 access_token：

```
Authorization: Bearer <access_token>
```

### Token 黑名单机制

服务端实现了基于 Redis 的 Token 黑名单机制：

- **登出时**：access_token 和 refresh_token 被加入黑名单，立即失效
- **刷新时**：旧的 refresh_token 被加入黑名单，防止重放攻击
- **自动清理**：黑名单记录在 token 原定过期时间后自动清除，不占用额外存储
- **影响**：已吊销的 token 请求接口会返回 `code=1005`，需重新登录

### Token 刷新流程（前端推荐）

```
1. 请求业务接口 → 返回 code=1005 或 1006
2. 用 refresh_token 调用 /auth/refresh
3. 成功 → 更新本地 access_token 和 refresh_token，重发原请求
4. 失败（refresh_token 已过期或已吊销）→ 清除登录态，跳转登录页
```

### 登出流程

```
1. 用户点击「退出登录」
2. 调用 POST /auth/logout { access_token, refresh_token }
3. 清除本地存储的所有 token
4. 跳转登录页
```

---

## 前端接入流程总览

### 注册流程

```
1. 用户输入邮箱 → 点击「获取验证码」
   POST /auth/send-code  { email, type: "register" }
2. 用户填写验证码 + 密码 + 确认密码 → 点击「注册」
   POST /auth/register  { email, password, confirm_password, code }
3. 注册成功 → 跳转登录页
```

### 登录流程

```
1. 调用获取滑块验证码
   GET /auth/captcha
2. 渲染滑块组件，用户拖动完成后
   POST /auth/login  { email, password, captcha_id, captcha_x }
3. 登录成功 → 存储 token → 进入主页
```

### 登出流程

```
1. 用户点击「退出登录」
   POST /auth/logout  { access_token, refresh_token }
2. 清除本地存储的 token
3. 跳转登录页
```

### 找回密码流程

```
1. 用户输入邮箱 → 点击「获取验证码」
   POST /auth/send-code  { email, type: "reset" }
2. 用户填写验证码 + 新密码 + 确认密码 → 点击「重置」
   POST /auth/reset-password  { email, code, new_password, confirm_password }
3. 重置成功 → 跳转登录页
```
