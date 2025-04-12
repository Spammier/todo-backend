package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"todolist/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	// 缓存过期时间
	cacheDuration = time.Minute * 10
)

// ---- 缓存 Key 生成函数 ----

// getUserTodosKey 生成用户待办事项列表的缓存Key
func getUserTodosKey(userID uint) string {
	return fmt.Sprintf("user:%d:todos", userID)
}

// getTodoKey 生成单个待办事项的缓存Key
func getTodoKey(todoID uint) string {
	return fmt.Sprintf("todo:%d", todoID)
}

// ---- 缓存清除函数 ----

// clearUserCache 清除指定用户的所有相关缓存
func clearUserCache(userID uint) {
	userTodosKey := getUserTodosKey(userID)
	models.Rdb.Del(models.Ctx, userTodosKey)
	// 如果有其他与用户相关的缓存，也在此处清除
}

// clearTodoCache 清除单个待办事项的缓存
func clearTodoCache(todoID uint) {
	models.Rdb.Del(models.Ctx, getTodoKey(todoID))
}

// GetAllTodos 返回当前用户的待办事项 (带缓存)
func GetAllTodos(c *gin.Context) {
	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	currentUserID := userID.(uint)

	// --- 缓存读取 ---
	cacheKey := getUserTodosKey(currentUserID)
	cachedTodos, err := models.Rdb.Get(models.Ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中
		var todos []models.Todo
		if json.Unmarshal([]byte(cachedTodos), &todos) == nil {
			c.JSON(http.StatusOK, todos)
			fmt.Println("Cache hit for key:", cacheKey) // 日志
			return
		} else {
			// JSON解析失败，可能缓存数据损坏，继续查询数据库
			fmt.Println("Cache data corrupted for key:", cacheKey)
		}
	} else if err != redis.Nil {
		// Redis查询出错 (非 'key not found')
		fmt.Printf("Redis Get error for key %s: %v\n", cacheKey, err)
		// 出错时可以选择直接查询数据库，或者返回错误
	}
	fmt.Println("Cache miss for key:", cacheKey) // 日志

	// --- 缓存未命中或出错，查询数据库 ---
	var todos []models.Todo
	result := models.DB.Where("user_id = ?", currentUserID).Find(&todos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取待办事项失败"})
		return
	}

	// --- 结果存入缓存 ---
	todosJSON, err := json.Marshal(todos)
	if err == nil {
		err = models.Rdb.Set(models.Ctx, cacheKey, todosJSON, cacheDuration).Err()
		if err != nil {
			fmt.Printf("Redis Set error for key %s: %v\n", cacheKey, err)
			// 缓存写入失败不应阻塞主流程，记录日志即可
		} else {
			fmt.Println("Cache set for key:", cacheKey)
		}
	} else {
		fmt.Printf("JSON Marshal error when caching todos for user %d: %v\n", currentUserID, err)
	}

	c.JSON(http.StatusOK, todos)
}

// GetTodoByID 根据ID获取当前用户的待办事项 (带缓存)
func GetTodoByID(c *gin.Context) {
	// 从上下文中获取当前用户ID
	userID, userExists := c.Get("user_id")
	if !userExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	currentUserID := userID.(uint)

	todoIDStr := c.Param("id")
	var todoID uint
	fmt.Sscan(todoIDStr, &todoID) // 简单转换，实际应处理错误
	if todoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的待办事项ID"})
		return
	}

	// --- 缓存读取 ---
	cacheKey := getTodoKey(todoID)
	cachedTodo, err := models.Rdb.Get(models.Ctx, cacheKey).Result()
	if err == nil {
		var todo models.Todo
		if json.Unmarshal([]byte(cachedTodo), &todo) == nil {
			// 检查缓存中的Todo是否属于当前用户
			if todo.UserID == currentUserID {
				c.JSON(http.StatusOK, todo)
				fmt.Println("Cache hit for key:", cacheKey) // 日志
				return
			} else {
				// 用户无权访问此缓存项
				c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该待办事项"})
				return
			}
		} else {
			fmt.Println("Cache data corrupted for key:", cacheKey)
		}
	} else if err != redis.Nil {
		fmt.Printf("Redis Get error for key %s: %v\n", cacheKey, err)
	}
	fmt.Println("Cache miss for key:", cacheKey)

	// --- 缓存未命中，查询数据库 ---
	var todo models.Todo
	if err := models.DB.Where("id = ? AND user_id = ?", todoID, currentUserID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "待办事项未找到或无权访问"})
		return
	}

	// --- 结果存入缓存 ---
	todoJSON, err := json.Marshal(todo)
	if err == nil {
		err = models.Rdb.Set(models.Ctx, cacheKey, todoJSON, cacheDuration).Err()
		if err != nil {
			fmt.Printf("Redis Set error for key %s: %v\n", cacheKey, err)
		} else {
			fmt.Println("Cache set for key:", cacheKey)
		}
	} else {
		fmt.Printf("JSON Marshal error when caching todo %d: %v\n", todoID, err)
	}

	c.JSON(http.StatusOK, todo)
}

