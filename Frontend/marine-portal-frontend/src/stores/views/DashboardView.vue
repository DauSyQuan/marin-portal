<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { useFleetStore } from "@/stores/fleet";
import { Modal } from "bootstrap";

const store = useFleetStore();

// ===== Modal =====
let modalInstance = null;
const modalRef = ref(null);

// ===== Form Add Vessel =====
const form = reactive({
  id: "",
  name: "",
  company: "VIMC Lines",
  type: "Container",
  ip: "10.10.x.x",
  satellite: "JCSAT-4B",
  beam: "Spot-1",
  lat: 10.0,
  lon: 106.0,
});

// ===== UI State =====
const isLoading = ref(true);

// ===== Stats fallback (nếu store.stats rỗng) =====
const safeStats = computed(() => {
  if (Array.isArray(store.stats) && store.stats.length > 0) return store.stats;

  // fallback: tính từ ships nếu store.stats chưa set
  const ships = store.ships || store.filteredShips || [];
  const total = ships.length;
  const online = ships.filter((s) => s.status === "Online").length;
  const warning = ships.filter((s) => s.status === "Warning").length;
  const offline = ships.filter((s) => s.status === "Offline").length;

  return [
    { label: "Total Vessels", value: total, color: "primary", icon: "fa-solid fa-ship" },
    { label: "Online", value: online, color: "success", icon: "fa-solid fa-signal" },
    { label: "Warning", value: warning, color: "warning", icon: "fa-solid fa-triangle-exclamation" },
    { label: "Offline", value: offline, color: "secondary", icon: "fa-solid fa-power-off" },
  ];
});

// ===== Helpers =====
const openAddModal = () => {
  Object.assign(form, {
    id: "",
    name: "",
    company: "VIMC Lines",
    type: "Container",
    ip: "10.10.x.x",
    satellite: "JCSAT-4B",
    beam: "Spot-1",
    lat: 10.0,
    lon: 106.0,
  });

  modalInstance?.show();
};

const saveShip = async () => {
  if (!form.id?.trim() || !form.name?.trim()) {
    alert("Vui lòng nhập ID và Tên tàu!");
    return;
  }

  const payload = {
    ...form,
    id: form.id.trim(),
    name: form.name.trim(),
    lat: Number(form.lat),
    lon: Number(form.lon),
    snr: 12.0,
    status: "Online",
  };

  const success = await store.addShip(payload);

  if (success) {
    modalInstance?.hide();
  }
};

const statusClass = (status) => {
  if (status === "Online") return "badge-online";
  if (status === "Warning") return "badge-warning";
  if (status === "Blockage") return "badge-danger";
  return "badge-offline";
};

const statusDotClass = (status) => {
  if (status === "Online") return "dot-online";
  if (status === "Warning") return "dot-warning";
  if (status === "Blockage") return "dot-danger";
  return "dot-offline";
};

const refresh = async () => {
  isLoading.value = true;
  await store.fetchFleet();
  isLoading.value = false;
};

onMounted(async () => {
  try {
    await store.fetchFleet();
  } finally {
    isLoading.value = false;
  }

  if (modalRef.value) modalInstance = new Modal(modalRef.value);
});
</script>

