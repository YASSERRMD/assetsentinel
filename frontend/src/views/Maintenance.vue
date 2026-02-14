<template>
  <div class="page">
    <div class="header">
      <h1>Maintenance Plans</h1>
      <button @click="showForm = true" class="btn-primary">Add Plan</button>
    </div>
    <table class="data-table">
      <thead><tr><th>Asset ID</th><th>Frequency (Days)</th><th>Next Date</th><th>Assigned Role</th><th>Actions</th></tr></thead>
      <tbody>
        <tr v-for="plan in plans" :key="plan.id">
          <td>{{ plan.asset_id }}</td>
          <td>{{ plan.frequency_days }}</td>
          <td>{{ plan.next_maintenance_date }}</td>
          <td>{{ plan.assigned_role || '-' }}</td>
          <td>
            <button @click="editPlan(plan)" class="btn-sm">Edit</button>
            <button @click="deletePlan(plan.id)" class="btn-sm danger">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="showForm" class="modal">
      <div class="modal-content">
        <h2>{{ editingId ? 'Edit' : 'Add' }} Plan</h2>
        <form @submit.prevent="savePlan">
          <input v-model="form.asset_id" type="number" placeholder="Asset ID" required />
          <input v-model="form.frequency_days" type="number" placeholder="Frequency (days)" required />
          <input v-model="form.estimated_duration_hours" type="number" placeholder="Est. Duration (hours)" />
          <input v-model="form.next_maintenance_date" type="date" required />
          <select v-model="form.assigned_role"><option value="">Select Role</option><option value="technician">Technician</option><option value="maintenance_manager">Maintenance Manager</option></select>
          <div class="modal-actions"><button type="button" @click="closeForm">Cancel</button><button type="submit" class="btn-primary">Save</button></div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { maintenance as maintenanceApi } from '../services/api'

const plans = ref([])
const showForm = ref(false)
const editingId = ref(null)
const form = ref({ asset_id: null, frequency_days: 30, estimated_duration_hours: null, next_maintenance_date: '', assigned_role: '' })

const fetchPlans = async () => { try { const { data } = await maintenanceApi.list(); plans.value = data.data } catch (err) { console.error(err) } }
const savePlan = async () => { try { if (editingId.value) await maintenanceApi.update(editingId.value, form.value); else await maintenanceApi.create(form.value); closeForm(); fetchPlans() } catch (err) { console.error(err) } }
const editPlan = (plan) => { editingId.value = plan.id; form.value = { ...plan }; showForm.value = true }
const deletePlan = async (id) => { if (confirm('Delete this plan?')) { await maintenanceApi.delete(id); fetchPlans() } }
const closeForm = () => { showForm.value = false; editingId.value = null; form.value = { asset_id: null, frequency_days: 30, estimated_duration_hours: null, next_maintenance_date: '', assigned_role: '' } }
onMounted(fetchPlans)
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.btn-primary { padding: 0.5rem 1rem; background: #667eea; color: white; border: none; border-radius: 6px; cursor: pointer; }
.data-table { width: 100%; background: white; border-collapse: collapse; border-radius: 8px; }
.data-table th, .data-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #eee; }
.btn-sm { padding: 0.25rem 0.5rem; margin-right: 0.25rem; border: none; border-radius: 4px; cursor: pointer; }
.btn-sm.danger { background: #dc3545; color: white; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: white; padding: 2rem; border-radius: 8px; width: 400px; }
.modal-content form { display: flex; flex-direction: column; gap: 1rem; }
.modal-content input, .modal-content select { padding: 0.75rem; border: 1px solid #ddd; border-radius: 4px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 1rem; }
</style>
