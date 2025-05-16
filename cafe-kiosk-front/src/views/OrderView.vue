<template>
  <div class="order-view-container">
    <div class="order-view">
      <!-- 오류 메시지 표시 -->
      <el-alert
        v-if="error"
        :title="error"
        type="error"
        :closable="false"
        show-icon
      />

      <!-- 카테고리 탭 -->
      <div class="category-container">
        <el-tabs 
          v-model="selectedCategory" 
          type="border-card"
          class="category-tabs"
          @tab-click="handleTabClick"
        >
          <el-tab-pane
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :name="category.id"
          >
            <!-- 탭 내용은 비워두고 아래 메뉴 섹션에서 컨텐츠 표시 -->
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 메뉴 리스트 -->
      <section class="menu-section">
        <el-empty 
          v-if="menuItems.length === 0" 
          description="이 카테고리에 메뉴가 없습니다."
        />
        
        <div v-else class="menu-container">
          <div class="menu-list">
            <el-card
              v-for="item in paginatedMenuItems"
              :key="item.id"
              class="menu-item"
              shadow="hover"
              @click="addToCart(item)"
            >
              <div class="menu-image">
                <el-image
                  :src="`http://localhost:8080${item.image_url}`"
                  :alt="item.name"
                  fit="contain"
                  loading="lazy"
                >
                  <template #error>
                    <div class="image-error">
                      <el-icon><picture-rounded /></el-icon>
                    </div>
                  </template>
                </el-image>
              </div>
              <div class="menu-content">
                <h3 class="menu-name">{{ item.name }}</h3>
                <div class="menu-price">{{ item.price.toLocaleString() }}원</div>
              </div>
            </el-card>
          </div>
          
          <!-- 페이지네이션 컨트롤 -->
          <div class="pagination-wrapper" v-if="menuItems.length > itemsPerPage">
            <el-pagination
              v-model:current-page="currentPage"
              :page-size="itemsPerPage"
              layout="prev, pager, next"
              :total="menuItems.length"
              background
              @current-change="handleCurrentChange"
            />
          </div>
        </div>
      </section>

      <!-- 장바구니 영역 -->
      <section class="cart-section">
        <div class="cart-row">
          <div class="cart-container">
            <div class="cart-header">
              <h2 class="cart-title">장바구니</h2>
              <el-badge :value="cartItems.length" :hidden="cartItems.length === 0" type="primary">
                <el-icon><shopping-cart /></el-icon>
              </el-badge>
            </div>
            
            <div v-if="cartItems.length === 0" class="empty-cart">
              <el-empty 
                description="장바구니가 비어있습니다."
                :image-size="80"
              />
            </div>
            
            <div v-else class="cart-items">
              <transition-group name="cart-item">
                <div v-for="(cartItem, index) in cartItems" :key="index" class="cart-item">
                  <div class="cart-item-info">
                    <h4 class="cart-item-name">{{ cartItem.item.name }}</h4>
                    <div class="cart-item-price">{{ cartItem.item.price.toLocaleString() }}원</div>
                  </div>
                  
                  <div class="cart-item-actions">
                    <div class="quantity-control">
                      <el-button
                        type="primary"
                        circle
                        size="small"
                        class="quantity-btn minus-btn"
                        @click.stop="decreaseQuantity(index)"
                      >
                        <span class="minus-icon">-</span>
                      </el-button>
                      <span class="quantity">{{ cartItem.quantity }}</span>
                      <el-button
                        type="primary"
                        circle
                        size="small"
                        class="quantity-btn plus-btn"
                        @click.stop="increaseQuantity(index)"
                      >
                        <el-icon><plus /></el-icon>
                      </el-button>
                    </div>
                    <el-button 
                      type="danger" 
                      circle 
                      size="small"
                      @click.stop="removeFromCart(index)"
                    >
                      <el-icon><delete /></el-icon>
                    </el-button>
                  </div>
                </div>
              </transition-group>
            </div>
          </div>
          
          <div class="summary-container">
            <el-card class="cart-summary" shadow="never">
              <div class="cart-total">
                <span>총 금액:</span>
                <span class="total-amount">{{ totalAmount.toLocaleString() }}원</span>
              </div>
              <div class="cart-actions">
                <div class="button-column">
                  <el-button 
                    @click="clearCart" 
                    class="clear-cart-btn" 
                    type="info"
                    :disabled="cartItems.length === 0"
                    plain
                  >
                    <el-icon class="button-icon"><delete /></el-icon>
                    전체 삭제
                  </el-button>
                  <el-button 
                    @click="placeOrder" 
                    class="order-btn" 
                    type="primary"
                    :disabled="cartItems.length === 0"
                  >
                    <el-icon class="button-icon"><shopping-cart-full /></el-icon>
                    주문하기
                  </el-button>
                </div>
              </div>
            </el-card>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import { CategoryAPI } from '../api/menu';
