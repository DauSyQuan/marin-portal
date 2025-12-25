<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useFleetStore } from '@/stores/fleet';
import BaseMap from '@/components/BaseMap.vue';
import VsatChart from '@/components/VsatChart.vue';

const route = useRoute();
const router = useRouter();
const store = useFleetStore();

const liveSnr = ref(0);
const mapMode = ref('internal');
let simulationInterval = null;

const ship = computed(() => store.ships.find(s => s.id === route.params.id));

const getMarineTrafficUrl = (ship) => {
    let id = ship.id.replace('IMO', '');
    if (id.length < 7) id = '9330264'; 
    return `https://www.marinetraffic.com/en/ais/embed/map/imo:${id}/latitude:${ship.lat}/longitude:${ship.lon}/zoom:9`;
};

const downloadPdf = () => {
    window.open(`http://localhost:8080/api/report/${route.params.id}`, '_blank');
};

const formatCoord = (val, type) => {
    return `${val?.toFixed(4)}° ${type === 'lat' ? (val>0?'N':'S') : (val>0?'E':'W')}`;
};

onMounted(async () => {
    if (store.ships.length === 0) await store.fetchFleet();
    if (ship.value) liveSnr.value = ship.value.snr;

    simulationInterval = setInterval(() => {
        if (ship.value) {
            const base = ship.value.status === 'Offline' ? 0 : ship.value.snr;
            liveSnr.value = parseFloat((base + (Math.random() - 0.5)).toFixed(1));
        }
    }, 1000);
});

onUnmounted(() => {
    if (simulationInterval) clearInterval(simulationInterval);
});
</script>

<template>
    <div class="container-fluid p-0">
        <!-- HEADER: BREADCRUMB & CONTROLS -->
        <div class="d-flex justify-content-between align-items-center mb-4">
            <div class="d-flex align-items-center">
                <button class="btn btn-outline-secondary border-0 me-2" @click="router.push('/')">
                    <i class="fa-solid fa-arrow-left"></i>
                </button>
                <div>
                    <div class="text-uppercase text-secondary" style="font-size: 0.65rem; letter-spacing: 1px;">Vessel Monitor</div>
                    <h3 class="fw-bold m-0 d-flex align-items-center">
                        {{ ship?.name }}
                        <span v-if="ship" class="badge ms-3 fs-6 rounded-pill" 
                            :class="ship.status === 'Online' ? 'bg-success' : 'bg-danger'">
                            {{ ship.status }}
                        </span>
                    </h3>
                </div>
            </div>
            <div class="d-flex gap-2">
                <button class="btn btn-white border shadow-sm fw-bold text-secondary">
                    <i class="fa-solid fa-envelope me-2"></i> Email Logs
                </button>
                <button class="btn btn-primary shadow fw-bold" @click="downloadPdf">
                    <i class="fa-solid fa-file-pdf me-2"></i> Voyage Report
                </button>
            </div>
        </div>

        <div v-if="!ship" class="alert alert-warning">Loading vessel data...</div>

        <div v-if="ship" class="row g-4">
            
            <!-- TELEMETRY STRIP (Dải thông số kỹ thuật) -->
            <div class="col-12">
                <div class="card border-0 shadow-sm bg-dark text-white">
                    <div class="card-body d-flex justify-content-between align-items-center py-3">
                        <div class="d-flex gap-5">
                            <div>
                                <small class="text-white-50 text-uppercase d-block" style="font-size: 0.7rem;">IMO Number</small>
                                <span class="font-monospace fw-bold">{{ ship.id }}</span>
                            </div>
                            <div>
                                <small class="text-white-50 text-uppercase d-block" style="font-size: 0.7rem;">Position (Lat/Lon)</small>
                                <span class="font-monospace text-info">
                                    <i class="fa-solid fa-location-crosshairs me-1"></i>
                                    {{ formatCoord(ship.lat, 'lat') }} / {{ formatCoord(ship.lon, 'lon') }}
                                </span>
                            </div>
                            <div>
                                <small class="text-white-50 text-uppercase d-block" style="font-size: 0.7rem;">Heading / Speed</small>
                                <span class="font-monospace text-warning">124° / 14.2 kn</span>
                            </div>
                        </div>
                        <div class="border-start border-white border-opacity-10 ps-4">
                            <small class="text-white-50 text-uppercase d-block" style="font-size: 0.7rem;">Last Contact</small>
                            <span class="text-white"><i class="fa-solid fa-satellite me-2"></i>Just now</span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- CỘT TRÁI: KẾT NỐI & BIỂU ĐỒ -->
            <div class="col-lg-4">
                <div class="card shadow-sm border-0 h-100">
                    <div class="card-header bg-white py-3 fw-bold border-bottom">
                        <i class="fa-solid fa-network-wired me-2 text-primary"></i>Connection Health
                    </div>
                    <div class="card-body d-flex flex-column">
                        <!-- Info Grid -->
                        <div class="row g-3 mb-4">
                            <div class="col-6">
                                <div class="p-3 bg-light rounded border">
                                    <small class="text-secondary d-block">IP Address</small>
                                    <strong class="font-monospace text-dark">{{ ship.ip }}</strong>
                                </div>
                            </div>
                            <div class="col-6">
                                <div class="p-3 bg-light rounded border">
                                    <small class="text-secondary d-block">Satellite Beam</small>
                                    <strong class="text-dark">{{ ship.satellite }}</strong>
                                </div>
                            </div>
                        </div>

                        <!-- Live Chart -->
                        <div class="mt-auto">
                            <div class="d-flex justify-content-between align-items-end mb-2">
                                <div>
                                    <small class="text-uppercase text-secondary fw-bold">Live SNR</small>
                                    <div class="display-4 fw-bold text-primary" style="line-height: 1;">{{ liveSnr }} <span class="fs-4 text-muted">dB</span></div>
                                </div>
                                <div class="spinner-grow text-success spinner-grow-sm" role="status"></div>
                            </div>
                            <div style="height: 200px; margin: 0 -15px;">
                                <VsatChart :snr="liveSnr" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- CỘT PHẢI: BẢN ĐỒ -->
            <div class="col-lg-8">
                <div class="card shadow-sm border-0 h-100 position-relative overflow-hidden">
                    <div class="position-absolute top-0 end-0 m-3 z-3 bg-white p-1 rounded shadow-sm border">
                        <div class="btn-group btn-group-sm">
                            <button class="btn fw-bold" :class="mapMode === 'internal' ? 'btn-primary' : 'btn-light'" @click="mapMode = 'internal'">GPS Mode</button>
                            <button class="btn fw-bold" :class="mapMode === 'marinetraffic' ? 'btn-primary' : 'btn-light'" @click="mapMode = 'marinetraffic'">Traffic View</button>
                        </div>
                    </div>
                    <div class="h-100 w-100" style="min-height: 500px;">
                        <BaseMap v-if="mapMode === 'internal'" :lat="ship.lat" :lon="ship.lon" />
                        <iframe v-else width="100%" height="100%" :src="getMarineTrafficUrl(ship)" style="border:0"></iframe>
                    </div>
                </div>
            </div>

        </div>
    </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&display=swap');
.font-monospace { font-family: 'JetBrains Mono', monospace; }
</style>