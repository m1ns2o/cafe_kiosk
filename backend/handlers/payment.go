package handlers

import (
    "fmt"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "kiosk/models"
)

// DepositStateInterface 인터페이스 정의 (main.go의 DepositState 구조체와 동일한 메서드)
type DepositStateInterface interface {
    UpdateAndCheckDeposit(expectedAmount int64) (bool, int64, error)
}

func ProcessPayment(c *gin.Context) {
    var req models.PaymentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // DepositState 인스턴스 가져오기
    depositStateInterface, exists := c.Get("depositState")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "deposit state not found"})
        return
    }

    depositState, ok := depositStateInterface.(DepositStateInterface)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid deposit state type"})
        return
    }

    // 50초 동안 2초 간격으로 체크 (총 25회)
    maxAttempts := 60
    interval := 3 * time.Second
    success := false
    var actualChange int64
    var err error

    for attempt := 1; attempt <= maxAttempts; attempt++ {
        success, actualChange, err = depositState.UpdateAndCheckDeposit(req.Amount)
        
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": fmt.Sprintf("error checking deposit: %v", err),
            })
            return
        }

        if success {
            break
        }

        if attempt < maxAttempts {
            time.Sleep(interval)
        }
    }

    if success {
        response := models.PaymentResponse{
            Success: true,
            Message: "Payment verified successfully",
            Details: map[string]interface{}{
                "expected_amount": req.Amount,
                "actual_change":   actualChange,
                "verified_at":     time.Now().Format(time.RFC3339),
            },
        }
        c.JSON(http.StatusOK, response)
    } else {
        response := models.PaymentResponse{
            Success: false,
            Message: "Payment verification timeout",
            Details: map[string]interface{}{
                "expected_amount": req.Amount,
                "actual_change":   actualChange,
                "timeout_after":   "50 seconds",
            },
        }
        c.JSON(http.StatusRequestTimeout, response)
    }
}