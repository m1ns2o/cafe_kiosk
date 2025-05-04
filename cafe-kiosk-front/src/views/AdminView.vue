<template>
	<div class="sidebar-container">
		<!-- 사이드바 메뉴 -->
		<el-menu
			default-active="2"
			class="el-menu-vertical-demo"
			:collapse="isCollapse"
			@open="handleOpen"
			@close="handleClose"
			text-color="var(--text-primary)"
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
			<el-menu-item index="2">
				<el-icon><icon-menu /></el-icon>
				<template #title>Navigator Two</template>
			</el-menu-item>
			<el-menu-item index="3">
				<el-icon><document /></el-icon>
				<template #title>Navigator Three</template>
			</el-menu-item>
			<el-menu-item index="4">
				<el-icon><setting /></el-icon>
				<template #title>Navigator Four</template>
			</el-menu-item>
		</el-menu>
	</div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import {
	Document,
	Menu as IconMenu,
	Setting,
	DArrowRight,
	DArrowLeft,
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
	color: var(--text-primary);
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

.el-menu-item .el-icon {
	color: var(--text-primary);
}

.el-menu-item.is-active .el-icon {
	color: var(--accent-color);
}
</style>
