<script setup>
import { onMounted, ref, onUnmounted, reactive, watch } from 'vue';
import Chart from 'chart.js/auto';
import apiClient from '@/services/apiClient';

// --- REF & STATE ---
const trafficCanvas = ref(null);
const radarCanvas = ref(null);
const spark1 = ref(null);
const spark2 = ref(null);
const spark3 = ref(null);
const spark4 = ref(null);

let trafficChart = null;
let radarChart = null;
let sparkChart1 = null;
let sparkChart2 = null;
let sparkChart3 = null;
let sparkChart4 = null;

// --- RANGE ---
const selectedRange = ref('24h'); // 24h | 7d | 30d
const isLoading = ref(false);

// --- DATA FROM API ---
const overview = reactive({
  total_consumption_tb: 0,
  total_consumption_change_pct: 0,
  estimated_cost_usd: 0,
  cost_projection_usd: 0,
  security_threats: 0,
  threats_change_pct: 0,
  fleet_uptime_pct: 0,
  uptime_change_pct: 0,
});

const topConsumers = ref([]);

// --- HÀM TẠO GRADIENT MÀU (giữ nguyên) ---
const createGradient = (ctx, colorStart, colorEnd) => {
  const gradient = ctx.createLinearGradient(0, 0, 0, 400);
  gradient.addColorStop(0, colorStart);
  gradient.addColorStop(1, colorEnd);
  return gradient;
};

// --- Helpers ---
const safeNum = (v, fallback = 0) => {
  const n = Number(v);
  return Number.isFinite(n) ? n : fallback;
};

const formatMoney = (v) => {
  const n = safeNum(v, 0);
  return new Intl.NumberFormat('en-US').format(n);
};

const formatSignedPercent = (v) => {
  const n = safeNum(v, 0);
  const sign = n > 0 ? '+' : '';
  return `${sign}${n.toFixed(1)}%`;
};

const getThreatLabel = () => {
  // Bạn có thể đổi rule theo thực tế
  if (overview.security_threats >= 20) return 'High Risk';
  if (overview.security_threats >= 10) return 'Medium Risk';
  return 'Low Risk';
};

const getUptimeLabel = () => {
  if (overview.fleet_uptime_pct >= 99.5) return 'SLA Met';
  if (overview.fleet_uptime_pct >= 98.0) return 'Stable';
  return 'Unstable';
};

// --- API Calls ---
const fetchOverview = async () => {
  const res = await apiClient.get('/api/analytics/overview', { params: { range: selectedRange.value } });
  const data = res.data || {};

  overview.total_consumption_tb = safeNum(data.total_consumption_tb);
  overview.total_consumption_change_pct = safeNum(data.total_consumption_change_pct);

  overview.estimated_cost_usd = safeNum(data.estimated_cost_usd);
  overview.cost_projection_usd = safeNum(data.cost_projection_usd);

  overview.security_threats = safeNum(data.security_threats);
  overview.threats_change_pct = safeNum(data.threats_change_pct);

  overview.fleet_uptime_pct = safeNum(data.fleet_uptime_pct);
  overview.uptime_change_pct = safeNum(data.uptime_change_pct);
};

const fetchTraffic = async () => {
  const res = await apiClient.get('/api/analytics/traffic', { params: { range: selectedRange.value } });
  const data = res.data || {};

  const labels = Array.isArray(data.labels) ? data.labels : [];
  const bandwidth = Array.isArray(data.bandwidth_mbps) ? data.bandwidth_mbps.map((x) => safeNum(x)) : [];
  const latency = Array.isArray(data.latency_ms) ? data.latency_ms.map((x) => safeNum(x)) : [];

  if (trafficChart) {
    trafficChart.data.labels = labels;
    trafficChart.data.datasets[0].data = bandwidth;
    trafficChart.data.datasets[1].data = latency;
    trafficChart.update();
  }
};

