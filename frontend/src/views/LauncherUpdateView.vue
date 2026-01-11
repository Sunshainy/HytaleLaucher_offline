<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/appStore'
import Logo from '@/components/Logo.vue'
import InstallationProgressBar from '@/components/InstallationProgressBar.vue'

const router = useRouter()
const appStore = useAppStore()
const completed = ref(false)

onMounted(async () => {
  try {
    await appStore.applyUpdates()
    completed.value = true
  } catch (error) {
    router.push({ name: 'error', query: { error: String(error) } })
  }
})
</script>

<template>
  <div class="launcher-update">
    <div class="launcher-update__container">
      <div class="launcher-update__top-content">
        <Logo class="launcher-update__logo" />
      </div>
      <InstallationProgressBar
        v-if="appStore.updateRunning"
        class="launcher-update__installation-progress-bar"
      />
      <label v-if="!appStore.updateRunning && !completed" class="launcher-update__status-label">
        {{ $t('launcher_update.preparing') }}
      </label>
      <label v-if="!appStore.updateRunning && completed" class="launcher-update__status-label">
        {{ $t('launcher_update.restarting') }}
      </label>
    </div>
  </div>
</template>

<style scoped>
.launcher-update {
  height: 100%;
  width: 100%;
}

.launcher-update__container {
  padding: 44px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
  height: 100%;
  width: 100%;
  position: relative;
}

.launcher-update__top-content {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.launcher-update__installation-progress-bar,
.launcher-update__status-label {
  margin-top: auto;
}

.launcher-update__logo :deep(img) {
  height: 146px;
}

.launcher-update__status-label {
  color: #8b949f;
  font-size: 14px;
}
</style>
