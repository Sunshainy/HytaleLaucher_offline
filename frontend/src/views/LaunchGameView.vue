<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/appStore'
import { useAuthStore } from '@/stores/authStore'
import { useNotificationStore } from '@/stores/notificationStore'
import { useI18n } from 'vue-i18n'
import Logo from '@/components/Logo.vue'
import HyDropdown from '@/components/HyDropdown.vue'
import HyButton from '@/components/HyButton.vue'
import InstallationProgressBar from '@/components/InstallationProgressBar.vue'
import NewsCarousel from '@/components/NewsCarousel.vue'
import { LaunchGame, GetPlayerName, IsGameInstalled, InstallGame } from '@wailsjs/go/app/App'
import { EventsOn } from '@wailsjs/runtime/runtime'

const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const { t } = useI18n()

// Player nickname
const playerNickname = ref('')
const isLaunching = ref(false)
const isInstalling = ref(false)
const installProgress = ref(0)
const gameInstalled = ref(false)

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
  // TODO: Temporarily disabled all checks
  if (appStore.isValidating) return STATE.VALIDATING
  if (appStore.isCancellingUpdate) return STATE.CANCELLING
  if (appStore.updateRunning) return STATE.INSTALLING
  if (updateInfo.value !== null) return STATE.INSTALL_AVAILABLE
  
  // Skip offline check - always show as ready to play
  return STATE.READY_TO_PLAY
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

async function install() {
  try {
    await appStore.applyUpdates()
  } catch (error) {
    router.push({ name: 'error', query: { error: String(error) } })
  }
}

async function play() {
  // Check if nickname is entered
  if (!playerNickname.value.trim()) {
    notificationStore.showError('Введите ник')
    return
  }

  // Check if game is installed
  try {
    const installed = await IsGameInstalled()
    if (!installed) {
      // Start installation
      notificationStore.showInfo('Начинается установка игры...')
      await installGame()
      return
    }
  } catch (error) {
    console.error('Failed to check game installation:', error)
  }

  // Launch the game
  try {
    isLaunching.value = true
    await LaunchGame({ playerName: playerNickname.value.trim() })
    notificationStore.showSuccess('Игра запущена!')
  } catch (error) {
    notificationStore.showError(String(error))
    console.error('Failed to launch game:', error)
  } finally {
    isLaunching.value = false
  }
}

async function installGame() {
  try {
    isInstalling.value = true
    installProgress.value = 0
    
    await InstallGame()
    
    gameInstalled.value = true
    notificationStore.showSuccess('Игра установлена!')
    
    // Auto-launch after installation
    try {
      isLaunching.value = true
      await LaunchGame({ playerName: playerNickname.value.trim() })
      notificationStore.showSuccess('Игра запущена!')
    } catch (error) {
      notificationStore.showError(String(error))
      console.error('Failed to launch game after install:', error)
    } finally {
      isLaunching.value = false
    }
  } catch (error) {
    notificationStore.showError('Ошибка установки: ' + String(error))
    console.error('Failed to install game:', error)
  } finally {
    isInstalling.value = false
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
  // Load saved player name
  try {
    const savedName = await GetPlayerName()
    if (savedName) {
      playerNickname.value = savedName
    }
  } catch (error) {
    console.log('No saved player name found')
  }

  // Check if game is installed
  try {
    gameInstalled.value = await IsGameInstalled()
  } catch (error) {
    console.error('Failed to check game installation:', error)
  }

  // Listen for installation progress events
  EventsOn('install:progress', (data: any) => {
    if (data.progress) {
      installProgress.value = Math.round(data.progress)
    }
  })

  EventsOn('install:complete', () => {
    gameInstalled.value = true
  })

  // TODO: Temporarily disabled
  // await appStore.fetchNewsFeed()
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
          :disabled="isLaunching || isInstalling"
          @click="play"
        >
          {{ isInstalling ? 'Установка...' : 'Играть' }}
        </HyButton>
        
        <!-- Installation progress -->
        <div v-if="isInstalling" class="play-hytale__install-progress">
          <span class="play-hytale__progress-text">{{ installProgress }}%</span>
        </div>
        
        <input
          v-model="playerNickname"
          type="text"
          placeholder="Введите ник"
          class="play-hytale__nickname-input"
          :disabled="isInstalling"
        />
        <span class="play-hytale__version-text hytale-version">
          Version: latest
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
  position: absolute;
  left: 90px;
  top: 235px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.launch-game__logo :deep(img) {
  height: 146px;
}

/* Install Hytale sub-component styles */
.install-hytale {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: auto;
}

.install-hytale__button {
  width: 220px;
}

.install-hytale__version-text-container {
  margin-top: 8px;
}

.install-hytale__version-text {
  margin-top: 8px;
  margin-right: 0;
  color: rgba(210, 217, 226, 0.5);
  font-size: 14px;
  text-align: center;
}

/* Play Hytale sub-component styles */
.play-hytale {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: auto;
}

.play-hytale__button {
  width: 220px;
}

.play-hytale__version-text {
  margin-top: 8px;
  text-align: center;
}

.play-hytale__nickname-input {
  width: 220px;
  padding: 8px 12px;
  margin-top: 12px;
  margin-bottom: 12px;
  background: rgba(116, 123, 132, 0.1);
  border: 1px solid rgba(116, 123, 132, 0.3);
  border-radius: 4px;
  color: rgba(210, 217, 226, 0.9);
  font-size: 14px;
  text-align: center;
  font-family: 'Nunito Sans', sans-serif;
  transition: all 0.2s ease;
}

.play-hytale__nickname-input:focus {
  outline: none;
  border-color: rgba(116, 123, 132, 0.6);
  background: rgba(116, 123, 132, 0.2);
  box-shadow: 0 0 8px rgba(116, 123, 132, 0.3);
}

.play-hytale__nickname-input::placeholder {
  color: rgba(210, 217, 226, 0.4);
}

.play-hytale__install-progress {
  margin-top: 12px;
  margin-bottom: 12px;
  text-align: center;
}

.play-hytale__progress-text {
  color: rgba(210, 217, 226, 0.9);
  font-size: 16px;
  font-weight: 600;
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
