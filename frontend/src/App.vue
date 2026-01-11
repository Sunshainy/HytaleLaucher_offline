<script lang="ts" setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
import WindowHeader from '@/components/WindowHeader.vue'
import CurseBackground from '@/components/CurseBackground.vue'
import NotificationContainer from '@/components/NotificationContainer.vue'
import OfflineFooter from '@/components/OfflineFooter.vue'
import DiscordLink from '@/components/DiscordLink.vue'
import HypixelStudiosLink from '@/components/HypixelStudiosLink.vue'

const route = useRoute()
const authStore = useAuthStore()

const showFooter = computed(() => {
  const hiddenRoutes = ['login', 'init', 'eula', 'error', 'validation-error']
  return !hiddenRoutes.includes(route.name as string)
})

const showGradient = computed(() => {
  return route.name === 'launch-game'
})
</script>

<template>
  <div class="app-container">
    <CurseBackground :show-gradient="showGradient" />
    <WindowHeader />
    <NotificationContainer />
    <div class="app-content">
      <router-view />
    </div>
    <div v-if="showFooter" class="app-footer">
      <DiscordLink />
      <HypixelStudiosLink />
    </div>
    <OfflineFooter v-if="authStore.isOffline" />
  </div>
</template>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  position: relative;
}

.app-content {
  flex: 1;
  overflow: auto;
  position: relative;
  z-index: 1;
}

.app-footer {
  position: absolute;
  bottom: 25px;
  right: 25px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 14px;
  z-index: 1;
}
</style>
