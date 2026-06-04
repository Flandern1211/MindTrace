# 表结构设计：笔记本、会话、消息、父块

> 日期：2026-06-04
> 状态：已确认

---

## 1. 设计概述

### 1.1 范围

本次设计四张核心表：
- **notebooks**（笔记本）— 用户的知识管理单元
- **conversations**（会话）— 笔记本下的 AI 对话
- **messages**（消息）— 对话中的单条消息
- **parent_blocks**（父块）— 资料来源的父子分块中的父块

资料来源表（sources）和笔记表（notes）不在本次设计范围。

### 1.2 设计原则

- **扁平关联**：表之间通过外键直接关联，不做过度抽象
- **int 主键**：所有 id 类字段统一使用 int 类型
- **软删除**：继承 BaseEntity，使用 gorm.DeletedAt 实现软删除
- **父块 + Milvus 子块**：父块存 MySQL 原文，子块向量存 Milvus，通过 `parent_block_id` 关联

### 1.3 表关系

```
users
  └── notebooks (user_id)
        ├── conversations (notebook_id)
        │     └── messages (conversation_id)
        └── sources (notebook_id)  ← 另行设计
              └── parent_blocks (source_id)
```

---

## 2. 表结构详情

### 2.1 notebooks（笔记本表）

```go
type Notebook struct {
    BaseEntity
    UserID int    `gorm:"index;not null;comment:所属用户ID"`
    Name   string `gorm:"type:varchar(100);not null;comment:笔记本名称"`
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | int | PK, AUTO_INCREMENT | 主键 |
| user_id | int | NOT NULL, INDEX | 所属用户 ID |
| name | varchar(100) | NOT NULL | 笔记本名称 |
| created_at | datetime | NOT NULL | 创建时间 |
| updated_at | datetime | NOT NULL | 更新时间 |
| deleted_at | datetime | INDEX | 软删除时间 |

**索引**：
- `idx_notebooks_user_id` ON (user_id)

### 2.2 conversations（会话表）

```go
type Conversation struct {
    BaseEntity
    NotebookID int    `gorm:"index;not null;comment:所属笔记本ID"`
    Title      string `gorm:"type:varchar(100);not null;comment:会话标题"`
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | int | PK, AUTO_INCREMENT | 主键 |
| notebook_id | int | NOT NULL, INDEX | 所属笔记本 ID |
| title | varchar(100) | NOT NULL | 会话标题，默认取首条消息前 20 字符 |
| created_at | datetime | NOT NULL | 创建时间 |
| updated_at | datetime | NOT NULL | 更新时间（用于排序） |
| deleted_at | datetime | INDEX | 软删除时间 |

**索引**：
- `idx_conversations_notebook_id` ON (notebook_id)

### 2.3 messages（消息表）

```go
type Message struct {
    BaseEntity
    ConversationID int    `gorm:"index;not null;comment:所属会话ID"`
    Role           string `gorm:"type:varchar(20);not null;comment:角色:user/assistant/system"`
    Content        string `gorm:"type:text;not null;comment:消息内容"`
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | int | PK, AUTO_INCREMENT | 主键 |
| conversation_id | int | NOT NULL, INDEX | 所属会话 ID |
| role | varchar(20) | NOT NULL | 角色：user / assistant / system |
| content | text | NOT NULL | 消息内容 |
| created_at | datetime | NOT NULL | 创建时间（用于排序） |
| updated_at | datetime | NOT NULL | 更新时间 |
| deleted_at | datetime | INDEX | 软删除时间 |

**索引**：
- `idx_messages_conversation_id` ON (conversation_id)

### 2.4 parent_blocks（父块表）

```go
type ParentBlock struct {
    BaseEntity
    SourceID   int    `gorm:"index;not null;comment:所属资料来源ID"`
    Heading    string `gorm:"type:varchar(255);comment:父块标题/小标题"`
    Content    string `gorm:"type:text;not null;comment:父块原文内容"`
    ChunkIndex int    `gorm:"not null;comment:父块在来源中的序号(从0开始)"`
    Metadata   string `gorm:"type:json;comment:元数据JSON(页码/章节等)"`
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | int | PK, AUTO_INCREMENT | 主键 |
| source_id | int | NOT NULL, INDEX | 所属资料来源 ID |
| heading | varchar(255) | NULL | 父块标题/小标题（章节名等） |
| content | text | NOT NULL | 父块原文内容 |
| chunk_index | int | NOT NULL | 父块在来源中的序号（从 0 开始） |
| metadata | json | NULL | 元数据（页码、章节等） |
| created_at | datetime | NOT NULL | 创建时间 |
| updated_at | datetime | NOT NULL | 更新时间 |
| deleted_at | datetime | INDEX | 软删除时间 |

**索引**：
- `idx_parent_blocks_source_id` ON (source_id)

**与 Milvus 子块的关联**：
- Milvus 子块集合中存储 `parent_block_id` 字段
- 检索流程：Milvus 语义检索 → 得到子块 → 通过 `parent_block_id` 回查 MySQL 获取父块原文

---

## 3. 需要同步修改的内容

### 3.1 BaseEntity 主键类型

当前 `BaseEntity.ID` 为 `uint`，需要改为 `int`：

```go
type BaseEntity struct {
    ID        int            `gorm:"primarykey;autoIncrement" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### 3.2 User 实体

`User.ID` 继承自 BaseEntity，会自动变为 int。相关引用字段（如 `FailedAttempts` 的关联逻辑）需检查兼容性。

### 3.3 现有代码影响

- `pkg/jwt/claims.go`：`UserID` 字段从 `uint` 改为 `int`
- `internal/middleware/auth.go`：`GetUserID` 返回值从 `uint` 改为 `int`
- `internal/repository/`：相关 Repository 方法签名需同步修改

---

## 4. Milvus 子块集合设计（参考）

子块向量集合建议结构：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT64 | 主键 |
| parent_block_id | INT64 | 关联 MySQL 的 parent_blocks.id |
| source_id | INT64 | 所属资料来源 ID |
| notebook_id | INT64 | 所属笔记本 ID（冗余，便于过滤） |
| chunk_index | INT64 | 子块在父块中的序号 |
| content | VARCHAR | 子块文本内容 |
| embedding | FLOAT_VECTOR | 向量表示 |

**检索流程**：
1. 用户提问 → 向量化
2. Milvus 语义检索 Top-K 子块
3. 通过 `parent_block_id` 回查 MySQL 获取父块原文
4. 将父块原文作为上下文传给 LLM

---

## 5. 文件清单

新建文件：
- `internal/model/entity/notebook.go`
- `internal/model/entity/conversation.go`
- `internal/model/entity/message.go`
- `internal/model/entity/parent_block.go`

修改文件：
- `internal/model/entity/base.go`（ID 类型改为 int）
- `internal/model/entity/user.go`（同步 ID 类型）
- `pkg/jwt/claims.go`（UserID 类型）
- `internal/middleware/auth.go`（GetUserID 返回类型）
- `internal/app/app.go`（AutoMigrate 添加新表）
