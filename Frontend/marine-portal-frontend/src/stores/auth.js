import { defineStore } from 'pinia'
import axios from 'axios'
import router from '@/router'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: JSON.parse(localStorage.getItem('user')) || null,
        token: localStorage.getItem('token') || null
    }),
    getters: {
        isAuthenticated: (state) => !!state.token
    },
    actions: {
        async login(username, password) {
            try {
                // Gọi API Go
                const res = await axios.post('http://localhost:8080/api/login', { username, password });
                
                // Lưu vào State và LocalStorage
                this.token = res.data.token;
                this.user = { username: res.data.username, role: res.data.role };
                
                localStorage.setItem('token', this.token);
                localStorage.setItem('user', JSON.stringify(this.user));
                
                // Chuyển hướng vào Dashboard
                router.push('/');
            } catch (error) {
                throw new Error("Login failed! Check your credentials.");
            }
        },
        logout() {
            this.token = null;
            this.user = null;
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            router.push('/login');
        }
    }
})