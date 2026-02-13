<template>
  <Layout v-if="isAuthenticated" />
  <router-view v-else />
</template>

<script setup>
import { computed } from 'vue'
import Layout from './views/Layout.vue'

const isAuthenticated = computed(() => {
  return localStorage.getItem('token')
})
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
.page { padding: 1rem; }
</style>
