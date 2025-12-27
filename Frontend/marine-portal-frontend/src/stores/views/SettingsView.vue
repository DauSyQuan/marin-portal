<script setup>
import { ref, reactive, onMounted } from 'vue';
import api from '@/services/api';

const activeTab = ref('network');
const isSaving = ref(false);
const auditLogs = ref([]);

// --- STATE CẤU HÌNH ---
const networkConfig = reactive({
  firewall: true,
  dpi: true,
  sdwan_mode: 'failover',
  primary_link: 'vsat_ka',
  secondary_link: '4g_lte',
  blocked_apps: {
    youtube: false,
    facebook: true,
    tiktok: true,
    torrent: true
  }
});

const alertsConfig = reactive({
  snr_threshold: 5.0,
  daily_quota_warning: 80,
  recipients: 'noc@marine-corp.com; fleet.mgr@marine-corp.com',
  channels: { email: true, sms: false, webhook: true }
});

const integrations = reactive({
  starlink_api: '',
  inmarsat_key: '',
  weather_api: true
});

// --- Helper: map state UI -> payload backend ---
const buildPayload = () => ({
  primary_link: networkConfig.primary_link,
  sdwan_mode: networkConfig.sdwan_mode,
  firewall: networkConfig.firewall,

  block_youtube: networkConfig.blocked_apps.youtube,
  block_facebook: networkConfig.blocked_apps.facebook,
  block_tiktok: networkConfig.blocked_apps.tiktok,

  snr_threshold: Number(alertsConfig.snr_threshold),
  quota_warning: Number(alertsConfig.daily_quota_warning),
  recipients: alertsConfig.recipients,

  starlink_api_key: integrations.starlink_api,
  weather_api: !!integrations.weather_api,
});

// --- Helper: map response backend -> state UI ---
const applyFromBackend = (data) => {
  if (!data) return;

  if (data.primary_link !== undefined) networkConfig.primary_link = data.primary_link;
  if (data.sdwan_mode !== undefined) networkConfig.sdwan_mode = data.sdwan_mode;
  if (data.firewall !== undefined) networkConfig.firewall = !!data.firewall;

  if (data.block_youtube !== undefined) networkConfig.blocked_apps.youtube = !!data.block_youtube;
  if (data.block_facebook !== undefined) networkConfig.blocked_apps.facebook = !!data.block_facebook;
  if (data.block_tiktok !== undefined) networkConfig.blocked_apps.tiktok = !!data.block_tiktok;

  if (data.snr_threshold !== undefined) alertsConfig.snr_threshold = Number(data.snr_threshold);
  if (data.quota_warning !== undefined) alertsConfig.daily_quota_warning = Number(data.quota_warning);
  if (data.recipients !== undefined) alertsConfig.recipients = data.recipients;

  if (data.starlink_api_key !== undefined) integrations.starlink_api = data.starlink_api_key;
  if (data.weather_api !== undefined) integrations.weather_api = !!data.weather_api;
};

const loadSettings = async () => {
  try {
    const res = await api.get('/api/settings');
    applyFromBackend(res.data);
  } catch (e) {
    console.error("loadSettings error:", e);
  }
};

const loadLogs = async () => {
  try {
    const res = await api.get('/api/audit-logs');
    auditLogs.value = Array.isArray(res.data) ? res.data : [];
  } catch (e) {
    console.error("loadLogs error:", e);
  }
};

onMounted(async () => {
  await loadSettings();
  await loadLogs();
});

const saveSettings = async () => {
  isSaving.value = true;
  try {
    await api.put('/api/settings', buildPayload());
    await loadSettings();
    await loadLogs();
  } catch (e) {
    alert("Lỗi lưu: " + (e?.response?.data?.error || e.message));
  } finally {
    setTimeout(() => { isSaving.value = false; }, 800);
  }
};

const formatDate = (d) => new Date(d).toLocaleString();
</script>

