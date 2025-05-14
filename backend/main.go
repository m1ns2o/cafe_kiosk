package main

import (
	"kiosk/database"
	"kiosk/routes"
	"kiosk/utils"
	"kiosk/handlers"
	"log"
	"time"
	"sync"
	// "net/http"
	"os"
	"strings"
	"context"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/static"
	"github.com/joho/godotenv"
)

var depositState *utils.DepositState

// refreshToken handles the periodic token refresh process
func refreshTokenPeriodically(ctx context.Context, kisApi *utils.KISApi, wg *sync.WaitGroup) {
	defer wg.Done()
	
	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			log.Println("토큰 재발급 시도 중...")
			success, err := kisApi.GetAccessToken()
			if !success || err != nil {
				log.Printf("토큰 재발급 실패: %v", err)
				// 실패해도 계속 진행
				continue
			}
			log.Println("KIS API 토큰 재발급 성공")
			
		case <-ctx.Done():
			log.Println("토큰 재발급 고루틴 종료")
			return
		}
	}
}

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

	// KIS API 설정
	appKey := os.Getenv("KIS_APP_KEY")
	appSecret := os.Getenv("KIS_APP_SECRET")
	accountNo := os.Getenv("KIS_ACCOUNT_NO")
	accountProdCode := os.Getenv("KIS_ACCOUNT_PROD_CODE")

	if appKey == "" || appSecret == "" || accountNo == "" {
		log.Fatal("KIS API credentials are not set in environment variables")
	}

	kisApi := utils.NewKISApi(appKey, appSecret, accountNo, accountProdCode)
	
	// 토큰 발급
	if success, err := kisApi.GetAccessToken(); !success || err != nil {
		log.Fatalf("Failed to get KIS API token: %v", err)
	}
	log.Println("KIS API token obtained successfully")
	
	// 토큰 주기적 재발급을 위한 컨텍스트 및 WaitGroup 설정
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	
	// 토큰 재발급 고루틴 시작
	wg.Add(1)
	go refreshTokenPeriodically(ctx, kisApi, &wg)

	// 예수금 상태 관리 초기화
	depositState = utils.NewDepositState(kisApi)
	if err := depositState.Initialize(); err != nil {
		log.Fatalf("예수금 상태 초기화 실패: %v", err)
	}
	log.Printf("현재 예수금: %s원", utils.FormatNumber(depositState.GetCurrentDeposit()))

    // 로그 시스템 초기화
    if err := handlers.InitLogSystem(); err != nil {
        log.Fatalf("로그 시스템 초기화 실패: %v", err)
    }
    defer handlers.CloseLogSystem()
    
    // SSE 브로드캐스터 시작
    handlers.StartSSEBroadcaster()

	// Gin 라우터 설정
	r := gin.Default()
	
	// 프록시 신뢰 설정 (Gin 경고 메시지 해결)
	r.SetTrustedProxies(nil) // 모든 프록시를 신뢰하지 않음 (개발 환경용)
	
	// KIS API 미들웨어 설정
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Set("kisApi", kisApi)
			c.Set("depositState", depositState)
		}
		c.Next()
	})
	
	// API 라우트 설정
	routes.SetupRoutes(r)
	
	// 서버 시작
	log.Println("서버가 http://localhost:8080 에서 시작되었습니다.")
	
	// 종료 처리를 위한 defer 설정
	defer func() {
		cancel() // 토큰 재발급 고루틴 종료 신호
		wg.Wait() // 고루틴이 종료될 때까지 대기
		log.Println("모든 고루틴이 정상적으로 종료되었습니다.")
	}()
	
	r.Run(":8080")
}