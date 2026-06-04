# 搜索Agent实施计划

> **智能体执行指南：** 必须使用子技能：superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务执行本计划。

**目标：** 实现搜索Agent服务，支持关键词搜索和URL直接导入两种模式。

**架构：** SearchAgentService 根据搜索类型路由：关键词搜索通过 ConfigService 获取搜索引擎执行搜索；URL 导入通过 ImporterService 直接抓取网页内容。

**技术栈：** Gin、Eino（Agent框架，桩代码）

**前置依赖：**
- `2026-06-04-01-foundation.md`（实体和仓储）
- `2026-06-04-03-importer.md`（ImporterService）
- `2026-06-04-04-config-service.md`（ConfigService）

---

## 文件结构

### 新增文件

```
internal/service/search_agent_interface.go          → SearchAgentService 接口
internal/service/search_agent_service.go            → SearchAgentService 实现
internal/model/dto/request/search.go                → 搜索请求 DTO
internal/model/dto/response/search.go               → 搜索响应 DTO
internal/api/v1/search/controller.go                → 搜索控制器
internal/api/v1/search/routes.go                    → 搜索路由
internal/agent/search/agent.go                      → Eino 搜索 Agent 桩代码
internal/agent/search/tools.go                      → Agent 工具定义
internal/agent/search/prompts.go                    → Agent 提示词
```

---

## 任务 1：搜索 DTO

**文件：**
- 新建：`internal/model/dto/request/search.go`
- 新建：`internal/model/dto/response/search.go`

- [ ] **步骤 1：创建搜索请求 DTO**

```go
// internal/model/dto/request/search.go
package request

// SearchRequest 搜索请求
type SearchRequest struct {
	Query string `json:"query" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=keyword url"`
}

// SearchImportRequest 搜索结果批量导入请求
type SearchImportRequest struct {
	URLs []string `json:"urls" binding:"required,min=1"`
}
```

- [ ] **步骤 2：创建搜索响应 DTO**

```go
// internal/model/dto/response/search.go
package response

// SearchResultResponse 搜索结果响应
type SearchResultResponse struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}
```

- [ ] **步骤 3：验证编译**

```bash
go build ./internal/model/dto/...
```

- [ ] **步骤 4：提交**

```bash
git add internal/model/dto/request/search.go internal/model/dto/response/search.go
git commit -m "feat(dto): 添加搜索请求/响应 DTO"
```

---

## 任务 2：SearchAgentService

**文件：**
- 新建：`internal/service/search_agent_interface.go`
- 新建：`internal/service/search_agent_service.go`

- [ ] **步骤 1：创建 SearchAgentService 接口**

```go
// internal/service/search_agent_interface.go
package service

import (
	"YoudaoNoteLm/internal/model/dto/response"
	"YoudaoNoteLm/internal/model/entity"
)

// SearchAgentService 搜索Agent服务接口
type SearchAgentService interface {
	Search(userID, notebookID uint, query string, searchType string) (interface{}, error)
	ImportFromURL(userID, notebookID uint, url string) (*entity.Source, error)
	ImportSearchResults(userID, notebookID uint, urls []string) (*entity.ImportTask, error)
}

// SearchResponse 搜索关键词响应
type SearchResponse struct {
	Results []response.SearchResultResponse `json:"results"`
}
```

- [ ] **步骤 2：创建 SearchAgentService 实现**

```go
// internal/service/search_agent_service.go
package service

import (
	"YoudaoNoteLm/internal/model/dto/response"
	"YoudaoNoteLm/internal/model/entity"
	bizerrors "YoudaoNoteLm/pkg/errors"
	"YoudaoNoteLm/pkg/logger"

	"go.uber.org/zap"
)

type searchAgentService struct {
	configService ConfigService
	importer      ImporterService
}

func NewSearchAgentService(configService ConfigService, importer ImporterService) SearchAgentService {
	return &searchAgentService{
		configService: configService,
		importer:      importer,
	}
}

func (s *searchAgentService) Search(userID, notebookID uint, query string, searchType string) (interface{}, error) {
	if searchType == "url" {
		return s.ImportFromURL(userID, notebookID, query)
	}

	engine, err := s.configService.GetSearchEngine(userID)
	if err != nil {
		return nil, err
	}

	results, err := engine.Search(query, 10)
	if err != nil {
		return nil, bizerrors.NewWithErr(bizerrors.CodeSearchQuotaExhausted, "搜索失败", err)
	}

	searchResults := make([]response.SearchResultResponse, 0, len(results))
	for _, r := range results {
		searchResults = append(searchResults, response.SearchResultResponse{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Snippet,
		})
	}

	logger.Info("搜索完成",
		zap.Uint("user_id", userID),
		zap.String("query", query),
		zap.Int("results", len(searchResults)),
	)

	return &SearchResponse{Results: searchResults}, nil
}

func (s *searchAgentService) ImportFromURL(userID, notebookID uint, url string) (*entity.Source, error) {
	task, err := s.importer.ImportSearchResults(userID, notebookID, []string{url})
	if err != nil {
		return nil, err
	}
	_ = task
	// TODO: 等待任务完成后返回 Source
	return nil, nil
}

