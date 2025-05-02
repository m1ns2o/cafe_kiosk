// src/api/menu.ts
import apiClient from './index';
import type { MenuItem, Category } from '../types/menuType';

export const CategoryAPI = {
  // 모든 카테고리 가져오기
  getAllCategories: () => {
    return apiClient.get<Category[]>('/categories');
  },
  
  // 특정 카테고리 가져오기
  getCategoryById: (id: number) => {
    return apiClient.get<Category>(`/categories/${id}`);
  },
  
  // 특정 카테고리의 메뉴 가져오기
  getMenus: (categoryId: number) => {
    return apiClient.get<MenuItem[]>(`/categories/${categoryId}/menus`);
  }
};