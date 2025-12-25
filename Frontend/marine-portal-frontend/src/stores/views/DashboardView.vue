<script setup>
import { onMounted, reactive, ref, nextTick } from 'vue';
import { useFleetStore } from '@/stores/fleet';
import { Modal } from 'bootstrap';

const store = useFleetStore();
let modalInstance = null;
const modalRef = ref(null);

// Form dữ liệu
const form = reactive({
    id: '', name: '', company: 'VIMC Lines', 
    type: 'Container', ip: '10.10.x.x', 
    satellite: 'JCSAT-4B', beam: 'Spot-1',
    lat: 10.0, lon: 106.0
});

// --- THAY ĐỔI QUAN TRỌNG Ở ĐÂY ---
onMounted(() => {
    // 1. Gọi API tải dữ liệu (KHÔNG DÙNG AWAIT)
    // Để giao diện hiện lên ngay lập tức, dữ liệu về sau thì tự điền vào
    store.fetchFleet(); 
    
    // 2. Khởi tạo Modal (Dùng nextTick để đảm bảo HTML đã vẽ xong)
    nextTick(() => {
        if(modalRef.value) {
            modalInstance = new Modal(modalRef.value);
        }
    });
});

// --- LOGIC MODAL ---
const openAddModal = () => {
    Object.assign(form, { 
        id: '', name: '', company: 'VIMC Lines', 
        type: 'Container', ip: '10.10.x.x', 
        satellite: 'JCSAT-4B', beam: 'Spot-1',
        lat: 10.0, lon: 106.0 
    });
    modalInstance.show();
};

const saveShip = async () => {
    if(!form.id || !form.name) return alert("Vui lòng nhập ID và Tên tàu!");
    
    const payload = {
        ...form,
        lat: parseFloat(form.lat || 0),
        lon: parseFloat(form.lon || 0),
        snr: 12.0, status: 'Online'
    };

    const success = await store.addShip(payload);
    if(success) {
        modalInstance.hide();
        Object.assign(form, { id: '', name: '', lat: 0, lon: 0 });
    }
};

// --- UTILS ---
const getShipIcon = (type) => {
    if(type?.includes('Container')) return 'fa-solid fa-ship';
    if(type?.includes('Tanker') || type?.includes('Oil')) return 'fa-solid fa-oil-well';
    if(type?.includes('Bulk')) return 'fa-solid fa-cubes-stacked';
    if(type?.includes('Drill')) return 'fa-solid fa-screwdriver-wrench';
    return 'fa-solid fa-anchor';
};

const getSignalPercent = (snr) => {
    const val = parseFloat(snr);
    if(isNaN(val) || val < 0) return 0;
    return Math.min((val / 20) * 100, 100);
};

const getSignalColor = (snr) => {
    if(snr >= 12) return 'bg-success';
    if(snr >= 5) return 'bg-warning';
    return 'bg-danger';
};

const formatCoord = (lat, lon) => {
    return `${lat?.toFixed(2)}°N, ${lon?.toFixed(2)}°E`;
};

const statusClass = (status) => {
    if(status === 'Online') return 'bg-success-subtle text-success border-success';
    if(status === 'Warning') return 'bg-warning-subtle text-warning border-warning';
    return 'bg-danger-subtle text-danger border-danger';
};
</script>

