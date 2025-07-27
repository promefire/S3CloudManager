package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

// 配置结构体
type config struct {
	Endpoint            string
	AccessKeyID         string
	SecretAccessKey     string
	Region              string
	UseSSL              bool
	SkipSSLVerification bool
	SignatureType       string
	Port                int
}

// 读取配置
func loadConfig() config {
	// 设置配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("⚠️  警告: 无法读取配置文件: %v", err)
	}

	// 也支持环境变量
	//viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("ENDPOINT", "s3.amazonaws.com")
	viper.SetDefault("USE_SSL", true)
	viper.SetDefault("SKIP_SSL_VERIFICATION", false)
	viper.SetDefault("SIGNATURE_TYPE", "V4")
	viper.SetDefault("PORT", 8080)

	return config{
		Endpoint:            viper.GetString("ENDPOINT"),
		AccessKeyID:         viper.GetString("ACCESS_KEY_ID"),
		SecretAccessKey:     viper.GetString("SECRET_ACCESS_KEY"),
		Region:              viper.GetString("REGION"),
		UseSSL:              viper.GetBool("USE_SSL"),
		SkipSSLVerification: viper.GetBool("SKIP_SSL_VERIFICATION"),
		SignatureType:       viper.GetString("SIGNATURE_TYPE"),
		Port:                viper.GetInt("PORT"),
	}
}

// 创建 S3 客户端
func createS3Client(cfg config) (*minio.Client, error) {
	// 验证必需参数
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("ACCESS_KEY_ID is required")
	}
	if cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("SECRET_ACCESS_KEY is required")
	}

	// 设置签名类型
	var signatureType credentials.SignatureType
	switch cfg.SignatureType {
	case "V2":
		signatureType = credentials.SignatureV2
	case "V4":
		signatureType = credentials.SignatureV4
	case "V4Streaming":
		signatureType = credentials.SignatureV4Streaming
	case "Anonymous":
		signatureType = credentials.SignatureAnonymous
	default:
		return nil, fmt.Errorf("invalid SIGNATURE_TYPE: %s", cfg.SignatureType)
	}

	// 创建客户端选项
	opts := &minio.Options{
		Creds:  credentials.NewStatic(cfg.AccessKeyID, cfg.SecretAccessKey, "", signatureType),
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

	// 测试列出存储桶（最基本的操作）
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("failed to list buckets: %w", err)
	}

	fmt.Printf("✅ S3 连接成功！\n")
	fmt.Printf("📊 找到 %d 个存储桶:\n", len(buckets))

	for _, bucket := range buckets {
		fmt.Printf("  - %s (创建时间: %s)\n", bucket.Name, bucket.CreationDate.Format("2006-01-02 15:04:05"))
	}

	return nil
}

func main() {
	fmt.Println("🔍 开始测试 S3 连接...")

	// 加载配置
	cfg := loadConfig()

	fmt.Printf("📋 配置信息:\n")
	fmt.Printf("  端点: %s\n", cfg.Endpoint)
	fmt.Printf("  区域: %s\n", cfg.Region)
	fmt.Printf("  使用 SSL: %t\n", cfg.UseSSL)
	fmt.Printf("  签名类型: %s\n", cfg.SignatureType)

	// 创建 S3 客户端
	client, err := createS3Client(cfg)
	if err != nil {
		log.Fatalf("❌ 创建 S3 客户端失败: %v", err)
	}

	// 测试连接
	if err := testS3Connection(client); err != nil {
		log.Fatalf("❌ S3 连接测试失败: %v", err)
	}

	fmt.Println("🎉 S3 连接测试通过！")

	// 启动 HTTP 服务器
	startServer(client, cfg)

}

// 启动 HTTP 服务器
func startServer(client *minio.Client, cfg config) {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 路由器
	router := gin.Default()

	// 添加中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// 设置路由
	setupRoutes(router, client)

	// 启动服务器
	port := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("🚀 启动 HTTP 服务器，端口: %s\n", port)
	fmt.Printf("📖 API 文档: http://localhost%s/api/v1/health\n", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("❌ 启动服务器失败: %v", err)
	}
}

// CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