const fetchAppUsage = async () => {
  const res = await apiClient.get('/api/analytics/app-usage', { params: { range: selectedRange.value } });
  const data = res.data || {};

  const labels = Array.isArray(data.labels) ? data.labels : [];
  const current = Array.isArray(data.current) ? data.current.map((x) => safeNum(x)) : [];
  const limit = Array.isArray(data.limit) ? data.limit.map((x) => safeNum(x)) : [];

  if (radarChart) {
    radarChart.data.labels = labels;
    radarChart.data.datasets[0].data = current;
    radarChart.data.datasets[1].data = limit;
    radarChart.update();
  }
};

const fetchTopConsumers = async () => {
  const res = await apiClient.get('/api/analytics/top-consumers', {
    params: { range: selectedRange.value, limit: 10 },
  });

  const data = res.data;
  topConsumers.value = Array.isArray(data) ? data.map((x) => ({
    // Map backend -> UI fields (giữ nguyên template)
    name: x.name ?? 'Unknown',
    type: x.type ?? 'N/A',
    data: safeNum(x.data_gb ?? x.data, 0),   // backend dùng data_gb
    cost: safeNum(x.cost_usd ?? x.cost, 0),  // backend dùng cost_usd
    signal: safeNum(x.signal, 0),
    status: x.status ?? 'Normal',
  })) : [];
};

const loadAll = async () => {
  isLoading.value = true;
  try {
    await Promise.all([
      fetchOverview(),
      fetchTraffic(),
      fetchAppUsage(),
      fetchTopConsumers(),
    ]);

    // Update sparklines sau khi overview có data
    updateSparklines();
  } catch (e) {
    console.error('Analytics load failed:', e);
  } finally {
    isLoading.value = false;
  }
};

// --- Chart creators ---
const createSpark = (canvasEl, color, data) => {
  if (!canvasEl) return null;
  return new Chart(canvasEl, {
    type: 'line',
    data: {
      labels: [1,2,3,4,5,6,7],
      datasets: [{ data, borderColor: color, borderWidth: 2, pointRadius: 0, tension: 0.4 }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false }, tooltip: { enabled: false } },
      scales: { x: { display: false }, y: { display: false } }
    }
  });
};

const destroyCharts = () => {
  if (trafficChart) trafficChart.destroy();
  if (radarChart) radarChart.destroy();
  if (sparkChart1) sparkChart1.destroy();
  if (sparkChart2) sparkChart2.destroy();
  if (sparkChart3) sparkChart3.destroy();
  if (sparkChart4) sparkChart4.destroy();
  trafficChart = null;
  radarChart = null;
  sparkChart1 = null;
  sparkChart2 = null;
  sparkChart3 = null;
  sparkChart4 = null;
};

// Build sparklines data from KPIs (nhẹ nhàng, đẹp, không random quá)
const buildSparkFromValue = (base, variance = 10) => {
  const b = safeNum(base, 0);
  // Tạo 7 điểm dựa trên base
  return Array.from({ length: 7 }, (_, i) => {
    const t = i - 3;
    return Math.max(0, Math.round(b + t * (variance * 0.5) + (Math.random() - 0.5) * variance));
  });
};

const updateSparklines = () => {
  // Nếu chưa init spark charts thì init
  if (!sparkChart1 && spark1.value) {
    sparkChart1 = createSpark(spark1.value, '#3b82f6', buildSparkFromValue(overview.total_consumption_tb * 10, 10));
  } else if (sparkChart1) {
    sparkChart1.data.datasets[0].data = buildSparkFromValue(overview.total_consumption_tb * 10, 10);
    sparkChart1.update();
  }

  if (!sparkChart2 && spark2.value) {
    sparkChart2 = createSpark(spark2.value, '#f59e0b', buildSparkFromValue(overview.estimated_cost_usd / 200, 15));
  } else if (sparkChart2) {
    sparkChart2.data.datasets[0].data = buildSparkFromValue(overview.estimated_cost_usd / 200, 15);
    sparkChart2.update();
  }

  if (!sparkChart3 && spark3.value) {
    sparkChart3 = createSpark(spark3.value, '#ef4444', buildSparkFromValue(overview.security_threats, 6));
  } else if (sparkChart3) {
    sparkChart3.data.datasets[0].data = buildSparkFromValue(overview.security_threats, 6);
    sparkChart3.update();
  }

  if (!sparkChart4 && spark4.value) {
    sparkChart4 = createSpark(spark4.value, '#10b981', buildSparkFromValue(overview.fleet_uptime_pct, 1));
  } else if (sparkChart4) {
    sparkChart4.data.datasets[0].data = buildSparkFromValue(overview.fleet_uptime_pct, 1);
    sparkChart4.update();
  }
};

