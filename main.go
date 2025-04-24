package main

import (
    "log"
    "strings"
    "github.com/gin-gonic/gin"
    "kiosk/database"
    "kiosk/routes"
)

func main() {
    // DB 연결
    if err := database.InitDB(); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // 기본 카테고리 생성
    database.InitializeCategories()

    // Gin 라우터 설정
    r := gin.Default()

    // API 라우트 설정
    routes.SetupRoutes(r)

    // 정적 파일 서빙 (Vue 등)
    r.Static("/static", "./static")
    r.StaticFile("/test", "./static/test.html")  // 테스트 페이지 직접 접근
    
    // Vue 앱 관련 설정
    r.StaticFile("/", "./dist/index.html")
    r.NoRoute(func(c *gin.Context) {
        // API 경로가 아닌 경우에만 index.html 반환
        if !strings.HasPrefix(c.Request.URL.Path, "/api") {
            c.File("./dist/index.html")
        }
    })

    // 서버 시작
    r.Run(":8080")
}