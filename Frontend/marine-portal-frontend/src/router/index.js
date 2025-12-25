import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Import các màn hình (Chỉ import 1 lần mỗi file)
import DashboardView from '../stores/views/DashboardView.vue'
import LoginView from '../stores/views/LoginView.vue'
import ShipDetailView from '../stores/views/ShipDetailView.vue'
import AnalyticsView from '../stores/views/AnalyticsView.vue'
import SettingsView from '../stores/views/SettingsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { 
      path: '/login', 
      component: LoginView 
    },
    { 
      path: '/', 
      component: DashboardView, 
      meta: { requiresAuth: true } 
    },
    { 
      path: '/ship/:id', 
      component: ShipDetailView, 
      meta: { requiresAuth: true } 
    },
    { 
      path: '/analytics', 
      component: AnalyticsView, 
      meta: { requiresAuth: true } 
    },
    { 
      path: '/settings', 
      component: SettingsView, 
      meta: { requiresAuth: true } 
    }
  ]
})

// Chặn truy cập nếu chưa đăng nhập
router.beforeEach((to, from, next) => {
  const auth = useAuthStore();
  const publicPages = ['/login'];
  const authRequired = !publicPages.includes(to.path);

  if (authRequired && !auth.isAuthenticated) {
    return next('/login');
  }

  // Đã login thì không cho vào trang login nữa
  if (to.path === '/login' && auth.isAuthenticated) {
    return next('/');
  }

  next();
});

export default router