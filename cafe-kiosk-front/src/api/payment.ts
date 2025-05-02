// src/api/payment.ts
import apiClient from './index';
import type { CartItem } from '../types/menuType';

export interface PaymentResponse {
  success: boolean;
  message: string;
  details?: {
    actual_change: number;
    expected_amount: number;
    timeout_after: string;
  };
}

export interface OrderRequest {
  items: Array<{
    menu_id: number;
    quantity: number;
  }>;
}

export const PaymentAPI = {
  // 결제 요청 보내기
  requestPayment: (amount: number) => {
    return apiClient.post<PaymentResponse>('/payment', { amount });
  },

  postOrder: (cartItems: CartItem[]) => {
    // 백엔드 API 형식에 맞게 데이터 변환
    const items = cartItems.map(item => ({
      menu_id: item.item.id,
      quantity: item.quantity
    }));
    
    return apiClient.post<OrderRequest>('/orders', { items });
  },
};
