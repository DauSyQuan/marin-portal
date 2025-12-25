import { defineStore } from "pinia";
import apiClient from "@/services/apiClient";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    token: localStorage.getItem("token") || null,

    // ✅ nếu không có user thì mặc định null, KHÔNG parse undefined
    user: (() => {
      try {
        const raw = localStorage.getItem("user");
        return raw ? JSON.parse(raw) : null;
      } catch {
        localStorage.removeItem("user");
        return null;
      }
    })(),

    loading: false,
    error: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
  },

  actions: {
    async login(username, password) {
      this.loading = true;
      this.error = null;

      try {
        const res = await apiClient.post("/api/login", { username, password });

        // ✅ backend trả token, username, role
        this.token = res.data.token;
        this.user = {
          username: res.data.username,
          role: res.data.role,
        };

        localStorage.setItem("token", this.token);
        localStorage.setItem("user", JSON.stringify(this.user));

        return true;
      } catch (err) {
        const status = err?.response?.status;

        if (status === 401) this.error = "Sai tài khoản hoặc mật khẩu.";
        else if (!err?.response) this.error = "Không kết nối được server.";
        else this.error = "Đăng nhập thất bại. Vui lòng thử lại.";

        return false;
      } finally {
        this.loading = false;
      }
    },

    logout() {
      this.token = null;
      this.user = null;
      localStorage.removeItem("token");
      localStorage.removeItem("user");
    },
  },
});
