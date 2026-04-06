package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	_ "bluebell/docs" // 导入swagger文档
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 路由
func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 配置 CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg":     "Bluebell API Server is running",
			"version": "v2.0.0",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	// 静态文件：上传的图片
	r.Static("/uploads", "./uploads")

	v1 := r.Group("/api/v1")

	// 认证相关（无需登录）
	v1.POST("/auth/login", controller.LoginHandler)
	v1.POST("/auth/signup", controller.SignUpHandler)
	v1.POST("/auth/refresh", controller.RefreshTokenHandler) // 新增刷新token接口

	// 文章相关（访客可访问）
	v1.GET("/articles", controller.GetArticleListHandler)
	v1.GET("/articles/featured", controller.GetFeaturedArticlesHandler)
	v1.GET("/articles/search", controller.SearchArticlesHandler)
	v1.GET("/articles/trending", controller.GetTrendingArticlesHandler) // 新增热门文章排行

	// 使用查询参数替代路径参数，避免Gin v1.6.3路由冲突
	v1.GET("/article-stats/daily", controller.GetArticleDailyStatsHandler)
	v1.GET("/article-stats/trend", controller.GetArticleTrendHandler) // 新增访问趋势
	v1.GET("/article-stats/batch", controller.BatchGetArticleStatsHandler) // 新增批量统计
	v1.POST("/articles/view", controller.RecordArticleViewWithAntiCheatHandler)
	v1.GET("/articles/:id", middlewares.OptionalJWTMiddleware(), controller.GetArticleDetailHandler)

	// 标签相关（访客可访问）
	v1.GET("/tags", controller.GetAllTagsHandler)
	v1.GET("/tags/:id/articles", controller.GetArticlesByTagHandler)

	// 栏目相关（访客可访问）
	v1.GET("/categories", controller.GetCategoryListHandler)
	v1.GET("/categories/:id", controller.GetCategoryByIdHandler)
	v1.GET("/categories/:id/articles", controller.GetArticlesByCategoryHandler)
	v1.GET("/articles/:id/categories", controller.GetCategoriesByArticleHandler)

	// 作者主页（公开接口）
	v1.GET("/authors/:username", controller.GetAuthorInfoHandler)
	v1.GET("/authors/:username/articles", controller.GetAuthorArticlesListHandler)

	// RSS订阅（公开接口）
	v1.GET("/rss", controller.GetRSSHandler)

	// 评论列表（公开接口，访客可查看）
	v1.GET("/comments", controller.GetCommentListHandler)

	// 需要登录的接口
	v1.Use(middlewares.JWTAuthMiddleware())

	{
		// 用户相关
		v1.GET("/auth/profile", controller.GetUserProfileHandler)
		v1.PUT("/auth/profile", controller.UpdateUserProfileHandler)

		// 作者统计（需要登录）
		v1.GET("/author/stats/trend", controller.GetAuthorTrendHandler)

		// 文章相关（需要作者权限）
		v1.POST("/articles", controller.CreateArticleHandler)
		v1.PUT("/author/articles/:id", controller.UpdateArticleHandler)
		v1.DELETE("/author/articles/:id", controller.DeleteArticleHandler)
		v1.PATCH("/author/articles/:id/status", controller.UpdateArticleStatusHandler)
		v1.PATCH("/author/articles/:id/featured", controller.UpdateArticleFeaturedHandler)
		v1.GET("/author/articles", controller.GetAuthorArticlesHandler)

		// 文章导出（需要作者权限）
		v1.GET("/author/articles/:id/export", controller.ExportArticleHandler)
		v1.POST("/author/articles/export", controller.ExportArticlesBatchHandler)

		// 标签相关（需要作者权限）
		v1.POST("/tags", controller.CreateTagHandler)
		v1.PUT("/tags/:id", controller.UpdateTagHandler)
		v1.DELETE("/tags/:id", controller.DeleteTagHandler)
		v1.GET("/author/tags", controller.GetAuthorTagsHandler)

		// 评论相关（需要阅读者权限）
		v1.POST("/comments", controller.CreateCommentHandler)
		v1.PUT("/comments/:id", controller.UpdateCommentHandler)
		v1.DELETE("/comments/:id", controller.DeleteCommentHandler)

		// 点赞相关（需要阅读者权限）
		v1.POST("/likes", controller.LikeHandler)
		v1.DELETE("/likes", controller.UnlikeHandler)
		v1.GET("/likes/status", controller.GetLikeStatusHandler)
		v1.GET("/user/likes", controller.GetUserLikesHandler)

		// 文件上传
		v1.POST("/upload/image", controller.UploadImageHandler)
		v1.POST("/upload/attachment", controller.UploadAttachmentHandler)

		// 管理员接口（需要管理员权限）
		admin := v1.Group("/admin")
		admin.Use(middlewares.AdminAuthMiddleware()) // 添加管理员权限验证
		{
			admin.GET("/articles", controller.GetAdminArticlesHandler)
			admin.PATCH("/articles/:id/featured", controller.AdminSetArticleFeaturedHandler)
			admin.GET("/users", controller.GetUserListHandler)
			admin.GET("/users/:id", controller.GetUserDetailHandler)
			admin.PATCH("/users/:id/status", controller.UpdateUserStatusHandler)
			admin.PATCH("/users/:id/role", controller.UpdateUserRoleHandler)
			admin.PATCH("/users/batch/status", controller.BatchUpdateUserStatusHandler)
			admin.GET("/stats/overview", controller.GetSystemOverviewHandler)
			admin.GET("/stats/daily", controller.GetSystemDailyStatsHandler)
			admin.GET("/metrics/realtime", controller.GetSystemMetricsHandler)
			admin.GET("/metrics/history", controller.GetSystemMetricsHistoryHandler)
		}

		// 栏目管理（需要登录）
		v1.POST("/categories", controller.CreateCategoryHandler)
		v1.PUT("/categories/:id", controller.UpdateCategoryHandler)
		v1.DELETE("/categories/:id", controller.DeleteCategoryHandler)
		v1.POST("/articles/:id/categories", controller.AddArticleToCategoriesHandler)
	}

	pprof.Register(r) // 注册pprof相关路由

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
