package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3 配置结构体（从前端请求头获取）
type s3Config struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
	UseSSL          bool   `json:"useSsl"`
	SignatureType   string `json:"signatureType"`
}

// 从前端请求创建 S3 客户端
func createS3ClientFromRequest(cfg s3Config) (*minio.Client, error) {
	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.Endpoint == "" {
		return nil, fmt.Errorf("缺少必要的 S3 配置参数")
	}

	// 设置默认值
	if cfg.Region == "" {
		cfg.Region = "auto"
	}
	if cfg.SignatureType == "" {
		cfg.SignatureType = "V4"
	}

	// 设置签名类型
	var signatureTypeEnum credentials.SignatureType
	switch cfg.SignatureType {
	case "V2":
		signatureTypeEnum = credentials.SignatureV2
	case "V4":
		signatureTypeEnum = credentials.SignatureV4
	case "V4Streaming":
		signatureTypeEnum = credentials.SignatureV4Streaming
	case "Anonymous":
		signatureTypeEnum = credentials.SignatureAnonymous
	default:
		return nil, fmt.Errorf("invalid SIGNATURE_TYPE: %s", cfg.SignatureType)
	}

	// 创建客户端选项
	opts := &minio.Options{
		Creds:  credentials.NewStatic(cfg.AccessKeyID, cfg.SecretAccessKey, "", signatureTypeEnum),
		Secure: cfg.UseSSL,
	}

	// 设置区域
	if cfg.Region != "" {
		opts.Region = cfg.Region
	}

	// 创建客户端
	client, err := minio.New(cfg.Endpoint, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return client, nil
}

// 测试 S3 连接
func testS3Connection(client *minio.Client) error {
	ctx := context.Background()
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("failed to list buckets: %w", err)
	}
	fmt.Printf("✅ S3 连接成功！找到 %d 个存储桶\n", len(buckets))
	for _, bucket := range buckets {
		fmt.Printf("  - %s\n", bucket.Name)
	}
	return nil
}

func main() {
	// 从环境变量读取端口（可选，默认 9000）
	port := 9000
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	fmt.Println("🔍 S3 云存储管理服务启动中...")
	fmt.Println("   配置模式：仅支持前端动态配置（无 config.env 依赖）")

	// 启动 HTTP 服务器
	startServer(port)
}

// 启动 HTTP 服务器
func startServer(port int) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 添加中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(s3ClientMiddleware())

	// 设置路由
	setupRoutes(router)

	// 启动服务器
	portStr := fmt.Sprintf(":%d", port)
	fmt.Printf("🚀 HTTP 服务器已启动，端口: %s\n", portStr)
	fmt.Printf("📖 API: http://localhost%s/api/health\n", portStr)

	if err := router.Run(portStr); err != nil {
		log.Fatalf("❌ 启动服务器失败: %v", err)
	}
}

// CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-S3-Config")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// S3 客户端中间件 - 仅从请求头获取配置
func s3ClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		// GET 健康检查不需要 S3 配置
		if path == "/api/health" && method == "GET" {
			c.Next()
			return
		}

		// 从请求头获取 S3 配置
		configHeader := c.GetHeader("X-S3-Config")
		if configHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "未配置 S3 连接信息",
				"message": "请在前端设置中配置 S3 连接信息",
				"code": "CONFIG_REQUIRED",
			})
			c.Abort()
			return
		}

		var s3Cfg s3Config
		if err := json.Unmarshal([]byte(configHeader), &s3Cfg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "S3 配置格式无效",
				"message": "请求头中的 X-S3-Config 格式错误",
			})
			c.Abort()
			return
		}

		// 创建 S3 客户端
		client, err := createS3ClientFromRequest(s3Cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "S3 配置无效: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 存储到上下文
		c.Set("s3Client", client)
		c.Next()
	}
}
