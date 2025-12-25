import { defineStore } from 'pinia'
import axios from 'axios' // Bắt buộc phải import axios

export const useFleetStore = defineStore('fleet', {
    state: () => ({
        ships: [],       // Danh sách tàu
        loading: false,  // Trạng thái tải
        filterText: ''   // Ô tìm kiếm
    }),
    
    getters: {
        // 1. Tính toán thống kê KPI cho Dashboard
        stats(state) {
            const total = state.ships.length;
            const online = state.ships.filter(s => s.status === 'Online').length;
            const offline = state.ships.filter(s => s.status === 'Offline').length;
            const warning = state.ships.filter(s => s.status === 'Warning' || s.status === 'Blockage').length;
            
            // Công thức tính điểm sức khỏe hệ thống (Health Score)
            const health = total > 0 ? ((online + (warning * 0.5)) / total * 100).toFixed(1) : 0;

            return [
                { label: 'Total Fleet', value: total, color: 'primary', icon: 'fa-solid fa-ship' },
                { label: 'System Health', value: health + '%', color: 'success', icon: 'fa-solid fa-heart-pulse' },
                { label: 'Critical Alerts', value: offline, color: 'danger', icon: 'fa-solid fa-triangle-exclamation' },
                { label: 'Warnings', value: warning, color: 'warning', icon: 'fa-solid fa-circle-exclamation' }
            ];
        },
        
        // 2. Logic lọc danh sách tàu (Search)
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
        // --- CÁC HÀM QUẢN LÝ TÀU (SHIP) ---

        // Lấy danh sách tàu từ Backend Go
        async fetchFleet() {
            this.loading = true;
            try {
                const response = await axios.get('http://localhost:8080/api/ships');
                // Nếu backend trả về null (do chưa có tàu nào), gán mảng rỗng []
                this.ships = response.data || [];
            } catch (error) {
                console.error("Lỗi tải dữ liệu tàu:", error);
                // Không alert ở đây để tránh spam popup khi mới vào trang
            } finally {
                this.loading = false;
            }
        },

        // Thêm tàu mới
        async addShip(newShipData) {
            try {
                const response = await axios.post('http://localhost:8080/api/ships', newShipData);
                if (response.data) {
                    // Thêm vào đầu danh sách (unshift) để hiện ra ngay lập tức
                    this.ships.unshift(response.data);
                    return true; // Trả về true để báo component đóng Modal
                }
                return false;
            } catch (error) {
                console.error("Lỗi thêm tàu:", error);
                // Hiển thị lỗi chi tiết từ Backend trả về
                const msg = error.response?.data?.error || error.message || "Lỗi kết nối Server!";
                alert("Không thể thêm tàu: " + msg);
                return false;
            }
        },

        // --- CÁC HÀM QUẢN LÝ THỦY THỦ (CREW) ---

        // Lấy danh sách thủy thủ theo ID tàu
        async fetchCrew(shipId) {
            try {
                const response = await axios.get(`http://localhost:8080/api/ships/${shipId}/crew`);
                return response.data || [];
            } catch (error) {
                console.error("Lỗi tải danh sách thủy thủ:", error);
                return [];
            }
        },

        // Thêm thủy thủ mới
        async addCrew(crewData) {
            try {
                await axios.post('http://localhost:8080/api/crew', crewData);
                return true;
            } catch (error) {
                console.error("Lỗi thêm thủy thủ:", error);
                alert("Không thể thêm thành viên: " + (error.response?.data?.error || error.message));
                return false;
            }
        }
    }
})