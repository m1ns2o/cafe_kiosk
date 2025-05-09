package handlers

import (
    "encoding/json"
    "fmt"
    "kiosk/models"
    "kiosk/utils"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "sync"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var (
    // 로그 파일 핸들러
    paymentLogFile *os.File
    // 로거 인스턴스
    paymentLogger *log.Logger
)

// 로그 시스템 초기화
func InitLogSystem() error {
    // 로그 디렉토리 생성
    logDir := "logs"
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return fmt.Errorf("로그 디렉토리 생성 실패: %v", err)
    }
    
    // 오늘 날짜로 로그 파일 생성
    today := time.Now().Format("2006-01-02")
    logFilePath := filepath.Join(logDir, fmt.Sprintf("payment_%s.log", today))
    
    // 로그 파일 열기 (없으면 생성, 있으면 추가 모드)
    file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("로그 파일 생성 실패: %v", err)
    }
    
    // 로거 설정
    paymentLogFile = file
    paymentLogger = log.New(file, "", log.LstdFlags)
    
    // 콘솔과 파일 모두에 로깅하기 위해 멀티 라이터 설정
    // log.SetOutput(io.MultiWriter(os.Stdout, file))
    
    paymentLogger.Println("결제 로그 시스템 초기화 완료")
    return nil
}

// 로그 시스템 종료
func CloseLogSystem() {
    if paymentLogFile != nil {
        paymentLogFile.Close()
    }
}

// 로그 메시지 기록 (콘솔 및 파일)
func logMessage(format string, args ...interface{}) {
    message := fmt.Sprintf(format, args...)
    log.Println(message)  // 콘솔에 출력
    
    if paymentLogger != nil {
        paymentLogger.Println(message)  // 파일에 기록
    }
}

// 웹소켓 연결 업그레이더 정의
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // CORS 허용 설정 (실제 환경에서는 보안 설정 필요)
    },
}

// 웹소켓 메시지 타입 정의
const (
    MsgTypePaymentRequest = "payment_request"
    MsgTypePaymentStatus  = "payment_status"
    MsgTypePaymentResult  = "payment_result"
    MsgTypeError          = "error"
    MsgTypeCancelRequest  = "cancel_request"
    MsgTypeCancelResult   = "cancel_result"
)

// 웹소켓 메시지 구조체
type WebSocketMessage struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

// PaymentStatus 구조체 - 결제 과정 중 상태 업데이트
type PaymentStatus struct {
    Attempt      int   `json:"attempt"`
    MaxAttempts  int   `json:"max_attempts"`
    ActualChange int64 `json:"actual_change,omitempty"`
}

// PaymentHandler 웹소켓 핸들러
func PaymentHandler(c *gin.Context) {
    // 웹소켓으로 업그레이드
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        logMessage("웹소켓 업그레이드 실패: %v", err)
        return
    }
    defer conn.Close()

    // DepositState 인스턴스 가져오기
    depositStateInterface, exists := c.Get("depositState")
    if !exists {
        sendError(conn, "deposit state not found")
        return
    }

    depositState, ok := depositStateInterface.(*utils.DepositState)
    if !ok {
        sendError(conn, "invalid deposit state type")
        return
    }

    // 클라이언트로부터 메시지 수신 대기
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("웹소켓 에러: %v \n", err)
            }
            break
        }

        // 메시지 파싱
        var wsMsg WebSocketMessage
        if err := json.Unmarshal(message, &wsMsg); err != nil {
            sendError(conn, fmt.Sprintf("메시지 파싱 실패: %v", err))
            continue
        }

        // 메시지 타입에 따른 처리
        switch wsMsg.Type {
        case MsgTypePaymentRequest:
            // PaymentRequest로 변환
            payloadBytes, err := json.Marshal(wsMsg.Payload)
            if err != nil {
                sendError(conn, "invalid payment request format")
                continue
            }

            var req models.PaymentRequest
            if err := json.Unmarshal(payloadBytes, &req); err != nil {
                sendError(conn, "invalid payment request data")
                continue
            }

            // 비동기로 결제 처리 시작
            go processPaymentWithWebSocket(conn, depositState, req)

        case MsgTypeCancelRequest:
            // 결제 취소 요청 처리
            sendMessage(conn, nil, MsgTypeCancelResult, gin.H{
                "success": true,
                "message": "결제가 취소되었습니다",
            })
            logMessage("결제 취소 요청 처리 완료")
            
        default:
            sendError(conn, "unknown message type")
        }
    }
}

