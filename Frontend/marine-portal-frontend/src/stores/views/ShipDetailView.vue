<script setup>
import { computed, onMounted, onUnmounted, ref, reactive, nextTick } from 'vue'; // Thêm nextTick
import { useRoute, useRouter } from 'vue-router';
import { useFleetStore } from '@/stores/fleet';
import BaseMap from '@/components/BaseMap.vue';
import VsatChart from '@/components/VsatChart.vue';
import { Modal } from 'bootstrap';
import axios from 'axios';

const route = useRoute();
const router = useRouter();
const store = useFleetStore();

// --- STATE CHUNG ---
const activeTab = ref('crew'); // Mặc định vào tab Crew để test luôn
const mapMode = ref('internal');
const liveSnr = ref(0);
let simulationInterval = null;

const ship = computed(() => store.ships.find(s => s.id === route.params.id));

// --- STATE: CREW / ICT USER ---
const crewList = ref([]);
const isEditing = ref(false);
const editingId = ref(null);

// QUAN TRỌNG: Khai báo ref cho Modal
const crewModalRef = ref(null);
let crewModalInstance = null;

const crewForm = reactive({
    full_name: '', username: '', rank: 'Captain', 
    nationality: 'Vietnam', data_plan: 'Basic (1GB)', status: 'Active'
});

const filter = reactive({ username: '', name: '', tab: 'Active' });

const filteredCrew = computed(() => {
    return crewList.value.filter(c => {
        if (filter.tab !== 'All' && c.status !== filter.tab) return false;
        if (filter.username && !c.username?.toLowerCase().includes(filter.username.toLowerCase())) return false;
        if (filter.name && !c.full_name?.toLowerCase().includes(filter.name.toLowerCase())) return false;
        return true;
    });
});

// --- STATE: ROUTER ---
const routerStats = ref({ cpu_load: 0, free_memory: 0, tx_rate: 0, rx_rate: 0, uptime: '', board_name: '' });
const routerConnected = ref(false);
const isSimulation = ref(false);
const routerIP = ref('');
let routerInterval = null;

// --- LIFECYCLE ---
onMounted(async () => {
    if (store.ships.length === 0) await store.fetchFleet();
    
    if (ship.value) {
        liveSnr.value = ship.value.snr;
        loadCrew();
        routerIP.value = ship.value.router_ip;
    }

    // --- FIX LỖI MODAL Ở ĐÂY ---
    // Dùng nextTick để đảm bảo HTML đã vẽ xong thì mới gắn Modal
    await nextTick();
    if(crewModalRef.value) {
        crewModalInstance = new Modal(crewModalRef.value);
    }

    simulationInterval = setInterval(() => {
        if (ship.value) {
            const base = ship.value.status === 'Offline' ? 0 : ship.value.snr;
            liveSnr.value = parseFloat((base + (Math.random() - 0.5)).toFixed(1));
        }
    }, 1000);

    fetchRouterStats();
    routerInterval = setInterval(fetchRouterStats, 3000);
});

onUnmounted(() => {
    if (simulationInterval) clearInterval(simulationInterval);
    if (routerInterval) clearInterval(routerInterval);
    // Ẩn modal nếu đang mở để tránh kẹt màn hình đen
    if (crewModalInstance) crewModalInstance.hide();
});

// --- LOGIC CREW ---
const loadCrew = async () => {
    if(ship.value) crewList.value = await store.fetchCrew(ship.value.id);
};

const openCreateModal = () => {
    // Kiểm tra lại lần nữa nếu modal chưa khởi tạo (phòng hờ)
    if (!crewModalInstance && crewModalRef.value) {
        crewModalInstance = new Modal(crewModalRef.value);
    }
    
    isEditing.value = false; 
    editingId.value = null;
    
    // Reset form
    Object.assign(crewForm, { 
        full_name: '', username: '', 
        rank: 'Captain', nationality: 'Vietnam', 
        data_plan: 'Basic (1GB)', status: 'Active' 
    });
    
    crewModalInstance?.show();
};

const openEditModal = (c) => {
    if (!crewModalInstance && crewModalRef.value) {
        crewModalInstance = new Modal(crewModalRef.value);
    }
    isEditing.value = true; 
    editingId.value = c.id;
    Object.assign(crewForm, c);
    crewModalInstance?.show();
};

const handleSaveCrew = async () => {
    if(!crewForm.username || !crewForm.full_name) return alert("Vui lòng điền đủ thông tin!");
    
    let success = false;
    if (isEditing.value) {
        success = await store.updateCrew(editingId.value, crewForm);
    } else {
        success = await store.addCrew({ ...crewForm, ship_id: ship.value.id });
    }

    if(success) {
        crewModalInstance.hide();
        loadCrew();
    }
};

