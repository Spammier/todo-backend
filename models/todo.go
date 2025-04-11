package models

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Todo 表示一个待办事项
type Todo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"` // 截止日期
	Completed   bool      `json:"completed" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DB 全局数据库连接
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error

	// 从环境变量获取数据库配置
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPass := getEnvOrDefault("DB_PASSWORD", "")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbName := getEnvOrDefault("DB_NAME", "todo")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移数据库表结构
	err = DB.AutoMigrate(&Todo{}, &User{})
	if err != nil {
		return err
	}

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
