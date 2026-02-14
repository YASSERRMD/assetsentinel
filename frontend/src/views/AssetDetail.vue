<template>
  <div class="page">
    <h1>Asset Details</h1>
    <div v-if="asset" class="detail-card">
      <div class="detail-row"><span>Name:</span><strong>{{ asset.name }}</strong></div>
      <div class="detail-row"><span>Category:</span><strong>{{ asset.category }}</strong></div>
      <div class="detail-row"><span>Serial Number:</span><strong>{{ asset.serial_number || '-' }}</strong></div>
      <div class="detail-row"><span>Location:</span><strong>{{ asset.location || '-' }}</strong></div>
      <div class="detail-row"><span>Purchase Cost:</span><strong>${{ asset.purchase_cost?.toLocaleString() }}</strong></div>
      <div class="detail-row"><span>Status:</span><span :class="`status ${asset.status}`">{{ asset.status }}</span></div>
      <div class="detail-row"><span>Installation Date:</span><strong>{{ asset.installation_date || '-' }}</strong></div>
      <div class="detail-row"><span>Warranty Expiry:</span><strong>{{ asset.warranty_expiry || '-' }}</strong></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { assets } from '../services/api'

const route = useRoute()
const asset = ref(null)

onMounted(async () => {
  try {
    const { data } = await assets.get(route.params.id)
    asset.value = data
  } catch (err) { console.error(err) }
})
</script>

<style scoped>
.detail-card { background: white; padding: 2rem; border-radius: 8px; margin-top: 1rem; }
.detail-row { display: flex; justify-content: space-between; padding: 1rem 0; border-bottom: 1px solid #eee; }
.status { padding: 0.25rem 0.5rem; border-radius: 4px; }
.status.active { background: #d4edda; color: #155724; }
</style>
