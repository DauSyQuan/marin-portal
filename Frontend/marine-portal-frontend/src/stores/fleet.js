import { defineStore } from "pinia";
import apiClient from "@/services/apiClient";

export const useFleetStore = defineStore("fleet", {
  state: () => ({
    ships: [],
    loading: false,
    error: null,
    searchQuery: "",
    addLoading: false,
    addError: null,
  }),

  getters: {
    filteredShips(state) {
      const q = state.searchQuery.trim().toLowerCase();
      if (!q) return state.ships;

      // Tối ưu: hạn chế toLowerCase nhiều lần
      return state.ships.filter((s) => (s._searchText || "").includes(q));
    },

    stats(state) {
      // Tối ưu: duyệt 1 lần
      let total = state.ships.length;
      let online = 0, offline = 0, warning = 0;

      for (const s of state.ships) {
        if (s.status === "Online") online++;
        else if (s.status === "Offline") offline++;
        else if (s.status === "Warning") warning++;
      }

      return { total, online, offline, warning };
    },
  },

  actions: {
    _decorateShips(list) {
      // tạo searchable text 1 lần
      return list.map((s) => ({
        ...s,
        _searchText: `${s.name || ""} ${s.id || ""} ${s.company || ""}`.toLowerCase(),
      }));
    },

    async fetchFleet() {
      this.loading = true;
      this.error = null;

      try {
        const res = await apiClient.get("/api/ships");
        this.ships = this._decorateShips(res.data);
      } catch (err) {
        this.error = err?.response?.data?.message || "Không tải được danh sách tàu.";
      } finally {
        this.loading = false;
      }
    },

    async addShip(newShip) {
      this.addLoading = true;
      this.addError = null;

      try {
        const res = await apiClient.post("/api/ships", newShip);
        this.ships.unshift({
          ...res.data,
          _searchText: `${res.data.name || ""} ${res.data.id || ""} ${res.data.company || ""}`.toLowerCase(),
        });
      } catch (err) {
        const status = err?.response?.status;
        if (status === 409) this.addError = "Trùng ID tàu. Vui lòng chọn ID khác.";
        else this.addError = err?.response?.data?.error || "Thêm tàu thất bại.";
        throw err;
      } finally {
        this.addLoading = false;
      }
    },
  },
});
