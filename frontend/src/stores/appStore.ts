import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as App from '@wailsjs/go/app/App'

export interface UpdateInfo {
  PrimaryAction: string
  GameVersion: string | null
}

export interface UpdateMessage {
  id: string
  params?: Record<string, any>
}

export interface UpdateStatus {
  message: UpdateMessage
  progress: number
  download_progress?: number
  download_total?: number
  download_bps?: number
}

export interface FeedArticle {
  id: string
  title: string
  description: string
  image_url?: string
  dest_url?: string
}

export const useAppStore = defineStore('app', () => {
  // State
  const currentChannel = ref<string>('')
  const allowedChannels = ref<string[]>([])
  const gameVersion = ref<string | null>(null)
  const lastKnownGoodVersion = ref<string | null>(null)
  const updateInfo = ref<UpdateInfo | null>(null)
  const updateRunning = ref(false)
  const updateStatus = ref<UpdateStatus>({
    message: { id: 'update_status.checking_for_updates' },
    progress: 0
  })
  const isCancellingUpdate = ref(false)
  const canCancel = ref(true)
  const cancellationStatus = ref<UpdateMessage>({ id: '' })
  const feedArticles = ref<FeedArticle[]>([])
  const isValidating = ref(false)

  // Getters
  const hasUpdate = computed(() => updateInfo.value !== null)

  // Actions
  async function fetchChannels() {
    try {
      allowedChannels.value = await App.GetUserChannels()
      const state = await App.GetState()
      if (state?.channel) {
        currentChannel.value = state.channel
      }
    } catch (error) {
      console.error('Failed to fetch channels:', error)
    }
  }

  async function setChannel(channel: string) {
    await App.SetChannel(channel)
    currentChannel.value = channel
  }

  async function fetchGameVersion() {
    try {
      const state = await App.GetState()
      if (state?.dependencies?.game) {
        gameVersion.value = state.dependencies.game.version
      }
    } catch (error) {
      console.error('Failed to fetch game version:', error)
    }
  }

  async function fetchLastKnownGoodVersion() {
    try {
      const state = await App.GetState()
      if (state?.dependencies?.lkg) {
        lastKnownGoodVersion.value = state.dependencies.lkg.version
      }
    } catch (error) {
      console.error('Failed to fetch LKG version:', error)
    }
  }

  async function fetchInstallInfo() {
    await fetchChannels()
    await fetchGameVersion()
    await fetchLastKnownGoodVersion()
  }

  async function checkForUpdates(force = false) {
    try {
      // CheckForUpdates returns a number status code
      const result = await App.CheckForUpdates(force)
      // Non-zero result indicates updates available
      if (result !== 0) {
        updateInfo.value = {
          PrimaryAction: 'Install',
          GameVersion: null
        }
      } else {
        updateInfo.value = null
      }
    } catch (error) {
      console.error('Failed to check for updates:', error)
      throw error
    }
  }

  async function checkForFreestandingLauncherUpdate(): Promise<boolean> {
    // Check if there's a launcher update available
    try {
      const result = await App.CheckForUpdates(true)
      // A specific code might indicate launcher update
      // For now, we'll assume code 2 is launcher update
      return result === 2
    } catch {
      return false
    }
  }

  async function applyUpdates() {
    updateRunning.value = true
    try {
      // Updates are applied via backend events
      // This function initiates the update process
    } catch (error) {
      updateRunning.value = false
      throw error
    }
  }

  async function cancelUpdates() {
    isCancellingUpdate.value = true
    cancellationStatus.value = { id: 'update_status.cancelling_updates' }
    // Cancel logic handled by backend
  }

  async function fetchNewsFeed() {
    try {
      // RefreshNewsFeed returns void - it updates via events
      await App.RefreshNewsFeed()
      // Feed articles will be populated via event handlers
    } catch (error) {
      console.error('Failed to fetch news feed:', error)
    }
  }

  async function uninstall(version: any) {
    // Uninstall logic
  }

  async function validateGameFiles(version: any) {
    // Validation logic
  }

  async function deleteUserData() {
    // Delete user data logic
  }

  return {
    // State
    currentChannel,
    allowedChannels,
    gameVersion,
    lastKnownGoodVersion,
    updateInfo,
    updateRunning,
    updateStatus,
    isCancellingUpdate,
    canCancel,
    cancellationStatus,
    feedArticles,
    isValidating,
    // Getters
    hasUpdate,
    // Actions
    fetchChannels,
    setChannel,
    fetchGameVersion,
    fetchLastKnownGoodVersion,
    fetchInstallInfo,
    checkForUpdates,
    checkForFreestandingLauncherUpdate,
    applyUpdates,
    cancelUpdates,
    fetchNewsFeed,
    uninstall,
    validateGameFiles,
    deleteUserData
  }
})