import type { MenuItem, Category, CartItem } from '../types/menuType';
import { useRouter } from 'vue-router';
import { 
  PictureRounded, 
  Delete, 
  ShoppingCart, 
  ShoppingCartFull,
  Plus
} from '@element-plus/icons-vue';

const router = useRouter();

const menuItems = ref<MenuItem[]>([]);
const selectedCategory = ref<number>(0);
const categories = ref<Category[]>([]);
const error = ref<string | null>(null);
const cartItems = ref<CartItem[]>([]);

// 페이지네이션 관련 변수
const currentPage = ref(1);
const itemsPerPage = 8;

// Element Plus 탭 클릭 핸들러
const handleTabClick = (tab: any) => {
  selectedCategory.value = parseInt(tab.props.name);
};

// 페이지네이션 핸들러
const handleCurrentChange = (page: number) => {
  currentPage.value = page;
};

// 현재 페이지에 표시할 메뉴 아이템
const paginatedMenuItems = computed(() => {
  const startIndex = (currentPage.value - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  return menuItems.value.slice(startIndex, endIndex);
});

// 총 페이지 수
const totalPages = computed(() => {
  return Math.ceil(menuItems.value.length / itemsPerPage);
});

// 카테고리가 변경되면 해당 카테고리의 메뉴를 불러오는 함수
const loadMenuItems = async (categoryId: number) => {
  if (!categoryId) return;
  
  try {
    error.value = null;
    const menuResponse = await CategoryAPI.getMenus(categoryId);
    menuItems.value = menuResponse.data;
    console.log('메뉴 아이템 로드:', menuItems.value);
  } catch (err) {
    console.error('메뉴를 불러오는 중 오류가 발생했습니다:', err);
    error.value = '메뉴를 불러오는 중 오류가 발생했습니다.';
    menuItems.value = [];
  }
};

// 장바구니에 아이템 추가
const addToCart = (item: MenuItem) => {
  const existingItem = cartItems.value.find(cartItem => cartItem.item.id === item.id);
  
  if (existingItem) {
    existingItem.quantity += 1;
  } else {
    cartItems.value.push({ item, quantity: 1 });
  }
};

// 장바구니에서 아이템 수량 증가
const increaseQuantity = (index: number) => {
  cartItems.value[index].quantity += 1;
};

// 장바구니에서 아이템 수량 감소
const decreaseQuantity = (index: number) => {
  if (cartItems.value[index].quantity > 1) {
    cartItems.value[index].quantity -= 1;
  }
};

// 장바구니에서 아이템 삭제
const removeFromCart = (index: number) => {
  cartItems.value.splice(index, 1);
};

// 장바구니 비우기
const clearCart = () => {
  cartItems.value = [];
};

// 총 주문 금액 계산
const totalAmount = computed(() => {
  return cartItems.value.reduce((total, cartItem) => {
    return total + (cartItem.item.price * cartItem.quantity);
  }, 0);
});

// 주문하기
const placeOrder = () => {
  if (cartItems.value.length === 0) {
    ElMessage.warning('장바구니가 비어있습니다.');
    return;
  }
  
  router.push({
    name: 'PaymentView',
    params: {
      totalAmount: totalAmount.value.toString(),
      cartItems: encodeURIComponent(JSON.stringify(cartItems.value))
    }
  });
};

// 카테고리가 변경될 때마다 메뉴 아이템 업데이트 및 페이지 초기화
watch(() => selectedCategory.value, async (newCategoryId) => {
  if (newCategoryId) {
    await loadMenuItems(newCategoryId);
    currentPage.value = 1;
  }
});

onMounted(async () => {
  try {
    error.value = null;
    
    // 카테고리 불러오기
    const response = await CategoryAPI.getAllCategories();
    categories.value = response.data;
    console.log('카테고리 로드:', categories.value);
    
    // 첫 번째 카테고리를 기본값으로 설정
    if (categories.value.length > 0) {
      selectedCategory.value = categories.value[0].id;
      // 초기 메뉴 아이템 로드
      await loadMenuItems(selectedCategory.value);
    }
  } catch (err) {
    console.error('카테고리를 불러오는 중 오류가 발생했습니다:', err);
    error.value = '카테고리를 불러오는 중 오류가 발생했습니다.';
    categories.value = [];
  }
});
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Gowun+Dodum&display=swap');

.order-view-container {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  font-family: 'Gowun Dodum', sans-serif;
}

.order-view {
  background: var(--el-bg-color);
  width: 100%;
  height: 100%;
  margin: 0 auto;
  position: relative;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 카테고리 탭 스타일 개선 */
.category-container {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  padding: 10px 15px 0;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  z-index: 10;
  flex-shrink: 0; /* 탭이 줄어들지 않도록 설정 */
}

.category-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
  border-bottom: none;
}

.category-tabs :deep(.el-tabs__nav) {
  border-radius: 10px 10px 0 0;
  overflow: hidden;
}

.category-tabs :deep(.el-tabs__item) {
  height: 50px;
  line-height: 50px;
  font-size: 16px;
  font-weight: 500;
  transition: all 0.3s ease;
  padding: 0 24px;
  background-color: rgba(var(--el-color-primary-rgb), 0.05);
  position: relative;
  letter-spacing: 0.5px;
  min-width: 200px;
}

.category-tabs :deep(.el-tabs__item.is-active) {
  color: #fff;
  background-color: var(--el-color-primary);
  font-weight: 600;
  font-size: 17px;
  box-shadow: 0 -3px 10px rgba(0, 0, 0, 0.1);
}

.category-tabs :deep(.el-tabs__item:not(.is-active):hover) {
  color: var(--el-color-primary);
  background-color: rgba(var(--el-color-primary-rgb), 0.1);
}

.category-tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.category-tabs :deep(.el-tabs__content) {
  display: none;
}

/* 메뉴 섹션 스타일 개선 */
.menu-section {
  padding: 15px;
  flex: 1;
  overflow-y: auto;
  background: var(--el-bg-color-page);
  display: flex;
  flex-direction: column;
}

.menu-container {
  display: flex;
  flex-direction: column;
  gap: 15px;
  height: 100%;
}

.menu-list {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 15px;
  width: 100%;
  flex: 1;
  min-height: 0; /* 그리드가 flex 내에서 제대로 작동하도록 */
}

.menu-item {
  border-radius: 12px;
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  /* width: 100%; */
  height: 250px; /* 고정 높이 설정 */
}

.menu-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 20px rgba(0, 0, 0, 0.08);
}

