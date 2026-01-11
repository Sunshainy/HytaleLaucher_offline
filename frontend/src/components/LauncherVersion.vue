<script lang="ts" setup>
import { ref, onMounted } from 'vue'

const version = ref('unknown')

onMounted(async () => {
  // Version would be fetched from Go backend
  // For now, use placeholder
  version.value = '1.0.0'
})
</script>

<template>
  <span class="launcher-version">{{ version }}</span>
</template>

<style scoped>
.launcher-version {
  color: #8b949f;
  font-size: 12px;
  font-family: 'Nunito Sans', sans-serif;
  font-weight: 500;
}
</style>
