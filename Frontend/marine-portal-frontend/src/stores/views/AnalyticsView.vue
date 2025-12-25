<script setup>
import { onMounted, ref } from 'vue';
import Chart from 'chart.js/auto';

// Dữ liệu giả lập cho báo cáo
const topConsumers = [
    { name: 'VIMC DIAMOND', usage: '450 GB', cost: '$2,100', trend: '+12%' },
    { name: 'EVER GIVEN', usage: '380 GB', cost: '$1,850', trend: '+5%' },
    { name: 'PACIFIC STAR', usage: '210 GB', cost: '$980', trend: '-2%' },
];

const pieCanvas = ref(null);
const lineCanvas = ref(null);

onMounted(() => {
    // 1. Biểu đồ tròn (Application Usage)
    new Chart(pieCanvas.value, {
        type: 'doughnut',
        data: {
            labels: ['Business (ERP)', 'Crew Wi-Fi', 'VoIP', 'Video Streaming'],
            datasets: [{
                data: [45, 30, 15, 10],
                backgroundColor: ['#2563eb', '#10b981', '#f59e0b', '#ef4444'],
                borderWidth: 0
            }]
        },
        options: { cutout: '75%', plugins: { legend: { position: 'right' } } }
    });

    // 2. Biểu đồ đường (Traffic Trend)
    new Chart(lineCanvas.value, {
        type: 'line',
        data: {
            labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
            datasets: [{
                label: 'Total Traffic (GB)',
                data: [120, 190, 150, 250, 220, 300, 350],
                borderColor: '#2563eb',
                backgroundColor: 'rgba(37, 99, 235, 0.1)',
                fill: true,
                tension: 0.4
            }]
        },
        options: { 
            maintainAspectRatio: false, 
            plugins: { legend: { display: false } },
            scales: { x: { grid: { display: false } }, y: { border: { display: false } } }
        }
    });
});
</script>

<template>
    <div class="container-fluid p-0">
        <h3 class="fw-bold text-dark mb-4">Fleet Intelligence Report</h3>

        <!-- KPI CARDS -->
        <div class="row g-4 mb-4">
            <div class="col-md-3">
                <div class="card border-0 shadow-sm p-3">
                    <small class="text-secondary text-uppercase fw-bold">Total Data (Month)</small>
                    <div class="d-flex justify-content-between align-items-center mt-2">
                        <h2 class="fw-bold m-0">15.4 TB</h2>
                        <span class="badge bg-success-subtle text-success"><i class="fa-solid fa-arrow-trend-up"></i> 12%</span>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card border-0 shadow-sm p-3">
                    <small class="text-secondary text-uppercase fw-bold">Est. Cost</small>
                    <div class="d-flex justify-content-between align-items-center mt-2">
                        <h2 class="fw-bold m-0 text-dark">$12,450</h2>
                        <span class="badge bg-warning-subtle text-warning">Pending</span>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card border-0 shadow-sm p-3">
                    <small class="text-secondary text-uppercase fw-bold">Cyber Threats</small>
                    <div class="d-flex justify-content-between align-items-center mt-2">
                        <h2 class="fw-bold m-0 text-danger">14</h2>
                        <span class="badge bg-danger-subtle text-danger">Blocked</span>
                    </div>
                </div>
            </div>
             <div class="col-md-3">
                <div class="card border-0 shadow-sm p-3 bg-primary text-white">
                    <small class="text-white-50 text-uppercase fw-bold">Fleet Uptime</small>
                    <div class="d-flex justify-content-between align-items-center mt-2">
                        <h2 class="fw-bold m-0">99.98%</h2>
                        <i class="fa-solid fa-award fa-2x opacity-50"></i>
                    </div>
                </div>
            </div>
        </div>

        <!-- CHARTS SECTION -->
        <div class="row g-4">
            <div class="col-lg-8">
                <div class="card border-0 shadow-sm h-100">
                    <div class="card-header bg-white py-3 fw-bold">Bandwidth Usage Trends</div>
                    <div class="card-body">
                        <div style="height: 300px;">
                            <canvas ref="lineCanvas"></canvas>
                        </div>
                    </div>
                </div>
            </div>
            <div class="col-lg-4">
                <div class="card border-0 shadow-sm h-100">
                    <div class="card-header bg-white py-3 fw-bold">Traffic by Application</div>
                    <div class="card-body d-flex align-items-center">
                        <div style="height: 250px; width: 100%;">
                            <canvas ref="pieCanvas"></canvas>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- TOP CONSUMERS TABLE -->
        <div class="card border-0 shadow-sm mt-4">
            <div class="card-header bg-white py-3 fw-bold">Top Data Consumers</div>
            <div class="table-responsive">
                <table class="table align-middle mb-0">
                    <thead class="bg-light">
                        <tr>
                            <th class="ps-4">Vessel Name</th>
                            <th>Data Usage</th>
                            <th>Est. Cost</th>
                            <th>Trend</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in topConsumers" :key="item.name">
                            <td class="ps-4 fw-bold text-primary">{{ item.name }}</td>
                            <td class="fw-bold">{{ item.usage }}</td>
                            <td>{{ item.cost }}</td>
                            <td><span :class="item.trend.includes('+') ? 'text-danger' : 'text-success'">{{ item.trend }}</span></td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>