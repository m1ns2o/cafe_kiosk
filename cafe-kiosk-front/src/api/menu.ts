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
  },

  // 카테고리 생성
  createCategory: (name: string) => {
    return apiClient.post<Category>('/categories', { name });
  },

  // 카테고리 삭제
  deleteCategory: (id: number) => {
    return apiClient.delete(`/categories/${id}`);
  },

  // 카테고리 수정
  updateCategory: (id: number, name: string) => {
    // console.log('updateCategory API 호출 (객체로 전달):', id, name);
    return apiClient.put<Category>(`/categories/${id}`, { name });
  },
};

export const MenuAPI = {
  // 전체 메뉴 조회 (카테고리별 필터링 지원)
  getMenus: (categoryId?: number) => {
    if (categoryId) {
      return apiClient.get<MenuItem[]>(`/categories/${categoryId}/menus`);
    }
    return apiClient.get<MenuItem[]>('/menus');
  },

  // 메뉴 생성 (이미지 포함)
  createMenu: (data: FormData) => {
    return apiClient.post<MenuItem>('/menus', data, { headers: { 'Content-Type': 'multipart/form-data' } });
  },

  // 메뉴 수정 (이미지 포함)
  updateMenu: (id: number, data: FormData) => {
    return apiClient.put<MenuItem>(`/menus/${id}`, data, { headers: { 'Content-Type': 'multipart/form-data' } });
  },

  // 메뉴 삭제
  deleteMenu: (id: number) => {
    return apiClient.delete(`/menus/${id}`);
  },
};

export const OrderAPI = {
  // 전체 주문 조회
  getOrders: () => {
    return apiClient.get('/orders');
  },

  // 단일 주문 조회
  getOrder: (id: number) => {
    return apiClient.get(`/orders/${id}`);
  },
};