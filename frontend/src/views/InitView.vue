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
    // TODO: Temporarily disabled all checks - just go to launch-game
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