const handleDeleteCrew = async (id) => {
    if(await store.deleteCrew(id)) loadCrew();
};

const resetFilter = () => { filter.username = ''; filter.name = ''; filter.tab = 'Active'; };

// --- LOGIC ROUTER ---
const fetchRouterStats = async () => {
    if (!ship.value) return;
    try {
        const res = await axios.get(`http://localhost:8080/api/ships/${route.params.id}/router/stats`);
        routerStats.value = res.data;
        if (res.data.connected) {
            routerConnected.value = true; isSimulation.value = false; routerIP.value = ship.value.router_ip;
        } else {
            routerConnected.value = false; isSimulation.value = true; routerIP.value = "Simulation Mode";
        }
    } catch(e) {
        routerConnected.value = false; isSimulation.value = false;
        routerStats.value = { cpu_load: 0, free_memory: 0, tx_rate: 0, rx_rate: 0 };
    }
};
const syncRouter = async () => { if(confirm("Đồng bộ danh sách User xuống Router?")) { try { const res = await axios.post(`http://localhost:8080/api/ships/${route.params.id}/router/sync`); alert(res.data.message); } catch(e) { alert("Lỗi: " + e.message); } } };
const rebootRouter = async () => { if(confirm("CẢNH BÁO: Khởi động lại Router?")) { try { await axios.post(`http://localhost:8080/api/ships/${route.params.id}/router/reboot`); alert(routerConnected.value ? "Đã gửi lệnh Reboot." : "Đã gửi lệnh Reboot (Giả lập)."); } catch(e) { alert("Lỗi: " + e.message); } } };
const blockInternet = async () => { if(confirm("KHẨN CẤP: Chặn toàn bộ Internet?")) alert("Đã kích hoạt Firewall Drop Rule!"); };

// --- UTILS ---
const getMarineTrafficUrl = (s) => `https://www.marinetraffic.com/en/ais/embed/map/imo:${s.id.replace('IMO', '') || '9330264'}/latitude:${s.lat}/longitude:${s.lon}/zoom:9`;
const downloadPdf = () => window.open(`http://localhost:8080/api/report/${route.params.id}`, '_blank');
const formatCoord = (val, type) => `${val?.toFixed(4)}° ${type === 'lat' ? (val>0?'N':'S') : (val>0?'E':'W')}`;
const formatDate = (d) => d ? new Date(d).toLocaleDateString('vi-VN') : 'Just now';
</script>

