<script setup>
import { ref, reactive } from 'vue';
import { useAuthStore } from '@/stores/auth';
import axios from 'axios';
import { Modal } from 'bootstrap';

const auth = useAuthStore();
const username = ref('');
const password = ref('');
const errorMsg = ref('');
const isLoading = ref(false);
const showPass = ref(false);

// --- LOGIC QUÊN MẬT KHẨU ---
const showResetModal = ref(false);
const resetForm = reactive({
    username: '',
    new_password: '',
    secret_key: ''
});
const resetMsg = ref({ type: '', text: '' });
const isResetting = ref(false);

// Hàm Đăng nhập
const handleLogin = async () => {
    isLoading.value = true;
    errorMsg.value = '';
    await new Promise(r => setTimeout(r, 500)); // Hiệu ứng mượt
    try {
        await auth.login(username.value, password.value);
    } catch (e) {
        errorMsg.value = "Thông tin đăng nhập không chính xác.";
    } finally {
        isLoading.value = false;
    }
};

// Hàm Xử lý Reset Password
const handleReset = async () => {
    if (!resetForm.username || !resetForm.new_password || !resetForm.secret_key) {
        resetMsg.value = { type: 'danger', text: 'Vui lòng điền đầy đủ thông tin.' };
        return;
    }

    isResetting.value = true;
    resetMsg.value = { type: '', text: '' };

    try {
        // Gọi API Backend
        await axios.post('http://localhost:8080/api/reset-password', resetForm);
        
        resetMsg.value = { type: 'success', text: 'Đổi mật khẩu thành công! Vui lòng đăng nhập lại.' };
        
        // Tự động đóng Modal sau 2 giây và điền username vào form login
        setTimeout(() => {
            showResetModal.value = false;
            username.value = resetForm.username;
            password.value = '';
            // Reset form
            Object.assign(resetForm, { username: '', new_password: '', secret_key: '' });
            resetMsg.value = { type: '', text: '' };
        }, 2000);

    } catch (error) {
        resetMsg.value = { type: 'danger', text: error.response?.data || "Lỗi kết nối Server!" };
    } finally {
        isResetting.value = false;
    }
};
</script>

