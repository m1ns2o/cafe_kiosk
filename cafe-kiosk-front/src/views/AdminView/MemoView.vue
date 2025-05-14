<template>
  <div class="order-container">
    <div class="order-header">
      <h2>주문 현황 <el-badge :value="sortedOrders.length" type="primary" v-if="sortedOrders.length > 0"/></h2>
      <div class="actions">
        <el-button type="danger" plain size="small" @click="clearAllOrders" v-if="sortedOrders.length > 0">
          <el-icon><Delete /></el-icon> 모든 주문 지우기
        </el-button>
      </div>
    </div>
    
    <el-empty description="새로운 주문을 기다리는 중입니다..." v-if="sortedOrders.length === 0" />
    
    <div class="order-grid" :class="{ 'two-row-layout': isTwoRowLayout }" v-else>
      <el-row :gutter="20">
        <!-- 동적 카드 크기 조절 사용 -->
        <el-col 
          v-for="order in sortedOrders" 
          :key="order.id" 
          :span="getColSpan(sortedOrders.length)"
          class="order-col"
        >
          <el-card 
            :class="['order-card', { 'new-order': isNewOrder(order.id) }]" 
            shadow="hover"
            :style="isTwoRowLayout ? { height: 'calc((100vh - 110px) / 2)' } : {}"
          >
            <template #header>
              <div class="card-header">
                <div class="order-info">
                  <el-tag effect="dark" size="large" class="order-id">#{{ order.id }}</el-tag>
                  <el-tag type="info" size="small">{{ formatDate(order.created_at) }}</el-tag>
                </div>
                <el-button 
                  
                  :icon="Check"
                  circle
                  class="complete-btn"
                  @click="completeOrder(order.id)"
                />
              </div>
            </template>
            
            <div class="order-items">
              <el-scrollbar :height="isTwoRowLayout ? 'calc((100vh - 320px) / 2)' : '180px'">
                <div v-for="item in order.order_items" :key="`${order.id}-${item.menu.name}`" class="item">
                  <div class="item-top">
                    <span class="item-name">{{ item.menu.name }}</span>
                    <el-badge :value="item.quantity" type="primary" class="item-quantity-badge" />
                  </div>
                  <div class="item-bottom">
                    <div class="item-price">{{ item.price.toLocaleString() }}원</div>
                    <div class="item-total">{{ (item.price * item.quantity).toLocaleString() }}원</div>
                  </div>
                </div>
              </el-scrollbar>
            </div>
            
            <div class="order-footer">
              <div class="divider"></div>
              <div class="order-summary">
                <div class="summary-item">
                  <span class="label">총 품목:</span>
                  <span class="value">{{ order.order_items.length }}개</span>
                </div>
                <div class="summary-item">
                  <span class="label">총 수량:</span>
                  <span class="value">{{ getTotalQuantity(order) }}개</span>
                </div>
                <div class="summary-item total">
                  <span class="label">결제 금액:</span>
                  <span class="value price">{{ order.total_price?.toLocaleString() }}원</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { Check, Delete } from '@element-plus/icons-vue';
import type { Order } from '../../api/orderApi'; // 타입만 가져옵니다
import { ElNotification } from 'element-plus';

// 주문 데이터와 상태 관리
const orders = ref<Order[]>([]);
const completedOrderIds = ref<number[]>([]);
const newOrderIds = ref<number[]>([]);
const processedOrderIds = ref<number[]>([]); // 이미 처리한 주문 ID 추적
let eventSource: EventSource | null = null;
let isFirstBatch = true; // 첫 번째 데이터 배치 여부 추적

// 카드 레이아웃이 2행 고정인지 여부 (8개 이하인 경우)
const isTwoRowLayout = computed(() => {
  return sortedOrders.value.length <= 8;
});

// 완료되지 않은 주문만 표시 (최신 순)
const sortedOrders = computed(() => {
  return [...orders.value]
    .filter(order => !completedOrderIds.value.includes(order.id))
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
});

