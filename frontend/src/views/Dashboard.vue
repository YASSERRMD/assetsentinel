<template>
  <div class="dashboard">
    <h1>Dashboard</h1>
    <div class="stats-grid">
      <div class="stat-card">
        <h3>Total Assets</h3>
        <p class="stat-value">{{ stats.asset_count || 0 }}</p>
      </div>
      <div class="stat-card warning">
        <h3>Overdue Maintenance</h3>
        <p class="stat-value">{{ stats.overdue_maintenance || 0 }}</p>
      </div>
      <div class="stat-card danger">
        <h3>Low Stock Items</h3>
        <p class="stat-value">{{ stats.low_stock || 0 }}</p>
      </div>
      <div class="stat-card">
        <h3>Total Costs</h3>
        <p class="stat-value">${{ (stats.total_costs || 0).toLocaleString() }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { dashboard } from '../services/api'

const stats = ref({})

const fetchStats = async () => {
  try {
    const { data } = await dashboard.stats()
    stats.value = data
  } catch (err) { console.error(err) }
}

onMounted(fetchStats)
</script>

<style scoped>
.dashboard h1 { margin-bottom: 2rem; }
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1.5rem; }
.stat-card { background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.stat-card h3 { margin: 0 0 0.5rem; color: #666; font-size: 0.9rem; }
.stat-value { margin: 0; font-size: 2rem; font-weight: bold; color: #2c3e50; }
.stat-card.warning .stat-value { color: #f39c12; }
.stat-card.danger .stat-value { color: #e74c3c; }
</style>
