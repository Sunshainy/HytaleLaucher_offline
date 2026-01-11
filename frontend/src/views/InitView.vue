<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/authStore'
import { useAppStore } from '@/stores/appStore'
import Logo from '@/components/Logo.vue'
import Spinner from '@/components/Spinner.vue'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()
const { t } = useI18n()

const message = ref(t('init.loading'))

// Placeholder functions for backend calls
async function checkEulaAccepted(): Promise<boolean> {
  // Would check EULA acceptance from backend
  return true
}

async function checkGameAvailable(): Promise<boolean> {
  // Would check game availability from backend
  return true
}

onMounted(async () => {
  try {
    // Check network mode
    if (!authStore.hasNetworkBeenChecked) {
      await authStore.checkNetworkMode(true, 'initial_network_check')
    }

    // Check for launcher updates (if online)
    if (!authStore.isOffline) {
      const hasLauncherUpdate = await appStore.checkForFreestandingLauncherUpdate()
      if (hasLauncherUpdate) {
        router.push({ name: 'launcher-update' })
        return
      }
    }

    // Load account data
    await authStore.load()
    message.value = t('init.fetching_account_data')

    // Check if logged in
    const isLoggedIn = await authStore.checkSessionInfo()
    if (!isLoggedIn) {
      setTimeout(() => {
        router.push({ name: 'login' })
      }, 1000)
      return
    }

    // Check EULA acceptance
    const eulaAccepted = await checkEulaAccepted()
    if (!eulaAccepted) {
      router.push({ name: 'eula' })
      return
    }

    // Check game availability
    const gameAvailable = await checkGameAvailable()
    if (!gameAvailable) {
      router.push({ name: 'game-unavailable' })
      return
    }

    // Fetch install info
    await appStore.fetchInstallInfo()

    // Handle offline mode without game installed
    if (authStore.isOffline && appStore.currentChannel === '') {
      router.push({ name: 'launch-game' })
      return
    }

    // Check for updates (if online)
    if (!authStore.isOffline) {
      message.value = t('init.checking_for_updates')
      await appStore.checkForUpdates()
    }

    router.push({ name: 'launch-game' })
  } catch (error) {
    router.push({ name: 'error', query: { error: String(error) } })
  }
})
</script>

<template>
  <div class="init-view">
    <Logo class="init-view__logo" />
    <Spinner :message="message" />
  </div>
</template>

<style scoped>
.init-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.init-view__logo {
  margin-bottom: 32px;
}

.init-view__logo :deep(img) {
  height: 228px;
}
</style>
