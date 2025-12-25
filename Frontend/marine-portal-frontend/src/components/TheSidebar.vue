<script setup>
import { useAuthStore } from "@/stores/auth";
import { useRouter } from "vue-router";

const auth = useAuthStore();
const router = useRouter();

const handleLogout = () => {
  auth.logout();
  router.push("/login");
};
</script>

<template>
  <aside class="sidebar">
    <!-- Brand -->
    <div class="brand">
      <div class="brand-icon">
        <i class="fa-solid fa-satellite"></i>
      </div>
      <div class="brand-text">
        <div class="brand-title">MARINE PRO</div>
        <div class="brand-sub">Fleet Operations</div>
      </div>
    </div>

    <!-- Menu -->
    <nav class="menu">
      <router-link to="/" class="menu-item" active-class="active">
        <i class="fa-solid fa-ship"></i>
        <span>Fleet Overview</span>
      </router-link>

      <router-link to="/analytics" class="menu-item" active-class="active">
        <i class="fa-solid fa-chart-pie"></i>
        <span>Analytics Report</span>
      </router-link>

      <router-link to="/settings" class="menu-item" active-class="active">
        <i class="fa-solid fa-gear"></i>
        <span>System Config</span>
      </router-link>
    </nav>

    <!-- Footer -->
    <div class="sidebar-footer">
      <div class="user-box">
        <div class="avatar">A</div>
        <div class="user-meta">
          <div class="user-name">Admin User</div>
          <div class="user-role">Fleet Manager</div>
        </div>
      </div>

      <button class="logout-btn" @click="handleLogout">
        <i class="fa-solid fa-right-from-bracket"></i>
        <span>Sign Out</span>
      </button>
    </div>
  </aside>
</template>

<style scoped>
/* ✅ Sidebar MUST be fully styled here to avoid global override */
.sidebar {
  width: 260px;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  z-index: 1000;

  display: flex;
  flex-direction: column;

  /* ✅ solid dark background */
  background: #0f172a !important;
  color: #ffffff !important;
  border-right: 1px solid rgba(255, 255, 255, 0.08);
}

/* Brand */
.brand {
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 18px 18px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.brand-icon {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  background: rgba(59, 130, 246, 0.18);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #3b82f6;
  font-size: 18px;
}

.brand-title {
  font-weight: 900;
  letter-spacing: 0.08em;
  font-size: 0.95rem;
}

.brand-sub {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.55);
  margin-top: 2px;
}

/* Menu */
.menu {
  padding: 12px 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;

  padding: 12px 14px;
  border-radius: 12px;

  color: rgba(255, 255, 255, 0.78) !important;
  text-decoration: none !important;
  font-weight: 700;

  transition: all 0.18s ease;
}

.menu-item i {
  width: 18px;
  text-align: center;
  opacity: 0.95;
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #ffffff !important;
}

.menu-item.active {
  background: rgba(59, 130, 246, 0.14);
  color: #ffffff !important;
  box-shadow: inset 0 0 0 1px rgba(59, 130, 246, 0.25);
}

/* Footer */
.sidebar-footer {
  padding: 16px 14px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.user-box {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 12px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 999px;
  background: rgba(59, 130, 246, 0.18);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 900;
}

.user-name {
  font-weight: 900;
  font-size: 0.9rem;
  color: #ffffff;
}

.user-role {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.55);
  margin-top: 2px;
}

/* Logout */
.logout-btn {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;

  border: 0;
  background: transparent;
  color: #f87171;

  padding: 12px 12px;
  border-radius: 12px;

  font-weight: 800;
  cursor: pointer;
  transition: all 0.18s ease;
}

.logout-btn:hover {
  background: rgba(248, 113, 113, 0.12);
}

/* Responsive: nếu bạn muốn sidebar ẩn trên mobile sau này */
@media (max-width: 992px) {
  .sidebar {
    position: relative;
    width: 100%;
    height: auto;
  }
}
</style>
