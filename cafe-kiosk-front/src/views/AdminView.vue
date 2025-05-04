<template>
	<div class="sidebar-container">
		<!-- 사이드바 메뉴 -->
		<el-menu
			default-active="1"
			class="el-menu-vertical-demo"
			:collapse="isCollapse"
			@open="handleOpen"
			@close="handleClose"
			text-color="var(--nordic-icon-gray)"
			active-text-color="var(--accent-color)"
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

			<!-- 기존 메뉴 아이템들 -->
			<el-menu-item index="1">
				<el-icon><PieChart /></el-icon>
				<template #title>Statistics</template>
			</el-menu-item>
			<el-menu-item index="2">
				<el-icon><list /></el-icon>
				<template #title>Category</template>
			</el-menu-item>
			<el-menu-item index="3">
				<el-icon><Coffee /></el-icon>
				<template #title>Menu</template>
			</el-menu-item>
			<el-menu-item index="4">
				<el-icon><Wallet /></el-icon>
				<template #title>Order</template>
			</el-menu-item>
		</el-menu>
	</div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import {
	Menu as IconMenu,
	DArrowRight,
	DArrowLeft,
	Coffee,
	PieChart,
	List,
	Wallet,
} from "@element-plus/icons-vue";

const isCollapse = ref(true);
const handleOpen = (key: string, keyPath: string[]) => {
	console.log(key, keyPath);
};
const handleClose = (key: string, keyPath: string[]) => {
	console.log(key, keyPath);
};
</script>

<style>
:root {
	/* 추가된 적절한 그레이 색상 */
	--nordic-icon-gray: #94a3b8; /* 중간 톤의 슬레이트 그레이 */
	--nordic-icon-light-gray: #cbd5e1; /* 밝은 톤의 슬레이트 그레이 */
}

.sidebar-container {
	display: flex;
	flex-direction: column;
	height: 100%;
	border-right: 1px solid var(--border-color);
}

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
	color: var(--accent-color);
	border-radius: 6px;
	cursor: pointer;
	transition: background-color 0.3s;
}

.toggle-icon-wrapper:hover {
	background-color: var(--nordic-light-blue);
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
	color: var(--nordic-icon-gray);
}

/* 활성화된 메뉴 아이템의 아이콘 색상 */
.el-menu-item.is-active .el-icon {
	color: var(--accent-color);
}
</style>
