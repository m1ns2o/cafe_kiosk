<template>
  <div class="menu-container">
    <div class="menu-content">
      <el-card shadow="hover" class="menu-card">
        <template #header>
          <div class="card-header">
            <span class="header-title">메뉴 관리</span>
            <!-- 메뉴 추가 버튼 -->
            <el-button 
              type="primary" 
              size="small" 
              @click="handleAdd"
              class="add-button"
            >
              <el-icon class="el-icon--left" style="color: white;"><Plus /></el-icon>
              메뉴 추가
            </el-button>
          </div>
        </template>
        
        <el-table 
          ref="tableRef" 
          row-key="id" 
          :data="pagedTableData" 
          header-cell-class-name="table-header"
          v-loading="loading"
          class="menu-table"
        >
          <!-- 이미지 열 -->
          <el-table-column label="이미지" class-name="image-column" align="left">
            <template #default="scope">
              <div class="image-container">
                <img v-if="scope.row.image_url" 
                    :src="scope.row.image_url" 
                    :alt="scope.row.name" 
                    class="menu-image" />
                <el-icon v-else><PictureRounded /></el-icon>
              </div>
            </template>
          </el-table-column>
          
          <!-- 메뉴 이름 열 -->
          <el-table-column prop="name" label="메뉴" class-name="name-column" align="left">
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
          <el-table-column prop="price" label="가격" class-name="price-column" align="left">
            <template #default="scope">
              <span>{{ formatPrice(scope.row.price) }}</span>
            </template>
          </el-table-column>
          
          <!-- 카테고리 열 -->
          <el-table-column
            prop="category_id"
            label="카테고리"
            class-name="category-column"
            align="left"
            :filters="categoryFilters"
            :filter-method="filterByCategory"
            filter-placement="bottom-end"
            @filter-change="onFilterChange"
          >
            <template #default="scope">
              <span>{{ getCategoryName(scope.row.category_id) }}</span>
            </template>
          </el-table-column>
          
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
                  @click="confirmDelete(scope.row)"
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

  <!-- 메뉴 추가/편집 다이얼로그 -->
  <el-dialog
    v-model="dialogVisible"
    :title="dialogType === 'add' ? '메뉴 추가' : '메뉴 편집'"
    width="40%"
    :close-on-click-modal="false"
    @close="resetForm"
  >
    <el-form 
      ref="formRef" 
      :model="menuForm" 
      :rules="rules" 
      label-width="120px"
      label-position="left"
    >
      <el-form-item label="메뉴 이름" prop="name">
        <el-input 
          v-model="menuForm.name" 
          placeholder="메뉴 이름을 입력하세요" 
          clearable
        />
      </el-form-item>

      <el-form-item label="가격" prop="price">
        <el-input-number 
          v-model="menuForm.price" 
          :min="0" 
          :precision="0" 
          :step="100"
          :controls="true"
          style="width: 100%;"
        />
      </el-form-item>

      <el-form-item label="카테고리" prop="category_id">
        <el-select
          v-model="menuForm.category_id"
          placeholder="카테고리 선택"
          style="width: 100%;"
          clearable
        >
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="이미지">
        <el-upload
          class="menu-image-upload"
          :auto-upload="false"
          :show-file-list="true"
          :limit="1"
          accept="image/*"
          list-type="picture-card"
          :on-change="handleImageChange"
          :on-remove="handleImageRemove"
          :file-list="fileList"
          :on-exceed="handleExceed"
          :class="{ 'upload-hidden': fileList.length >= 1 }"
        >
          <div v-show="fileList.length < 1">
            <el-icon><Plus /></el-icon>
          </div>
          <template #file="{ file }">
            <div>
              <img class="el-upload-list__item-thumbnail" :src="file.url" alt="" />
              <span class="el-upload-list__item-actions">
                <span class="el-upload-list__item-preview" @click="previewImage(file)">
                  <el-icon><ZoomIn /></el-icon>
                </span>
                <span class="el-upload-list__item-delete" @click="handleImageRemove(file)">
                  <el-icon><Delete /></el-icon>
                </span>
              </span>
            </div>
          </template>
        </el-upload>
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

  <!-- 이미지 미리보기 -->
  <el-dialog v-model="previewVisible" title="이미지 미리보기">
    <img :src="previewImageUrl" alt="Preview" style="max-width: 100%;" />
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, watch, reactive } from 'vue'
import { PictureRounded, Edit, Delete, Plus, ZoomIn } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import type { TableInstance, FormInstance, FormRules, UploadFile, UploadUserFile } from 'element-plus'
import { MenuAPI, CategoryAPI } from '../../api/menu'
import type { MenuItem, Category } from '../../types/menuType'

