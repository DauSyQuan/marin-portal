<script setup>
import { useAuthStore } from '@/stores/auth';
const auth = useAuthStore();
</script>

<template>
  <div class="sidebar">
    <!-- 1. LOGO BRAND -->
    <div class="sidebar-brand">
      <div class="d-flex align-items-center">
        <div class="logo-icon me-3">
          <i class="fa-solid fa-anchor"></i>
        </div>
        <div>
          <h5 class="fw-bold m-0 text-white tracking-wide">MARINE PRO</h5>
          <small class="text-white-50" style="font-size: 0.65rem;">FLEET OPERATIONS</small>
        </div>
      </div>
    </div>
    
    <!-- 2. NAVIGATION MENU -->
    <nav class="mt-4 flex-grow-1 px-3">
      <div class="menu-label">Main Menu</div>
      
      <router-link to="/" class="nav-link" active-class="active">
        <i class="fa-solid fa-ship"></i>
        <span>Fleet Overview</span>
      </router-link>
      
      <router-link to="/analytics" class="nav-link" active-class="active">
        <i class="fa-solid fa-chart-pie"></i>
        <span>Analytics Report</span>
      </router-link>

      <div class="menu-label mt-4">Management</div>

      <router-link to="/usage" class="nav-link" active-class="active">
        <i class="fa-solid fa-file-invoice"></i>
        <span>User Monthly Usage</span>
      </router-link>

      <router-link to="/online-users" class="nav-link" active-class="active">
        <i class="fa-solid fa-laptop-code me-2"></i> <!-- Icon laptop -->
        <span>Online Users</span>
      </router-link>
      <router-link to="/vouchers" class="nav-link" active-class="active">
        <i class="fa-solid fa-ticket me-2"></i>
        <span>Vouchers</span>
      </router-link>
      <router-link to="/bandwidth" class="nav-link" active-class="active">
        <i class="fa-solid fa-bolt me-2"></i> <!-- Icon tia sét -->
        <span>Bandwidth Plans</span>
      </router-link>

      <router-link to="/settings" class="nav-link" active-class="active">
        <i class="fa-solid fa-gear"></i>
        <span>System Config</span>
      </router-link>
    </nav>

    <!-- 3. USER PROFILE FOOTER -->
    <div class="sidebar-footer">
      <div class="user-card d-flex align-items-center mb-3">
        <div class="avatar me-3">
          {{ auth.user?.username.charAt(0).toUpperCase() || 'A' }}
        </div>
        <div class="flex-grow-1 overflow-hidden">
          <div class="text-white fw-bold text-truncate" style="font-size: 0.9rem;">
            {{ auth.user?.username || 'Admin' }}
          </div>
          <div class="text-success small" style="font-size: 0.75rem;">
            <i class="fa-solid fa-circle fa-2xs me-1"></i>Online
          </div>
        </div>
      </div>
      <button @click="auth.logout()" class="btn btn-danger w-100 btn-sm fw-bold shadow-sm">
        <i class="fa-solid fa-right-from-bracket me-2"></i> Sign Out
      </button>
    </div>
  </div>
</template>

<style scoped>
/* --- CẤU TRÚC CHÍNH (MÀU NỀN ĐẬM) --- */
.sidebar { 
    width: 260px; 
    background-color: #0f172a; /* Màu Xanh Than Đậm (Dark Navy) - Rất quan trọng */
    height: 100vh; 
    position: fixed; 
    z-index: 1000; 
    display: flex; 
    flex-direction: column; 
    border-right: 1px solid rgba(255,255,255,0.05);
    box-shadow: 4px 0 24px rgba(0,0,0,0.4); /* Bóng đổ đậm hơn */
}

/* --- LOGO --- */
.sidebar-brand { 
    padding: 24px; 
    background: rgba(0,0,0,0.2); /* Nền tối hơn chút cho logo */
    border-bottom: 1px solid rgba(255,255,255,0.1); 
}
.logo-icon {
    width: 40px; height: 40px;
    background: linear-gradient(135deg, #2563eb, #1d4ed8);
    color: white;
    border-radius: 8px;
    display: flex; align-items: center; justify-content: center;
    font-size: 1.2rem;
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}
.tracking-wide { letter-spacing: 1px; }

/* --- MENU LINKS --- */
.menu-label {
    color: #475569;
    font-size: 0.7rem;
    font-weight: 700;
    text-transform: uppercase;
    margin-bottom: 8px;
    padding-left: 12px;
    letter-spacing: 0.5px;
}

.nav-link { 
    color: #94a3b8; /* Màu chữ xám sáng (dễ đọc trên nền đen) */
    padding: 12px 16px; 
    font-weight: 500; 
    cursor: pointer; 
    border-radius: 8px;
    transition: all 0.2s ease;
    text-decoration: none;
    display: flex; align-items: center;
    margin-bottom: 4px;
}

.nav-link i { 
    width: 24px; 
    margin-right: 12px; 
    text-align: center; 
    font-size: 1.1rem;
}

/* Hiệu ứng Hover & Active */
.nav-link:hover { 
    color: white; 
    background: rgba(255,255,255,0.08); 
    transform: translateX(4px);
}

.nav-link.active { 
    color: white; 
    background: linear-gradient(90deg, #2563eb 0%, #1d4ed8 100%); /* Gradient xanh nổi bật */
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3); /* Bóng đổ phát sáng */
    font-weight: 600;
}

/* --- FOOTER USER --- */
.sidebar-footer {
    padding: 20px;
    background: rgba(0,0,0,0.2);
    border-top: 1px solid rgba(255,255,255,0.1);
}
.avatar {
    width: 40px; height: 40px;
    background-color: #3b82f6;
    color: white;
    font-weight: bold;
    border-radius: 10px;
    display: flex; align-items: center; justify-content: center;
    font-size: 1.2rem;
}
</style>