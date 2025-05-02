// src/api/payment.ts
import apiClient from './index';

export interface PaymentResponse {
  success: boolean;
  message: string;
  details?: {
    actual_change: number;
    expected_amount: number;
    timeout_after: string;
  };
}

export const PaymentAPI = {
  // 결제 요청 보내기
  requestPayment: (amount: number) => {
    return apiClient.post<PaymentResponse>('/payment', { amount });
  }
};
