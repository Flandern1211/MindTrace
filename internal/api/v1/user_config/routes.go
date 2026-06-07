package user_config

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册用户配置路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	cfg := r.Group("/user/config")
	{
		cfg.GET("/search", ctrl.ListSearchConfigs)
		cfg.POST("/search", ctrl.CreateSearchConfig)
		cfg.PUT("/search/:id", ctrl.UpdateSearchConfig)
		cfg.DELETE("/search/:id", ctrl.DeleteSearchConfig)

		cfg.GET("/asr", ctrl.ListASRConfigs)
		cfg.POST("/asr", ctrl.CreateASRConfig)
		cfg.PUT("/asr/:id", ctrl.UpdateASRConfig)
		cfg.DELETE("/asr/:id", ctrl.DeleteASRConfig)

		cfg.GET("/embedding", ctrl.ListEmbeddingConfigs)
		cfg.POST("/embedding", ctrl.CreateEmbeddingConfig)
		cfg.PUT("/embedding/:id", ctrl.UpdateEmbeddingConfig)
		cfg.DELETE("/embedding/:id", ctrl.DeleteEmbeddingConfig)
	}
}