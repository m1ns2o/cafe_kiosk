<script setup lang="ts">
import { ref, onMounted, computed, onBeforeUnmount } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import QRCode from 'qrcode';
import { PaymentAPI } from '../api/payment';

const route = useRoute();
const router = useRouter();

// 결제 상태
const paymentStatus = ref<'pending' | 'success' | 'failed' | 'cancelled'>('pending');
const statusMessage = ref<string>('결제를 진행 중입니다...');
const progressInfo = ref<string>('');
const paymentAttempt = ref<number>(0);
const maxAttempts = ref<number>(0);
const paymentID = ref<string>('');

// QR 코드 관련
const qrCodeDataUrl = ref<string>('');
const totalAmount = ref<number>(0);
const isProcessingCancel = ref<boolean>(false);

// 웹소켓 연결
const socket = ref<WebSocket | null>(null);
let redirectTimer: ReturnType<typeof setTimeout>;

// 웹소켓 연결 설정
const setupWebSocket = () => {
  // 웹소켓 서버 URL (실제 환경에 맞게 수정해야 함)
  const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/api/ws/payment';
  socket.value = new WebSocket(wsUrl);

  // 웹소켓 이벤트 핸들러 등록
  socket.value.onopen = () => {
    console.log('웹소켓 연결 성공');
    // 연결 후 결제 요청 전송
    sendPaymentRequest();
  };

  socket.value.onmessage = (event) => {
    handleWebSocketMessage(event);
  };

  socket.value.onclose = (event) => {
    console.log('웹소켓 연결 종료:', event);
  };

  socket.value.onerror = (error) => {
    console.error('웹소켓 오류:', error);
    paymentStatus.value = 'failed';
    statusMessage.value = '결제 서버와 연결 중 오류가 발생했습니다.';
    
    // 5초 후 주문 화면으로 이동
    redirectTimer = setTimeout(() => {
      router.push({ name: 'OrderView' });
    }, 5000);
  };
};

// 웹소켓 메시지 처리


// 결제 결과 처리
const handlePaymentResult = (result: any) => {
  if (result.success) {
    paymentStatus.value = 'success';
    statusMessage.value = '결제가 완료되었습니다!';
    
    // 결제 성공 시 주문 데이터를 백엔드로 전송
    submitOrderToBackend().then(() => {
      // 성공 페이지로 이동 (지연 추가)
      setTimeout(() => {
        router.push({ name: 'PaymentSuccessView' });
      }, 1000);
    });
  } else {
    paymentStatus.value = 'failed';
    statusMessage.value = result.message || '결제에 실패했습니다.';
    
    redirectTimer = setTimeout(() => {
      router.push({ name: 'OrderView' });
    }, 5000);
  }
};

// 결제 요청 전송 (웹소켓)
const sendPaymentRequest = () => {
  if (!socket.value || socket.value.readyState !== WebSocket.OPEN) {
    console.error('웹소켓이 연결되지 않았습니다.');
    return;
  }
  
  const paymentRequest = {
    type: 'payment_request',
    payload: {
      amount: totalAmount.value,
      order_id: generateOrderId(), // 주문 ID 생성
      timestamp: new Date().toISOString()
    }
  };
  
  socket.value.send(JSON.stringify(paymentRequest));
  console.log('결제 요청 전송 완료');
};

