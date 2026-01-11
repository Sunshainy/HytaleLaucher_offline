import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Notification {
  id: number
  message: string
  type: 'success' | 'error' | 'info'
  duration?: number
}

let notificationId = 0

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])

  function show(message: string, type: 'success' | 'error' | 'info' = 'info', duration = 5000) {
    const id = ++notificationId
    const notification: Notification = { id, message, type, duration }
    notifications.value.push(notification)

    if (duration > 0) {
      setTimeout(() => {
        remove(id)
      }, duration)
    }

    return id
  }

  function showSuccess(message: string, duration = 5000) {
    return show(message, 'success', duration)
  }

  function showError(message: string, duration = 5000) {
    return show(message, 'error', duration)
  }

  function showInfo(message: string, duration = 5000) {
    return show(message, 'info', duration)
  }

  function remove(id: number) {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index !== -1) {
      notifications.value.splice(index, 1)
    }
  }

  function clear() {
    notifications.value = []
  }

  return {
    notifications,
    show,
    showSuccess,
    showError,
    showInfo,
    remove,
    clear
  }
})