.menu-image {
  height: 160px; /* 고정 이미지 높이 */
  overflow: hidden;
  position: relative;
  background-color: #f8f9fa;
  display: flex;
  align-items: center;
  justify-content: center;
}

.menu-image :deep(.el-image) {
  width: 100%;
  height: 100%;
  transition: transform 0.4s ease;
}

.menu-image :deep(.el-image__inner) {
  object-fit: contain !important;
  max-height: 100%;
  max-width: 100%;
  padding: 5px;
}

.menu-item:hover .menu-image :deep(.el-image) {
  transform: scale(1.05);
}

.image-error {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-secondary);
  font-size: 24px;
}

.menu-content {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  position: relative;
  flex: 1; /* 남은 공간 차지 */
}

.menu-name {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.menu-price {
  font-size: 14px;
  font-weight: 700;
  color: var(--el-color-primary);
  margin-top: auto; /* 가격이 항상 하단에 위치하도록 */
}

/* 페이지네이션 스타일 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 10px;
  padding: 5px 0;
  flex-shrink: 0; /* 페이지네이션이 압축되지 않도록 */
}

.pagination-wrapper :deep(.el-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
  /* background-color: var(--el-color-primary); */
  font-weight: 700;
  font-size: 14px;
}

.pagination-wrapper :deep(.el-pagination) {
  padding: 5px;
  border-radius: 30px;
  /* background: rgba(var(--el-color-primary-rgb), 0.05); */
  display: inline-flex;
}

/* 카트 영역 스타일 개선 */
.cart-section {
  background: linear-gradient(to bottom, #f8f9fa, #edf1f5);
  border-top: 1px solid var(--el-border-color-light);
  padding: 12px;
  height: 200px; /* 고정된 높이로 설정 */
  flex-shrink: 0; /* 카트 영역이 압축되지 않도록 */
  overflow: hidden; /* 내부 스크롤만 허용 */
}

.cart-row {
  display: flex;
  gap: 15px;
  height: 100%;
}

.cart-container {
  flex: 2;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.cart-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  flex-shrink: 0; /* 헤더가 압축되지 않도록 */
}

.cart-title {
  font-size: 16px;
  margin: 0;
  color: var(--el-text-color-primary);
  display: flex;
  align-items: center;
  gap: 8px;
}

.cart-title::before {
  content: '';
  display: block;
  width: 4px;
  height: 16px;
  background-color: var(--el-color-primary);
  border-radius: 2px;
}

/* 빈 장바구니 중앙 정렬 */
.empty-cart {
  display: flex;
  align-items: center;
  justify-content: center;
  height: calc(100% - 40px); /* 헤더 높이 제외한 높이 */
}

.cart-items {
  display: flex;
  flex-direction: column;
  gap: 10px;
  height: calc(100% - 40px); /* 헤더 높이 제외한 높이 */
  overflow-y: auto;
  padding-right: 8px;
}

.cart-items::-webkit-scrollbar {
  width: 6px;
}

.cart-items::-webkit-scrollbar-thumb {
  background: var(--el-color-primary-light-7);
  border-radius: 3px;
}

.cart-items::-webkit-scrollbar-track {
  background: var(--el-fill-color-light);
  border-radius: 3px;
}

.cart-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--el-bg-color);
  padding: 10px 12px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  transition: all 0.3s ease;
}

