<script setup>
import { computed, onMounted, onUnmounted, ref, reactive, nextTick } from 'vue'; // Thêm nextTick
import { useRoute, useRouter } from 'vue-router';
import { useFleetStore } from '@/stores/fleet';
import BaseMap from '@/components/BaseMap.vue';
import VsatChart from '@/components/VsatChart.vue';
import CrewManagementModal from '@/components/CrewManagementModal.vue';
import api from '@/services/api';

const route = useRoute();
const router = useRouter();
const fleetStore = useFleetStore();

const ship = computed(() => fleetStore.getShipById(route.params.id));
const routerHealth = ref(null);
const isLoading = ref(false);

// Crew Management Modal
const showCrewModal = ref(false);

// Performance Chart Data
let chartUpdateInterval = null;

// Download PDF Report
const downloadPdf = () =>
  window.open(
    `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/report/${route.params.id}`,
    '_blank'
  );

// --- Fetch Router Stats ---
const fetchRouterStats = async () => {
  try {
    const res = await api.get(`/api/ships/${route.params.id}/router/stats`);
    routerHealth.value = res.data;
  } catch (e) {
    console.error(e);
    routerHealth.value = null;
  }
};

// --- Sync Crew to Router ---
const syncCrewToRouter = async () => {
  isLoading.value = true;
  try {
    await api.post(`/api/ships/${route.params.id}/router/sync`);
    alert("✅ Sync Crew thành công!");
  } catch (e) {
    console.error(e);
    alert("❌ Sync Crew thất bại: " + (e?.response?.data?.error || e.message));
  } finally {
    isLoading.value = false;
  }
};

// --- Reboot Router ---
const rebootRouter = async () => {
  if (!confirm("Bạn có chắc muốn reboot router?")) return;
  isLoading.value = true;
  try {
    await api.post(`/api/ships/${route.params.id}/router/reboot`);
    alert("✅ Router đang reboot...");
  } catch (e) {
    console.error(e);
    alert("❌ Reboot thất bại: " + (e?.response?.data?.error || e.message));
  } finally {
    isLoading.value = false;
  }
};

// --- Mock Performance simulation ---
const performance = reactive({
  bandwidth: 120,
  latency: 45,
  packet_loss: 0.2,
  snr: 10
});

const simulatePerformance = () => {
  performance.bandwidth = Math.max(20, performance.bandwidth + (Math.random() * 20 - 10));
  performance.latency = Math.max(10, performance.latency + (Math.random() * 10 - 5));
  performance.packet_loss = Math.max(0, performance.packet_loss + (Math.random() * 1 - 0.5));
  performance.snr = Math.max(1, performance.snr + (Math.random() * 2 - 1));
};

onMounted(async () => {
  // Fetch ship list if needed
  if (!fleetStore.ships.length) await fleetStore.fetchShips();

  await fetchRouterStats();

  // Start realtime performance simulation
  chartUpdateInterval = setInterval(() => {
    simulatePerformance();
  }, 2000);

  // If your chart needs a nextTick
  await nextTick();
});

onUnmounted(() => {
  if (chartUpdateInterval) clearInterval(chartUpdateInterval);
});
</script>

<template>
  <div class="container-fluid p-0 fade-in-tab">
    <div v-if="!ship" class="text-white-50 text-center py-5">
      Loading ship detail...
    </div>

    <div v-else>
      <!-- HEADER -->
      <div class="d-flex justify-content-between align-items-end mb-4">
        <div>
          <small class="text-info fw-bold tracking-wide">VESSEL DETAIL</small>
          <h3 class="fw-bold text-white m-0">{{ ship.name }}</h3>
          <small class="text-white-50">IMO: {{ ship.imo_number }} | Type: {{ ship.type }}</small>
        </div>

        <div class="d-flex gap-2 align-items-center">
          <button class="btn btn-outline-light btn-sm fw-bold" @click="router.back()">
            <i class="fa-solid fa-arrow-left me-2"></i>Back
          </button>

          <button class="btn btn-outline-info btn-sm fw-bold" @click="downloadPdf">
            <i class="fa-solid fa-file-pdf me-2"></i>Report
          </button>

          <button class="btn btn-outline-light btn-sm fw-bold" @click="syncCrewToRouter" :disabled="isLoading">
            <i class="fa-solid fa-rotate me-2"></i>Sync Crew
          </button>

          <button class="btn btn-danger btn-sm fw-bold" @click="rebootRouter" :disabled="isLoading">
            <i class="fa-solid fa-power-off me-2"></i>Reboot
          </button>
        </div>
      </div>

      <!-- BODY -->
      <div class="row g-4">
        <!-- Left: Map -->
        <div class="col-lg-8">
          <div class="card glass-panel p-4 h-100">
            <h6 class="fw-bold text-white mb-3">
              Live Tracking
            </h6>
            <BaseMap :ship="ship" style="height: 380px;" />
          </div>
        </div>

        <!-- Right: Router Health -->
        <div class="col-lg-4">
          <div class="card glass-panel p-4 h-100">
            <div class="d-flex justify-content-between align-items-center mb-3">
              <h6 class="fw-bold text-white m-0">Router Health</h6>
              <button class="btn btn-sm btn-outline-light" @click="fetchRouterStats">
                <i class="fa-solid fa-rotate"></i>
              </button>
            </div>

            <div v-if="!routerHealth" class="text-white-50">
              No stats available
            </div>

            <div v-else class="text-white-50 small">
              <div class="mb-2">
                <b class="text-white">CPU:</b> {{ routerHealth.cpu_usage }}%
              </div>
              <div class="mb-2">
                <b class="text-white">RAM:</b> {{ routerHealth.memory_usage }}%
              </div>
              <div class="mb-2">
                <b class="text-white">Uptime:</b> {{ routerHealth.uptime }}
              </div>
              <div class="mb-2">
                <b class="text-white">Status:</b>
                <span class="badge bg-success bg-opacity-25 text-success ms-2">
                  {{ routerHealth.status }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- PERFORMANCE -->
      <div class="card glass-panel p-4 mt-4">
        <div class="d-flex justify-content-between align-items-center mb-3">
          <h6 class="fw-bold text-white m-0">Performance Metrics</h6>
          <button class="btn btn-sm btn-outline-secondary" @click="showCrewModal = true">
            <i class="fa-solid fa-users me-2"></i>Manage Crew
          </button>
        </div>

        <VsatChart :performance="performance" />
      </div>

      <!-- Crew Modal -->
      <CrewManagementModal v-if="showCrewModal" :shipId="route.params.id" @close="showCrewModal = false" />
    </div>
  </div>
</template>

<style scoped>
.tracking-wide { letter-spacing: 1px; }
.fade-in-tab { animation: fadeIn 0.4s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>
