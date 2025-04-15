package main

import (
	"fmt"
	"log"
	"os"
	"todolist/handlers"
	"todolist/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // 取消注释，重新导入 godotenv 库
)

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// 取消注释，在程序启动时加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		// 如果 .env 文件不存在，打印日志，但允许继续，依赖系统环境变量或默认值
		log.Println("未找到 .env 文件，将使用系统环境变量或默认值")
	}

	// 初始化数据库连接
	if err := models.InitDB(); err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// API基础路由组
	api := r.Group("/api")
	{
		// 公开路由，不需要认证
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		// 需要认证的路由组
		auth := api.Group("")
		auth.Use(handlers.AuthMiddleware())
		{
			// 用户相关路由
			auth.POST("/change-password", handlers.ChangePassword)

			// Todo相关路由
			todos := auth.Group("/todos")
			{
				todos.GET("", handlers.GetAllTodos)
				todos.GET("/:id", handlers.GetTodoByID)
				todos.POST("", handlers.CreateTodo)
				todos.PUT("/:id", handlers.UpdateTodo)
				todos.DELETE("/:id", handlers.DeleteTodo)
			}
		}
	}

	// 获取端口配置
	port := getEnvOrDefault("PORT", "8080")

	// 启动服务器
	fmt.Printf("服务器已启动，监听端口 %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
