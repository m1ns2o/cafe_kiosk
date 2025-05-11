<template>
  <div class="statistics-container">
    <el-row :gutter="24" class="statistics-content">
      <!-- 통계 정보 (좌측) -->
      <el-col :xs="24" :sm="24" :md="14" :lg="14" class="statistics-col">
        <el-card shadow="hover" class="statistics-card">
          <template #header>
            <div class="card-header">판매 통계</div>
          </template>
          
          <div class="statistics-card-content">
            <!-- 항목별 판매비율 파이차트 -->
            <div class="chart-container pie-container">
              <h3 class="chart-title">항목별 판매비율</h3>
              <div class="chart-wrapper pie-chart-wrapper">
                <div v-if="isLoading" class="chart-loading">
                  <el-skeleton animated :rows="3" />
                </div>
                <div v-else-if="menuSalesData.labels.length === 0" class="chart-empty">
                  <el-empty description="판매 데이터가 없습니다." />
                </div>
                <div v-else class="chart-container-inner">
                  <canvas ref="pieChartRef"></canvas>
                </div>
              </div>
            </div>
            
            <!-- 최근 이틀 매출 바차트 -->
            <div class="chart-container bar-container">
              <h3 class="chart-title">최근 이틀 매출</h3>
              <div class="chart-wrapper bar-chart-wrapper">
                <div v-if="isLoading" class="chart-loading">
                  <el-skeleton animated :rows="3" />
                </div>
                <div v-else-if="dailySalesData.labels.length === 0" class="chart-empty">
                  <el-empty description="매출 데이터가 없습니다." />
                </div>
                <div v-else class="chart-container-inner">
                  <canvas ref="barChartRef"></canvas>
                </div>
              </div>
            </div>
            
            <!-- 총 매출 -->
            <div class="total-sales-container">
              <div class="total-sales-card">
                <h3 class="total-sales-title">총 매출</h3>
                <div v-if="isLoading" class="total-loading">
                  <el-skeleton animated :rows="1" />
                </div>
                <div v-else class="total-sales-amount">
                  {{ totalSales.toLocaleString() }}원
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 주문 내역 (우측) -->
      <!-- 주문 내역 (우측) -->
<el-col :xs="24" :sm="24" :md="10" :lg="10" class="order-history-col">
  <el-card shadow="hover" class="order-history-card">
    <template #header>
      <div class="card-header">
        주문 내역
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="~"
          start-placeholder="시작일"
          end-placeholder="종료일"
          format="YYYY-MM-DD"
          size="small"
          style="width: 200px; margin-left: 10px;"
          @change="fetchOrderHistory"
        />
      </div>
    </template>
    
    <div class="order-table-wrapper">
      <div v-if="isLoading" class="table-loading">
        <el-skeleton animated :rows="10" />
      </div>
      <div v-else-if="orderHistory.length === 0" class="table-empty">
        <el-empty description="주문 내역이 없습니다." />
      </div>
      <el-table
        v-else
        :data="orderHistory"
        style="width: 100%"
        class="order-table"
        :header-cell-style="{ fontSize: '15px', fontWeight: 'bold', background: '#f5f7fa', height: '50px' }"
        :cell-style="{ fontSize: '14px', height: '50px' }"
      >
        <el-table-column type="expand">
          <template #default="props">
            <div class="order-detail-wrapper">
              <h4 class="order-detail-title">주문 상세 정보</h4>
              <el-table
                :data="props.row.order_items"
                style="width: 100%"
                class="order-detail-table"
                :show-header="true"
                size="small"
                border
              >
                <el-table-column label="상품명" min-width="120">
                  <template #default="scope">
                    {{ scope.row.menu.name || '알 수 없는 메뉴' }}
                  </template>
                </el-table-column>
                <el-table-column prop="quantity" label="수량" width="70" align="center" />
                <el-table-column label="단가" width="100" align="right">
                  <template #default="scope">
                    {{ scope.row.price.toLocaleString() }}원
                  </template>
                </el-table-column>
                <el-table-column label="소계" width="100" align="right">
                  <template #default="scope">
                    {{ (scope.row.price * scope.row.quantity).toLocaleString() }}원
                  </template>
                </el-table-column>
              </el-table>
              <div class="order-total">
                <strong>합계: {{ props.row.total_price.toLocaleString() }}원</strong>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="주문번호" width="80" align="center"/>
        <el-table-column prop="created_at" label="주문시간" min-width="80">
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="total_price" label="주문금액" width="120" align="right">
          <template #default="scope">
            {{ scope.row.total_price.toLocaleString() }}원
          </template>
        </el-table-column>
      </el-table>
    </div>
  </el-card>
