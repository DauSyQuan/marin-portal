import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Import tất cả các màn hình giao diện
import DashboardView from '../stores/views/DashboardView.vue'
import LoginView from '../stores/views/LoginView.vue'
import ShipDetailView from '../stores/views/ShipDetailView.vue'
import AnalyticsView from '../stores/views/AnalyticsView.vue'
import SettingsView from '../stores/views/SettingsView.vue'
import UsageView from '../views/UsageView.vue' // <--- MỚI: Báo cáo lưu lượng
import OnlineUsersView from '../views/OnlineUsersView.vue' // <--- Import
import VoucherView from '../views/VoucherView.vue'
import BandwidthView from '../views/BandwidthView.vue' // <--- Import
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { 
      path: '/login', 
      name: 'login', 
      component: LoginView 
    },
    { 
      path: '/', 
      name: 'dashboard', 
      component: DashboardView, 
      meta: { requiresAuth: true } 
    },
    { 
      path: '/ship/:id', 
      name: 'ship-detail', 
      component: ShipDetailView, 
      meta: { requiresAuth: true } 
    },
    { 
      path: '/analytics', 
      name: 'analytics', 
      component: AnalyticsView, 
      meta: { requiresAuth: true } 
    },
    { 
      path: '/usage', 
      name: 'usage', 
      component: UsageView, 
      meta: { requiresAuth: true } // <--- Route mới cho trang Usage
    },
    { 
      path: '/settings', 
      name: 'settings', 
      component: SettingsView, 
      meta: { requiresAuth: true } 
    },
    { 
    path: '/online-users', 
    name: 'online-users', 
    component: OnlineUsersView, 
    meta: { requiresAuth: true } 
},
{ 
    path: '/vouchers', 
    name: 'vouchers', 
    component: VoucherView, 
    meta: { requiresAuth: true } 
},
{ 
    path: '/bandwidth', 
    name: 'bandwidth', 
    component: BandwidthView, 
    meta: { requiresAuth: true } 
},
  ]
})

// Chặn truy cập nếu chưa đăng nhập (Navigation Guard)
router.beforeEach((to, from, next) => {
  const auth = useAuthStore();
  const publicPages = ['/login'];
  const authRequired = !publicPages.includes(to.path);

  // 1. Nếu trang yêu cầu login mà chưa có token -> Đá về login
  if (authRequired && !auth.isAuthenticated) {
    return next('/login');
  }

  // 2. Nếu đã login mà cố vào trang login -> Đá về dashboard
  if (to.path === '/login' && auth.isAuthenticated) {
    return next('/');
  }

  next();
});

export default router