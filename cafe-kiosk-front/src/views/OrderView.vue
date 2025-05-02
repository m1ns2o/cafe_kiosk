<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import { CategoryAPI } from '../api/menu';
import type { MenuItem, Category, CartItem } from '../types/menuType';
import { useRouter } from 'vue-router'; // 라우터 추가

const router = useRouter(); // 라우터 인스턴스 생성

const menuItems = ref<MenuItem[]>([]);
const selectedCategory = ref<number>(0);
const categories = ref<Category[]>([]);
// const isLoading = ref<boolean>(true);
const error = ref<string | null>(null);
const cartItems = ref<CartItem[]>([]);

// 페이지네이션 관련 변수
const currentPage = ref(1);
const itemsPerPage = 6; // 한 페이지에 표시할 아이템 수

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

// 페이지 변경 함수
const changePage = (page: number) => {
  currentPage.value = page;
};

// 이전 페이지로 이동
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
};

// 다음 페이지로 이동
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
};

// 카테고리가 변경되면 해당 카테고리의 메뉴를 불러오는 함수
const loadMenuItems = async (categoryId: number) => {
  if (!categoryId) return;
  
  try {
    // isLoading.value = true;
    error.value = null;
    const menuResponse = await CategoryAPI.getMenus(categoryId);
    menuItems.value = menuResponse.data;
    console.log('메뉴 아이템 로드:', menuItems.value);
  } catch (err) {
    console.error('메뉴를 불러오는 중 오류가 발생했습니다:', err);
    error.value = '메뉴를 불러오는 중 오류가 발생했습니다.';
    menuItems.value = [];
  } finally {
    // isLoading.value = false;
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
    alert('장바구니가 비어있습니다.');
    return;
  }
  
  // 라우터로 결제 페이지로 이동하면서 장바구니 데이터와 총액 전달
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
    currentPage.value = 1; // 카테고리 변경 시 첫 페이지로 이동
  }
});

onMounted(async () => {
  try {
    // isLoading.value = true;
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
  } finally {
    // isLoading.value = false;
  }
});
</script>

<template>
	<div class="order-view-container">
		<div class="order-view">
			<!-- 오류 메시지 표시 -->
			<div v-if="error" class="error-message">
				{{ error }}
			</div>

			<!-- 로딩 표시 -->
			<!-- <div v-if="isLoading" class="loading">
				데이터를 불러오는 중...
			</div> -->

			<!-- 카테고리 탭 -->
			<nav class="category-tabs">
				<button 
          v-for="category in categories" 
          :key="category.id"
          :class="{ active: selectedCategory === category.id }"
          @click="selectedCategory = category.id"
        >
          {{ category.name }}
        </button>
			</nav>

			<!-- 메뉴 리스트 -->
			<section class="menu-list">
				<div v-if="menuItems.length === 0" class="empty-menu">
					이 카테고리에 메뉴가 없습니다.
				</div>
				<div v-else class="menu-grid">
					<div v-for="item in paginatedMenuItems" :key="item.id" class="menu-item" @click="addToCart(item)">
						<div class="menu-image-container">
							<img :src="`http://localhost:8080${item.image_url}`" :alt="item.name" />
						</div>
						<div class="menu-name">{{ item.name }}</div>
						<div class="menu-price">{{ item.price }}원</div>
					</div>
				</div>
				
				<!-- 페이지네이션 컨트롤 -->
				<div v-if="menuItems.length > itemsPerPage" class="pagination-controls">
					<button 
						class="pagination-btn" 
						:disabled="currentPage === 1" 
						@click="prevPage"
					>
						&lt;
					</button>
					
					<div class="pagination-pages">
						<button 
							v-for="page in totalPages" 
							:key="page" 
							:class="['page-btn', { active: currentPage === page }]"
							@click="changePage(page)"
						>
							{{ page }}
						</button>
					</div>
					
					<button 
						class="pagination-btn" 
						:disabled="currentPage === totalPages" 
						@click="nextPage"
					>
						&gt;
					</button>
				</div>
			</section>

			<!-- 장바구니 영역 -->
			<section class="cart-section">
				<div class="cart-row">
					<div class="cart-container">
						<h2 class="cart-title">장바구니</h2>
						
						<div v-if="cartItems.length === 0" class="empty-cart">
							장바구니가 비어있습니다.
						</div>
						
						<div v-else class="cart-items">
							<div v-for="(cartItem, index) in cartItems" :key="index" class="cart-item">
								<div class="cart-item-info">
									<div class="cart-item-name">{{ cartItem.item.name }}</div>
									<div class="cart-item-price">{{ cartItem.item.price }}원</div>
								</div>
								
								<div class="cart-item-actions">
									<div class="quantity-control">
										<button class="quantity-btn minus-btn" @click.stop="decreaseQuantity(index)">
											<span class="material-icon">remove</span>
										</button>
										<span class="quantity">{{ cartItem.quantity }}</span>
										<button class="quantity-btn plus-btn" @click.stop="increaseQuantity(index)">
											<span class="material-icon">add</span>
										</button>
									</div>
									<button class="remove-btn" @click.stop="removeFromCart(index)">
										<span class="material-icon">delete</span>
									</button>
								</div>
							</div>
						</div>
					</div>
					
					<div class="summary-container">
						<div class="cart-summary">
							<div class="cart-total">
								<span>총 금액:</span>
								<span class="total-amount">{{ totalAmount.toLocaleString() }}원</span>
							</div>
							<div class="cart-actions">
								<button class="clear-cart-btn" @click="clearCart">
									<span class="material-icon mr-1">delete_sweep</span>
									전체 삭제
								</button>
								<button class="order-btn" @click="placeOrder">
									<span class="material-icon mr-1">shopping_cart_checkout</span>
									주문하기
								</button>
							</div>
						</div>
					</div>
				</div>
			</section>

			<!-- 주문/결제 영역 -->
			
		</div>
	</div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/icon?family=Material+Icons');

