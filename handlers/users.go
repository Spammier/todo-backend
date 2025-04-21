package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"todolist/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// jwtKey 从环境变量获取，不再提供默认值
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值 (此函数不再用于jwtKey)
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Register 用户注册
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("注册JSON绑定失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	fmt.Printf("收到注册请求: username=%s, password=%s, password_length=%d\n",
		req.Username, req.Password, len(req.Password))

	// 检查密码是否为空
	if len(req.Password) == 0 {
		fmt.Println("密码为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := models.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		fmt.Printf("用户名已存在: %s\n", req.Username)
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("密码加密失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	// 创建用户
	if err := models.DB.Create(&user).Error; err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功", "user": user})
}

// Login 用户登录
func Login(c *gin.Context) {
	var loginReq models.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		fmt.Printf("JSON绑定失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	fmt.Printf("收到登录请求: username=%s, password=%s, password_length=%d\n",
		loginReq.Username, loginReq.Password, len(loginReq.Password))

	// 查找用户
	var user models.User
	if err := models.DB.Where("username = ?", loginReq.Username).First(&user).Error; err != nil {
		fmt.Printf("用户查找失败: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	fmt.Printf("找到用户: id=%d, username=%s\n", user.ID, user.Username)
	fmt.Printf("数据库中的密码: %s, 长度: %d\n", user.Password, len(user.Password))
	fmt.Printf("用户输入的密码: %s, 长度: %d\n", loginReq.Password, len(loginReq.Password))

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		fmt.Printf("密码验证失败: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	fmt.Println("密码验证成功，生成JWT令牌")

	// 生成JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Printf("JWT令牌生成失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	fmt.Println("登录成功，返回JWT令牌")

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: tokenString,
		User:  user,
	})
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 移除"Bearer "前缀
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌声明"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))

		c.Next()
	}
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangePassword 修改用户密码
func ChangePassword(c *gin.Context) {
	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 解析请求体
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("修改密码JSON绑定失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 检查新密码是否为空
	if len(req.NewPassword) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "新密码不能为空"})
		return
	}

	// 获取用户信息
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		fmt.Printf("旧密码验证失败: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "旧密码不正确"})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("新密码加密失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "新密码加密失败"})
		return
	}

	// 更新密码
	if err := models.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
