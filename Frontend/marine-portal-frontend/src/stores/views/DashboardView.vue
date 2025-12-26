<script setup>
import { onMounted, reactive, ref, nextTick } from 'vue';
import { useFleetStore } from '@/stores/fleet';
import { Modal } from 'bootstrap';

const store = useFleetStore();
let modalInstance = null;
const modalRef = ref(null);

const form = reactive({ id: '', name: '', company: 'VIMC Lines', type: 'Container', ip: '10.10.x.x', satellite: 'JCSAT-4B', beam: 'Spot-1', lat: 10.0, lon: 106.0 });

onMounted(() => {
    store.fetchFleet(); 
    nextTick(() => { if(modalRef.value) modalInstance = new Modal(modalRef.value); });
});

const openAddModal = () => { Object.assign(form, { id: '', name: '', company: 'VIMC Lines', type: 'Container', ip: '10.10.x.x', satellite: 'JCSAT-4B', beam: 'Spot-1', lat: 10.0, lon: 106.0 }); modalInstance.show(); };
const saveShip = async () => {
    if(!form.id || !form.name) return alert("Vui lòng nhập ID và Tên tàu!");
    const payload = { ...form, lat: parseFloat(form.lat||0), lon: parseFloat(form.lon||0), snr: 12.0, status: 'Online' };
    if(await store.addShip(payload)) { modalInstance.hide(); Object.assign(form, { id: '', name: '', lat: 0, lon: 0 }); }
};

const getShipIcon = (type) => type?.includes('Container') ? 'fa-solid fa-ship' : type?.includes('Tanker') ? 'fa-solid fa-oil-well' : 'fa-solid fa-anchor';
const getSignalPercent = (snr) => Math.min((parseFloat(snr||0)/20)*100, 100);
const getSignalColor = (snr) => snr >= 12 ? 'bg-success' : snr >= 5 ? 'bg-warning' : 'bg-danger';
const formatCoord = (lat, lon) => `${lat?.toFixed(2)}°N, ${lon?.toFixed(2)}°E`;
const statusClass = (status) => status === 'Online' ? 'text-success border-success bg-success bg-opacity-10' : status === 'Warning' ? 'text-warning border-warning bg-warning bg-opacity-10' : 'text-danger border-danger bg-danger bg-opacity-10';
</script>

