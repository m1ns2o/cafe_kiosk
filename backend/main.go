package main

import (
	"kiosk/database"
	"kiosk/routes"
	"kiosk/utils"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// DepositState 예수금 상태를 관리하는 구조체
type DepositState struct {
    mu             sync.RWMutex
    currentDeposit int64
    lastUpdateTime time.Time
    kisClient      *utils.KISApi
}

// NewDepositState 새로운 DepositState 인스턴스 생성
func NewDepositState(kisClient *utils.KISApi) *DepositState {
    return &DepositState{
        kisClient: kisClient,
    }
}

// Initialize 초기 예수금 값 설정
func (ds *DepositState) Initialize() error {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    depositAmount, err := ds.kisClient.GetDepositAmount()
    if err != nil {
        return err
    }

    ds.currentDeposit = depositAmount
    ds.lastUpdateTime = time.Now()
    log.Printf("초기 예수금 설정 완료: %s원", utils.FormatNumber(depositAmount))
    return nil
}

// GetCurrentDeposit 현재 예수금 조회 (읽기 전용)
func (ds *DepositState) GetCurrentDeposit() int64 {
    ds.mu.RLock()
    defer ds.mu.RUnlock()
    return ds.currentDeposit
}

// UpdateAndCheckDeposit 예수금 업데이트 및 변동 확인
func (ds *DepositState) UpdateAndCheckDeposit(expectedAmount int64) (bool, int64, error) {
    ds.mu.Lock()
    defer ds.mu.Unlock()

    // 최신 예수금 조회
    newDepositAmount, err := ds.kisClient.GetDepositAmount()
    if err != nil {
        return false, 0, err
    }

    // 실제 변동액 계산
    actualChange := newDepositAmount - ds.currentDeposit

    // 상태 업데이트
    ds.currentDeposit = newDepositAmount
    ds.lastUpdateTime = time.Now()

    // 예상 변동액과 실제 변동액 비교
    return actualChange == expectedAmount, actualChange, nil
}

// 전역 상태 관리 변수
var depositState *DepositState

func main() {
    // .env 파일 로드
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found")
    }

    // DB 연결
    if err := database.InitDB(); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // 기본 카테고리 생성
    database.InitializeCategories()

    // KIS API 토큰 설정 (환경 변수에서 값 로드)
    appKey := os.Getenv("KIS_APP_KEY")
    appSecret := os.Getenv("KIS_APP_SECRET")
    accountNo := os.Getenv("KIS_ACCOUNT_NO")
    accountProdCode := os.Getenv("KIS_ACCOUNT_PROD_CODE")

    if appKey == "" || appSecret == "" || accountNo == "" {
        log.Fatal("KIS API credentials are not set in environment variables")
    }

    // KIS API 클라이언트 생성
    kisApi := utils.NewKISApi(appKey, appSecret, accountNo, accountProdCode)
    
    // 토큰 발급
    success, err := kisApi.GetAccessToken()
    if err != nil {
        log.Fatalf("Failed to get KIS API token: %v", err)
    }
    
    if !success {
        log.Fatal("Failed to get KIS API token")
    }
    
    log.Println("KIS API token obtained successfully")

    // 예수금 상태 관리 초기화
    depositState = NewDepositState(kisApi)
    if err := depositState.Initialize(); err != nil {
        log.Fatalf("예수금 상태 초기화 실패: %v", err)
    }

    // Gin 라우터 설정
    r := gin.Default()

    // KIS API 클라이언트와 DepositState를 라우터 컨텍스트에 저장
    r.Use(func(c *gin.Context) {
        c.Set("kisApi", kisApi)
        c.Set("depositState", depositState)
        c.Next()
    })

    // API 라우트 설정
    routes.SetupRoutes(r)

    // 정적 파일 서빙 (Vue 등)
    // 정적 파일 서빙 설정
    // 정적 파일 서빙 설정
r.StaticFS("/assets", http.Dir("./static/dist/assets"))
r.StaticFile("/vite.svg", "./static/dist/vite.svg")
r.StaticFile("/test.html", "./static/test.html")

// uploads 폴더는 별도로 정의 - 라우터 그룹으로 처리
uploads := r.Group("/uploads")
{
    uploads.StaticFS("/", http.Dir("./uploads"))
}

// 메인 페이지
r.GET("/", func(c *gin.Context) {
    c.File("./static/dist/index.html")
})

// SPA를 위한 fallback - 명시적으로 uploads 경로 제외
r.NoRoute(func(c *gin.Context) {
    path := c.Request.URL.Path
    
    // uploads 경로에 대한 요청은 처리하지 않음 (404 반환)
    if strings.HasPrefix(path, "/uploads/") {
        c.AbortWithStatus(http.StatusNotFound)
        return
    }
    
    // API 및 정적 파일 경로가 아니면 index.html 제공
    if !strings.HasPrefix(path, "/api") && 
       !strings.HasPrefix(path, "/assets") && 
       path != "/test" {
        c.File("./static/dist/index.html")
    }
})
    // 서버 시작
    r.Run(":8080")
}