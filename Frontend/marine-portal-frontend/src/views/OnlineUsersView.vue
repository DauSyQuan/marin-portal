<script setup>
import { ref, onMounted, computed } from 'vue';
import api from '@/services/api';

// State
const users = ref([]);
const loading = ref(false);
const searchQuery = ref('');

// Fetch Data
const fetchOnlineUsers = async () => {
    loading.value = true;
    try {
        const res = await api.get('/api/online-users');
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
        (u.framed_ip || '').includes(searchQuery.value)
    );
});

// Kick User
const kickUser = async (username) => {
    try {
        await api.post(`/api/online-users/${username}/kick`);
        await fetchOnlineUsers();
    } catch (e) {
        console.error(e);
        alert('Không thể kick user: ' + (e?.response?.data?.error || e.message));
    }
};

onMounted(fetchOnlineUsers);
</script>

<template>
    <div class="container-fluid p-0 fade-in-tab">
        <div class="d-flex justify-content-between align-items-end mb-4">
            <div>
                <small class="text-info fw-bold tracking-wide">LIVE CONNECTIONS</small>
                <h3 class="fw-bold text-white m-0">Online Users</h3>
            </div>
            <div class="d-flex gap-2 align-items-center">
                <input
                    v-model="searchQuery"
                    type="text"
                    class="form-control form-control-sm bg-dark text-white border-secondary"
                    placeholder="Search username / IP..."
                    style="width: 230px;"
                />
                <button class="btn btn-outline-light btn-sm fw-bold" @click="fetchOnlineUsers" :disabled="loading">
                    <i class="fa-solid fa-rotate me-2"></i>Refresh
                </button>
            </div>
        </div>

        <div class="card glass-panel p-0">
            <div class="table-responsive">
                <table class="table table-hover align-middle mb-0">
                    <thead>
                        <tr>
                            <th class="ps-4">Username</th>
                            <th>Vessel</th>
                            <th>IP Address</th>
                            <th>Mac</th>
                            <th>Login Time</th>
                            <th class="text-end pe-4">Action</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="loading">
                            <td colspan="6" class="text-center py-4 text-white-50">Loading...</td>
                        </tr>

                        <tr v-else-if="filteredUsers.length === 0">
                            <td colspan="6" class="text-center py-4 text-white-50">No online users found.</td>
                        </tr>

                        <tr v-else v-for="u in filteredUsers" :key="u.username">
                            <td class="ps-4 fw-bold text-white">{{ u.username }}</td>
                            <td class="text-white-50">{{ u.ship_name || 'N/A' }}</td>
                            <td class="font-monospace text-info">{{ u.framed_ip || 'N/A' }}</td>
                            <td class="font-monospace text-white-50">{{ u.mac_address || 'N/A' }}</td>
                            <td class="text-white-50">{{ u.login_time || 'N/A' }}</td>
                            <td class="text-end pe-4">
                                <button class="btn btn-sm btn-outline-danger" @click="kickUser(u.username)">
                                    <i class="fa-solid fa-ban me-2"></i>Kick
                                </button>
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
.fade-in-tab { animation: fadeIn 0.4s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>
