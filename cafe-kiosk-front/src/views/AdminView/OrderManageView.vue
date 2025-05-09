<template>
  <div class="order-container">
    <el-row :gutter="24" class="order-content">
      <!-- 주문 리스트 (좌측) -->
      <el-col :xs="24" :sm="24" :md="8" :lg="8">
        <el-card shadow="hover" class="order-list-card">
          <template #header>
            <div class="card-header">주문 목록</div>
          </template>
          <div class="order-list-scroll">
            <el-table
              :data="sortedOrders"
              style="width: 100%"
              @row-click="selectOrder"
              highlight-current-row
              :row-class-name="rowClassName"
              :header-cell-style="{ fontSize: '16px', fontWeight: 'bold', background: '#f5f7fa', height: '50px' }"
              :cell-style="{ fontSize: '15px', height: '60px' }"
              ref="orderTable"
            >
              <el-table-column prop="id" label="주문번호" width="100" align="center"/>
              <el-table-column prop="created_at" label="주문시간" min-width="150">
                <template #default="scope">
                  {{ formatDate(scope.row.created_at) }}
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-card>
      </el-col>

      <!-- 주문 상세 (우측) -->
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <el-card shadow="hover" class="order-detail-card">
          <template #header>
            <div class="card-header">
              주문 상세
              <span v-if="selectedOrder">(주문 #{{ selectedOrder.id }})</span>
            </div>
          </template>
          <el-empty v-if="!selectedOrder" description="주문을 선택하세요." />
          <template v-else>
            <el-table
              :data="selectedOrderItems"
              style="width: 100%"
              :header-cell-style="{ fontSize: '16px', fontWeight: 'bold', background: '#f5f7fa', height: '50px' }"
              :cell-style="{ fontSize: '15px', height: '60px' }"
            >
              <el-table-column prop="menu.name" label="메뉴명" min-width="180" />
              <el-table-column prop="quantity" label="수량" width="100" align="center" />
              <el-table-column prop="price" label="가격" width="120" align="right">
                <template #default="scope">
                  {{ scope.row.price.toLocaleString() }}원
                </template>
              </el-table-column>
            </el-table>
            <div class="order-total" v-if="selectedOrder.total_price">
              총 주문금액: <span class="price">{{ selectedOrder.total_price.toLocaleString() }}원</span>
            </div>
          </template>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue';
import type { Order } from '../../api/orderApi';
import { getOrders } from '../../api/orderApi';
import { ElNotification, ElTable } from 'element-plus';

// 타입 정의 (api에서 import가 불가능할 경우 여기서 선언)
// interface Order {
//   id: number;
//   created_at: string;
//   order_items: OrderItem[];
//   total_price?: number;
// }
// 
// interface OrderItem {
//   menu: {
//     name: string;
//   };
//   quantity: number;
//   price: number;
// }

const orders = ref<Order[]>([]);
const selectedOrderId = ref<number | null>(null);
const orderTable = ref<InstanceType<typeof ElTable> | null>(null);
const previousOrderIds = ref<number[]>([]);
const isFirstLoad = ref(true);
let eventSource: EventSource | null = null;

const sortedOrders = computed(() => {
  return [...orders.value].sort((a, b) =>
    new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  );
});

const selectedOrder = computed(() =>
  orders.value.find(order => order.id === selectedOrderId.value) || null
);

const selectedOrderItems = computed(() =>
  selectedOrder.value ? selectedOrder.value.order_items : []
);

// 초기 주문 목록 로드
async function fetchInitialOrders(): Promise<void> {
  try {
    const newOrders = await getOrders();
    orders.value = newOrders;
    
    // 주문이 있으면 첫 번째 주문 선택 (옵션)
    // if (newOrders.length > 0) {
    //   selectedOrderId.value = newOrders[0].id;
    // }
    
    // 이전 주문 ID 목록 초기화
    previousOrderIds.value = newOrders.map(order => order.id);
    
    // 첫 로드 완료 표시
    isFirstLoad.value = false;
    
  } catch (error) {
    console.error('주문 목록을 불러오는데 실패했습니다:', error);
  }
}

// SSE 연결 설정
function setupSSEConnection(): void {
  // 기존 연결이 있으면 닫기
  if (eventSource !== null) {
    eventSource.close();
  }
  
  // SSE 연결 생성
  eventSource = new EventSource('http://localhost:8080/api/orders/stream');
  
  // 연결 수립 이벤트
  eventSource.onopen = (event) => {
    console.log('SSE 연결이 수립되었습니다.');
  };
  
  // 메시지 수신 이벤트
  eventSource.onmessage = (event) => {
    try {
      // 수신된 주문 데이터 파싱
      const order = JSON.parse(event.data);
      
      // 새 주문인지 확인
      const isNewOrder = !previousOrderIds.value.includes(order.id);
      
      // 기존 주문 업데이트 또는 새 주문 추가
      const orderIndex = orders.value.findIndex(o => o.id === order.id);
      
      if (orderIndex !== -1) {
        // 기존 주문 업데이트
        orders.value[orderIndex] = order;
      } else {
        // 새 주문 추가
        orders.value.push(order);
        
        // 새 주문 알림 표시
        if (!isFirstLoad.value && isNewOrder) {
          showNewOrderNotification([order]);
          
          // 새 주문 선택 및 스크롤
          selectedOrderId.value = order.id;
          nextTick(() => {
            scrollToSelectedOrder();
          });
        }
      }
      
      // 이전 주문 ID 목록 업데이트
      if (isNewOrder) {
        previousOrderIds.value.push(order.id);
      }
      
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
    
    // 3초 후 재연결 시도
    setTimeout(() => {
      console.log('SSE 재연결 시도...');
      setupSSEConnection();
    }, 3000);
  };
}

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

function scrollToSelectedOrder(): void {
  // 선택된 주문이 테이블에서 보이도록 스크롤
  if (orderTable.value) {
    // Element Plus의 테이블 메서드 사용
    orderTable.value.scrollTo(
      sortedOrders.value.findIndex(order => order.id === selectedOrderId.value)
    );
  }
}

function selectOrder(row: Order): void {
  selectedOrderId.value = row.id;
}

function rowClassName(params: { row: Order }): string {
  return params.row.id === selectedOrderId.value ? 'current-row' : '';
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleString('ko-KR', {
    year: '2-digit',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

onMounted(() => {
  // 초기 주문 목록 로드 후 SSE 연결 설정
  fetchInitialOrders().then(() => {
    setupSSEConnection();
  });
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
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  min-height: 100%;
  padding: 20px;
  position: relative;
}

.order-content {
  width: 100%;
  max-width: 1200px;
}

.order-list-card,
.order-detail-card {
  height: 100%;
  margin-bottom: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
}

.order-list-scroll {
  max-height: calc(100vh - 200px);
  overflow-y: auto;
  scrollbar-width: none; /* Firefox */
}

.order-list-scroll::-webkit-scrollbar {
  display: none; /* Chrome, Safari, Opera */
}

/* 주문 상세 카드 고정 스타일 */
.el-col:nth-child(2) {
  position: sticky;
  top: 20px;
}

.card-header {
  font-weight: bold;
  font-size: 18px;
  padding: 5px 0;
}

.current-row {
  background: #f0f9eb !important;
  transition: background-color 0.3s;
}

.current-row:hover {
  background: #e6f7d9 !important;
}

.order-total {
  margin-top: 24px;
  padding: 16px;
  text-align: right;
  font-weight: bold;
  background: #f8f9fc;
  border-radius: 6px;
  font-size: 16px;
}

.price {
  color: #409eff;
  font-size: 22px;
  margin-left: 10px;
  font-weight: 700;
}

/* 테이블 내부 스타일 커스터마이징 */
:deep(.el-table__row) {
  cursor: pointer;
}

:deep(.el-table .cell) {
  line-height: 1.5;
  padding: 8px 12px;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar) {
  display: none;
}

:deep(.el-table__body-wrapper) {
  scrollbar-width: none;
}

/* 새로운 주문 행 스타일 - 잠시 깜빡임 효과 */
@keyframes newOrderHighlight {
  0% { background-color: rgba(64, 158, 255, 0.2); }
  100% { background-color: transparent; }
}

.new-order-row {
  animation: newOrderHighlight 2s ease-in-out;
}

/* 모바일 반응형 스타일 */
@media (max-width: 768px) {
  .order-container {
    padding: 10px;
    flex-direction: column;
  }
  
  .el-row {
    display: flex;
    flex-direction: column;
  }
  
  .el-col {
    width: 100%;
    max-width: 100%;
  }
  
  .el-col:nth-child(2) {
    position: relative;
    top: 0;
  }
  
  .order-list-scroll {
    max-height: 400px;
  }
  
  .order-total {
    padding: 12px;
  }
}
</style>