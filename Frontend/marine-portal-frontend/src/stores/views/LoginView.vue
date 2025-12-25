<script setup>
import { ref, reactive } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { Modal } from 'bootstrap';

const auth = useAuthStore();
const username = ref('');
const password = ref('');
const errorMsg = ref('');
const isLoading = ref(false);
const showPass = ref(false);

// Logic Reset Password
const showResetModal = ref(false);
const resetForm = reactive({ username: '', new_password: '', secret_key: '' });
const resetMsg = ref({ type: '', text: '' });
const isResetting = ref(false);

const handleLogin = async () => {
    isLoading.value = true;
    errorMsg.value = '';
    await new Promise(r => setTimeout(r, 800)); // Delay giả lập
    try {
        await auth.login(username.value, password.value);
    } catch (e) {
        errorMsg.value = "Thông tin đăng nhập không chính xác.";
    } finally {
        isLoading.value = false;
    }
};

// Logic Reset (Demo)
const handleReset = async () => {
    // ... (Giữ nguyên logic cũ)
    showResetModal.value = false;
};
</script>

<template>
    <!-- WRAPPER CHÍNH: FULL MÀN HÌNH -->
    <div class="login-wrapper">
        
        <!-- 1. BACKGROUND HÌNH ẢNH (ĐÃ THAY LINK CỦA BẠN) -->
        <div class="bg-image"></div>
        
        <!-- Lớp phủ nhẹ để làm dịu nền -->
        <div class="bg-overlay"></div>

        <!-- 2. TẤM KÍNH MỜ (LIQUID GLASS CARD) -->
        <div class="glass-card">
            
            <!-- Logo & Header -->
            <div class="text-center mb-5">
                <div class="logo-glow mb-3">
                    <i class="fa-solid fa-anchor fa-2x text-white"></i>
                </div>
                <h2 class="fw-bold text-white tracking-wide">MARINE PRO</h2>
                <p class="text-white-50 small text-uppercase letter-spacing-2">Fleet Command Access</p>
            </div>

            <!-- Form Login -->
            <form @submit.prevent="handleLogin">
                
                <!-- Input Username -->
                <div class="glass-input-group mb-3">
                    <div class="icon"><i class="fa-solid fa-user text-white-50"></i></div>
                    <input v-model="username" type="text" class="glass-input" placeholder="Username" required>
                </div>

                <!-- Input Password -->
                <div class="glass-input-group mb-4">
                    <div class="icon"><i class="fa-solid fa-lock text-white-50"></i></div>
                    <input v-model="password" :type="showPass ? 'text' : 'password'" class="glass-input" placeholder="Password" required>
                    <span class="eye-btn" @click="showPass = !showPass">
                        <i :class="showPass ? 'fa-solid fa-eye-slash' : 'fa-solid fa-eye'"></i>
                    </span>
                </div>

                <!-- Options -->
                <div class="d-flex justify-content-between align-items-center mb-4 text-white small">
                    <div class="form-check custom-checkbox">
                        <input class="form-check-input" type="checkbox" id="remember">
                        <label class="form-check-label text-white-50" for="remember">Remember me</label>
                    </div>
                    <a href="#" @click.prevent="showResetModal = true" class="text-info text-decoration-none fw-bold glow-link">Forgot password?</a>
                </div>

                <!-- Button -->
                <button type="submit" class="btn-liquid w-100 py-3 fw-bold text-white" :disabled="isLoading">
                    <span v-if="isLoading" class="spinner-border spinner-border-sm me-2"></span>
                    {{ isLoading ? 'CONNECTING...' : 'LOGIN DASHBOARD' }}
                </button>

                <!-- Error -->
                <div v-if="errorMsg" class="glass-alert mt-4">
                    <i class="fa-solid fa-triangle-exclamation me-2"></i> {{ errorMsg }}
                </div>
            </form>

            <div class="text-center mt-5 pt-4 border-top border-white border-opacity-10">
                <small class="text-white-50">Authorized Personnel Only</small>
            </div>
        </div>

        <!-- MODAL RESET (Glass Style) -->
        <div v-if="showResetModal" class="modal-backdrop-glass d-flex align-items-center justify-content-center">
            <div class="glass-card p-4" style="width: 400px; animation: popIn 0.3s ease;">
                <h5 class="text-white fw-bold mb-3 text-center">Reset Access</h5>
                <input v-model="resetForm.username" class="glass-input-simple mb-3" placeholder="Username">
                <input v-model="resetForm.secret_key" type="password" class="glass-input-simple mb-4" placeholder="Master Key">
                <button class="btn-liquid w-100 mb-2" @click="handleReset">Confirm Reset</button>
                <button class="btn btn-link text-white-50 w-100 text-decoration-none small" @click="showResetModal = false">Cancel</button>
            </div>
        </div>

    </div>
