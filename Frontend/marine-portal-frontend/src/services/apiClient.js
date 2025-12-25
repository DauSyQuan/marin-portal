import axios from "axios";

// Dùng biến môi trường Vite
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "http://localhost:8080",
  timeout: 15000,
});

// Gắn token tự động
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle lỗi tập trung
apiClient.interceptors.response.use(
  (res) => res,
  (err) => {
    // Nếu backend trả 401 => token invalid/expired
    if (err?.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      // redirect nhẹ nhàng (tránh import router trong file này)
      window.location.href = "/login";
    }
    return Promise.reject(err);
  }
);

export default apiClient;