func (s *searchAgentService) ImportSearchResults(userID, notebookID uint, urls []string) (*entity.ImportTask, error) {
	return s.importer.ImportSearchResults(userID, notebookID, urls)
}
```

- [ ] **步骤 3：验证编译**

```bash
go build ./internal/service/...
```

- [ ] **步骤 4：提交**

```bash
git add internal/service/search_agent_interface.go internal/service/search_agent_service.go
git commit -m "feat(service): 添加 SearchAgentService"
```

---

## 任务 3：搜索控制器 + 路由

**文件：**
- 新建：`internal/api/v1/search/controller.go`
- 新建：`internal/api/v1/search/routes.go`

- [ ] **步骤 1：创建搜索控制器**

```go
// internal/api/v1/search/controller.go
package search

import (
	"YoudaoNoteLm/internal/middleware"
	"YoudaoNoteLm/internal/model/dto/request"
	"YoudaoNoteLm/internal/service"
	"YoudaoNoteLm/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	searchService service.SearchAgentService
}

func NewController(searchService service.SearchAgentService) *Controller {
	return &Controller{searchService: searchService}
}

func (ctrl *Controller) Search(c *gin.Context) {
	userID := middleware.GetUserID(c)
	nbID, err := strconv.ParseUint(c.Param("nbId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的笔记本ID")
		return
	}

	var req request.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := ctrl.searchService.Search(userID, uint(nbID), req.Query, req.Type)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

func (ctrl *Controller) ImportResults(c *gin.Context) {
	userID := middleware.GetUserID(c)
	nbID, err := strconv.ParseUint(c.Param("nbId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的笔记本ID")
		return
	}

	var req request.SearchImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	task, err := ctrl.searchService.ImportSearchResults(userID, uint(nbID), req.URLs)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, task)
}
```

- [ ] **步骤 2：创建路由**

```go
// internal/api/v1/search/routes.go
package search

import "github.com/gin-gonic/gin"

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	notebooks := r.Group("/notebooks/:nbId/search")
	{
		notebooks.POST("", ctrl.Search)
		notebooks.POST("/import", ctrl.ImportResults)
	}
}
```

- [ ] **步骤 3：验证编译**

```bash
go build ./internal/api/v1/search/...
```

- [ ] **步骤 4：提交**

```bash
mkdir -p internal/api/v1/search
git add internal/api/v1/search/
git commit -m "feat(api): 添加 Search 控制器和路由"
```

---

## 任务 4：搜索 Agent 桩代码

**文件：**
- 新建：`internal/agent/search/agent.go`
- 新建：`internal/agent/search/tools.go`
- 新建：`internal/agent/search/prompts.go`

- [ ] **步骤 1：创建搜索 Agent 桩代码**

```go
// internal/agent/search/prompts.go
package search

const SearchSystemPrompt = `你是一个网络搜索助手。你的职责是：
1. 理解用户的搜索意图
2. 使用搜索引擎查找相关信息
3. 如果用户提供了URL，直接导入该URL的内容
4. 返回结构化的搜索结果`
```

```go
// internal/agent/search/tools.go
package search

type SearchTool struct {
	Name        string
	Description string
	Parameters  map[string]any
}

func NewSearchTool() *SearchTool {
	return &SearchTool{
		Name:        "web_search",
		Description: "搜索网络内容，输入搜索关键词，返回搜索结果列表",
		Parameters: map[string]any{
			"query": map[string]any{
				"type":        "string",
				"description": "搜索关键词",
			},
		},
	}
}

type URLImportTool struct {
	Name        string
	Description string
	Parameters  map[string]any
}

func NewURLImportTool() *URLImportTool {
	return &URLImportTool{
		Name:        "import_url",
		Description: "导入指定URL的网页内容，输入URL地址，返回导入结果",
		Parameters: map[string]any{
			"url": map[string]any{
				"type":        "string",
				"description": "要导入的URL地址",
			},
		},
	}
}
```

```go
// internal/agent/search/agent.go
package search

// SearchAgent 搜索Agent（Eino Agent 封装）
// TODO: 接入 Eino Agent SDK 后替换为实际实现
type SearchAgent struct {
	prompt string
	tools  []interface{}
}

func NewSearchAgent() *SearchAgent {
	return &SearchAgent{
		prompt: SearchSystemPrompt,
		tools: []interface{}{
			NewSearchTool(),
			NewURLImportTool(),
		},
	}
}
```

- [ ] **步骤 2：验证编译**

```bash
go build ./internal/agent/search/...
```

- [ ] **步骤 3：提交**

```bash
mkdir -p internal/agent/search
git add internal/agent/search/
git commit -m "feat(agent): 添加搜索 Agent 桩代码"
```

---

## 验证

```bash
go build ./internal/service/... ./internal/api/v1/search/... ./internal/agent/search/...
go vet ./internal/service/... ./internal/api/v1/search/... ./internal/agent/search/...
```

**API 端点：**
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/notebooks/:nbId/search | 搜索（关键词/URL） |
| POST | /api/v1/notebooks/:nbId/search/import | 批量导入搜索结果 |
