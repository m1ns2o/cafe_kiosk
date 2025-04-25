package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "kiosk/database"
    "kiosk/models"
)

func GetMenus(c *gin.Context) {
    var menus []models.Menu
    query := database.DB.Preload("Category")
    
    // 카테고리별 필터링
    if categoryID := c.Query("category_id"); categoryID != "" {
        query = query.Where("category_id = ?", categoryID)
    }
    
    // 카테고리 이름으로 필터링 (옵션)
    if categoryName := c.Query("category"); categoryName != "" {
        query = query.Joins("JOIN categories ON categories.id = menus.category_id").
            Where("categories.name = ?", categoryName)
    }
    
    if err := query.Find(&menus).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, menus)
}

func GetMenusByCategory(c *gin.Context) {
    categoryID := c.Param("id")  // :category_id를 :id로 변경
    var menus []models.Menu
    
    if err := database.DB.Where("category_id = ?", categoryID).Find(&menus).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, menus)
}

func GetMenu(c *gin.Context) {
    id := c.Param("id")
    var menu models.Menu
    if err := database.DB.Preload("Category").First(&menu, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
        return
    }
    c.JSON(http.StatusOK, menu)
}

func CreateMenu(c *gin.Context) {
    var req models.CreateMenuRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    menu := models.Menu{
        CategoryID: req.CategoryID,
        Name:       req.Name,
        Price:      req.Price,
    }

    if err := database.DB.Create(&menu).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, menu)
}

func UpdateMenu(c *gin.Context) {
    id := c.Param("id")
    var menu models.Menu
    if err := database.DB.First(&menu, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
        return
    }

    var req models.CreateMenuRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    menu.CategoryID = req.CategoryID
    menu.Name = req.Name
    menu.Price = req.Price

    if err := database.DB.Save(&menu).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, menu)
}

func DeleteMenu(c *gin.Context) {
    id := c.Param("id")
    if err := database.DB.Delete(&models.Menu{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
}