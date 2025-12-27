<script setup>
import { onMounted, ref, onUnmounted, watch } from 'vue';
import Chart from 'chart.js/auto';
import api from '@/services/api';

const trafficCanvas = ref(null);
const radarCanvas = ref(null);
const spark1 = ref(null);
const spark2 = ref(null);
const spark3 = ref(null);
const spark4 = ref(null);

const range = ref('24h');
const isLoading = ref(false);

// ✅ Overview state (để template không bị undefined)
const overview = ref({
  total_consumption_tb: 0,
  total_consumption_change_pct: 0,
  estimated_cost_usd: 0,
  cost_projection_usd: 0,
  security_threats: 0,
  fleet_uptime_pct: 0,
});

let trafficChart = null;
let radarChart = null;
let sparkCharts = [];
let updateInterval = null;

const topConsumers = ref([]);
const appUsage = ref([]);

// ✅ Helpers (để template không báo lỗi)
const formatMoney = (n) => {
  const num = Number(n || 0);
  return num.toLocaleString('en-US', { maximumFractionDigits: 0 });
};

const formatSignedPercent = (n) => {
  const num = Number(n || 0);
  const sign = num >= 0 ? '+' : '';
  return `${sign}${num.toFixed(1)}%`;
};

const getThreatLabel = () => {
  const t = Number(overview.value?.security_threats || 0);
  if (t >= 20) return 'Critical';
  if (t >= 10) return 'High Risk';
  if (t >= 5) return 'Medium';
  return 'Low';
};

const getUptimeLabel = () => {
  const u = Number(overview.value?.fleet_uptime_pct || 0);
  if (u >= 99.95) return 'SLA Met';
  if (u >= 99.0) return 'Stable';
  return 'Investigate';
};

// -------------------------
// Charts
// -------------------------
const destroyCharts = () => {
  if (trafficChart) trafficChart.destroy();
  if (radarChart) radarChart.destroy();
  sparkCharts.forEach((c) => c.destroy());
  sparkCharts = [];
  trafficChart = null;
  radarChart = null;
};

const renderTrafficChart = (data) => {
  if (!trafficCanvas.value) return;
  if (trafficChart) trafficChart.destroy();

  // ✅ hỗ trợ 2 kiểu response:
  // 1) { labels:[], bandwidth:[], latency:[] }
  // 2) [ {time:'', bandwidth_mbps:.., latency_ms:..}, ...]
  let labels = [];
  let bandwidth = [];
  let latency = [];

  if (Array.isArray(data)) {
    labels = data.map((x) => x.time);
    bandwidth = data.map((x) => x.bandwidth_mbps);
    latency = data.map((x) => x.latency_ms);
  } else {
    labels = data?.labels || [];
    bandwidth = data?.bandwidth || [];
    latency = data?.latency || [];
  }

  trafficChart = new Chart(trafficCanvas.value, {
    type: 'line',
    data: {
      labels,
      datasets: [
        {
          label: 'Bandwidth (Mbps)',
          data: bandwidth,
          borderWidth: 3,
          tension: 0.4,
          fill: false,
          pointRadius: 0,
        },
        {
          label: 'Latency (ms)',
          type: 'bar',
          data: latency,
          borderRadius: 4,
          barThickness: 10,
        },
      ],
    },
    options: { responsive: true, maintainAspectRatio: false },
  });
};

const renderRadarChart = (apps) => {
  if (!radarCanvas.value) return;
  if (radarChart) radarChart.destroy();

  // ✅ backend có thể trả:
  // [{app:'VoIP', current:85, limit:60}]
  const labels = (apps || []).map((x) => x.app);
  const current = (apps || []).map((x) => x.current);
  const limit = (apps || []).map((x) => x.limit);

  radarChart = new Chart(radarCanvas.value, {
    type: 'radar',
    data: {
      labels,
      datasets: [
        { label: 'Current Usage', data: current, borderWidth: 2 },
        { label: 'Limit', data: limit, borderWidth: 1, borderDash: [5, 5] },
      ],
    },
    options: { responsive: true, maintainAspectRatio: false },
  });
};

// ✅ Sparkline (optional – nếu bạn muốn hiển thị)
const renderSparks = () => {
  // Hiện tại bạn chưa dùng spark chart bằng data thật
  // Nếu muốn mình sẽ build y chang UI cũ bạn gửi (spark1-4).
};