</el-col>


    </el-row>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue';
import { getOrdersByPeriod } from '../../api/orderApi';
import type { Order } from '../../api/orderApi';
import Chart from 'chart.js/auto';
// 화면 크기 상태
const windowWidth = ref(window.innerWidth);
const windowHeight = ref(window.innerHeight);

// 차트 크기 계산
const getChartHeight = computed(() => {
  // 화면 높이에 따라 동적으로 차트 높이 계산
  const baseHeight = windowHeight.value;
  
  if (baseHeight < 600) {
    return 180; // 작은 화면
  } else if (baseHeight < 800) {
    return 220; // 중간 화면
  } else {
    return 250; // 큰 화면
  }
});

// 상태 변수
const isLoading = ref(true);
const orderHistory = ref<Order[]>([]);
const totalSales = ref(0);

// 차트 참조
const pieChartRef = ref<HTMLCanvasElement | null>(null);
const barChartRef = ref<HTMLCanvasElement | null>(null);
let pieChart: Chart | null = null;
let barChart: Chart | null = null;

// 차트 데이터
const menuSalesData = ref({
  labels: [] as string[],
  datasets: [{
    data: [] as number[],
    backgroundColor: [] as string[],
  }]
});

const dailySalesData = ref({
  labels: [] as string[],
  datasets: [{
    label: '매출',
    data: [] as number[],
    backgroundColor: '#78A1BB',
  }]
});

// 날짜 범위 선택 (기본값: 최근 일주일)
const dateRange = ref([
  new Date(new Date().setDate(new Date().getDate() - 7)),
  new Date()
]);

// 차트 색상 - 블루와 그레이 계열 6개
const COLORS = [
  '#78A1BB', // 기준 파스텔 블루
  '#5D88A5', // 더 진한 블루
  '#96B6CB', // 밝은 파스텔 블루
  '#A0ACBD', // 블루 틴트가 있는 그레이
  '#8495A8', // 진한 블루-그레이
  '#6E7F91'  // 어두운 블루-그레이
];

// 화면 크기 변경 감지 함수
function handleResize() {
  windowWidth.value = window.innerWidth;
  windowHeight.value = window.innerHeight;
  
  // 차트 다시 그리기 - 약간의 지연을 줘서 DOM이 업데이트된 후 실행
  if (!isLoading.value) {
    nextTick(() => {
      setTimeout(() => {
        createOrUpdateCharts();
      }, 100);
    });
  }
}

