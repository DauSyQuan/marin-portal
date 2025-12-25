<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const auth = useAuthStore();
const router = useRouter();

const username = ref('');
const password = ref('');
const errorMsg = ref('');
const isLoading = ref(false);

const handleLogin = async () => {
    isLoading.value = true;
    errorMsg.value = '';
    
    // Giả lập độ trễ nhỏ để tạo cảm giác xử lý mượt mà
    await new Promise(r => setTimeout(r, 600));

    const success = await auth.login(username.value, password.value);
    
    if (success) {
        router.push('/');
    } else {
        errorMsg.value = auth.error || "Authentication failed.";
    }
    isLoading.value = false;
};
</script>

<template>
    <div class="login-wrapper">
        <div class="login-container">
            <!-- Left Side: Visual -->
            <div class="visual-side">
                <div class="visual-content">
                    <div class="logo-box">
                        <i class="fa-solid fa-anchor fa-2x"></i>
                    </div>
                    <h1>Marine Portal</h1>
                    <p>Global Fleet Management System</p>
                    <div class="illustration-placeholder">
                        <i class="fa-solid fa-earth-americas fa-10x opacity-25"></i>
                    </div>
                </div>
                <div class="overlay"></div>
            </div>

            <!-- Right Side: Form -->
            <div class="form-side">
                <div class="form-header">
                    <h2>Welcome Back</h2>
                    <p class="text-secondary">Please enter your credentials to access the command center.</p>
                </div>

                <form @submit.prevent="handleLogin" class="login-form">
                    <div class="form-group">
                        <label>Username</label>
                        <div class="input-wrapper">
                            <i class="fa-solid fa-user input-icon"></i>
                            <input v-model="username" type="text" placeholder="e.g. admin" required>
                        </div>
                    </div>

                    <div class="form-group">
                        <label>Password</label>
                        <div class="input-wrapper">
                            <i class="fa-solid fa-lock input-icon"></i>
                            <input v-model="password" type="password" placeholder="••••••" required>
                        </div>
                    </div>

                    <div v-if="errorMsg" class="error-alert">
                        <i class="fa-solid fa-circle-exclamation me-2"></i> {{ errorMsg }}
                    </div>

                    <button type="submit" class="btn-login" :disabled="isLoading">
                        <span v-if="!isLoading">Sign In</span>
                        <div v-else class="spinner"></div>
                    </button>
                </form>

                <div class="form-footer">
                    <span>Protected by MarineGuard™ Security</span>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.login-wrapper {
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f1f5f9;
}

.login-container {
    display: flex;
    width: 1000px;
    height: 600px;
    background: white;
    border-radius: 24px;
    box-shadow: 0 20px 60px rgba(15, 23, 42, 0.1);
    overflow: hidden;
}

/* Visual Side */
.visual-side {
    flex: 1;
    background: url('https://images.unsplash.com/photo-1549495163-95d35d9472eb?q=80&w=1000&auto=format&fit=crop') center/cover;
    position: relative;
    color: white;
    display: flex;
    flex-direction: column;
    justify-content: center;
    padding: 40px;
}

.overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, rgba(15, 23, 42, 0.9), rgba(59, 130, 246, 0.7));
    z-index: 1;
}

.visual-content {
    position: relative;
    z-index: 2;
}

.logo-box {
    width: 60px;
    height: 60px;
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(10px);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 24px;
}

.visual-side h1 { font-weight: 800; font-size: 2.5rem; margin-bottom: 8px; }
.visual-side p { font-size: 1.1rem; opacity: 0.9; font-weight: 300; }

.illustration-placeholder {
    margin-top: 60px;
    display: flex;
    justify-content: center;
    opacity: 0.6;
}

/* Form Side */
.form-side {
    flex: 1;
    padding: 60px;
    display: flex;
    flex-direction: column;
    justify-content: center;
}

.form-header h2 { font-weight: 800; color: #0f172a; margin-bottom: 8px; }

.form-group { margin-bottom: 24px; }
.form-group label { display: block; font-size: 0.85rem; font-weight: 600; color: #64748b; margin-bottom: 8px; }

.input-wrapper { position: relative; }
.input-icon { position: absolute; left: 16px; top: 50%; transform: translateY(-50%); color: #94a3b8; }
.input-wrapper input {
    width: 100%;
    padding: 14px 14px 14px 44px;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    font-size: 1rem;
    transition: all 0.2s;
    background: #f8fafc;
}
.input-wrapper input:focus {
    outline: none;
    border-color: #3b82f6;
    background: white;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.btn-login {
    width: 100%;
    padding: 14px;
    background: #0f172a;
    color: white;
    border: none;
    border-radius: 12px;
    font-weight: 600;
    font-size: 1rem;
    cursor: pointer;
    transition: background 0.2s;
    display: flex;
    justify-content: center;
    align-items: center;
}
.btn-login:hover { background: #334155; }
.btn-login:disabled { opacity: 0.7; cursor: not-allowed; }

.error-alert {
    background: #fef2f2;
    color: #ef4444;
    padding: 12px;
    border-radius: 8px;
    font-size: 0.9rem;
    margin-bottom: 20px;
    display: flex;
    align-items: center;
}

.form-footer { margin-top: auto; text-align: center; font-size: 0.75rem; color: #cbd5e1; }

.spinner {
    width: 20px; height: 20px;
    border: 2px solid rgba(255,255,255,0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>