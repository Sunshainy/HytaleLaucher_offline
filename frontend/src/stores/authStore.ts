import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as App from '@wailsjs/go/app/App'
import { account as accountNs } from '@wailsjs/go/models'

export interface UserProfile {
  uuid: string
  username: string
}

export const useAuthStore = defineStore('auth', () => {
  // State
  const isOffline = ref(false)
  const hasNetworkBeenChecked = ref(false)
  const account = ref<accountNs.Account | null>(null)
  const userProfiles = ref<UserProfile[]>([])
  const currentProfileUuid = ref<string | null>(null)

  // Getters
  const isLoggedIn = computed(() => account.value !== null)
  const getUserProfile = computed(() => {
    if (!currentProfileUuid.value || userProfiles.value.length === 0) return null
    return userProfiles.value.find(p => p.uuid === currentProfileUuid.value) || null
  })
  const getUserProfiles = computed(() => userProfiles.value)

  // Actions
  async function checkNetworkMode(force = false, reason = '') {
    if (!force && hasNetworkBeenChecked.value) return

    try {
      const result = await App.CheckNetworkMode(force, reason)
      // CheckNetworkMode returns boolean - true means offline
      isOffline.value = result === true
      hasNetworkBeenChecked.value = true
    } catch (error) {
      console.error('Failed to check network mode:', error)
      isOffline.value = true
      hasNetworkBeenChecked.value = true
    }
  }

  async function load() {
    try {
      const accountData = await App.GetAccount()
      if (accountData) {
        account.value = accountData
        // Map Profile[] to UserProfile[] - Profile has 'name' not 'username'
        userProfiles.value = (accountData.profiles || []).map(p => ({
          uuid: p.uuid,
          username: p.name
        }))
        if (accountData.selected_profile) {
          currentProfileUuid.value = accountData.selected_profile
        }
      }
    } catch (error) {
      console.error('Failed to load account:', error)
    }
  }

  async function checkSessionInfo(): Promise<boolean> {
    try {
      const loggedIn = await App.IsLoggedIn()
      if (loggedIn) {
        await load()
        return true
      }
      return false
    } catch (error) {
      console.error('Failed to check session:', error)
      return false
    }
  }

  async function logout() {
    await App.Logout()
    account.value = null
    userProfiles.value = []
    currentProfileUuid.value = null
  }

  async function setUserProfile(uuid: string) {
    await App.SetUserProfile(uuid)
    currentProfileUuid.value = uuid
  }

  return {
    // State
    isOffline,
    hasNetworkBeenChecked,
    account,
    userProfiles,
    currentProfileUuid,
    // Getters
    isLoggedIn,
    getUserProfile,
    getUserProfiles,
    // Actions
    checkNetworkMode,
    load,
    checkSessionInfo,
    logout,
    setUserProfile
  }
})
