<template>
  <div class="category-container">
    <div class="category-content">
      <el-card shadow="hover" class="category-card">
        <template #header>
          <div class="card-header">
            <span class="header-title">카테고리 관리</span>
            <!-- 필요한 경우 여기에 추가 버튼 등을 배치할 수 있습니다 -->
          </div>
        </template>
        
        <el-table 
          ref="tableRef" 
          row-key="id" 
          :data="pagedTableData" 
          header-cell-class-name="table-header"
          v-loading="loading"
          class="category-table"
          :header-cell-style="{ fontSize: '16px', fontWeight: 'bold', background: '#f5f7fa', height: '50px' }"
          :cell-style="{ fontSize: '15px', height: '60px' }"
        >
          <!-- 카테고리 이름 열 -->
          <el-table-column prop="name" label="카테고리" class-name="name-column" align="left" />
          
          <!-- 작업 열 -->
          <el-table-column label="" class-name="action-column" align="right">
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
                  @click="confirmDeleteCategory(scope.row)"
                >
                  <Delete />
                </el-icon>
              </div>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 페이지네이션 추가 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            :page-size="pageSize"
            layout="total, prev, pager, next"
            :total="filteredTableData.length"
            @current-change="handleCurrentChange"
            :pager-count="5"
            background
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { Edit, Delete } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import type { TableInstance } from 'element-plus'
import { CategoryAPI } from '../../api/menu'
import type { Category } from '../../types/menuType'

const tableRef = ref<TableInstance>()
const tableData = ref<Category[]>([])
const loading = ref(false)

// 페이지네이션 관련 변수
const currentPage = ref(1)
const pageSize = ref(8) // 페이지당 8개 항목

// 필터링된 테이블 데이터 
const filteredTableData = computed(() => tableData.value)

// 현재 페이지에 표시할 데이터
const pagedTableData = computed(() => {
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return filteredTableData.value.slice(startIndex, endIndex)
})

// 컴포넌트 마운트 시 데이터 로드
onMounted(async () => {
  await fetchData()
})

// API에서 데이터 가져오기
const fetchData = async () => {
  loading.value = true;
  try {
    const categoryResponse = await CategoryAPI.getAllCategories();
    tableData.value = categoryResponse.data;
  } catch (error) {
    console.error('카테고리 데이터 로딩 오류:', error);
    ElMessage.error('카테고리 데이터를 불러오는 중 오류가 발생했습니다.');
  } finally {
    loading.value = false;
  }
}

// 페이지 변경 처리
const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

// 수정 핸들러
const handleEdit = (index: number, row: Category) => {
  console.log('수정:', index, row)
  // 여기에서 수정 기능을 구현하거나 부모 컴포넌트에 이벤트를 발생시킵니다
}

// 삭제 확인 다이얼로그
const confirmDeleteCategory = (row: Category) => {
  ElMessageBox.confirm(
    `"${row.name}" 카테고리를 삭제하시겠습니까?`,
    '삭제 확인',
    {
      confirmButtonText: '삭제',
      cancelButtonText: '취소',
      type: 'warning',
    }
  )
    .then(() => {
      deleteCategory(row.id)
    })
    .catch(() => {
      // 사용자가 취소를 클릭한 경우, 아무 작업도 수행하지 않음
    })
}

// 카테고리 항목 삭제
const deleteCategory = async (id: number) => {
  loading.value = true
  try {
    await CategoryAPI.deleteCategory(id)
    ElMessage.success('카테고리가 삭제되었습니다.')
    // 카테고리 목록 새로고침
    await fetchData()
  } catch (error) {
    console.error('카테고리 삭제 오류:', error)
    ElMessage.error('카테고리 삭제 중 오류가 발생했습니다.')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.category-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  min-height: 100vh;
  padding: 20px;
}

.category-content {
  width: 100%;
  max-width: 1200px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.category-card {
  width: 95%;
  height: 100%;
  margin: 0 auto;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px 0;
}

.header-title {
  font-weight: bold;
  font-size: 18px;
}

.category-table {
  width: 100%;
  margin-bottom: 10px;
}

.pagination-container {
  margin-top: 20px;
  margin-bottom: 10px;
  display: flex;
  justify-content: center;
}

/* 테이블 헤더와 셀 정렬 */
.table-header {
  text-align: left !important;
  font-weight: bold;
  background-color: #f5f7fa !important;
}

:deep(.el-table__header-wrapper .el-table__header th) {
  text-align: left !important;
}

:deep(.el-table__row) {
  height: 60px;
  font-size: 14px;
}



/* 개선된 컬럼 너비 설정 */
:deep(.name-column) {
  width: calc(100% - 80px) !important; 
}

:deep(.action-column) {
  width: 80px !important;
}

:deep(.el-table__row .name-column) {
  text-align: left !important;
}

:deep(.el-table__row .action-column) {
  text-align: right !important;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  width: 100%;
}

.action-icon {
  font-size: 18px;
  cursor: pointer;
  transition: transform 0.2s, color 0.2s;
}

.action-icon:hover {
  transform: scale(1.2);
}

.edit-icon {
  color: #a0a0a0; /* --nordic-icon-light-gray */
}

.edit-icon:hover {
  color: #409eff; /* Element Plus primary color */
}

.delete-icon {
  color: #f56c6c; /* Element Plus danger color */
}

.delete-icon:hover {
  color: #e41e1e; /* Darker red on hover */
}

/* 미디어 쿼리 */
@media screen and (max-width: 768px) {
  .category-container {
    padding: 10px;
  }
  
  .category-card {
    width: 100%;
  }
  
  :deep(.action-column) {
    width: 70px !important;
  }
}

@media screen and (max-width: 576px) {
  .category-container {
    padding: 10px 5px;
  }
  
  .header-title {
    font-size: 16px;
  }
  
  :deep(.action-column) {
    width: 60px !important;
  }
  
  :deep(.el-table__row) {
    height: 50px;
    font-size: 13px;
  }
}
</style>