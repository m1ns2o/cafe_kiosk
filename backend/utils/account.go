// account.go
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// KISApi represents the KIS API client
type KISApi struct {
	AppKey          string
	AppSecret       string
	AccountNo       string
	AccountProdCode string
	BaseURL         string
	AccessToken     string
	TokenExpireTime time.Time
	ApprovalKey     string
}

// TokenRequest represents the token request payload
type TokenRequest struct {
	GrantType string `json:"grant_type"`
	AppKey    string `json:"appkey"`
	AppSecret string `json:"appsecret"`
}

// TokenResponse represents the token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

// BalanceResponse represents the balance inquiry response
type BalanceResponse struct {
	RtCd    string           `json:"rt_cd"`
	Msg1    string           `json:"msg1"`
	Output1 []BalanceItem    `json:"output1"` // 개별 종목 정보 (현재 사용하지 않지만 API 응답 구조상 포함)
	Output2 []BalanceSummary `json:"output2"`
}

// BalanceItem represents individual stock balance information
type BalanceItem struct {
	// 현재 사용하지 않으므로 비워둠
}

// BalanceSummary represents account balance summary
type BalanceSummary struct {
	DncaTotAmt string `json:"dnca_tot_amt"` // 예수금 총액
}

type DepositState struct {
	mu             sync.RWMutex
	currentDeposit int64
	lastUpdateTime time.Time
	kisClient      *KISApi
}

// NewKISApi creates a new KIS API client
func NewKISApi(appKey, appSecret, accountNo string, accountProdCode ...string) *KISApi {
	prodCode := "01"
	if len(accountProdCode) > 0 {
		prodCode = accountProdCode[0]
	}
	return &KISApi{
		AppKey:          appKey,
		AppSecret:       appSecret,
		AccountNo:       accountNo,
		AccountProdCode: prodCode,
		BaseURL:         "https://openapi.koreainvestment.com:9443", // 실전투자 URL
	}
}

// GetAccessToken obtains an OAuth token
func (k *KISApi) GetAccessToken() (bool, error) {
	url := fmt.Sprintf("%s/oauth2/tokenP", k.BaseURL)
	tokenRequest := TokenRequest{
		GrantType: "client_credentials",
		AppKey:    k.AppKey,
		AppSecret: k.AppSecret,
	}

	jsonData, err := json.Marshal(tokenRequest)
	if err != nil {
		return false, fmt.Errorf("토큰 요청 데이터 마샬링 실패: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("토큰 요청 생성 실패: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("토큰 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var tokenResp TokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			return false, fmt.Errorf("토큰 응답 디코딩 실패: %v", err)
		}
		k.AccessToken = tokenResp.AccessToken
		// 토큰 만료 시간 설정 (현재 시간 + 23시간)
		k.TokenExpireTime = time.Now().Add(23 * time.Hour)
		fmt.Printf("토큰 발급 성공: %s...\n", k.AccessToken[:10])
		return true, nil
	}

	return false, fmt.Errorf("토큰 발급 실패: %d", resp.StatusCode)
}

// GetBalance retrieves stock balance information
func (k *KISApi) GetBalance() (*BalanceResponse, error) {
	// Check token validity
	if k.AccessToken == "" || time.Now().After(k.TokenExpireTime) {
		success, err := k.GetAccessToken()
		if !success || err != nil {
			return nil, fmt.Errorf("토큰 갱신 실패: %v", err)
		}
	}

	baseURL := fmt.Sprintf("%s/uapi/domestic-stock/v1/trading/inquire-balance", k.BaseURL)
	
	// Create query parameters
	params := url.Values{}
	params.Add("CANO", k.AccountNo)
	params.Add("ACNT_PRDT_CD", k.AccountProdCode)
	params.Add("AFHR_FLPR_YN", "N")
	params.Add("OFL_YN", "")
	params.Add("INQR_DVSN", "02")
	params.Add("UNPR_DVSN", "01")
	params.Add("FUND_STTL_ICLD_YN", "N")
	params.Add("FNCG_AMT_AUTO_RDPT_YN", "N")
	params.Add("PRCS_DVSN", "00")
	params.Add("CTX_AREA_FK100", "")
	params.Add("CTX_AREA_NK100", "")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("잔고 조회 요청 생성 실패: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", k.AccessToken))
	req.Header.Set("appkey", k.AppKey)
	req.Header.Set("appsecret", k.AppSecret)
	req.Header.Set("tr_id", "TTTC8434R")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("잔고 조회 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		var balanceResp BalanceResponse
		if err := json.Unmarshal(body, &balanceResp); err != nil {
			return nil, fmt.Errorf("잔고 응답 디코딩 실패: %v", err)
		}

		if balanceResp.RtCd == "0" {
			return &balanceResp, nil
		}

		return nil, fmt.Errorf("잔고 조회 실패: %s", balanceResp.Msg1)
	}

	return nil, fmt.Errorf("잔고 조회 실패: %d", resp.StatusCode)
}

// GetDepositAmount retrieves only the deposit amount (예수금 총액)
func (k *KISApi) GetDepositAmount() (int64, error) {
	balanceResp, err := k.GetBalance()
	if err != nil {
		return 0, err
	}

	if len(balanceResp.Output2) > 0 {
		dncaTot, err := strconv.ParseInt(balanceResp.Output2[0].DncaTotAmt, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("예수금 총액 파싱 실패: %v", err)
		}
		return dncaTot, nil
	}

	return 0, fmt.Errorf("예수금 정보가 없습니다")
}

// FormatNumber formats a number with commas
func FormatNumber(n int64) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	str := strconv.FormatInt(n, 10)
	result := ""
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	return sign + result
}

// 여기서부터 main.go에서 옮긴 코드 시작

// DepositState 예수금 상태를 관리하는 구조체


// NewDepositState 새로운 DepositState 인스턴스 생성
func NewDepositState(kisClient *KISApi) *DepositState {
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
	log.Printf("초기 예수금 설정 완료: %s원", FormatNumber(depositAmount))
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