// 결제 취소 요청 전송 (웹소켓)
// 웹소켓 메시지 처리 함수 전체
const handleWebSocketMessage = (event: MessageEvent) => {
  try {
    const message = JSON.parse(event.data);
    
    switch (message.type) {
      case 'payment_initiated':
        // 결제 ID 저장
        paymentID.value = message.payload.payment_id;
        console.log('결제 ID 수신:', paymentID.value);
        break;
        
      case 'payment_status':
        // 결제 진행 상태 업데이트
        paymentAttempt.value = message.payload.attempt;
        maxAttempts.value = message.payload.max_attempts;
        progressInfo.value = `결제 확인 중... (${paymentAttempt.value}/${maxAttempts.value})`;
        break;
        
      case 'payment_result':
        // 최종 결제 결과 처리
        handlePaymentResult(message.payload);
        break;
        
      case 'error':
        // 오류 처리
        paymentStatus.value = 'failed';
        statusMessage.value = message.payload.error || '결제 처리 중 오류가 발생했습니다.';
        
        redirectTimer = setTimeout(() => {
          router.push({ name: 'OrderView' });
        }, 5000);
        break;
        
      case 'cancel_result':
        // 취소 결과 처리
        if (message.payload.success) {
          paymentStatus.value = 'cancelled';
          statusMessage.value = '결제가 취소되었습니다.';
          isProcessingCancel.value = false;
          
          // 3초 후 주문 화면으로 이동
          redirectTimer = setTimeout(() => {
            router.push({ name: 'OrderView' });
          }, 3000);
        } else {
          isProcessingCancel.value = false;
          statusMessage.value = message.payload.message || '결제 취소에 실패했습니다.';
        }
        break;
    }
  } catch (error) {
    console.error('웹소켓 메시지 파싱 오류:', error);
  }
};

// 결제 취소 요청 함수 전체
const cancelPayment = () => {
  if (!socket.value || socket.value.readyState !== WebSocket.OPEN) {
    console.error('웹소켓이 연결되지 않았습니다.');
    return;
  }
  
  isProcessingCancel.value = true;
  statusMessage.value = '결제를 취소 중입니다...';
  
  const cancelRequest = {
    type: 'cancel_request',
    payload: {
      payment_id: paymentID.value,
      amount: totalAmount.value,
      timestamp: new Date().toISOString()
    }
  };
  
  socket.value.send(JSON.stringify(cancelRequest));
  console.log('결제 취소 요청 전송 완료 - 결제 ID:', paymentID.value);
};

// 결제 재시도
const retryPayment = () => {
  clearTimeout(redirectTimer);
  paymentStatus.value = 'pending';
  statusMessage.value = '결제를 다시 시도 중입니다...';
  progressInfo.value = '';
  
  // 웹소켓이 열려있는지 확인
  if (socket.value && socket.value.readyState === WebSocket.OPEN) {
    sendPaymentRequest();
  } else {
    // 웹소켓이 닫혀있다면 다시 연결
    setupWebSocket();
  }
};

// 주문 ID 생성 헬퍼 함수
const generateOrderId = () => {
  return 'ORD-' + Date.now() + '-' + Math.floor(Math.random() * 1000);
};

// 주문 데이터를 백엔드로 전송
const submitOrderToBackend = async () => {
  try {
    // 백엔드로 주문 데이터 전송
    await PaymentAPI.postOrder(cartItems.value);
    console.log('주문 데이터가 성공적으로 전송되었습니다.');
    return true;
  } catch (error) {
    console.error('주문 데이터 전송 중 오류 발생:', error);
    return false;
  }
};

// QR 코드 생성
const generateQRCode = async () => {
  try {
    // 토스 결제 URL 생성 (amount 값을 총액으로 교체)
    const tossPaymentUrl = `supertoss://send?amount=${totalAmount.value}&bank=%ED%95%9C%EA%B5%AD%ED%88%AC%EC%9E%90%EC%A6%9D%EA%B6%8C&accountNo=6961147001&origin=qr`;
    
    // QR 코드 생성 (qrcode 패키지 사용)
    qrCodeDataUrl.value = await QRCode.toDataURL(tossPaymentUrl, {
      width: 300,
      margin: 2,
      color: {
        dark: '#000000',
        light: '#FFFFFF'
      }
    });
  } catch (error) {
    console.error('QR 코드 생성 중 오류 발생:', error);
    // 오류 발생 시 대체 이미지 생성
    const canvas = document.createElement('canvas');
    canvas.width = 300;
    canvas.height = 300;
    const ctx = canvas.getContext('2d');
    if (ctx) {
      ctx.fillStyle = '#FFFFFF';
      ctx.fillRect(0, 0, 300, 300);
      ctx.fillStyle = '#000000';
      ctx.font = '16px Arial';
      ctx.textAlign = 'center';
      ctx.fillText('결제 금액: ' + totalAmount.value.toLocaleString() + '원', 150, 150);
      qrCodeDataUrl.value = canvas.toDataURL();
    }
  }
};

