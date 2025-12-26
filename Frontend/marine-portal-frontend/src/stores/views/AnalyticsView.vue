<script setup>
import { onMounted, ref, onUnmounted } from 'vue';
import Chart from 'chart.js/auto';

// --- REF & STATE ---
const trafficCanvas = ref(null);
const radarCanvas = ref(null);
const spark1 = ref(null);
const spark2 = ref(null);
const spark3 = ref(null);
const spark4 = ref(null);

let trafficChart = null;
let updateInterval = null;

// --- DỮ LIỆU GIẢ LẬP CAO CẤP ---
const topConsumers = [
    { name: 'VIMC DIAMOND', type: 'Cargo', data: 450, cost: 2100, signal: 98, status: 'High' },
    { name: 'PV DRILLING V', type: 'Drill', data: 380, cost: 1850, signal: 85, status: 'Normal' },
    { name: 'EVER GIVEN', type: 'Container', data: 310, cost: 1400, signal: 42, status: 'Warning' },
    { name: 'MAERSK HANOI', type: 'Container', data: 150, cost: 750, signal: 92, status: 'Normal' },
];

// --- HÀM TẠO GRADIENT MÀU (Bí mật của giao diện đẹp) ---
const createGradient = (ctx, colorStart, colorEnd) => {
    const gradient = ctx.createLinearGradient(0, 0, 0, 400);
    gradient.addColorStop(0, colorStart);
    gradient.addColorStop(1, colorEnd);
    return gradient;
};

onMounted(() => {
    // 1. BIỂU ĐỒ CHÍNH: TRAFFIC VS LATENCY (Mixed Chart)
    if (trafficCanvas.value) {
        const ctx = trafficCanvas.value.getContext('2d');
        const gradientFill = createGradient(ctx, 'rgba(59, 130, 246, 0.5)', 'rgba(59, 130, 246, 0.0)');
        const gradientBar = createGradient(ctx, '#10b981', '#059669');

        trafficChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'],
                datasets: [
                    {
                        label: 'Bandwidth (Mbps)',
                        data: [120, 190, 150, 250, 220, 300, 350],
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
                        data: [40, 50, 45, 80, 60, 90, 110],
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

    // 2. BIỂU ĐỒ RADAR (Phân tích ứng dụng)
    if (radarCanvas.value) {
        new Chart(radarCanvas.value, {
            type: 'radar',
            data: {
                labels: ['VoIP', 'Video', 'ERP', 'Email', 'Social', 'IoT'],
                datasets: [{
                    label: 'Current Usage',
                    data: [85, 90, 60, 40, 70, 30],
                    backgroundColor: 'rgba(245, 158, 11, 0.2)',
                    borderColor: '#f59e0b',
                    pointBackgroundColor: '#f59e0b',
                    borderWidth: 2,
                }, {
                    label: 'Limit',
                    data: [60, 70, 90, 90, 50, 100],
                    backgroundColor: 'rgba(255, 255, 255, 0.05)',
                    borderColor: '#64748b',
                    borderWidth: 1,
                    borderDash: [5, 5]
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    r: {
                        angleLines: { color: 'rgba(255, 255, 255, 0.1)' },
                        grid: { color: 'rgba(255, 255, 255, 0.1)' },
                        pointLabels: { color: '#cbd5e1', font: { size: 11, family: 'Inter' } },
                        ticks: { display: false }
                    }
                },
                plugins: { legend: { display: false } }
            }
        });
    }

    // 3. TẠO CÁC SPARKLINES (Biểu đồ mini trong thẻ KPI)
    const createSpark = (ref, color, data) => {
        if(!ref) return;
        new Chart(ref, {
            type: 'line',
            data: { labels: [1,2,3,4,5,6,7], datasets: [{ data: data, borderColor: color, borderWidth: 2, pointRadius: 0, tension: 0.4 }] },
            options: { responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false }, tooltip: { enabled: false } }, scales: { x: { display: false }, y: { display: false } } }
        });
    };

    createSpark(spark1.value, '#3b82f6', [10, 25, 20, 40, 35, 50, 60]); // Total Data
    createSpark(spark2.value, '#f59e0b', [50, 40, 60, 55, 70, 65, 80]); // Cost
    createSpark(spark3.value, '#ef4444', [5, 2, 8, 1, 0, 10, 14]);     // Threats
    createSpark(spark4.value, '#10b981', [98, 99, 98, 99, 99, 100, 100]); // Uptime

    // 4. GIẢ LẬP REAL-TIME (Làm biểu đồ động đậy)
    updateInterval = setInterval(() => {
        if(trafficChart) {
            const newData = Math.floor(Math.random() * 100) + 200;
            const newLatency = Math.floor(Math.random() * 50) + 30;
            
            // Xóa đầu, thêm đuôi
            trafficChart.data.datasets[0].data.shift();
            trafficChart.data.datasets[0].data.push(newData);
            trafficChart.data.datasets[1].data.shift();
            trafficChart.data.datasets[1].data.push(newLatency);
            trafficChart.update('none'); // Update mượt không animation lại
        }
    }, 2000);
});

onUnmounted(() => { if(updateInterval) clearInterval(updateInterval); });
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
                <select class="form-select bg-dark text-white border-secondary form-select-sm" style="width: auto;">
                    <option>Last 24 Hours</option>
                    <option>Last 7 Days</option>
                    <option>This Month</option>
                </select>
                <button class="btn btn-primary btn-sm fw-bold"><i class="fa-solid fa-download me-2"></i>Export</button>
            </div>
        </div>

        <!-- 1. KPI CARDS VỚI SPARKLINES (ĐIỂM NHẤN) -->
        <div class="row g-4 mb-4">
            <div class="col-md-3">
                <div class="card glass-panel h-100 p-3 position-relative overflow-hidden">
                    <div class="d-flex justify-content-between">
                        <div>
                            <small class="text-white-50 text-uppercase fw-bold" style="font-size: 0.65rem;">Total Consumption</small>
                            <h2 class="fw-bold text-white mt-1 mb-0">15.4 TB</h2>
                            <small class="text-success"><i class="fa-solid fa-arrow-trend-up"></i> +12.5%</small>
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
                            <h2 class="fw-bold text-white mt-1 mb-0">$12,450</h2>
                            <small class="text-warning">Proj: $15k</small>
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
                            <h2 class="fw-bold text-white mt-1 mb-0">14</h2>
                            <small class="text-danger">High Risk</small>
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
                            <h2 class="fw-bold text-white mt-1 mb-0">99.98%</h2>
                            <small class="text-success">SLA Met</small>
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
                        <!-- Center Info -->
                        <div class="position-absolute top-50 start-50 translate-middle text-center" style="pointer-events: none;">
                            <!-- <i class="fa-solid fa-network-wired fa-2x text-white opacity-25"></i> -->
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 3. TOP CONSUMERS TABLE (CHI TIẾT HƠN) -->
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
                                    :class="item.status === 'High' ? 'text-danger border-danger bg-danger bg-opacity-10' : 'text-success border-success bg-success bg-opacity-10'">
                                    {{ item.status }}
                                </span>
                            </td>
                            <td class="text-end pe-4 text-white-50">
                                <i class="fa-solid fa-chart-line"></i> Analyze
                            </td>
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