// src/api/categories.ts
import apiClient from './index';

// 타입 정의
export interface Category {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Menu {
  id: number;
  name: string;
  price: number;
  description?: string;
  // 기타 필요한 속성들...
}

export const CategoryAPI = {
  // 모든 카테고리 가져오기
  getAll: () => {
    return apiClient.get<Category[]>('/categories');
  },
  
  // 특정 카테고리 가져오기
  getById: (id: number) => {
    return apiClient.get<Category>(`/categories/${id}`);
  },
  
  // 특정 카테고리의 메뉴 가져오기
  getMenus: (categoryId: number) => {
    return apiClient.get<Menu[]>(`/categories/${categoryId}/menus`);
  }
};