// -------------------------
// Fetch
// -------------------------
const fetchAll = async () => {
  isLoading.value = true;
  try {
    const [overviewRes, trafficRes, topRes, appRes] = await Promise.all([
      api.get(`/api/analytics/overview?range=${range.value}`),
      api.get(`/api/analytics/traffic?range=${range.value}`),
      api.get(`/api/analytics/top-consumers?limit=10&range=${range.value}`),
      api.get(`/api/analytics/app-usage?range=${range.value}`),
    ]);

    // ✅ áp dữ liệu overview vào state (có fallback)
    const o = overviewRes.data || {};
    overview.value = {
      total_consumption_tb: Number(o.total_consumption_tb ?? o.total_consumption ?? 0),
      total_consumption_change_pct: Number(o.total_consumption_change_pct ?? o.total_consumption_change ?? 0),
      estimated_cost_usd: Number(o.estimated_cost_usd ?? o.estimated_cost ?? 0),
      cost_projection_usd: Number(o.cost_projection_usd ?? o.cost_projection ?? 0),
      security_threats: Number(o.security_threats ?? o.cyber_threats ?? 0),
      fleet_uptime_pct: Number(o.fleet_uptime_pct ?? o.fleet_uptime ?? 0),
    };

    topConsumers.value = Array.isArray(topRes.data) ? topRes.data : [];
    appUsage.value = Array.isArray(appRes.data) ? appRes.data : [];

    renderTrafficChart(trafficRes.data);
    renderRadarChart(appUsage.value);
    renderSparks();
  } catch (e) {
    console.error('Analytics fetchAll error:', e);
  } finally {
    isLoading.value = false;
  }
};

// -------------------------
// Lifecycle
// -------------------------
onMounted(async () => {
  await fetchAll();

  // Optional realtime refresh traffic chart (dev demo)
  updateInterval = setInterval(async () => {
    try {
      const trafficRes = await api.get(`/api/analytics/traffic?range=${range.value}`);
      renderTrafficChart(trafficRes.data);
    } catch (e) {}
  }, 6000);
});

onUnmounted(() => {
  if (updateInterval) clearInterval(updateInterval);
  destroyCharts();
});

watch(range, async () => {
  destroyCharts();
  await fetchAll();
});
</script>