<template>
    <div class="container-fluid p-0 fade-in-tab">
        <h3 class="fw-bold text-white mb-4">System Configuration</h3>

        <div class="row g-4">
            <!-- SIDEBAR SETTINGS -->
            <div class="col-md-3">
                <div class="card glass-panel h-100">
                    <div class="list-group list-group-flush p-2">
                        <button class="list-group-item list-group-item-action py-3 rounded mb-1" 
                            :class="{ 'active fw-bold': activeTab === 'network' }" @click="activeTab = 'network'">
                            <i class="fa-solid fa-network-wired me-3 w-20"></i> SD-WAN & Network
                        </button>
                        <button class="list-group-item list-group-item-action py-3 rounded mb-1" 
                            :class="{ 'active fw-bold': activeTab === 'alerts' }" @click="activeTab = 'alerts'">
                            <i class="fa-solid fa-bell me-3 w-20"></i> Alerts & Notifications
                        </button>
                        <button class="list-group-item list-group-item-action py-3 rounded mb-1" 
                            :class="{ 'active fw-bold': activeTab === 'integrations' }" @click="activeTab = 'integrations'">
                            <i class="fa-solid fa-plug me-3 w-20"></i> Integrations (API)
                        </button>
                        <button class="list-group-item list-group-item-action py-3 rounded mb-1" 
                            :class="{ 'active fw-bold': activeTab === 'audit' }" @click="activeTab = 'audit'">
                            <i class="fa-solid fa-shield-halved me-3 w-20"></i> Audit Logs
                        </button>
                    </div>
                </div>
            </div>

            <!-- CONTENT AREA -->
            <div class="col-md-9">
                <!-- ✅ CHỈ THÊM class settings-card ĐỂ CHỪA CHỖ NÚT SAVE -->
                <div class="card glass-panel h-100 p-4 position-relative overflow-hidden settings-card">
                    
                    <!-- TAB 1: SD-WAN & NETWORK -->
                    <div v-if="activeTab === 'network'">
                        <h5 class="fw-bold text-white mb-4 border-bottom border-white border-opacity-10 pb-2">Global Traffic Steering (SD-WAN)</h5>
                        
                        <!-- SD-WAN Priority -->
                        <div class="row g-4 mb-4">
                            <div class="col-md-6">
                                <label class="small fw-bold text-secondary mb-2">Primary Uplink</label>
                                <select class="form-select" v-model="networkConfig.primary_link">
                                    <option value="vsat_ka">High Speed VSAT (Ka-Band)</option>
                                    <option value="l_band">L-Band (Backup)</option>
                                    <option value="4g_lte">4G/LTE Coastal</option>
                                    <option value="starlink">Starlink Maritime</option>
                                </select>
                            </div>
                            <div class="col-md-6">
                                <label class="small fw-bold text-secondary mb-2">Failover Mode</label>
                                <div class="d-flex gap-2">
                                    <div class="form-check custom-radio-box flex-fill">
                                        <input class="form-check-input" type="radio" name="sdwan" id="failover" value="failover" v-model="networkConfig.sdwan_mode">
                                        <label class="form-check-label w-100" for="failover"><i class="fa-solid fa-arrow-right-arrow-left me-2"></i>Failover</label>
                                    </div>
                                    <div class="form-check custom-radio-box flex-fill">
                                        <input class="form-check-input" type="radio" name="sdwan" id="load_balance" value="load_balance" v-model="networkConfig.sdwan_mode">
                                        <label class="form-check-label w-100" for="load_balance"><i class="fa-solid fa-scale-balanced me-2"></i>Load Balance</label>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <h6 class="fw-bold text-white mt-4 mb-3">Firewall & Deep Packet Inspection (DPI)</h6>
                        <div class="glass-inset p-3 rounded mb-3">
                            <div class="d-flex justify-content-between align-items-center mb-3">
                                <div><div class="fw-bold text-white">Master Firewall</div><small class="text-secondary">Enforce policies across all fleet routers.</small></div>
                                <div class="form-check form-switch"><input class="form-check-input" type="checkbox" v-model="networkConfig.firewall"></div>
                            </div>
                            <hr class="border-white border-opacity-10">
                            <div class="row g-3">
                                <div class="col-6 d-flex justify-content-between">
                                    <span class="text-white-50"><i class="fa-brands fa-youtube me-2 text-danger"></i> YouTube</span>
                                    <div class="form-check form-switch"><input class="form-check-input" type="checkbox" v-model="networkConfig.blocked_apps.youtube"></div>
                                </div>
                                <div class="col-6 d-flex justify-content-between">
                                    <span class="text-white-50"><i class="fa-brands fa-tiktok me-2 text-white"></i> TikTok</span>
                                    <div class="form-check form-switch"><input class="form-check-input" type="checkbox" v-model="networkConfig.blocked_apps.tiktok"></div>
                                </div>
                                <div class="col-6 d-flex justify-content-between">
                                    <span class="text-white-50"><i class="fa-brands fa-facebook me-2 text-primary"></i> Social Media</span>
                                    <div class="form-check form-switch"><input class="form-check-input" type="checkbox" v-model="networkConfig.blocked_apps.facebook"></div>
                                </div>
                                <div class="col-6 d-flex justify-content-between">
                                    <span class="text-white-50"><i class="fa-solid fa-download me-2 text-warning"></i> P2P / Torrent</span>
                                    <div class="form-check form-switch"><input class="form-check-input" type="checkbox" v-model="networkConfig.blocked_apps.torrent"></div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- TAB 2: ALERTS -->
                    <div v-if="activeTab === 'alerts'">
                        <h5 class="fw-bold text-white mb-4 border-bottom border-white border-opacity-10 pb-2">Smart Alerting System</h5>
                        
                        <div class="row g-4 mb-4">
                            <div class="col-md-6">
                                <label class="fw-bold text-white mb-2">Signal Quality Threshold (SNR)</label>
                                <div class="d-flex align-items-center gap-3">
                                    <input type="range" class="form-range" min="0" max="15" step="0.5" v-model="alertsConfig.snr_threshold">
                                    <span class="badge bg-danger fs-6">{{ alertsConfig.snr_threshold }} dB</span>
                                </div>
                                <div class="form-text text-secondary">Alert when satellite signal drops below this value.</div>
                            </div>
                            <div class="col-md-6">
                                <label class="fw-bold text-white mb-2">Daily Data Quota Warning</label>
                                <div class="input-group">
                                    <input type="number" class="form-control" v-model="alertsConfig.daily_quota_warning">
                                    <span class="input-group-text">% used</span>
                                </div>
                            </div>
                        </div>

                        <div class="mb-3">
                            <label class="fw-bold text-white mb-2">Notification Recipients</label>
                            <input type="text" class="form-control" v-model="alertsConfig.recipients">
                        </div>

                        <div class="d-flex gap-4 mt-4">
                            <div class="form-check"><input class="form-check-input" type="checkbox" v-model="alertsConfig.channels.email" id="c1"><label class="form-check-label text-white" for="c1">Email</label></div>
                            <div class="form-check"><input class="form-check-input" type="checkbox" v-model="alertsConfig.channels.sms" id="c2"><label class="form-check-label text-white" for="c2">SMS (Critical Only)</label></div>
                            <div class="form-check"><input class="form-check-input" type="checkbox" v-model="alertsConfig.channels.webhook" id="c3"><label class="form-check-label text-white" for="c3">Slack/Teams Webhook</label></div>
                        </div>
                    </div>

                    <!-- TAB 3: INTEGRATIONS -->
                    <div v-if="activeTab === 'integrations'">
                        <h5 class="fw-bold text-white mb-4 border-bottom border-white border-opacity-10 pb-2">3rd Party Integrations</h5>
                        
                        <div class="glass-inset p-4 rounded mb-3">
                            <div class="d-flex align-items-center mb-3">
                                <i class="fa-solid fa-satellite me-3 fa-2x text-white"></i>
                                <div class="flex-grow-1">
                                    <h6 class="fw-bold text-white m-0">Starlink Maritime API</h6>
                                    <small class="text-success"><i class="fa-solid fa-circle fa-2xs me-1"></i>Connected</small>
                                </div>
                                <button class="btn btn-sm btn-outline-light">Configure</button>
                            </div>

                            <!-- ✅ Không thay đổi UI, chỉ đổi thành v-model -->
                            <input type="password" class="form-control" v-model="integrations.starlink_api" placeholder="Enter Starlink API key">
                        </div>

                        <div class="glass-inset p-4 rounded">
                            <div class="d-flex align-items-center mb-3">
                                <i class="fa-solid fa-cloud-sun-rain me-3 fa-2x text-info"></i>
                                <div class="flex-grow-1">
                                    <h6 class="fw-bold text-white m-0">Global Weather Data</h6>
                                    <small class="text-secondary">Overlay weather maps on fleet view.</small>
                                </div>
                                <div class="form-check form-switch">
                                    <input class="form-check-input" type="checkbox" v-model="integrations.weather_api">
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- TAB 4: AUDIT LOGS -->
                    <div v-if="activeTab === 'audit'">
                        <div class="d-flex justify-content-between align-items-center mb-4 border-bottom border-white border-opacity-10 pb-2">
                            <h5 class="fw-bold text-white m-0">Audit Logs</h5>
                            <button class="btn btn-sm btn-outline-secondary"><i class="fa-solid fa-download me-2"></i>Export CSV</button>
                        </div>
                        <div class="table-responsive">
                            <table class="table table-hover align-middle">
                                <thead><tr><th>Time</th><th>User</th><th>Action</th><th>IP Address</th><th>Status</th></tr></thead>
                                <tbody>
                                    <tr v-for="log in auditLogs" :key="log.id">
                                        <td class="text-secondary font-monospace small">{{ formatDate(log.created_at) }}</td>
                                        <td class="fw-bold text-white">{{ log.user }}</td>
                                        <td class="text-white">{{ log.action }}</td>
                                        <td class="text-secondary font-monospace small">{{ log.ip_address }}</td>
                                        <td>
                                          <span
                                            class="badge"
                                            :class="log.status==='Success' ? 'bg-success bg-opacity-25 text-success' : 'bg-warning bg-opacity-25 text-warning'"
                                          >
                                            {{ log.status }}
                                          </span>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>

                    <!-- SAVE BUTTON -->
                    <div class="position-absolute bottom-0 end-0 p-4">
                        <button class="btn btn-primary px-4 fw-bold shadow-lg" @click="saveSettings" :disabled="isSaving">
                            <span v-if="isSaving"><i class="fa-solid fa-spinner fa-spin me-2"></i>Saving...</span>
                            <span v-else><i class="fa-solid fa-floppy-disk me-2"></i>Save Changes</span>
                        </button>
                    </div>

                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Sidebar Styles */
