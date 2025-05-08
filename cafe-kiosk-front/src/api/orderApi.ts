import apiClient from './index';

export interface Menu {
  id: number;
  name: string;
  price: number;
  image_url?: string;
}

export interface OrderItem {
  id: number;
  order_id: number;
  menu_id: number;
  menu: Menu;
  quantity: number;
  price: number;
}

export interface Order {
  id: number;
  total_price: number;
  created_at: string;
  updated_at: string;
  order_items: OrderItem[];
}

// 주문 목록 조회 (Order + OrderItems + Menu)
export async function getOrders(): Promise<Order[]> {
  const res = await apiClient.get('/orders');
  return res.data;
}

// 주문 상세 조회 (Order + OrderItems + Menu)
export async function getOrder(orderId: number): Promise<Order> {
  const res = await apiClient.get(`/orders/${orderId}`);
  return res.data;
}

