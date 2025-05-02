<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { CategoryAPI } from '../api/menu';
import type { MenuItem, Category } from '../types/menuType';

const menuItems = ref<MenuItem[]>([]);
const selectedCategory = ref<number>(0);
const categories = ref<Category[]>([]);
const isLoading = ref<boolean>(true);
const error = ref<string | null>(null);

// 카테고리가 변경되면 해당 카테고리의 메뉴를 불러오는 함수
const loadMenuItems = async (categoryId: number) => {
  if (!categoryId) return;
  
  try {
    isLoading.value = true;
    error.value = null;
    const menuResponse = await CategoryAPI.getMenus(categoryId);
    menuItems.value = menuResponse.data;
    console.log('메뉴 아이템 로드:', menuItems.value);
  } catch (err) {
    console.error('메뉴를 불러오는 중 오류가 발생했습니다:', err);
    error.value = '메뉴를 불러오는 중 오류가 발생했습니다.';
    menuItems.value = [];
  } finally {
    isLoading.value = false;
  }
};

// 카테고리가 변경될 때마다 메뉴 아이템 업데이트
watch(() => selectedCategory.value, async (newCategoryId) => {
  if (newCategoryId) {
    await loadMenuItems(newCategoryId);
  }
});

onMounted(async () => {
  try {
    isLoading.value = true;
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
    isLoading.value = false;
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
			<div v-if="isLoading" class="loading">
				데이터를 불러오는 중...
			</div>

			<!-- 카테고리 탭 -->
			<nav v-if="!isLoading && categories.length > 0" class="category-tabs">
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
			<section v-if="!isLoading && !error" class="menu-list">
				<div v-if="menuItems.length === 0" class="empty-menu">
					이 카테고리에 메뉴가 없습니다.
				</div>
				<div v-else v-for="item in menuItems" :key="item.id" class="menu-item">
					<img :src="item.image_url || '/placeholder.png'" :alt="item.name" />
					<div class="menu-name">{{ item.name }}</div>
					<div class="menu-price">{{ item.price }}원</div>
				</div>
			</section>

			<!-- 주문/결제 영역 -->
			<aside class="order-panel">
				<button class="clear-btn">전체삭제</button>
				<div class="selected-title">선택한 상품</div>
				<div class="selected-items">
					<!-- 선택된 상품 리스트 -->
				</div>
				<button class="pay-btn">결제하기</button>
			</aside>
		</div>
	</div>
</template>

<style scoped>
.order-view-container {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
}

.order-view {
	background: #fff;
	width: 100%;
	height: 100%;
	/* max-width: 1200px; */
	margin: 0 auto;
	border-radius: 10px;
	overflow: hidden;
	box-shadow: 0 2px 8px rgba(0,0,0,0.1);
	font-family: 'Noto Sans KR', sans-serif;
	position: relative;
	display: flex;
	flex-direction: column;
}

.header {
	background: #ff6f41;
	color: #fff;
	display: flex;
	align-items: center;
	padding: 16px;
	font-size: 1.5rem;
	justify-content: space-between;
}

.category-tabs {
	display: flex;
	background: #f6f6f6;
	border-bottom: 1px solid #eee;
}

.category-tabs button {
	flex: 1;
	padding: 12px;
	font-size: 1.1rem;
	background: none;
	border: none;
	cursor: pointer;
}

.category-tabs .active {
	background: #fff;
	border-bottom: 2px solid #ff6f41;
	color: #ff6f41;
	font-weight: bold;
}

.menu-list {
	display: flex;
	flex-wrap: wrap;
	gap: 16px;
	padding: 16px;
	justify-content: space-between;
}

.menu-item {
	width: 180px;
	background: #fafafa;
	border-radius: 8px;
	text-align: center;
	padding: 12px;
	box-shadow: 0 1px 4px rgba(0,0,0,0.06);
	cursor: pointer;
	transition: box-shadow 0.2s;
}

.menu-item:hover {
	box-shadow: 0 4px 12px rgba(255,111,65,0.12);
}

.menu-item img {
	width: 100px;
	height: 120px;
	object-fit: contain;
	margin-bottom: 8px;
}

.menu-name {
	font-size: 1rem;
	margin: 4px 0;
}

.menu-price {
	color: #ff6f41;
	font-weight: bold;
}

.order-panel {
	position: absolute;
	right: 16px;
	top: 120px;
	width: 220px;
	background: #f9f9f9;
	border-radius: 8px;
	padding: 16px;
	box-shadow: 0 2px 8px rgba(0,0,0,0.05);
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.timer {
	color: #333;
	font-size: 1.1rem;
}

.time {
	color: #ff6f41;
	font-weight: bold;
	font-size: 1.2rem;
}

.clear-btn, .pay-btn {
	background: #444;
	color: #fff;
	border: none;
	border-radius: 6px;
	padding: 10px;
	font-size: 1rem;
	cursor: pointer;
}

.pay-btn {
	background: #ff6f41;
	margin-top: 12px;
}

.selected-title {
	margin-top: 12px;
	font-weight: bold;
}

.selected-items {
	min-height: 60px;
	background: #fff;
	border-radius: 6px;
	padding: 8px;
	margin-bottom: 8px;
}

.loading {
	display: flex;
	justify-content: center;
	align-items: center;
	padding: 2rem;
	font-size: 1.1rem;
	color: #666;
}

.error-message {
	background-color: #ffebee;
	color: #c62828;
	padding: 1rem;
	margin: 1rem;
	border-radius: 4px;
	text-align: center;
}

.empty-menu {
	width: 100%;
	padding: 2rem;
	text-align: center;
	color: #666;
	font-size: 1.1rem;
}
</style>