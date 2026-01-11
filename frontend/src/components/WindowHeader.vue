<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { WindowMinimise, Quit, Environment } from '@wailsjs/runtime/runtime'

const router = useRouter()
const isMacOS = ref(false)

onMounted(async () => {
  try {
    const env = await Environment()
    isMacOS.value = env.platform === 'darwin'
  } catch {
    isMacOS.value = false
  }
})

function openSettings() {
  router.push({ name: 'settings' })
}

function minimize() {
  WindowMinimise()
}

function close() {
  Quit()
}
</script>

<template>
  <div class="window-header" :class="{ 'window-header--macos': isMacOS }">
    <div class="window-header__cog-container">
      <button class="window-header__cog" @click="openSettings">
        <svg class="window-header__cog-image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#747B84" stroke-width="2">
          <circle cx="12" cy="12" r="3"/>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
      </button>
    </div>
    <div class="window-header__controls">
      <button class="window-header__button" @click="minimize">
        <svg xmlns="http://www.w3.org/2000/svg" width="12" height="2" viewBox="0 0 12 2" fill="none">
          <rect width="12" height="2" fill="#747B84"/>
        </svg>
      </button>
      <button class="window-header__button" @click="close">
        <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 12 12" fill="none">
          <path d="M0 1.33333L1.33331 1.46109e-06L11.9998 10.6667L10.6665 12L0 1.33333Z" fill="#747B84"/>
          <path d="M10.6667 0L12 1.33333L1.3335 12L0.000184319 10.6667L10.6667 0Z" fill="#747B84"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.window-header {
  position: fixed;
  width: 100%;
  height: 32px;
  display: flex;
  flex-direction: row-reverse;
  justify-content: space-between;
  align-items: center;
  padding: 0 8px;
  --wails-draggable: drag;
  user-select: none;
  z-index: 10;
}

.window-header--macos {
  flex-direction: row;
}

.window-header__cog-container {
  margin-right: 5px;
  margin-top: 12px;
  --wails-draggable: no-drag;
}

.window-header:not(.window-header--macos) .window-header__cog-container {
  margin-right: 0;
  margin-left: 5px;
}

.window-header__cog {
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
  --wails-draggable: no-drag;
  opacity: 0.9;
}

.window-header__cog:hover {
  opacity: 1;
  background-color: rgba(116, 123, 132, 0.15);
}

.window-header__cog:active {
  background-color: rgba(255, 255, 255, 0.15);
}

.window-header__cog-image {
  height: 18px;
  aspect-ratio: 1/1;
  pointer-events: none;
}

.window-header__controls {
  display: flex;
  gap: 4px;
  --wails-draggable: no-drag;
}

.window-header:not(.window-header--macos) .window-header__controls {
  flex-direction: row-reverse;
}

.window-header__button {
  width: 32px;
  height: 24px;
  border: none;
  background: transparent;
  color: #747b84;
  opacity: 0.9;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
  --wails-draggable: no-drag;
}

.window-header__button:hover {
  opacity: 1;
  background-color: rgba(116, 123, 132, 0.15);
}

.window-header__button:active {
  background-color: rgba(255, 255, 255, 0.15);
}

.window-header__button svg {
  pointer-events: none;
}
</style>
