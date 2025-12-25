<script setup>
import { ref, onMounted, reactive, computed } from 'vue';
import { useBandwidthStore } from '@/stores/bandwidth';

const store = useBandwidthStore();
const viewMode = ref('list'); // 'list' hoặc 'create'

// --- FORM DATA ---
const form = reactive({
    name: '',
    upload_val: 10, upload_unit: 'Mbps',
    download_val: 10, download_unit: 'Mbps',
    burst_limit: '', burst_threshold: '', burst_time: '',
    priority: 8, limit_at: '',
    status: 'Active'
});

// --- LOGIC ---
onMounted(() => store.fetchPlans());

const switchView = (mode) => viewMode.value = mode;

// Hàm tạo tên tự động (Optional)
const autoName = computed(() => {
    return `Plan ${form.upload_val}${form.upload_unit}/${form.download_val}${form.download_unit}`;
});

// Xử lý Lưu
const handleSave = async () => {
    // Quy đổi về Kbps để lưu DB (Chuẩn kỹ thuật)
    const up = form.upload_unit === 'Mbps' ? form.upload_val * 1024 : form.upload_val;
    const down = form.download_unit === 'Mbps' ? form.download_val * 1024 : form.download_val;

    const payload = {
        name: form.name || autoName.value,
        upload_speed: up,
        download_speed: down,
        burst_limit: form.burst_limit,
        burst_threshold: form.burst_threshold,
        burst_time: form.burst_time,
        priority: form.priority,
        status: form.status
    };

    if(await store.createPlan(payload)) {
        await store.fetchPlans();
        switchView('list');
        // Reset form
        Object.assign(form, { name: '', upload_val: 10, download_val: 10, burst_limit: '' }); 
    }
};

// --- PRESET ACTIONS (Cột bên phải giống hình) ---
const applyPreset = (type) => {
    if (type === '2x') {
        form.burst_limit = `${form.upload_val * 2}M/${form.download_val * 2}M`;
        form.burst_threshold = `${form.upload_val}M/${form.download_val}M`;
        form.burst_time = '16/16';
    } else if (type === '1M') {
        form.upload_val = 1; form.upload_unit = 'Mbps';
        form.download_val = 1; form.download_unit = 'Mbps';
    }
    // Thêm các logic preset khác tùy ý
};

// Format hiển thị
const formatSpeed = (kbps) => {
    if (kbps >= 1024) return (kbps / 1024).toFixed(0) + ' Mbps';
    return kbps + ' Kbps';
};
</script>

<template>
    <div class="container-fluid p-0">
        
        <!-- HEADER -->
        <h3 class="fw-bold text-dark mb-4">Bandwidth Plans</h3>

        <!-- VIEW 1: LIST -->
        <div v-if="viewMode === 'list'">
            <!-- FILTER -->
            <div class="card border-0 shadow-sm mb-4">
                <div class="card-body py-4 d-flex gap-3">
                    <input class="form-control w-25" placeholder="Name">
                    <select class="form-select w-25"><option>Active</option></select>
                    <button class="btn btn-primary px-4">Query</button>
                    <button class="btn btn-light border">Reset</button>
                </div>
            </div>

            <!-- TABLE -->
            <div class="card border-0 shadow-sm">
                <div class="card-header bg-white py-3 d-flex justify-content-between">
                    <h5 class="m-0 fw-bold">Bandwidth List</h5>
                    <button class="btn btn-primary fw-bold" @click="switchView('create')">
                        <i class="fa-solid fa-plus me-1"></i> Create
                    </button>
                </div>
                <div class="table-responsive">
                    <table class="table table-hover align-middle mb-0">
                        <thead class="bg-light text-secondary small text-uppercase">
                            <tr>
                                <th class="ps-4">ID</th><th>Name</th><th>Download</th><th>Upload</th><th>Burst</th><th>Status</th><th>Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="p in store.plans" :key="p.id">
                                <td class="ps-4 text-muted">#{{ p.id }}</td>
                                <td class="fw-bold">{{ p.name }}</td>
                                <td class="text-primary">{{ formatSpeed(p.download_speed) }}</td>
                                <td class="text-success">{{ formatSpeed(p.upload_speed) }}</td>
                                <td class="small text-muted">{{ p.burst_limit || '-' }}</td>
                                <td><span class="badge bg-success bg-opacity-10 text-success border border-success border-opacity-25">{{ p.status }}</span></td>
                                <td>
                                    <button class="btn btn-sm btn-link text-danger" @click="store.deletePlan(p.id)"><i class="fa-solid fa-trash"></i></button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <!-- VIEW 2: CREATE FORM (Giống hình 2) -->
        <div v-if="viewMode === 'create'" class="row g-4">
            <!-- Cột Trái: Form -->
            <div class="col-md-8">
                <div class="card border-0 shadow-sm p-4">
                    <h5 class="fw-bold mb-4">Create Plan</h5>
                    
                    <div class="mb-3">
                        <label class="form-label fw-bold">Name</label>
                        <input v-model="form.name" class="form-control" :placeholder="autoName">
                    </div>

                    <div class="row mb-3">
                        <div class="col-md-6">
                            <label class="form-label fw-bold text-danger">* Upload Speed</label>
                            <div class="input-group">
                                <input v-model="form.upload_val" type="number" class="form-control">
                                <select v-model="form.upload_unit" class="form-select bg-light" style="max-width: 100px;">
                                    <option>Mbps</option><option>Kbps</option>
                                </select>
                            </div>
                        </div>
                        <div class="col-md-6">
                            <label class="form-label fw-bold text-danger">* Download Speed</label>
                            <div class="input-group">
                                <input v-model="form.download_val" type="number" class="form-control">
                                <select v-model="form.download_unit" class="form-select bg-light" style="max-width: 100px;">
                                    <option>Mbps</option><option>Kbps</option>
                                </select>
                            </div>
                        </div>
                    </div>

                    <h6 class="fw-bold mt-4 mb-3">Optional (Burst)</h6>
                    <div class="mb-3"><label class="small">Burst Limit</label><input v-model="form.burst_limit" class="form-control"></div>
                    <div class="mb-3"><label class="small">Burst Threshold</label><input v-model="form.burst_threshold" class="form-control"></div>
                    <div class="mb-3"><label class="small">Burst Time</label><input v-model="form.burst_time" class="form-control"></div>
                    
                    <div class="mt-4">
                        <button class="btn btn-primary px-4 fw-bold me-2" @click="handleSave">Save</button>
                        <button class="btn btn-light border" @click="switchView('list')">Cancel</button>
                    </div>
                </div>
            </div>

            <!-- Cột Phải: Presets -->
            <div class="col-md-4">
                <div class="card border-0 shadow-sm p-4">
                    <h5 class="fw-bold mb-3">Burst Preset</h5>
                    <div class="row g-2">
                        <div class="col-6"><button class="btn btn-primary w-100" @click="applyPreset('2x')">2x Burst</button></div>
                        <div class="col-6"><button class="btn btn-outline-primary w-100" @click="applyPreset('1M')">upto 1M</button></div>
                        <div class="col-6"><button class="btn btn-outline-secondary w-100">3M to 6M</button></div>
                        <div class="col-6"><button class="btn btn-outline-secondary w-100">5M to 10M</button></div>
                    </div>
                    <div class="alert alert-info mt-4 small border-0 bg-info bg-opacity-10">
                        <i class="fa-solid fa-circle-info me-1"></i>
                        Presets help you quickly configure burst rates for better user experience during short browsing sessions.
                    </div>
                </div>
            </div>
        </div>

    </div>
</template>