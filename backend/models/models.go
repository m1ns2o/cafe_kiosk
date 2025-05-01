package models

import (
	"time"
)

// 모델 정의
type Category struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"not null" json:"name"` // "음료" 또는 "디저트"
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Menus     []Menu    `json:"menus,omitempty"`
}

type Menu struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    CategoryID uint      `json:"category_id"`
    Name       string    `gorm:"not null" json:"name"`
    Price      int       `gorm:"not null" json:"price"`
    ImageURL   string    `json:"image_url,omitempty"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    // Category   Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

type Order struct {
    ID         uint        `gorm:"primaryKey" json:"id"`
    TotalPrice int         `gorm:"not null" json:"total_price"`
    CreatedAt  time.Time   `json:"created_at"`
    UpdatedAt  time.Time   `json:"updated_at"`
    OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

type OrderItem struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    OrderID   uint      `gorm:"index" json:"order_id"`
    MenuID    uint      `gorm:"index" json:"menu_id"`
    Quantity  int       `gorm:"not null;default:1" json:"quantity"`
    Price     int       `gorm:"not null" json:"price"` // 주문 시점의 가격
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Menu      Menu      `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
    Order     Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// 요청 구조체
// type CreateMenuRequest struct {
//     CategoryID uint   `json:"category_id" binding:"required"`
//     Name       string `json:"name" binding:"required"`
//     Price      int    `json:"price" binding:"required,min=0"`
// }

type CreateMenuRequest struct {
    CategoryID uint   `form:"category_id" json:"category_id" binding:"required"`
    Name       string `form:"name" json:"name" binding:"required"`
    Price      int    `form:"price" json:"price" binding:"required,min=0"`
}


type CreateOrderRequest struct {
    Items []OrderItemRequest `json:"items" binding:"required,min=1"`
}

type OrderItemRequest struct {
    MenuID   uint `json:"menu_id" binding:"required"`
    Quantity int  `json:"quantity" binding:"required,min=1"`
}

type PaymentRequest struct {
    Amount int64 `json:"amount" binding:"required"`
}

// PaymentResponse 결제 응답 구조체
type PaymentResponse struct {
    Success bool                   `json:"success"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}