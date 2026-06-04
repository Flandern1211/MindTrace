# 配置服务实施计划

> **智能体执行指南：** 必须使用子技能：superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务执行本计划。

**目标：** 实现统一的服务配置路由和降级策略，为搜索、ASR、Embedding 等服务提供配置查找能力。

**架构：** ConfigService 按优先级查找服务配置：用户自定义配置 → sys_config 内置配置 → 兜底默认实现（如 DuckDuckGo）。

**技术栈：** GORM

**前置依赖：** `2026-06-04-01-foundation.md`（仓储层）

---

## 文件结构

### 新增文件

```
internal/service/external/search_engine_interface.go → SearchEngine 接口
internal/service/external/duckduckgo_engine.go      → DuckDuckGo 搜索引擎实现
internal/service/external/custom_engine.go          → 自定义搜索引擎实现
internal/service/config_service.go                  → ConfigService 实现
```

---

## 任务 1：SearchEngine 接口 + 实现

**文件：**
- 新建：`internal/service/external/search_engine_interface.go`
- 新建：`internal/service/external/duckduckgo_engine.go`
- 新建：`internal/service/external/custom_engine.go`

- [ ] **步骤 1：创建 SearchEngine 接口**

```go
// internal/service/external/search_engine_interface.go
package external

// SearchResultItem 搜索结果项
type SearchResultItem struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

// SearchEngine 搜索引擎接口
type SearchEngine interface {
	Search(query string, limit int) ([]SearchResultItem, error)
	Name() string
}
```

- [ ] **步骤 2：创建 DuckDuckGoEngine**

```go
// internal/service/external/duckduckgo_engine.go
package external

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type duckDuckGoEngine struct {
	httpClient *http.Client
}

func NewDuckDuckGoEngine() SearchEngine {
	return &duckDuckGoEngine{
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (e *duckDuckGoEngine) Name() string {
	return "duckduckgo"
}

func (e *duckDuckGoEngine) Search(query string, limit int) ([]SearchResultItem, error) {
	if limit <= 0 {
		limit = 10
	}

	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("DuckDuckGo请求失败: %w", err)
	}
	defer resp.Body.Close()

	body := make([]byte, 1024*1024)
	n, _ := resp.Body.Read(body)
	html := string(body[:n])

	results := parseDuckDuckGoResults(html, limit)
	return results, nil
}

func parseDuckDuckGoResults(html string, limit int) []SearchResultItem {
	var results []SearchResultItem

	titleRe := regexp.MustCompile(`<a[^>]*class="result__a"[^>]*href="([^"]*)"[^>]*>([^<]*)</a>`)
	snippetRe := regexp.MustCompile(`<a[^>]*class="result__snippet"[^>]*>([^<]*(?:<[^>]*>[^<]*)*)</a>`)

	titles := titleRe.FindAllStringSubmatch(html, -1)
	snippets := snippetRe.FindAllStringSubmatch(html, -1)

	for i, match := range titles {
		if len(results) >= limit {
			break
		}
		if len(match) < 3 {
			continue
		}

		item := SearchResultItem{
			URL:   strings.TrimSpace(match[1]),
			Title: strings.TrimSpace(stripHTML(match[2])),
		}
		if i < len(snippets) && len(snippets[i]) > 1 {
			item.Snippet = strings.TrimSpace(stripHTML(snippets[i][1]))
		}

		if item.URL != "" && item.Title != "" {
			results = append(results, item)
		}
	}

	return results
}

func stripHTML(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(s, "")
}
```

- [ ] **步骤 3：创建 CustomEngine**

```go
// internal/service/external/custom_engine.go
package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type customEngine struct {
	name   string
	apiURL string
	apiKey string
	client *http.Client
}

