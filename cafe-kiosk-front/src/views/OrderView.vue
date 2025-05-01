<script setup lang="ts">
import { ref } from 'vue'

const activeTab = ref('hot')
const remainingTime = ref(300)
const menuItems = ref([
  {
    id: 1,
    name: '아메리카노',
    price: 3000,
    image: 'https://via.placeholder.com/100x120'
  },
  {
    id: 2,
    name: '카페라떼',
    price: 4000,
    image: 'https://via.placeholder.com/100x120'
  },
  {
    id: 3,
    name: '카푸치노',
    price: 4500,
    image: 'https://via.placeholder.com/100x120'
  }
])
</script>

<template>
	<div class="order-view-container">
		<div class="order-view">
			<!-- 상단 바 -->
			<header class="header">
				<span class="icon home"></span>
				<span class="title">Easy KIOSK</span>
				<span class="icon settings"></span>
			</header>

			<!-- 카테고리 탭 -->
			<nav class="category-tabs">
				<button :class="{ active: activeTab === 'season' }" @click="activeTab = 'season'">시즌 메뉴</button>
				<button :class="{ active: activeTab === 'hot' }" @click="activeTab = 'hot'">커피(HOT)</button>
				<button :class="{ active: activeTab === 'ice' }" @click="activeTab = 'ice'">커피(ICE)</button>
			</nav>

			<!-- 메뉴 리스트 -->
			<section class="menu-list">
				<div v-for="item in menuItems" :key="item.id" class="menu-item">
					<img :src="item.image" :alt="item.name" />
					<div class="menu-name">{{ item.name }}</div>
					<div class="menu-price">{{ item.price }}원</div>
				</div>
			</section>

			<!-- 주문/결제 영역 -->
			<aside class="order-panel">
				<div class="timer">
					남은시간 <span class="time">{{ remainingTime }}</span>초
				</div>
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
</style>