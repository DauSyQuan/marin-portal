<script setup>
import { ref, onMounted, computed } from 'vue';
import api from '@/services/api';

const data = ref([]);
const loading = ref(false);
const searchQuery = ref('');

const fetchMonthlyUsage = async () => {
    loading.value = true;
    try {
        const res = await api.get('/api/usage-report');
        data.value = Array.isArray(res.data) ? res.data : [];
    } catch (e) {
        console.error(e);
    } finally {
        loading.value = false;
    }
};

const filteredData = computed(() => {
    if (!searchQuery.value) return data.value;
    return data.value.filter(x =>
        (x.ship_name || '').toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        (x.imo_number || '').includes(searchQuery.value)
    );
});

onMounted(fetchMonthlyUsage);
</script>

<template>
    <div class="container-fluid p-0 fade-in-tab">
        <div class="d-flex justify-content-between align-items-end mb-4">
            <div>
                <small class="text-info fw-bold tracking-wide">BILLING MODULE</small>
                <h3 class="fw-bold text-white m-0">Monthly Usage</h3>
            </div>
            <div class="d-flex gap-2 align-items-center">
                <input
                    v-model="searchQuery"
                    type="text"
                    class="form-control form-control-sm bg-dark text-white border-secondary"
                    placeholder="Search ship / IMO..."
                    style="width: 230px;"
                />
                <button class="btn btn-outline-light btn-sm fw-bold" @click="fetchMonthlyUsage" :disabled="loading">
                    <i class="fa-solid fa-rotate me-2"></i>Refresh
                </button>
            </div>
        </div>

        <div class="card glass-panel p-0">
            <div class="table-responsive">
                <table class="table table-hover align-middle mb-0">
                    <thead>
                        <tr>
                            <th class="ps-4">Vessel</th>
                            <th>IMO</th>
                            <th>Month</th>
                            <th>Total Usage</th>
                            <th>Cost</th>
                            <th class="text-end pe-4">Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="loading">
                            <td colspan="6" class="text-center py-4 text-white-50">Loading...</td>
                        </tr>

                        <tr v-else-if="filteredData.length === 0">
                            <td colspan="6" class="text-center py-4 text-white-50">No report found.</td>
                        </tr>

                        <tr v-else v-for="row in filteredData" :key="row.id">
                            <td class="ps-4 fw-bold text-white">{{ row.ship_name }}</td>
                            <td class="text-white-50 font-monospace">{{ row.imo_number }}</td>
                            <td class="text-white-50">{{ row.month }}</td>
                            <td class="fw-bold text-white">{{ row.total_usage_gb }} GB</td>
                            <td class="font-monospace text-warning">${{ row.cost }}</td>
                            <td class="text-end pe-4">
                                <span class="badge bg-success bg-opacity-25 text-success">OK</span>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>

<style scoped>
.tracking-wide { letter-spacing: 1px; }
.fade-in-tab { animation: fadeIn 0.4s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>
