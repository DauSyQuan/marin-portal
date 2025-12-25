import { defineStore } from 'pinia'
import axios from 'axios'

export const useFleetStore = defineStore('fleet', {
    state: () => ({
        ships: [],
        loading: false,
        filterText: ''
    }),
    
    getters: {
        // Tính toán thống kê từ dữ liệu thật
        stats(state) {
            const total = state.ships.length;
            const online = state.ships.filter(s => s.status === 'Online').length;
            const offline = state.ships.filter(s => s.status === 'Offline').length;
            const warning = state.ships.filter(s => s.status === 'Warning' || s.status === 'Blockage').length;
            
            // Công thức tính điểm sức khỏe hệ thống
            const health = total > 0 ? ((online + (warning * 0.5)) / total * 100).toFixed(1) : 0;

            return [
                { label: 'Total Fleet', value: total, color: 'primary', icon: 'fa-solid fa-ship' },
                { label: 'System Health', value: health + '%', color: 'success', icon: 'fa-solid fa-heart-pulse' },
                { label: 'Critical Alerts', value: offline, color: 'danger', icon: 'fa-solid fa-triangle-exclamation' },
                { label: 'Warnings', value: warning, color: 'warning', icon: 'fa-solid fa-circle-exclamation' }
            ];
        },
        
        // Lọc danh sách tàu
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
        async fetchFleet() {
            this.loading = true;
            try {
                // GỌI API TỪ SERVER GO (Port 8080)
                const response = await axios.get('http://localhost:8080/api/ships');
                
                // Gán dữ liệu thật vào biến ships
                this.ships = response.data;
                console.log("Đã tải dữ liệu từ Backend:", this.ships.length, "tàu");
                
            } catch (error) {
                console.error("Lỗi kết nối Backend:", error);
                alert("Không thể kết nối tới Server Go! Hãy kiểm tra xem Backend chạy chưa.");
            } finally {
                this.loading = false;
            }
        }
    }
})