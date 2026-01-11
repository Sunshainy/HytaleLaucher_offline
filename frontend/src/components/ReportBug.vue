<script lang="ts" setup>
import { BrowserOpenURL } from '@wailsjs/runtime/runtime'
import bugReportImage from '@/assets/images/bug-report-image.png'

const props = withDefaults(defineProps<{
  url?: string
}>(), {
  url: 'https://accounts.hytale.com/feedback'
})

function openFeedback() {
  BrowserOpenURL(props.url)
}
</script>

<template>
  <button class="report-bug" @click="openFeedback" :title="$t('common.report_bug')">
    <img :src="bugReportImage" :alt="$t('common.report_bug')" class="report-bug__image" draggable="false" />
  </button>
</template>

<style scoped>
.report-bug {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.2s ease;
}

.report-bug:hover {
  opacity: 0.8;
}

.report-bug:active {
  opacity: 1;
}

.report-bug__image {
  width: 130px;
  height: 130px;
  display: block;
  object-fit: contain;
}
</style>
