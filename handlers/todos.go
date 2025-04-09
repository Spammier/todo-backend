package handlers

import (
	"net/http"
	"todolist/models"

	"github.com/gin-gonic/gin"
)

// GetAllTodos 返回所有待办事项
func GetAllTodos(c *gin.Context) {
	var todos []models.Todo
	result := models.DB.Find(&todos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取待办事项失败"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// GetTodoByID 根据ID获取待办事项
func GetTodoByID(c *gin.Context) {
	var todo models.Todo
	if err := models.DB.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "待办事项未找到"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// CreateTodo 创建待办事项（支持单个和批量创建）
func CreateTodo(c *gin.Context) {
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
			if err := models.DB.Create(payload.Single).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建待办事项失败"})
				return
			}
			c.JSON(http.StatusCreated, payload.Single)
			return
		}

		// 批量创建
		if len(payload.Batch) > 0 {
			if err := models.DB.Create(&payload.Batch).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "批量创建待办事项失败"})
				return
			}
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

// UpdateTodo 更新待办事项
func UpdateTodo(c *gin.Context) {
	var todo models.Todo
	if err := models.DB.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "待办事项未找到"})
		return
	}

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
	// 显式更新completed字段，即使是false
	updates["completed"] = updatedTodo.Completed

	if err := models.DB.Model(&todo).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新待办事项失败"})
		return
	}

	// 重新获取更新后的待办事项
	models.DB.First(&todo, todo.ID)
	c.JSON(http.StatusOK, todo)
}

// DeleteTodo 删除待办事项
func DeleteTodo(c *gin.Context) {
	var todo models.Todo
	if err := models.DB.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "待办事项未找到"})
		return
	}

	if err := models.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除待办事项失败"})
		return
	}

	c.Status(http.StatusNoContent)
}
