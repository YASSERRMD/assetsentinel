import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Dashboard from '../views/Dashboard.vue'
import Assets from '../views/Assets.vue'
import AssetDetail from '../views/AssetDetail.vue'
import Maintenance from '../views/Maintenance.vue'
import WorkOrders from '../views/WorkOrders.vue'
import Inventory from '../views/Inventory.vue'
import Reports from '../views/Reports.vue'

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', name: 'Login', component: Login },
  { path: '/dashboard', name: 'Dashboard', component: Dashboard, meta: { requiresAuth: true } },
  { path: '/assets', name: 'Assets', component: Assets, meta: { requiresAuth: true } },
  { path: '/assets/:id', name: 'AssetDetail', component: AssetDetail, meta: { requiresAuth: true } },
  { path: '/maintenance', name: 'Maintenance', component: Maintenance, meta: { requiresAuth: true } },
  { path: '/work-orders', name: 'WorkOrders', component: WorkOrders, meta: { requiresAuth: true } },
  { path: '/inventory', name: 'Inventory', component: Inventory, meta: { requiresAuth: true } },
  { path: '/reports', name: 'Reports', component: Reports, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
