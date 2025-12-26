<script setup>
import { useAuthStore } from '@/stores/auth';
const auth = useAuthStore();
</script>

<template>
  <div class="sidebar glass-sidebar">
    
    <!-- 1. BRAND LOGO (Hiệu ứng phát sáng) -->
    <div class="sidebar-header">
      <div class="logo-container">
        <i class="fa-solid fa-anchor logo-icon"></i>
      </div>
      <div class="brand-text">
        <h5 class="m-0 fw-bold text-white tracking-wide">MARINE PRO</h5>
        <small class="text-blue-300" style="font-size: 0.65rem; letter-spacing: 2px;">FLEET OPS</small>
      </div>
    </div>
    
    <!-- 2. NAVIGATION MENU -->
    <nav class="sidebar-nav">
      
      <!-- Group 1 -->
      <div class="nav-label">Main Menu</div>
      <router-link to="/" class="nav-item" active-class="active">
        <i class="fa-solid fa-ship"></i>
        <span>Fleet Overview</span>
      </router-link>
      <router-link to="/analytics" class="nav-item" active-class="active">
        <i class="fa-solid fa-chart-pie"></i>
        <span>Analytics Report</span>
      </router-link>

      <!-- Group 2 -->
      <div class="nav-label mt-4">Management</div>
      <router-link to="/usage" class="nav-item" active-class="active">
        <i class="fa-solid fa-file-invoice"></i>
        <span>Monthly Usage</span>
      </router-link>
      <router-link to="/online-users" class="nav-item" active-class="active">
        <i class="fa-solid fa-laptop-code"></i>
        <span>Online Users</span>
      </router-link>
      <router-link to="/vouchers" class="nav-item" active-class="active">
        <i class="fa-solid fa-ticket"></i>
        <span>Vouchers</span>
      </router-link>
      <router-link to="/bandwidth" class="nav-item" active-class="active">
        <i class="fa-solid fa-bolt"></i>
        <span>Bandwidth Plans</span>
      </router-link>

      <!-- Group 3 -->
      <div class="nav-label mt-4">System</div>
      <router-link to="/settings" class="nav-item" active-class="active">
        <i class="fa-solid fa-sliders"></i>
        <span>Configuration</span>
      </router-link>
    </nav>

    <!-- 3. USER PROFILE (Glass Card nhỏ gọn) -->
    <div class="sidebar-footer">
      <div class="user-glass-card">
        <div class="d-flex align-items-center gap-3">
          <!-- Avatar -->
          <div class="user-avatar">
            {{ (auth.user?.username || 'A').charAt(0).toUpperCase() }}
            <span class="status-dot"></span>
          </div>
          <!-- Info -->
          <div class="user-info">
            <div class="fw-bold text-white text-truncate" style="max-width: 90px;">
              {{ auth.user?.username || 'Admin' }}
            </div>
            <div class="text-white-50 small" style="font-size: 0.7rem;">Manager</div>
          </div>
          <!-- Logout Icon Button (Thay cho nút đỏ to) -->
          <button @click="auth.logout()" class="btn-logout" title="Sign Out">
            <i class="fa-solid fa-right-from-bracket"></i>
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
/* --- CẤU TRÚC CHÍNH (GLASSMORPHISM) --- */
.glass-sidebar {
    width: 260px;
    height: 100vh;
    display: flex;
    flex-direction: column;
    /* Hiệu ứng kính mờ tối */
    background: rgba(15, 23, 42, 0.65); 
    backdrop-filter: blur(16px);
    -webkit-backdrop-filter: blur(16px);
    border-right: 1px solid rgba(255, 255, 255, 0.08);
}

/* --- 1. HEADER --- */
.sidebar-header {
    padding: 30px 24px;
    display: flex;
    align-items: center;
    gap: 15px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.logo-container {
    width: 42px; height: 42px;
    background: linear-gradient(135deg, #3b82f6, #2563eb);
    border-radius: 12px;
    display: flex; align-items: center; justify-content: center;
    box-shadow: 0 0 20px rgba(59, 130, 246, 0.4); /* Glow xanh */
}
.logo-icon { color: white; font-size: 1.2rem; }
.text-blue-300 { color: #93c5fd; }
.tracking-wide { letter-spacing: 1px; }

/* --- 2. MENU --- */
.sidebar-nav {
    flex-grow: 1;
    padding: 24px 16px;
    overflow-y: auto;
}

.nav-label {
    color: rgba(255, 255, 255, 0.4);
    font-size: 0.7rem;
    font-weight: 700;
    text-transform: uppercase;
    margin-bottom: 12px;
    padding-left: 12px;
    letter-spacing: 1px;
}

/* Style cho từng mục menu */
.nav-item {
    display: flex; align-items: center;
    padding: 12px 16px;
    margin-bottom: 4px;
    color: #94a3b8; /* Màu xám sáng */
    text-decoration: none;
    font-weight: 500;
    border-radius: 12px;
    transition: all 0.3s ease;
    border: 1px solid transparent; /* Giữ chỗ cho border */
}

.nav-item i { width: 24px; margin-right: 12px; text-align: center; font-size: 1.1rem; transition: color 0.3s; }

/* Hover */
.nav-item:hover {
    color: white;
    background: rgba(255, 255, 255, 0.05);
}

/* Active State (Glow Effect) */
.nav-item.active {
    color: white;
    background: linear-gradient(90deg, rgba(59, 130, 246, 0.15) 0%, rgba(59, 130, 246, 0) 100%);
    border-left: 3px solid #3b82f6; /* Viền trái xanh */
    border-radius: 4px 12px 12px 4px; /* Bo góc phải nhiều hơn */
}
.nav-item.active i { color: #3b82f6; text-shadow: 0 0 10px rgba(59, 130, 246, 0.6); }

/* --- 3. FOOTER PROFILE --- */
.sidebar-footer {
    padding: 20px;
    border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.user-glass-card {
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.05);
    border-radius: 16px;
    padding: 12px;
    transition: background 0.3s;
}
.user-glass-card:hover { background: rgba(0, 0, 0, 0.3); }

.user-avatar {
    width: 36px; height: 36px;
    background: #475569;
    color: white;
    font-weight: bold;
    border-radius: 10px;
    display: flex; align-items: center; justify-content: center;
    position: relative;
}

.status-dot {
    width: 8px; height: 8px;
    background-color: #10b981; /* Xanh lá online */
    border: 2px solid #1e293b;
    border-radius: 50%;
    position: absolute;
    bottom: -2px; right: -2px;
}

.user-info { flex-grow: 1; }

.btn-logout {
    background: transparent;
    border: none;
    color: #94a3b8;
    cursor: pointer;
    padding: 8px;
    border-radius: 8px;
    transition: all 0.2s;
}
.btn-logout:hover {
    background: rgba(239, 68, 68, 0.15); /* Nền đỏ nhạt */
    color: #ef4444; /* Icon đỏ */
}

/* Custom Scrollbar */
.sidebar-nav::-webkit-scrollbar { width: 4px; }
.sidebar-nav::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 4px; }
</style>
