package main

import (
	"kiosk/database"
	"kiosk/routes"
	"kiosk/utils"
	"kiosk/handlers"
	"log"
	// "net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/static"
	"github.com/joho/godotenv"
)

var depositState *utils.DepositState

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
	
	// PWA 관련 파일 명시적 서빙
	// r.StaticFile("/manifest.webmanifest", "./static/dist/manifest.webmanifest")
	// r.StaticFile("/manifest.json", "./static/dist/manifest.json") // 두 가지 이름으로 제공
	// r.StaticFile("/sw.js", "./static/dist/sw.js") // 서비스 워커
	// r.StaticFile("/registerSW.js", "./static/dist/registerSW.js") // Vite PWA에서 생성하는 등록 스크립트
	// r.StaticFile("/favicon.ico", "./static/dist/favicon.ico")
	
	// // 정적 파일들 서빙
	// r.Use(static.Serve("/", static.LocalFile("./static/dist", false)))
	// r.StaticFS("/uploads", http.Dir("./uploads"))
	// r.StaticFS("/assets", http.Dir("./static/dist/assets"))
	
	// // SPA Fallback 핸들러 (존재하지 않는 경로 처리)
	// r.NoRoute(func(c *gin.Context) {
	// 	path := c.Request.URL.Path
		
	// 	// 정적 파일 확장자가 아닌 경우에만 index.html 제공
	// 	if !isStaticFileRequest(path) {
	// 		c.File("./static/dist/index.html")
	// 		return
	// 	}
	// 	c.Status(404)
	// })

	// 서버 시작
	log.Println("서버가 http://localhost:8080 에서 시작되었습니다.")
	r.Run(":8080")
}

// 정적 파일 요청 여부 확인
// func isStaticFileRequest(path string) bool {
// 	staticExtensions := []string{".js", ".css", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webmanifest", ".json", ".ico", ".woff", ".woff2", ".ttf"}
// 	for _, ext := range staticExtensions {
// 		if strings.HasSuffix(path, ext) {
// 			return true
// 		}
// 	}
// 	return false
// }