<template>
    <div class="container-fluid p-0">
        <!-- HEADER -->
        <div class="d-flex justify-content-between align-items-center mb-4 pb-3 border-bottom border-white border-opacity-10">
            <div class="d-flex align-items-center">
                <button class="btn btn-outline-light border-0 me-3" @click="router.push('/')"><i class="fa-solid fa-arrow-left"></i></button>
                <div>
                    <div class="text-uppercase text-info small fw-bold tracking-wide">Vessel Monitor</div>
                    <h3 class="fw-bold m-0 text-white d-flex align-items-center">
                        {{ ship?.name }}
                        <span v-if="ship" class="badge ms-3 fs-6 border border-opacity-25" :class="ship.status === 'Online' ? 'bg-success bg-opacity-25 text-success border-success' : 'bg-danger bg-opacity-25 text-danger border-danger'">{{ ship.status }}</span>
                    </h3>
                </div>
            </div>
            <div class="d-flex gap-2">
                <button class="btn btn-primary fw-bold px-4 shadow-sm" @click="downloadPdf"><i class="fa-solid fa-file-pdf me-2"></i> Voyage Report</button>
            </div>
        </div>

        <div v-if="!ship" class="text-center py-5"><div class="spinner-border text-primary"></div></div>

        <div v-if="ship">
            <!-- TABS -->
            <ul class="nav nav-tabs mb-4 border-bottom border-white border-opacity-10">
                <li class="nav-item"><a class="nav-link cursor-pointer" :class="{ active: activeTab === 'overview' }" @click="activeTab = 'overview'"><i class="fa-solid fa-chart-line me-2"></i>Overview</a></li>
                <li class="nav-item"><a class="nav-link cursor-pointer" :class="{ active: activeTab === 'crew' }" @click="activeTab = 'crew'"><i class="fa-solid fa-users me-2"></i>ICT User & Internet</a></li>
                <li class="nav-item"><a class="nav-link cursor-pointer" :class="{ active: activeTab === 'router' }" @click="activeTab = 'router'"><i class="fa-solid fa-server me-2"></i>Router & Gateway</a></li>
            </ul>

            <!-- TAB 1: OVERVIEW -->
            <div v-if="activeTab === 'overview'" class="row g-4 fade-in-tab">
                <!-- (Nội dung Overview giữ nguyên) -->
                <div class="col-12">
                    <div class="card glass-panel p-4">
                        <div class="d-flex justify-content-between align-items-center">
                            <div class="d-flex gap-5">
                                <div><small class="text-secondary d-block uppercase-label">IMO Number</small><span class="font-monospace fw-bold text-white">{{ ship.id }}</span></div>
                                <div><small class="text-secondary d-block uppercase-label">Position</small><span class="font-monospace text-info">{{ formatCoord(ship.lat, 'lat') }} / {{ formatCoord(ship.lon, 'lon') }}</span></div>
                                <div><small class="text-secondary d-block uppercase-label">Last Contact</small><span class="text-white"><i class="fa-solid fa-satellite me-2 text-success"></i>Just now</span></div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="col-lg-4">
                    <div class="card glass-panel h-100 p-4">
                        <h6 class="fw-bold text-white mb-4 border-bottom border-white border-opacity-10 pb-2">Connection Health</h6>
                        <div class="row g-3 mb-4">
                            <div class="col-6"><div class="p-3 rounded border border-white border-opacity-10 bg-dark bg-opacity-25"><small class="text-secondary d-block">IP Address</small><strong class="font-monospace text-white">{{ ship.ip }}</strong></div></div>
                            <div class="col-6"><div class="p-3 rounded border border-white border-opacity-10 bg-dark bg-opacity-25"><small class="text-secondary d-block">Beam ID</small><strong class="text-white">{{ ship.satellite }}</strong></div></div>
                        </div>
                        <div class="d-flex justify-content-between align-items-end mb-2">
                            <div><small class="text-uppercase text-secondary fw-bold">Live SNR</small><div class="display-4 fw-bold text-primary">{{ liveSnr }} <span class="fs-6 text-muted">dB</span></div></div>
                            <div class="spinner-grow text-success spinner-grow-sm" role="status"></div>
                        </div>
                        <div style="height: 180px;"><VsatChart :snr="liveSnr" /></div>
                    </div>
                </div>
                <div class="col-lg-8">
                    <div class="card glass-panel h-100 p-0 overflow-hidden position-relative">
                        <div class="position-absolute top-0 end-0 m-3 z-3 bg-dark p-1 rounded shadow-sm border border-secondary">
                            <div class="btn-group btn-group-sm">
                                <button class="btn" :class="mapMode==='internal'?'btn-primary':'btn-outline-secondary'" @click="mapMode='internal'">GPS</button>
                                <button class="btn" :class="mapMode==='marinetraffic'?'btn-primary':'btn-outline-secondary'" @click="mapMode='marinetraffic'">Traffic</button>
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
            <div v-if="activeTab === 'crew'" class="fade-in-tab">
                <div class="card glass-panel mb-4 p-4">
                    <div class="row g-3 align-items-end">
                        <div class="col-md-4"><label class="form-label small fw-bold text-secondary">Username:</label><input type="text" class="form-control" v-model="filter.username" placeholder="Search..."></div>
                        <div class="col-md-4"><label class="form-label small fw-bold text-secondary">Name:</label><input type="text" class="form-control" v-model="filter.name" placeholder="Search..."></div>
                        <div class="col-md-4 text-end"><button class="btn btn-outline-light me-2" @click="resetFilter">Reset</button><button class="btn btn-primary px-4" @click="loadCrew"><i class="fa-solid fa-magnifying-glass me-2"></i> Query</button></div>
                    </div>
                </div>

                <div class="card glass-panel p-0">
                    <div class="card-header border-bottom border-white border-opacity-10 py-3 d-flex justify-content-between align-items-center">
                        <ul class="nav nav-pills custom-pills">
                            <li class="nav-item"><a class="nav-link cursor-pointer" :class="{ active: filter.tab === 'Active' }" @click="filter.tab = 'Active'">Active</a></li>
                            <li class="nav-item"><a class="nav-link cursor-pointer" :class="{ active: filter.tab === 'Inactive' }" @click="filter.tab = 'Inactive'">Inactive</a></li>
                            <li class="nav-item"><a class="nav-link cursor-pointer" :class="{ active: filter.tab === 'All' }" @click="filter.tab = 'All'">All</a></li>
                        </ul>
                        <div class="d-flex gap-2">
                            <!-- NÚT CREATE USER: ĐÃ KIỂM TRA -->
                            <button class="btn btn-primary btn-sm fw-bold" @click="openCreateModal"><i class="fa-solid fa-user-plus me-2"></i> Create User</button>
                            <button class="btn btn-outline-light btn-sm" @click="loadCrew"><i class="fa-solid fa-rotate"></i></button>
                        </div>
                    </div>
                    <div class="table-responsive">
                        <table class="table table-hover align-middle mb-0">
                            <thead><tr><th class="ps-4">Username</th><th>Name / Rank</th><th>Plan</th><th>Usage</th><th>Status</th><th>Created</th><th class="text-end pe-4">Action</th></tr></thead>
                            <tbody>
                                <tr v-for="c in filteredCrew" :key="c.id">
                                    <td class="ps-4 fw-bold text-white">{{ c.username || 'N/A' }}</td>
                                    <td><div class="text-white">{{ c.full_name }}</div><small class="text-muted">{{ c.rank }}</small></td>
                                    <td><span class="badge bg-info bg-opacity-25 text-info border border-info border-opacity-25">{{ c.data_plan }}</span></td>
                                    <td><i class="fa-solid fa-chart-pie me-2 text-secondary"></i>{{ c.data_usage }} MB</td>
                                    <td><span class="badge rounded-pill" :class="c.status === 'Active' ? 'bg-success bg-opacity-25 text-success' : 'bg-secondary bg-opacity-25'">{{ c.status }}</span></td>
                                    <td class="text-secondary small">{{ formatDate(c.created_at) }}</td>
                                    <td class="text-end pe-4">
                                        <button class="btn btn-sm btn-link text-info" @click="openEditModal(c)"><i class="fa-solid fa-pen-to-square"></i></button>
                                        <button class="btn btn-sm btn-link text-danger" @click="handleDeleteCrew(c.id)"><i class="fa-solid fa-trash"></i></button>
                                    </td>
                                </tr>
                                <tr v-if="filteredCrew.length === 0"><td colspan="7" class="text-center py-5 text-muted">No data found.</td></tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>

            <!-- TAB 3: ROUTER -->
            <div v-if="activeTab === 'router'" class="fade-in-tab">
                <div class="card glass-panel mb-4 p-4">
                    <div class="d-flex justify-content-between align-items-center">
                        <div class="d-flex align-items-center">
                            <div class="p-3 bg-white rounded me-3"><img src="https://mikrotik.com/img/logo/logo-blue.svg" style="height: 24px;"></div>
                            <div><h5 class="fw-bold text-white m-0">{{ routerStats.board_name || 'Connecting...' }}</h5><small class="text-secondary font-monospace">{{ routerStats.version || 'RouterOS' }}</small></div>
                        </div>
                        <div class="text-end">
                            <div v-if="routerConnected" class="badge bg-success px-3 py-2"><i class="fa-solid fa-link me-2"></i> ONLINE</div>
                            <div v-else-if="isSimulation" class="badge bg-warning text-dark px-3 py-2"><i class="fa-solid fa-flask me-2"></i> SIMULATION</div>
                            <div v-else class="badge bg-danger px-3 py-2"><i class="fa-solid fa-unlink me-2"></i> DISCONNECTED</div>
                            <div class="text-secondary small mt-1 font-monospace">{{ routerIP || 'Unknown IP' }}</div>
                        </div>
                    </div>
                </div>
                <!-- ... (Giữ nguyên nội dung Router cũ) ... -->
                <div class="row g-4">
                    <div class="col-md-4">
                        <div class="card glass-panel h-100 p-4">
                            <h6 class="fw-bold text-white mb-4 border-bottom border-white border-opacity-10 pb-2">System Resources</h6>
                            <div class="mb-4"><div class="d-flex justify-content-between text-white mb-1"><small>CPU Load</small><small class="fw-bold">{{ routerStats.cpu_load }}%</small></div><div class="progress bg-dark" style="height: 8px;"><div class="progress-bar bg-info" :style="{width: routerStats.cpu_load + '%'}"></div></div></div>
                            <div class="mb-4"><div class="d-flex justify-content-between text-white mb-1"><small>Memory</small><small class="fw-bold">{{ routerStats.free_memory }}MB Free</small></div><div class="progress bg-dark" style="height: 8px;"><div class="progress-bar bg-warning" style="width: 60%"></div></div></div>
                            <div class="d-flex justify-content-between bg-dark p-3 rounded border border-white border-opacity-10"><span class="text-secondary small">Uptime</span><span class="text-white font-monospace">{{ routerStats.uptime || '--:--:--' }}</span></div>
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="card glass-panel h-100 p-4 text-center">
                            <h6 class="fw-bold text-white mb-4 border-bottom border-white border-opacity-10 pb-2">WAN Traffic</h6>
                            <div class="py-2">
                                <div class="mb-3"><small class="text-secondary text-uppercase fw-bold">Download</small><h3 class="text-success fw-bold m-0"><i class="fa-solid fa-arrow-down me-2"></i>{{ (routerStats.rx_rate / 1000000).toFixed(2) }} Mbps</h3></div>
                                <hr class="border-white border-opacity-10">
                                <div class="mb-3"><small class="text-secondary text-uppercase fw-bold">Upload</small><h3 class="text-primary fw-bold m-0"><i class="fa-solid fa-arrow-up me-2"></i>{{ (routerStats.tx_rate / 1000000).toFixed(2) }} Mbps</h3></div>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="card glass-panel h-100 p-4">
                            <h6 class="fw-bold text-white mb-4 border-bottom border-white border-opacity-10 pb-2">Management Tools</h6>
                            <button class="btn btn-outline-light w-100 mb-3 text-start py-3" @click="syncRouter"><i class="fa-solid fa-users-gear me-3 text-warning"></i> Sync Crew to Hotspot</button>
                            <button class="btn btn-outline-light w-100 mb-3 text-start py-3" @click="rebootRouter"><i class="fa-solid fa-power-off me-3 text-danger"></i> Reboot Router</button>
                            <button class="btn btn-outline-light w-100 text-start py-3" @click="blockInternet"><i class="fa-solid fa-ban me-3 text-danger"></i> Emergency Block</button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- MODAL CREW: ĐÃ ĐẶT ĐÚNG VỊ TRÍ, CÓ ref -->
            <div class="modal fade" id="crewModal" tabindex="-1" ref="crewModalRef">
                <div class="modal-dialog">
                    <div class="modal-content">
                        <div class="modal-header border-bottom border-white border-opacity-10"><h5 class="modal-title fw-bold">{{ isEditing ? 'Edit User' : 'Create User' }}</h5><button class="btn-close" data-bs-dismiss="modal"></button></div>
                        <div class="modal-body">
                            <div class="row g-3">
                                <div class="col-6"><label class="small fw-bold">Username</label><input v-model="crewForm.username" class="form-control"></div>
                                <div class="col-6"><label class="small fw-bold">Name</label><input v-model="crewForm.full_name" class="form-control"></div>
                                <div class="col-6"><label class="small fw-bold">Rank</label><select v-model="crewForm.rank" class="form-select"><option>Captain</option><option>Chief Officer</option><option>Crew</option></select></div>
                                <div class="col-6"><label class="small fw-bold">Plan</label><select v-model="crewForm.data_plan" class="form-select"><option>Basic (1GB)</option><option>Standard (5GB)</option><option>Unlimited</option></select></div>
                                <div class="col-12"><div class="form-check form-switch"><input class="form-check-input" type="checkbox" v-model="crewForm.status" true-value="Active" false-value="Inactive"><label class="form-check-label">Enable Access</label></div></div>
                            </div>
                            <div class="mt-4 text-end"><button class="btn btn-link text-white-50 me-2" data-bs-dismiss="modal">Cancel</button><button class="btn btn-primary fw-bold px-4" @click="handleSaveCrew">Save User</button></div>
                        </div>
                    </div>
                </div>
            </div>

        </div>
    </div>
</template>

<style scoped>
.cursor-pointer { cursor: pointer; }
.nav-link { color: #94a3b8; font-weight: 600; padding: 10px 20px; transition: all 0.3s; }
.nav-link:hover { color: white; }
.nav-link.active { color: #38bdf8; border-bottom: 2px solid #38bdf8; background: transparent; }
.custom-pills .nav-link { border-radius: 8px; border: 1px solid transparent; }
.custom-pills .nav-link.active { background: rgba(56, 189, 248, 0.15); color: #38bdf8; border: 1px solid rgba(56, 189, 248, 0.3); border-bottom: 1px solid rgba(56, 189, 248, 0.3); }
.uppercase-label { font-size: 0.7rem; letter-spacing: 1px; font-weight: 700; text-transform: uppercase; }
.tracking-wide { letter-spacing: 1px; }
.fade-in-tab { animation: fadeIn 0.4s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>