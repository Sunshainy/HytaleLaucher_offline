<script lang="ts" setup>
import { computed } from 'vue'
import { useAppStore } from '@/stores/appStore'
import { formatBytes, formatSpeed } from '@/composables/formatBytes'
import HyButton from './HyButton.vue'

const props = defineProps<{
  onCancel?: (() => void) | null
}>()

const appStore = useAppStore()

const percentage = computed(() => Math.round(appStore.updateStatus.progress * 100))

const downloadedDisplay = computed(() => {
  const bytes = appStore.updateStatus.download_progress
  return bytes ? formatBytes(bytes) : '0 B'
})

const totalDisplay = computed(() => {
  const bytes = appStore.updateStatus.download_total
  return bytes ? formatBytes(bytes) : null
})

const speedDisplay = computed(() => {
  const bps = appStore.updateStatus.download_bps
  return bps ? formatSpeed(bps) : '0 B/s'
})
</script>

<template>
  <div class="installation-progress-bar">
    <div class="installation-progress-bar__top-container">
      <div class="installation-progress-bar__info">
        <span class="installation-progress-bar__percentage">{{ percentage }}%</span>
        <span class="installation-progress-bar__status">
          {{ $t('update_status.' + appStore.updateStatus.message.id, appStore.updateStatus.message.params || {}) }}
        </span>
        <span v-if="totalDisplay" class="installation-progress-bar__status"> - </span>
        <span v-if="totalDisplay" class="installation-progress-bar__download-progress">
          {{ downloadedDisplay }}/{{ totalDisplay }} - {{ speedDisplay }}
        </span>
      </div>
      <div class="installation-progress-bar__actions">
        <HyButton
          v-if="onCancel && !appStore.isCancellingUpdate && appStore.canCancel"
          type="secondary"
          small
          class="installation-progress-bar__button--cancel"
          @click="onCancel"
        >
          {{ $t('common.cancel') }}
        </HyButton>
      </div>
    </div>
    <div class="installation-progress-bar__bar-container">
      <div class="installation-progress-bar__bar">
        <div class="installation-progress-bar__bar-fill" :style="{ width: `${percentage}%` }"></div>
        <div class="installation-progress-bar__bar-mask" :style="{ clipPath: `inset(0 ${100 - percentage}% 0 0)` }"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.installation-progress-bar {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.installation-progress-bar__top-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.installation-progress-bar__info {
  display: flex;
  align-items: baseline;
}

.installation-progress-bar__percentage {
  font-size: 22px;
  font-weight: 800;
  color: #d2d9e2;
  font-family: 'Nunito Sans', sans-serif;
  line-height: 1;
  margin-right: 12px;
}

.installation-progress-bar__status {
  font-size: 14px;
  color: #8b949f;
  font-family: 'Nunito Sans', sans-serif;
  font-weight: 500;
  margin-right: 4px;
}

.installation-progress-bar__download-progress {
  font-size: 14px;
  color: #d1d1d1;
  font-family: 'Nunito Sans', sans-serif;
  font-weight: 500;
}

.installation-progress-bar__bar-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.installation-progress-bar__bar {
  flex: 1;
  height: 16px;
  background-color: #2a2f38;
  border-radius: 2px;
  overflow: hidden;
  position: relative;
}

.installation-progress-bar__bar-fill {
  height: 100%;
  background: linear-gradient(to right, #4B41B0, #6395CD, #9AADEA);
  background-size: 200% 100%;
  border-radius: 2px;
  transition: width 0.2s ease;
  position: absolute;
  top: 0;
  left: 0;
  z-index: 1;
  animation: gradient-shift 3s ease infinite;
}

.installation-progress-bar__bar-mask {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: url('/assets/progress-bar-masked.2231806c.png');
  background-size: 200% 100%;
  background-position: 0% center;
  background-repeat: repeat-x;
  mix-blend-mode: overlay;
  pointer-events: none;
  z-index: 2;
  transition: clip-path 0.2s ease;
  animation: mask-scroll 10s linear infinite, pulse 3s ease-in-out infinite;
}

.installation-progress-bar__actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.installation-progress-bar__button--pause {
  width: 32px;
  padding: 0;
  position: relative;
}

.installation-progress-bar__button--pause svg {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.installation-progress-bar__button--cancel {
  width: 150px;
}

@keyframes mask-scroll {
  0% {
    background-position: -100% center;
  }
  100% {
    background-position: 100% center;
  }
}

@keyframes pulse {
  0%, 100% {
    transform: scaleY(1);
  }
  50% {
    transform: scaleY(1.1);
  }
}

@keyframes gradient-shift {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}
</style>
