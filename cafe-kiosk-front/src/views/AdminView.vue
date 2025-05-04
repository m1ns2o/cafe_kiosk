<script lang="ts" setup>
import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import {
	DArrowRight,
	DArrowLeft,
	Coffee,
	PieChart,
	List,
	Wallet,
} from "@element-plus/icons-vue";

const isCollapse = ref(true);
const route = useRoute();

// 현재 라우트 경로를 기반으로 활성화된 메뉴 아이템 설정
const activeIndex = computed(() => {
	return route.path;
});

const handleOpen = (key: string, keyPath: string[]) => {
	console.log(key, keyPath);
};

const handleClose = (key: string, keyPath: string[]) => {
	console.log(key, keyPath);
};

// 브라우저 창 크기 변경 시 자동으로 사이드바 상태 조정
const handleResize = () => {
	isCollapse.value = window.innerWidth < 768;
};

onMounted(() => {
	// 초기 로드 시 화면 크기에 따라 사이드바 상태 설정
	handleResize();
	// 화면 크기 변경 이벤트 리스너 추가
	// window.addEventListener("resize", handleResize);
});
</script>

<template>
	<div class="layout-container" :class="{ 'sidebar-collapsed': isCollapse }">
		<div class="sidebar-container">
			<!-- 사이드바 메뉴 -->
			<el-menu
				:default-active="activeIndex"
				class="el-menu-vertical-demo"
				:collapse="isCollapse"
				@open="handleOpen"
				@close="handleClose"
				text-color="var(--nordic-icon-light-gray)"
				active-text-color="var(--accent-color)"
				router
			>
				<!-- 맨 위에 토글 버튼 추가 -->
				<div class="toggle-item">
					<div
						class="toggle-content"
						:class="{ 'toggle-content-collapsed': isCollapse }"
					>
						<div
							class="toggle-icon-wrapper"
							@click.stop="isCollapse = !isCollapse"
						>
							<el-icon :size="20">
								<d-arrow-right v-if="isCollapse" />
								<d-arrow-left v-else />
							</el-icon>
						</div>
					</div>
				</div>

				<!-- 라우터 링크가 적용된 메뉴 아이템들 -->
				<el-menu-item index="/admin/statistics">
					<el-icon><PieChart /></el-icon>
					<template #title>Statistics</template>
				</el-menu-item>
				<el-menu-item index="/admin/category">
					<el-icon><list /></el-icon>
					<template #title>Category</template>
				</el-menu-item>
				<el-menu-item index="/admin/menu">
					<el-icon><Coffee /></el-icon>
					<template #title>Menu</template>
				</el-menu-item>
				<el-menu-item index="/admin/order">
					<el-icon><Wallet /></el-icon>
					<template #title>Order</template>
				</el-menu-item>
			</el-menu>
		</div>

		<!-- 내부 라우터 뷰 - 자식 라우트 컴포넌트가 여기에 렌더링됨 -->
		<div class="content-container">
			<router-view />
		</div>
	</div>
</template>

<style>
:root {
	--nordic-icon-gray: #94a3b8; /* 중간 톤의 슬레이트 그레이 */
	--nordic-icon-light-gray: #a0a0a0; /* 밝은 톤의 슬레이트 그레이 */
	/* --border-color: #e2e8f0;
  --nordic-light-blue: #e0f2fe;
  --accent-color: #0ea5e9; */
}

.layout-container {
	display: flex;
	height: 100vh;
	width: 100%;
	transition: all 0.3s ease;
}

.sidebar-container {
	display: flex;
	flex-direction: column;
	height: 100%;
	border-right: 1px solid var(--border-color);
	transition: all 0.3s ease;
}

.content-container {
	flex: 1;
	/* padding: 20px; */
	overflow-y: auto;
	transition: margin-left 0.3s ease;
}

/* 사이드바가 접혔을 때 컨텐츠 영역 조정 */
/* .sidebar-collapsed .content-container {
	margin-left: 64px; 
} */

/* 사이드바가 펼쳐졌을 때 컨텐츠 영역 조정 */
/* .layout-container:not(.sidebar-collapsed) .content-container {
	margin-left: 200px; /
} */

.toggle-item {
	height: 56px;
	display: flex;
	align-items: center;
	padding: 0;
	color: var(--nordic-icon-gray);
	position: relative;
}

.toggle-content {
	display: flex;
	width: 100%;
	justify-content: flex-end;
	padding-right: 20px;
}

/* 접힌 상태일 때 중앙 정렬 */
.toggle-content-collapsed {
	justify-content: center;
	padding-right: 0;
}

.toggle-icon-wrapper {
	display: flex;
	justify-content: center;
	align-items: center;
	width: 32px;
	height: 32px;
	color: var(--nordic-icon-light-gray);
	border-radius: 6px;
	cursor: pointer;
	transition: background-color 0.3s;
}

.toggle-icon-wrapper:hover {
	background-color: var(--nordic-light-blue);
	color: var(--accent-color);
}

.el-menu-vertical-demo:not(.el-menu--collapse) {
	width: 200px;
}

.el-menu-vertical-demo {
	height: 100%;
}

.el-menu-item:hover {
	background-color: var(--nordic-light-blue) !important;
}

/* 선택되지 않은 메뉴 아이템의 아이콘 색상 변경 */
.el-menu-item .el-icon {
	color: var(--nordic-icon-light-gray);
}

/* 활성화된 메뉴 아이템의 아이콘 색상 */
.el-menu-item.is-active .el-icon {
	color: var(--accent-color);
}

/* 반응형 스타일 */
/* @media (max-width: 768px) {
	.layout-container:not(.sidebar-collapsed) .content-container {
		margin-left: 64px; 
	}
} */
</style>
