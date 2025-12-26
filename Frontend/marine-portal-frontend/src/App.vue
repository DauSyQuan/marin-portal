<script setup>
import TheSidebar from './components/TheSidebar.vue'
import { useRoute } from 'vue-router'
const route = useRoute()
</script>

<template>
  <div class="app-layout">
    <!-- Sidebar -->
    <div v-if="route.path !== '/login'" class="sidebar-wrapper glass-sidebar">
      <TheSidebar />
    </div>

    <!-- Nội dung chính -->
    <div class="main-wrapper" :class="{ 'login-mode': route.path === '/login' }">
      <div class="main-content" :class="{ 'login-content': route.path === '/login' }">
        <router-view v-slot="{ Component }">
          <transition name="pure-fade">
            <component :is="Component" :key="route.fullPath" />
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<style>
/* Layout */
body { overflow: hidden; }
.app-layout { display: flex; width: 100vw; height: 100vh; }

/* Sidebar Kính */
.sidebar-wrapper {
    width: 260px; height: 100%; position: fixed; z-index: 20; left: 0; top: 0;
    background: rgba(15, 23, 42, 0.6); 
    backdrop-filter: blur(10px);
    border-right: 1px solid rgba(255,255,255,0.1);
}

/* ✅ Language switch: né header dashboard */
.top-lang {
  position: fixed;
  top: 72px;       /* ✅ đẩy xuống dưới header */
  right: 26px;     /* ✅ né cụm nút bên phải */
  z-index: 60;
}

.main-wrapper {
    flex-grow: 1; margin-left: 260px;
    padding: 20px;
    height: 100vh;
}

.main-content {
    background-color: transparent; 
    width: 100%; height: 100%;
    min-height: calc(100vh - 40px);
    border-radius: 32px;
    position: relative;
    overflow-y: auto;
    overflow-x: hidden;
}

.main-wrapper.login-mode { 
    margin-left: 0; 
    padding: 0; 
}

/* ✅ Login mode: full screen, không bo góc */
.main-content.login-content {
    border-radius: 0 !important;
    min-height: 100vh !important;
    height: 100vh !important;
    overflow: hidden !important;
}

/* Scrollbar */
::-webkit-scrollbar { width: 6px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.2); border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: rgba(255,255,255,0.4); }

/* Animation Fade */
.pure-fade-enter-active, .pure-fade-leave-active { 
    transition: opacity 0.5s ease; position: absolute; 
    top: 0; left: 0; width: 100%; 
}
.pure-fade-enter-from { opacity: 0; z-index: 2; }
.pure-fade-leave-to { opacity: 0; z-index: 0; }
</style>
