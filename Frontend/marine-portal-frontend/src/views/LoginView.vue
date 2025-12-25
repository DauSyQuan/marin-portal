<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const username = ref('');
const password = ref('');
const errorMsg = ref('');
const isLoading = ref(false);

const handleLogin = async () => {
    isLoading.value = true;
    errorMsg.value = '';
    try {
        await auth.login(username.value, password.value);
    } catch (e) {
        errorMsg.value = "Incorrect Username or Password!";
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div class="vh-100 d-flex align-items-center justify-content-center bg-light">
        <div class="card p-5 border-0 shadow-lg" style="width: 400px;">
            <div class="text-center mb-5">
                <i class="fa-solid fa-earth-asia fa-3x text-primary mb-3"></i>
                <h3 class="fw-bold">MARINE PORTAL</h3>
                <p class="text-secondary small">SECURE LOGIN</p>
            </div>
            <form @submit.prevent="handleLogin">
                <div class="mb-3">
                    <input v-model="username" type="text" class="form-control py-2" placeholder="Username (admin)" required>
                </div>
                <div class="mb-4">
                    <input v-model="password" type="password" class="form-control py-2" placeholder="Password (123)" required>
                </div>
                <button type="submit" class="btn btn-primary w-100 py-2 fw-bold" :disabled="isLoading">
                    {{ isLoading ? 'Signing in...' : 'Sign In' }}
                </button>
                <div v-if="errorMsg" class="text-danger text-center mt-3 small">{{ errorMsg }}</div>
            </form>
        </div>
    </div>
</template>