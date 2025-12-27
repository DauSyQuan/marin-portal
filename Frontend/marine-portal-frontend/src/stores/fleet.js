import { defineStore } from 'pinia'
import api from '@/services/api'   // ✅ dùng instance có token

export const useFleetStore = defineStore('fleet', {
  state: () => ({
    ships: [],
    loading: false,
    filterText: ''
  }),

  getters: {
    stats(state) {
      const total = state.ships.length;
      const online = state.ships.filter(s => s.status === 'Online').length;
      const offline = state.ships.filter(s => s.status === 'Offline').length;
      const warning = state.ships.filter(s => s.status === 'Warning' || s.status === 'Blockage').length;
      const health = total > 0 ? ((online + (warning * 0.5)) / total * 100).toFixed(1) : 0;

      return [
        { label: 'Total Fleet', value: total, color: 'primary', icon: 'fa-solid fa-ship' },
        { label: 'System Health', value: health + '%', color: 'success', icon: 'fa-solid fa-heart-pulse' },
        { label: 'Critical Alerts', value: offline, color: 'danger', icon: 'fa-solid fa-triangle-exclamation' },
        { label: 'Warnings', value: warning, color: 'warning', icon: 'fa-solid fa-circle-exclamation' }
      ];
    },

    filteredShips(state) {
      if (!state.filterText) return state.ships;
      const term = state.filterText.toLowerCase();
      return state.ships.filter(s =>
        s.name.toLowerCase().includes(term) ||
        s.id.toLowerCase().includes(term) ||
        s.company.toLowerCase().includes(term)
      );
    }
  },

  actions: {
    // --- QUẢN LÝ TÀU ---
    async fetchFleet() {
      this.loading = true;
      try {
        // ✅ interceptor trong api.js sẽ gắn Authorization: Bearer <token>
        const response = await api.get('/api/ships');
        this.ships = response.data || [];
      } catch (error) {
        console.error("Lỗi tải tàu:", error);
      } finally {
        this.loading = false;
      }
    },

    async addShip(newShipData) {
      try {
        const response = await api.post('/api/ships', newShipData);
        if (response.data) {
          this.ships.unshift(response.data);
          return true;
        }
        return false;
      } catch (error) {
        alert("Lỗi thêm tàu: " + (error.response?.data?.error || error.message));
        return false;
      }
    },

    // --- QUẢN LÝ THỦY THỦ (CREW / ICT USER) ---

    // 1. Lấy danh sách
    async fetchCrew(shipId) {
      try {
        const response = await api.get(`/api/ships/${shipId}/crew`);
        return response.data || [];
      } catch (error) {
        console.error("Lỗi tải Crew:", error);
        return [];
      }
    },

    // 2. Thêm mới
    async addCrew(crewData) {
      try {
        await api.post('/api/crew', crewData);
        return true;
      } catch (error) {
        alert("Lỗi thêm: " + (error.response?.data?.error || error.message));
        return false;
      }
    },

    // 3. Cập nhật (Sửa)
    async updateCrew(id, updatedData) {
      try {
        await api.put(`/api/crew/${id}`, updatedData);
        return true;
      } catch (error) {
        alert("Lỗi cập nhật: " + (error.response?.data?.error || error.message));
        return false;
      }
    },

    // 4. Xóa
    async deleteCrew(id) {
      try {
        if (!confirm("Bạn có chắc chắn muốn xóa User này không?")) return false;
        await api.delete(`/api/crew/${id}`);
        return true;
      } catch (error) {
        alert("Lỗi xóa: " + (error.response?.data?.error || error.message));
        return false;
      }
    }
  }
})
