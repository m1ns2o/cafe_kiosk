package handlers

import (
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"
    "github.com/gin-gonic/gin"
    "kiosk/database"
    "kiosk/models"
)

func GetMenus(c *gin.Context) {
    var menus []models.Menu
    query := database.DB
    
    // 카테고리별 필터링
    if categoryID := c.Query("category_id"); categoryID != "" {
        query = query.Where("category_id = ?", categoryID)
    }
    
    // 카테고리 이름으로 필터링 (옵션)
    if categoryName := c.Query("category"); categoryName != "" {
        query = query.Joins("JOIN categories ON categories.id = menus.category_id").
            Where("categories.name = ?", categoryName)
    }
    
    if err := query.Find(&menus).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, menus)
}

func GetMenusByCategory(c *gin.Context) {
    categoryID := c.Param("id")  // :category_id를 :id로 변경
    var menus []models.Menu
    
    if err := database.DB.Where("category_id = ?", categoryID).Find(&menus).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, menus)
}

func GetMenu(c *gin.Context) {
    id := c.Param("id")
    var menu models.Menu
    if err := database.DB.First(&menu, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
        return
    }
    c.JSON(http.StatusOK, menu)
}

func CreateMenu(c *gin.Context) {
    // 폼 데이터 파싱
    var req models.CreateMenuRequest
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 이미지 파일 처리
    file, err := c.FormFile("image")
    if err != nil && err != http.ErrMissingFile {
        c.JSON(http.StatusBadRequest, gin.H{"error": "이미지 업로드 오류: " + err.Error()})
        return
    }

    var imageURL string
    if file != nil {
        // 파일 확장자 가져오기
        ext := filepath.Ext(file.Filename)
        
        // 고유한 파일명 생성
        fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
        
        // 저장 경로 설정
        // filePath := filepath.Join("static", "uploads", fileName)
        filePath := filepath.Join("uploads", fileName)
        
        // static/uploads 디렉토리가 없으면 생성
        if err := os.MkdirAll(filepath.Join("uploads"), 0755); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "디렉토리 생성 오류: " + err.Error()})
            return
        }
        
        // 파일 저장
        if err := c.SaveUploadedFile(file, filePath); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "파일 저장 오류: " + err.Error()})
            return
        }
        
        // 웹에서 접근할 수 있는 URL 경로 설정
        imageURL = "/" + filePath
    }

    // 메뉴 객체 생성
    menu := models.Menu{
        CategoryID: req.CategoryID,
        Name:       req.Name,
        Price:      req.Price,
        ImageURL:   imageURL,
    }

    // 데이터베이스에 저장
    if err := database.DB.Create(&menu).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, menu)
}

func UpdateMenu(c *gin.Context) {
    id := c.Param("id")
    
    // 기존 메뉴 조회
    var menu models.Menu
    if err := database.DB.First(&menu, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
        return
    }

    // 폼 데이터 파싱
    var req models.CreateMenuRequest
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 이미지 파일 처리
    file, err := c.FormFile("image")
    if err != nil && err != http.ErrMissingFile {
        c.JSON(http.StatusBadRequest, gin.H{"error": "이미지 업로드 오류: " + err.Error()})
        return
    }

    if file != nil {
        // 기존 이미지가 있으면 삭제 (선택 사항)
        if menu.ImageURL != "" {
            oldFilePath := filepath.Join(".", menu.ImageURL)
            if err := os.Remove(oldFilePath); err != nil && !os.IsNotExist(err) {
                // 파일 삭제 실패는 치명적 오류가 아니므로 로그만 남기고 계속 진행
                log.Printf("기존 이미지 파일 삭제 실패: %v", err)
            }
        }

        // 새 파일명 생성
        ext := filepath.Ext(file.Filename)
        fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
        filePath := filepath.Join("uploads", fileName)
        
        // 파일 저장
        if err := c.SaveUploadedFile(file, filePath); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "파일 저장 오류: " + err.Error()})
            return
        }
        
        // 이미지 URL 업데이트
        menu.ImageURL = "/" + filePath
    }

    // 다른 필드 업데이트
    menu.CategoryID = req.CategoryID
    menu.Name = req.Name
    menu.Price = req.Price

    // 데이터베이스에 저장
    if err := database.DB.Save(&menu).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, menu)
}

func DeleteMenu(c *gin.Context) {
    id := c.Param("id")
    
    // 삭제할 메뉴 정보 조회
    var menu models.Menu
    if err := database.DB.First(&menu, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
        return
    }
    
    // 이미지 파일이 있으면 삭제
    if menu.ImageURL != "" {
        filePath := filepath.Join(".", menu.ImageURL)
        if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
            // 파일 삭제 실패는 로그만 남기고 계속 진행
            log.Printf("이미지 파일 삭제 실패: %v", err)
        }
    }

    // 데이터베이스에서 메뉴 삭제
    if err := database.DB.Delete(&models.Menu{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
}