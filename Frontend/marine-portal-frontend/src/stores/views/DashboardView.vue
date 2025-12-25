<script setup>
import { onMounted, reactive, ref } from 'vue';
import { useFleetStore } from '@/stores/fleet';
import { Modal } from 'bootstrap'; // Import Bootstrap Modal

const store = useFleetStore();

// Biến quản lý Modal
let modalInstance = null;
const modalRef = ref(null);

// Form dữ liệu thêm mới
const form = reactive({
    id: '', name: '', company: 'VIMC Lines', 
    type: 'Container', ip: '10.10.x.x', 
    satellite: 'JCSAT-4B', beam: 'Spot-1',
    lat: 10.0, lon: 106.0
});

onMounted(() => {
    store.fetchFleet();
    // Khởi tạo Modal Bootstrap khi trang tải xong
    if(modalRef.value) {
        modalInstance = new Modal(modalRef.value);
    }
});

// 1. HÀM MỞ MODAL (Reset form và hiện bảng nhập)
const openAddModal = () => {
    Object.assign(form, { 
        id: '', name: '', company: 'VIMC Lines', 
        type: 'Container', ip: '10.10.x.x', 
        satellite: 'JCSAT-4B', beam: 'Spot-1',
        lat: 10.0, lon: 106.0 
    });
    modalInstance.show();
};

// 2. HÀM LƯU TÀU (Gửi dữ liệu lên Backend)
const saveShip = async () => {
    // Validate cơ bản
    if(!form.id || !form.name) return alert("Vui lòng nhập ID và Tên tàu!");
    
    // Ép kiểu dữ liệu sang số thực (Float) để Backend Go hiểu
    const payload = {
        ...form,
        lat: parseFloat(form.lat || 0),
        lon: parseFloat(form.lon || 0),
        snr: 12.0,       // Mặc định
        status: 'Online' // Mặc định
    };

    console.log("Đang gửi dữ liệu:", payload);

    // Gọi Store để post API
    const success = await store.addShip(payload);
    
    if(success) {
        modalInstance.hide(); // Tắt modal nếu thành công
        // Reset form để tránh lỗi lưu đè lần sau
        Object.assign(form, { id: '', name: '', lat: 0, lon: 0 });
    }
};

// Hàm CSS trạng thái
const statusClass = (status) => {
    if(status === 'Online') return 'bg-success text-success border-success';
    if(status === 'Warning') return 'bg-warning text-warning border-warning';
    if(status === 'Blockage') return 'bg-danger text-danger border-danger';
    return 'bg-secondary text-secondary border-secondary';
};
</script>