// 주문 정보
const cartItems = computed(() => {
  if (!route.params.cartItems) return [];
  try {
    const decodedData = decodeURIComponent(route.params.cartItems as string);
    const parsedData = JSON.parse(decodedData);
    console.log('장바구니 아이템:', parsedData);
    return parsedData;
  } catch (error) {
    console.error('장바구니 데이터 파싱 오류:', error);
    return [];
  }
});

onMounted(async () => {
  // 라우터에서 전달받은 총액 설정
  if (route.params.totalAmount) {
    totalAmount.value = parseInt(route.params.totalAmount as string);
  }
  
  // QR 코드 생성
  await generateQRCode();
  
  // 웹소켓 연결 설정
  setupWebSocket();
});

onBeforeUnmount(() => {
  // 타이머 정리
  if (redirectTimer) {
    clearTimeout(redirectTimer);
  }
  
  // 웹소켓 연결 종료
  if (socket.value) {
    socket.value.close();
    socket.value = null;
  }
});
</script>

<template>
  <div class="payment-view-container">
    <div class="payment-view">
      <!-- 가로 레이아웃을 위한 구조 변경 -->
      <div class="payment-layout">
        <!-- 좌측 영역: QR 코드 -->
        <div class="payment-left">
          <div class="qr-code-container">
            <div class="qr-code">
              <img v-if="qrCodeDataUrl" :src="qrCodeDataUrl" alt="결제 QR 코드" />
              <div v-else class="loading-qr">QR 코드 생성 중...</div>
            </div>
            <p class="qr-instruction">QR 코드를 스캔하여 결제를 진행해주세요.</p>
            <p class="account-instruction">토스 앱이 없다면 한국투자증권 6961147001로 입금해주세요</p>
          </div>
        </div>
        
        <!-- 우측 영역: 결제 정보 및 상태 -->
        <div class="payment-right">
          <div class="amount-display">
            <h2>결제 금액</h2>
            <div class="total-amount">{{ totalAmount.toLocaleString() }}원</div>
          </div>
          
          <!-- 장바구니 아이템 표시 -->
          <div v-if="cartItems.length > 0" class="cart-summary">
            <h2>주문 내역</h2>
            <div class="cart-items-list">
              <div v-for="(cartItem, index) in cartItems" :key="index" class="cart-item-summary">
                <div class="item-name">{{ cartItem.item.name }}</div>
                <div class="item-quantity">{{ cartItem.quantity }}개</div>
                <div class="item-price">{{ (cartItem.item.price * cartItem.quantity).toLocaleString() }}원</div>
              </div>
            </div>
          </div>
          
          <div class="payment-status" :class="paymentStatus">
            <div class="status-icon">
              <span v-if="paymentStatus === 'pending'" class="material-icons">hourglass_empty</span>
              <span v-else-if="paymentStatus === 'success'" class="material-icons">check_circle</span>
              <span v-else-if="paymentStatus === 'cancelled'" class="material-icons">cancel</span>
              <span v-else class="material-icons">error</span>
            </div>
            <div class="status-message">{{ statusMessage }}</div>
            
            <!-- 진행 상태 표시 -->
            <div v-if="paymentStatus === 'pending' && progressInfo" class="progress-info">
              {{ progressInfo }}
            </div>
            
            <!-- 결제 실패 시 재시도 버튼 표시 -->
            <button v-if="paymentStatus === 'failed'" class="retry-btn" @click="retryPayment">
              <span class="material-icons mr-1">refresh</span>
              결제 재시도
            </button>
            
            <!-- 결제 취소 버튼 -->
            <button v-if="paymentStatus === 'pending' && !isProcessingCancel" class="cancel-btn" @click="cancelPayment">
              <span class="material-icons mr-1">close</span>
              결제 취소
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.payment-view-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--background-primary, #f5f5f5);
  /* padding: 100px; */
  padding-bottom: 80px;
  padding-top: 80px;
  overflow: hidden;
}

