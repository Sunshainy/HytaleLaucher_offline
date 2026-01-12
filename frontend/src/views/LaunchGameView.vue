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
import NewsCarousel from '@/components/NewsCarousel.vue'

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
  READY_TO_PLAY: 'ready_to_play',
  VALIDATING: 'validating'
} as const

const updateInfo = computed(() => appStore.updateInfo)

const currentState = computed(() => {
  if (appStore.isValidating) return STATE.VALIDATING
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

const newsArticles = computed(() => {
  return appStore.feedArticles.map(article => ({
    id: article.id || String(Math.random()),
    title: article.title,
    description: article.description || '',
    imageUrl: article.image_url,
    link: article.dest_url
  }))
})

const gameVersionText = computed(() => {
  if (updateInfo.value?.GameVersion) {
    return `${t('launch_game.game_version')}: ${updateInfo.value.GameVersion}`
  }
  return ''
})

const installedVersionText = computed(() => {
  const version = appStore.gameVersion
  return version || 'Unknown'
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
  if (window.go?.main?.App?.LaunchGame) {
    window.go.main.App.LaunchGame()
  }
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

function handleNewsDetails(article: { link?: string }) {
  if (article.link && window.runtime?.BrowserOpenURL) {
    window.runtime.BrowserOpenURL(article.link)
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
          v-if="profileOptions.length > 0"
          :model-value="selectedProfile"
          @update:model-value="setProfile"
          :options="profileOptions"
          class="launch-game__profile-dropdown"
        />
      </div>

      <!-- News carousel -->
      <NewsCarousel
        :articles="newsArticles"
        class="launch-game__carousel"
        @details="handleNewsDetails"
      />

      <!-- Install Hytale section -->
      <div v-if="currentState === STATE.INSTALL_AVAILABLE" class="launch-game__install-hytale install-hytale">
        <HyButton
          type="primary"
          class="install-hytale__button"
          @click="install"
        >
          {{ updateInfo?.PrimaryAction || 'Install' }}
        </HyButton>
        <div v-if="gameVersionText" class="install-hytale__version-text-container">
          <span class="install-hytale__version-text">{{ gameVersionText }}</span>
        </div>
      </div>

      <!-- Play Hytale section -->
      <div v-if="currentState === STATE.READY_TO_PLAY" class="launch-game__play-hytale play-hytale">
        <HyButton
          type="primary"
          class="play-hytale__button"
          :disabled="isLaunching"
          @click="play"
        >
          {{ $t('common.play') }}
        </HyButton>
        <span class="play-hytale__version-text hytale-version">
          Version: {{ installedVersionText }}
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

      <!-- Validating label -->
      <label v-if="currentState === STATE.VALIDATING" class="launch-game__validating-label">
        {{ $t('update_status.validating_patch') }}
      </label>

      <!-- Offline not installed label -->
      <label v-if="currentState === STATE.OFFLINE_NOT_INSTALLED" class="launch-game__not-installed-label">
        {{ $t('error.offline_not_installed') }}
      </label>
    </div>
  </div>
</template>

<style scoped>
.launch-game__container {
  padding: 44px 44px 25px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
  height: 100%;
  width: 100%;
  position: relative;
}

.launch-game__logo :deep(img) {
  width: 287px;
}

.launch-game__top-content {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.launch-game__profile-dropdown {
  width: 180px;
  position: absolute;
  right: 14px;
  top: 55px;
}

.launch-game__carousel {
  margin-top: 24px;
}

.launch-game__install-hytale,
.launch-game__play-hytale,
.launch-game__installation-progress-bar,
.launch-game__cancelling-label,
.launch-game__not-installed-label,
.launch-game__validating-label {
  margin-top: 64px;
}

.launch-game__logo :deep(img) {
  height: 146px;
}

/* Install Hytale sub-component styles */
.install-hytale {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
  width: 100%;
}

.install-hytale__button {
  width: 220px;
}

.install-hytale__version-text-container {
  margin-top: 8px;
}

.install-hytale__version-text {
  margin-top: 8px;
  margin-right: 8px;
  color: rgba(210, 217, 226, 0.5);
  font-size: 14px;
}

/* Play Hytale sub-component styles */
.play-hytale {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
  width: 100%;
}

.play-hytale__button {
  width: 220px;
}

.play-hytale__version-text {
  margin-top: 8px;
}

.hytale-version {
  color: rgba(210, 217, 226, 0.5);
  font-size: 14px;
}

/* Labels */
.launch-game__cancelling-label,
.launch-game__not-installed-label,
.launch-game__validating-label {
  color: #8b949f;
  font-size: 14px;
  font-family: 'Nunito Sans', sans-serif;
}
</style>
