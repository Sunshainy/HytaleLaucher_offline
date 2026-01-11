<script lang="ts" setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/appStore'
import { useAuthStore } from '@/stores/authStore'
import { useI18n } from 'vue-i18n'
import Logo from '@/components/Logo.vue'
import HyDropdown from '@/components/HyDropdown.vue'
import HyButton from '@/components/HyButton.vue'
import InstallationProgressBar from '@/components/InstallationProgressBar.vue'

const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()
const { t } = useI18n()

// State enum
const STATE = {
  CANCELLING: 'cancelling',
  INSTALLING: 'installing',
  INSTALL_AVAILABLE: 'install_available',
  OFFLINE_NOT_INSTALLED: 'offline_not_installed',
  READY_TO_PLAY: 'ready_to_play'
} as const

const updateInfo = computed(() => appStore.updateInfo)

const currentState = computed(() => {
  if (appStore.isCancellingUpdate) return STATE.CANCELLING
  if (appStore.updateRunning) return STATE.INSTALLING
  if (updateInfo.value !== null) return STATE.INSTALL_AVAILABLE
  if (updateInfo.value === null && appStore.currentChannel) {
    const gameVersion = appStore.gameVersion
    if (authStore.isOffline && !gameVersion) {
      return STATE.OFFLINE_NOT_INSTALLED
    }
    return STATE.READY_TO_PLAY
  }
  return STATE.OFFLINE_NOT_INSTALLED
})

const profileOptions = computed(() => {
  return authStore.getUserProfiles.map(p => ({
    label: p.username,
    value: p.uuid
  }))
})

const selectedProfile = computed(() => {
  const profile = authStore.getUserProfile
  return profile ? profile.uuid : ''
})

const cancellationMessage = computed(() => {
  const status = appStore.cancellationStatus
  return status.id ? t(status.id, status.params || {}) : t('update_status.cancelling_updates')
})

// Launch composable
const isLaunching = computed(() => false) // Would track launching state

async function install() {
  try {
    await appStore.applyUpdates()
  } catch (error) {
    router.push({ name: 'error', query: { error: String(error) } })
  }
}

function play() {
  // Would call backend to launch game
  console.log('Launching game...')
}

async function setProfile(uuid: string | number) {
  try {
    await authStore.setUserProfile(String(uuid))
  } catch (error) {
    console.error(`Failed to set user profile: ${error}`)
    router.push({ name: 'error', query: { error: String(error) } })
  }
}

async function cancelUpdates() {
  try {
    await appStore.cancelUpdates()
  } catch (error) {
    console.error(`Failed to cancel updates: ${error}`)
    router.push({ name: 'error', query: { error: String(error) } })
  }
}

onMounted(async () => {
  await appStore.fetchNewsFeed()
})
</script>

<template>
  <div class="launch-game">
    <div class="launch-game__container">
      <div class="launch-game__top-content">
        <Logo class="launch-game__logo" />
        <HyDropdown
          :model-value="selectedProfile"
          @update:model-value="setProfile"
          :options="profileOptions"
          class="launch-game__profile-dropdown"
        />
      </div>

      <!-- News carousel placeholder -->
      <div class="launch-game__carousel">
        <div v-for="article in appStore.feedArticles" :key="article.id" class="launch-game__article">
          {{ article.title }}
        </div>
      </div>

      <!-- Install button -->
      <div v-if="currentState === STATE.INSTALL_AVAILABLE" class="launch-game__install-hytale">
        <HyButton type="primary" @click="install">
          {{ updateInfo?.PrimaryAction || 'Install' }}
        </HyButton>
        <span v-if="updateInfo?.GameVersion" class="launch-game__version-text">
          {{ $t('launch_game.game_version') }}: {{ updateInfo.GameVersion }}
        </span>
      </div>

      <!-- Play button -->
      <div v-if="currentState === STATE.READY_TO_PLAY" class="launch-game__play-hytale">
        <HyButton type="primary" @click="play" :disabled="isLaunching">
          {{ $t('common.play') }}
        </HyButton>
        <span class="launch-game__version-text">
          Version: {{ appStore.gameVersion || 'Unknown' }}
        </span>
      </div>

      <!-- Installation progress -->
      <InstallationProgressBar
        v-if="currentState === STATE.INSTALLING"
        class="launch-game__installation-progress-bar"
        :on-cancel="cancelUpdates"
      />

      <!-- Cancelling label -->
      <label v-if="currentState === STATE.CANCELLING" class="launch-game__cancelling-label">
        {{ cancellationMessage }}
      </label>

      <!-- Offline not installed label -->
      <label v-if="currentState === STATE.OFFLINE_NOT_INSTALLED" class="launch-game__not-installed-label">
        {{ $t('error.offline_not_installed') }}
      </label>
    </div>
  </div>
</template>

<style scoped>
.launch-game {
  height: 100%;
  width: 100%;
}

.launch-game__container {
  padding: 44px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
  height: 100%;
  width: 100%;
  position: relative;
}

.launch-game__top-content {
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  justify-content: space-between;
  width: 100%;
}

.launch-game__logo :deep(img) {
  height: 146px;
}

.launch-game__profile-dropdown {
  width: 200px;
}

.launch-game__carousel {
  flex: 1;
  width: 100%;
  margin: 24px 0;
}

.launch-game__article {
  color: #d2d9e2;
  padding: 16px;
  background: rgba(22, 33, 47, 0.5);
  border-radius: 4px;
  margin-bottom: 8px;
}

.launch-game__install-hytale,
.launch-game__play-hytale {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  margin-top: auto;
}

.launch-game__version-text {
  color: #8b949f;
  font-size: 12px;
}

.launch-game__installation-progress-bar {
  margin-top: auto;
  width: 100%;
}

.launch-game__cancelling-label,
.launch-game__not-installed-label {
  margin-top: auto;
  color: #8b949f;
  font-size: 14px;
}
</style>
