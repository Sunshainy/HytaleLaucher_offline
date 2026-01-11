<script lang="ts" setup>
import { useNotificationStore } from '@/stores/notificationStore'

const notificationStore = useNotificationStore()
</script>

<template>
  <div class="notification-container">
    <TransitionGroup name="notification">
      <div
        v-for="notification in notificationStore.notifications"
        :key="notification.id"
        class="notification"
        :class="`notification--${notification.type}`"
      >
        <span class="notification__message">{{ notification.message }}</span>
        <button class="notification__close" @click="notificationStore.remove(notification.id)">
          <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 12 12" fill="none">
            <path d="M0 1.33333L1.33331 0L11.9998 10.6667L10.6665 12L0 1.33333Z" fill="currentColor"/>
            <path d="M10.6667 0L12 1.33333L1.3335 12L0 10.6667L10.6667 0Z" fill="currentColor"/>
          </svg>
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.notification-container {
  position: fixed;
  top: 14px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 3000;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  pointer-events: none;
}

.notification {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  border-radius: 4px;
  min-width: 332px;
  max-width: 400px;
  pointer-events: auto;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  color: #121b25;
  background-size: cover;
  background-position: center;
}

.notification__message {
  font-family: 'Nunito Sans', sans-serif;
  font-size: 14px;
  font-weight: 800;
  flex: 1;
  line-height: 1.4;
}

.notification__close {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 4px;
  margin-left: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.2s ease;
  flex-shrink: 0;
}

.notification__close:hover {
  opacity: 0.8;
}

.notification--success {
  background-color: #3dd390;
}

.notification--error {
  background-color: #f2486a;
}

.notification--info {
  background-color: #598ac3;
  color: #d2d9e2;
}

.notification-enter-active {
  transition: all 0.3s ease-out;
}

.notification-leave-active {
  transition: all 0.3s ease-in;
}

.notification-enter-from {
  opacity: 0;
  transform: translateY(-20px) scale(0.95);
}

.notification-enter-to,
.notification-leave-from {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.notification-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(0.95);
}

.notification-move {
  transition: transform 0.3s ease;
}
</style>
