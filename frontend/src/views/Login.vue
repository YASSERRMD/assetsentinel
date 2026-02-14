<template>
  <div class="login-container">
    <div class="login-card">
      <h1>AssetSentinel</h1>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" required />
        </div>
        <div class="form-group">
          <label>Password</label>
          <input v-model="password" type="password" required />
        </div>
        <button type="submit">{{ loading ? 'Logging in...' : 'Login' }}</button>
        <p v-if="error" class="error">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { auth } from '../services/api'

const router = useRouter()
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  loading.value = true
  error.value = ''
  try {
    const { data } = await auth.login(email.value, password.value)
    localStorage.setItem('token', data.token)
    localStorage.setItem('user', JSON.stringify(data.user))
    router.push('/dashboard')
  } catch (err) {
    error.value = err.response?.data?.error || 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #667eea, #764ba2); }
.login-card { background: white; padding: 2rem; border-radius: 12px; width: 100%; max-width: 400px; }
h1 { text-align: center; color: #333; }
.form-group { margin-bottom: 1rem; }
label { display: block; margin-bottom: 0.5rem; }
input { width: 100%; padding: 0.75rem; border: 1px solid #ddd; border-radius: 6px; }
button { width: 100%; padding: 0.75rem; background: #667eea; color: white; border: none; border-radius: 6px; cursor: pointer; }
.error { color: red; text-align: center; margin-top: 1rem; }
</style>
