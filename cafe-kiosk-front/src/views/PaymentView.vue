<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { PaymentAPI } from '../api/payment';
import QRCode from 'qrcode';

const route = useRoute();
// const router = useRouter();

// 결제 상태
const paymentStatus = ref<'pending' | 'success' | 'failed'>('pending');
const statusMessage = ref<string>('결제를 진행 중입니다...');
const paymentDetails = ref<any>(null);

// QR 코드 관련
const qrCodeDataUrl = ref<string>('');
const totalAmount = ref<number>(0);

// 결제 요청 및 결과 확인
const requestPayment = async () => {
  try {
    const response = await PaymentAPI.requestPayment(totalAmount.value);
    
    if (response.data.success) {
      paymentStatus.value = 'success';
      statusMessage.value = '결제가 완료되었습니다!';
    } else {
      paymentStatus.value = 'failed';
      statusMessage.value = response.data.message || '결제에 실패했습니다.';
      paymentDetails.value = response.data.details;
    }
  } catch (error) {
    console.error('결제 요청 중 오류 발생:', error);
    paymentStatus.value = 'failed';
    statusMessage.value = '결제 처리 중 오류가 발생했습니다.';
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
    // 장바구니 아이템 목록 출력
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
  
  // 결제 요청 시작
  requestPayment();
});
</script>

<template>
  <div class="payment-view-container">
    <div class="payment-view">
      <div class="payment-info">
        <div class="amount-display">
          <h2>결제 금액</h2>
          <div class="total-amount">{{ totalAmount.toLocaleString() }}원</div>
        </div>
        
        <div class="qr-code-container">
          <div class="qr-code">
            <img v-if="qrCodeDataUrl" :src="qrCodeDataUrl" alt="결제 QR 코드" />
            <div v-else class="loading-qr">QR 코드 생성 중...</div>
          </div>
          <p class="qr-instruction">QR 코드를 스캔하여 결제를 진행해주세요.</p>
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
            <span v-if="paymentStatus === 'pending'" class="material-icon">hourglass_empty</span>
            <span v-else-if="paymentStatus === 'success'" class="material-icon">check_circle</span>
            <span v-else class="material-icon">error</span>
          </div>
          <div class="status-message">{{ statusMessage }}</div>
          <div v-if="paymentDetails && paymentStatus === 'failed'" class="payment-details">
            <p>예상 금액: {{ paymentDetails.expected_amount.toLocaleString() }}원</p>
            <p>실제 변동액: {{ paymentDetails.actual_change.toLocaleString() }}원</p>
            <p>타임아웃: {{ paymentDetails.timeout_after }}</p>
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
  min-height: 100vh;
  background-color: var(--background-primary, #f5f5f5);
  padding: 20px;
}

.payment-view {
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.payment-title {
  font-size: 1.8rem;
  color: var(--text-primary, #333);
  margin-bottom: 20px;
  text-align: center;
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.payment-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
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
}

.qr-code-container h2 {
  font-size: 1.2rem;
  margin-bottom: 15px;
  color: var(--text-secondary, #555);
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
  max-width: 300px;
}

/* 장바구니 요약 스타일 */
.cart-summary {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 15px;
  max-height: 25vh;
  overflow-y: auto;
}

.cart-summary::-webkit-scrollbar {
  display: none;
}

.cart-summary h2 {
  font-size: 1.2rem;
  margin-bottom: 10px;
  color: var(--text-secondary, #555);
}

.cart-items-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cart-item-summary {
  display: flex;
  justify-content: space-between;
  padding: 8px;
  background-color: white;
  border-radius: 4px;
  border: 1px solid #eee;
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
  color: var(--button-primary, #4caf50);
  font-weight: 500;
}

.payment-status {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 15px;
  border-radius: 8px;
  background-color: #f9f9f9;
  margin-top: 10px;
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

.status-icon {
  font-size: 2rem;
  margin-bottom: 10px;
}

.status-icon .material-icon {
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

.status-message {
  font-size: 1.2rem;
  font-weight: 500;
  text-align: center;
  margin-bottom: 10px;
}

.payment-details {
  font-size: 0.9rem;
  color: #666;
  text-align: center;
  margin-top: 10px;
}

.payment-actions {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-top: 30px;
}

.cancel-btn, .home-btn {
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
}

.cancel-btn {
  background-color: #f1f1f1;
  color: #555;
}

.home-btn {
  background-color: var(--button-primary, #4caf50);
  color: white;
}

.cancel-btn:hover {
  background-color: #e5e5e5;
}

.home-btn:hover {
  opacity: 0.9;
}

.mr-1 {
  margin-right: 5px;
}

@media (max-width: 768px) {
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
  
  .payment-actions {
    flex-direction: column;
  }
  
  .cancel-btn, .home-btn {
    width: 100%;
  }
}
</style>
