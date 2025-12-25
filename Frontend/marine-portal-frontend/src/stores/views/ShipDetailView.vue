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

const liveSnr = ref(0);
const mapMode = ref('internal');
const activeTab = ref('overview'); // State quản lý Tab
const crewList = ref([]); // Danh sách thủy thủ
let simulationInterval = null;
let crewModalInstance = null;
const crewModalRef = ref(null);

const ship = computed(() => store.ships.find(s => s.id === route.params.id));

// Form thêm Crew
const crewForm = reactive({
    full_name: '', rank: 'Captain', nationality: 'Vietnam', status: 'Onboard'
});

onMounted(async () => {
    if (store.ships.length === 0) await store.fetchFleet();
    if (ship.value) {
        liveSnr.value = ship.value.snr;
        loadCrew(); // Tải danh sách thủy thủ
    }
    
    // Khởi tạo Modal
    if(crewModalRef.value) crewModalInstance = new Modal(crewModalRef.value);

    simulationInterval = setInterval(() => {
        if (ship.value) {
            const base = ship.value.status === 'Offline' ? 0 : ship.value.snr;
            liveSnr.value = parseFloat((base + (Math.random() - 0.5)).toFixed(1));
        }
    }, 1000);
});

const loadCrew = async () => {
    if(ship.value) crewList.value = await store.fetchCrew(ship.value.id);
};

const openCrewModal = () => {
    Object.assign(crewForm, { full_name: '', rank: 'Captain', nationality: 'Vietnam', status: 'Onboard' });
    crewModalInstance.show();
};

const saveCrew = async () => {
    if(!crewForm.full_name) return alert("Nhập tên!");
    const success = await store.addCrew({ ...crewForm, ship_id: ship.value.id });
    if(success) {
        crewModalInstance.hide();
        loadCrew(); // Reload lại danh sách
    }
};

// Utils
const getMarineTrafficUrl = (ship) => {
    let id = ship.id.replace('IMO', '');
    if (id.length < 7) id = '9330264'; 
    return `https://www.marinetraffic.com/en/ais/embed/map/imo:${id}/latitude:${ship.lat}/longitude:${ship.lon}/zoom:9`;
};
const downloadPdf = () => window.open(`http://localhost:8080/api/report/${route.params.id}`, '_blank');
const formatCoord = (val, type) => `${val?.toFixed(4)}° ${type === 'lat' ? (val>0?'N':'S') : (val>0?'E':'W')}`;

onUnmounted(() => { if (simulationInterval) clearInterval(simulationInterval); });
</script>

<template>
    <div class="container-fluid p-0">
        <!-- HEADER -->
        <div class="d-flex justify-content-between align-items-center mb-4">
            <div class="d-flex align-items-center">
                <button class="btn btn-outline-secondary border-0 me-2" @click="router.push('/')">
                    <i class="fa-solid fa-arrow-left"></i>
                </button>
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
            <!-- TABS NAVIGATION -->
            <ul class="nav nav-tabs mb-4">
                <li class="nav-item">
                    <a class="nav-link cursor-pointer" :class="{ active: activeTab === 'overview' }" @click="activeTab = 'overview'">
                        <i class="fa-solid fa-chart-line me-2"></i>Overview
                    </a>
                </li>
                <li class="nav-item">
                    <a class="nav-link cursor-pointer" :class="{ active: activeTab === 'crew' }" @click="activeTab = 'crew'">
                        <i class="fa-solid fa-users me-2"></i>Crew List
                    </a>
                </li>
            </ul>

            <!-- TAB 1: OVERVIEW (Nội dung cũ) -->
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

            <!-- TAB 2: CREW LIST (Mới) -->
            <div v-if="activeTab === 'crew'" class="card border-0 shadow-sm">
                <div class="card-header bg-white py-3 d-flex justify-content-between align-items-center">
                    <h5 class="m-0 fw-bold">Onboard Crew Members</h5>
                    <button class="btn btn-primary btn-sm fw-bold" @click="openCrewModal"><i class="fa-solid fa-user-plus me-2"></i>Add Crew</button>
                </div>
                <div class="table-responsive">
                    <table class="table align-middle mb-0">
                        <thead class="bg-light">
                            <tr><th>Name</th><th>Rank</th><th>Nationality</th><th>Status</th></tr>
                        </thead>
                        <tbody>
                            <tr v-for="c in crewList" :key="c.id">
                                <td class="fw-bold">{{ c.full_name }}</td>
                                <td><span class="badge bg-info text-dark">{{ c.rank }}</span></td>
                                <td>{{ c.nationality }}</td>
                                <td><span class="badge bg-success-subtle text-success">{{ c.status }}</span></td>
                            </tr>
                            <tr v-if="crewList.length === 0"><td colspan="4" class="text-center py-4 text-muted">No crew data available.</td></tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <!-- MODAL ADD CREW -->
        <div class="modal fade" id="crewModal" tabindex="-1" ref="crewModalRef">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header"><h5 class="modal-title fw-bold">Add Crew Member</h5><button class="btn-close" data-bs-dismiss="modal"></button></div>
                    <div class="modal-body">
                        <div class="mb-3"><label class="small fw-bold">Full Name</label><input v-model="crewForm.full_name" class="form-control"></div>
                        <div class="mb-3">
                            <label class="small fw-bold">Rank</label>
                            <select v-model="crewForm.rank" class="form-select">
                                <option>Captain</option><option>Chief Officer</option><option>Chief Engineer</option><option>Able Seaman</option><option>Cook</option>
                            </select>
                        </div>
                        <div class="mb-3"><label class="small fw-bold">Nationality</label><input v-model="crewForm.nationality" class="form-control"></div>
                        <button class="btn btn-primary w-100 fw-bold" @click="saveCrew">Add Member</button>
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
</style>