<template>
  <div class="dashboard-page">
    <!-- Header -->
    <div class="header-row">
      <div>
        <div class="page-kicker">DASHBOARD</div>
        <h1 class="page-title">Fleet Command Center</h1>
      </div>

      <div class="header-actions">
        <button class="btn btn-light border shadow-sm" @click="refresh">
          <i class="fa-solid fa-rotate me-2 text-secondary"></i> Refresh
        </button>

        <button class="btn btn-primary shadow-sm" @click="openAddModal">
          <i class="fa-solid fa-plus me-2"></i> Add Vessel
        </button>
      </div>
    </div>

    <!-- Stats -->
    <div class="stats-grid">
      <template v-if="isLoading">
        <div class="stat-skeleton" v-for="i in 4" :key="i"></div>
      </template>

      <template v-else>
        <div class="stat-card" v-for="stat in safeStats" :key="stat.label">
          <div class="stat-top">
            <div class="stat-label">{{ stat.label }}</div>
            <div class="stat-icon" :class="`icon-${stat.color}`">
              <i :class="stat.icon"></i>
            </div>
          </div>
          <div class="stat-value">{{ stat.value }}</div>
        </div>
      </template>
    </div>

    <!-- Fleet Table -->
    <div class="card shadow-sm border-0 mt-4">
      <div class="table-header">
        <h6 class="m-0 fw-bold text-dark">
          <i class="fa-solid fa-list-ul me-2 text-secondary"></i> Live Fleet Status
        </h6>

        <div class="search-wrap">
          <i class="fa-solid fa-search search-icon"></i>
          <input
            type="text"
            class="form-control search-input"
            v-model="store.filterText"
            placeholder="Search IMO, Name..."
          />
        </div>
      </div>

      <div class="table-responsive">
        <table class="table table-hover align-middle mb-0">
          <thead class="thead">
            <tr>
              <th class="ps-4">IMO / ID</th>
              <th>Vessel & Satellite</th>
              <th>Owner</th>
              <th>Mgmt IP</th>
              <th>Status</th>
              <th class="text-end pe-4">Action</th>
            </tr>
          </thead>

          <tbody v-if="isLoading">
            <tr v-for="i in 5" :key="i">
              <td class="ps-4"><div class="row-skeleton"></div></td>
              <td><div class="row-skeleton"></div></td>
              <td><div class="row-skeleton"></div></td>
              <td><div class="row-skeleton"></div></td>
              <td><div class="row-skeleton"></div></td>
              <td class="text-end pe-4"><div class="row-skeleton"></div></td>
            </tr>
          </tbody>

          <tbody v-else>
            <tr v-if="store.filteredShips?.length === 0">
              <td colspan="6" class="empty-state">
                <div class="empty-title">No vessels found</div>
                <div class="empty-desc">Try another keyword or add a new vessel.</div>
                <button class="btn btn-primary btn-sm mt-2" @click="openAddModal">
                  <i class="fa-solid fa-plus me-2"></i> Add Vessel
                </button>
              </td>
            </tr>

            <tr v-for="ship in store.filteredShips" :key="ship.id">
              <td class="ps-4 font-monospace text-muted small">{{ ship.id }}</td>

              <td>
                <div class="fw-bold text-dark">{{ ship.name }}</div>
                <div class="small text-muted">
                  <i class="fa-solid fa-satellite-dish me-1 text-primary"></i>{{ ship.satellite }}
                  <span class="ms-2 text-muted">•</span>
                  <span class="ms-2">{{ ship.beam }}</span>
                </div>
              </td>

              <td>
                <span class="owner-pill">{{ ship.company }}</span>
              </td>

              <td class="font-monospace text-primary small">
                {{ ship.ip }}
              </td>

              <td>
                <div class="status-badge" :class="statusClass(ship.status)">
                  <span class="status-dot" :class="statusDotClass(ship.status)"></span>
                  {{ ship.status }}
                </div>

                <div
                  v-if="ship.status !== 'Offline'"
                  class="small text-muted mt-1"
                >
                  SNR: <strong>{{ Number(ship.snr).toFixed(1) }} dB</strong>
                </div>
              </td>

              <td class="text-end pe-4">
                <button class="btn btn-sm btn-outline-primary" @click="$router.push('/ship/' + ship.id)">
                  <i class="fa-solid fa-chart-line me-1"></i> Monitor
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL ADD VESSEL -->
    <div class="modal fade" id="addVesselModal" tabindex="-1" ref="modalRef">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content modal-modern">
          <div class="modal-header">
            <h5 class="modal-title fw-bold">Commission New Vessel</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>

          <div class="modal-body">
            <form @submit.prevent="saveShip">
              <div class="mb-3">
                <label class="form-label small text-secondary fw-bold">IMO Number</label>
                <input v-model="form.id" class="form-control" placeholder="e.g. IMO999888" required />
              </div>

              <div class="mb-3">
                <label class="form-label small text-secondary fw-bold">Vessel Name</label>
                <input v-model="form.name" class="form-control" placeholder="e.g. PACIFIC STAR" required />
              </div>

              <div class="row g-3 mb-3">
                <div class="col-6">
                  <label class="form-label small text-secondary fw-bold">Company</label>
                  <select v-model="form.company" class="form-select">
                    <option>VIMC Lines</option>
                    <option>Bien Dong POC</option>
                    <option>Maersk</option>
                    <option>New Client</option>
                  </select>
                </div>

                <div class="col-6">
                  <label class="form-label small text-secondary fw-bold">Type</label>
                  <select v-model="form.type" class="form-select">
                    <option>Container</option>
                    <option>Bulk Carrier</option>
                    <option>Tanker</option>
                    <option>Tug</option>
                  </select>
                </div>
              </div>

              <div class="row g-3 mb-3">
                <div class="col">
                  <label class="form-label small text-secondary fw-bold">Latitude</label>
                  <input v-model.number="form.lat" type="number" step="any" class="form-control" />
                </div>
                <div class="col">
                  <label class="form-label small text-secondary fw-bold">Longitude</label>
                  <input v-model.number="form.lon" type="number" step="any" class="form-control" />
                </div>
              </div>

              <div class="text-end border-top pt-3 mt-4">
                <button type="button" class="btn btn-light me-2" data-bs-dismiss="modal">Cancel</button>
                <button type="submit" class="btn btn-primary px-4 fw-bold">
                  Add Vessel
                </button>
              </div>
            </form>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ===== Layout ===== */
