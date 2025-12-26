import { defineStore } from "pinia";
import { ref, computed } from "vue";

export const useI18nStore = defineStore("i18n", () => {
  const lang = ref(localStorage.getItem("lang") || "vi");

  const dict = {
    vi: {
      login_title: "MARINE PRO",
      login_subtitle: "HỆ THỐNG CHỈ HUY HẠM ĐỘI",
      username: "Tên người dùng",
      password: "Mật khẩu",
      remember: "Ghi nhớ",
      forgot: "Quên mật khẩu?",
      login_btn: "ĐĂNG NHẬP",
      connecting: "ĐANG KẾT NỐI...",

      main_menu: "Menu chính",
      management: "Quản lý",
      system: "Hệ thống",
      fleet_overview: "Tổng quan hạm đội",
      analytics: "Báo cáo phân tích",
      monthly_usage: "Sử dụng hàng tháng",
      online_users: "Người dùng online",
      vouchers: "Voucher",
      bandwidth: "Gói băng thông",
      settings: "Cấu hình",
      sign_out: "Đăng xuất",
    },
    en: {
      login_title: "MARINE PRO",
      login_subtitle: "FLEET COMMAND SYSTEM",
      username: "Username",
      password: "Password",
      remember: "Remember me",
      forgot: "Forgot password?",
      login_btn: "SIGN IN",
      connecting: "CONNECTING...",

      main_menu: "Main Menu",
      management: "Management",
      system: "System",
      fleet_overview: "Fleet Overview",
      analytics: "Analytics Report",
      monthly_usage: "Monthly Usage",
      online_users: "Online Users",
      vouchers: "Vouchers",
      bandwidth: "Bandwidth Plans",
      settings: "Configuration",
      sign_out: "Sign out",
    },
  };

  const t = (key) => dict[lang.value]?.[key] || key;

  const setLang = (newLang) => {
    lang.value = newLang;
    localStorage.setItem("lang", newLang);
  };

  return {
    lang,
    t,
    setLang,
  };
});