.cart-item:hover {
  box-shadow: 0 4px 12px rgba(var(--el-color-primary-rgb), 0.15);
  transform: translateX(4px);
}

.cart-item-enter-active,
.cart-item-leave-active {
  transition: all 0.3s ease;
}

.cart-item-enter-from,
.cart-item-leave-to {
  opacity: 0;
  transform: translateY(30px);
}

.cart-item-info {
  flex: 1;
  overflow: hidden;
}

.cart-item-name {
  margin: 0 0 4px 0;
  font-size: 14px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cart-item-price {
  color: var(--el-color-primary);
  font-size: 13px;
  font-weight: 500;
}

.cart-item-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 수량 조절 컨트롤 */
.quantity-control {
  display: flex;
  align-items: center;
  gap: 5px;
}

.quantity-btn {
  width: 24px !important;
  height: 24px !important;
  padding: 0 !important;
  font-size: 12px;
}

/* 마이너스 아이콘 스타일 */
.minus-icon {
  font-size: 14px;
  line-height: 1;
  font-weight: bold;
}

.quantity {
  display: inline-block;
  min-width: 24px;
  text-align: center;
  font-size: 14px;
  font-weight: 500;
}

.summary-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.cart-summary {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--el-bg-color);
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.cart-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
  padding: 10px 0;
  border-bottom: 1px solid var(--el-border-color-light);
  margin-bottom: 10px;
  font-size: 14px;
  flex-shrink: 0; /* 압축되지 않도록 */
}

.total-amount {
  color: var(--el-color-primary);
  font-size: 18px;
  font-weight: 700;
}

.cart-actions {
  display: flex;
  flex-direction: column;
  margin-top: auto;
  flex: 1;
}

/* 버튼 정렬 수정 - 세로 배치 */
.button-column {
  display: flex;
  flex-direction: column;
  gap: 10px;
  height: 100%;
  justify-content: flex-end; /* 버튼을 하단에 배치 */
}

.clear-cart-btn, 
.order-btn {
  height: 36px;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
}

.button-icon {
  margin-right: 5px;
}

.order-btn {
  background: linear-gradient(135deg, var(--el-color-primary) 0%, var(--el-color-primary-light-3) 100%);
  margin-left:0;
}

.order-btn:hover {
  background: linear-gradient(135deg, var(--el-color-primary-dark-2) 0%, var(--el-color-primary) 100%);
}

/* 반응형 스타일 */
@media (max-width: 1200px) {
  .menu-list {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 992px) {
  .menu-list {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .cart-row {
    flex-direction: column;
  }
  
  .summary-container {
    margin-top: 10px;
    height: auto;
  }
  
  .cart-summary {
    height: auto;
  }
  
  .cart-section {
    height: 300px; /* 모바일에서 높이 증가 */
  }
  
  .cart-items {
    max-height: 150px;
  }
  
  .button-column {
    flex-direction: column;
  }
}

@media (max-width: 768px) {
  .category-tabs :deep(.el-tabs__item) {
    height: 46px;
    line-height: 46px;
    font-size: 15px;
    padding: 0 16px;
  }
}

@media (max-width: 576px) {
  .menu-list {
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }
  
  .menu-section {
    padding: 10px;
  }
  
  .cart-section {
    padding: 10px;
  }
  
  .category-tabs :deep(.el-tabs__item) {
    height: 40px;
    line-height: 40px;
    font-size: 14px;
    padding: 0 12px;
  }
  
  .menu-item {
    height: 200px; /* 모바일에서 높이 줄임 */
  }
  
  .menu-image {
    height: 120px; /* 모바일에서 이미지 높이 줄임 */
  }
  
  .menu-content {
    padding: 8px;
  }
  
  .cart-item {
    padding: 8px 10px;
  }
}

@media (max-width: 480px) {
  .menu-list {
    grid-template-columns: 1fr;
  }
  
  .cart-item {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .cart-item-actions {
    width: 100%;
    margin-top: 8px;
    justify-content: space-between;
  }
}
</style>