<script setup>
import { computed, onMounted, onUnmounted, ref, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useFleetStore } from '@/stores/fleet';
import BaseMap from '@/components/BaseMap.vue';
import VsatChart from '@/components/VsatChart.vue';
import { Modal } from 'bootstrap';

const route = useRoute();
const router = useRouter();
const store = useFleetStore();

// State
const liveSnr = ref(0);
const mapMode = ref('internal');
const activeTab = ref('crew'); 
const crewList = ref([]);
let simulationInterval = null;
let crewModalInstance = null;
const crewModalRef = ref(null);

const ship = computed(() => store.ships.find(s => s.id === route.params.id));

// State Bộ lọc
const filter = reactive({ username: '', name: '', tab: 'Active' });

// --- STATE FORM (Dùng chung cho Thêm & Sửa) ---
const isEditing = ref(false); // Cờ đánh dấu đang sửa hay thêm
const editingId = ref(null);  // Lưu ID người đang sửa

const crewForm = reactive({
    full_name: '', username: '', rank: 'Captain', 
    nationality: 'Vietnam', data_plan: 'Basic (1GB)', status: 'Active'
});

onMounted(async () => {
    if (store.ships.length === 0) await store.fetchFleet();
    if (ship.value) {
        liveSnr.value = ship.value.snr;
        loadCrew();
    }
    if(crewModalRef.value) crewModalInstance = new Modal(crewModalRef.value);

    simulationInterval = setInterval(() => {
        if (ship.value) {
            const base = ship.value.status === 'Offline' ? 0 : ship.value.snr;
            liveSnr.value = parseFloat((base + (Math.random() - 0.5)).toFixed(1));
        }
    }, 1000);
});

onUnmounted(() => { if (simulationInterval) clearInterval(simulationInterval); });

// --- LOGIC CRUD ---

const loadCrew = async () => {
    if(ship.value) crewList.value = await store.fetchCrew(ship.value.id);
};

// 1. Mở Modal để THÊM MỚI
const openCreateModal = () => {
    isEditing.value = false;
    editingId.value = null;
    // Reset form
    Object.assign(crewForm, { 
        full_name: '', username: '', rank: 'Captain', 
        nationality: 'Vietnam', data_plan: 'Basic (1GB)', status: 'Active' 
    });
    crewModalInstance.show();
};

// 2. Mở Modal để SỬA (Đổ dữ liệu cũ vào)
const openEditModal = (crew) => {
    isEditing.value = true;
    editingId.value = crew.id;
    // Copy dữ liệu từ dòng đã chọn vào form
    Object.assign(crewForm, { 
        full_name: crew.full_name, 
        username: crew.username, 
        rank: crew.rank, 
        nationality: crew.nationality, 
        data_plan: crew.data_plan, 
        status: crew.status 
    });
    crewModalInstance.show();
};

// 3. Xử lý Lưu (Phân loại Thêm hoặc Sửa)
const handleSave = async () => {
    if(!crewForm.username || !crewForm.full_name) return alert("Vui lòng điền đủ thông tin!");
    
    let success = false;
    if (isEditing.value) {
        // Gọi hàm Sửa
        success = await store.updateCrew(editingId.value, crewForm);
    } else {
        // Gọi hàm Thêm
        success = await store.addCrew({ ...crewForm, ship_id: ship.value.id });
    }

    if(success) {
        crewModalInstance.hide();
        loadCrew(); // Tải lại danh sách
    }
};

// 4. Xử lý Xóa
const handleDelete = async (id) => {
    const success = await store.deleteCrew(id);
    if(success) loadCrew();
};

// Filter Logic
const filteredCrew = computed(() => {
    return crewList.value.filter(c => {
        if (filter.tab !== 'All' && c.status !== filter.tab) return false;
        if (filter.username && !c.username?.toLowerCase().includes(filter.username.toLowerCase())) return false;
        if (filter.name && !c.full_name?.toLowerCase().includes(filter.name.toLowerCase())) return false;
        return true;
    });
});
const resetFilter = () => { filter.username = ''; filter.name = ''; filter.tab = 'Active'; };

