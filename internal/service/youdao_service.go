package service

import (
	"context"
	"fmt"
	"sync"

	"YoudaoNoteLm/internal/model/entity"
	"YoudaoNoteLm/internal/repository"
	"YoudaoNoteLm/internal/service/external"
	"YoudaoNoteLm/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type youdaoService struct {
	cli         external.YoudaoCLI
	bindingRepo repository.YoudaoBindingRepository
	sourceRepo  repository.SourceRepository
	embedding   EmbeddingService
	cancelFuncs sync.Map // taskID -> context.CancelFunc
}

// NewYoudaoService 创建有道云笔记服务
func NewYoudaoService(
	cli external.YoudaoCLI,
	bindingRepo repository.YoudaoBindingRepository,
	sourceRepo repository.SourceRepository,
	embedding EmbeddingService,
) YoudaoService {
	return &youdaoService{
		cli:         cli,
		bindingRepo: bindingRepo,
		sourceRepo:  sourceRepo,
		embedding:   embedding,
	}
}

// getAPIKey 获取用户的有道 API Key（内部辅助方法）
func (s *youdaoService) getAPIKey(userID uint) (string, error) {
	binding, err := s.bindingRepo.FindByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("查询绑定信息失败: %w", err)
	}
	if binding == nil || binding.Status != "active" {
		return "", fmt.Errorf("请先绑定有道云笔记账号")
	}
	return binding.APIKey, nil
}

// Bind 绑定有道 API Key
func (s *youdaoService) Bind(userID uint, apiKey string) error {
	// 1. 检查 CLI 是否可用
	if err := s.cli.CheckAvailable(); err != nil {
		return fmt.Errorf("youdaonote CLI 不可用: %w", err)
	}

	// 2. 验证 Key 有效性（调用 list 测试）
	_, err := s.cli.List(apiKey, "")
	if err != nil {
		return fmt.Errorf("API Key 无效，请检查后重试")
	}

	// 3. 保存或更新绑定
	existing, err := s.bindingRepo.FindByUserID(userID)
	if err != nil {
		return fmt.Errorf("查询绑定信息失败: %w", err)
	}

	if existing != nil {
		existing.APIKey = apiKey
		existing.Status = "active"
		return s.bindingRepo.Update(existing)
	}

	binding := &entity.YoudaoBinding{
		UserID: userID,
		APIKey: apiKey,
		Status: "active",
	}
	return s.bindingRepo.Create(binding)
}

// Unbind 解绑有道账号
func (s *youdaoService) Unbind(userID uint) error {
	return s.bindingRepo.Delete(userID)
}

// GetBinding 获取绑定信息
func (s *youdaoService) GetBinding(userID uint) (*entity.YoudaoBinding, error) {
	return s.bindingRepo.FindByUserID(userID)
}

// ListNotes 浏览有道云笔记目录
func (s *youdaoService) ListNotes(userID uint, folderID string) ([]external.YoudaoNoteItem, error) {
	apiKey, err := s.getAPIKey(userID)
	if err != nil {
		return nil, err
	}

	items, err := s.cli.List(apiKey, folderID)
	if err != nil {
		return nil, fmt.Errorf("获取笔记列表失败: %w", err)
	}

	return items, nil
}

// ImportNote 导入单篇有道云笔记到本系统
func (s *youdaoService) ImportNote(userID uint, notebookID uint, fileID string) (*entity.Source, error) {
	apiKey, err := s.getAPIKey(userID)
	if err != nil {
		return nil, err
	}

	// 1. 读取笔记内容
	readResult, err := s.cli.Read(apiKey, fileID)
	if err != nil {
		return nil, fmt.Errorf("读取笔记内容失败: %w", err)
	}

	// 2. 通过 list 获取笔记名称
	noteName := fileID // 降级使用 fileID
	items, listErr := s.cli.List(apiKey, "")
	if listErr == nil {
		for _, item := range items {
			if item.ID == fileID {
				noteName = item.Name
				break
			}
		}
	}

	// 3. 创建 Source 记录
	source := &entity.Source{
		UserID:          userID,
		NotebookID:      notebookID,
		Name:            noteName,
		Type:            "youdao",
		ExternalID:      fileID,
		MarkdownContent: readResult.Content,
		Status:          "ready",
	}

	if err := s.sourceRepo.Create(source); err != nil {
		return nil, fmt.Errorf("创建 Source 记录失败: %w", err)
	}

	// 4. 异步向量化
	if s.embedding != nil {
		go func() {
			if err := s.embedding.Vectorize(source.ID, readResult.Content); err != nil {
				logger.Warn("有道笔记向量化失败", zap.Uint("source_id", source.ID), zap.Error(err))
			} else {
				if err := s.sourceRepo.SetVectorized(source.ID); err != nil {
					logger.Warn("标记向量化状态失败", zap.Uint("source_id", source.ID), zap.Error(err))
				}
			}
		}()
	}

	logger.Info("有道笔记导入成功",
		zap.Uint("user_id", userID),
		zap.String("file_id", fileID),
		zap.String("name", noteName),
	)

	return source, nil
}

