<script setup>
import { ref, onMounted, reactive } from 'vue';
import axios from 'axios';

// State
const reports = ref([]);
const loading = ref(false);

// Filter mặc định: Từ đầu tháng đến nay
const today = new Date();
const firstDay = new Date(today.getFullYear(), today.getMonth(), 1);

const filter = reactive({
    username: '',
    startDate: firstDay.toISOString().split('T')[0],
    endDate: today.toISOString().split('T')[0]
});

// Hàm lấy dữ liệu
const fetchReport = async () => {
    loading.value = true;
    try {
        const res = await axios.get('http://localhost:8080/api/usage-report', {
            params: {
                username: filter.username,
                start_date: filter.startDate,
                end_date: filter.endDate
            }
        });
        reports.value = res.data;
    } catch (e) {
        console.error(e);
    } finally {
        loading.value = false;
    }
};

const resetFilter = () => {
    filter.username = '';
    filter.startDate = firstDay.toISOString().split('T')[0];
    filter.endDate = today.toISOString().split('T')[0];
    fetchReport();
};

// --- UTILS FORMAT ---
const formatBytes = (mb) => {
    if (mb === 0) return '0 MB';
    if (mb < 1) return (mb * 1024).toFixed(2) + ' KB';
    if (mb > 1024) return (mb / 1024).toFixed(2) + ' GB';
    return mb.toFixed(2) + ' MB';
};

const formatTime = (seconds) => {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    // Format: 120:54:24
    return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
};

onMounted(() => {
    fetchReport();
});
</script>

<template>
    <div class="container-fluid p-0">
        <h3 class="fw-bold text-dark mb-4">User Monthly Usage</h3>

        <!-- 1. FILTER BAR -->
        <div class="card border-0 shadow-sm mb-4">
            <div class="card-body py-4">
                <div class="row g-3 align-items-center">
                    <div class="col-md-4">
                        <div class="d-flex align-items-center">
                            <label class="me-3 fw-bold text-secondary" style="min-width: 80px;">Username:</label>
                            <input type="text" class="form-control" v-model="filter.username" placeholder="Please enter">
                        </div>
                    </div>
                    <div class="col-md-5">
                        <div class="d-flex align-items-center">
                            <label class="me-3 fw-bold text-secondary" style="min-width: 90px;">Select Range:</label>
                            <div class="input-group">
                                <input type="date" class="form-control" v-model="filter.startDate">
                                <span class="input-group-text bg-white border-start-0 border-end-0">→</span>
                                <input type="date" class="form-control" v-model="filter.endDate">
                            </div>
                        </div>
                    </div>
                    <div class="col-md-3 text-end">
                        <button class="btn btn-primary px-4 me-2" @click="fetchReport">
                            <i class="fa-solid fa-magnifying-glass me-2"></i>Query
                        </button>
                        <button class="btn btn-light border" @click="resetFilter">Reset</button>
                    </div>
                </div>
            </div>
        </div>

        <!-- 2. DATA TABLE -->
        <div class="card border-0 shadow-sm">
            <div class="card-header bg-white py-3">
                <h6 class="m-0 fw-bold">Data Usage Report</h6>
            </div>
            <div class="table-responsive">
                <table class="table table-hover align-middle mb-0">
                    <thead class="bg-light text-secondary small text-uppercase">
                        <tr>
                            <th class="ps-4 py-3">Username</th>
                            <th>Download</th>
                            <th>Upload</th>
                            <th>Total Data</th>
                            <th>Time Used</th>
                            <th class="text-end pe-4">Action</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in reports" :key="item.id">
                            <td class="ps-4 fw-bold text-dark">{{ item.username || 'Unknown' }}</td>
                            <td>{{ formatBytes(item.download) }}</td>
                            <td>{{ formatBytes(item.upload) }}</td>
                            <td class="fw-bold text-primary">{{ formatBytes(item.total_data) }}</td>
                            <td class="font-monospace text-secondary">{{ formatTime(item.time_used) }}</td>
                            <td class="text-end pe-4">
                                <button class="btn btn-sm btn-light border rounded-circle" title="Detail">
                                    <i class="fa-regular fa-clock"></i>
                                </button>
                            </td>
                        </tr>
                        <tr v-if="reports.length === 0">
                            <td colspan="6" class="text-center py-5 text-muted">No usage data found for this period.</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <!-- Pagination Fake giống hình -->
            <div class="card-footer bg-white py-3 d-flex justify-content-end align-items-center">
                <small class="text-secondary me-3">1-{{ reports.length }} of {{ reports.length }} items</small>
                <nav>
                    <ul class="pagination pagination-sm mb-0">
                        <li class="page-item disabled"><a class="page-link" href="#">&lt;</a></li>
                        <li class="page-item active"><a class="page-link" href="#">1</a></li>
                        <li class="page-item disabled"><a class="page-link" href="#">&gt;</a></li>
                    </ul>
                </nav>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Font số cho đồng hồ */
.font-monospace { font-family: 'Courier New', Courier, monospace; }
</style>