// 날짜 포맷 함수
function formatDateTime(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleString('ko-KR', {
    year: '2-digit',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

function formatDate(date: Date): string {
  // 'YYYY-MM-DD' 형식으로 변환 (백엔드 API 요구사항)
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  
  return `${year}-${month}-${day}`;
}
// 주문 내역 가져오기
// 주문 내역 가져오기
async function fetchOrderHistory(): Promise<void> {
  isLoading.value = true;

  try {
    // 날짜 포맷 변환
    const startDate = formatDate(dateRange.value[0]);
    const endDate = formatDate(dateRange.value[1]);
    
    // 백엔드 API를 사용하여 기간별 주문 조회
    const response = await getOrdersByPeriod(startDate, endDate, {
      sortBy: 'created_at',
      order: 'desc'
    });
    
    // 응답에서 주문 목록 추출
    orderHistory.value = response.orders;
    
    // 통계 데이터 계산
    calculateStatistics(response.orders);
    
    // 차트 업데이트/생성
    nextTick(() => {
      setTimeout(() => {
        createOrUpdateCharts();
      }, 100);
    });
  } catch (error) {
    console.error('주문 내역을 불러오는데 실패했습니다:', error);
  } finally {
    isLoading.value = false;
  }
}
// 통계 데이터 계산
function calculateStatistics(orders: Order[]): void {
  // 1. 메뉴별 판매량 계산 (파이차트용)
  const menuSales = new Map<string, number>();
  
  orders.forEach(order => {
    order.order_items.forEach(item => {
      const menuName = item.menu.name;
      const quantity = item.quantity;
      menuSales.set(menuName, (menuSales.get(menuName) || 0) + quantity);
    });
  });
  
  // 메뉴 판매 데이터를 차트.js 형식으로 변환
  const menuNames: string[] = [];
  const menuQuantities: number[] = [];
  const menuColors: string[] = [];
  
  // 색상 동적 생성 함수 (항목이 COLORS 배열보다 많을 경우 사용)
  const generateColor = (index: number): string => {
    // 기본 색상 사용
    if (index < COLORS.length) {
      return COLORS[index];
    }
    
    // 기본 색상 #78A1BB를 HSL로 변환 (여기서는 대략적인 값 사용)
    // 실제 #78A1BB는 HSL로 약 (203, 33%, 66%)
    const baseHue = 203;
    const baseSaturation = 33;
    const baseLightness = 66;
    
    // 항목 수에 따라 색조(hue)만 골고루 분포
    const hueStep = 360 / (index + COLORS.length);
    const newHue = (baseHue + (index * hueStep)) % 360;
    
    // 채도와 명도는 비슷하게 유지하되 약간의 변화를 줌
    const newSaturation = baseSaturation + (index % 3) * 5;
    const newLightness = baseLightness + (index % 4 - 2) * 5;
    
    return `hsl(${newHue}, ${newSaturation}%, ${newLightness}%)`;
  };
  
  let i = 0;
  menuSales.forEach((quantity, name) => {
    menuNames.push(name);
    menuQuantities.push(quantity);
    menuColors.push(generateColor(i));
    i++;
  });
  
  // 차트 데이터 업데이트
  menuSalesData.value = {
    labels: menuNames,
    datasets: [{
      data: menuQuantities,
      backgroundColor: menuColors,
    }]
  };
  
  // 2. 최근 이틀 매출 계산 (바차트용)
  const dailySales = new Map<string, number>();
  
  // 현재 날짜와 어제 날짜 계산
  const today = new Date();
  const yesterday = new Date(today);
  yesterday.setDate(yesterday.getDate() - 1);
  
  const todayStr = formatDate(today);
  const yesterdayStr = formatDate(yesterday);
  
  // 초기값 설정 (데이터가 없는 날도 표시)
  dailySales.set(yesterdayStr, 0);
  dailySales.set(todayStr, 0);
  
  // 3. 총 매출 계산 - 더 효율적인 방법 (한 번의 루프로 두 가지 계산)
  let totalSalesAmount = 0;
  
  orders.forEach(order => {
    // 총 매출에 더하기
    totalSalesAmount += order.total_price;
    
    // 날짜별 매출에 더하기
    const orderDate = formatDate(new Date(order.created_at));
    if (orderDate === todayStr || orderDate === yesterdayStr) {
      dailySales.set(orderDate, (dailySales.get(orderDate) || 0) + order.total_price);
    }
  });
  
  // 총 매출 업데이트
  totalSales.value = totalSalesAmount;
  
  // 날짜순으로 정렬된 데이터 생성
  const sortedDates = Array.from(dailySales.keys()).sort();
  const salesData = sortedDates.map(date => dailySales.get(date) || 0);
  
  // 차트 데이터 업데이트
  dailySalesData.value = {
    labels: sortedDates,
    datasets: [{
      label: '매출',
      data: salesData,
      backgroundColor: '#78A1BB',  // 메인 색상 사용
    }]
  };
}

// 차트 생성 또는 업데이트
function createOrUpdateCharts(): void {
  // 파이 차트 생성/업데이트
  if (pieChartRef.value) {
    // 기존 차트가 있으면 제거
    if (pieChart) {
      pieChart.destroy();
    }
    
    // 새 차트 생성
    pieChart = new Chart(pieChartRef.value, {
      type: 'pie',
      data: menuSalesData.value,
      options: {
        responsive: true,
        maintainAspectRatio: false, // 컨테이너에 맞게 크기 조정
        plugins: {
          legend: {
            position: windowWidth.value < 768 ? 'bottom' : 'right', // 화면 크기에 따라 범례 위치 변경
            labels: {
              boxWidth: windowWidth.value < 768 ? 10 : 12, // 화면 크기에 따라 범례 크기 조정
              font: {
                size: windowWidth.value < 768 ? 10 : 12
              }
            }
          },
          tooltip: {
            callbacks: {
              label: function(context) {
                const label = context.label || '';
                const value = context.raw as number;
                const total = (context.chart.data.datasets[0].data as number[]).reduce((a, b) => (a as number) + (b as number), 0) as number;
                const percentage = Math.round((value / total) * 100);
              
                return `${label}: ${value}개 (${percentage}%)`;
              }
            }
          }
        }
      }
    });
  }
  
  // 바 차트 생성/업데이트
  if (barChartRef.value) {
    // 기존 차트가 있으면 제거
    if (barChart) {
      barChart.destroy();
    }
    
    // 새 차트 생성
    barChart = new Chart(barChartRef.value, {
      type: 'bar',
      data: dailySalesData.value,
      options: {
        responsive: true,
        maintainAspectRatio: false, // 컨테이너에 맞게 크기 조정
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              callback: function(value) {
                return value.toLocaleString() + '원';
              },
              font: {
                size: windowWidth.value < 768 ? 10 : 12
              }
            }
          },
          x: {
            ticks: {
              font: {
                size: windowWidth.value < 768 ? 10 : 12
              }
            }
          }
        },
        plugins: {
          legend: {
            display: windowWidth.value >= 768, // 작은 화면에서는 범례 숨김
            labels: {
              font: {
                size: 12
              }
            }
          },
          tooltip: {
            callbacks: {
              label: function(context) {
                const value = context.raw as number;
                return `매출: ${value.toLocaleString()}원`;
              }
            }
          }
        }
      }
    });
  }
  
  console.log('차트 생성/업데이트 완료:', 
              '파이차트:', pieChart ? 'O' : 'X', 
              '바차트:', barChart ? 'O' : 'X');
}

