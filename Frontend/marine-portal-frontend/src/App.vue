<script setup>
import TheSidebar from './components/TheSidebar.vue'
import { useRoute } from 'vue-router'

const route = useRoute()
</script>

<template>
  <div class="app-layout">
    <!-- Sidebar chỉ hiện khi không phải trang login -->
    <TheSidebar v-if="route.path !== '/login'" />

    <div class="main-content" :style="{ marginLeft: route.path === '/login' ? '0' : '260px' }">
      
      <!-- THÊM HIỆU ỨNG TRANSITION Ở ĐÂY -->
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>

    </div>
  </div>
</template>

<style>
/* Reset mặc định */
.main-content { 
    padding: 0; 
    min-height: 100vh;
    background-color: var(--bg-body); /* Đảm bảo màu nền đồng nhất */
}

/* CSS cho hiệu ứng mờ dần (Fade) */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>