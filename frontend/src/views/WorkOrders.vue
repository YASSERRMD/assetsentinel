<template>
  <div class="page">
    <div class="header">
      <h1>Work Orders</h1>
      <button @click="showForm = true" class="btn-primary">Create Work Order</button>
    </div>
    <div class="filters">
      <select v-model="filters.status" @change="fetchOrders">
        <option value="">All Status</option>
        <option value="pending">Pending</option>
        <option value="in_progress">In Progress</option>
        <option value="completed">Completed</option>
        <option value="closed">Closed</option>
      </select>
    </div>
    <table class="data-table">
      <thead><tr><th>Title</th><th>Asset ID</th><th>Status</th><th>Priority</th><th>Created</th><th>Actions</th></tr></thead>
      <tbody>
        <tr v-for="wo in workOrders" :key="wo.id">
          <td>{{ wo.title }}</td>
          <td>{{ wo.asset_id }}</td>
          <td><span :class="`status ${wo.status}`">{{ wo.status }}</span></td>
          <td><span :class="`priority ${wo.priority}`">{{ wo.priority }}</span></td>
          <td>{{ new Date(wo.created_at).toLocaleDateString() }}</td>
          <td>
            <button @click="editOrder(wo)" class="btn-sm">Edit</button>
            <button @click="deleteOrder(wo.id)" class="btn-sm danger">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="showForm" class="modal">
      <div class="modal-content">
        <h2>{{ editingId ? 'Edit' : 'Create' }} Work Order</h2>
        <form @submit.prevent="saveOrder">
          <input v-model="form.title" placeholder="Title" required />
          <input v-model="form.asset_id" type="number" placeholder="Asset ID" required />
          <input v-model="form.technician_id" type="number" placeholder="Technician ID" />
          <textarea v-model="form.description" placeholder="Description"></textarea>
          <select v-model="form.status"><option value="pending">Pending</option><option value="in_progress">In Progress</option><option value="completed">Completed</option><option value="closed">Closed</option></select>
          <select v-model="form.priority"><option value="low">Low</option><option value="medium">Medium</option><option value="high">High</option><option value="critical">Critical</option></select>
          <input v-model="form.scheduled_start" type="datetime-local" />
          <input v-model="form.scheduled_end" type="datetime-local" />
          <input v-model="form.total_cost" type="number" placeholder="Total Cost" />
          <textarea v-model="form.notes" placeholder="Notes"></textarea>
          <div class="modal-actions"><button type="button" @click="closeForm">Cancel</button><button type="submit" class="btn-primary">Save</button></div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { workOrders as workOrdersApi } from '../services/api'

const workOrders = ref([])
const showForm = ref(false)
const editingId = ref(null)
const filters = ref({ status: '' })
const form = ref({ title: '', asset_id: null, technician_id: null, description: '', status: 'pending', priority: 'medium', scheduled_start: '', scheduled_end: '', total_cost: 0, notes: '' })

const fetchOrders = async () => { try { const { data } = await workOrdersApi.list({ ...filters.value }); workOrders.value = data.data } catch (err) { console.error(err) } }
const saveOrder = async () => { try { if (editingId.value) await workOrdersApi.update(editingId.value, form.value); else await workOrdersApi.create(form.value); closeForm(); fetchOrders() } catch (err) { console.error(err) } }
const editOrder = (wo) => { editingId.value = wo.id; form.value = { ...wo }; showForm.value = true }
const deleteOrder = async (id) => { if (confirm('Delete this work order?')) { await workOrdersApi.delete(id); fetchOrders() } }
const closeForm = () => { showForm.value = false; editingId.value = null; form.value = { title: '', asset_id: null, technician_id: null, description: '', status: 'pending', priority: 'medium', scheduled_start: '', scheduled_end: '', total_cost: 0, notes: '' } }
onMounted(fetchOrders)
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.filters { margin-bottom: 1rem; }
.filters select { padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
.btn-primary { padding: 0.5rem 1rem; background: #667eea; color: white; border: none; border-radius: 6px; cursor: pointer; }
.data-table { width: 100%; background: white; border-collapse: collapse; border-radius: 8px; }
.data-table th, .data-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #eee; }
.status { padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.85rem; }
.status.pending { background: #fff3cd; color: #856404; }
.status.in_progress { background: #cce5ff; color: #004085; }
.status.completed { background: #d4edda; color: #155724; }
.status.closed { background: #e2e3e5; color: #383d41; }
.priority { padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.85rem; }
.priority.low { background: #d4edda; }
.priority.medium { background: #fff3cd; }
.priority.high { background: #f8d7da; }
.priority.critical { background: #dc3545; color: white; }
.btn-sm { padding: 0.25rem 0.5rem; margin-right: 0.25rem; border: none; border-radius: 4px; cursor: pointer; }
.btn-sm.danger { background: #dc3545; color: white; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: white; padding: 2rem; border-radius: 8px; width: 500px; max-height: 80vh; overflow-y: auto; }
.modal-content form { display: flex; flex-direction: column; gap: 0.75rem; }
.modal-content input, .modal-content select, .modal-content textarea { padding: 0.75rem; border: 1px solid #ddd; border-radius: 4px; }
.modal-content textarea { min-height: 80px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem; }
</style>
