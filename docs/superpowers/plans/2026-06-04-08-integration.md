# 最终集成实施计划

> **智能体执行指南：** 必须使用子技能：superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务执行本计划。

**目标：** 将所有已完成的模块接入应用路由和依赖注入，完成数据库迁移，验证整体构建。

**架构：** 更新 Router 注册所有新控制器，在 App.initDependencies() 中组装所有依赖，更新 AutoMigrate 添加新实体表。

**技术栈：** Gin、GORM

**前置依赖：** 所有前序计划（01-07）必须完成

---

## 文件结构

### 修改文件

```
internal/api/router.go                              → 新增所有控制器和路由注册
internal/app/app.go                                 → 新增依赖注入和数据库迁移
```

---

## 任务 1：更新 Router

**文件：**
- 修改：`internal/api/router.go`

- [ ] **步骤 1：更新 Router 结构体和路由注册**

将 `internal/api/router.go` 替换为：

```go
package api

import (
	"YoudaoNoteLm/internal/api/v1/admin"
	"YoudaoNoteLm/internal/api/v1/auth"
	"YoudaoNoteLm/internal/api/v1/importn"
	"YoudaoNoteLm/internal/api/v1/search"
	"YoudaoNoteLm/internal/api/v1/source"
	"YoudaoNoteLm/internal/api/v1/user"
	"YoudaoNoteLm/internal/api/v1/user_config"
	"YoudaoNoteLm/internal/api/v1/youdao"
	"YoudaoNoteLm/internal/middleware"
	"YoudaoNoteLm/internal/service"
	"github.com/gin-gonic/gin"
)

type Router struct {
	userCtrl    *user.Controller
	authCtrl    *auth.Controller
	sourceCtrl  *source.Controller
	importCtrl  *importn.Controller
	searchCtrl  *search.Controller
	adminCtrl   *admin.Controller
	youdaoCtrl  *youdao.Controller
	userCfgCtrl *user_config.Controller
}

func NewRouter(
	userService service.UserService,
	authService service.AuthService,
	sourceService service.SourceService,
	importerService service.ImporterService,
	searchService service.SearchAgentService,
	adminService service.AdminService,
	youdaoService service.YoudaoAgentService,
	userConfigService service.UserConfigService,
) *Router {
	return &Router{
		userCtrl:    user.NewController(userService),
		authCtrl:    auth.NewController(authService, userService),
		sourceCtrl:  source.NewController(sourceService),
		importCtrl:  importn.NewController(importerService),
		searchCtrl:  search.NewController(searchService),
		adminCtrl:   admin.NewController(adminService),
		youdaoCtrl:  youdao.NewController(youdaoService),
		userCfgCtrl: user_config.NewController(userConfigService),
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	// 全局中间件
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	// 健康检查
	engine.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "YoudaoNoteLM API is running",
		})
	})

	// API v1 路由组
	v1 := engine.Group("/api/v1")
	{
		// 公开路由
		r.authCtrl.RegisterRoutes(v1)

		// 需要认证的路由
		authRequired := v1.Group("")
		authRequired.Use(middleware.Auth())
		{
			r.userCtrl.RegisterRoutes(authRequired)
			r.sourceCtrl.RegisterRoutes(authRequired)
			r.importCtrl.RegisterRoutes(authRequired)
			r.searchCtrl.RegisterRoutes(authRequired)
			r.adminCtrl.RegisterRoutes(authRequired)
			r.youdaoCtrl.RegisterRoutes(authRequired)
			r.userCfgCtrl.RegisterRoutes(authRequired)
		}
	}
}
```

- [ ] **步骤 2：验证编译**

```bash
go build ./internal/api/...
```

- [ ] **步骤 3：提交**

```bash
git add internal/api/router.go
git commit -m "feat(wiring): 更新 Router 注册所有新模块路由"
```

---

## 任务 2：更新 App 依赖注入

**文件：**
- 修改：`internal/app/app.go`

- [ ] **步骤 1：更新 initDependencies() 和 initDatabase()**

在 `internal/app/app.go` 中：

1. 添加新的 import：
```go
import (
	// ... existing imports ...
	"YoudaoNoteLm/internal/service/external"
)
```

2. 替换 `initDatabase()` 中的 AutoMigrate：
```go
if err := a.mysqlDB.AutoMigrate(
	&entity.User{},
	&entity.Source{},
	&entity.ImportTask{},
	&entity.AudioPreview{},
	&entity.SysConfig{},
	&entity.YoudaoBinding{},
	&entity.UserSearchConfig{},
	&entity.UserASRConfig{},
	&entity.UserEmbeddingConfig{},
); err != nil {
	logger.Warn("数据库迁移警告", zap.Error(err))
} else {
	logger.Info("数据库迁移完成")
}
```

