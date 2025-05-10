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
const itemsPerPage = 8; // 한 페이지에 표시할 아이템 수

// 장바구니 높이 조절 관련 변수 추가
const cartSectionHeight = ref<number>(25); // 기본값 25vh
const isDragging = ref<boolean>(false);
const startY = ref<number>(0);
const startHeight = ref<number>(0);
const isResizeHandleVisible = ref<boolean>(true); // 드래그 바 표시 여부

// 장바구니 높이 관련 함수
const startResize = (e: MouseEvent) => {
  isDragging.value = true;
  startY.value = e.clientY;
  startHeight.value = cartSectionHeight.value;
  
  // 마우스 이벤트 리스너 추가
  document.addEventListener('mousemove', resizeCart);
  document.addEventListener('mouseup', stopResize);
};

const resizeCart = (e: MouseEvent) => {
  if (!isDragging.value) return;
  
  // 마우스 이동 거리 계산 (위로 이동하면 값이 커짐)
  const deltaY = startY.value - e.clientY;
  
  // vh 단위로 환산 (윈도우 높이의 1%가 1vh)
  const deltaVh = (deltaY / window.innerHeight) * 100;
  
  // 새 높이 계산
  let newHeight = startHeight.value + deltaVh;
  
  // 최소/최대 높이 제한
  newHeight = Math.max(15, Math.min(60, newHeight));
  
  cartSectionHeight.value = newHeight;
};

const stopResize = () => {
  isDragging.value = false;
  
  // 마우스 이벤트 리스너 제거
  document.removeEventListener('mousemove', resizeCart);
  document.removeEventListener('mouseup', stopResize);
  
  // 브라우저에 높이 저장
  localStorage.setItem('cartSectionHeight', cartSectionHeight.value.toString());
};

// 드래그 바 토글 함수
const toggleResizeHandle = () => {
  isResizeHandleVisible.value = !isResizeHandleVisible.value;
  // 설정 저장
  localStorage.setItem('isResizeHandleVisible', isResizeHandleVisible.value.toString());
};

