import axios from "axios";
import router from "@/router";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080",
  timeout: 15000,
});

// ✅ Attach token automatically
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) config.headers.Authorization = `Bearer ${token}`;
    return config;
  },
  (error) => Promise.reject(error)
);

// ✅ Handle 401 globally
api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const status = error?.response?.status;

    if (status === 401) {
      // clear auth
      localStorage.removeItem("token");
      localStorage.removeItem("user");

      // avoid infinite loop if already at /login
      if (router.currentRoute.value.path !== "/login") {
        router.push("/login");
      }
    }

    return Promise.reject(error);
  }
);

export default api;