// --- Init charts (giữ style đẹp của bạn) ---
const initCharts = () => {
  // 1) Main chart: Mixed Traffic vs Latency
  if (trafficCanvas.value) {
    const ctx = trafficCanvas.value.getContext('2d');
    const gradientFill = createGradient(ctx, 'rgba(59, 130, 246, 0.5)', 'rgba(59, 130, 246, 0.0)');
    const gradientBar = createGradient(ctx, '#10b981', '#059669');

    trafficChart = new Chart(ctx, {
      type: 'line',
      data: {
        labels: [],
        datasets: [
          {
            label: 'Bandwidth (Mbps)',
            data: [],
            borderColor: '#3b82f6',
            backgroundColor: gradientFill,
            borderWidth: 3,
            tension: 0.4,
            fill: true,
            yAxisID: 'y',
            pointRadius: 0,
            pointHoverRadius: 6
          },
          {
            label: 'Latency (ms)',
            type: 'bar',
            data: [],
            backgroundColor: gradientBar,
            borderRadius: 4,
            barThickness: 10,
            yAxisID: 'y1'
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: { mode: 'index', intersect: false },
        plugins: { legend: { labels: { color: '#94a3b8', font: { family: 'JetBrains Mono' } } } },
        scales: {
          x: { grid: { color: 'rgba(255,255,255,0.05)' }, ticks: { color: '#94a3b8' } },
          y: {
            type: 'linear', display: true, position: 'left',
            grid: { color: 'rgba(255,255,255,0.05)' }, ticks: { color: '#3b82f6' },
            title: { display: true, text: 'Throughput (Mbps)', color: '#3b82f6' }
          },
          y1: {
            type: 'linear', display: true, position: 'right',
            grid: { drawOnChartArea: false }, ticks: { color: '#10b981' },
            title: { display: true, text: 'Latency (ms)', color: '#10b981' }
          }
        }
      }
    });
  }

  // 2) Radar chart: App usage
  if (radarCanvas.value) {
    radarChart = new Chart(radarCanvas.value, {
      type: 'radar',
      data: {
        labels: [],
        datasets: [
          {
            label: 'Current Usage',
            data: [],
            backgroundColor: 'rgba(245, 158, 11, 0.2)',
            borderColor: '#f59e0b',
            pointBackgroundColor: '#f59e0b',
            borderWidth: 2,
          },
          {
            label: 'Limit',
            data: [],
            backgroundColor: 'rgba(255, 255, 255, 0.05)',
            borderColor: '#64748b',
            borderWidth: 1,
            borderDash: [5, 5]
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          r: {
            angleLines: { color: 'rgba(255, 255, 255, 0.1)' },
            grid: { color: 'rgba(255, 255, 255, 0.1)' },
            pointLabels: { color: '#cbd5e1', font: { size: 11, family: 'Inter' } },
            ticks: { display: false },
            suggestedMin: 0,
            suggestedMax: 100,
          }
        },
        plugins: { legend: { display: false } }
      }
    });
  }

  // 3) Sparklines init (data sẽ update sau khi load overview)
  if (spark1.value) sparkChart1 = createSpark(spark1.value, '#3b82f6', [10, 25, 20, 40, 35, 50, 60]);
  if (spark2.value) sparkChart2 = createSpark(spark2.value, '#f59e0b', [50, 40, 60, 55, 70, 65, 80]);
  if (spark3.value) sparkChart3 = createSpark(spark3.value, '#ef4444', [5, 2, 8, 1, 0, 10, 14]);
  if (spark4.value) sparkChart4 = createSpark(spark4.value, '#10b981', [98, 99, 98, 99, 99, 100, 100]);
};

// --- Lifecycle ---
onMounted(async () => {
  initCharts();
  await loadAll();
});

watch(selectedRange, async () => {
  await loadAll();
});

onUnmounted(() => {
  destroyCharts();
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
                <!-- ✅ BIND RANGE -->
                <select class="form-select bg-dark text-white border-secondary form-select-sm" style="width: auto;" v-model="selectedRange">
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
                            <small class="text-white-50 text-uppercase fw-bold" style="font-size: 0.65rem;">Total Consumption</small>
                            <h2 class="fw-bold text-white mt-1 mb-0">{{ overview.total_consumption_tb.toFixed(1) }} TB</h2>
                            <small :class="overview.total_consumption_change_pct >= 0 ? 'text-success' : 'text-danger'">
                              <i class="fa-solid" :class="overview.total_consumption_change_pct >= 0 ? 'fa-arrow-trend-up' : 'fa-arrow-trend-down'"></i>
                              {{ formatSignedPercent(overview.total_consumption_change_pct) }}
                            </small>
                        </div>
                        <div style="width: 80px; height: 50px;"><canvas ref="spark1"></canvas></div>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card glass-panel h-100 p-3">
                    <div class="d-flex justify-content-between">
                        <div>
                            <small class="text-white-50 text-uppercase fw-bold" style="font-size: 0.65rem;">Est. Cost</small>
                            <h2 class="fw-bold text-white mt-1 mb-0">${{ formatMoney(overview.estimated_cost_usd) }}</h2>
                            <small class="text-warning">Proj: ${{ formatMoney(overview.cost_projection_usd) }}</small>
                        </div>
                        <div style="width: 80px; height: 50px;"><canvas ref="spark2"></canvas></div>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card glass-panel h-100 p-3">
                    <div class="d-flex justify-content-between">
                        <div>
                            <small class="text-white-50 text-uppercase fw-bold" style="font-size: 0.65rem;">Cyber Threats</small>
                            <h2 class="fw-bold text-white mt-1 mb-0">{{ overview.security_threats }}</h2>
                            <small class="text-danger">{{ getThreatLabel() }}</small>
                        </div>
                        <div style="width: 80px; height: 50px;"><canvas ref="spark3"></canvas></div>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card glass-panel h-100 p-3">
                    <div class="d-flex justify-content-between">
                        <div>
                            <small class="text-white-50 text-uppercase fw-bold" style="font-size: 0.65rem;">Fleet Uptime</small>
                            <h2 class="fw-bold text-white mt-1 mb-0">{{ overview.fleet_uptime_pct.toFixed(2) }}%</h2>
                            <small class="text-success">{{ getUptimeLabel() }}</small>
                        </div>
                        <div style="width: 80px; height: 50px;"><canvas ref="spark4"></canvas></div>
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
                            <small class="text-white-50"><i class="fa-solid fa-circle text-primary fa-2xs me-1"></i> Bandwidth</small>
                            <small class="text-white-50"><i class="fa-solid fa-circle text-success fa-2xs me-1"></i> Latency</small>
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
                        <div class="position-absolute top-50 start-50 translate-middle text-center" style="pointer-events: none;">
                            <!-- Center Info -->
                        </div>
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
                                    <div class="progress-bar bg-primary" :style="{width: (item.data/500)*100 + '%'}"></div>
                                </div>
                            </td>
                            <td class="font-monospace text-warning">${{ item.cost }}</td>
                            <td>
                                <span :class="item.signal > 90 ? 'text-success' : 'text-warning'">
                                    <i class="fa-solid fa-wifi me-1"></i> {{ item.signal }}%
                                </span>
                            </td>
                            <td>
                                <span class="badge border"
                                    :class="item.status === 'High'
                                      ? 'text-danger border-danger bg-danger bg-opacity-10'
                                      : (item.status === 'Warning'
                                        ? 'text-warning border-warning bg-warning bg-opacity-10'
                                        : 'text-success border-success bg-success bg-opacity-10')">
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
