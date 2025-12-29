<script setup>
import { ref, onMounted, reactive, nextTick } from 'vue';
import axios from 'axios';
import { Modal } from 'bootstrap';

// State
const vouchers = ref([]);
const ships = ref([]); 
const crews = ref([]); 
const loading = ref(false);

// Filter
const filter = reactive({ status: '', assignTo: '' });

// Forms
const createForm = reactive({ data_plan: '1GB', valid_days: 30 });
const assignForm = reactive({ voucherId: null, shipId: '', crewId: null, crewName: '' });

// Modal Refs
const createModalRef = ref(null);
const assignModalRef = ref(null);
let createModalObj = null;
let assignModalObj = null;

// --- 1. API CALLS ---
const fetchVouchers = async () => {
    loading.value = true;
    try {
        const res = await axios.get('http://localhost:8080/api/vouchers', {
            params: { status: filter.status, assign_to: filter.assignTo }
        });
        vouchers.value = res.data || [];
    } catch(e) {
        console.error("Lỗi tải voucher:", e);
    } finally { loading.value = false; }
};

const fetchShips = async () => {
    try {
        const res = await axios.get('http://localhost:8080/api/ships');
        ships.value = res.data || [];
    } catch(e) { console.error(e); }
};

const onShipSelect = async () => {
    if (!assignForm.shipId) return;
    crews.value = []; // Reset crew list
    try {
        const res = await axios.get(`http://localhost:8080/api/ships/${assignForm.shipId}/crew`);
        crews.value = res.data || [];
    } catch(e) { console.error(e); }
};

// --- 2. MODAL ACTIONS ---

// Mở Modal Tạo Mới (Đảm bảo Modal đã khởi tạo)
const openCreateModal = () => {
    if (!createModalObj && createModalRef.value) {
        createModalObj = new Modal(createModalRef.value);
    }
    // Reset form
    createForm.data_plan = '1GB';
    createForm.valid_days = 30;
    createModalObj?.show();
};

const createVoucher = async () => {
    try {
        await axios.post('http://localhost:8080/api/vouchers', createForm);
        createModalObj.hide();
        fetchVouchers(); // Refresh lại bảng
    } catch(e) {
        alert("Lỗi tạo voucher: " + e.message);
    }
};

// Mở Modal Gán (Assign)
const openAssignModal = (voucher) => {
    if (!assignModalObj && assignModalRef.value) {
        assignModalObj = new Modal(assignModalRef.value);
    }
    
    // Reset Form Gán
    assignForm.voucherId = voucher.id;
    assignForm.shipId = '';
    assignForm.crewId = null;
    assignForm.crewName = '';
    crews.value = [];
    
    fetchShips(); // Tải danh sách tàu để chọn
    assignModalObj?.show();
};

const assignVoucher = async () => {
    const selectedCrew = crews.value.find(c => c.id === assignForm.crewId);
    if (!selectedCrew) return alert("Vui lòng chọn thủy thủ!");

    try {
        await axios.put(`http://localhost:8080/api/vouchers/${assignForm.voucherId}/assign`, {
            crew_id: selectedCrew.id,
            crew_name: selectedCrew.full_name
        });
        assignModalObj.hide();
        fetchVouchers();
    } catch (e) {
        alert("Lỗi gán voucher: " + e.message);
    }
};

// --- LIFECYCLE ---
onMounted(async () => {
    await fetchVouchers();
    // Khởi tạo modal sau khi DOM đã render xong
    await nextTick();
    if(createModalRef.value) createModalObj = new Modal(createModalRef.value);
    if(assignModalRef.value) assignModalObj = new Modal(assignModalRef.value);
});
</script>