// Utils
const getMarineTrafficUrl = (ship) => `https://www.marinetraffic.com/en/ais/embed/map/imo:${ship.id.replace('IMO', '') || '9330264'}/latitude:${ship.lat}/longitude:${ship.lon}/zoom:9`;
const downloadPdf = () => window.open(`http://localhost:8080/api/report/${route.params.id}`, '_blank');
const formatCoord = (val, type) => `${val?.toFixed(4)}° ${type === 'lat' ? (val>0?'N':'S') : (val>0?'E':'W')}`;
const formatDate = (d) => d ? new Date(d).toLocaleDateString('vi-VN') : 'Just now';
</script>

<template>
    <div class="container-fluid p-0">
        <!-- HEADER -->
        <div class="d-flex justify-content-between align-items-center mb-4">
            <div class="d-flex align-items-center">
                <button class="btn btn-outline-secondary border-0 me-2" @click="router.push('/')"><i class="fa-solid fa-arrow-left"></i></button>
                <div>
                    <div class="text-uppercase text-secondary" style="font-size: 0.65rem; letter-spacing: 1px;">Vessel Monitor</div>
                    <h3 class="fw-bold m-0 d-flex align-items-center">
                        {{ ship?.name }}
                        <span v-if="ship" class="badge ms-3 fs-6 rounded-pill" :class="ship.status === 'Online' ? 'bg-success' : 'bg-danger'">{{ ship.status }}</span>
                    </h3>
                </div>
            </div>
            <div class="d-flex gap-2">
                <button class="btn btn-primary shadow fw-bold" @click="downloadPdf"><i class="fa-solid fa-file-pdf me-2"></i> Report</button>
            </div>
        </div>

        <div v-if="!ship" class="alert alert-warning">Loading vessel data...</div>

        <div v-if="ship">
            <ul class="nav nav-tabs mb-4">
                <li class="nav-item"><a class="nav-link cursor-pointer" :class="{ active: activeTab === 'overview' }" @click="activeTab = 'overview'"><i class="fa-solid fa-chart-line me-2"></i>Overview</a></li>
                <li class="nav-item"><a class="nav-link cursor-pointer" :class="{ active: activeTab === 'crew' }" @click="activeTab = 'crew'"><i class="fa-solid fa-wifi me-2"></i>ICT User & Internet</a></li>
            </ul>

            <!-- TAB 1: OVERVIEW -->
            <div v-if="activeTab === 'overview'" class="row g-4">
                <div class="col-12">
                    <div class="card border-0 shadow-sm bg-dark text-white">
                        <div class="card-body d-flex justify-content-between align-items-center py-3">
                            <div class="d-flex gap-5">
                                <div><small class="text-white-50 d-block">IMO</small><span class="font-monospace fw-bold">{{ ship.id }}</span></div>
                                <div><small class="text-white-50 d-block">Position</small><span class="font-monospace text-info">{{ formatCoord(ship.lat, 'lat') }} / {{ formatCoord(ship.lon, 'lon') }}</span></div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="col-lg-4">
                    <div class="card shadow-sm border-0 h-100">
                        <div class="card-header bg-white fw-bold">Connection Health</div>
                        <div class="card-body">
                            <div class="row g-3 mb-4">
                                <div class="col-6"><div class="p-2 bg-light border rounded"><small>IP</small><div class="fw-bold">{{ ship.ip }}</div></div></div>
                                <div class="col-6"><div class="p-2 bg-light border rounded"><small>Beam</small><div class="fw-bold">{{ ship.satellite }}</div></div></div>
                            </div>
                            <div class="d-flex justify-content-between align-items-end mb-2">
                                <div class="display-4 fw-bold text-primary">{{ liveSnr }} <span class="fs-6 text-muted">dB</span></div>
                            </div>
                            <div style="height: 200px; margin: 0 -10px;"><VsatChart :snr="liveSnr" /></div>
                        </div>
                    </div>
                </div>
                <div class="col-lg-8">
                    <div class="card shadow-sm border-0 h-100 position-relative">
                        <div class="position-absolute top-0 end-0 m-3 z-3 bg-white p-1 rounded shadow-sm border">
                            <div class="btn-group btn-group-sm">
                                <button class="btn" :class="mapMode==='internal'?'btn-primary':'btn-light'" @click="mapMode='internal'">GPS</button>
                                <button class="btn" :class="mapMode==='marinetraffic'?'btn-primary':'btn-light'" @click="mapMode='marinetraffic'">Traffic</button>
                            </div>
                        </div>
                        <div class="h-100 w-100" style="min-height: 500px;">
                            <BaseMap v-if="mapMode==='internal'" :lat="ship.lat" :lon="ship.lon" />
                            <iframe v-else width="100%" height="100%" :src="getMarineTrafficUrl(ship)" style="border:0"></iframe>
                        </div>
                    </div>
                </div>
            </div>

            <!-- TAB 2: ICT USER & INTERNET -->
            <div v-if="activeTab === 'crew'">
                <div class="card border-0 shadow-sm mb-4">
                    <div class="card-body py-4">
                        <div class="row g-3 align-items-end">
                            <div class="col-md-4"><label class="form-label small fw-bold text-secondary">Username:</label><input type="text" class="form-control" v-model="filter.username" placeholder="Please enter"></div>
                            <div class="col-md-4"><label class="form-label small fw-bold text-secondary">Name:</label><input type="text" class="form-control" v-model="filter.name" placeholder="Please enter"></div>
                            <div class="col-md-4 text-end"><button class="btn btn-light border me-2" @click="resetFilter">Reset</button><button class="btn btn-primary px-4" @click="loadCrew"><i class="fa-solid fa-magnifying-glass me-2"></i> Query</button></div>
                        </div>
                    </div>
                </div>

                <div class="card border-0 shadow-sm">
                    <div class="card-body p-0">
                        <div class="d-flex justify-content-between align-items-center p-3 border-bottom">
                            <ul class="nav nav-pills custom-pills">
                                <li class="nav-item"><a class="nav-link cursor-pointer" :class="{ active: filter.tab === 'Active' }" @click="filter.tab = 'Active'">Active</a></li>
                                <li class="nav-item"><a class="nav-link cursor-pointer" :class="{ active: filter.tab === 'Inactive' }" @click="filter.tab = 'Inactive'">Inactive</a></li>
                                <li class="nav-item"><a class="nav-link cursor-pointer" :class="{ active: filter.tab === 'All' }" @click="filter.tab = 'All'">All</a></li>
                            </ul>
                            <div class="d-flex gap-2">
                                <button class="btn btn-primary" @click="openCreateModal"><i class="fa-solid fa-plus me-2"></i>Create</button>
                                <button class="btn btn-light border" @click="loadCrew"><i class="fa-solid fa-rotate"></i></button>
                            </div>
                        </div>

                        <div class="table-responsive">
                            <table class="table table-hover align-middle mb-0">
                                <thead class="bg-light text-secondary small text-uppercase">
                                    <tr>
                                        <th class="ps-4">ID</th><th>Username</th><th>Name / Rank</th><th>Active Data Plan</th><th>Usage</th><th>Status</th><th>Created At</th><th class="text-end pe-4">Action</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr v-for="c in filteredCrew" :key="c.id">
                                        <td class="ps-4 text-muted">#{{ c.id }}</td>
                                        <td class="fw-bold text-primary">{{ c.username || 'N/A' }}</td>
                                        <td><div class="fw-bold">{{ c.full_name }}</div><small class="text-muted">{{ c.rank }}</small></td>
                                        <td><span class="badge bg-info bg-opacity-10 text-info border border-info border-opacity-25">{{ c.data_plan }}</span></td>
                                        <td><div class="d-flex align-items-center"><i class="fa-solid fa-chart-pie me-2 text-secondary"></i>{{ c.data_usage }} MB</div></td>
                                        <td><span class="badge rounded-pill" :class="c.status === 'Active' ? 'bg-success' : 'bg-secondary'">{{ c.status }}</span></td>
                                        <td class="text-secondary small">{{ formatDate(c.created_at) }}</td>
                                        <td class="text-end pe-4">
                                            <!-- NÚT SỬA -->
                                            <button class="btn btn-sm btn-link text-secondary" @click="openEditModal(c)"><i class="fa-solid fa-pen-to-square"></i></button>
                                            <!-- NÚT XÓA -->
                                            <button class="btn btn-sm btn-link text-danger" @click="handleDelete(c.id)"><i class="fa-solid fa-trash"></i></button>
                                        </td>
                                    </tr>
                                    <tr v-if="filteredCrew.length === 0"><td colspan="8" class="text-center py-5 text-muted">No data found.</td></tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>

            <!-- MODAL (DÙNG CHUNG CHO CREATE & EDIT) -->
            <div class="modal fade" id="crewModal" tabindex="-1" ref="crewModalRef">
                <div class="modal-dialog">
                    <div class="modal-content">
                        <div class="modal-header bg-primary text-white">
                            <h5 class="modal-title fw-bold">{{ isEditing ? 'Edit User' : 'Create User' }}</h5>
                            <button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
                        </div>
                        <div class="modal-body">
                            <div class="row g-3">
                                <div class="col-6"><label class="small fw-bold">Username</label><input v-model="crewForm.username" class="form-control" placeholder="e.g. user01"></div>
                                <div class="col-6"><label class="small fw-bold">Full Name</label><input v-model="crewForm.full_name" class="form-control" placeholder="e.g. Nguyen Van A"></div>
                                <div class="col-6">
                                    <label class="small fw-bold">Rank</label>
                                    <select v-model="crewForm.rank" class="form-select">
                                        <option>Captain</option><option>Chief Officer</option><option>Engineer</option><option>Crew</option>
                                    </select>
                                </div>
                                <div class="col-6">
                                    <label class="small fw-bold">Data Plan</label>
                                    <select v-model="crewForm.data_plan" class="form-select">
                                        <option>Basic (1GB)</option><option>Standard (5GB)</option><option>Premium (Unlimited)</option><option>VIP</option>
                                    </select>
                                </div>
                                <div class="col-12">
                                    <div class="form-check form-switch p-3 bg-light rounded border">
                                        <input class="form-check-input" type="checkbox" id="statusSwitch" v-model="crewForm.status" true-value="Active" false-value="Inactive">
                                        <label class="form-check-label fw-bold" for="statusSwitch">Enable Internet Access</label>
                                    </div>
                                </div>
                            </div>
                            <div class="mt-4 text-end">
                                <button class="btn btn-light me-2" data-bs-dismiss="modal">Cancel</button>
                                <button class="btn btn-primary fw-bold px-4" @click="handleSave">{{ isEditing ? 'Update User' : 'Create User' }}</button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

        </div>
    </div>
</template>

<style scoped>
.cursor-pointer { cursor: pointer; }
.nav-link { color: #64748b; font-weight: 600; }
.nav-link.active { color: #2563eb; border-bottom: 2px solid #2563eb; background: transparent; }
.custom-pills .nav-link { color: #64748b; font-weight: 500; padding: 6px 16px; border-radius: 4px; border: none; }
.custom-pills .nav-link:hover { background-color: #f1f5f9; }
.custom-pills .nav-link.active { background-color: transparent; color: #2563eb; font-weight: 700; border-bottom: 2px solid #2563eb; border-radius: 0; }
</style>