<template>
  <div class="container-fluid p-0">
    <!-- HEADER DASHBOARD -->
    <div class="d-flex justify-content-between align-items-end mb-4">
        <div>
            <h6 class="text-secondary text-uppercase mb-1 fw-bold" style="font-size: 0.7rem;">Dashboard</h6>
            <h2 class="fw-bold text-dark m-0">Fleet Command Center</h2>
        </div>
        <div>
            <button class="btn btn-white border me-2 shadow-sm bg-white" @click="store.fetchFleet">
                <i class="fa-solid fa-rotate me-2 text-secondary"></i> Refresh
            </button>
            <!-- NÚT NÀY GỌI openAddModal (MỞ FORM) -->
            <button class="btn btn-primary shadow-sm" @click="openAddModal">
                <i class="fa-solid fa-plus me-2"></i> Add Vessel
            </button>
        </div>
    </div>

    <!-- STATS CARDS -->
    <div class="row mb-4 g-4">
        <div class="col-md-3" v-for="stat in store.stats" :key="stat.label">
            <div class="card h-100 border-0 border-start border-4 shadow-sm" :class="'border-'+stat.color">
                <div class="card-body">
                    <div class="stat-card-label mb-2 text-uppercase fw-bold text-secondary" style="font-size: 0.7rem;">{{ stat.label }}</div>
                    <div class="d-flex justify-content-between align-items-center">
                        <div class="fs-2 fw-bold" :class="'text-'+stat.color">{{ stat.value }}</div>
                        <div :class="`bg-${stat.color} bg-opacity-10 p-3 rounded-circle`">
                            <i :class="stat.icon + ' text-'+stat.color + ' fa-lg'"></i>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- DANH SÁCH TÀU -->
    <div class="card shadow-sm border-0">
        <div class="card-header d-flex justify-content-between align-items-center bg-white py-3">
            <h6 class="m-0 fw-bold text-dark"><i class="fa-solid fa-list-ul me-2 text-secondary"></i> Live Fleet Status</h6>
            <div class="input-group w-25">
                <span class="input-group-text bg-white"><i class="fa-solid fa-search text-secondary"></i></span>
                <input type="text" class="form-control" v-model="store.filterText" placeholder="Search IMO, Name...">
            </div>
        </div>
        <div class="table-responsive">
            <table class="table table-hover align-middle mb-0">
                <thead class="bg-light">
                    <tr>
                        <th class="ps-4 text-secondary text-uppercase small">IMO / ID</th>
                        <th class="text-secondary text-uppercase small">Vessel & Satellite</th>
                        <th class="text-secondary text-uppercase small">Owner</th>
                        <th class="text-secondary text-uppercase small">Mgmt IP</th>
                        <th class="text-secondary text-uppercase small">Status</th>
                        <th class="text-end pe-4 text-secondary text-uppercase small">Action</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="ship in store.filteredShips" :key="ship.id">
                        <td class="ps-4 font-monospace text-secondary small">{{ ship.id }}</td>
                        <td>
                            <div class="fw-bold text-dark">{{ ship.name }}</div>
                            <div class="small text-muted" style="font-size: 0.75rem;">
                                <i class="fa-solid fa-satellite-dish me-1 text-primary"></i>{{ ship.satellite }} 
                            </div>
                        </td>
                        <td><span class="badge bg-secondary bg-opacity-10 text-secondary border border-secondary border-opacity-25">{{ ship.company }}</span></td>
                        <td class="font-monospace text-primary small">{{ ship.ip }}</td>
                        <td>
                            <span class="badge bg-opacity-10 border border-opacity-25" :class="statusClass(ship.status)">
                                <i class="fa-solid fa-circle fa-2xs me-1"></i>{{ ship.status }}
                            </span>
                            <!-- Đã thêm toFixed(1) để làm tròn số -->
                            <div v-if="ship.status !== 'Offline'" class="small text-secondary mt-1" style="font-size: 0.7rem;">
                                SNR: <strong>{{ Number(ship.snr).toFixed(1) }} dB</strong>
                            </div>
                        </td>
                        <td class="text-end pe-4">
                            <button class="btn btn-sm btn-outline-primary" @click="$router.push('/ship/' + ship.id)">
                                <i class="fa-solid fa-chart-line me-1"></i>Monitor
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>

    <!-- MODAL FORM NHẬP LIỆU -->
    <div class="modal fade" id="addVesselModal" tabindex="-1" ref="modalRef">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title fw-bold">Commission New Vessel</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body">
                    <form @submit.prevent="saveShip">
                        <div class="mb-3">
                            <label class="form-label small text-secondary fw-bold">IMO Number</label>
                            <input v-model="form.id" class="form-control" placeholder="e.g. IMO999888" required>
                        </div>
                        <div class="mb-3">
                            <label class="form-label small text-secondary fw-bold">Vessel Name</label>
                            <input v-model="form.name" class="form-control" placeholder="e.g. PACIFIC STAR" required>
                        </div>
                        <div class="row mb-3">
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
                        <div class="row mb-3">
                             <div class="col"><label class="form-label small text-secondary">Latitude</label><input v-model.number="form.lat" type="number" step="any" class="form-control"></div>
                             <div class="col"><label class="form-label small text-secondary">Longitude</label><input v-model.number="form.lon" type="number" step="any" class="form-control"></div>
                        </div>
                        
                        <div class="text-end border-top pt-3 mt-4">
                            <button type="button" class="btn btn-light me-2" data-bs-dismiss="modal">Cancel</button>
                            <!-- NÚT NÀY GỌI saveShip (LƯU) -->
                            <button type="button" @click="saveShip" class="btn btn-primary px-4 fw-bold">Add Vessel</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>

  </div>
</template>