// CreateTodo 创建待办事项（支持单个和批量创建）(带缓存清除)
func CreateTodo(c *gin.Context) {
	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	currentUserID := userID.(uint)

	contentType := c.GetHeader("Content-Type")
	if contentType == "application/json" {
		var payload struct {
			Single *models.Todo  `json:"todo"`
			Batch  []models.Todo `json:"todos"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		// 单个创建
		if payload.Single != nil {
			payload.Single.UserID = currentUserID
			if err := models.DB.Create(payload.Single).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建待办事项失败"})
				return
			}
			// --- 清除用户列表缓存 ---
			clearUserCache(currentUserID)
			fmt.Println("Cache cleared for user:", currentUserID) // 日志
			c.JSON(http.StatusCreated, payload.Single)
			return
		}

		// 批量创建
		if len(payload.Batch) > 0 {
			for i := range payload.Batch {
				payload.Batch[i].UserID = currentUserID
			}
			if err := models.DB.Create(&payload.Batch).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "批量创建待办事项失败"})
				return
			}
			// --- 清除用户列表缓存 ---
			clearUserCache(currentUserID)
			fmt.Println("Cache cleared for user:", currentUserID) // 日志
			c.JSON(http.StatusCreated, gin.H{
				"message": "批量创建成功",
				"todos":   payload.Batch,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体必须包含 todo 或 todos 字段"})
		return
	}

	c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "仅支持 application/json"})
}

// UpdateTodo 更新当前用户的待办事项 (带缓存清除)
func UpdateTodo(c *gin.Context) {
	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	currentUserID := userID.(uint)

	var todo models.Todo
	// 查找属于当前用户的待办事项
	if err := models.DB.Where("id = ? AND user_id = ?", c.Param("id"), currentUserID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "待办事项未找到或无权更新"})
		return
	}
	originalTodoID := todo.ID // 保存原始ID用于缓存清除

	var updatedTodo models.Todo
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 更新具体字段，包括零值字段
	updates := map[string]interface{}{}
	if updatedTodo.Title != "" {
		updates["title"] = updatedTodo.Title
	}
	if updatedTodo.Description != "" {
		updates["description"] = updatedTodo.Description
	}
	updates["completed"] = updatedTodo.Completed
	if err := models.DB.Model(&todo).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新待办事项失败"})
		return
	}

	// --- 清除相关缓存 ---
	clearUserCache(currentUserID)                                                        // 清除用户列表缓存
	clearTodoCache(originalTodoID)                                                       // 清除单个待办事项缓存
	fmt.Printf("Cache cleared for user %d and todo %d\n", currentUserID, originalTodoID) // 日志

	// 重新获取更新后的待办事项 (这一步会触发缓存写入)
	models.DB.First(&todo, originalTodoID)
	c.JSON(http.StatusOK, todo)
}

// DeleteTodo 删除当前用户的待办事项 (带缓存清除)
func DeleteTodo(c *gin.Context) {
	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	currentUserID := userID.(uint)

	var todo models.Todo
	// 查找属于当前用户的待办事项
	if err := models.DB.Where("id = ? AND user_id = ?", c.Param("id"), currentUserID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "待办事项未找到或无权删除"})
		return
	}
	deletedTodoID := todo.ID // 保存ID用于缓存清除

	if err := models.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除待办事项失败"})
		return
	}

	// --- 清除相关缓存 ---
	clearUserCache(currentUserID)                                                       // 清除用户列表缓存
	clearTodoCache(deletedTodoID)                                                       // 清除单个待办事项缓存
	fmt.Printf("Cache cleared for user %d and todo %d\n", currentUserID, deletedTodoID) // 日志

	c.Status(http.StatusNoContent)
}
