package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "kiosk/database"
    "kiosk/models"
)

func GetCategories(c *gin.Context) {
    var categories []models.Category
    if err := database.DB.Preload("Menus").Find(&categories).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, categories)
}

func GetCategory(c *gin.Context) {
    id := c.Param("id")
    var category models.Category
    if err := database.DB.Preload("Menus").First(&category, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
        return
    }
    c.JSON(http.StatusOK, category)
}

// 카테고리 생성
func CreateCategory(c *gin.Context) {
    var req struct {
        Name string `json:"name" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    category := models.Category{
        Name: req.Name,
    }
    
    if err := database.DB.Create(&category).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, category)
}

// 카테고리 삭제
func DeleteCategory(c *gin.Context) {
    id := c.Param("id")
    
    // 해당 카테고리의 메뉴가 있는지 확인
    var menuCount int64
    database.DB.Model(&models.Menu{}).Where("category_id = ?", id).Count(&menuCount)
    
    if menuCount > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with existing menus"})
        return
    }
    
    if err := database.DB.Delete(&models.Category{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}