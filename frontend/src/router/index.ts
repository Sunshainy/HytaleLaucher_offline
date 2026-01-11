import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/init'
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginView.vue')
  },
  {
    path: '/init',
    name: 'init',
    component: () => import('@/views/InitView.vue')
  },
  {
    path: '/eula',
    name: 'eula',
    component: () => import('@/views/EULAView.vue')
  },
  {
    path: '/launch-game',
    name: 'launch-game',
    component: () => import('@/views/LaunchGameView.vue')
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/views/SettingsView.vue')
  },
  {
    path: '/error',
    name: 'error',
    component: () => import('@/views/ErrorView.vue')
  },
  {
    path: '/validation-error',
    name: 'validation-error',
    component: () => import('@/views/ValidationErrorView.vue')
  },
  {
    path: '/launcher-update',
    name: 'launcher-update',
    component: () => import('@/views/LauncherUpdateView.vue')
  },
  {
    path: '/uninstall',
    name: 'uninstall',
    component: () => import('@/views/UninstallView.vue')
  },
  {
    path: '/game-unavailable',
    name: 'game-unavailable',
    component: () => import('@/views/GameUnavailableView.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
