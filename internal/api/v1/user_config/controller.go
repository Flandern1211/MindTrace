package user_config

import (
	"YoudaoNoteLm/internal/middleware"
	"YoudaoNoteLm/internal/model/dto/request"
	"YoudaoNoteLm/internal/model/entity"
	"YoudaoNoteLm/internal/service"
	"YoudaoNoteLm/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	configService service.UserConfigService
}

func NewController(configService service.UserConfigService) *Controller {
	return &Controller{configService: configService}
}

// ===== Search Config =====

func (ctrl *Controller) ListSearchConfigs(c *gin.Context) {
	userID := middleware.GetUserID(c)
	configs, err := ctrl.configService.ListSearchConfigs(userID)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, configs)
}

func (ctrl *Controller) CreateSearchConfig(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req request.UserConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	config := &entity.UserConfig{
		Name: req.Name, Provider: req.Provider, APIKey: req.APIKey,
		APIURL: req.APIURL, DailyQuota: req.DailyQuota, Enabled: true,
	}

	if err := ctrl.configService.CreateSearchConfig(userID, config); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, config)
}

func (ctrl *Controller) UpdateSearchConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req request.UserConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	config := &entity.UserConfig{
		Name: req.Name, Provider: req.Provider, APIKey: req.APIKey,
		APIURL: req.APIURL, DailyQuota: req.DailyQuota,
	}

	if err := ctrl.configService.UpdateSearchConfig(uint(id), config); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, config)
}

func (ctrl *Controller) DeleteSearchConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := ctrl.configService.DeleteSearchConfig(uint(id)); err != nil {
		response.BizError(c, err)
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ===== ASR Config =====

func (ctrl *Controller) ListASRConfigs(c *gin.Context) {
	userID := middleware.GetUserID(c)
	configs, err := ctrl.configService.ListASRConfigs(userID)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, configs)
}

func (ctrl *Controller) CreateASRConfig(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req request.UserConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	config := &entity.UserConfig{
		Name: req.Name, Provider: req.Provider, APIKey: req.APIKey,
		APIURL: req.APIURL, ExtraConfig: string(req.ExtraConfig), Enabled: true,
	}

	if err := ctrl.configService.CreateASRConfig(userID, config); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, config)
}

func (ctrl *Controller) UpdateASRConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req request.UserConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	config := &entity.UserConfig{
		Name: req.Name, Provider: req.Provider, APIKey: req.APIKey,
		APIURL: req.APIURL, ExtraConfig: string(req.ExtraConfig),
	}

	if err := ctrl.configService.UpdateASRConfig(uint(id), config); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, config)
}

func (ctrl *Controller) DeleteASRConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := ctrl.configService.DeleteASRConfig(uint(id)); err != nil {
		response.BizError(c, err)
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ===== Embedding Config =====

func (ctrl *Controller) ListEmbeddingConfigs(c *gin.Context) {
	userID := middleware.GetUserID(c)
	configs, err := ctrl.configService.ListEmbeddingConfigs(userID)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, configs)
}

func (ctrl *Controller) CreateEmbeddingConfig(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req request.UserConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	config := &entity.UserConfig{
		Name: req.Name, Provider: req.Provider, APIKey: req.APIKey,
		APIURL: req.APIURL, Model: req.Model, Dimensions: req.Dimensions, Enabled: true,
	}

	if err := ctrl.configService.CreateEmbeddingConfig(userID, config); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, config)
}

func (ctrl *Controller) UpdateEmbeddingConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req request.UserConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	config := &entity.UserConfig{
		Name: req.Name, Provider: req.Provider, APIKey: req.APIKey,
		APIURL: req.APIURL, Model: req.Model, Dimensions: req.Dimensions,
	}

	if err := ctrl.configService.UpdateEmbeddingConfig(uint(id), config); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, config)
}

func (ctrl *Controller) DeleteEmbeddingConfig(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := ctrl.configService.DeleteEmbeddingConfig(uint(id)); err != nil {
		response.BizError(c, err)
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}