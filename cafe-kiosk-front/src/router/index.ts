import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import About from '../views/About.vue'
import OrderView from '../views/OrderView.vue'

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
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router