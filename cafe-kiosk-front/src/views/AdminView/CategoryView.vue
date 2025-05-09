<template>
  <div class="category-container">
    <div class="category-content">
      <el-card shadow="hover" class="category-card">
        <template #header>
          <div class="card-header">
            <span class="header-title">카테고리 관리</span>
            <!-- 추가 버튼 -->
            <el-button 
              type="primary" 
              size="small" 
              @click="handleAdd"
              class="add-button"
            >
              <el-icon class="el-icon--left" style="color: white;"><Plus /></el-icon>
              항목 추가
            </el-button>
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

  <!-- 카테고리 추가/편집 다이얼로그 -->
  <el-dialog
    v-model="dialogVisible"
    :title="dialogType === 'add' ? '카테고리 추가' : '카테고리 편집'"
    width="30%"
    :close-on-click-modal="false"
    @close="resetForm"
  >
    <el-form 
      ref="formRef" 
      :model="categoryForm" 
      :rules="rules" 
      label-width="100px"
      label-position="top"
    >
      <el-form-item label="카테고리 이름" prop="name">
        <el-input 
          v-model="categoryForm.name" 
          placeholder="카테고리 이름을 입력하세요" 
          clearable
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">취소</el-button>
        <el-button 
          type="primary" 
          @click="submitForm" 
          :loading="submitLoading"
        >
          {{ dialogType === 'add' ? '추가' : '저장' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { Edit, Delete, Plus } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import type { TableInstance, FormInstance, FormRules } from 'element-plus'
import { CategoryAPI } from '../../api/menu'
import type { Category } from '../../types/menuType'

// API 요청 시 에러 처리를 위한 유틸리티 함수
const handleApiError = (error: any, message: string) => {
  console.error(message, error)
  let errorMessage = message
  
  // API 응답에서 오류 메시지가 있으면 추가
  if (error.response && error.response.data && error.response.data.message) {
    errorMessage += `: ${error.response.data.message}`
  }
  
  ElMessage.error(errorMessage)
}

const tableRef = ref<TableInstance>()
const tableData = ref<Category[]>([])
const loading = ref(false)
const submitLoading = ref(false)

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

// 다이얼로그 관련 변수
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const formRef = ref<FormInstance>()
const categoryForm = reactive({
  id: 0,
  name: ''
})

// 폼 검증 규칙
const rules = reactive<FormRules>({
  name: [
    { required: true, message: '카테고리 이름을 입력해주세요', trigger: 'blur' },
    { min: 1, max: 50, message: '카테고리 이름은 1-50자 이내로 입력해주세요', trigger: 'blur' }
  ]
})

// 컴포넌트 마운트 시 데이터 로드
onMounted(async () => {
  await fetchData()
})

// API에서 데이터 가져오기
const fetchData = async () => {
  loading.value = true;
  try {
    const response = await CategoryAPI.getAllCategories();
    if (response && response.data) {
      tableData.value = response.data;
    }
  } catch (error) {
    handleApiError(error, '카테고리 데이터를 불러오는 중 오류가 발생했습니다');
  } finally {
    loading.value = false;
  }
}

// 페이지 변경 처리
const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

// 추가 버튼 핸들러
const handleAdd = () => {
  dialogType.value = 'add'
  resetForm()
  dialogVisible.value = true
}

// 수정 핸들러
const handleEdit = (_index: number, row: Category) => {
  dialogType.value = 'edit'
  categoryForm.id = row.id
  categoryForm.name = row.name
  dialogVisible.value = true
}

// 폼 초기화
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  categoryForm.id = 0
  categoryForm.name = ''
}

// 폼 제출 처리
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (dialogType.value === 'add') {
          // 객체가 아닌 문자열로 name을 전달
          await CategoryAPI.createCategory(categoryForm.name)
          ElMessage.success('카테고리가 추가되었습니다.')
        } else {
          // 두 번째 매개변수를 객체로 전달
          await CategoryAPI.updateCategory(categoryForm.id, categoryForm.name)
          ElMessage.success('카테고리가 수정되었습니다.')
        }
        
        dialogVisible.value = false
        await fetchData()
      } catch (error) {
        handleApiError(error, `카테고리를 ${dialogType.value === 'add' ? '추가' : '수정'}하는 중 오류가 발생했습니다`)
      } finally {
        submitLoading.value = false
      }
    }
  })
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
    handleApiError(error, '카테고리 삭제 중 오류가 발생했습니다')
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

.add-button {
  margin-left: auto;
}

.add-button :deep(.el-icon) {
  color: white;
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

/* 다이얼로그 스타일 */
.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
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
  
  .header-title {
    font-size: 16px;
  }
}

@media screen and (max-width: 576px) {
  .category-container {
    padding: 10px 5px;
  }
  
  .header-title {
    font-size: 16px;
  }
  
  .add-button {
    padding: 6px 10px;
    font-size: 12px;
  }
  
  :deep(.action-column) {
    width: 60px !important;
  }
  
  :deep(.el-table__row) {
    height: 50px;
    font-size: 13px;
  }
  
  .el-dialog {
    width: 90% !important;
    max-width: 320px;
  }
}
</style>