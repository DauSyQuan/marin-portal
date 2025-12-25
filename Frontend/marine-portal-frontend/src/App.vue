<script setup>
import TheSidebar from './components/TheSidebar.vue'
import { useRoute } from 'vue-router'

const route = useRoute()
</script>

<template>
  <div class="app-layout">
    
    <!-- 1. HIỆU ỨNG SIDEBAR (Trượt từ trái vào) -->
    <transition name="sidebar-slide">
      <TheSidebar v-if="route.path !== '/login'" />
    </transition>

    <!-- 2. WRAPPER BAO NGOÀI -->
    <div class="main-wrapper" :class="{ 'login-mode': route.path === '/login' }">
      <div class="main-content">
        
        <!-- 3. HIỆU ỨNG NỘI DUNG (Phóng to & Mờ dần) -->
        <router-view v-slot="{ Component }">
          <transition name="book-reveal" mode="out-in">
            <component :is="Component" :key="route.fullPath" />
          </transition>
        </router-view>

      </div>
    </div>
  </div>
</template>

<style>
/* --- CẤU TRÚC CƠ BẢN --- */
body {
    background-color: #0f172a;
    margin: 0;
    overflow-x: hidden; /* Chặn thanh cuộn ngang khi animation chạy */
}

.app-layout {
    display: flex;
    width: 100%;
    min-height: 100vh;
}

.main-wrapper {
    flex-grow: 1;
    margin-left: 260px;
    background-color: #0f172a;
    padding: 12px 12px 12px 0;
    min-height: 100vh;
    transition: margin-left 0.5s cubic-bezier(0.25, 0.8, 0.25, 1); /* Mượt mà khi Login/Logout */
}

.main-content {
    background-color: #f1f5f9;
    width: 100%;
    height: 100%;
    min-height: calc(100vh - 24px);
    border-radius: 32px;
    padding: 30px;
    box-shadow: -10px 0 30px rgba(0, 0, 0, 0.2);
    position: relative;
    overflow: hidden;
}

/* Chế độ Login */
.main-wrapper.login-mode {
    margin-left: 0;
    padding: 0;
}
.main-wrapper.login-mode .main-content {
    border-radius: 0;
    padding: 0;
}

/* --- A. HIỆU ỨNG SIDEBAR (Trượt như gáy sách) --- */
.sidebar-slide-enter-active,
.sidebar-slide-leave-active {
    transition: all 0.6s cubic-bezier(0.16, 1, 0.3, 1); /* Ease-out Quint */
}

.sidebar-slide-enter-from,
.sidebar-slide-leave-to {
    transform: translateX(-100%); /* Trượt hẳn ra ngoài màn hình */
    opacity: 0;
}

/* --- B. HIỆU ỨNG TRANG SÁCH (Content Reveal) --- */
.book-reveal-enter-active {
    transition: all 0.8s cubic-bezier(0.34, 1.56, 0.64, 1); /* Có độ nảy nhẹ (Bounce) */
}

.book-reveal-leave-active {
    transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
}

/* Trạng thái bắt đầu khi vào trang (Enter) */
.book-reveal-enter-from {
    opacity: 0;
    transform: translateY(30px) scale(0.95); /* Tụt xuống dưới và nhỏ lại */
    filter: blur(10px); /* Hiệu ứng mờ ảo */
}

/* Trạng thái kết thúc khi rời trang (Leave) */
.book-reveal-leave-to {
    opacity: 0;
    transform: scale(1.05); /* Phóng to nhẹ khi biến mất */
    filter: blur(5px);
}
</style>