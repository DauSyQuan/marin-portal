<script setup>
import { onMounted, ref, watch, onUnmounted } from 'vue';
import Chart from 'chart.js/auto';

// Nhận giá trị SNR mới nhất từ cha truyền vào
const props = defineProps(['snr']);

const chartCanvas = ref(null);
let chartInstance = null;

// Khởi tạo dữ liệu mẫu ban đầu (30 điểm dữ liệu rỗng)
const initialData = Array(30).fill(0).map(() => 10 + Math.random() * 2);

onMounted(() => {
    if (chartCanvas.value) {
        const ctx = chartCanvas.value.getContext('2d');
        
        // Tạo Gradient màu xanh cho đẹp
        const gradient = ctx.createLinearGradient(0, 0, 0, 400);
        gradient.addColorStop(0, 'rgba(37, 99, 235, 0.5)'); // Xanh đậm
        gradient.addColorStop(1, 'rgba(37, 99, 235, 0.0)'); // Mờ dần

        chartInstance = new Chart(ctx, {
            type: 'line',
            data: {
                labels: Array(30).fill(''), // Nhãn rỗng để ẩn trục X
                datasets: [{
                    label: 'Signal Quality (SNR)',
                    data: initialData,
                    borderColor: '#2563eb', // Màu đường kẻ
                    borderWidth: 2,
                    backgroundColor: gradient,
                    fill: true,
                    tension: 0.4, // Đường cong mềm mại
                    pointRadius: 0 // Ẩn dấu chấm tròn cho mượt
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: false, // Tắt animation mặc định để real-time mượt hơn
                interaction: { intersect: false },
                scales: {
                    y: {
                        min: 0,
                        max: 20,
                        grid: { color: '#f1f5f9' },
                        title: { display: true, text: 'dB' }
                    },
                    x: { display: false } // Ẩn lưới trục X
                },
                plugins: { legend: { display: false } }
            }
        });
    }
});

// Khi cha truyền SNR mới vào, cập nhật biểu đồ
watch(() => props.snr, (newVal) => {
    if (chartInstance) {
        const data = chartInstance.data.datasets[0].data;
        
        // Kỹ thuật "Sliding Window": Bỏ số đầu, thêm số cuối
        data.shift(); 
        data.push(newVal);
        
        chartInstance.update(); // Vẽ lại
    }
});

onUnmounted(() => {
    if (chartInstance) chartInstance.destroy();
});
</script>

<template>
    <div style="height: 300px; width: 100%;">
        <canvas ref="chartCanvas"></canvas>
    </div>
</template>