// 웹소켓을 통한 결제 처리
func processPaymentWithWebSocket(conn *websocket.Conn, depositState *utils.DepositState, req models.PaymentRequest) {
    // 웹소켓 메시지 전송을 위한 뮤텍스 (동시 전송 방지)
    var mutex sync.Mutex

    // 초기 예수금 조회 및 로깅
    initialDeposit := depositState.GetCurrentDeposit()
    log.Printf("결제 요청 시작 - 요청 금액: %s원, 초기 예수금: %s원\n", 
        utils.FormatNumber(req.Amount), utils.FormatNumber(initialDeposit))

    // 결제 시작 전 한 번 더 예수금 상태 조회 및 업데이트
    // 기존 예수금 잔고와 값이 동일한지 확인
    latestDepositAmount, err := depositState.GetKISDepositAmount()
    if err != nil {
        logMessage("초기 예수금 재확인 실패: %v", err)
        sendError(conn, fmt.Sprintf("초기 예수금 조회 오류: %v", err))
        return
    }

    // 기존 예수금 값과 최신 예수금 값 비교
    if initialDeposit != latestDepositAmount {
        logMessage("[중요] 예수금 일치하지 않음: 내부 상태 (%s원) vs KIS API (%s원), 차이: %s원",
            utils.FormatNumber(initialDeposit),
            utils.FormatNumber(latestDepositAmount),
            utils.FormatNumber(latestDepositAmount-initialDeposit))
            
        // 예수금 상태 강제 업데이트 (내부 상태 조정)
        depositState.SetCurrentDeposit(latestDepositAmount)
        initialDeposit = latestDepositAmount
        
        logMessage("예수금 상태를 KIS API 값으로 업데이트 완료: %s원", utils.FormatNumber(latestDepositAmount))
    } else {
        log.Printf("예수금 상태 일치 확인 완료: %s원\n", utils.FormatNumber(initialDeposit))
    }

    // 결제 처리 파라미터 설정
    maxAttempts := 20
    interval := 1 * time.Second
    success := false
    var actualChange int64
    
    // 이전 예수금 상태 유지
    previousDeposit := initialDeposit

    // 결제 처리 시작 시간
    startTime := time.Now()

    for attempt := 1; attempt <= maxAttempts; attempt++ {
        // 최신 예수금 조회를 통한 예수금 업데이트
        success, actualChange, err = depositState.UpdateAndCheckDeposit(req.Amount)
        
        if err != nil {
            logMessage("결제 검증 오류: %v", err)
            sendError(conn, fmt.Sprintf("error checking deposit: %v", err))
            return
        }

        // 현재 예수금 조회 (변동 후)
        currentDeposit := depositState.GetCurrentDeposit()

        // 상태 업데이트 전송
        status := PaymentStatus{
            Attempt:      attempt,
            MaxAttempts:  maxAttempts,
            ActualChange: actualChange,
        }
        sendMessage(conn, &mutex, MsgTypePaymentStatus, status)

        // 상태 로깅
        log.Printf("결제 확인 시도 #%d - 예상 증가액: %s원, 실제 증가액: %s원, 현재 예수금: %s원\n", 
            attempt, utils.FormatNumber(req.Amount), utils.FormatNumber(actualChange), 
            utils.FormatNumber(currentDeposit))

        // 결제 금액 검증 - 예상 금액과 실제 변동액 비교
        if success {
            log.Printf("결제 성공 - 요청 금액: %s원, 실제 변동액: %s원, 소요 시간: %v\n", 
                utils.FormatNumber(req.Amount), utils.FormatNumber(actualChange), time.Since(startTime))
            break
        } else if(actualChange !=0){
            logMessage("결제 실패 - 요청 금액: %s원, 최종 변동액: %s원, 타임아웃: %v초", 
                utils.FormatNumber(req.Amount), utils.FormatNumber(actualChange), time.Since(startTime))
                
        }

        // 이전 예수금과 현재 예수금(변동액 제외) 비교
        currentDepositWithoutChange := currentDeposit - actualChange
        if previousDeposit != currentDepositWithoutChange {
            logMessage("[주의] 예수금 불일치 감지 - 이전: %s원, 현재(변동 전): %s원, 차이: %s원", 
                utils.FormatNumber(previousDeposit), 
                utils.FormatNumber(currentDepositWithoutChange), 
                utils.FormatNumber(currentDepositWithoutChange - previousDeposit))
        }
        
        // 이전 예수금 값 업데이트
        previousDeposit = currentDeposit

        if attempt < maxAttempts {
            time.Sleep(interval)
        }
    }

    // 결과 전송
    if success {
        response := models.PaymentResponse{
            Success: true,
            Message: "결제가 성공적으로 확인되었습니다",
            Details: map[string]interface{}{
                "expected_amount": req.Amount,
                "actual_change":   actualChange,
                "verified_at":     time.Now().Format(time.RFC3339),
                "elapsed_time":    time.Since(startTime).String(),
            },
        }
        sendMessage(conn, &mutex, MsgTypePaymentResult, response)
    } else {
        // 결제 실패 로깅
        logMessage("[중요] 결제 실패 - 요청 금액: %s원, 최종 변동액: %s원, 타임아웃: %v초", 
            utils.FormatNumber(req.Amount), utils.FormatNumber(actualChange), 
            maxAttempts*int(interval/time.Second))

        response := models.PaymentResponse{
            Success: false,
            Message: "결제 확인 시간 초과",
            Details: map[string]interface{}{
                "expected_amount": req.Amount,
                "actual_change":   actualChange,
                "timeout_after":   fmt.Sprintf("%d초", maxAttempts*int(interval/time.Second)),
                "elapsed_time":    time.Since(startTime).String(),
            },
        }
        sendMessage(conn, &mutex, MsgTypePaymentResult, response)
    }
}

// 웹소켓 에러 메시지 전송
func sendError(conn *websocket.Conn, errorMsg string) {
    var mutex sync.Mutex
    logMessage("에러: %s", errorMsg)
    sendMessage(conn, &mutex, MsgTypeError, gin.H{"error": errorMsg})
}

// 웹소켓 메시지 전송 (스레드 안전)
func sendMessage(conn *websocket.Conn, mutex *sync.Mutex, msgType string, payload interface{}) {
    msg := WebSocketMessage{
        Type:    msgType,
        Payload: payload,
    }
    
    if mutex != nil {
        mutex.Lock()
        defer mutex.Unlock()
    }
    
    if err := conn.WriteJSON(msg); err != nil {
        logMessage("웹소켓 메시지 전송 실패: %v", err)
    }
}