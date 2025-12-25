import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "../stores/auth";

// Lazy imports từ src/stores/views/
const LoginView = () => import("../stores/views/LoginView.vue");
const DashboardView = () => import("../stores/views/DashboardView.vue");
const ShipDetailView = () => import("../stores/views/ShipDetailView.vue");

// ✅ thêm 2 view match sidebar
const AnalyticsView = () => import("../stores/views/AnalyticsView.vue");
const SettingsView = () => import("../stores/views/SettingsView.vue");

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/login", name: "Login", component: LoginView },

    { path: "/", name: "Dashboard", component: DashboardView, meta: { requiresAuth: true } },

    { path: "/ship/:id", name: "ShipDetail", component: ShipDetailView, meta: { requiresAuth: true } },

    // ✅ NEW: fix /analytics & /settings warnings
    { path: "/analytics", name: "Analytics", component: AnalyticsView, meta: { requiresAuth: true } },

    { path: "/settings", name: "Settings", component: SettingsView, meta: { requiresAuth: true } },

    // optional fallback
    { path: "/:pathMatch(.*)*", redirect: "/" },
  ],
});

// Router guard: chặn nếu chưa login
router.beforeEach((to, from, next) => {
  const auth = useAuthStore();

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return next({ name: "Login" });
  }

  // Nếu đã login mà vẫn vào /login thì đá về dashboard
  if (to.name === "Login" && auth.isAuthenticated) {
    return next({ name: "Dashboard" });
  }

  next();
});

export default router;
