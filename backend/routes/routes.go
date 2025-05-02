package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "kiosk/handlers"
)

func SetupRoutes(r *gin.Engine) {
    r.Use(cors.New(cors.Config{
        AllowAllOrigins:  true,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
    }))
    api := r.Group("/api")
    {
        // 카테고리 관련
        api.GET("/categories", handlers.GetCategories)
        api.GET("/categories/:id", handlers.GetCategory)
        api.POST("/categories", handlers.CreateCategory)  // 카테고리 추가
        api.DELETE("/categories/:id", handlers.DeleteCategory)  // 카테고리 삭제
        api.GET("/categories/:id/menus", handlers.GetMenusByCategory)

        // 메뉴 관련
        api.GET("/menus", handlers.GetMenus)
        api.GET("/menus/:id", handlers.GetMenu)
        api.POST("/menus", handlers.CreateMenu)
        api.PUT("/menus/:id", handlers.UpdateMenu)
        api.DELETE("/menus/:id", handlers.DeleteMenu)

        // 주문 관련
        api.GET("/orders", handlers.GetOrders)
        api.GET("/orders/:id", handlers.GetOrder)
        api.POST("/orders", handlers.CreateOrder)

        // 결제 관련
        api.POST("/payment", handlers.ProcessPayment)
    }
}