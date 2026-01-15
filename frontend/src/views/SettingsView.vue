<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/appStore'
import { useAuthStore } from '@/stores/authStore'
import { useNotificationStore } from '@/stores/notificationStore'
import PanelView from '@/components/PanelView.vue'
import HyDropdown from '@/components/HyDropdown.vue'
import HyButton from '@/components/HyButton.vue'
import LauncherVersion from '@/components/LauncherVersion.vue'
import { OpenGameDirectory, StartServer, StopServer, IsServerRunning } from '@wailsjs/go/app/App'
import { EventsOn } from '@wailsjs/runtime/runtime'

const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const { t } = useI18n()

const isCheckingUpdates = ref(false)
const serverRunning = ref(false)
const serverStarting = ref(false)
const isTogglingServer = ref(false)

const openInIcon = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABEAAAARCAYAAAA7bUf6AAAACXBIWXMAAAsTAAALEwEAmpwYAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAC8SURBVHgBrZK9EQIhEIX5Cwgp4SJmCCnBCmzFDjxLsAM7sQTMSC3hKgB3Ax08gT3v7iXsDPDxeLs8xjiybz2dczcscE8IcWaERG8TYGNK6cIIqfJCCwSOWM9R10kJyjlfN0HAycA5P66GIAC+codyWAWpATDoedjqX8C7AWXYTSdSSgOLqQFQZfubEGvtA8I8QDnNAT8gnMrK1H4UQjCMkKIOeO8n6gES0pPW+oTromGjtAuE90Jdql2cvAClzFdGDZMFsAAAAABJRU5ErkJggg=='

const allowedChannels = computed(() => appStore.allowedChannels)
const currentChannel = computed(() => appStore.currentChannel)

const checkingDisabled = computed(() => {
  return isCheckingUpdates.value || authStore.isOffline || appStore.updateRunning
})

const serverButtonDisabled = computed(() => {
  return isTogglingServer.value || appStore.updateRunning || serverStarting.value
})

const serverButtonText = computed(() => {
  if (serverStarting.value) {
    return 'Сервер запускается...'
  }
  if (isTogglingServer.value) {
    return serverRunning.value ? 'Останавливается...' : 'Запускается...'
  }
  return serverRunning.value ? 'Остановить сервер' : 'Запустить сервер'
})

const serverStatusText = computed(() => {
  if (serverStarting.value) {
    return 'Сервер запускается...'
  }
  return serverRunning.value ? 'Сервер запущен' : 'Сервер остановлен'
})

const channelOptions = computed(() => {
  return allowedChannels.value.map(c => ({ label: c, value: c }))
})

const canPerformActions = computed(() => {
  return !appStore.updateRunning
})

async function setChannel(channel: string | number) {
  try {
    await appStore.setChannel(String(channel))
  } catch (error) {
    router.push({ name: 'error', query: { error: String(error) } })
  }
}

async function openDirectory() {
  try {
    await OpenGameDirectory()
  } catch (error) {
    console.error('Failed to open directory:', error)
    notificationStore.showError(t('settings.failed_to_open_directory'))
  }
}

async function checkForUpdates() {
  isCheckingUpdates.value = true
  notificationStore.showInfo(t('settings.checking_for_updates'))

  try {
    const hasLauncherUpdate = await appStore.checkForFreestandingLauncherUpdate()
    if (hasLauncherUpdate) {
      notificationStore.showSuccess(t('settings.new_updates_available'))
      router.push({ name: 'launcher-update' })
      return
    }

    await appStore.checkForUpdates(true)

    if (appStore.updateInfo) {
      notificationStore.showSuccess(t('settings.new_updates_available'))
    } else {
      notificationStore.showInfo(t('settings.no_updates_available'))
    }
  } catch (error) {
    router.push({ name: 'error', query: { error: String(error) } })
    notificationStore.showError(t('settings.failed_to_check_for_updates'))
  } finally {
    isCheckingUpdates.value = false
  }
}

async function logout() {
  await authStore.logout()
  router.push({ name: 'login' })
}

async function toggleServer() {
  isTogglingServer.value = true
  try {
    if (serverRunning.value) {
      await StopServer()
      serverRunning.value = false
      serverStarting.value = false
    } else {
      serverStarting.value = true
      await StartServer()
    }
  } catch (error) {
    console.error('Failed to toggle server:', error)
    notificationStore.showError(`Ошибка: ${error}`)
    serverStarting.value = false
  } finally {
    isTogglingServer.value = false
  }
}

async function checkServerStatus() {
  try {
    serverRunning.value = await IsServerRunning()
  } catch (error) {
    console.error('Failed to check server status:', error)
  }
}