<template>
    <div class="container-fluid p-0">
        <!-- HEADER & FILTER -->
        <div class="card border-0 shadow-sm mb-4">
            <div class="card-body py-4">
                <div class="row g-3 align-items-center">
                    <div class="col-md-3">
                        <div class="d-flex align-items-center">
                            <label class="me-2 fw-bold text-secondary">Status:</label>
                            <select class="form-select" v-model="filter.status">
                                <option value="">All</option>
                                <option value="Unused">Unused</option>
                                <option value="Assigned">Assigned</option>
                            </select>
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="d-flex align-items-center">
                            <label class="me-2 fw-bold text-secondary" style="min-width: 80px;">Assign to:</label>
                            <input type="text" class="form-control" v-model="filter.assignTo" placeholder="Search name...">
                        </div>
                    </div>
                    <div class="col-md-5 text-end">
                        <button class="btn btn-light border me-2" @click="filter.status='';filter.assignTo='';fetchVouchers()">Reset</button>
                        <button class="btn btn-primary px-4" @click="fetchVouchers"><i class="fa-solid fa-magnifying-glass me-2"></i> Query</button>
                    </div>
                </div>
            </div>
        </div>

        <!-- TABLE -->
        <div class="card border-0 shadow-sm">
            <div class="card-header bg-white py-3 d-flex justify-content-between align-items-center">
                <h5 class="m-0 fw-bold">Vouchers</h5>
                <!-- NÚT BẤM GỌI HÀM openCreateModal -->
                <button class="btn btn-primary fw-bold" @click="openCreateModal">
                    <i class="fa-solid fa-plus me-2"></i> Create
                </button>
            </div>
            <div class="table-responsive">
                <table class="table table-hover align-middle mb-0">
                    <thead class="bg-light small text-uppercase text-secondary">
                        <tr>
                            <th class="ps-4">Status</th><th>Created By</th><th>Data Plan</th><th>Code</th><th>Assign To</th><th>Validity</th><th>Action</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="v in vouchers" :key="v.id">
                            <td class="ps-4">
                                <span class="badge" :class="v.status==='Unused'?'bg-success':'bg-secondary'">{{ v.status }}</span>
                            </td>
                            <td>{{ v.created_by }}</td>
                            <td class="fw-bold text-primary">{{ v.data_plan }}</td>
                            <!-- Code hiển thị đẹp -->
                            <td><span class="font-monospace bg-light border px-2 py-1 rounded text-dark">{{ v.code }}</span></td>
                            <td>
                                <div v-if="v.assign_to" class="fw-bold text-dark"><i class="fa-solid fa-user me-1"></i> {{ v.assign_to }}</div>
                                <span v-else class="text-muted small">--</span>
                            </td>
                            <td>{{ v.valid_days }} Days</td>
                            <td>
                                <!-- Chỉ hiện nút Assign nếu chưa dùng -->
                                <button v-if="v.status === 'Unused'" class="btn btn-sm btn-outline-primary fw-bold" @click="openAssignModal(v)">
                                    <i class="fa-solid fa-hand-holding-hand me-1"></i> Assign
                                </button>
                            </td>
                        </tr>
                        <tr v-if="vouchers.length === 0">
                            <td colspan="7" class="text-center py-5">
                                <i class="fa-solid fa-ticket fa-3x text-secondary opacity-25 mb-3"></i>
                                <div class="text-muted">No Data. Create a voucher to start.</div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <!-- MODAL 1: CREATE VOUCHER (Chỉ tạo mã, chưa chọn người) -->
        <div class="modal fade" ref="createModalRef" tabindex="-1">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header bg-primary text-white">
                        <h5 class="modal-title fw-bold">Generate Voucher Code</h5>
                        <button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
                    </div>
                    <div class="modal-body">
                        <div class="alert alert-info small border-0 bg-info bg-opacity-10">
                            <i class="fa-solid fa-circle-info me-2"></i>
                            Tạo mã Voucher trước, sau đó gán cho thủy thủ ở danh sách bên dưới.
                        </div>
                        <div class="mb-3">
                            <label class="form-label fw-bold">Data Plan</label>
                            <select v-model="createForm.data_plan" class="form-select">
                                <option>1GB</option><option>5GB</option><option>10GB</option><option>Unlimited</option>
                            </select>
                        </div>
                        <div class="mb-4">
                            <label class="form-label fw-bold">Validity (Days)</label>
                            <input type="number" v-model="createForm.valid_days" class="form-control">
                        </div>
                        <div class="text-end">
                            <button class="btn btn-light me-2" data-bs-dismiss="modal">Cancel</button>
                            <button class="btn btn-primary fw-bold px-4" @click="createVoucher">Generate</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- MODAL 2: ASSIGN TO CREW (Chọn người ở đây) -->
        <div class="modal fade" ref="assignModalRef" tabindex="-1">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header bg-success text-white">
                        <h5 class="modal-title fw-bold">Assign Voucher to Crew</h5>
                        <button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
                    </div>
                    <div class="modal-body">
                        <div class="mb-3">
                            <label class="form-label small fw-bold text-secondary">Step 1: Select Vessel</label>
                            <select v-model="assignForm.shipId" @change="onShipSelect" class="form-select">
                                <option value="" disabled>-- Choose a Ship --</option>
                                <option v-for="s in ships" :key="s.id" :value="s.id">{{ s.name }}</option>
                            </select>
                        </div>
                        <div class="mb-4">
                            <label class="form-label small fw-bold text-secondary">Step 2: Select Crew Member</label>
                            <select v-model="assignForm.crewId" class="form-select" :disabled="!crews.length">
                                <option value="" disabled>-- Choose Crew --</option>
                                <option v-for="c in crews" :key="c.id" :value="c.id">{{ c.full_name }} ({{ c.rank }})</option>
                            </select>
                            <div v-if="assignForm.shipId && crews.length === 0" class="form-text text-danger mt-2">
                                <i class="fa-solid fa-circle-exclamation me-1"></i> No crew found on this ship.
                            </div>
                        </div>
                        <div class="text-end">
                            <button class="btn btn-light me-2" data-bs-dismiss="modal">Cancel</button>
                            <button class="btn btn-success text-white fw-bold px-4" @click="assignVoucher" :disabled="!assignForm.crewId">Confirm Assign</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

    </div>
</template>

<style scoped>
.font-monospace { font-family: 'Consolas', monospace; letter-spacing: 1px; }
</style>