<template>
  <div class="container-fluid p-0 fade-in-static">
    <!-- TITLE HEADER -->
    <div class="d-flex justify-content-between align-items-end mb-4">
      <div>
        <small class="text-info fw-bold tracking-wide">INTELLIGENCE MODULE</small>
        <h3 class="fw-bold text-white m-0">Advanced Analytics</h3>
      </div>
      <div class="d-flex gap-2">
        <!-- ✅ FIX: dùng v-model="range" thay vì selectedRange -->
        <select
          class="form-select bg-dark text-white border-secondary form-select-sm"
          style="width: auto;"
          v-model="range"
        >
          <option value="24h">Last 24 Hours</option>
          <option value="7d">Last 7 Days</option>
          <option value="30d">This Month</option>
        </select>

        <button class="btn btn-primary btn-sm fw-bold" :disabled="isLoading">
          <i class="fa-solid fa-download me-2"></i>Export
        </button>
      </div>
    </div>

    <!-- 1. KPI CARDS VỚI SPARKLINES -->
    <div class="row g-4 mb-4">
      <div class="col-md-3">
        <div class="card glass-panel h-100 p-3 position-relative overflow-hidden">
          <div class="d-flex justify-content-between">
            <div>
              <small class="text-white-50 text-uppercase fw-bold" style="font-size: 0.65rem;">
                Total Consumption
              </small>
              <h2 class="fw-bold text-white mt-1 mb-0">
                {{ Number(overview.total_consumption_tb).toFixed(1) }} TB
              </h2>
              <small
                :class="overview.total_consumption_change_pct >= 0 ? 'text-success' : 'text-danger'"
              >
                <i
                  class="fa-solid"
                  :class="overview.total_consumption_change_pct >= 0 ? 'fa-arrow-trend-up' : 'fa-arrow-trend-down'"
                ></i>
                {{ formatSignedPercent(overview.total_consumption_change_pct) }}
              </small>
            </div>
            <div style="width: 80px; height: 50px;">
              <canvas ref="spark1"></canvas>
            </div>
          </div>
        </div>
      </div>

      <div class="col-md-3">
        <div class="card glass-panel h-100 p-3">
          <div class="d-flex justify-content-between">
            <div>
              <small class="text-white-50 text-uppercase fw-bold" style="font-size: 0.65rem;">
                Est. Cost
              </small>
              <h2 class="fw-bold text-white mt-1 mb-0">
                ${{ formatMoney(overview.estimated_cost_usd) }}
              </h2>
              <small class="text-warning">
                Proj: ${{ formatMoney(overview.cost_projection_usd) }}
              </small>
            </div>
            <div style="width: 80px; height: 50px;">
              <canvas ref="spark2"></canvas>
            </div>
          </div>
        </div>
      </div>

      <div class="col-md-3">
        <div class="card glass-panel h-100 p-3">
          <div class="d-flex justify-content-between">
            <div>
              <small class="text-white-50 text-uppercase fw-bold" style="font-size: 0.65rem;">
                Cyber Threats
              </small>
              <h2 class="fw-bold text-white mt-1 mb-0">{{ overview.security_threats }}</h2>
              <small class="text-danger">{{ getThreatLabel() }}</small>
            </div>
            <div style="width: 80px; height: 50px;">
              <canvas ref="spark3"></canvas>
            </div>
          </div>
        </div>
      </div>

      <div class="col-md-3">
        <div class="card glass-panel h-100 p-3">
          <div class="d-flex justify-content-between">
            <div>
              <small class="text-white-50 text-uppercase fw-bold" style="font-size: 0.65rem;">
                Fleet Uptime
              </small>
              <h2 class="fw-bold text-white mt-1 mb-0">
                {{ Number(overview.fleet_uptime_pct).toFixed(2) }}%
              </h2>
              <small class="text-success">{{ getUptimeLabel() }}</small>
            </div>
            <div style="width: 80px; height: 50px;">
              <canvas ref="spark4"></canvas>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 2. MAIN CHARTS AREA -->
    <div class="row g-4 mb-4">
      <!-- Biểu đồ kép: Traffic & Latency -->
      <div class="col-lg-8">
        <div class="card glass-panel h-100 p-4">
          <div class="d-flex justify-content-between align-items-center mb-3">
            <h6 class="fw-bold text-white m-0">Network Performance Analysis</h6>
            <div class="d-flex align-items-center gap-3">
              <small class="text-white-50">
                <i class="fa-solid fa-circle text-primary fa-2xs me-1"></i> Bandwidth
              </small>
              <small class="text-white-50">
                <i class="fa-solid fa-circle text-success fa-2xs me-1"></i> Latency
              </small>
            </div>
          </div>
          <div style="height: 350px;">
            <canvas ref="trafficCanvas"></canvas>
          </div>
        </div>
      </div>

      <!-- Biểu đồ Radar: Application -->
      <div class="col-lg-4">
        <div class="card glass-panel h-100 p-4">
          <h6 class="fw-bold text-white mb-3">Application Usage Profile</h6>
          <div style="height: 350px; position: relative;">
            <canvas ref="radarCanvas"></canvas>
            <div
              class="position-absolute top-50 start-50 translate-middle text-center"
              style="pointer-events: none;"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 3. TOP CONSUMERS TABLE -->
    <div class="card glass-panel p-0">
      <div class="card-header border-bottom border-white border-opacity-10 py-3">
        <h6 class="fw-bold text-white m-0">High Bandwidth Consumption</h6>
      </div>
      <div class="table-responsive">
        <table class="table table-hover align-middle mb-0">
          <thead>
            <tr>
              <th class="ps-4">Vessel Name</th>
              <th>Total Usage</th>
              <th>Cost Impact</th>
              <th>Signal Health</th>
              <th>Status</th>
              <th class="text-end pe-4">Trend</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in topConsumers" :key="item.name">
              <td class="ps-4">
                <div class="d-flex align-items-center">
                  <div class="p-2 bg-white bg-opacity-10 rounded me-3">
                    <i class="fa-solid fa-ship text-info"></i>
                  </div>
                  <div>
                    <div class="fw-bold text-white">{{ item.name }}</div>
                    <small class="text-white-50">{{ item.type }}</small>
                  </div>
                </div>
              </td>
              <td>
                <div class="fw-bold text-white">{{ item.data }} GB</div>
                <div class="progress bg-dark mt-1" style="height: 4px; width: 80px;">
                  <div class="progress-bar bg-primary" :style="{ width: (item.data / 500) * 100 + '%' }"></div>
                </div>
              </td>
              <td class="font-monospace text-warning">${{ item.cost }}</td>
              <td>
                <span :class="item.signal > 90 ? 'text-success' : 'text-warning'">
                  <i class="fa-solid fa-wifi me-1"></i> {{ item.signal }}%
                </span>
              </td>
              <td>
                <span
                  class="badge border"
                  :class="
                    item.status === 'High'
                      ? 'text-danger border-danger bg-danger bg-opacity-10'
                      : item.status === 'Warning'
                        ? 'text-warning border-warning bg-warning bg-opacity-10'
                        : 'text-success border-success bg-success bg-opacity-10'
                  "
                >
                  {{ item.status }}
                </span>
              </td>
              <td class="text-end pe-4 text-white-50">
                <i class="fa-solid fa-chart-line"></i> Analyze
              </td>
            </tr>

            <tr v-if="topConsumers.length === 0 && !isLoading">
              <td colspan="6" class="text-center text-white-50 py-4">No data</td>
            </tr>
            <tr v-if="isLoading">
              <td colspan="6" class="text-center text-white-50 py-4">Loading...</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tracking-wide { letter-spacing: 1px; }
.fade-in-static { animation: simpleFade 0.6s ease-out; }
@keyframes simpleFade { from { opacity: 0; } to { opacity: 1; } }
</style>