function close() {
  router.back()
}

function openUninstall() {
  router.push({ name: 'uninstall' })
}

onMounted(async () => {
  await appStore.fetchChannels()
  await checkServerStatus()
  
  // Listen for server events
  EventsOn('server:starting', () => {
    serverStarting.value = true
    serverRunning.value = false
  })
  
  EventsOn('server:ready', () => {
    serverStarting.value = false
    serverRunning.value = true
  })
  
  EventsOn('server:stopped', () => {
    serverRunning.value = false
    serverStarting.value = false
  })
  
  EventsOn('server:boot_timeout', () => {
    serverStarting.value = false
    notificationStore.showError('⚠️ Сервер не ответил. Проверьте логи.')
  })
})
</script>

<template>
  <PanelView
    :title="$t('settings.title')"
    :show-close-button="true"
    :show-report-bug="true"
    :esc-handler="close"
    @close="close"
  >
    <div class="settings__section">
      <h2 class="settings__label">{{ $t('settings.patchline') }}</h2>
      <HyDropdown
        :model-value="currentChannel"
        @update:model-value="setChannel"
        :options="channelOptions"
        class="settings__dropdown"
        :disabled="!canPerformActions"
      />
    </div>

    <div class="settings__section">
      <h2 class="settings__label">{{ $t('settings.directory') }}</h2>
      <HyButton small type="tertiary" class="settings__directory-button" @click="openDirectory">
        <span>{{ $t('settings.open_directory') }}</span>
        <img :src="openInIcon" :alt="$t('settings.open')" class="settings__directory-icon" draggable="false" />
      </HyButton>
    </div>

    <div class="settings__section">
      <h2 class="settings__label">{{ $t('settings.launcher_version') }}</h2>
      <LauncherVersion class="settings__version" />
    </div>

    <div class="settings__section">
      <h2 class="settings__label">
        Сетевая игра
        <span class="settings__version settings__server-status" :class="{ 
          'settings__server-status--active': serverRunning,
          'settings__server-status--starting': serverStarting 
        }">
          ({{ serverStatusText }})
        </span>
      </h2>
      <div class="settings__server-controls">
        <HyButton
          class="settings__action-button"
          :type="serverRunning ? 'primary' : 'tertiary'"
          @click="toggleServer"
          :disabled="serverButtonDisabled"
        >
          {{ serverButtonText }}
        </HyButton>
        <div v-if="serverStarting" class="settings__server-spinner">
          <div class="spinner-small"></div>
          <span class="settings__spinner-text">Загрузка сервера...</span>
        </div>
      </div>
    </div>

    <div class="settings__actions">
      <HyButton
        class="settings__action-button"
        type="tertiary"
        @click="checkForUpdates"
        :disabled="checkingDisabled"
      >
        {{ $t('settings.check_for_updates') }}
      </HyButton>
      <HyButton
        class="settings__action-button"
        type="tertiary"
        @click="openUninstall"
        :disabled="!canPerformActions"
      >
        {{ $t('settings.uninstall') }}
      </HyButton>
      <HyButton
        class="settings__action-button"
        type="tertiary"
        @click="logout"
        :disabled="!canPerformActions"
      >
        {{ $t('settings.logout') }}
      </HyButton>
    </div>
  </PanelView>
</template>

<style scoped>
.settings__section {
  margin-bottom: 20px;
}

.settings__label {
  font-size: 14px;
  font-weight: 600;
  color: #8b949f;
  margin: 0 0 8px;
  text-transform: uppercase;
}

.settings__version {
  font-weight: 400;
  color: #d2d9e2;
}

.settings__server-status {
  transition: color 0.3s ease;
}

.settings__server-status--active {
  color: #4ade80;
  font-weight: 600;
}

.settings__server-status--starting {
  color: #fbbf24;
  font-weight: 600;
}

.settings__server-controls {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.settings__server-spinner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  background: rgba(251, 191, 36, 0.1);
  border-radius: 6px;
  border: 1px solid rgba(251, 191, 36, 0.3);
}

.spinner-small {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(251, 191, 36, 0.3);
  border-top: 2px solid #fbbf24;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.settings__spinner-text {
  color: #fbbf24;
  font-size: 13px;
  font-weight: 500;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.settings__dropdown {
  width: 100%;
}

.settings__directory-button {
  display: flex;
  align-items: center;
  gap: 8px;
}

.settings__directory-icon {
  height: 14px;
}

.settings__actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 24px;
}

.settings__action-button {
  width: 100%;
}
</style>
