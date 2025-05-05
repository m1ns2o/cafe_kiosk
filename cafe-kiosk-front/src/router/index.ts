import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import About from '../views/About.vue'
import OrderView from '../views/OrderView.vue'
import PaymentView from '../views/PaymentView.vue'
import PaymentSuccessView from '../views/PaymentSuccessView.vue'
import AdminView from '../views/AdminView.vue'
import CategoryView from '../views/AdminView/CategoryView.vue'
import MenuView from '../views/AdminView/MenuView.vue'
// import AdminDashboard from '../views/admin/AdminDashboard.vue'
// import AdminStatistics from '../views/admin/AdminStatistics.vue'
// import AdminCategory from '../views/admin/AdminCategory.vue'
// import AdminMenu from '../views/admin/AdminMenu.vue'
// import AdminOrder from '../views/admin/AdminOrder.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    component: About
  },
  {
    path: '/order',
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
      // {
      //   path: '', // /admin으로 접근 시 기본 대시보드 페이지로 리다이렉트
      //   name: 'AdminDashboard',
      //   component: AdminDashboard
      // },
      // {
      //   path: 'statistics', // /admin/statistics
      //   name: 'AdminStatistics',
      //   component: AdminStatistics
      // },
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
      // {
      //   path: 'order', // /admin/order
      //   name: 'AdminOrder',
      //   component: AdminOrder
      // }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router