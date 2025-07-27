package app

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// Handler 路由处理器结构体
type Handler struct {
	s3Client *minio.Client
}

// NewHandler 创建新的处理器
func NewHandler(client *minio.Client) *Handler {
	return &Handler{
		s3Client: client,
	}
}

// HealthCheck 健康检查
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "图片床服务运行正常",
	})
}

// ListBuckets 列出所有存储桶
func (h *Handler) ListBuckets(c *gin.Context) {
	ctx := c.Request.Context()
	buckets, err := h.s3Client.ListBuckets(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取存储桶列表失败: " + err.Error(),
		})
		return
	}

	var bucketList []gin.H
	for _, bucket := range buckets {
		bucketList = append(bucketList, gin.H{
			"name":         bucket.Name,
			"creationDate": bucket.CreationDate.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"buckets": bucketList,
		"count":   len(bucketList),
	})
}

// CreateBucket 创建存储桶
func (h *Handler) CreateBucket(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Region string `json:"region"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证存储桶名称格式
	if err := validateBucketName(req.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "存储桶名称格式不正确: " + err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// 检查存储桶是否已存在
	exists, err := h.s3Client.BucketExists(ctx, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "检查存储桶存在性失败: " + err.Error(),
		})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "存储桶已存在",
		})
		return
	}

	err = h.s3Client.MakeBucket(ctx, req.Name, minio.MakeBucketOptions{
		Region: req.Region,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建存储桶失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "存储桶创建成功",
		"name":    req.Name,
	})
}

// DeleteBucket 删除存储桶
func (h *Handler) DeleteBucket(c *gin.Context) {
	bucketName := c.Param("bucket")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "存储桶名称不能为空",
		})
		return
	}

	ctx := c.Request.Context()

	// 检查存储桶是否存在
	exists, err := h.s3Client.BucketExists(ctx, bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "检查存储桶存在性失败: " + err.Error(),
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "存储桶不存在",
		})
		return
	}

	err = h.s3Client.RemoveBucket(ctx, bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除存储桶失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "存储桶删除成功",
		"name":    bucketName,
	})
}

// ListObjects 列出存储桶中的对象（支持分页和前缀过滤，支持文件夹浏览）
func (h *Handler) ListObjects(c *gin.Context) {
	bucketName := c.Param("bucket")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "存储桶名称不能为空",
		})
		return
	}

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	prefix := c.Query("prefix")       // 支持前缀过滤
	delimiter := c.Query("delimiter") // 支持分隔符，用于文件夹浏览

	// 如果没有指定分隔符但请求文件夹浏览，默认使用 "/"
	if delimiter == "" && c.Query("folder") == "true" {
		delimiter = "/"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	ctx := c.Request.Context()

	// 检查存储桶是否存在
	exists, err := h.s3Client.BucketExists(ctx, bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "检查存储桶存在性失败: " + err.Error(),
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "存储桶不存在",
		})
		return
	}

	// 使用分隔符进行文件夹式浏览
	listOptions := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: delimiter == "", // 如果没有分隔符，则递归列出所有对象
	}

	objectsCh := h.s3Client.ListObjects(ctx, bucketName, listOptions)

	// 🔄 修改：定义文件信息结构体，包含时间字段用于排序
	type FileItem struct {
		Data         gin.H
		LastModified time.Time
		IsFolder     bool
	}

	var allFileItems []FileItem
	var folders []string
	folderSet := make(map[string]bool)

	for object := range objectsCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "列出对象失败: " + object.Err.Error(),
			})
			return
		}

		// 如果使用分隔符，需要处理文件夹
		if delimiter != "" {
			// 移除前缀，获取相对路径
			relativePath := strings.TrimPrefix(object.Key, prefix)

			// 如果包含分隔符，说明这是一个子文件夹中的文件
			if delimiterIndex := strings.Index(relativePath, delimiter); delimiterIndex != -1 {
				// 提取文件夹名
				folderName := relativePath[:delimiterIndex+1]
				fullFolderPath := prefix + folderName

				if !folderSet[fullFolderPath] {
					folderSet[fullFolderPath] = true
					folders = append(folders, fullFolderPath)
				}
				continue // 跳过子文件夹中的文件
			}
		}

		// 🔄 修改：添加文件对象到FileItem结构体中
		fileData := gin.H{
			"name":         object.Key,
			"size":         object.Size,
			"lastModified": object.LastModified.Format("2006-01-02 15:04:05"),
			"etag":         strings.Trim(object.ETag, "\""),
			"url":          fmt.Sprintf("/api/v1/buckets/%s/browse/%s", bucketName, object.Key),
			"type":         "file",
		}

		allFileItems = append(allFileItems, FileItem{
			Data:         fileData,
			LastModified: object.LastModified,
			IsFolder:     false,
		})
	}

	// 🔄 修改：添加文件夹对象到FileItem结构体中
	for _, folder := range folders {
		folderData := gin.H{
			"name":         folder,
			"size":         0,
			"lastModified": "",
			"etag":         "",
			"url":          fmt.Sprintf("/api/v1/buckets/%s/browse/%s", bucketName, folder),
			"type":         "folder",
		}

		allFileItems = append(allFileItems, FileItem{
			Data:         folderData,
			LastModified: time.Time{}, // 文件夹使用零时间
			IsFolder:     true,
		})
	}

	// 🔄 修改：新的排序逻辑：文件夹在前，然后文件按时间倒序排列
	sort.Slice(allFileItems, func(i, j int) bool {
		itemI := allFileItems[i]
		itemJ := allFileItems[j]

		// 文件夹始终排在文件前面
		if itemI.IsFolder != itemJ.IsFolder {
			return itemI.IsFolder // 文件夹排在前面
		}

		// 如果都是文件夹，按名称排序
		if itemI.IsFolder && itemJ.IsFolder {
			nameI := itemI.Data["name"].(string)
			nameJ := itemJ.Data["name"].(string)
			return nameI < nameJ
		}

		// 🆕 新增：如果都是文件，按最后修改时间倒序排列（最新的在前）
		if !itemI.IsFolder && !itemJ.IsFolder {
			return itemI.LastModified.After(itemJ.LastModified)
		}

		return false
	})

	// 🔄 修改：转换回gin.H格式
	var allItems []gin.H
	for _, item := range allFileItems {
		allItems = append(allItems, item.Data)
	}

	totalCount := len(allItems)
	totalPages := (totalCount + pageSize - 1) / pageSize

	// 计算分页范围
	offset := (page - 1) * pageSize
	start := offset
	end := offset + pageSize
	if start >= totalCount {
		start = totalCount
	}
	if end > totalCount {
		end = totalCount
	}

	// 获取当前页的对象
	var itemList []gin.H
	if start < totalCount {
		itemList = allItems[start:end]
	}

	// 构建面包屑导航
	breadcrumbs := buildBreadcrumbs(prefix, bucketName)

	c.JSON(http.StatusOK, gin.H{
		"objects": itemList,
		"pagination": gin.H{
			"current_page": page,
			"page_size":    pageSize,
			"total_count":  totalCount,
			"total_pages":  totalPages,
			"has_next":     page < totalPages,
			"has_previous": page > 1,
		},
		"bucket":      bucketName,
		"prefix":      prefix,
		"delimiter":   delimiter,
		"breadcrumbs": breadcrumbs,
	})
}

// buildBreadcrumbs 构建面包屑导航
func buildBreadcrumbs(prefix, bucketName string) []gin.H {
	var breadcrumbs []gin.H

	// 根目录
	breadcrumbs = append(breadcrumbs, gin.H{
		"name": bucketName,
		"path": "",
		"url":  fmt.Sprintf("/api/v1/buckets/%s/objects", bucketName),
	})

	if prefix == "" {
		return breadcrumbs
	}

	// 分割路径
	parts := strings.Split(strings.TrimSuffix(prefix, "/"), "/")
	currentPath := ""

	for _, part := range parts {
		if part == "" {
			continue
		}

		currentPath += part + "/"
		breadcrumbs = append(breadcrumbs, gin.H{
			"name": part,
			"path": currentPath,
			"url":  fmt.Sprintf("/api/v1/buckets/%s/browse/%s", bucketName, currentPath),
		})
	}

	return breadcrumbs
}

// UploadObject 上传对象
func (h *Handler) UploadObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "存储桶名称不能为空",
		})
		return
	}

	ctx := c.Request.Context()

	// 检查存储桶是否存在
	exists, err := h.s3Client.BucketExists(ctx, bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "检查存储桶存在性失败: " + err.Error(),
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "存储桶不存在",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "获取文件失败: " + err.Error(),
		})
		return
	}

	// 文件大小限制（例如100MB）
	const maxFileSize = 100 * 1024 * 1024
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件大小超过限制（100MB）",
		})
		return
	}

	// 自定义对象名称
	objectName := c.PostForm("object_name")
	if objectName == "" {
		objectName = file.Filename
	}

	// 验证文件名
	if err := validateObjectName(objectName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "对象名称格式不正确: " + err.Error(),
		})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "打开文件失败: " + err.Error(),
		})
		return
	}
	defer src.Close()

	// 获取内容类型
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(objectName))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	_, err = h.s3Client.PutObject(ctx, bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "上传文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "文件上传成功",
		"object_name":  objectName,
		"size":         file.Size,
		"bucket":       bucketName,
		"content_type": contentType,
		"url":          fmt.Sprintf("/api/v1/buckets/%s/browse/%s", bucketName, objectName),
	})
}

// UpdateObject 更新对象（实际上是重新上传）
func (h *Handler) UpdateObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	objectName := c.Param("object")

	if bucketName == "" || objectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "存储桶名称或对象名称不能为空",
		})
		return
	}

	ctx := c.Request.Context()

	// 检查对象是否存在
	_, err := h.s3Client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "对象不存在",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "获取文件失败: " + err.Error(),
		})
		return
	}

	// 文件大小限制
	const maxFileSize = 100 * 1024 * 1024
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件大小超过限制（100MB）",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "打开文件失败: " + err.Error(),
		})
		return
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(objectName))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	_, err = h.s3Client.PutObject(ctx, bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "文件更新成功",
		"object_name":  objectName,
		"size":         file.Size,
		"bucket":       bucketName,
		"content_type": contentType,
	})
}

// GetObjectInfo 获取对象信息
func (h *Handler) GetObjectInfo(c *gin.Context) {
	bucketName := c.Param("bucket")
	objectName := c.Param("object")

	if bucketName == "" || objectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "存储桶名称或对象名称不能为空",
		})
		return
	}

	ctx := c.Request.Context()

	stat, err := h.s3Client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "对象不存在或获取信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"object_name":   objectName,
		"bucket":        bucketName,
		"size":          stat.Size,
		"last_modified": stat.LastModified.Format("2006-01-02 15:04:05"),
		"etag":          strings.Trim(stat.ETag, "\""),
		"content_type":  stat.ContentType,
		"url":           fmt.Sprintf("/api/v1/buckets/%s/browse/%s", bucketName, objectName),
	})
}

// DownloadOrListObjects 根据路径是文件还是文件夹，进行下载或列出内容
func (h *Handler) DownloadOrListObjects(c *gin.Context) {
	bucketName := c.Param("bucket")
	filepath := strings.TrimPrefix(c.Param("filepath"), "/")

	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "存储桶名称不能为空",
		})
		return
	}

	ctx := c.Request.Context()

	// 如果路径以 / 结尾，说明是文件夹，列出内容
	if strings.HasSuffix(filepath, "/") || filepath == "" {
		// 模拟文件夹浏览
		originalQuery := c.Request.URL.RawQuery
		if originalQuery != "" {
			c.Request.URL.RawQuery = fmt.Sprintf("%s&prefix=%s&folder=true", originalQuery, filepath)
		} else {
			c.Request.URL.RawQuery = fmt.Sprintf("prefix=%s&folder=true", filepath)
		}
		h.ListObjects(c)
		return
	}

	// 否则尝试作为文件下载
	obj, err := h.s3Client.GetObject(ctx, bucketName, filepath, minio.GetObjectOptions{})
	if err != nil {
		// 如果文件不存在，尝试作为文件夹处理
		originalQuery := c.Request.URL.RawQuery
		if originalQuery != "" {
			c.Request.URL.RawQuery = fmt.Sprintf("%s&prefix=%s/&folder=true", originalQuery, filepath)
		} else {
			c.Request.URL.RawQuery = fmt.Sprintf("prefix=%s/&folder=true", filepath)
		}
		h.ListObjects(c)
		return
	}
	defer obj.Close()

	// 获取对象信息
	stat, err := obj.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "对象不存在或获取信息失败: " + err.Error(),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", strconv.FormatInt(stat.Size, 10))

	// 支持在线预览或强制下载
	disposition := c.Query("disposition")
	if disposition == "attachment" {
		c.Header("Content-Disposition", "attachment; filename=\""+filepath+"\"")
	} else {
		c.Header("Content-Disposition", "inline; filename=\""+filepath+"\"")
	}

	// 流式传输文件
	c.DataFromReader(http.StatusOK, stat.Size, stat.ContentType, obj, nil)
}

// CreateFolder 创建文件夹（实际上是创建一个以/结尾的空对象）
func (h *Handler) CreateFolder(c *gin.Context) {
	bucketName := c.Param("bucket")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "存储桶名称不能为空",
		})
		return
	}

	var req struct {
		FolderName string `json:"folder_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 确保文件夹名以 / 结尾
	folderName := strings.TrimSuffix(req.FolderName, "/") + "/"

	ctx := c.Request.Context()

	// 检查存储桶是否存在
	exists, err := h.s3Client.BucketExists(ctx, bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "检查存储桶存在性失败: " + err.Error(),
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "存储桶不存在",
		})
		return
	}

	// 创建一个空的文件夹对象
	_, err = h.s3Client.PutObject(ctx, bucketName, folderName, strings.NewReader(""), 0, minio.PutObjectOptions{
		ContentType: "application/x-directory",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建文件夹失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "文件夹创建成功",
		"folder_name": folderName,
		"bucket":      bucketName,
		"url":         fmt.Sprintf("/api/v1/buckets/%s/browse/%s", bucketName, folderName),
	})
}
func (h *Handler) DownloadObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	objectName := c.Param("object")

	if bucketName == "" || objectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "存储桶名称或对象名称不能为空",
		})
		return
	}

	ctx := c.Request.Context()
	obj, err := h.s3Client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取对象失败: " + err.Error(),
		})
		return
	}
	defer obj.Close()

	// 获取对象信息
	stat, err := obj.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "对象不存在或获取信息失败: " + err.Error(),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", strconv.FormatInt(stat.Size, 10))

	// 支持在线预览或强制下载
	disposition := c.Query("disposition")
	if disposition == "attachment" {
		c.Header("Content-Disposition", "attachment; filename=\""+objectName+"\"")
	} else {
		c.Header("Content-Disposition", "inline; filename=\""+objectName+"\"")
	}

	// 流式传输文件
	c.DataFromReader(http.StatusOK, stat.Size, stat.ContentType, obj, nil)
}

