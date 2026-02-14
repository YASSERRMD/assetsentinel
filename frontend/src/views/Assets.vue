<template>
  <div class="assets-page">
    <div class="header">
      <h1>Assets</h1>
      <button @click="showForm = true" class="btn-primary">Add Asset</button>
    </div>
    <table class="data-table">
      <thead><tr><th>Name</th><th>Category</th><th>Status</th><th>Actions</th></tr></thead>
      <tbody>
        <tr v-for="asset in assets" :key="asset.id">
          <td>{{ asset.name }}</td>
          <td>{{ asset.category }}</td>
          <td><span :class="`status ${asset.status}`">{{ asset.status }}</span></td>
          <td><button @click="deleteAsset(asset.id)" class="btn-sm danger">Delete</button></td>
        </tr>
      </tbody>
    </table>
    <div v-if="showForm" class="modal">
      <div class="modal-content">
        <h2>Add Asset</h2>
        <form @submit.prevent="saveAsset">
          <input v-model="form.name" placeholder="Name" required />
          <input v-model="form.category" placeholder="Category" required />
          <input v-model.number="form.purchase_cost" type="number" placeholder="Purchase Cost" />
          <select v-model="form.status"><option value="active">Active</option><option value="under_maintenance">Under Maintenance</option><option value="retired">Retired</option></select>
          <div class="modal-actions">
            <button type="button" @click="showForm = false">Cancel</button>
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
const showForm = ref(false)
const form = ref({ name: '', category: '', purchase_cost: 0, status: 'active' })

const fetchAssets = async () => {
  try { const { data } = await assetsApi.list(); assets.value = data.data } 
  catch (err) { console.error(err) }
}

const saveAsset = async () => {
  try { await assetsApi.create(form.value); showForm.value = false; form.value = { name: '', category: '', purchase_cost: 0, status: 'active' }; fetchAssets() } 
  catch (err) { console.error(err) }
}

const deleteAsset = async (id) => { if (confirm('Delete?')) { await assetsApi.delete(id); fetchAssets() } }

onMounted(fetchAssets)
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.btn-primary { padding: 0.5rem 1rem; background: #667eea; color: white; border: none; border-radius: 6px; cursor: pointer; }
.data-table { width: 100%; background: white; border-collapse: collapse; border-radius: 8px; }
.data-table th, .data-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #eee; }
.status { padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.85rem; }
.status.active { background: #d4edda; color: #155724; }
.btn-sm { padding: 0.25rem 0.5rem; border: none; border-radius: 4px; cursor: pointer; }
.btn-sm.danger { background: #dc3545; color: white; }
.modal { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: white; padding: 2rem; border-radius: 8px; width: 400px; }
.modal-content form { display: flex; flex-direction: column; gap: 1rem; }
.modal-content input, .modal-content select { padding: 0.75rem; border: 1px solid #ddd; border-radius: 4px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 1rem; }
</style>