.list-group-item {
    background: transparent; color: #94a3b8; border: none; transition: all 0.2s;
}
.list-group-item:hover { background: rgba(255,255,255,0.05); color: white; }
.list-group-item.active {
    background: linear-gradient(90deg, rgba(59, 130, 246, 0.2) 0%, rgba(59, 130, 246, 0) 100%);
    color: #3b82f6; border-left: 3px solid #3b82f6; border-radius: 4px;
}
.w-20 { width: 20px; text-align: center; }

/* Custom Radio Box for SD-WAN */
.custom-radio-box {
    position: relative;
}
.custom-radio-box input {
    position: absolute; opacity: 0; cursor: pointer;
}
.custom-radio-box label {
    display: block; padding: 12px; border: 1px solid rgba(255,255,255,0.1); border-radius: 8px;
    background: rgba(0,0,0,0.2); cursor: pointer; transition: all 0.2s; text-align: center; color: #94a3b8; font-weight: 600;
}
.custom-radio-box input:checked + label {
    background: rgba(59, 130, 246, 0.2); border-color: #3b82f6; color: #3b82f6;
}

/* Glass Inset (Lõm xuống) */
.glass-inset {
    background: rgba(0, 0, 0, 0.2); border: 1px solid rgba(255,255,255,0.05);
}

.fade-in-tab { animation: fadeIn 0.4s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }

/* ✅ Chừa khoảng trống đáy để nút Save Changes không đè nội dung */
.settings-card {
  padding-bottom: 110px !important;
}
</style>
