<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue';
import L from 'leaflet';

// Nhận vĩ độ (lat), kinh độ (lon) từ trang cha truyền vào
const props = defineProps(['lat', 'lon']);

const mapContainer = ref(null);
let mapInstance = null;

onMounted(() => {
    if (mapContainer.value) {
        // 1. Khởi tạo bản đồ
        mapInstance = L.map(mapContainer.value).setView([props.lat, props.lon], 6);

        // 2. Thêm lớp nền bản đồ (Miễn phí từ CartoDB)
        L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png', {
            attribution: '&copy; Marine Pro',
            maxZoom: 18
        }).addTo(mapInstance);

        // 3. Vẽ Marker vị trí tàu
        drawMarker(props.lat, props.lon);
    }
});

// Hàm vẽ marker tròn
const drawMarker = (lat, lon) => {
    if (!mapInstance) return;
    
    // Xóa marker cũ nếu có (để tránh vẽ chồng lên nhau)
    mapInstance.eachLayer((layer) => {
        if (layer instanceof L.CircleMarker) {
            mapInstance.removeLayer(layer);
        }
    });

    // Vẽ marker mới
    L.circleMarker([lat, lon], {
        color: '#2563eb',      // Viền xanh
        fillColor: '#2563eb',  // Nền xanh
        fillOpacity: 0.8,
        radius: 10
    }).addTo(mapInstance);
};

// Theo dõi nếu tọa độ thay đổi thì cập nhật bản đồ
watch(() => [props.lat, props.lon], ([newLat, newLon]) => {
    if (mapInstance) {
        mapInstance.setView([newLat, newLon], 6);
        drawMarker(newLat, newLon);
    }
});

// Hủy bản đồ khi rời trang để giải phóng bộ nhớ
onUnmounted(() => {
    if (mapInstance) mapInstance.remove();
});
</script>

<template>
    <!-- Khung chứa bản đồ -->
    <div ref="mapContainer" class="h-100 w-100" style="min-height: 400px; border-radius: 8px; z-index: 1;"></div>
</template>