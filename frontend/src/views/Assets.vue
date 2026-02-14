<template>
  <div class="assets-page">
    <div class="header">
      <h1>Assets</h1>
      <button @click="showForm = true" class="btn-primary">Add Asset</button>
    </div>
    
    <div class="filters">
      <select v-model="filters.status" @change="fetchAssets">
        <option value="">All Status</option>
        <option value="active">Active</option>
        <option value="under_maintenance">Under Maintenance</option>
        <option value="retired">Retired</option>
      </select>
      <input v-model="filters.category" placeholder="Category" @change="fetchAssets" />
    </div>

    <table class="data-table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Category</th>
          <th>Serial Number</th>
          <th>Location</th>
          <th>Status</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="asset in assets" :key="asset.id">
          <td><router-link :to="`/assets/${asset.id}`">{{ asset.name }}</router-link></td>
          <td>{{ asset.category }}</td>
          <td>{{ asset.serial_number || '-' }}</td>
          <td>{{ asset.location || '-' }}</td>
          <td><span :class="`status ${asset.status}`">{{ asset.status }}</span></td>
          <td>
            <button @click="editAsset(asset)" class="btn-sm">Edit</button>
            <button @click="deleteAsset(asset.id)" class="btn-sm danger">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div class="pagination">
      <button @click="changePage(-1)" :disabled="page === 1">Prev</button>
      <span>Page {{ page }}</span>
      <button @click="changePage(1)" :disabled="assets.length < pageSize">Next</button>
    </div>

    <div v-if="showForm" class="modal">
      <div class="modal-content">
        <h2>{{ editingId ? 'Edit' : 'Add' }} Asset</h2>
        <form @submit.prevent="saveAsset">
          <input v-model="form.name" placeholder="Name" required />
          <input v-model="form.category" placeholder="Category" required />
          <input v-model="form.serial_number" placeholder="Serial Number" />
          <input v-model="form.location" placeholder="Location" />
          <input v-model="form.purchase_cost" type="number" placeholder="Purchase Cost" />
          <input v-model="form.installation_date" type="date" placeholder="Installation Date" />
          <input v-model="form.warranty_expiry" type="date" placeholder="Warranty Expiry" />
          <select v-model="form.status">
            <option value="active">Active</option>
            <option value="under_maintenance">Under Maintenance</option>
            <option value="retired">Retired</option>
          </select>
          <div class="modal-actions">
            <button type="button" @click="closeForm">Cancel</button>
            <button type="submit" class="btn-primary">Save</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { assets as assetsApi } from '../services/api'

const assets = ref([])
const page = ref(1)
const pageSize = ref(10)
const filters = ref({ status: '', category: '' })
const showForm = ref(false)
const editingId = ref(null)
const form = ref({ name: '', category: '', serial_number: '', location: '', purchase_cost: 0, installation_date: '', warranty_expiry: '', status: 'active' })

const fetchAssets = async () => {
  try {
    const { data } = await assetsApi.list({ page: page.value, page_size: pageSize.value, ...filters.value })
    assets.value = data.data
  } catch (err) { console.error(err) }
}

const saveAsset = async () => {
  try {
    if (editingId.value) {
      await assetsApi.update(editingId.value, form.value)
    } else {
      await assetsApi.create(form.value)
    }
    closeForm()
    fetchAssets()
  } catch (err) { console.error(err) }
}

const editAsset = (asset) => {
  editingId.value = asset.id
  form.value = { ...asset }
  showForm.value = true
}

const deleteAsset = async (id) => {
  if (confirm('Delete this asset?')) {
    await assetsApi.delete(id)
    fetchAssets()
  }
}

const closeForm = () => { showForm.value = false; editingId.value = null; form.value = { name: '', category: '', serial_number: '', location: '', purchase_cost: 0, installation_date: '', warranty_expiry: '', status: 'active' } }
const changePage = (delta) => { page.value += delta; fetchAssets() }
onMounted(fetchAssets)
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.filters { display: flex; gap: 1rem; margin-bottom: 1rem; }
.filters input, .filters select { padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
.data-table { width: 100%; background: white; border-collapse: collapse; border-radius: 8px; overflow: hidden; }
.data-table th, .data-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #eee; }
.data-table th { background: #f8f9fa; font-weight: 600; }
.status { padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.85rem; }
.status.active { background: #d4edda; color: #155724; }
.status.under_maintenance { background: #fff3cd; color: #856404; }
.status.retired { background: #f8d7da; color: #721c24; }
.btn-sm { padding: 0.25rem 0.5rem; margin-right: 0.25rem; border: none; border-radius: 4px; cursor: pointer; }
.btn-sm.danger { background: #dc3545; color: white; }
.btn-primary { padding: 0.5rem 1rem; background: #667eea; color: white; border: none; border-radius: 6px; cursor: pointer; }
.pagination { display: flex; justify-content: center; align-items: center; gap: 1rem; margin-top: 1rem; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: white; padding: 2rem; border-radius: 8px; width: 500px; max-height: 80vh; overflow-y: auto; }
.modal-content form { display: flex; flex-direction: column; gap: 1rem; }
.modal-content input, .modal-content select { padding: 0.75rem; border: 1px solid #ddd; border-radius: 4px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem; }
</style>
