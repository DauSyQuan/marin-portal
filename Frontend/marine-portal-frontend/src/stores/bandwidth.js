import { defineStore } from 'pinia'
import axios from 'axios'

export const useBandwidthStore = defineStore('bandwidth', {
    state: () => ({
        plans: [],
        loading: false
    }),
    actions: {
        async fetchPlans() {
            this.loading = true;
            try {
                const res = await axios.get('http://localhost:8080/api/bandwidth-plans');
                this.plans = res.data || [];
            } finally { this.loading = false; }
        },
        async createPlan(data) {
            try {
                // Convert đơn vị về Kbps trước khi gửi (nếu cần) hoặc gửi thô
                await axios.post('http://localhost:8080/api/bandwidth-plans', data);
                return true;
            } catch (e) {
                alert(e.message);
                return false;
            }
        },
        async deletePlan(id) {
            if(!confirm("Xóa gói này?")) return;
            await axios.delete(`http://localhost:8080/api/bandwidth-plans/${id}`);
            this.fetchPlans();
        }
    }
})