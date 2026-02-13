<template>
  <div class="layout">
    <aside class="sidebar">
      <h2>AssetSentinel</h2>
      <nav>
        <router-link to="/dashboard">Dashboard</router-link>
        <router-link to="/assets">Assets</router-link>
        <router-link to="/maintenance">Maintenance</router-link>
        <router-link to="/work-orders">Work Orders</router-link>
        <router-link to="/inventory">Inventory</router-link>
        <router-link to="/reports">Reports</router-link>
      </nav>
      <button @click="logout" class="logout-btn">Logout</button>
    </aside>
    <main class="content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { ws } from '../services/api'

const router = useRouter()

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}

ws.connect()
</script>

<style scoped>
.layout { display: flex; min-height: 100vh; }
.sidebar {
  width: 250px;
  background: #2c3e50;
  color: white;
  padding: 1rem;
}
.sidebar h2 { padding: 1rem 0; border-bottom: 1px solid #34495e; }
.sidebar nav { display: flex; flex-direction: column; gap: 0.5rem; margin-top: 1rem; }
.sidebar a {
  padding: 0.75rem 1rem;
  color: #bdc3c7;
  text-decoration: none;
  border-radius: 6px;
  transition: all 0.2s;
}
.sidebar a:hover, .sidebar a.router-link-active {
  background: #34495e;
  color: white;
}
.logout-btn {
  margin-top: auto;
  width: 100%;
  padding: 0.75rem;
  background: #e74c3c;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}
.content { flex: 1; padding: 2rem; background: #f5f6fa; }
</style>
