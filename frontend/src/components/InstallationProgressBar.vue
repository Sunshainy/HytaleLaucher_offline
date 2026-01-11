<script lang="ts" setup>
import { computed } from 'vue'
import { useAppStore } from '@/stores/appStore'
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

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatSpeed(bps: number): string {
  return formatBytes(bps) + '/s'
}
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
}

.installation-progress-bar__top-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.installation-progress-bar__info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.installation-progress-bar__percentage {
  font-size: 18px;
  font-weight: 800;
  color: #d2d9e2;
}

.installation-progress-bar__status {
  font-size: 14px;
  color: #8b949f;
}

.installation-progress-bar__download-progress {
  font-size: 14px;
  color: #8b949f;
}

.installation-progress-bar__bar-container {
  width: 100%;
}

.installation-progress-bar__bar {
  height: 8px;
  background-color: rgba(67, 78, 101, 0.5);
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.installation-progress-bar__bar-fill {
  height: 100%;
  background: linear-gradient(to right, #465DA9, #78A1FF);
  border-radius: 4px;
  transition: width 0.3s ease;
}

.installation-progress-bar__bar-mask {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: url('@/assets/images/progress-bar-masked.png');
  background-size: cover;
}
</style>
