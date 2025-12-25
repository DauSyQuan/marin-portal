<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useFleetStore } from '@/stores/fleet';
import BaseMap from '@/components/BaseMap.vue';
import VsatChart from '@/components/VsatChart.vue';

const route = useRoute();
const router = useRouter();
const store = useFleetStore();
const downloadPdf = () => {
    // Mở tab mới gọi thẳng vào API backend để tải file
    window.open(`http://localhost:8080/api/report/${route.params.id}`, '_blank');
};
// Biến cho biểu đồ và map
const liveSnr = ref(0);
const mapMode = ref('internal'); // 'internal' hoặc 'marinetraffic'
let simulationInterval = null;

// Tìm thông tin tàu
const ship = computed(() => store.ships.find(s => s.id === route.params.id));

// Hàm tạo link MarineTraffic Embed
const getMarineTrafficUrl = (ship) => {
    let id = ship.id.replace('IMO', '');
    // Nếu ID quá ngắn (dữ liệu ảo), dùng ID thật của tàu VIMC Diamond để demo
    if (id.length < 7) id = '9330264'; 
    return `https://www.marinetraffic.com/en/ais/embed/map/imo:${id}/latitude:${ship.lat}/longitude:${ship.lon}/zoom:9`;
};

onMounted(async () => {
    // 1. Nếu chưa có dữ liệu thì gọi API tải về
    if (store.ships.length === 0) await store.fetchFleet();
    
    // 2. Thiết lập giá trị ban đầu
    if (ship.value) {
        liveSnr.value = ship.value.snr;
    }

    // 3. Giả lập SNR nhảy múa (Simulation)
    simulationInterval = setInterval(() => {
        if (ship.value) {
            const base = ship.value.status === 'Offline' ? 0 : ship.value.snr;
            const fluctuation = (Math.random() - 0.5); 
            liveSnr.value = parseFloat((base + fluctuation).toFixed(1));
        }
    }, 1000);
});

onUnmounted(() => {
    if (simulationInterval) clearInterval(simulationInterval);
});
</script>

<template>
    <div class="container-fluid p-0">
        <!-- HEADER: Back Button & Title -->
        <div class="d-flex align-items-center mb-4">
            <button class="btn btn-white border me-3 shadow-sm bg-white" @click="router.push('/')">
                <i class="fa-solid fa-arrow-left me-1"></i> Back
            </button>
            <div class="me-auto"> <!-- Thêm div bọc tiêu đề để đẩy nút kia sang phải -->
        <div v-if="ship">
             <!-- ... code hiển thị tên tàu cũ ... -->
        </div>
    </div>
    
    <!-- NÚT DOWNLOAD PDF MỚI -->
    <button v-if="ship" class="btn btn-outline-primary shadow-sm" @click="downloadPdf">
        <i class="fa-solid fa-file-pdf me-2"></i>Download Report
    </button>
            <div v-if="ship">
                <div class="d-flex align-items-center">
                    <h3 class="fw-bold text-dark m-0 me-3">{{ ship.name }}</h3>
                    <span class="badge" :class="ship.status === 'Online' ? 'bg-success' : 'bg-danger'">
                        {{ ship.status }}
                    </span>
                </div>
                <span class="text-secondary small">IMO: {{ ship.id }} • {{ ship.type }}</span>
            </div>
        </div>

        <!-- Cảnh báo nếu không tìm thấy tàu -->
        <div v-if="!ship && !store.loading" class="alert alert-danger">
            Vessel Not Found. Please return to dashboard.
        </div>

        <!-- LAYOUT CHÍNH: 2 CỘT -->
        <div v-if="ship" class="row g-4">
            
            <!-- CỘT TRÁI: THÔNG TIN & BIỂU ĐỒ -->
            <div class="col-md-4">
                <div class="card shadow-sm border-0 h-100">
                    <div class="card-header bg-white fw-bold py-3">Vessel Information</div>
                    <div class="card-body d-flex flex-column">
                        <ul class="list-group list-group-flush mb-3">
                            <li class="list-group-item px-0 d-flex justify-content-between">
                                <span class="text-secondary">Owner</span> <strong class="text-dark">{{ ship.company }}</strong>
                            </li>
                            <li class="list-group-item px-0 d-flex justify-content-between">
                                <span class="text-secondary">IP Address</span> <strong class="font-monospace text-primary">{{ ship.ip }}</strong>
                            </li>
                            <li class="list-group-item px-0 d-flex justify-content-between">
                                <span class="text-secondary">Satellite</span> <strong class="text-primary">{{ ship.satellite }}</strong>
                            </li>
                        </ul>
                        
                        <!-- Chart Section -->
                        <div class="mt-auto">
                            <h6 class="small fw-bold text-secondary text-uppercase mb-3 border-top pt-3">
                                <i class="fa-solid fa-wave-square me-2"></i>Live VSAT Monitor
                            </h6>
                            <div class="d-flex align-items-end mb-2">
                                <div class="display-4 fw-bold text-primary me-2">{{ liveSnr }}</div>
                                <div class="mb-2 text-secondary">dB</div>
                                <div class="ms-auto mb-2 badge bg-success bg-opacity-10 text-success animate-pulse">
                                    <i class="fa-solid fa-circle fa-2xs me-1"></i> Live
                                </div>
                            </div>
                            <div style="height: 180px; margin-left: -10px; margin-right: -10px;">
                                <VsatChart :snr="liveSnr" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- CỘT PHẢI: BẢN ĐỒ -->
            <div class="col-md-8">
                <div class="card shadow-sm border-0 h-100 position-relative">
                    
                    <!-- Map Control Buttons -->
                    <div class="position-absolute top-0 end-0 m-3 z-3 bg-white p-1 rounded shadow-sm">
                        <div class="btn-group btn-group-sm">
                            <button class="btn" :class="mapMode === 'internal' ? 'btn-primary' : 'btn-light'" @click="mapMode = 'internal'">Internal GPS</button>
                            <button class="btn" :class="mapMode === 'marinetraffic' ? 'btn-primary' : 'btn-light'" @click="mapMode = 'marinetraffic'">MarineTraffic</button>
                        </div>
                    </div>

                    <div class="card-body p-0 h-100" style="min-height: 500px; overflow: hidden; border-radius: 8px;">
                        <!-- Mode 1: Internal Map -->
                        <BaseMap v-if="mapMode === 'internal'" :lat="ship.lat" :lon="ship.lon" />
                        
                        <!-- Mode 2: MarineTraffic -->
                        <iframe v-else
                            width="100%" height="100%" scrolling="no" frameborder="0"
                            :src="getMarineTrafficUrl(ship)"
                            style="border:0">
                        </iframe>
                    </div>
                </div>
            </div>
            
        </div>
    </div>
</template>