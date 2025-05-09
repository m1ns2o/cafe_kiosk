import { createRouter, createWebHistory } from 'vue-router'
import OrderView from '../views/OrderView.vue'
import PaymentView from '../views/PaymentView.vue'
import PaymentSuccessView from '../views/PaymentSuccessView.vue'
import AdminView from '../views/AdminView.vue'
import CategoryView from '../views/AdminView/CategoryView.vue'
import MenuView from '../views/AdminView/MenuView.vue'
import OrderManageView from '../views/AdminView/OrderManageView.vue'
import StatisticView from '../views/AdminView/StatisticView.vue'

const routes = [
  {
    path: '/',
    name: 'OrderView',
    component: OrderView
  },
  {
    path: '/payment/:totalAmount/:cartItems',
    name: 'PaymentView',
    component: PaymentView
  },
  {
    path: '/success',
    name: 'PaymentSuccessView',
    component: PaymentSuccessView 
  },
  {
    path: '/admin',
    name: 'Admin',
    component: AdminView,
    // 관리자 페이지 중첩 라우트 설정
    children: [
      {
        path: '', // /admin으로 접근 시 기본 대시보드 페이지로 리다이렉트
        redirect: '/admin/statistics'
      },
      {
        path: 'statistics', // /admin/statistics
        name: 'AdminStatistics',
        component: StatisticView
      },
      {
        path: 'category', // /admin/category
        name: 'Category',
        component: CategoryView
      },
      {
        path: 'menu', // /admin/menu
        name: 'AdminMenu',
        component: MenuView
      },
      {
        path: 'order', // /admin/order
        name: 'AdminOrder',
        component: OrderManageView
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router