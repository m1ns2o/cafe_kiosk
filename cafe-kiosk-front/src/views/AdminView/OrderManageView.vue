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
            >
              <el-table-column prop="id" label="주문번호" width="100" />
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
import { ref, computed, onMounted } from 'vue';
import type { Order, OrderItem } from '../../api/orderApi';
import { getOrders } from '../../api/orderApi';

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

async function fetchOrders(): Promise<void> {
  try {
    orders.value = await getOrders();
  } catch (error) {
    console.error('주문 목록을 불러오는데 실패했습니다:', error);
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
  fetchOrders();
});
</script>

<style scoped>
.order-container {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  width: 100%;
  min-height: 80vh;
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