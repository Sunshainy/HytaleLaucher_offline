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

const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const { t } = useI18n()

const isCheckingUpdates = ref(false)

const openInIcon = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABEAAAARCAYAAAA7bUf6AAAACXBIWXMAAAsTAAALEwEAmpwYAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAC8SURBVHgBrZK9EQIhEIX5Cwgp4SJmCCnBCmzFDjxLsAM7sQTMSC3hKgB3Ax08gT3v7iXsDPDxeLs8xjiybz2dczcscE8IcWaERG8TYGNK6cIIqfJCCwSOWM9R10kJyjlfN0HAycA5P66GIAC+codyWAWpATDoedjqX8C7AWXYTSdSSgOLqQFQZfubEGvtA8I8QDnNAT8gnMrK1H4UQjCMkKIOeO8n6gES0pPW+oTromGjtAuE90Jdql2cvAClzFdGDZMFsAAAAABJRU5ErkJggg=='

const allowedChannels = computed(() => appStore.allowedChannels)
const currentChannel = computed(() => appStore.currentChannel)

const checkingDisabled = computed(() => {
  return isCheckingUpdates.value || authStore.isOffline || appStore.updateRunning
})

const lkgDisabled = computed(() => {
  return !appStore.lastKnownGoodVersion || appStore.updateRunning
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
  // Would call backend to open directory
  console.log('Opening directory...')
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

function launchLKG() {
  // Would call backend to launch LKG version
  console.log('Launching LKG...')
}

function close() {
  router.back()
}

function openUninstall() {
  router.push({ name: 'uninstall' })
}

onMounted(async () => {
  await appStore.fetchLastKnownGoodVersion()
  await appStore.fetchChannels()
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
        {{ $t('settings.last_known_good_version') }}
        <span class="settings__version">
          ({{ appStore.lastKnownGoodVersion ? appStore.lastKnownGoodVersion : $t('settings.not_available') }})
        </span>
      </h2>
      <HyButton
        class="settings__action-button"
        type="tertiary"
        @click="launchLKG"
        :disabled="lkgDisabled"
      >
        {{ $t('settings.launch_lkg') }}
      </HyButton>
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