// SSE 연결 설정
function setupSSEConnection(): void {
  // 기존 연결이 있으면 닫기
  if (eventSource !== null) {
    eventSource.close();
  }
  
  // SSE 연결 생성
  eventSource = new EventSource('/api/orders/stream');
  
  // 연결 수립 이벤트
  eventSource.onopen = () => {
    console.log('SSE 연결 성공');
    
    // 연결 성공 알림
    ElNotification({
      title: '실시간 연결됨',
      message: '주문 시스템에 실시간으로 연결되었습니다.',
      type: 'success',
      duration: 3000,
    });
    
    // 최초 연결 시 초기화
    isFirstBatch = true;
    
    // 3초 후 초기 데이터 로드 완료로 간주
    setTimeout(() => {
      isFirstBatch = false;
    }, 3000);
  };
  
  // 메시지 수신 이벤트
  eventSource.onmessage = (event) => {
    try {
      // 수신된 주문 데이터 파싱
      const order = JSON.parse(event.data);
      
      // 이미 처리한 주문인지 확인 (중복 처리 방지)
      if (processedOrderIds.value.includes(order.id)) {
        return; // 이미 처리된 주문은 무시
      }
      
      // 초기 데이터 배치인 경우 기록하고 무시
      if (isFirstBatch) {
        processedOrderIds.value.push(order.id);
        return; // 초기 데이터 배치는 무시
      }
      
      // 기존 주문 업데이트 또는 새 주문 추가
      const orderIndex = orders.value.findIndex(o => o.id === order.id);
      
      if (orderIndex !== -1) {
        // 기존 주문 업데이트
        orders.value[orderIndex] = order;
      } else {
        // 새 주문 추가
        orders.value.push(order);
        
        // 새 주문 알림 표시
        showNewOrderNotification([order]);
        
        // 새로운 주문 ID 기록 (애니메이션용)
        newOrderIds.value.push(order.id);
        setTimeout(() => {
          // 5초 후 애니메이션 제거
          newOrderIds.value = newOrderIds.value.filter(id => id !== order.id);
        }, 5000);
      }
      
      // 처리된 주문 ID 기록
      processedOrderIds.value.push(order.id);
      
    } catch (error) {
      console.error('SSE 메시지 처리 중 오류 발생:', error);
    }
  };
  
  // 에러 처리
  eventSource.onerror = (error) => {
    console.error('SSE 연결 오류:', error);
    
    // 연결 종료
    if (eventSource) {
      eventSource.close();
      eventSource = null;
    }
    
    // 연결 오류 알림
    ElNotification({
      title: '연결 오류',
      message: '주문 시스템 연결이 끊겼습니다. 재연결을 시도합니다.',
      type: 'error',
      duration: 5000,
    });
    
    // 3초 후 재연결 시도
    setTimeout(() => {
      console.log('SSE 재연결 시도...');
      setupSSEConnection();
    }, 3000);
  };
}

// 새 주문 알림 표시
function showNewOrderNotification(newOrders: Order[]): void {
  // 새 주문 정보 요약
  const orderList = newOrders.map(order => {
    const items = order.order_items.map(item => `${item.menu.name} ${item.quantity}개`).join(', ');
    return `주문 #${order.id}: ${items}`;
  }).join('\n');
  
  ElNotification({
    title: `새로운 주문 ${newOrders.length}건이 접수되었습니다`,
    message: orderList,
    type: 'success',
    duration: 5000,
    position: 'top-right'
  });
}

// 주문 완료 처리 (프론트엔드에서만)
function completeOrder(orderId: number): void {
  // 완료된 주문 ID 목록에 추가
  completedOrderIds.value.push(orderId);
  
  // 알림 표시
  ElNotification({
    title: '주문 완료',
    message: `주문 #${orderId}이(가) 완료 처리되었습니다.`,
    type: 'success',
    duration: 3000,
  });
}

