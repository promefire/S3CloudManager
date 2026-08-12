package main

import (
	"Picture_bed/app"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// 设置路由
func setupRoutes(router *gin.Engine) {
	// 健康检查（GET - 不需要 S3 配置，返回服务状态）
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
			"config_required": true,
		})
	})

	// 测试连接（POST - 前端传入 S3 配置进行测试）
	router.POST("/api/health", func(c *gin.Context) {
		clientObj, exists := c.Get("s3Client")
		if !exists || clientObj == nil {
			c.JSON(400, gin.H{
				"error": "S3 配置无效",
			})
			return
		}
		s3Client, ok := clientObj.(*minio.Client)
		if !ok {
			c.JSON(400, gin.H{
				"error": "S3 客户端类型错误",
			})
			return
		}
		// 实际测试 S3 连接
		ctx := c.Request.Context()
		buckets, err := s3Client.ListBuckets(ctx)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "连接失败: " + err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "连接成功",
			"buckets": len(buckets),
		})
	})

	// API v1 路由
	api := router.Group("/api/v1")
	{
		// 存储桶管理
		api.GET("/buckets", func(c *gin.Context) {
			client, exists := c.Get("s3Client")
			if !exists || client == nil {
				c.JSON(400, gin.H{"error": "未配置 S3 连接信息"})
				return
			}
			handler := app.NewHandler(client.(*minio.Client))
			handler.ListBuckets(c)
		})

		api.POST("/buckets", func(c *gin.Context) {
			client, exists := c.Get("s3Client")
			if !exists || client == nil {
				c.JSON(400, gin.H{"error": "未配置 S3 连接信息"})
				return
			}
			handler := app.NewHandler(client.(*minio.Client))
			handler.CreateBucket(c)
		})

		api.DELETE("/buckets/:bucket", func(c *gin.Context) {
			client, exists := c.Get("s3Client")
			if !exists || client == nil {
				c.JSON(400, gin.H{"error": "未配置 S3 连接信息"})
				return
			}
			handler := app.NewHandler(client.(*minio.Client))
			handler.DeleteBucket(c)
		})

		// 对象管理
		buckets := api.Group("/buckets/:bucket")
		{
			buckets.GET("/objects", func(c *gin.Context) {
				client, exists := c.Get("s3Client")
				if !exists || client == nil {
					c.JSON(400, gin.H{"error": "未配置 S3 连接信息"})
					return
				}
				handler := app.NewHandler(client.(*minio.Client))
				handler.ListObjects(c)
			})

			buckets.POST("/objects", func(c *gin.Context) {
				client, exists := c.Get("s3Client")
				if !exists || client == nil {
					c.JSON(400, gin.H{"error": "未配置 S3 连接信息"})
					return
				}
				handler := app.NewHandler(client.(*minio.Client))
				handler.UploadObject(c)
			})

			buckets.POST("/folders", func(c *gin.Context) {
				client, exists := c.Get("s3Client")
				if !exists || client == nil {
					c.JSON(400, gin.H{"error": "未配置 S3 连接信息"})
					return
				}
				handler := app.NewHandler(client.(*minio.Client))
				handler.CreateFolder(c)
			})

			buckets.POST("/objects/batch-delete", func(c *gin.Context) {
				client, exists := c.Get("s3Client")
				if !exists || client == nil {
					c.JSON(400, gin.H{"error": "未配置 S3 连接信息"})
					return
				}
				handler := app.NewHandler(client.(*minio.Client))
				handler.BatchDeleteObjects(c)
			})

			objectsApi := buckets.Group("/api/objects")
			{
				objectsApi.GET("/:object/info", func(c *gin.Context) {
					client, exists := c.Get("s3Client")
					if !exists || client == nil {
						c.JSON(400, gin.H{"error": "未配置 S3 连接信息"})
						return
					}
					handler := app.NewHandler(client.(*minio.Client))
					handler.GetObjectInfo(c)
				})

				objectsApi.PUT("/:object", func(c *gin.Context) {
					client, exists := c.Get("s3Client")
					if !exists || client == nil {
						c.JSON(400, gin.H{"error": "未配置 S3 连接信息"})
						return
					}
					handler := app.NewHandler(client.(*minio.Client))
					handler.UpdateObject(c)
				})

				objectsApi.DELETE("/:object", func(c *gin.Context) {
					client, exists := c.Get("s3Client")
					if !exists || client == nil {
						c.JSON(400, gin.H{"error": "未配置 S3 连接信息"})
						return
					}
					handler := app.NewHandler(client.(*minio.Client))
					handler.DeleteObject(c)
				})
			}

			buckets.GET("/browse/*filepath", func(c *gin.Context) {
				client, exists := c.Get("s3Client")
				if !exists || client == nil {
					c.JSON(400, gin.H{"error": "未配置 S3 连接信息"})
					return
				}
				handler := app.NewHandler(client.(*minio.Client))
				handler.DownloadOrListObjects(c)
			})
		}
	}
}
