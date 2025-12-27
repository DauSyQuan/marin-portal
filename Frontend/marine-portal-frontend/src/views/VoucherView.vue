<script setup>
import { ref, onMounted, computed } from 'vue';
import api from '@/services/api';

const vouchers = ref([]);
const loading = ref(false);
const searchQuery = ref('');

const fetchVouchers = async () => {
    loading.value = true;
    try {
        const res = await api.get('/api/vouchers');
        vouchers.value = Array.isArray(res.data) ? res.data : [];
    } catch (e) {
        console.error(e);
    } finally {
        loading.value = false;
    }
};

const filteredVouchers = computed(() => {
    if (!searchQuery.value) return vouchers.value;
    return vouchers.value.filter(v =>
        (v.code || '').toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        (v.status || '').toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        (v.assigned_to || '').toLowerCase().includes(searchQuery.value.toLowerCase())
    );
});

// Create voucher
const newVoucher = ref({
    code: '',
    data_limit_gb: 10,
    expires_at: '',
});

const createVoucher = async () => {
    try {
        await api.post('/api/vouchers', newVoucher.value);
        newVoucher.value = { code: '', data_limit_gb: 10, expires_at: '' };
        await fetchVouchers();
    } catch (e) {
        console.error(e);
        alert("Tạo voucher lỗi: " + (e?.response?.data?.error || e.message));
    }
};

// Assign voucher
const assignInfo = ref({
    voucher_id: null,
    assigned_to: '',
});

const assignVoucher = async () => {
    if (!assignInfo.value.voucher_id) return alert("Chọn voucher cần assign!");
    try {
        await api.put(`/api/vouchers/${assignInfo.value.voucher_id}/assign`, {
            assigned_to: assignInfo.value.assigned_to
        });
        assignInfo.value = { voucher_id: null, assigned_to: '' };
        await fetchVouchers();
    } catch (e) {
        console.error(e);
        alert("Assign voucher lỗi: " + (e?.response?.data?.error || e.message));
    }
};

onMounted(fetchVouchers);
</script>

<template>
    <div class="container-fluid p-0 fade-in-tab">
        <div class="d-flex justify-content-between align-items-end mb-4">
            <div>
                <small class="text-info fw-bold tracking-wide">ACCESS CONTROL</small>
                <h3 class="fw-bold text-white m-0">Vouchers</h3>
            </div>
            <div class="d-flex gap-2 align-items-center">
                <input
                    v-model="searchQuery"
                    type="text"
                    class="form-control form-control-sm bg-dark text-white border-secondary"
                    placeholder="Search voucher / status / user..."
                    style="width: 260px;"
                />
                <button class="btn btn-outline-light btn-sm fw-bold" @click="fetchVouchers" :disabled="loading">
                    <i class="fa-solid fa-rotate me-2"></i>Refresh
                </button>
            </div>
        </div>

        <!-- CREATE + ASSIGN -->
        <div class="row g-4 mb-4">
            <div class="col-lg-6">
                <div class="card glass-panel p-4 h-100">
                    <h6 class="fw-bold text-white mb-3">Create Voucher</h6>
                    <div class="mb-3">
                        <label class="small text-white-50">Code</label>
                        <input v-model="newVoucher.code" class="form-control bg-dark text-white border-secondary" placeholder="VOU-001" />
                    </div>
                    <div class="mb-3">
                        <label class="small text-white-50">Data Limit (GB)</label>
                        <input v-model.number="newVoucher.data_limit_gb" type="number" class="form-control bg-dark text-white border-secondary" />
                    </div>
                    <div class="mb-3">
                        <label class="small text-white-50">Expires At</label>
                        <input v-model="newVoucher.expires_at" type="date" class="form-control bg-dark text-white border-secondary" />
                    </div>
                    <button class="btn btn-primary fw-bold" @click="createVoucher">
                        <i class="fa-solid fa-plus me-2"></i>Create
                    </button>
                </div>
            </div>

            <div class="col-lg-6">
                <div class="card glass-panel p-4 h-100">
                    <h6 class="fw-bold text-white mb-3">Assign Voucher</h6>
                    <div class="mb-3">
                        <label class="small text-white-50">Voucher</label>
                        <select v-model="assignInfo.voucher_id" class="form-select bg-dark text-white border-secondary">
                            <option :value="null" disabled>Choose voucher...</option>
                            <option v-for="v in vouchers" :key="v.id" :value="v.id">
                                {{ v.code }} ({{ v.status }})
                            </option>
                        </select>
                    </div>
                    <div class="mb-3">
                        <label class="small text-white-50">Assign To</label>
                        <input v-model="assignInfo.assigned_to" class="form-control bg-dark text-white border-secondary" placeholder="username" />
                    </div>
                    <button class="btn btn-outline-info fw-bold" @click="assignVoucher">
                        <i class="fa-solid fa-user-check me-2"></i>Assign
                    </button>
                </div>
            </div>
        </div>

        <!-- TABLE -->
        <div class="card glass-panel p-0">
            <div class="table-responsive">
                <table class="table table-hover align-middle mb-0">
                    <thead>
                        <tr>
                            <th class="ps-4">Code</th>
                            <th>Limit</th>
                            <th>Status</th>
                            <th>Assigned To</th>
                            <th>Expires</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="loading">
                            <td colspan="5" class="text-center py-4 text-white-50">Loading...</td>
                        </tr>

                        <tr v-else-if="filteredVouchers.length === 0">
                            <td colspan="5" class="text-center py-4 text-white-50">No voucher found.</td>
                        </tr>

                        <tr v-else v-for="v in filteredVouchers" :key="v.id">
                            <td class="ps-4 fw-bold text-white">{{ v.code }}</td>
                            <td class="text-white-50">{{ v.data_limit_gb }} GB</td>
                            <td>
                                <span class="badge"
                                  :class="v.status === 'available'
                                    ? 'bg-success bg-opacity-25 text-success'
                                    : 'bg-warning bg-opacity-25 text-warning'">
                                    {{ v.status }}
                                </span>
                            </td>
                            <td class="text-white-50">{{ v.assigned_to || '-' }}</td>
                            <td class="text-white-50">{{ v.expires_at || '-' }}</td>
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
