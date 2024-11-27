package main

import (
	"log"
	"net/http"
	"time"

	"grading-api/config"
	"grading-api/types"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

func init() {
	// 加载配置文件
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatal("Failed to load config:", err)
	}
}

func main() {
	r := gin.Default()
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		api.POST("/grade", handleGrading(nil))
		api.GET("/grade/history", getGradingHistory)
		api.DELETE("/grade/history", clearGradingHistory)
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// createOpenAIClient 创建OpenAI客户端
func createOpenAIClient(customConfig types.OpenAIConfig) *openai.Client {
	// 获取配置（合并默认配置和自定义配置）
	cfg := config.GetOpenAIConfig(customConfig.APIKey, customConfig.BaseURL)

	clientConfig := openai.DefaultConfig(cfg.DefaultAPIKey)

	if cfg.DefaultBaseURL != "" {
		clientConfig.BaseURL = cfg.DefaultBaseURL
	}

	clientConfig.HTTPClient = &http.Client{
		Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
	}

	return openai.NewClientWithConfig(clientConfig)
}