// 일정 시간 후에만 리사이즈 이벤트 실행 (디바운스)
let resizeTimeout: ReturnType<typeof setTimeout> | null = null;
function debouncedResize() {
  if (resizeTimeout) {
    clearTimeout(resizeTimeout);
  }
  resizeTimeout = setTimeout(() => {
    handleResize();
  }, 200);
}

// 컴포넌트 마운트 시 데이터 로드 및 이벤트 리스너 등록
onMounted(() => {
  fetchOrderHistory();
  window.addEventListener('resize', debouncedResize);
  
  // 초기 화면 크기 설정
  handleResize();
});

// 컴포넌트 언마운트 시 차트 인스턴스 제거 및 이벤트 리스너 해제
onUnmounted(() => {
  if (pieChart) {
    pieChart.destroy();
    pieChart = null;
  }
  
  if (barChart) {
    barChart.destroy();
    barChart = null;
  }
  
  window.removeEventListener('resize', debouncedResize);
  
  if (resizeTimeout) {
    clearTimeout(resizeTimeout);
  }
});

// 화면 크기 변경 시 차트 높이 자동 조정
watch(getChartHeight, () => {
  const pieWrapper = document.querySelector('.pie-chart-wrapper') as HTMLElement;
  const barWrapper = document.querySelector('.bar-chart-wrapper') as HTMLElement;
  
  if (pieWrapper) {
    pieWrapper.style.height = `${getChartHeight.value}px`;
  }
  
  if (barWrapper) {
    barWrapper.style.height = `${getChartHeight.value}px`;
  }
});
</script>

<style scoped>

.statistics-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
}

.statistics-content {
  width: 100%;
  max-width: 1400px;
  height: 100%;
}

