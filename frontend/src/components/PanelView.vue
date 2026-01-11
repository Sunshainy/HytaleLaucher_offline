<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import ReportBug from './ReportBug.vue'

const props = withDefaults(defineProps<{
  title: string
  panelWidth?: string
  panelMaxHeight?: string
  showCloseButton?: boolean
  showReportBug?: boolean
  escHandler?: (() => void) | null
}>(), {
  panelWidth: '400px',
  panelMaxHeight: '',
  showCloseButton: false,
  showReportBug: false,
  escHandler: null
})

const emit = defineEmits<{
  close: []
}>()

const panelRef = ref<HTMLElement | null>(null)

function close() {
  emit('close')
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    if (props.escHandler) {
      props.escHandler()
    } else if (props.showCloseButton) {
      close()
    }
  }
}

onMounted(() => {
  if (panelRef.value) {
    panelRef.value.focus()
  }
})
</script>

<template>
  <div ref="panelRef" class="panel-view" @keydown="handleKeydown" tabindex="-1">
    <div class="panel-view__container">
      <h1 class="panel-view__title">{{ title }}</h1>
      <div class="panel-view__panel" :style="{ width: panelWidth, maxHeight: panelMaxHeight }">
        <div v-if="showCloseButton" class="panel-view__panel-header">
          <button class="panel-view__close-button" @click="close" :title="$t('common.close')">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 12 12" fill="none">
              <path d="M0 1.33333L1.33331 1.46109e-06L11.9998 10.6667L10.6665 12L0 1.33333Z" fill="#747B84"/>
              <path d="M10.6667 0L12 1.33333L1.3335 12L0.000184319 10.6667L10.6667 0Z" fill="#747B84"/>
            </svg>
          </button>
        </div>
        <slot />
      </div>
    </div>
    <ReportBug v-if="showReportBug" class="panel-view__report-bug" />
  </div>
</template>

<style scoped>
.panel-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
  position: relative;
  outline: none;
}

.panel-view__container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

.panel-view__title {
  font-size: 24px;
  font-weight: 800;
  color: #d2d9e2;
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.panel-view__panel {
  background-color: rgba(22, 33, 47, 0.95);
  border: 2px solid #434E65;
  border-radius: 4px;
  padding: 24px;
  position: relative;
}

.panel-view__panel-header {
  position: absolute;
  top: 8px;
  right: 8px;
}

.panel-view__close-button {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: background-color 0.2s ease;
}

.panel-view__close-button:hover {
  background-color: rgba(116, 123, 132, 0.15);
}

.panel-view__report-bug {
  position: absolute;
  bottom: 20px;
  right: 20px;
}
</style>