<template>
  <div class="container-fluid p-0">
    <!-- HEADER -->
    <div class="d-flex justify-content-between align-items-center mb-4 pb-3 border-bottom border-white border-opacity-10">
        <div>
            <div class="d-flex align-items-center gap-2 mb-1">
                <span class="badge bg-primary glow-badge">LIVE SATELLITE</span>
                <small class="text-info fw-bold text-uppercase tracking-wide">Global Fleet Ops</small>
            </div>
            <h3 class="fw-bold text-white m-0">Command Center</h3>
        </div>
        <div class="d-flex gap-2">
            <button class="btn btn-outline-light border-opacity-25 shadow-sm fw-bold" @click="store.fetchFleet"><i class="fa-solid fa-rotate me-2"></i> Sync</button>
            <button class="btn btn-primary fw-bold px-4 btn-glow" @click="openAddModal"><i class="fa-solid fa-plus me-2"></i> Commission</button>
        </div>
    </div>

    <!-- KPI CARDS (ÁP DỤNG GLASS PANEL) -->
    <div class="row mb-4 g-4">
        <div class="col-md-3">
            <div class="card h-100 glass-panel p-3">
                <div class="d-flex justify-content-between mb-2">
                    <div class="text-white-50 small text-uppercase fw-bold">Total Vessels</div>
                    <span class="badge bg-primary bg-opacity-25 text-primary border border-primary">+2 New</span>
                </div>
                <h2 class="display-5 fw-bold text-white mb-0">{{ store.stats ? store.stats[0].value : 0 }}</h2>
            </div>
        </div>
        <div class="col-md-3">
            <div class="card h-100 glass-panel p-3">
                <div class="d-flex justify-content-between mb-2">
                    <div class="text-white-50 small text-uppercase fw-bold">Network Health</div>
                    <i class="fa-solid fa-heart-pulse text-success"></i>
                </div>
                <h2 class="display-5 fw-bold text-success mb-0">{{ store.stats ? store.stats[1].value : '0%' }}</h2>
                <div class="progress mt-2 bg-dark" style="height: 4px;"><div class="progress-bar bg-success" style="width: 95%"></div></div>
            </div>
        </div>
        <div class="col-md-3" v-for="(stat, idx) in (store.stats ? store.stats.slice(2) : [])" :key="idx">
            <div class="card h-100 glass-panel p-3">
                <div class="d-flex justify-content-between mb-2">
                    <div class="text-white-50 small text-uppercase fw-bold">{{ stat.label }}</div>
                    <i :class="stat.icon + ' text-' + stat.color"></i>
                </div>
                <h2 class="display-5 fw-bold mb-0" :class="'text-' + stat.color">{{ stat.value }}</h2>
            </div>
        </div>
    </div>

    <!-- TABLE (ÁP DỤNG GLASS PANEL) -->
    <div class="card glass-panel">
        <div class="card-header border-bottom border-white border-opacity-10 py-3 d-flex justify-content-between align-items-center">
            <div class="d-flex align-items-center gap-2">
                <h5 class="fw-bold m-0 text-white">Vessel Monitoring</h5>
                <span class="badge bg-white bg-opacity-10 text-white border border-white border-opacity-10">{{ store.ships.length }} Units</span>
            </div>
            <div class="input-group" style="width: 300px;">
                <span class="input-group-text"><i class="fa-solid fa-search"></i></span>
                <input type="text" class="form-control" v-model="store.filterText" placeholder="Search ID, Name...">
            </div>
        </div>

        <div class="table-responsive">
            <table class="table table-hover align-middle mb-0">
                <thead>
                    <tr>
                        <th class="ps-4">Vessel Info</th><th>Telemetry</th><th>Network</th><th>Signal (SNR)</th><th class="text-end pe-4">Controls</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="store.loading"><td colspan="5" class="text-center py-5 text-white-50">Loading Data...</td></tr>
                    <tr v-else v-for="ship in store.filteredShips" :key="ship.id">
                        <td class="ps-4 py-3">
                            <div class="d-flex align-items-center">
                                <div class="avatar-glass me-3"><i :class="getShipIcon(ship.type) + ' fa-lg text-info'"></i></div>
                                <div>
                                    <div class="fw-bold text-white">{{ ship.name }}</div>
                                    <div class="small text-white-50 font-monospace">{{ ship.id }}</div>
                                </div>
                            </div>
                        </td>
                        <td>
                            <div class="d-flex flex-column">
                                <span class="font-monospace text-info small mb-1"><i class="fa-solid fa-crosshairs me-2"></i>{{ formatCoord(ship.lat, ship.lon) }}</span>
                                <span class="small text-white-50">Updated just now</span>
                            </div>
                        </td>
                        <td>
                            <span class="badge border mb-1" :class="statusClass(ship.status)">● {{ ship.status }}</span>
                            <div class="small font-monospace text-white-50">{{ ship.ip }}</div>
                        </td>
                        <td style="width: 200px;">
                            <div class="d-flex justify-content-between align-items-end mb-1">
                                <span class="fw-bold small text-white">{{ Number(ship.snr).toFixed(1) }} dB</span>
                                <i class="fa-solid fa-signal small" :class="ship.snr > 10 ? 'text-success' : 'text-warning'"></i>
                            </div>
                            <div class="progress bg-dark border border-white border-opacity-10" style="height: 6px;">
                                <div class="progress-bar" :class="getSignalColor(ship.snr)" :style="{ width: getSignalPercent(ship.snr) + '%' }"></div>
                            </div>
                        </td>
                        <td class="text-end pe-4">
                            <button class="btn btn-sm btn-outline-info" @click="$router.push('/ship/' + ship.id)">Monitor</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>

    <!-- MODAL (Giữ nguyên logic) -->
    <div class="modal fade" id="addVesselModal" tabindex="-1" ref="modalRef">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header border-bottom border-white border-opacity-10">
                    <h5 class="modal-title fw-bold">Commission New Vessel</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body">
                    <form @submit.prevent="saveShip">
                        <div class="mb-3"><label class="small fw-bold">IMO Number</label><input v-model="form.id" class="form-control" placeholder="e.g. IMO999"></div>
                        <div class="mb-3"><label class="small fw-bold">Vessel Name</label><input v-model="form.name" class="form-control" placeholder="e.g. PACIFIC STAR"></div>
                        <div class="text-end mt-4">
                            <button type="button" class="btn btn-link text-white-50 text-decoration-none me-2" data-bs-dismiss="modal">Cancel</button>
                            <button type="button" @click="saveShip" class="btn btn-primary fw-bold px-4 btn-glow">Add Vessel</button>
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
.avatar-glass {
    width: 44px; height: 44px; display: flex; align-items: center; justify-content: center;
    background: rgba(255,255,255,0.05); border-radius: 12px; border: 1px solid rgba(255,255,255,0.1);
}
.btn-glow { box-shadow: 0 0 15px rgba(56, 189, 248, 0.4); }
.glow-badge { box-shadow: 0 0 10px rgba(37, 99, 235, 0.5); }
.tracking-wide { letter-spacing: 1px; }
</style>