.order-view-container {
  height: 100vh; /* 뷰포트 높이 사용 */
  width: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden; /* 전체 컨테이너 오버플로우 방지 */
}

.order-view {
  background: var(--background-primary);
  width: 100%;
  height: 100%;
  margin: 0 auto;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  font-family: 'Noto Sans KR', sans-serif;
  position: relative;
  display: flex;
  flex-direction: column;
  overflow: hidden; /* 자식 요소 오버플로우 방지 */
}

.category-tabs {
  display: flex;
  background: #f6f6f6;
  border-bottom: 1px solid #eee;
  min-height: 46px; /* 최소 높이 설정 */
}

.category-tabs button {
  flex: 1;
  padding: 8px; /* 패딩 줄임 */
  font-size: 1rem; /* 폰트 크기 줄임 */
  background: none;
  border: none;
  cursor: pointer;
}

.category-tabs .active {
  background: #fff;
  border-bottom: 2px solid var(--button-primary);
  color: var(--button-primary);
  font-weight: bold;
}

.menu-list {
  padding: 12px; /* 패딩 줄임 */
  width: 100%;
  flex: 1; /* 남은 공간 차지 */
  overflow-y: none; /* 수직 스크롤 추가 */
  display: flex;
  flex-direction: column;
}

.menu-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px; /* 간격 줄임 */
  width: 100%;
  margin-bottom: 10px; /* 마진 줄임 */
}

.menu-item {
  background: white;
  border-radius: 8px; /* 라운딩 줄임 */
  text-align: center;
  padding: 12px; /* 패딩 줄임 */
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
  cursor: pointer;
  transition: box-shadow 0.2s;
  display: flex;
  flex-direction: column;
  align-items: center;
  max-width: 300px;
  margin: 0 auto;
  width: 100%;
	/* height: auto; */
	max-height: 340px;
  min-height: 0; /* 높이 자동 조절 */
}

.menu-item:hover {
  box-shadow: 0 4px 12px rgba(255,111,65,0.12);
}

/* 이미지 컨테이너 추가 */
.menu-image-container {
  width: 100%;
  margin-bottom: 8px;
  aspect-ratio: 1 / 1;
  overflow: hidden;
  border-radius: 6px;
}

.menu-item img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.menu-name {
  font-size: 0.9rem; /* 폰트 크기 줄임 */
  margin: 2px 0; /* 마진 줄임 */
  line-height: 1.2;
}

.menu-price {
  color: var(--button-primary);
  font-weight: bold;
  font-size: 0.9rem; /* 폰트 크기 줄임 */
}

.cart-section {
  background: #f8f9fa;
  border-top: 1px solid #eee;
  padding: 12px; /* 패딩 줄임 */
  max-height: 25vh; /* 최대 높이 제한 */
  overflow-y: auto; /* 필요시 스크롤 */
}

.cart-row {
  display: flex;
  flex-direction: row;
  gap: 16px; /* 간격 줄임 */
}

.cart-container {
  flex: 2;
}

.summary-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
}

.cart-title {
  font-size: 1.1rem; /* 폰트 크기 줄임 */
  margin-bottom: 8px; /* 마진 줄임 */
  font-weight: bold;
}

.empty-cart {
  text-align: center;
  color: #888;
  padding: 8px; /* 패딩 줄임 */
  height: 100px; /* 높이 줄임 */
  display: flex;
  align-items: center;
  justify-content: center;
}