// DeleteObject 删除对象
func (h *Handler) DeleteObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	objectName := c.Param("object")

	if bucketName == "" || objectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "存储桶名称或对象名称不能为空",
		})
		return
	}

	ctx := c.Request.Context()

	// 检查对象是否存在
	_, err := h.s3Client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "对象不存在",
		})
		return
	}

	err = h.s3Client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除对象失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "对象删除成功",
		"bucket":  bucketName,
		"object":  objectName,
	})
}

// BatchDeleteObjects 批量删除对象
func (h *Handler) BatchDeleteObjects(c *gin.Context) {
	bucketName := c.Param("bucket")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "存储桶名称不能为空",
		})
		return
	}

	var req struct {
		Objects []string `json:"objects" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	if len(req.Objects) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "对象列表不能为空",
		})
		return
	}

	ctx := c.Request.Context()

	// 创建删除通道
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for _, objName := range req.Objects {
			objectsCh <- minio.ObjectInfo{Key: objName}
		}
	}()

	// 批量删除
	errorCh := h.s3Client.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{})

	var failedObjects []string
	for rmErr := range errorCh {
		if rmErr.Err != nil {
			failedObjects = append(failedObjects, rmErr.ObjectName)
		}
	}

	if len(failedObjects) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{
			"message":        "部分对象删除失败",
			"failed_objects": failedObjects,
			"success_count":  len(req.Objects) - len(failedObjects),
			"total_count":    len(req.Objects),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "所有对象删除成功",
		"delete_count": len(req.Objects),
	})
}

// GetBucketsForWeb 为Web界面获取存储桶列表
func (h *Handler) GetBucketsForWeb(c *gin.Context) {
	ctx := c.Request.Context()
	buckets, err := h.s3Client.ListBuckets(ctx)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "layout", gin.H{
			"Buckets": []gin.H{},
			"Error":   "获取存储桶列表失败: " + err.Error(),
		})
		return
	}

	var bucketList []gin.H
	for _, bucket := range buckets {
		bucketList = append(bucketList, gin.H{
			"Name":         bucket.Name,
			"CreationDate": bucket.CreationDate.Format("2006-01-02 15:04:05"),
		})
	}

	c.HTML(http.StatusOK, "layout", gin.H{
		"Buckets": bucketList,
	})
}

// GetBucketObjectsForWeb 为Web界面获取存储桶中的对象
func (h *Handler) GetBucketObjectsForWeb(c *gin.Context) {
	// 统一使用 bucket 参数名
	bucketName := c.Param("bucket")
	if bucketName == "" {
		c.HTML(http.StatusBadRequest, "bucket.html.tmpl", gin.H{
			"BucketName":  "",
			"Objects":     []gin.H{},
			"CurrentPath": "",
			"Paths":       []string{},
			"AllowDelete": true,
			"Error":       "存储桶名称不能为空",
		})
		return
	}

	ctx := c.Request.Context()
	objectsCh := h.s3Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})

	var objectList []gin.H
	for object := range objectsCh {
		if object.Err != nil {
			c.HTML(http.StatusInternalServerError, "bucket.html.tmpl", gin.H{
				"BucketName":  bucketName,
				"Objects":     []gin.H{},
				"CurrentPath": "",
				"Paths":       []string{},
				"AllowDelete": true,
				"Error":       "获取对象列表失败: " + object.Err.Error(),
			})
			return
		}

		// 判断是否为文件夹
		isFolder := strings.HasSuffix(object.Key, "/")
		icon := "insert_drive_file"
		if isFolder {
			icon = "folder"
		}

		objectList = append(objectList, gin.H{
			"Key":          object.Key,
			"DisplayName":  object.Key,
			"Size":         object.Size,
			"Owner":        "owner",
			"LastModified": object.LastModified.Format("2006-01-02 15:04:05"),
			"IsFolder":     isFolder,
			"Icon":         icon,
		})
	}

	// 添加调试信息
	fmt.Printf("Debug: 访问存储桶 %s，找到 %d 个对象\n", bucketName, len(objectList))

	c.HTML(http.StatusOK, "layout", gin.H{
		"BucketName":  bucketName,
		"Objects":     objectList,
		"CurrentPath": "",
		"Paths":       []string{},
		"AllowDelete": true,
	})
}

// 验证存储桶名称
func validateBucketName(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return fmt.Errorf("存储桶名称长度必须在3-63字符之间")
	}
	// 这里可以添加更多的验证规则
	return nil
}

// 验证对象名称
func validateObjectName(name string) error {
	if name == "" {
		return fmt.Errorf("对象名称不能为空")
	}
	if len(name) > 1024 {
		return fmt.Errorf("对象名称长度不能超过1024字符")
	}
	// 这里可以添加更多的验证规则
	return nil
}