const tableRef = ref<TableInstance>()
const tableData = ref<MenuItem[]>([])
const categories = ref<Category[]>([])
const loading = ref(false)
const submitLoading = ref(false)

// 페이지네이션 관련 변수
const currentPage = ref(1)
const pageSize = ref(8) // 페이지당 8개 항목

// 필터링된 테이블 데이터 (카테고리 필터 적용)
const categoryFilterValue = ref<number[] | null>(null)

// 다이얼로그 관련 변수
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const formRef = ref<FormInstance>()
const menuForm = reactive<{
  id: number;
  name: string;
  price: number;
  category_id: number | null;
  image_url: string;
}>({
  id: 0,
  name: '',
  price: 0,
  category_id: null,
  image_url: ''
})

// 이미지 업로드 관련
const fileList = ref<UploadUserFile[]>([])
const uploadedFile = ref<File | null>(null)
const previewVisible = ref(false)
const previewImageUrl = ref('')

// 폼 검증 규칙
const rules = reactive<FormRules>({
  name: [
    { required: true, message: '메뉴 이름을 입력해주세요', trigger: 'blur' },
    { min: 1, max: 50, message: '메뉴 이름은 1-50자 이내로 입력해주세요', trigger: 'blur' }
  ],
  price: [
    { required: true, message: '가격을 입력해주세요', trigger: 'blur' },
    { type: 'number', min: 0, message: '가격은 0 이상이어야 합니다', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '카테고리를 선택해주세요', trigger: 'change' }
  ]
})

// API 요청 시 에러 처리를 위한 유틸리티 함수
const handleApiError = (error: any, message: string) => {
  console.error(message, error)
  let errorMessage = message
  
  // API 응답에서 오류 메시지가 있으면 추가
  if (error.response && error.response.data) {
    if (error.response.data.error) {
      errorMessage += `: ${error.response.data.error}`
    } else if (typeof error.response.data === 'string') {
      errorMessage += `: ${error.response.data}`
    } else if (error.response.data.message) {
      errorMessage += `: ${error.response.data.message}`
    }
  }
  
  ElMessage.error(errorMessage)
}

const filteredTableData = computed(() => {
  if (!categoryFilterValue.value || categoryFilterValue.value.length === 0) {
    return tableData.value
  }
  return tableData.value.filter(row => categoryFilterValue.value!.includes(row.category_id))
})

// 현재 페이지에 표시할 데이터
const pagedTableData = computed(() => {
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return filteredTableData.value.slice(startIndex, endIndex)
})

// 필터 변경 시 첫 페이지로 돌아가기
watch(() => filteredTableData.value.length, () => {
  currentPage.value = 1
})

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
    
    if (menuResponse.data) {
      tableData.value = menuResponse.data
    }
    
    if (categoryResponse.data) {
      categories.value = categoryResponse.data
    }
  } catch (error) {
    handleApiError(error, '데이터를 불러오는 중 오류가 발생했습니다')
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

// 필터 변경 이벤트 핸들러
const onFilterChange = (filters: Record<string, number[]>) => {
  categoryFilterValue.value = filters.category_id || []
}

// 페이지 변경 처리
const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

// 가격 포맷팅 (원화)
const formatPrice = (price: number): string => {
  return new Intl.NumberFormat('ko-KR', { 
    style: 'currency', 
    currency: 'KRW',
    minimumFractionDigits: 0
  }).format(price)
}

// 폼 초기화
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  menuForm.id = 0
  menuForm.name = ''
  menuForm.price = 0
  menuForm.category_id = null
  menuForm.image_url = ''
  fileList.value = []
  uploadedFile.value = null
}

// 추가 버튼 핸들러
const handleAdd = () => {
  dialogType.value = 'add'
  resetForm()
  dialogVisible.value = true
}

