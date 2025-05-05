<template>
  <div>
    <el-table 
      ref="tableRef" 
      row-key="id" 
      :data="tableData" 
      style="width: 100%" 
      header-cell-class-name="table-header"
      v-loading="loading"
    >
      <!-- 이미지 열 -->
      <el-table-column label="이미지" width="100">
        <template #default="scope">
          <div class="image-container">
            <img v-if="scope.row.image_url" 
                 :src="scope.row.image_url" 
                 :alt="scope.row.name" 
                 class="menu-image" />
            <el-icon v-else><picture-rounded /></el-icon>
          </div>
        </template>
      </el-table-column>
      
      <!-- 메뉴 이름 열 -->
      <el-table-column prop="name" label="메뉴" width="180">
        <template #default="scope">
          <el-popover effect="light" trigger="hover" placement="top" width="auto">
            <template #default>
              <div>이름: {{ scope.row.name }}</div>
              <div>가격: {{ formatPrice(scope.row.price) }}</div>
              <div>카테고리: {{ getCategoryName(scope.row.category_id) }}</div>
            </template>
            <template #reference>
              <span>{{ scope.row.name }}</span>
            </template>
          </el-popover>
        </template>
      </el-table-column>
      
      <!-- 가격 열 -->
      <el-table-column prop="price" label="가격" width="120">
        <template #default="scope">
          <span>{{ formatPrice(scope.row.price) }}</span>
        </template>
      </el-table-column>
      
      <!-- 카테고리 열 -->
      <el-table-column
        prop="category_id"
        label="카테고리"
        width="150"
        :filters="categoryFilters"
        :filter-method="filterByCategory"
        filter-placement="bottom-end"
      >
        <template #default="scope">
          <span>{{ getCategoryName(scope.row.category_id) }}</span>
        </template>
      </el-table-column>
      
      <!-- 작업 열 -->
      <el-table-column label="" width="120">
        <template #default="scope">
          <div class="action-buttons">
            <el-icon 
              class="action-icon edit-icon" 
              @click="handleEdit(scope.$index, scope.row)"
            >
              <Edit />
            </el-icon>
            <el-icon
              class="action-icon delete-icon"
              @click="confirmDelete(scope.row)"
            >
              <Delete />
            </el-icon>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { PictureRounded, Edit, Delete } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import type { TableInstance } from 'element-plus'
import { MenuAPI, CategoryAPI } from '../../api/menu'
import type { MenuItem, Category } from '../../types/menuType'

const tableRef = ref<TableInstance>()
const tableData = ref<MenuItem[]>([])
const categories = ref<Category[]>([])
const loading = ref(false)

// 컴포넌트 마운트 시 데이터 로드
onMounted(async () => {
  await fetchData()
})

// API에서 데이터 가져오기
const fetchData = async () => {
  loading.value = true
  try {
    const [menuResponse, categoryResponse] = await Promise.all([
      MenuAPI.getMenus(),
      CategoryAPI.getAllCategories()
    ])
    tableData.value = menuResponse.data
    categories.value = categoryResponse.data
  } catch (error) {
    console.error('데이터 로딩 오류:', error)
    ElMessage.error('메뉴 데이터를 불러오는 중 오류가 발생했습니다.')
  } finally {
    loading.value = false
  }
}

// 카테고리 필터 옵션 생성
const categoryFilters = computed(() => {
  return categories.value.map(category => ({
    text: category.name,
    value: category.id
  }))
})

// 카테고리 이름 가져오기
const getCategoryName = (categoryId: number): string => {
  const category = categories.value.find(c => c.id === categoryId)
  return category ? category.name : '미분류'
}

// 카테고리 필터링
const filterByCategory = (value: number, row: MenuItem) => {
  return row.category_id === value
}

// 가격 포맷팅 (원화)
const formatPrice = (price: number): string => {
  return new Intl.NumberFormat('ko-KR', { 
    style: 'currency', 
    currency: 'KRW',
    minimumFractionDigits: 0
  }).format(price)
}

// 수정 핸들러
const handleEdit = (index: number, row: MenuItem) => {
  console.log('수정:', index, row)
  // 여기에서 수정 기능을 구현하거나 부모 컴포넌트에 이벤트를 발생시킵니다
}

// 삭제 확인 다이얼로그
const confirmDelete = (row: MenuItem) => {
  ElMessageBox.confirm(
    `"${row.name}" 메뉴를 삭제하시겠습니까?`,
    '삭제 확인',
    {
      confirmButtonText: '삭제',
      cancelButtonText: '취소',
      type: 'warning',
    }
  )
    .then(() => {
      deleteMenu(row.id)
    })
    .catch(() => {
      // 사용자가 취소를 클릭한 경우, 아무 작업도 수행하지 않음
    })
}

// 메뉴 항목 삭제
const deleteMenu = async (id: number) => {
  loading.value = true
  try {
    await MenuAPI.deleteMenu(id)
    ElMessage.success('메뉴가 삭제되었습니다.')
    // 메뉴 목록 새로고침
    await fetchData()
  } catch (error) {
    console.error('메뉴 삭제 오류:', error)
    ElMessage.error('메뉴 삭제 중 오류가 발생했습니다.')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.controls {
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-end;
}

.image-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50px;
}

.menu-image {
  max-height: 50px;
  max-width: 50px;
  object-fit: contain;
}

.table-header {
  text-align: center !important;
  font-weight: bold;
}

.action-buttons {
  display: flex;
  justify-content: flex-start;
  gap: 16px;
}

.action-icon {
  font-size: 18px;
  cursor: pointer;
}

.edit-icon {
  color: #a0a0a0; /* --nordic-icon-light-gray */
}

.delete-icon {
  color: #f56c6c; /* Element Plus danger color */
}
</style>