<template>
    <div class="container-fluid vh-100 p-0 overflow-hidden">
        <div class="row h-100 g-0">
            
            <!-- CỘT TRÁI: HÌNH ẢNH -->
            <div class="col-lg-7 d-none d-lg-block position-relative bg-dark">
                <div class="bg-image"></div>
                <div class="bg-overlay"></div>
                <div class="position-absolute bottom-0 start-0 p-5 text-white z-2" style="max-width: 600px;">
                    <div class="mb-4"><i class="fa-solid fa-satellite-dish fa-3x text-primary mb-3"></i></div>
                    <h1 class="display-5 fw-bold mb-3">Global Fleet Intelligence</h1>
                    <p class="lead text-white-50">Hệ thống giám sát vận tải biển thế hệ mới. Kết nối dữ liệu vệ tinh Real-time và AI.</p>
                </div>
            </div>

            <!-- CỘT PHẢI: FORM LOGIN -->
            <div class="col-lg-5 col-md-12 bg-white d-flex align-items-center justify-content-center position-relative">
                <div class="w-100 p-5" style="max-width: 500px;">
                    <div class="mb-5">
                        <h2 class="fw-bold text-dark display-6">Welcome back</h2>
                        <p class="text-secondary">Please enter your details to sign in.</p>
                    </div>

                    <form @submit.prevent="handleLogin">
                        <div class="form-floating mb-3">
                            <input v-model="username" type="text" class="form-control" placeholder="name@example.com" required>
                            <label class="text-secondary"><i class="fa-solid fa-user me-2"></i>Username</label>
                        </div>

                        <div class="form-floating mb-3 position-relative">
                            <input v-model="password" :type="showPass ? 'text' : 'password'" class="form-control" placeholder="Password" required>
                            <label class="text-secondary"><i class="fa-solid fa-lock me-2"></i>Password</label>
                            <span @click="showPass = !showPass" class="position-absolute top-50 end-0 translate-middle-y me-3 cursor-pointer text-secondary">
                                <i :class="showPass ? 'fa-solid fa-eye-slash' : 'fa-solid fa-eye'"></i>
                            </span>
                        </div>

                        <div class="d-flex justify-content-between align-items-center mb-4">
                            <div class="form-check">
                                <input class="form-check-input" type="checkbox" id="rememberMe">
                                <label class="form-check-label text-secondary small" for="rememberMe">Remember me</label>
                            </div>
                            <!-- LINK QUÊN MẬT KHẨU: Mở Modal -->
                            <a href="#" @click.prevent="showResetModal = true" class="small text-decoration-none fw-bold">Forgot password?</a>
                        </div>

                        <button type="submit" class="btn btn-primary w-100 py-3 fw-bold shadow-sm d-flex align-items-center justify-content-center transition-all" :disabled="isLoading">
                            <span v-if="isLoading" class="spinner-border spinner-border-sm me-2"></span>
                            {{ isLoading ? 'Authenticating...' : 'Sign in to Dashboard' }}
                        </button>

                        <div v-if="errorMsg" class="alert alert-danger border-0 bg-danger bg-opacity-10 text-danger mt-4 d-flex align-items-center">
                            <i class="fa-solid fa-circle-exclamation me-2"></i> {{ errorMsg }}
                        </div>
                    </form>
                    
                    <div class="text-center mt-5"><p class="small text-muted">© 2025 Marine Pro. Enterprise Edition v2.0</p></div>
                </div>
            </div>
        </div>

        <!-- MODAL QUÊN MẬT KHẨU (Giao diện tối) -->
        <div v-if="showResetModal" class="modal-backdrop-custom d-flex align-items-center justify-content-center">
            <div class="card shadow-lg border-0 p-4" style="width: 400px; animation: slideDown 0.3s ease;">
                <div class="text-center mb-4">
                    <div class="bg-warning bg-opacity-10 text-warning rounded-circle d-inline-flex p-3 mb-2">
                        <i class="fa-solid fa-key fa-xl"></i>
                    </div>
                    <h5 class="fw-bold">Reset Password</h5>
                    <p class="small text-secondary">Enter your username and the system master key.</p>
                </div>

                <form @submit.prevent="handleReset">
                    <div class="mb-3">
                        <label class="form-label small fw-bold">Username</label>
                        <input v-model="resetForm.username" class="form-control" placeholder="e.g. admin" required>
                    </div>
                    <div class="mb-3">
                        <label class="form-label small fw-bold">New Password</label>
                        <input v-model="resetForm.new_password" type="password" class="form-control" placeholder="New secure password" required>
                    </div>
                    <div class="mb-4">
                        <label class="form-label small fw-bold text-danger">Master Secret Key</label>
                        <input v-model="resetForm.secret_key" type="password" class="form-control border-danger" placeholder="System Admin Code" required>
                        <div class="form-text small fst-italic">Default key is: <strong>marine_admin</strong></div>
                    </div>

                    <button type="submit" class="btn btn-dark w-100 py-2 fw-bold" :disabled="isResetting">
                        {{ isResetting ? 'Processing...' : 'Reset Password' }}
                    </button>
                    <button type="button" @click="showResetModal = false" class="btn btn-link text-decoration-none text-secondary w-100 mt-2 small">Cancel</button>
                    
                    <div v-if="resetMsg.text" :class="`alert alert-${resetMsg.type} mt-3 small text-center`">
                        {{ resetMsg.text }}
                    </div>
                </form>
            </div>
        </div>

    </div>
</template>

<style scoped>
/* Style cũ giữ nguyên */
.bg-image {
    position: absolute; top: 0; left: 0; width: 100%; height: 100%;
    background-image: url('https://images.unsplash.com/photo-1547895740-da9516c5e824?q=80&w=1974&auto=format&fit=crop'); 
    background-size: cover; background-position: center; z-index: 0;
    animation: zoomEffect 20s infinite alternate; 
}
.bg-overlay {
    position: absolute; top: 0; left: 0; width: 100%; height: 100%;
    background: linear-gradient(135deg, rgba(15, 23, 42, 0.85) 0%, rgba(30, 64, 175, 0.7) 100%);
    z-index: 1;
}
.form-control { border: 2px solid #f1f5f9; border-radius: 12px; padding-left: 15px; transition: all 0.2s; }
.form-control:focus { border-color: var(--bs-primary); box-shadow: 0 0 0 4px rgba(37, 99, 235, 0.1); }
.btn-primary { border-radius: 12px; background: #2563eb; border: none; transition: all 0.2s; }
.btn-primary:hover { background: #1d4ed8; transform: translateY(-2px); box-shadow: 0 10px 20px -5px rgba(37, 99, 235, 0.4); }
@keyframes zoomEffect { from { transform: scale(1); } to { transform: scale(1.1); } }
.cursor-pointer { cursor: pointer; }

/* STYLE CHO MODAL RIÊNG */
.modal-backdrop-custom {
    position: fixed; top: 0; left: 0; width: 100%; height: 100%;
    background: rgba(0,0,0,0.6); backdrop-filter: blur(4px);
    z-index: 9999;
}
@keyframes slideDown { from { transform: translateY(-20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
</style>