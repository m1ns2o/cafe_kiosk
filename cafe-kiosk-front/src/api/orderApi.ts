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

// 기간별 주문 조회
export async function getOrdersByPeriod(
  startDate: string, 
  endDate: string, 
  options?: {
    minAmount?: number;
    maxAmount?: number;
    menuId?: number;
    categoryId?: number;
    sortBy?: string;
    order?: 'asc' | 'desc';
  }
): Promise<{
  start_date: string;
  end_date: string;
  count: number;
  orders: Order[];
}> {
  let url = `/orders/period?start_date=${startDate}&end_date=${endDate}`;
  
  // 추가 옵션 쿼리 파라미터
  if (options) {
    if (options.minAmount) url += `&min_amount=${options.minAmount}`;
    if (options.maxAmount) url += `&max_amount=${options.maxAmount}`;
    if (options.menuId) url += `&menu_id=${options.menuId}`;
    if (options.categoryId) url += `&category_id=${options.categoryId}`;
    if (options.sortBy) url += `&sort_by=${options.sortBy}`;
    if (options.order) url += `&order=${options.order}`;
  }
  
  const res = await apiClient.get(url);
  return res.data;
}

