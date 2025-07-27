package main

import (
	"Picture_bed/app"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// 设置路由
func setupRoutes(router *gin.Engine, s3Client *minio.Client) {
	handler := app.NewHandler(s3Client)
	
	// 健康检查
	router.GET("/api/health", handler.HealthCheck)

	// API v1 路由
	api := router.Group("/api/v1")
	{
		// 存储桶管理
		api.GET("/buckets", handler.ListBuckets)             // 列出所有存储桶
		api.POST("/buckets", handler.CreateBucket)           // 创建存储桶
		api.DELETE("/buckets/:bucket", handler.DeleteBucket) // 删除存储桶

		// 对象管理 - 使用不同的路径来避免冲突
		buckets := api.Group("/buckets/:bucket")
		{
			// 列出对象和上传
			buckets.GET("/objects", handler.ListObjects)                      // 列出对象（支持分页和文件夹浏览）
			buckets.POST("/objects", handler.UploadObject)                    // 上传对象
			buckets.POST("/folders", handler.CreateFolder)                    // 创建文件夹
			buckets.POST("/objects/batch-delete", handler.BatchDeleteObjects) // 批量删除对象

			// 对象操作 - 使用专门的路径前缀来避免通配符冲突
			objectsApi := buckets.Group("/api/objects")
			{
				objectsApi.GET("/:object/info", handler.GetObjectInfo) // 获取对象信息
				objectsApi.PUT("/:object", handler.UpdateObject)       // 更新对象
				objectsApi.DELETE("/:object", handler.DeleteObject)    // 删除对象
			}

			// 文件下载和文件夹浏览 - 放在最后，使用通配符
			buckets.GET("/browse/*filepath", handler.DownloadOrListObjects) // 浏览文件夹或下载文件
		}
	}

}
