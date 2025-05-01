package handlers

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "kiosk/database"
    "kiosk/models"
)

func GetOrders(c *gin.Context) {
    var orders []models.Order
    if err := database.DB.Preload("OrderItems.Menu").Find(&orders).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, orders)
}

func GetOrder(c *gin.Context) {
    id := c.Param("id")
    var order models.Order
    if err := database.DB.Preload("OrderItems.Menu").First(&order, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        return
    }
    c.JSON(http.StatusOK, order)
}

func CreateOrder(c *gin.Context) {
    var req models.CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 트랜잭션 시작
    tx := database.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 주문 항목 데이터 준비 및 총액 계산
    var totalPrice int
    orderItems := make([]models.OrderItem, 0, len(req.Items))

    for _, item := range req.Items {
        var menu models.Menu
        if err := tx.First(&menu, item.MenuID).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Menu ID %d not found", item.MenuID)})
            return
        }

        // 아직 OrderID는 설정하지 않음 (주문 생성 후 설정)
        orderItem := models.OrderItem{
            MenuID:   item.MenuID,
            Quantity: item.Quantity,
            Price:    menu.Price,
        }
        orderItems = append(orderItems, orderItem)
        totalPrice += menu.Price * item.Quantity
    }

    // 총액이 계산된 후 주문 생성
    order := models.Order{
        TotalPrice: totalPrice,
    }
    if err := tx.Create(&order).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 주문 ID 설정 및 주문 항목 저장
    for i := range orderItems {
        orderItems[i].OrderID = order.ID
    }

    // 주문 항목 일괄 생성
    if err := tx.Create(&orderItems).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 트랜잭션 커밋
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 생성된 주문 조회 (Preload 사용, Category 제외)
    var completeOrder models.Order
    if err := database.DB.Preload("OrderItems.Menu").First(&completeOrder, order.ID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, completeOrder)
}