.payment-view {
  max-width: 80%; /* 가로 레이아웃을 위해 최대 너비 증가 */
  margin: 0 auto;
  width: 100%;
  height: 100%;
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  overflow: hidden;
}

/* 가로 레이아웃을 위한 새로운 컨테이너 스타일 */
.payment-layout {
  display: flex;
  flex-direction: row; /* 가로 방향으로 배치 */
  height: 100%;
  gap: 20px;
}

/* 좌측 QR 코드 영역 */
.payment-left {
  flex: 0 0 40%; /* 고정 너비 (전체의 40%) */
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

/* 우측 결제 정보 영역 */
.payment-right {
  flex: 1; /* 남은 공간 모두 차지 */
  display: flex;
  flex-direction: column;
  gap: 20px;
  overflow-y: auto;
}

.payment-title {
  font-size: 1.8rem;
  color: var(--text-primary, #333);
  margin-bottom: 20px;
  text-align: center;
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.amount-display {
  text-align: center;
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 8px;
}

.amount-display h2 {
  font-size: 1.2rem;
  margin-bottom: 10px;
  color: var(--text-secondary, #555);
}

.total-amount {
  font-size: 2rem;
  font-weight: bold;
  color: var(--button-primary, #4caf50);
}

.qr-code-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  background-color: white;
  border: 1px solid #eee;
  border-radius: 8px;
  width: 100%;
}

.qr-code {
  width: 300px;
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 15px;
}

.qr-code img {
  max-width: 100%;
  max-height: 100%;
}

.loading-qr {
  font-size: 1rem;
  color: #888;
}

.qr-instruction {
  text-align: center;
  color: #666;
  font-size: 0.9rem;
}

.account-instruction {
  margin-top: 3px;
  text-align: center;
  color: #999;
  font-size: 0.9rem;
}

/* 장바구니 요약 스타일 */
.cart-summary {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 15px;
  /* max-height: 40vh; */
  flex:1;
  overflow-y: auto;
}

.cart-summary::-webkit-scrollbar {
  display: none;
}

.cart-summary h2 {
  font-size: 1.2rem;
  margin-bottom: 10px;
  color: var(--text-secondary);
}

.cart-items-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow-y: auto;
}

.cart-item-summary {
  display: flex;
  justify-content: space-between;
  padding: 8px;
  background-color: white;
  border-radius: 4px;
  border: 1px solid #eee;
  overflow: hidden;
}

.item-name {
  flex: 2;
  font-weight: 500;
}

.item-quantity {
  flex: 1;
  text-align: center;
}

.item-price {
  flex: 1;
  text-align: right;
  color: var(--button-primary);
  font-weight: 500;
}

.payment-status {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 15px;
  border-radius: 8px;
  background-color: #f9f9f9;
}

.payment-status.pending {
  background-color: #fff9e6;
}

.payment-status.success {
  background-color: #e6f7e6;
}

.payment-status.failed {
  background-color: #ffe6e6;
}

.payment-status.cancelled {
  background-color: #f0f0f0;
}

.status-icon {
  font-size: 2rem;
  margin-bottom: 10px;
}

.status-icon .material-icons {
  font-family: 'Material Icons';
  font-weight: normal;
  font-style: normal;
  font-size: 36px;
  line-height: 1;
  letter-spacing: normal;
  text-transform: none;
  display: inline-block;
  white-space: nowrap;
  word-wrap: normal;
  direction: ltr;
  font-feature-settings: 'liga';
  -webkit-font-feature-settings: 'liga';
  -webkit-font-smoothing: antialiased;
}

.payment-status.pending .status-icon {
  color: #f9a825;
}

.payment-status.success .status-icon {
  color: #4caf50;
}

.payment-status.failed .status-icon {
  color: #f44336;
}

.payment-status.cancelled .status-icon {
  color: #757575;
}

.status-message {
  font-size: 1.2rem;
  font-weight: 500;
  text-align: center;
  margin-bottom: 10px;
}

.progress-info {
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 10px;
}

.payment-details {
  font-size: 0.9rem;
  color: #666;
  text-align: center;
  margin-top: 10px;
}

.cancel-btn, .home-btn, .retry-btn {
  padding: 10px 20px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: background-color 0.2s;
  margin-top: 10px;
}

.cancel-btn {
  background-color: #f44336;
  color: white;
}

.home-btn {
  background-color: var(--button-primary, #4caf50);
  color: white;
}

.retry-btn {
  background-color: #2196f3;
  color: white;
}

.cancel-btn:hover {
  background-color: #e53935;
}

.home-btn:hover, .retry-btn:hover {
  opacity: 0.9;
}

.mr-1 {
  margin-right: 5px;
}

/* 반응형 스타일링 - 화면 크기에 따라 레이아웃 변경 */
@media (max-width: 768px) {
  /* 태블릿 이하 크기에서는 세로 레이아웃으로 변경 */
  .payment-layout {
    flex-direction: column;
  }
  
  .payment-left, .payment-right {
    flex: none;
    width: 100%;
  }
  
  .payment-view {
    padding: 15px;
  }
  
  .qr-code {
    width: 250px;
    height: 250px;
  }
  
  .total-amount {
    font-size: 1.8rem;
  }
}

@media (max-width: 480px) {
  .payment-view-container {
    padding: 10px;
  }
  
  .payment-title {
    font-size: 1.5rem;
  }
  
  .qr-code {
    width: 200px;
    height: 200px;
  }
  
  .cancel-btn, .home-btn, .retry-btn {
    width: 100%;
  }
}

/* 높이 기반 미디어 쿼리 */
@media (max-height: 800px) {
  .qr-code {
    width: 220px;
    height: 220px;
  }
  
  .cart-summary {
    max-height: 18vh;
  }
  
  .payment-status {
    padding: 10px;
  }
  
  .status-icon .material-icons {
    font-size: 30px;
  }
  
  .status-message {
    font-size: 1rem;
  }
}

@media (max-height: 700px) {
  .qr-code {
    width: 180px;
    height: 180px;
  }
  
  .qr-instruction {
    font-size: 0.8rem;
    margin-bottom: 0;
  }
  
  .cart-summary {
    max-height: 15vh;
    padding: 10px;
  }
  
  .cart-summary h2 {
    font-size: 1rem;
    margin-bottom: 8px;
  }
  
  .amount-display {
    padding: 10px;
  }
  
  .amount-display h2 {
    font-size: 1rem;
    margin-bottom: 5px;
  }
  
  .total-amount {
    font-size: 1.5rem;
  }
  
  .payment-status {
    padding: 8px;
  }
  
  .status-icon .material-icons {
    font-size: 24px;
    margin-bottom: 5px;
  }
  
  .status-message {
    font-size: 0.9rem;
    margin-bottom: 5px;
  }
}

/* 극단적으로 작은 화면을 위한 미디어 쿼리 */
@media (max-width: 480px) and (max-height: 700px) {
  .qr-code {
    width: 150px;
    height: 150px;
  }
  
  .cart-summary {
    max-height: 12vh;
  }
  
  .payment-view-container {
    padding: 5px;
  }
  
  .payment-view {
    padding: 10px;
  }
  
  .amount-display, .payment-status {
    padding: 8px;
  }
}
</style>