.cart-items {
  display: flex;
  flex-direction: column;
  gap: 8px; /* 간격 줄임 */
  max-height: 100px; /* 최대 높이 제한 */
  overflow-y: auto;
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE and Edge */
}

.cart-items::-webkit-scrollbar {
  display: none; /* Chrome, Safari, Opera */
}

.cart-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 8px; /* 패딩 줄임 */
  border-radius: 6px; /* 라운딩 줄임 */
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  min-height: 36px; /* 최소 높이 줄임 */
}

.cart-item:hover {
  box-shadow: 0 4px 12px rgba(255,111,65,0.12);
}

.cart-item-info {
  flex: 1;
}

.cart-item-name {
  font-weight: 500;
  font-size: 0.85rem; /* 폰트 크기 줄임 */
}

.cart-item-price {
  color: var(--button-primary);
  font-size: 0.8rem; /* 폰트 크기 줄임 */
}

.cart-item-actions {
  display: flex;
  align-items: center;
  gap: 8px; /* 간격 줄임 */
}

.quantity-control {
  display: flex;
  align-items: center;
  gap: 6px; /* 간격 줄임 */
}

.quantity-btn {
  width: 24px; /* 크기 줄임 */
  height: 24px; /* 크기 줄임 */
  border-radius: 50%;
  border: 1px solid #ddd;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--button-primary);
  transition: background-color 0.2s;
}

.quantity-btn:hover {
  background-color: #f0f0f0;
}

.quantity-btn.minus-btn {
  font-size: 16px; /* 폰트 크기 줄임 */
}

.quantity-btn.plus-btn {
  font-size: 16px; /* 폰트 크기 줄임 */
}

.material-icon {
  font-family: 'Material Icons';
  font-weight: normal;
  font-style: normal;
  font-size: 18px; /* 폰트 크기 줄임 */
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
  -moz-osx-font-smoothing: grayscale;
}

.quantity {
  min-width: 18px; /* 크기 줄임 */
  text-align: center;
  font-size: 0.85rem; /* 폰트 크기 줄임 */
}

.remove-btn {
  background: transparent;
  border: none;
  color: var(--button-primary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s;
}

.remove-btn:hover {
  color: #ff5252;
}

.cart-summary {
  background: white;
  border-radius: 6px; /* 라운딩 줄임 */
  padding: 12px; /* 패딩 줄임 */
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.cart-total {
  display: flex;
  justify-content: space-between;
  font-weight: bold;
  padding: 6px 0; /* 패딩 줄임 */
  border-bottom: 1px solid #eee;
  margin-bottom: 8px; /* 마진 줄임 */
  font-size: 0.9rem; /* 폰트 크기 줄임 */
}

.total-amount {
  color: var(--button-primary);
  font-size: 1rem; /* 폰트 크기 줄임 */
}

.cart-actions {
  display: flex;
  flex-direction: column;
  gap: 6px; /* 간격 줄임 */
}

.clear-cart-btn, .order-btn {
  padding: 8px 12px; /* 패딩 줄임 */
  border-radius: 6px; /* 라운딩 줄임 */
  border: none;
  cursor: pointer;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem; /* 폰트 크기 줄임 */
}

.mr-1 {
  margin-right: 3px; /* 마진 줄임 */
}

.clear-cart-btn {
  background: #f1f1f1;
  color: #555;
}

.order-btn {
  background: var(--button-primary);
  color: white;
}

.order-btn:hover {
  opacity: 0.9;
}

/* 페이지네이션 스타일 */
.pagination-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 10px; /* 마진 줄임 */
  gap: 8px; /* 간격 줄임 */
  padding: 8px 0; /* 상하 패딩 */
  background: var(--background-primary);
  position: sticky;
  bottom: 0;
  z-index: 10;
}

.pagination-btn {
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 4px 8px; /* 패딩 줄임 */
  cursor: pointer;
  font-size: 0.9rem; /* 폰트 크기 줄임 */
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-pages {
  display: flex;
  gap: 4px; /* 간격 줄임 */
}

.page-btn {
  width: 26px; /* 크기 줄임 */
  height: 26px; /* 크기 줄임 */
  border-radius: 4px;
  border: 1px solid #ddd;
  background: white;
  cursor: pointer;
  font-size: 0.85rem; /* 폰트 크기 줄임 */
}

.page-btn.active {
  background: var(--button-primary);
  color: white;
  border-color: var(--button-primary);
}

/* 반응형 레이아웃 */
@media (max-width: 768px) {
  .menu-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .cart-row {
    flex-direction: column;
  }
  
  .cart-section {
    max-height: 40vh;
  }
}

@media (max-width: 480px) {
  .menu-grid {
    grid-template-columns: repeat(1, 1fr);
  }
}
</style>