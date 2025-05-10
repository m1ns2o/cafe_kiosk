package handlers

import (
	"encoding/json"
	"fmt"
	"kiosk/database"
	"kiosk/models"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// SSE 클라이언트 관리를 위한 전역 변수
var (
	// 클라이언트 채널을 저장하는 맵
	clients = make(map[chan models.Order]bool)
	// 클라이언트 맵을 안전하게 수정하기 위한 뮤텍스
	clientsMutex sync.Mutex
	// 주문 업데이트를 브로드캐스트하는 채널
	broadcaster = make(chan models.Order)
)

// StartSSEBroadcaster는 SSE 브로드캐스터를 초기화합니다
func StartSSEBroadcaster() {
	go func() {
		for {
			// 주문 업데이트를 기다립니다
			order := <-broadcaster
			
			// 클라이언트 맵에 접근하기 전에 잠금
			clientsMutex.Lock()
			
			// 모든 연결된 클라이언트에게 업데이트 전송
			for clientChan := range clients {
				// 논블로킹 전송
				select {
				case clientChan <- order:
				default:
					// 클라이언트가 수신할 수 없으면 제거
					delete(clients, clientChan)
					close(clientChan)
				}
			}
			
			clientsMutex.Unlock()
		}
	}()
}

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

// GetOrdersByPeriod는 지정된 기간 내의 주문을 조회합니다
func GetOrdersByPeriod(c *gin.Context) {
    // 요청에서 시작일과 종료일 매개변수 가져오기
    startDate := c.Query("start_date")
    endDate := c.Query("end_date")
    
    // 날짜 매개변수 유효성 검사
    if startDate == "" || endDate == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "시작일(start_date)과 종료일(end_date) 매개변수가 모두 필요합니다"})
        return
    }
    
    // 날짜 문자열을 time.Time 객체로 파싱
    start, err := time.Parse("2006-01-02", startDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 시작일 형식. YYYY-MM-DD 형식을 사용하세요"})
        return
    }
    
    end, err := time.Parse("2006-01-02", endDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 종료일 형식. YYYY-MM-DD 형식을 사용하세요"})
        return
    }
    
    // 종료일의 전체 날짜를 포함하기 위해 하루 추가
    end = end.Add(24 * time.Hour)
    
    // 추가 필터링 옵션
    minAmount := c.Query("min_amount")
    maxAmount := c.Query("max_amount")
    menuID := c.Query("menu_id")
    categoryID := c.Query("category_id")
    sortBy := c.Query("sort_by") // 정렬 필드 (created_at, total_price)
    order := c.Query("order")    // 정렬 순서 (asc, desc)
    
    // 기본 쿼리 설정
    query := database.DB.Preload("OrderItems.Menu").Where("orders.created_at BETWEEN ? AND ?", start, end)
    
    // 금액 범위 필터 적용
    if minAmount != "" {
        if minVal, err := strconv.Atoi(minAmount); err == nil {
            query = query.Where("total_price >= ?", minVal)
        }
    }
    
    if maxAmount != "" {
        if maxVal, err := strconv.Atoi(maxAmount); err == nil {
            query = query.Where("total_price <= ?", maxVal)
        }
    }
    
    // 특정 메뉴 ID로 필터링
    if menuID != "" {
        query = query.Joins("JOIN order_items ON orders.id = order_items.order_id").
            Where("order_items.menu_id = ?", menuID).
            Group("orders.id") // 중복 제거
    }
    
    // 특정 카테고리 ID로 필터링
    if categoryID != "" {
        query = query.Joins("JOIN order_items ON orders.id = order_items.order_id").
            Joins("JOIN menus ON order_items.menu_id = menus.id").
            Where("menus.category_id = ?", categoryID).
            Group("orders.id") // 중복 제거
    }
    
    // 정렬 적용
    if sortBy != "" {
        // 허용된 정렬 필드인지 확인
        allowedFields := map[string]bool{
            "created_at":  true,
            "updated_at":  true,
            "total_price": true,
        }
        
        if allowedFields[sortBy] {
            if order != "" && (order == "desc" || order == "asc") {
                query = query.Order(fmt.Sprintf("%s %s", sortBy, order))
            } else {
                // 기본값은 내림차순
                query = query.Order(fmt.Sprintf("%s desc", sortBy))
            }
        }
    } else {
        // 기본 정렬: 생성일 내림차순 (최신순)
        query = query.Order("created_at desc")
    }
    
    // 쿼리 실행
    var orders []models.Order
    if err := query.Find(&orders).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "start_date": startDate,
        "end_date":   endDate,
        "count":      len(orders),
        "orders":     orders,
    })
}

// OrdersEventStream은 주문 업데이트를 위한 SSE 연결을 처리합니다
func OrdersEventStream(c *gin.Context) {
	// SSE를 위한 헤더 설정
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	
	// 이 클라이언트를 위한 채널 생성
	clientChan := make(chan models.Order, 10)
	
	// 클라이언트 등록
	clientsMutex.Lock()
	clients[clientChan] = true
	clientsMutex.Unlock()
	
	// 연결이 닫힐 때 클라이언트 제거
	c.Writer.Flush()
	notify := c.Writer.CloseNotify()
	go func() {
		<-notify
		clientsMutex.Lock()
		delete(clients, clientChan)
		close(clientChan)
		clientsMutex.Unlock()
	}()
	
	// 초기 데이터 전송 - 현재 모든 주문
	var orders []models.Order
	if err := database.DB.Preload("OrderItems.Menu").Find(&orders).Error; err == nil {
		for _, order := range orders {
			data, _ := json.Marshal(order)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()
		}
	}
	
	// 새 업데이트를 클라이언트에게 스트리밍
	for {
		select {
		case order, ok := <-clientChan:
			if !ok {
				return
			}
			
			// 주문을 이벤트로 전송
			data, _ := json.Marshal(order)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()
		case <-time.After(30 * time.Second):
			// 연결 타임아웃을 방지하기 위한 keepalive 주석 전송
			fmt.Fprintf(c.Writer, ": keepalive\n\n")
			c.Writer.Flush()
		}
	}
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

    // 새 주문을 모든 연결된 클라이언트에게 브로드캐스트
    broadcaster <- completeOrder

    c.JSON(http.StatusCreated, completeOrder)
}

// UpdateOrder 함수 (주문 업데이트를 위한 추가 함수)
func UpdateOrder(c *gin.Context) {
    id := c.Param("id")
    var order models.Order
    
    // 기존 주문 조회
    if err := database.DB.Preload("OrderItems.Menu").First(&order, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        return
    }
    
    // 업데이트 로직 구현
    // ... (예: 상태 변경 등)
    
    // 데이터베이스에 업데이트 저장
    if err := database.DB.Save(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // 업데이트된 주문 브로드캐스트
    broadcaster <- order
    
    c.JSON(http.StatusOK, order)
}