.dashboard-page {
  margin-left: 260px; /* sidebar width */
  padding: 28px 32px;
  min-height: 100vh;
  background: var(--bg-body, #f1f5f9);
}

/* Responsive: khi màn nhỏ */
@media (max-width: 992px) {
  .dashboard-page {
    margin-left: 0;
    padding: 18px;
  }
}

/* ===== Header ===== */
.header-row {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
}

.page-kicker {
  font-size: 0.72rem;
  letter-spacing: 0.12em;
  font-weight: 800;
  color: #64748b;
}

.page-title {
  margin: 0;
  font-weight: 900;
  letter-spacing: -0.02em;
  color: #0f172a;
  font-size: 2rem;
}

.header-actions {
  display: flex;
  gap: 10px;
}

/* ===== Stats ===== */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

@media (max-width: 1200px) {
  .stats-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

@media (max-width: 576px) {
  .stats-grid { grid-template-columns: 1fr; }
}

.stat-card {
  background: #fff;
  border-radius: 14px;
  padding: 16px 16px 14px;
  border: 1px solid rgba(15, 23, 42, 0.06);
  box-shadow: 0 6px 24px rgba(15, 23, 42, 0.06);
}

.stat-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.stat-label {
  font-size: 0.72rem;
  letter-spacing: 0.12em;
  font-weight: 800;
  color: #64748b;
  text-transform: uppercase;
}

.stat-value {
  font-size: 2.1rem;
  font-weight: 900;
  color: #0f172a;
  letter-spacing: -0.03em;
}

.stat-icon {
  width: 42px;
  height: 42px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-primary { background: rgba(59, 130, 246, 0.12); color: #3b82f6; }
.icon-success { background: rgba(34, 197, 94, 0.12); color: #22c55e; }
.icon-warning { background: rgba(245, 158, 11, 0.14); color: #f59e0b; }
.icon-secondary { background: rgba(100, 116, 139, 0.12); color: #64748b; }

.stat-skeleton {
  height: 108px;
  border-radius: 14px;
  background: linear-gradient(90deg, #ffffff, #f1f5f9, #ffffff);
  background-size: 200% 200%;
  animation: shimmer 1.2s infinite;
  border: 1px solid rgba(15, 23, 42, 0.06);
  box-shadow: 0 6px 24px rgba(15, 23, 42, 0.04);
}

@keyframes shimmer {
  0% { background-position: 0% 50%; }
  100% { background-position: 100% 50%; }
}

/* ===== Table ===== */
.table-header {
  padding: 14px 16px;
  border-bottom: 1px solid rgba(15, 23, 42, 0.07);
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.search-wrap {
  position: relative;
  width: min(360px, 100%);
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #64748b;
}

.search-input {
  padding-left: 36px;
  border-radius: 12px;
}

.thead th {
  background: #f8fafc;
  color: #64748b;
  font-size: 0.78rem;
  letter-spacing: 0.08em;
  font-weight: 800;
  text-transform: uppercase;
  border-bottom: 1px solid rgba(15, 23, 42, 0.07);
}

/* Owner pill */
.owner-pill {
  display: inline-flex;
  align-items: center;
  padding: 6px 10px;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 700;
  color: #475569;
  background: rgba(100, 116, 139, 0.1);
  border: 1px solid rgba(100, 116, 139, 0.2);
}

/* Status badge */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 800;
  border: 1px solid transparent;
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 999px;
}

.badge-online { background: rgba(34, 197, 94, 0.12); color: #16a34a; border-color: rgba(34, 197, 94, 0.3); }
.badge-warning { background: rgba(245, 158, 11, 0.14); color: #d97706; border-color: rgba(245, 158, 11, 0.32); }
.badge-danger { background: rgba(239, 68, 68, 0.12); color: #dc2626; border-color: rgba(239, 68, 68, 0.32); }
.badge-offline { background: rgba(100, 116, 139, 0.12); color: #475569; border-color: rgba(100, 116, 139, 0.25); }

.dot-online { background: #22c55e; }
.dot-warning { background: #f59e0b; }
.dot-danger { background: #ef4444; }
.dot-offline { background: #64748b; }

/* Empty state */
.empty-state {
  padding: 42px 20px !important;
  text-align: center;
  background: #fff;
}

.empty-title {
  font-weight: 900;
  color: #0f172a;
  font-size: 1.05rem;
}

.empty-desc {
  color: #64748b;
  margin-top: 6px;
  font-size: 0.92rem;
}

/* Skeleton rows */
.row-skeleton {
  height: 14px;
  width: 80%;
  border-radius: 8px;
  background: linear-gradient(90deg, #ffffff, #f1f5f9, #ffffff);
  background-size: 200% 200%;
  animation: shimmer 1.2s infinite;
}

/* Modal */
.modal-modern {
  border-radius: 16px;
  border: 0;
  box-shadow: 0 18px 60px rgba(15, 23, 42, 0.22);
}
</style>
