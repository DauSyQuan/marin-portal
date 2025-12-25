<script setup>
import { ref } from 'vue';

const activeTab = ref('network');
const snrThreshold = ref(5);
const firewallEnabled = ref(true);
const emailAlerts = ref(true);

const saveSettings = () => {
    // Giả lập lưu
    const btn = document.getElementById('saveBtn');
    const originalText = btn.innerHTML;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Saving...';
    setTimeout(() => {
        btn.innerHTML = '<i class="fa-solid fa-check me-2"></i>Saved!';
        btn.classList.replace('btn-primary', 'btn-success');
        setTimeout(() => {
            btn.innerHTML = originalText;
            btn.classList.replace('btn-success', 'btn-primary');
        }, 2000);
    }, 1000);
};
</script>

<template>
    <div class="container-fluid p-0">
        <h3 class="fw-bold text-dark mb-4">System Configuration</h3>

        <div class="row g-4">
            <!-- SETTINGS SIDEBAR -->
            <div class="col-md-3">
                <div class="card border-0 shadow-sm">
                    <div class="list-group list-group-flush rounded-3">
                        <button class="list-group-item list-group-item-action py-3 d-flex align-items-center"
                            :class="{ 'active fw-bold': activeTab === 'network' }" @click="activeTab = 'network'">
                            <i class="fa-solid fa-network-wired me-3 w-20"></i> Network Policies
                        </button>
                        <button class="list-group-item list-group-item-action py-3 d-flex align-items-center"
                            :class="{ 'active fw-bold': activeTab === 'alerts' }" @click="activeTab = 'alerts'">
                            <i class="fa-solid fa-bell me-3 w-20"></i> Alert Thresholds
                        </button>
                        <button class="list-group-item list-group-item-action py-3 d-flex align-items-center"
                            :class="{ 'active fw-bold': activeTab === 'users' }" @click="activeTab = 'users'">
                            <i class="fa-solid fa-users-gear me-3 w-20"></i> User Management
                        </button>
                    </div>
                </div>
            </div>

            <!-- SETTINGS CONTENT -->
            <div class="col-md-9">
                <div class="card border-0 shadow-sm">
                    <div class="card-body p-4">
                        
                        <!-- TAB: NETWORK -->
                        <div v-if="activeTab === 'network'">
                            <h5 class="fw-bold mb-4">Global Network Firewall</h5>
                            <div class="d-flex justify-content-between align-items-center p-3 border rounded mb-3 bg-light">
                                <div>
                                    <div class="fw-bold">Master Firewall</div>
                                    <small class="text-secondary">Enable/Disable deep packet inspection for all vessels.</small>
                                </div>
                                <div class="form-check form-switch">
                                    <input class="form-check-input" type="checkbox" v-model="firewallEnabled" style="width: 3em; height: 1.5em;">
                                </div>
                            </div>

                            <h6 class="fw-bold mt-4 mb-3">Application Control</h6>
                            <ul class="list-group">
                                <li class="list-group-item d-flex justify-content-between align-items-center">
                                    <span><i class="fa-brands fa-youtube text-danger me-2"></i> YouTube Streaming</span>
                                    <div class="form-check form-switch"><input class="form-check-input" type="checkbox"></div>
                                </li>
                                <li class="list-group-item d-flex justify-content-between align-items-center">
                                    <span><i class="fa-brands fa-facebook text-primary me-2"></i> Social Media</span>
                                    <div class="form-check form-switch"><input class="form-check-input" type="checkbox" checked></div>
                                </li>
                            </ul>
                        </div>

                        <!-- TAB: ALERTS -->
                        <div v-if="activeTab === 'alerts'">
                            <h5 class="fw-bold mb-4">Notification Rules</h5>
                            <div class="mb-4">
                                <label class="form-label fw-bold">Low Signal Warning (SNR)</label>
                                <div class="d-flex align-items-center gap-3">
                                    <input type="range" class="form-range" min="0" max="10" step="0.5" v-model="snrThreshold">
                                    <span class="badge bg-danger">{{ snrThreshold }} dB</span>
                                </div>
                                <div class="form-text">Trigger alert when signal drops below this value.</div>
                            </div>
                            <div class="form-check mb-3">
                                <input class="form-check-input" type="checkbox" v-model="emailAlerts" id="emailCheck">
                                <label class="form-check-label" for="emailCheck">Send Email Notifications</label>
                            </div>
                            <div class="mb-3">
                                <input type="email" class="form-control" value="admin@marine-corp.com" :disabled="!emailAlerts">
                            </div>
                        </div>

                        <!-- TAB: USERS -->
                        <div v-if="activeTab === 'users'">
                            <div class="d-flex justify-content-between align-items-center mb-4">
                                <h5 class="fw-bold m-0">System Administrators</h5>
                                <button class="btn btn-sm btn-primary">Add User</button>
                            </div>
                            <table class="table align-middle">
                                <thead class="bg-light"><tr><th>User</th><th>Role</th><th>Status</th></tr></thead>
                                <tbody>
                                    <tr>
                                        <td><div class="fw-bold">admin</div><small>admin@system.com</small></td>
                                        <td><span class="badge bg-dark">Super Admin</span></td>
                                        <td><span class="text-success"><i class="fa-solid fa-circle fa-2xs"></i> Active</span></td>
                                    </tr>
                                    <tr>
                                        <td><div class="fw-bold">manager</div><small>fleet@system.com</small></td>
                                        <td><span class="badge bg-secondary">Viewer</span></td>
                                        <td><span class="text-success"><i class="fa-solid fa-circle fa-2xs"></i> Active</span></td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>

                        <!-- SAVE BUTTON -->
                        <div class="mt-4 pt-3 border-top text-end">
                            <button id="saveBtn" class="btn btn-primary px-4 fw-bold" @click="saveSettings">
                                Save Configuration
                            </button>
                        </div>

                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.w-20 { width: 20px; text-align: center; }
</style>