</template>

<style scoped>
/* --- 1. BACKGROUND --- */
.login-wrapper {
    position: relative;
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    background-color: #000;
}

.bg-image {
    position: absolute; top: 0; left: 0; width: 100%; height: 100%;
    /* LINK HÌNH BẠN GỬI */
    background-image: url('https://c0.wallpaperflare.com/preview/888/665/86/background-blue-calm-gradients.jpg');
    background-size: cover; background-position: center;
    z-index: 0;
    /* Hiệu ứng trôi nhẹ */
    animation: drift 60s infinite linear; 
}

.bg-overlay {
    position: absolute; top: 0; left: 0; width: 100%; height: 100%;
    /* Gradient tối giúp chữ nổi bật */
    background: radial-gradient(circle at center, rgba(0,0,0,0.2) 0%, rgba(0,0,0,0.6) 100%);
    z-index: 1;
}

/* --- 2. GLASS CARD (HIỆU ỨNG KÍNH LỎNG) --- */
.glass-card {
    position: relative;
    z-index: 10;
    width: 100%;
    max-width: 450px;
    padding: 50px 40px;
    border-radius: 24px;
    
    /* CÔNG THỨC GLASSMORPHISM CHUẨN */
    background: rgba(255, 255, 255, 0.05); /* Rất trong suốt */
    backdrop-filter: blur(16px) saturate(180%); /* Mờ nền phía sau */
    -webkit-backdrop-filter: blur(16px) saturate(180%);
    
    border: 1px solid rgba(255, 255, 255, 0.15); /* Viền kính sáng */
    box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.5); /* Bóng đổ sâu */
}

/* --- 3. LOGO GLOW --- */
.logo-glow {
    width: 80px; height: 80px;
    margin: 0 auto;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 50%;
    display: flex; align-items: center; justify-content: center;
    border: 1px solid rgba(255,255,255,0.2);
    box-shadow: 0 0 30px rgba(59, 130, 246, 0.4); /* Phát sáng xanh */
}

/* --- 4. INPUTS --- */
.glass-input-group {
    display: flex; align-items: center;
    background: rgba(0, 0, 0, 0.25); /* Nền đen mờ */
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.05);
    transition: all 0.3s;
}
.glass-input-group:focus-within {
    border-color: rgba(255, 255, 255, 0.5);
    box-shadow: 0 0 15px rgba(255, 255, 255, 0.1);
    background: rgba(0, 0, 0, 0.4);
}

.icon { padding: 0 15px; }
.glass-input {
    width: 100%; background: transparent; border: none; color: white;
    padding: 16px 10px 16px 0; outline: none; font-size: 1rem;
}
.glass-input::placeholder { color: rgba(255, 255, 255, 0.3); }

.eye-btn { padding: 0 15px; cursor: pointer; color: rgba(255,255,255,0.5); }
.eye-btn:hover { color: white; }

/* Input Simple (Cho Modal) */
.glass-input-simple {
    width: 100%; background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.2);
    color: white; padding: 12px; border-radius: 8px; outline: none;
}
.glass-input-simple:focus { background: rgba(255,255,255,0.2); }

/* --- 5. BUTTON LIQUID --- */
.btn-liquid {
    background: linear-gradient(90deg, #2563eb, #3b82f6);
    border: none; border-radius: 12px;
    transition: all 0.3s; letter-spacing: 1px;
    box-shadow: 0 4px 15px rgba(37, 99, 235, 0.4);
}
.btn-liquid:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(37, 99, 235, 0.6);
    filter: brightness(1.2);
}

/* --- 6. EXTRAS --- */
.tracking-wide { letter-spacing: 2px; }
.glow-link:hover { text-shadow: 0 0 10px rgba(59, 130, 246, 0.8); }

.glass-alert {
    background: rgba(220, 38, 38, 0.2); border: 1px solid rgba(220, 38, 38, 0.5);
    color: #fca5a5; padding: 12px; border-radius: 8px; font-size: 0.9rem;
}

.modal-backdrop-glass {
    position: fixed; top: 0; left: 0; width: 100%; height: 100%;
    background: rgba(0,0,0,0.7); backdrop-filter: blur(8px); z-index: 9999;
}

@keyframes drift { from { transform: scale(1.1); } to { transform: scale(1); } }
@keyframes popIn { from { transform: scale(0.9); opacity: 0; } to { transform: scale(1); opacity: 1; } }
</style>