// 모든 주문 지우기
function clearAllOrders(): void {
  const currentOrderIds = sortedOrders.value.map(order => order.id);
  completedOrderIds.value.push(...currentOrderIds);
  
  ElNotification({
    title: '주문 초기화',
    message: '모든 주문이 완료 처리되었습니다.',
    type: 'info',
    duration: 3000,
  });
}

// 주문의 총 메뉴 수량 계산
function getTotalQuantity(order: Order): number {
  return order.order_items.reduce((total, item) => total + item.quantity, 0);
}

// 주문 카드 열 크기 계산
function getColSpan(orderCount: number): number {
  if (orderCount <= 8) {
    return 6; // 4개 열 (24/6=4 columns)
  } else if (orderCount <= 16) {
    return 4; // 6개 열 (24/4=6 columns)
  } else {
    return 3; // 8개 열 (24/3=8 columns)
  }
}

// 주문 날짜 포맷
function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleString('ko-KR', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

// 새 주문 여부 확인 (애니메이션용)
function isNewOrder(orderId: number): boolean {
  return newOrderIds.value.includes(orderId);
}

onMounted(() => {
  // SSE 연결 즉시 설정
  setupSSEConnection();
});

onUnmounted(() => {
  // 컴포넌트 언마운트 시 SSE 연결 종료
  if (eventSource) {
    console.log('SSE 연결 종료');
    eventSource.close();
    eventSource = null;
  }
});
</script>

<style scoped>
.order-container {
  width: 100%;
  padding: 20px;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.order-header h2 {
  margin: 0;
  font-size: 24px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.order-grid {
  width: 100%;
  max-width: 1600px;
  margin: 0 auto;
  overflow-y: auto;
  max-height: calc(100vh - 80px); /* 상단 여백 고려 */
}

.order-col {
  margin-bottom: 20px;
}

.order-card {
  height: 100%;
  transition: all 0.3s;
  border-radius: 12px;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.order-id {
  font-size: 16px;
}

.order-items {
  margin-bottom: 15px;
  height: calc(100% - 180px); /* 헤더와 푸터를 제외한 나머지 공간 */
  display: flex;
  flex-direction: column;
}

.item {
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 8px;
  background-color: #f9f9f9;
  transition: all 0.2s;
}

.item:hover {
  background-color: #f5f5f5;
  transform: translateY(-2px);
}

.item-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.item-name {
  font-weight: 600;
  font-size: 15px;
  word-break: break-word;
  overflow-wrap: break-word;
  max-width: 70%;
}

.item-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.item-price {
  color: #606266;
  font-size: 14px;
}

.item-total {
  font-weight: 600;
  color: #409eff;
  font-size: 14px;
}

.divider {
  height: 1px;
  background-color: #ebeef5;
  margin: 15px 0;
}

.order-footer {
  padding-top: 5px;
}

.order-summary {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.summary-item {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.summary-item.total {
  margin-top: 5px;
  font-weight: 600;
  font-size: 16px;
}

.label {
  color: #606266;
}

.value {
  font-weight: 500;
}

.value.price {
  color: #409eff;
  font-size: 18px;
  font-weight: 700;
}

.complete-btn {
  font-size: 16px;
}

/* 새 주문 애니메이션 */
.new-order {
  animation: newOrderEffect 1s ease-in-out;
}

@keyframes newOrderEffect {
  0% {
    transform: translateY(20px);
    opacity: 0;
    box-shadow: 0 0 0 0 rgba(64, 158, 255, 0.7);
  }
  50% {
    transform: translateY(0);
    opacity: 1;
    box-shadow: 0 0 20px 0 rgba(64, 158, 255, 0.7);
  }
  100% {
    transform: translateY(0);
    opacity: 1;
    box-shadow: 0 0 0 0 rgba(64, 158, 255, 0);
  }
}

/* 반응형 스타일 */
@media (max-width: 768px) {
  .order-container {
    padding: 10px;
  }
  
  .el-col {
    width: 100% !important;
    max-width: 100% !important;
    flex: 0 0 100% !important;
  }
  
  .order-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .actions {
    width: 100%;
  }
}
</style>