// 브라우저에서 저장된 설정 로드
const loadSavedSettings = () => {
  // 높이 로드
  const savedHeight = localStorage.getItem('cartSectionHeight');
  if (savedHeight) {
    cartSectionHeight.value = parseFloat(savedHeight);
  }
  
  // 드래그 바 표시 설정 로드
  const savedHandleVisible = localStorage.getItem('isResizeHandleVisible');
  if (savedHandleVisible !== null) {
    isResizeHandleVisible.value = savedHandleVisible === 'true';
  }
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
    
    // 저장된 설정 로드
    loadSavedSettings();
    
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

      <!-- 장바구니 높이 조절 핸들 -->
      <div v-if="isResizeHandleVisible" class="cart-resize-handle" @mousedown="startResize">
        <div class="resize-indicator"></div>
        <button class="toggle-handle-btn" @click.stop="toggleResizeHandle" title="드래그 바 숨기기">
          <span class="material-icon">visibility_off</span>
        </button>
      </div>
      <div v-else class="show-handle-btn-container">
        <button class="show-handle-btn" @click="toggleResizeHandle" title="드래그 바 표시하기">
          <span class="material-icon">tune</span>
        </button>
      </div>

			<!-- 장바구니 영역 -->
			<section class="cart-section" :style="{ maxHeight: `${cartSectionHeight}vh` }">
				<div class="cart-row">
					<div class="cart-container">
						<h2 class="cart-title">장바구니</h2>
						
						<div v-if="cartItems.length === 0" class="empty-cart">
							장바구니가 비어있습니다.
						</div>
						
						<div v-else class="cart-items" :style="{ maxHeight: `calc(${cartSectionHeight}vh - 80px)` }">
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
						<div class="cart-summary" :style="{ height: `calc(${cartSectionHeight}vh - 54px)` }">
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
/* @import url('https://fonts.googleapis.com/css2?family=Jua&display=swap'); */
@import url('https://fonts.googleapis.com/css2?family=Gowun+Dodum&display=swap');

.order-view-container {
  height: 100%; /* Using 100% instead of 100vh */
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

/* 크롬 스타일 탭으로 변경 */
.category-tabs {
  display: flex;
  background: #dee1e6;
  min-height: 46px;
  flex-shrink: 0;
  padding: 20px 6px 0 6px;
  gap: 5px;
  position: relative;
  border-bottom: none;
}

.category-tabs button {
  position: relative;
  /* flex: 1 1 auto; */
  min-width: 200px;
  /* max-width: 500px; */
  flex:1;
  /* gap:10px; */
  padding: 10px 16px;
  font-size: 1.5rem;
  background: #bdc1c9;
  color: #5f6368;
  border: none;
  border-radius: 8px 8px 0 0;
  cursor: pointer;
  margin-right: -4px; /* 탭 간 약간 겹치는 효과 */
  transition: all 0.2s ease;
  z-index: 1;
}

/* .category-tabs button:hover {
  background: #d0d3d9;
} */

.category-tabs .active {
  background: var(--button-primary);
  color: #fff;
  /* text-color: var(--text-primary); */
  /* font-weight: 500; */
  z-index: 2; /* 활성 탭을 앞으로 가져오기 */
  padding-top: 12px; /* 활성 탭을 약간 더 높게 */
  margin-top: -2px;
  box-shadow: 0 -2px 6px rgba(0,0,0,0.1);
}

/* 활성 탭 아래에 흰색 선 추가 */
/* .category-tabs .active::after {
  content: "";
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 4px;
  background: #fff;
  z-index: 3;
} */

/* 탭 컨테이너 아래 구분선 */
.category-tabs::after {
  content: "";
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: #dadce0;
  z-index: 1;
}
.menu-list {
  padding: 12px; /* 패딩 줄임 */
  width: 100%;
  flex: 1; /* 남은 공간 차지 */
  overflow-y: auto; /* Enable vertical scrolling */
  display: flex;
  flex-direction: column;
}

.menu-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
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
  /* justify-content: flex-start; */
  max-width: 300px;
  margin: 0 auto;
  width: 100%;
  /* Calculate height based on available space and dynamic cart height */
  height: 280px;
  /* max-height: 340px; Maximum height limit */
  min-height: 180px; /* Minimum height */
  gap:10px;
}

/* 이미지 컨테이너 추가 */
.menu-image-container {
  width: 100%;
  height: 65%; /* Percentage of the menu item height */
  margin-bottom: 8px;
  overflow: hidden;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.menu-item img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.menu-name {
  font-weight: 600;
  font-size: 1.1rem; /* 폰트 크기 줄임 */
  margin: 2px 0; /* 마진 줄임 */
  line-height: 1.2;
  /* Ensure text doesn't overflow */
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  /* margin-bottom:10px; */
}

.menu-price {
  color: var(--button-primary);
  font-weight: bold;
  font-size: 1rem; /* 폰트 크기 줄임 */

  /* margin-top: auto; Push to bottom of flex container */
}

/* 높이 조절 핸들 스타일 */
.cart-resize-handle {
  height: 10px;
  width: 100%;
  background-color: #f0f0f0;
  cursor: ns-resize;
  display: flex;
  justify-content: center;
  align-items: center;
  border-top: 1px solid #ddd;
  border-bottom: 1px solid #ddd;
  position: relative;
  z-index: 11;
}

.resize-indicator {
  width: 40px;
  height: 4px;
  background-color: #ccc;
  border-radius: 2px;
}

.toggle-handle-btn {
  position: absolute;
  right: 10px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: none;
  background-color: #f0f0f0;
  color: #888;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.toggle-handle-btn:hover {
  background-color: #e0e0e0;
}

.show-handle-btn-container {
  position: relative;
  height: 0;
  z-index: 11;
  display: flex;
  justify-content: flex-end;
}

.show-handle-btn {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 1px solid #ddd;
  background-color: white;
  color: #888;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
  margin-top: -12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  transition: background-color 0.2s;
}

.show-handle-btn:hover {
  background-color: #f8f8f8;
}

/* 카트 영역 - 동적 높이 적용 */
.cart-section {
  background: #f8f9fa;
  border-top: 1px solid #eee;
  padding: 12px; /* 패딩 줄임 */
  max-height: v-bind('cartSectionHeight + "vh"'); /* 동적 높이 */
  overflow-y: auto; /* 필요시 스크롤 */
  transition: max-height 0.1s;
}

.cart-row {
  display: flex;
  flex-direction: row;
  gap: 16px; /* 간격 줄임 */
  height: 100%;
}

.cart-container {
  flex: 2;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.summary-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
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
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cart-items {
  display: flex;
  flex-direction: column;
  gap: 8px; /* 간격 줄임 */
  max-height: calc(v-bind('cartSectionHeight + "vh"') - 80px); /* 동적 높이에서 헤더 등 높이 빼기 */
  overflow-y: auto;
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE and Edge */
  flex: 1;
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
  height: calc(v-bind('cartSectionHeight + "vh"') - 54px); /* 동적 높이 */
  display: flex;
  flex-direction: column;
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
  margin-top: auto; /* 버튼이 항상 아래에 위치 */
  flex: 1; /* 남은 공간 차지 */
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
  flex: 1; /* 버튼이 각각 동일한 높이 비율로 공간 차지 */
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
    max-height: v-bind('cartSectionHeight + "vh"');
  }
  
  .menu-item {
    height: 180px; /* 태블릿에서의 고정 높이 */
  }
  
  .cart-summary {
    height: auto;
    min-height: 100px;
    margin-top: 10px;
  }
}

@media (max-width: 480px) {
  .menu-grid {
    grid-template-columns: repeat(1, 1fr);
  }
  
  .menu-item {
    height: 180px; /* 모바일에서의 고정 높이 */
    max-height: 280px;
  }
}

/* Add height-based media queries */
@media (max-height: 800px) {
  .menu-item {
    height: 200px; /* 화면 높이에 따른 고정 높이 */
  }
}

@media (max-height: 600px) {
  .menu-item {
    height: 150px; /* 작은 화면에서의 고정 높이 */
    min-height: 150px;
  }
  
  .category-tabs button {
    padding: 6px;
    font-size: 0.9rem;
  }
  
  .menu-name {
    font-size: 0.8rem;
  }
  
  .menu-price {
    font-size: 0.8rem;
  }
}
</style>