func NewCustomEngine(name, apiURL, apiKey string) SearchEngine {
	return &customEngine{
		name: name, apiURL: apiURL, apiKey: apiKey,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (e *customEngine) Name() string {
	return e.name
}

func (e *customEngine) Search(query string, limit int) ([]SearchResultItem, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"query": query,
		"limit": limit,
	})

	req, err := http.NewRequest("POST", e.apiURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if e.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.apiKey)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("自定义搜索API请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("搜索API返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Results []SearchResultItem `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("解析搜索结果失败: %w", err)
	}

	return apiResp.Results, nil
}
```

- [ ] **步骤 4：验证编译**

```bash
go build ./internal/service/external/...
```

- [ ] **步骤 5：提交**

```bash
git add internal/service/external/search_engine_interface.go internal/service/external/duckduckgo_engine.go internal/service/external/custom_engine.go
git commit -m "feat(external): 添加 SearchEngine 接口和 DuckDuckGo/Custom 实现"
```

---

## 任务 2：ConfigService

**文件：**
- 新建：`internal/service/config_service.go`

- [ ] **步骤 1：创建 ConfigService**

```go
// internal/service/config_service.go
package service

import (
	"YoudaoNoteLm/internal/repository"
	"YoudaoNoteLm/internal/service/external"
	bizerrors "YoudaoNoteLm/pkg/errors"
	"YoudaoNoteLm/pkg/logger"

	"go.uber.org/zap"
)

// ConfigService 配置路由服务接口
type ConfigService interface {
	GetSearchEngine(userID uint) (external.SearchEngine, error)
	GetASRService(userID uint) (external.ASRService, error)
}

type configService struct {
	sysConfigRepo       repository.SysConfigRepository
	userSearchConfigRepo repository.UserSearchConfigRepository
	userASRConfigRepo   repository.UserASRConfigRepository
}

func NewConfigService(
	sysConfigRepo repository.SysConfigRepository,
	userSearchConfigRepo repository.UserSearchConfigRepository,
	userASRConfigRepo repository.UserASRConfigRepository,
) ConfigService {
	return &configService{
		sysConfigRepo:       sysConfigRepo,
		userSearchConfigRepo: userSearchConfigRepo,
		userASRConfigRepo:   userASRConfigRepo,
	}
}

func (s *configService) GetSearchEngine(userID uint) (external.SearchEngine, error) {
	// 1. 查用户搜索配置
	userConfigs, err := s.userSearchConfigRepo.FindByUser(userID)
	if err == nil {
		for _, cfg := range userConfigs {
			if !cfg.Enabled {
				continue
			}
			engine := external.NewCustomEngine(cfg.Provider, cfg.APIURL, cfg.APIKey)
			return engine, nil
		}
	}

	// 2. 降级到系统内置配置
	builtins, err := s.sysConfigRepo.FindByGroup("search")
	if err == nil {
		for _, builtin := range builtins {
			if builtin.Enabled {
				logger.Info("使用系统内置搜索配置", zap.String("key", builtin.ConfigKey))
				break
			}
		}
	}

	// 3. DuckDuckGo 兜底
	logger.Info("使用 DuckDuckGo 兜底搜索引擎")
	return external.NewDuckDuckGoEngine(), nil
}

func (s *configService) GetASRService(userID uint) (external.ASRService, error) {
	// 1. 查用户ASR配置
	userConfigs, err := s.userASRConfigRepo.FindByUser(userID)
	if err == nil {
		for _, cfg := range userConfigs {
			if !cfg.Enabled {
				continue
			}
			return external.NewASRService(cfg.APIURL, cfg.APIKey), nil
		}
	}

	// 2. 降级到系统内置配置
	builtins, err := s.sysConfigRepo.FindByGroup("asr")
	if err == nil {
		for _, builtin := range builtins {
			if builtin.Enabled {
				logger.Info("使用系统内置ASR配置", zap.String("key", builtin.ConfigKey))
				break
			}
		}
	}

	// 3. 无可用配置
	return nil, bizerrors.New(bizerrors.CodeASTranscriptionFailed, "未配置ASR服务")
}
```

- [ ] **步骤 2：验证编译**

```bash
go build ./internal/service/...
```

- [ ] **步骤 3：提交**

```bash
git add internal/service/config_service.go
git commit -m "feat(service): 添加 ConfigService 服务路由降级"
```

---

## 验证

```bash
go build ./internal/service/...
go vet ./internal/service/...
```

**降级逻辑：**
1. 用户自定义配置（enabled=true，按 priority 排序）
2. sys_config 内置配置
3. 兜底默认实现（DuckDuckGo 搜索 / 返回错误）
