import { defineStore } from "pinia";
import api from "@/services/api";
import router from "@/router";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    user: JSON.parse(localStorage.getItem("user")) || null,
    token: localStorage.getItem("token") || null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
  },

  actions: {
    async login(username, password) {
      try {
        const res = await api.post("/api/login", { username, password });

        // ✅ token bắt buộc
        if (!res.data?.token) {
          throw new Error("Missing token in login response");
        }

        this.token = res.data.token;

        // ✅ user có thể backend trả khác -> fallback
        this.user = {
          username: res.data.username || res.data.user?.username || username,
          role: res.data.role || res.data.user?.role || "user",
        };

        localStorage.setItem("token", this.token);
        localStorage.setItem("user", JSON.stringify(this.user));

        router.push("/");
      } catch (error) {
        console.error("Login error:", error);
        throw new Error("Login failed! Check your credentials.");
      }
    },

    logout() {
      this.token = null;
      this.user = null;

      localStorage.removeItem("token");
      localStorage.removeItem("user");

      router.push("/login");
    },
  },
});
