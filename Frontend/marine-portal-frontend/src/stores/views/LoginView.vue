<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const username = ref('');
const password = ref('');
const errorMsg = ref('');
const isLoading = ref(false);
const showPass = ref(false);

const handleLogin = async () => {
    if (isLoading.value) return; // ✅ chặn spam click

    isLoading.value = true;
    errorMsg.value = '';

    // ✅ fake delay để UI đẹp
    await new Promise(r => setTimeout(r, 800));

    try {
        await auth.login(username.value, password.value);
    } catch (e) {
        errorMsg.value = "Sai tài khoản hoặc mật khẩu";
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <!-- WRAPPER FULL MÀN HÌNH -->
    <div class="login-container">
        
        <!-- 1. HÌNH NỀN -->
        <div class="bg-image"></div>
        
        <!-- 2. FORM ĐĂNG NHẬP -->
        <div class="glass-form">
            
            <!-- Logo -->
            <div class="text-center mb-4">
                <div class="logo-box mb-3">
                    <i class="fa-solid fa-anchor fa-2x text-white"></i>
                </div>
                <h3 class="fw-bold text-white tracking-wide">MARINE PRO</h3>
                <p class="text-white-50 small text-uppercase">Fleet Command System</p>
            </div>

            <form @submit.prevent="handleLogin">
                <div class="mb-3">
                    <input v-model="username" class="glass-input" placeholder="Username" required :disabled="isLoading">
                </div>

                <div class="mb-4 position-relative">
                    <input
                        v-model="password"
                        :type="showPass ? 'text' : 'password'"
                        class="glass-input"
                        placeholder="Password"
                        required
                        :disabled="isLoading"
                    >
                    <span class="eye-icon" @click="showPass = !showPass">
                        <i :class="showPass ? 'fa-solid fa-eye-slash' : 'fa-solid fa-eye'"></i>
                    </span>
                </div>

                <div class="d-flex justify-content-between mb-4 text-white-50 small px-1">
                    <div class="form-check">
                        <input class="form-check-input" type="checkbox" id="rem" :disabled="isLoading">
                        <label class="form-check-label" for="rem">Ghi nhớ</label>
                    </div>
                    <a href="#" class="text-white text-decoration-none">Quên mật khẩu?</a>
                </div>

                <!-- ✅ BUTTON LOADING RING XOAY -->
                <button
                  class="btn-login w-100"
                  :class="{ 'btn-rotating': isLoading }"
                  :disabled="isLoading"
                >
                    <span v-if="isLoading" class="d-flex align-items-center justify-content-center gap-2">
                        <i class="fa-solid fa-spinner fa-spin"></i>
                        Đang kết nối...
                    </span>
                    <span v-else>
                        ĐĂNG NHẬP
                    </span>
                </button>

                <div v-if="errorMsg" class="alert-glass mt-3">
                    <i class="fa-solid fa-circle-exclamation me-2"></i> {{ errorMsg }}
                </div>
            </form>
        </div>

    </div>
</template>

<style scoped>
/* BACKGROUND */
.login-container {
    position: relative; width: 100vw; height: 100vh;
    display: flex; align-items: center; justify-content: center;
    background-color: #000; overflow: hidden;
}
.bg-image {
    position: absolute; inset: 0;
    background-image: url('https://c0.wallpaperflare.com/preview/888/665/86/background-blue-calm-gradients.jpg');
    background-size: cover; background-position: center;
    filter: blur(0px);
    z-index: 0;
}

/* GLASS FORM */
.glass-form {
    position: relative; z-index: 10;
    width: 400px; padding: 40px;
    border-radius: 20px;
    background: rgba(255, 255, 255, 0.1); 
    backdrop-filter: blur(15px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

/* INPUT */
.glass-input {
    width: 100%; padding: 14px 16px;
    background: rgba(0, 0, 0, 0.25);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 10px; color: white; outline: none;
    transition: all 0.3s;
}
.glass-input:focus {
    background: rgba(0, 0, 0, 0.4);
    border-color: rgba(255, 255, 255, 0.5);
}
.glass-input::placeholder { color: rgba(255, 255, 255, 0.4); }

/* BUTTON */
.btn-login {
    padding: 14px; border: none; border-radius: 12px;
    background: linear-gradient(135deg, #3b82f6, #2563eb);
    color: white; font-weight: bold; letter-spacing: 1px;
    transition: all 0.3s;
    box-shadow: 0 4px 15px rgba(37, 99, 235, 0.4);
    display: flex;
    justify-content: center;
    align-items: center;
    position: relative;
    overflow: hidden;
}
.btn-login:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(37, 99, 235, 0.6);
}
.btn-login:disabled {
    opacity: 0.75;
    cursor: not-allowed;
    transform: none;
}

/* ✅ LOADING RING XOAY QUANH NÚT */
.btn-rotating::after {
    content: "";
    position: absolute;
    inset: -55%;
    border-radius: 999px;
    border: 4px solid rgba(255, 255, 255, 0.14);
    border-top-color: rgba(255, 255, 255, 0.6);
    animation: ringSpin 1s linear infinite;
    pointer-events: none;
}

/* Ring spin animation */
@keyframes ringSpin {
    to { transform: rotate(360deg); }
}

/* EXTRAS */
.logo-box {
    width: 60px; height: 60px; margin: 0 auto;
    background: rgba(255,255,255,0.1); border-radius: 50%;
    display: flex; align-items: center; justify-content: center;
    box-shadow: 0 0 20px rgba(59, 130, 246, 0.5);
}
.eye-icon {
    position: absolute; right: 15px; top: 15px;
    color: rgba(255,255,255,0.5); cursor: pointer;
}
.alert-glass {
    background: rgba(220, 38, 38, 0.25); color: #ffadad;
    padding: 10px; border-radius: 8px; font-size: 0.9rem;
    border: 1px solid rgba(220, 38, 38, 0.4);
}
.tracking-wide { letter-spacing: 2px; }
</style>
