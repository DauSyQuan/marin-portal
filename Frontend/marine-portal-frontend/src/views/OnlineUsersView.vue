<script setup>
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';

// State
const users = ref([]);
const loading = ref(false);
const searchQuery = ref('');

// Fetch Data
const fetchOnlineUsers = async () => {
    loading.value = true;
    try {
        const res = await axios.get('http://localhost:8080/api/online-users');
        users.value = res.data;
    } catch (e) {
        console.error(e);
    } finally {
        loading.value = false;
    }
};

// Filter Logic
const filteredUsers = computed(() => {
    if (!searchQuery.value) return users.value;
    return users.value.filter(u => 
        u.username.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        u.framed_ip.includes(searchQuery.value)
    );
});

// Kick Action
const kickUser = async (username) => {
    if(!confirm(`Bạn có chắc chắn muốn ngắt kết nối user: ${username}?`)) return;
    
    try {
        await axios.post(`http://localhost:8080/api/online-users/${username}/kick`);
        // Xóa tạm khỏi danh sách local để tạo cảm giác mượt
        users.value = users.value.filter(u => u.username !== username);
        alert(`Đã kick user ${username} thành công!`);
    } catch (e) {
        alert("Lỗi khi kick user: " + e.message);
    }
};

// Format Bytes (MB/GB)
const formatBytes = (mb) => {
    if (!mb) return '0 B';
    if (mb > 1024) return (mb / 1024).toFixed(2) + ' GB';
    return mb.toFixed(2) + ' MB';
};

onMounted(() => {
    fetchOnlineUsers();
});
</script>

<template>
    <div class="container-fluid p-0">
        <!-- Title & Breadcrumb placeholder if needed -->
        <!-- FILTER BAR -->
        <div class="d-flex justify-content-between align-items-center mb-4">
            <div class="d-flex gap-2">
                <div class="input-group" style="width: 300px;">
                    <span class="input-group-text bg-white border-end-0"><i class="fa-solid fa-magnifying-glass text-secondary"></i></span>
                    <input type="text" class="form-control border-start-0 ps-0" v-model="searchQuery" placeholder="Tìm theo username...">
                </div>
                <button class="btn btn-white border fw-bold text-secondary" @click="fetchOnlineUsers">Refresh</button>
            </div>
            <div class="d-flex gap-2">
                <button class="btn btn-light border rounded-circle"><i class="fa-solid fa-rotate-right"></i></button>
                <button class="btn btn-light border rounded-circle"><i class="fa-solid fa-expand"></i></button>
                <button class="btn btn-light border rounded-circle"><i class="fa-solid fa-gear"></i></button>
            </div>
        </div>

        <!-- MAIN TABLE -->
        <div class="card border-0 shadow-sm">
            <div class="table-responsive">
                <table class="table table-hover align-middle mb-0">
                    <thead class="bg-light text-secondary small text-uppercase fw-bold">
                        <tr>
                            <th class="ps-4 py-3">Username</th>
                            <th>NAS IP</th>
                            <th>Framed IP</th>
                            <th>Uptime</th>
                            <th>Upload</th>
                            <th>Download</th>
                            <th>Total</th>
                            <th class="text-center">Hành động</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="user in filteredUsers" :key="user.username">
                            <td class="ps-4 fw-bold text-dark">{{ user.username }}</td>
                            <td class="text-secondary">{{ user.nas_ip }}</td>
                            <td class="font-monospace text-dark">{{ user.framed_ip }}</td>
                            <td class="text-secondary">{{ user.uptime }}</td>
                            <td class="text-secondary">{{ formatBytes(user.upload) }}</td>
                            <td class="text-secondary">{{ formatBytes(user.download) }}</td>
                            <td class="fw-bold text-primary">{{ formatBytes(user.total) }}</td>
                            <td class="text-center">
                                <button class="btn btn-danger btn-sm text-white px-3 fw-bold" @click="kickUser(user.username)">
                                    <i class="fa-solid fa-power-off me-1"></i> Kick
                                </button>
                            </td>
                        </tr>
                        <tr v-if="filteredUsers.length === 0">
                            <td colspan="8" class="text-center py-5 text-muted">
                                <div v-if="loading" class="spinner-border spinner-border-sm text-primary"></div>
                                <span v-else>Không có user nào đang online.</span>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
            
            <!-- Pagination Fake -->
            <div class="card-footer bg-white py-3 d-flex justify-content-end align-items-center">
                <small class="text-secondary me-3">1-{{ filteredUsers.length }} of {{ filteredUsers.length }} items</small>
                <nav>
                    <ul class="pagination pagination-sm mb-0">
                        <li class="page-item disabled"><a class="page-link" href="#">&lt;</a></li>
                        <li class="page-item active"><a class="page-link" href="#">1</a></li>
                        <li class="page-item disabled"><a class="page-link" href="#">&gt;</a></li>
                    </ul>
                </nav>
            </div>
        </div>
    </div>
</template>

<style scoped>
.font-monospace { font-family: 'Consolas', 'Courier New', monospace; }
/* Nút Kick màu đỏ giống hình */
.btn-danger {
    background-color: #f56c6c;
    border-color: #f56c6c;
    border-radius: 4px;
    font-size: 0.8rem;
}
.btn-danger:hover {
    background-color: #f78989;
    border-color: #f78989;
}
</style>