<template>
  <div class="layout">
    <aside class="sidebar">
      <h2>AssetSentinel</h2>
      <nav>
        <router-link to="/dashboard">Dashboard</router-link>
        <router-link to="/assets">Assets</router-link>
      </nav>
      <button @click="logout" class="logout">Logout</button>
    </aside>
    <main class="content"><router-view /></main>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
const router = useRouter()
const logout = () => { localStorage.removeItem('token'); router.push('/login') }
</script>

<style scoped>
.layout { display: flex; min-height: 100vh; }
.sidebar { width: 250px; background: #2c3e50; color: white; padding: 1rem; }
.sidebar nav { display: flex; flex-direction: column; gap: 0.5rem; margin-top: 1rem; }
.sidebar a { padding: 0.75rem; color: #bdc3c7; text-decoration: none; border-radius: 6px; }
.sidebar a:hover, .sidebar a.router-link-active { background: #34495e; color: white; }
.logout { margin-top: auto; width: 100%; padding: 0.75rem; background: #e74c3c; color: white; border: none; border-radius: 6px; cursor: pointer; }
.content { flex: 1; padding: 2rem; background: #f5f6fa; }
</style>
