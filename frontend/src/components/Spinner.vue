<script lang="ts" setup>
defineProps<{
  message?: string
}>()
</script>

<template>
  <div class="spinner">
    <div class="spinner__loader"></div>
    <p class="spinner__message">{{ message || 'Loading...' }}</p>
  </div>
</template>

<style scoped>
.spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 200px;
}

.spinner__loader {
  width: 24px;
  height: 24px;
  border: 3px solid rgba(210, 217, 226, 0.2);
  border-top: 3px solid #78A1FF;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 14px;
}

.spinner__message {
  color: #8b949f;
  font-family: 'Nunito Sans', sans-serif;
  font-size: 14px;
  font-weight: 500;
  margin: 0;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
