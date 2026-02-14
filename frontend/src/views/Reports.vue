<template>
  <div class="page">
    <h1>Reports</h1>
    <div class="tabs">
      <button @click="activeTab = 'costs'" :class="{ active: activeTab === 'costs' }">Cost Summary</button>
      <button @click="activeTab = 'depreciation'" :class="{ active: activeTab === 'depreciation' }">Depreciation</button>
    </div>
    
    <div v-if="activeTab === 'costs'" class="tab-content">
      <h2>Asset Cost Summary</h2>
      <table class="data-table">
        <thead><tr><th>Asset ID</th><th>Asset Name</th><th>Total Maintenance Cost</th></tr></thead>
        <tbody>
          <tr v-for="item in costs" :key="item.asset_id">
            <td>{{ item.asset_id }}</td>
            <td>{{ item.asset_name }}</td>
            <td>${{ item.total_cost?.toLocaleString() }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <div v-if="activeTab === 'depreciation'" class="tab-content">
      <h2>Asset Depreciation</h2>
      <select v-model="selectedAsset" @change="fetchDepreciation">
        <option value="">Select Asset</option>
        <option v-for="asset in assets" :key="asset.id" :value="asset.id">{{ asset.name }}</option>
      </select>
      <table v-if="depreciation.length" class="data-table">
        <thead><tr><th>Year</th><th>Depreciation Amount</th></tr></thead>
        <tbody>
          <tr v-for="d in depreciation" :key="d.year">
            <td>{{ d.year }}</td>
            <td>${{ d.depreciation_amount?.toLocaleString() }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { reports, assets as assetsApi } from '../services/api'

const activeTab = ref('costs')
const costs = ref([])
const depreciation = ref([])
const selectedAsset = ref('')
const assets = ref([])

const fetchCosts = async () => { try { const { data } = await reports.costsAll(); costs.value = data } catch (err) { console.error(err) } }
const fetchAssets = async () => { try { const { data } = await assetsApi.list({ page_size: 100 }); assets.value = data.data } catch (err) { console.error(err) } }
const fetchDepreciation = async () => { 
  if (!selectedAsset.value) { depreciation.value = []; return }
  try { const { data } = await reports.depreciation(selectedAsset.value); depreciation.value = data } catch (err) { console.error(err) } 
}
onMounted(() => { fetchCosts(); fetchAssets() })
</script>

<style scoped>
.tabs { display: flex; gap: 0.5rem; margin-bottom: 1.5rem; }
.tabs button { padding: 0.75rem 1.5rem; border: none; background: white; cursor: pointer; border-radius: 6px; }
.tabs button.active { background: #667eea; color: white; }
.tab-content { background: white; padding: 1.5rem; border-radius: 8px; }
.tab-content select { padding: 0.5rem; margin-bottom: 1rem; border: 1px solid #ddd; border-radius: 4px; }
.data-table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
.data-table th, .data-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #eee; }
</style>