// ImportNotesBatch 批量导入有道云笔记
func (s *youdaoService) ImportNotesBatch(userID uint, notebookID uint, fileIDs []string) (string, []uint, error) {
	apiKey, err := s.getAPIKey(userID)
	if err != nil {
		return "", nil, err
	}

	// 去重
	seen := make(map[string]struct{}, len(fileIDs))
	uniqueIDs := make([]string, 0, len(fileIDs))
	for _, id := range fileIDs {
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}

	sourceIDs := make([]uint, 0, len(uniqueIDs))

	// 为每个 fileID 创建 pending 状态的 Source
	for _, fileID := range uniqueIDs {
		source := &entity.Source{
			UserID:     userID,
			NotebookID: notebookID,
			Name:       fileID,
			Type:       "youdao",
			ExternalID: fileID,
			Status:     "pending",
		}
		if err := s.sourceRepo.Create(source); err != nil {
			logger.Error("创建待导入有道笔记Source失败", zap.String("file_id", fileID), zap.Error(err))
			continue
		}
		sourceIDs = append(sourceIDs, source.ID)
	}

	if len(sourceIDs) == 0 {
		return "", nil, fmt.Errorf("创建导入记录失败")
	}

	// 创建可取消的 context
	taskID := uuid.New().String()
	taskCtx, cancel := context.WithCancel(context.Background())
	s.cancelFuncs.Store(taskID, cancel)

	// 异步处理
	go s.processBatch(taskCtx, taskID, apiKey, sourceIDs, uniqueIDs)

	return taskID, sourceIDs, nil
}

// processBatch 批量处理有道笔记导入
func (s *youdaoService) processBatch(taskCtx context.Context, taskID string, apiKey string, sourceIDs []uint, fileIDs []string) {
	defer s.cancelFuncs.Delete(taskID)

	concurrency := 3
	if len(fileIDs) < concurrency {
		concurrency = len(fileIDs)
	}

	type task struct {
		sourceID uint
		fileID   string
	}

	taskCh := make(chan task, concurrency)
	doneCh := make(chan struct{}, len(fileIDs))

	// 启动 worker
	for i := 0; i < concurrency; i++ {
		go func() {
			for t := range taskCh {
				if taskCtx.Err() != nil {
					doneCh <- struct{}{}
					continue
				}
				s.processSingleNote(taskCtx, apiKey, t.sourceID, t.fileID)
				doneCh <- struct{}{}
			}
		}()
	}

	// 分发任务
	go func() {
		for i, fileID := range fileIDs {
			if taskCtx.Err() != nil {
				break
			}
			taskCh <- task{sourceID: sourceIDs[i], fileID: fileID}
		}
		close(taskCh)
	}()

	// 等待完成
	for i := 0; i < len(fileIDs); i++ {
		<-doneCh
	}

	// 处理被取消的 pending 任务
	if taskCtx.Err() != nil {
		for _, sourceID := range sourceIDs {
			src, err := s.sourceRepo.FindByID(sourceID)
			if err != nil || src == nil {
				continue
			}
			if src.Status == "pending" {
				s.sourceRepo.UpdateStatus(sourceID, "cancelled", "任务已取消")
			}
		}
	}
}

// processSingleNote 处理单篇有道笔记导入
func (s *youdaoService) processSingleNote(taskCtx context.Context, apiKey string, sourceID uint, fileID string) {
	if taskCtx.Err() != nil {
		return
	}

	// 更新状态为 processing
	s.sourceRepo.UpdateStatus(sourceID, "processing", "")

	// 读取笔记内容
	readResult, err := s.cli.Read(apiKey, fileID)
	if err != nil {
		if taskCtx.Err() != nil {
			return
		}
		s.sourceRepo.UpdateStatus(sourceID, "failed", fmt.Sprintf("读取失败: %v", err))
		return
	}

	// 检查 Source 是否还存在
	existing, _ := s.sourceRepo.FindByID(sourceID)
	if existing == nil {
		return
	}

	// 更新内容和状态
	existing.MarkdownContent = readResult.Content
	existing.Status = "ready"
	if err := s.sourceRepo.Update(existing); err != nil {
		s.sourceRepo.UpdateStatus(sourceID, "failed", fmt.Sprintf("保存失败: %v", err))
		return
	}

	// 异步向量化
	if s.embedding != nil {
		go func() {
			if err := s.embedding.Vectorize(sourceID, readResult.Content); err != nil {
				logger.Warn("有道笔记批量导入向量化失败", zap.Uint("source_id", sourceID), zap.Error(err))
			} else {
				s.sourceRepo.SetVectorized(sourceID)
			}
		}()
	}
}