3. 替换 `initDependencies()`：
```go
func (a *App) initDependencies() {
	// === 仓储 ===
	userRepo := repository.NewUserRepository(a.mysqlDB)
	sourceRepo := repository.NewSourceRepository(a.mysqlDB)
	taskRepo := repository.NewImportTaskRepository(a.mysqlDB)
	previewRepo := repository.NewAudioPreviewRepository(a.mysqlDB)
	sysConfigRepo := repository.NewSysConfigRepository(a.mysqlDB)
	youdaoBindingRepo := repository.NewYoudaoBindingRepository(a.mysqlDB)
	userSearchCfgRepo := repository.NewUserSearchConfigRepository(a.mysqlDB)
	userASRCfgRepo := repository.NewUserASRConfigRepository(a.mysqlDB)
	userEmbedCfgRepo := repository.NewUserEmbeddingConfigRepository(a.mysqlDB)

	// === 外部服务 ===
	// TODO: 从配置中读取这些地址
	markitdownClient := external.NewMarkitdownClient("http://localhost:8081")
	fileStorage := external.NewMinIOStorage("localhost:9000", "", "", "youdaonotelm")
	asrService := external.NewASRService("http://localhost:8082", "")

	// === 服务 ===
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo, userSvc)
	sourceSvc := service.NewSourceService(sourceRepo)
	importerSvc := service.NewImporterService(markitdownClient, asrService, fileStorage, sourceRepo, taskRepo, previewRepo, nil)
	configSvc := service.NewConfigService(sysConfigRepo, userSearchCfgRepo, userASRCfgRepo)
	searchAgentSvc := service.NewSearchAgentService(configSvc, importerSvc)
	adminSvc := service.NewAdminService(userRepo, sysConfigRepo)
	youdaoAgentSvc := service.NewYoudaoAgentService(nil, youdaoBindingRepo, importerSvc)
	userCfgSvc := service.NewUserConfigService(userSearchCfgRepo, userASRCfgRepo, userEmbedCfgRepo)

	// === 路由 ===
	a.router = api.NewRouter(
		userSvc, authSvc, sourceSvc, importerSvc,
		searchAgentSvc, adminSvc, youdaoAgentSvc, userCfgSvc,
	)
}
```

- [ ] **步骤 2：验证完整构建**

```bash
go build ./cmd/server
```

预期：编译成功，无错误。

- [ ] **步骤 3：验证 vet**

```bash
go vet ./...
```

预期：无警告。

- [ ] **步骤 4：提交**

```bash
git add internal/app/app.go
git commit -m "feat(wiring): 完成所有模块依赖注入和数据库迁移"
```

---

## 任务 3：最终验证

- [ ] **步骤 1：完整构建验证**

```bash
go build ./...
```

- [ ] **步骤 2：vet 检查**

```bash
go vet ./...
```

- [ ] **步骤 3：列出所有新增 API 端点**

```
公开路由：
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh

认证路由：
GET    /api/v1/user/profile
PUT    /api/v1/user/profile
POST   /api/v1/user/password
GET    /api/v1/user/list

GET    /api/v1/notebooks/:nbId/sources
GET    /api/v1/sources/:id
PUT    /api/v1/sources/:id
DELETE /api/v1/sources/:id
POST   /api/v1/sources/batch-delete
GET    /api/v1/sources/:id/content
GET    /api/v1/sources/:id/original

POST   /api/v1/notebooks/:nbId/import/file
POST   /api/v1/notebooks/:nbId/import/audio/preview
POST   /api/v1/import/audio/confirm
GET    /api/v1/import/tasks/:taskId

POST   /api/v1/notebooks/:nbId/search
POST   /api/v1/notebooks/:nbId/search/import

GET    /api/v1/admin/users
PUT    /api/v1/admin/users/:id/status
GET    /api/v1/admin/config/status
GET    /api/v1/admin/config/:group
POST   /api/v1/admin/config/:group
PUT    /api/v1/admin/config/:group/:key

POST   /api/v1/youdao/bind
DELETE /api/v1/youdao/unbind
GET    /api/v1/youdao/status
GET    /api/v1/youdao/notes
POST   /api/v1/youdao/notes/preview
POST   /api/v1/notebooks/:nbId/youdao/import

GET    /api/v1/user/config/search
POST   /api/v1/user/config/search
PUT    /api/v1/user/config/search/:id
DELETE /api/v1/user/config/search/:id
GET    /api/v1/user/config/asr
POST   /api/v1/user/config/asr
PUT    /api/v1/user/config/asr/:id
DELETE /api/v1/user/config/asr/:id
GET    /api/v1/user/config/embedding
POST   /api/v1/user/config/embedding
PUT    /api/v1/user/config/embedding/:id
DELETE /api/v1/user/config/embedding/:id
```

- [ ] **步骤 4：最终提交**

```bash
git add -A
git commit -m "feat: 完成导入、搜索、有道、后台管理模块全部实现"
```

---

## 已知 TODO 项

1. **Eino Agent SDK 集成** — 搜索和有道 Agent 已提供桩代码，需接入实际 SDK
2. **MinIO SDK** — 安装 `github.com/minio/minio-go` 并实现实际上传/下载/删除
3. **EmbeddingService** — 由其他模块实现，当前传入 nil
4. **YoudaoNoteSkill** — 外部依赖，需实现有道 API 客户端
5. **配置文件** — 在 `configs/config.yaml` 中添加 MinIO、MarkItDown、ASR 地址配置
6. **API Key 加密** — 有道 API Key 和用户配置 API Key 应使用 AES-256 加密存储

---

## 计划文件索引

| 文件 | 内容 |
|------|------|
| `2026-06-04-01-foundation.md` | 实体 + 仓储 + 错误码 |
| `2026-06-04-02-source.md` | 资料来源管理 |
| `2026-06-04-03-importer.md` | 导入模块 |
| `2026-06-04-04-config-service.md` | 配置服务 |
| `2026-06-04-05-search-agent.md` | 搜索Agent |
| `2026-06-04-06-admin.md` | 后台管理 + 用户配置 |
| `2026-06-04-07-youdao.md` | 有道Agent |
| `2026-06-04-08-integration.md` | 最终集成（本文件） |