<template>
  <div class="container-fluid p-0">
    <!-- HEADER -->
    <div class="d-flex justify-content-between align-items-center mb-4 pb-3 border-bottom">
        <div>
            <div class="d-flex align-items-center gap-2 mb-1">
                <span class="badge bg-primary text-white">LIVE</span>
                <small class="text-secondary fw-bold text-uppercase tracking-wide">Global Fleet Ops</small>
            </div>
            <h3 class="fw-bold text-dark m-0">Command Center <span class="text-primary">.</span></h3>
        </div>
        <div class="d-flex gap-2">
            <button class="btn btn-light border shadow-sm fw-bold text-secondary" @click="store.fetchFleet">
                <i class="fa-solid fa-rotate me-2"></i> Sync Data
            </button>
            <button class="btn btn-primary shadow fw-bold px-4" @click="openAddModal">
                <i class="fa-solid fa-plus me-2"></i> Commission Vessel
            </button>
        </div>
    </div>

    <!-- NỘI DUNG CHÍNH (Đã bỏ v-if isReady để hiện luôn) -->
    
    <!-- 1. KPI CARDS -->
    <div class="row mb-4 g-4">
        <div class="col-md-3">
            <div class="card border-0 shadow-sm h-100 overflow-hidden position-relative">
                <div class="card-body">
                    <div class="d-flex justify-content-between mb-3">
                        <div class="text-secondary fw-bold small text-uppercase">Total Vessels</div>
                        <span class="badge bg-primary-subtle text-primary rounded-pill">+2 New</span>
                    </div>
                    <!-- Thêm check store.stats có dữ liệu chưa -->
                    <h2 class="display-5 fw-bold text-dark mb-0">{{ store.stats ? store.stats[0].value : 0 }}</h2>
                    <small class="text-secondary">Active in fleet</small>
                </div>
                <i class="fa-solid fa-ship position-absolute bottom-0 end-0 text-primary opacity-10" style="font-size: 5rem; margin: -10px;"></i>
            </div>
        </div>

        <div class="col-md-3">
            <div class="card border-0 shadow-sm h-100">
                <div class="card-body">
                    <div class="d-flex justify-content-between mb-3">
                        <div class="text-secondary fw-bold small text-uppercase">Network Health</div>
                        <i class="fa-solid fa-heart-pulse text-success"></i>
                    </div>
                    <h2 class="display-5 fw-bold text-success mb-0">{{ store.stats ? store.stats[1].value : '0%' }}</h2>
                    <div class="progress mt-2" style="height: 4px;">
                        <div class="progress-bar bg-success" style="width: 95%"></div>
                    </div>
                    <small class="text-secondary mt-1 d-block">SLA Met</small>
                </div>
            </div>
        </div>

        <div class="col-md-3" v-for="(stat, idx) in (store.stats ? store.stats.slice(2) : [])" :key="idx">
            <div class="card border-0 shadow-sm h-100">
                <div class="card-body">
                    <div class="d-flex justify-content-between mb-3">
                        <div class="text-secondary fw-bold small text-uppercase">{{ stat.label }}</div>
                        <i :class="stat.icon + ' text-' + stat.color"></i>
                    </div>
                    <h2 class="display-5 fw-bold mb-0" :class="'text-' + stat.color">{{ stat.value }}</h2>
                    <small class="text-muted">Requires attention</small>
                </div>
            </div>
        </div>
    </div>

    <!-- 2. MAIN TABLE -->
    <div class="card border-0 shadow-sm">
        <div class="card-header bg-white py-3 d-flex justify-content-between align-items-center border-bottom-0">
            <div class="d-flex align-items-center gap-2">
                <h5 class="fw-bold m-0">Vessel Monitoring</h5>
                <span class="badge bg-light text-secondary border">{{ store.ships.length }} Units</span>
            </div>
            <div class="input-group" style="width: 300px;">
                <span class="input-group-text bg-light border-end-0"><i class="fa-solid fa-search text-secondary"></i></span>
                <input type="text" class="form-control bg-light border-start-0 ps-0" v-model="store.filterText" placeholder="Search by ID, Name or IP...">
            </div>
        </div>

        <div class="table-responsive">
            <table class="table table-hover align-middle mb-0 custom-table">
                <thead class="bg-light text-uppercase">
                    <tr>
                        <th class="ps-4 text-secondary small fw-bold" style="width: 250px;">Vessel Info</th>
                        <th class="text-secondary small fw-bold">Telemetry / Position</th>
                        <th class="text-secondary small fw-bold">Network Link</th>
                        <th class="text-secondary small fw-bold">Signal Quality (SNR)</th>
                        <th class="text-end pe-4 text-secondary small fw-bold">Controls</th>
                    </tr>
                </thead>
                <tbody>
                    <!-- HIỆN SKELETON LOADING KHI CHƯA CÓ DỮ LIỆU -->
                    <tr v-if="store.loading">
                        <td colspan="5" class="text-center py-5 text-muted">
                            <div class="spinner-border text-primary mb-2" role="status"></div>
                            <div>Loading Fleet Data...</div>
                        </td>
                    </tr>

                    <tr v-else v-for="ship in store.filteredShips" :key="ship.id">
                        <td class="ps-4 py-3">
                            <div class="d-flex align-items-center">
                                <div class="avatar-icon bg-light text-primary rounded p-3 me-3">
                                    <i :class="getShipIcon(ship.type) + ' fa-lg'"></i>
                                </div>
                                <div>
                                    <div class="fw-bold text-dark">{{ ship.name }}</div>
                                    <div class="small text-muted font-monospace">{{ ship.id }}</div>
                                    <div class="badge bg-light text-secondary border mt-1" style="font-size: 0.65rem;">{{ ship.type }}</div>
                                </div>
                            </div>
                        </td>
                        <td>
                            <div class="d-flex flex-column">
                                <span class="font-monospace text-dark small mb-1">
                                    <i class="fa-solid fa-location-crosshairs text-secondary me-2"></i>
                                    {{ formatCoord(ship.lat, ship.lon) }}
                                </span>
                                <span class="small text-muted"><i class="fa-regular fa-clock me-2"></i>Updated just now</span>
                            </div>
                        </td>
                        <td>
                            <div class="mb-1">
                                <span class="badge border" :class="statusClass(ship.status)">
                                    <i class="fa-solid fa-circle fa-2xs me-1"></i> {{ ship.status }}
                                </span>
                            </div>
                            <div class="small font-monospace text-secondary">IP: {{ ship.ip }}</div>
                        </td>
                        <td style="width: 200px;">
                            <div class="d-flex justify-content-between align-items-end mb-1">
                                <span class="fw-bold small text-dark">{{ Number(ship.snr).toFixed(1) }} dB</span>
                                <i class="fa-solid fa-signal small" :class="ship.snr > 10 ? 'text-success' : 'text-warning'"></i>
                            </div>
                            <div class="progress" style="height: 6px;">
                                <div class="progress-bar" :class="getSignalColor(ship.snr)" :style="{ width: getSignalPercent(ship.snr) + '%' }"></div>
                            </div>
                            <small class="text-muted" style="font-size: 0.7rem;">Beam: {{ ship.beam }}</small>
                        </td>
                        <td class="text-end pe-4">
                            <button class="btn btn-sm btn-outline-dark fw-bold" @click="$router.push('/ship/' + ship.id)">
                                Monitor <i class="fa-solid fa-arrow-right ms-1"></i>
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>

    <!-- MODAL -->
    <div class="modal fade" id="addVesselModal" tabindex="-1" ref="modalRef">
        <div class="modal-dialog">
            <div class="modal-content border-0 shadow-lg">
                <div class="modal-header bg-dark text-white">
                    <h5 class="modal-title fw-bold">Commission New Vessel</h5>
                    <button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body bg-light">
                    <form @submit.prevent="saveShip">
                        <div class="mb-3"><label class="form-label small fw-bold">IMO Number</label><input v-model="form.id" class="form-control" placeholder="e.g. IMO999888"></div>
                        <div class="mb-3"><label class="form-label small fw-bold">Vessel Name</label><input v-model="form.name" class="form-control" placeholder="e.g. PACIFIC STAR"></div>
                        <div class="text-end mt-4">
                            <button type="button" class="btn btn-link text-secondary text-decoration-none me-2" data-bs-dismiss="modal">Cancel</button>
                            <button type="button" @click="saveShip" class="btn btn-primary fw-bold px-4">Commission Vessel</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>

  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500&display=swap');
.font-monospace { font-family: 'JetBrains Mono', monospace; }

.custom-table thead th {
    font-size: 0.75rem; background-color: #f8fafc; border-bottom: 2px solid #e2e8f0; padding: 1rem;
}
.custom-table tbody tr:hover {
    background-color: #f1f5f9; transform: translateY(-1px); box-shadow: 0 4px 6px -1px rgba(0,0,0,0.05); position: relative; z-index: 1;
}
.avatar-icon {
    width: 48px; height: 48px; display: flex; align-items: center; justify-content: center; background-color: #e0f2fe;
}
.bg-success-subtle { background-color: #dcfce7 !important; color: #166534 !important; }
.bg-warning-subtle { background-color: #fef9c3 !important; color: #854d0e !important; }
.bg-danger-subtle { background-color: #fee2e2 !important; color: #991b1b !important; }
.tracking-wide { letter-spacing: 0.1em; }
</style>