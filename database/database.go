package database

import (
    "github.com/glebarez/sqlite"  // CGO-free SQLite driver
    "gorm.io/gorm"
    "kiosk/models"
)

var DB *gorm.DB

func InitDB() error {
    var err error
    DB, err = gorm.Open(sqlite.Open("kiosk.db"), &gorm.Config{})
    if err != nil {
        return err
    }

    // 테이블 자동 생성
    err = DB.AutoMigrate(&models.Category{}, &models.Menu{}, &models.Order{}, &models.OrderItem{})
    if err != nil {
        return err
    }

    return nil
}

func InitializeCategories() {
    categories := []models.Category{
        {Name: "음료"},
        {Name: "디저트"},
    }

    for _, category := range categories {
        var count int64
        DB.Model(&models.Category{}).Where("name = ?", category.Name).Count(&count)
        if count == 0 {
            DB.Create(&category)
        }
    }
}