.statistics-col,
.order-history-col {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.statistics-card,
.order-history-card {
  height: 100%;
  margin-bottom: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
}

.statistics-card-content {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: auto;
}

:deep(.el-card__body) {
  flex: 1;
  overflow: auto;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  align-items: center;
  font-weight: bold;
  font-size: 18px;
  padding: 5px 0;
}

.chart-container {
  margin-bottom: 15px;
  flex-shrink: 0;
}

.pie-container, .bar-container {
  flex: 1;
  min-height: 0; /* 중요: flexbox에서 오버플로우를 방지 */
  display: flex;
  flex-direction: column;
}

.chart-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #606266;
}

.chart-wrapper {
  flex: 1;
  overflow: hidden;
  background-color: #f8f9fb;
  border-radius: 8px;
  padding: 10px;
  position: relative;
  min-height: 0; /* 중요: flexbox에서 오버플로우를 방지 */
}

.pie-chart-wrapper,
.bar-chart-wrapper {
  transition: height 0.3s ease; /* 높이 변경 시 애니메이션 효과 */
}

.chart-container-inner {
  width: 100%;
  height: 100%;
  position: relative;
}

.total-sales-container {
  margin-top: 5px;
  padding: 10px;
  background-color: #f0f9ff;
  border-radius: 8px;
  text-align: center;
  flex-shrink: 0;
}

.total-sales-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.total-sales-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 10px;
  color: #303133;
}

.total-sales-amount {
  font-size: 28px;
  font-weight: 700;
  color: #409eff;
}

.order-table-wrapper {
  flex: 1;
  position: relative;
  min-height: 400px; /* 고정된 최소 높이로 복구 */
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.order-table {
  flex: 1;
  height: 100%;
  overflow: auto;
}
.order-detail-wrapper {
  padding: 10px 20px 15px;
  background-color: #f8fafc;
  border-radius: 4px;
}

.order-detail-title {
  font-size: 14px;
  color: #606266;
  margin: 5px 0 10px;
}

.order-detail-table {
  margin-bottom: 10px;
}

.order-total {
  text-align: right;
  color: #409eff;
  font-size: 15px;
  padding: 5px 5px 0;
  border-top: 1px dashed #e0e0e0;
}

/* 확장 행의 화살표 스타일 */
:deep(.el-table__expand-icon) {
  font-size: 16px;
  color: #409eff;
  transition: transform 0.2s ease;
}

:deep(.el-table__expand-icon--expanded) {
  transform: rotate(90deg);
}

/* 확장 영역의 배경색 설정 */
:deep(.el-table__expanded-cell) {
  background-color: #f8fafc;
}

/* 확장행 호버 효과 강화 */
:deep(.el-table__row:hover) {
  cursor: pointer;
  background-color: #f0f9ff !important;
}

/* 모바일 대응 */
@media (max-width: 768px) {
  .order-detail-wrapper {
    padding: 10px;
  }
  
  .order-detail-title {
    font-size: 13px;
    margin-bottom: 8px;
  }
  
  .order-total {
    font-size: 14px;
  }
}


.chart-loading,
.table-loading,
.total-loading {
  padding: 15px;
}

.chart-empty,
.table-empty {
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px 0;
}

/* 모바일 반응형 스타일 */
@media (max-width: 768px) {
  .statistics-container {
    padding: 10px;
    height: auto;
  }
  
  .statistics-content {
    height: auto;
  }
  
  .chart-container {
    margin-bottom: 10px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    font-size: 16px;
  }
  
  .chart-title {
    font-size: 14px;
    margin-bottom: 5px;
  }
  
  .total-sales-title {
    font-size: 16px;
    margin-bottom: 5px;
  }
  
  .total-sales-amount {
    font-size: 24px;
  }
  
  .el-date-picker {
    margin-left: 0 !important;
    margin-top: 10px;
    width: 100% !important;
  }
  
  /* 모바일에서의 차트 컨테이너 크기는 computed 속성으로 동적 조절 */
}

/* 화면 높이가 매우 낮을 때 (작은 디스플레이) */
@media (max-height: 600px) {
  .chart-container {
    margin-bottom: 8px;
  }
  
  .chart-title {
    margin-bottom: 3px;
  }
  
  .total-sales-container {
    padding: 5px;
  }
  
  .total-sales-title {
    margin-bottom: 3px;
  }
}
</style>