<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { CategoryAPI } from '../api/menu';
import type { MenuItem, Category } from '../types/menuType';

const menuItems = ref<MenuItem[]>([]);
const selectedCategory = ref<number>(0);
const categories = ref<Category[]>([]);
// const isLoading = ref<boolean>(true);
const error = ref<string | null>(null);

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

// 카테고리가 변경될 때마다 메뉴 아이템 업데이트
watch(() => selectedCategory.value, async (newCategoryId) => {
  if (newCategoryId) {
    await loadMenuItems(newCategoryId);
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
					<div v-for="item in menuItems" :key="item.id" class="menu-item">
						<img :src="`http://localhost:8080${item.image_url}`" :alt="item.name" />
						<div class="menu-name">{{ item.name }}</div>
						<div class="menu-price">{{ item.price }}원</div>
					</div>
				</div>
			</section>

			<!-- 주문/결제 영역 -->
			
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
	background: var(--background-primary);
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
	border-bottom: 2px solid var(--button-primary);
	color: var(--button-primary);
	font-weight: bold;
}

.menu-list {
	padding: 16px;
	width: 100%;
}

.menu-grid {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: 16px;
	width: 100%;
}

.menu-item {
	background: white;
	border-radius: 12px;
	text-align: center;
	padding: 20px;
	box-shadow: 0 1px 4px rgba(0,0,0,0.06);
	cursor: pointer;
	transition: box-shadow 0.2s;
	display: flex;
	flex-direction: column;
	align-items: center;
}

.menu-item:hover {
	box-shadow: 0 4px 12px rgba(255,111,65,0.12);
}

.menu-item img {
	margin-bottom: 8px;
	width: 100%; 
	aspect-ratio: 1 / 1; 
	object-fit: contain; 
}

.menu-name {
	font-size: 1rem;
	margin: 4px 0;
}

.menu-price {
	color: var(--button-primary);
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

.clear-btn, .pay-btn {
	background: var(--button-primary);
	color: #fff;
	border: none;
	border-radius: 6px;
	padding: 10px;
	font-size: 1rem;
	cursor: pointer;
}

.pay-btn {
	background: var(--button-primary);
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

/* Responsive layout */
@media (max-width: 768px) {
  .menu-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .menu-grid {
    grid-template-columns: repeat(1, 1fr);
  }
}
</style>