// 수정 핸들러
const handleEdit = (_index: number, row: MenuItem) => {
  dialogType.value = 'edit'
  
  // 폼 데이터 설정
  menuForm.id = row.id
  menuForm.name = row.name
  menuForm.price = row.price
  menuForm.category_id = row.category_id
  menuForm.image_url = row.image_url || ''
  
  // 파일 목록 초기화 후 기존 이미지가 있으면 파일 목록에 추가
  fileList.value = []
  uploadedFile.value = null
  
  if (row.image_url) {
    // 기존 이미지를 fileList에 추가하여 업로드 컴포넌트에 표시
    fileList.value = [{
      name: row.name + '-이미지',
      url: row.image_url
    }]
  }
  
  console.log('수정할 메뉴 데이터:', menuForm, '파일 목록:', fileList.value)
  dialogVisible.value = true
}

// 이미지 변경 핸들러
const handleImageChange = (file: UploadFile) => {
  console.log('이미지 변경:', file, '현재 목록:', fileList.value)
  uploadedFile.value = file.raw as File
  const reader = new FileReader()
  reader.onload = (e) => {
    if (e.target && e.target.result) {
      file.url = e.target.result as string
    }
  }
  if (file.raw) {
    reader.readAsDataURL(file.raw as Blob)
  }
}

// 이미지 개수 초과 핸들러
const handleExceed = () => {
  ElMessage.warning('이미지는 하나만 선택할 수 있습니다.')
}

// 이미지 제거 핸들러
const handleImageRemove = (file: UploadFile) => {
  console.log('이미지 제거:', file)
  fileList.value = []
  uploadedFile.value = null
}

// 이미지 미리보기
const previewImage = (file: UploadFile) => {
  previewImageUrl.value = file.url || ''
  previewVisible.value = true
}

// 폼 제출 처리
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        // FormData 객체 생성
        const formData = new FormData()
        formData.append('name', menuForm.name)
        formData.append('price', menuForm.price.toString())
        
        // 카테고리 ID가 null이 아닌 경우에만 추가
        if (menuForm.category_id !== null) {
          formData.append('category_id', menuForm.category_id.toString())
        }
        
        // 파일이 있으면 추가
        if (uploadedFile.value) {
          formData.append('image', uploadedFile.value)
        }
        
        // 이미지 변경 여부 처리 (편집 모드에서만 적용)
        const imageChanged = dialogType.value === 'edit' && 
                           ((menuForm.image_url && fileList.value.length === 0) || 
                            (uploadedFile.value !== null));
        
        if (imageChanged) {
          formData.append('image_changed', 'true')
        }
        
        if (dialogType.value === 'add') {
          console.log('메뉴 추가 요청:', Object.fromEntries(formData))
          const response = await MenuAPI.createMenu(formData)
          console.log('메뉴 추가 응답:', response)
          ElMessage.success('메뉴가 추가되었습니다.')
        } else {
          console.log('메뉴 수정 요청:', Object.fromEntries(formData))
          const response = await MenuAPI.updateMenu(menuForm.id, formData)
          console.log('메뉴 수정 응답:', response)
          ElMessage.success('메뉴가 수정되었습니다.')
        }
        
        dialogVisible.value = false
        await fetchData()
      } catch (error) {
        console.error('API 오류 세부 정보:', error)
        handleApiError(error, `메뉴를 ${dialogType.value === 'add' ? '추가' : '수정'}하는 중 오류가 발생했습니다`)
      } finally {
        submitLoading.value = false
      }
    }
  })
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
    handleApiError(error, '메뉴 삭제 중 오류가 발생했습니다')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.menu-container {
  display: flex;
  justify-content: center;
  align-items: center; /* flex-start에서 center로 변경 */
  width: 100%;
  min-height: 100%; /* 80vh에서 100vh로 변경하여 전체 높이 활용 */
  padding: 20px;
}

