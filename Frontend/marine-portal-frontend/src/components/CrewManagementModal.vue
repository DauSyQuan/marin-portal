<script setup>
import { ref, onMounted } from 'vue';
import api from '@/services/api';

const props = defineProps({
  shipId: { type: [String, Number], required: true }
});

const emit = defineEmits(['close']);

const loading = ref(false);
const crew = ref([]);
const newCrew = ref({
  name: '',
  role: '',
  phone: '',
});

const fetchCrew = async () => {
  loading.value = true;
  try {
    const res = await api.get(`/api/ships/${props.shipId}/crew`);
    crew.value = Array.isArray(res.data) ? res.data : [];
  } catch (e) {
    console.error(e);
    crew.value = [];
  } finally {
    loading.value = false;
  }
};

const addCrew = async () => {
  if (!newCrew.value.name) return alert("Nhập tên thuyền viên");
  try {
    await api.post(`/api/ships/${props.shipId}/crew`, newCrew.value);
    newCrew.value = { name: '', role: '', phone: '' };
    await fetchCrew();
  } catch (e) {
    console.error(e);
    alert("Thêm crew lỗi: " + (e?.response?.data?.error || e.message));
  }
};

const removeCrew = async (id) => {
  if (!confirm("Xóa crew này?")) return;
  try {
    await api.delete(`/api/ships/${props.shipId}/crew/${id}`);
    await fetchCrew();
  } catch (e) {
    console.error(e);
    alert("Xóa crew lỗi: " + (e?.response?.data?.error || e.message));
  }
};

onMounted(fetchCrew);
</script>

<template>
  <div class="modal-backdrop">
    <div class="modal-card glass-panel">
      <div class="d-flex justify-content-between align-items-center mb-3">
        <h6 class="fw-bold text-white m-0">Crew Management</h6>
        <button class="btn btn-sm btn-outline-light" @click="emit('close')">
          <i class="fa-solid fa-xmark"></i>
        </button>
      </div>

      <!-- ADD CREW -->
      <div class="row g-2 mb-3">
        <div class="col-5">
          <input v-model="newCrew.name" class="form-control bg-dark text-white border-secondary" placeholder="Name">
        </div>
        <div class="col-4">
          <input v-model="newCrew.role" class="form-control bg-dark text-white border-secondary" placeholder="Role">
        </div>
        <div class="col-3">
          <button class="btn btn-primary w-100 fw-bold" @click="addCrew">
            <i class="fa-solid fa-plus"></i>
          </button>
        </div>
      </div>

      <!-- LIST CREW -->
      <div class="crew-list">
        <div v-if="loading" class="text-white-50 text-center py-3">Loading...</div>
        <div v-else-if="crew.length === 0" class="text-white-50 text-center py-3">No crew found.</div>

        <div v-else v-for="c in crew" :key="c.id" class="crew-row">
          <div>
            <div class="fw-bold text-white">{{ c.name }}</div>
            <small class="text-white-50">{{ c.role || 'N/A' }}</small>
          </div>
          <button class="btn btn-sm btn-outline-danger" @click="removeCrew(c.id)">
            <i class="fa-solid fa-trash"></i>
          </button>
        </div>
      </div>

      <div class="text-white-50 small mt-3">
        Sync crew sẽ đồng bộ user hotspot lên router.
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-backdrop{
  position: fixed;
  inset: 0;
  z-index: 999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0,0,0,0.5);
}
.modal-card{
  width: 520px;
  max-width: calc(100vw - 40px);
  padding: 18px;
  border-radius: 16px;
}
.crew-list{
  max-height: 320px;
  overflow: auto;
}
.crew-row{
  display:flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 12px;
  margin-bottom: 8px;
  background: rgba(0,0,0,0.2);
}
</style>
