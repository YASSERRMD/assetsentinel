<template>
  <div class="page">
    <div class="header">
      <h1>Inventory</h1>
      <button @click="showForm = true" class="btn-primary">Add Part</button>
    </div>
    <table class="data-table">
      <thead><tr><th>Name</th><th>SKU</th><th>Quantity</th><th>Min Threshold</th><th>Cost/Unit</th><th>Location</th><th>Actions</th></tr></thead>
      <tbody>
        <tr v-for="part in inventory" :key="part.id" :class="{ 'low-stock': part.quantity <= part.min_threshold }">
          <td>{{ part.name }}</td>
          <td>{{ part.sku }}</td>
          <td>{{ part.quantity }}</td>
          <td>{{ part.min_threshold }}</td>
          <td>${{ part.cost_per_unit?.toFixed(2) }}</td>
          <td>{{ part.location || '-' }}</td>
          <td>
            <button @click="editPart(part)" class="btn-sm">Edit</button>
            <button @click="deletePart(part.id)" class="btn-sm danger">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="showForm" class="modal">
      <div class="modal-content">
        <h2>{{ editingId ? 'Edit' : 'Add' }} Part</h2>
        <form @submit.prevent="savePart">
          <input v-model="form.name" placeholder="Name" required />
          <input v-model="form.sku" placeholder="SKU" required />
          <input v-model="form.quantity" type="number" placeholder="Quantity" required />
          <input v-model="form.min_threshold" type="number" placeholder="Min Threshold" required />
          <input v-model="form.cost_per_unit" type="number" step="0.01" placeholder="Cost per Unit" required />
          <input v-model="form.location" placeholder="Location" />
          <div class="modal-actions"><button type="button" @click="closeForm">Cancel</button><button type="submit" class="btn-primary">Save</button></div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { inventory as inventoryApi } from '../services/api'

const inventory = ref([])
const showForm = ref(false)
const editingId = ref(null)
const form = ref({ name: '', sku: '', quantity: 0, min_threshold: 0, cost_per_unit: 0, location: '' })

const fetchInventory = async () => { try { const { data } = await inventoryApi.list(); inventory.value = data.data } catch (err) { console.error(err) } }
const savePart = async () => { try { if (editingId.value) await inventoryApi.update(editingId.value, form.value); else await inventoryApi.create(form.value); closeForm(); fetchInventory() } catch (err) { console.error(err) } }
const editPart = (part) => { editingId.value = part.id; form.value = { ...part }; showForm.value = true }
const deletePart = async (id) => { if (confirm('Delete this part?')) { await inventoryApi.delete(id); fetchInventory() } }
const closeForm = () => { showForm.value = false; editingId.value = null; form.value = { name: '', sku: '', quantity: 0, min_threshold: 0, cost_per_unit: 0, location: '' } }
onMounted(fetchInventory)
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.btn-primary { padding: 0.5rem 1rem; background: #667eea; color: white; border: none; border-radius: 6px; cursor: pointer; }
.data-table { width: 100%; background: white; border-collapse: collapse; border-radius: 8px; }
.data-table th, .data-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #eee; }
.data-table tr.low-stock { background: #fff3cd; }
.btn-sm { padding: 0.25rem 0.5rem; margin-right: 0.25rem; border: none; border-radius: 4px; cursor: pointer; }
.btn-sm.danger { background: #dc3545; color: white; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: white; padding: 2rem; border-radius: 8px; width: 400px; }
.modal-content form { display: flex; flex-direction: column; gap: 1rem; }
.modal-content input { padding: 0.75rem; border: 1px solid #ddd; border-radius: 4px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 1rem; }
</style>