.menu-content {
  width: 100%;
  max-width: 1200px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.menu-card {
  width: 95%; /* 너비 제한 추가 */
  height: 100%;
  margin: 0 auto; /* 상하 여백 제거하고 좌우 자동 마진 설정 */
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

.menu-table {
  width: 100%;
  margin-bottom: 10px;
}

.pagination-container {
  margin-top: 20px;
  margin-bottom: 10px;
  display: flex;
  justify-content: center;
}

.image-container {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  height: 50px;
}

.menu-image {
  max-height: 50px;
  max-width: 50px;
  object-fit: contain;
  border-radius: 4px;
}

/* 테이블 헤더와 셀 정렬 */
.table-header {
  text-align: left !important;
  font-weight: bold;
  background-color: #f5f7fa !important;
  height: 50px;
  font-size: 15px;
}

:deep(.el-table__header-wrapper .el-table__header th) {
  text-align: left !important;
}

:deep(.el-table__row) {
  height: 60px;
  font-size: 14px;
}


/* 개선된 컬럼 너비 설정 - 백분율 기반 */
:deep(.image-column) {
  width: 100px !important;
}

:deep(.name-column) {
  width: 35% !important; 
}

:deep(.price-column) {
  width: 15% !important;
}

:deep(.category-column) {
  width: 20% !important;
}

:deep(.action-column) {
  width: 80px !important;
}

:deep(.el-table__row .image-column) {
  text-align: left !important;
}

:deep(.el-table__row .name-column) {
  text-align: left !important;
}

:deep(.el-table__row .price-column) {
  text-align: left !important;
}

:deep(.el-table__row .category-column) {
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

/* 다이얼로그 관련 스타일 */
.menu-image-upload {
  width: 100%;
}

.current-image-container {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.menu-image-upload {
  width: 100%;
}

.upload-hidden :deep(.el-upload--picture-card) {
  display: none !important;
}

.menu-image-upload :deep(.el-upload-list--picture-card) {
  margin-top: 10px;
}

.menu-image-upload :deep(.el-upload-list__item-actions) {
  opacity: 0;
  transition: opacity 0.3s;
}

.menu-image-upload :deep(.el-upload-list__item):hover .el-upload-list__item-actions {
  opacity: 1;
}

.menu-image-upload :deep(.el-upload-list__item-thumbnail) {
  object-fit: cover;
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 미디어 쿼리 */
@media screen and (min-width: 1201px) {
  :deep(.image-column) {
    width: 100px !important;
  }
  
  :deep(.name-column) {
    width: 35% !important;
  }
  
  :deep(.price-column) {
    width: 15% !important;
  }
  
  :deep(.category-column) {
    width: 20% !important;
  }
  
  :deep(.action-column) {
    width: 80px !important;
  }
}

@media screen and (min-width: 993px) and (max-width: 1200px) {
  :deep(.image-column) {
    width: 90px !important;
  }
  
  :deep(.name-column) {
    width: 35% !important;
  }
  
  :deep(.price-column) {
    width: 15% !important;
  }
  
  :deep(.category-column) {
    width: 20% !important;
  }
  
  :deep(.action-column) {
    width: 80px !important;
  }
}

@media screen and (min-width: 769px) and (max-width: 992px) {
  :deep(.image-column) {
    width: 80px !important;
  }
  
  :deep(.name-column) {
    width: 35% !important;
  }
  
  :deep(.price-column) {
    width: 15% !important;
  }
  
  :deep(.category-column) {
    width: 20% !important;
  }
  
  :deep(.action-column) {
    width: 80px !important;
  }
}

@media screen and (min-width: 577px) and (max-width: 768px) {
  .menu-container {
    padding: 10px;
  }
  
  :deep(.image-column) {
    width: 70px !important;
  }
  
  :deep(.name-column) {
    width: 35% !important;
  }
  
  :deep(.price-column) {
    width: 15% !important;
  }
  
  :deep(.category-column) {
    width: 20% !important;
  }
  
  :deep(.action-column) {
    width: 70px !important;
  }
  
  .el-dialog {
    width: 95% !important;
  }
}

@media screen and (min-width: 401px) and (max-width: 576px) {
  .menu-container {
    padding: 10px;
  }
  
  :deep(.image-column) {
    width: 60px !important;
  }
  
  :deep(.name-column) {
    width: 35% !important;
  }
  
  :deep(.price-column) {
    width: 15% !important;
  }
  
  :deep(.category-column) {
    width: 20% !important;
  }
  
  :deep(.action-column) {
    width: 60px !important;
  }
  
  .el-dialog {
    width: 95% !important;
  }
}

@media screen and (max-width: 400px) {
  .menu-container {
    padding: 5px;
  }
  
  :deep(.image-column) {
    width: 55px !important;
  }
  
  :deep(.name-column) {
    width: 35% !important;
  }
  
  :deep(.price-column) {
    width: 15% !important;
  }
  
  :deep(.category-column) {
    width: 15% !important;
  }
  
  :deep(.action-column) {
    width: 55px !important;
  }
  
  .el-dialog {
    width: 95% !important;
  }
  
  .el-form-item {
    margin-bottom: 15px;
  }
}
</style>