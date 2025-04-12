package models

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Todo 表示一个待办事项
type Todo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DB 全局数据库连接
var DB *gorm.DB

// Rdb 全局Redis客户端连接
var Rdb *redis.Client

// Ctx Redis操作的上下文
var Ctx = context.Background()

// InitDB 初始化数据库和Redis连接
func InitDB() error {
	var err error

	// 初始化数据库连接
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPass := getEnvOrDefault("DB_PASSWORD", "")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbName := getEnvOrDefault("DB_NAME", "todo")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 自动迁移数据库表结构
	err = DB.AutoMigrate(&Todo{}, &User{})
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 初始化Redis连接
	redisAddr := getEnvOrDefault("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnvOrDefault("REDIS_PASSWORD", "")
	redisDBStr := getEnvOrDefault("REDIS_DB", "0")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		redisDB = 0 // 如果转换失败，使用默认值
		fmt.Printf("Redis DB 配置无效，使用默认值 0: %v\n", err)
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// 测试Redis连接
	_, err = Rdb.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis 连接失败: %w", err)
	}
	fmt.Println("成功连接到 Redis")

	return nil
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
