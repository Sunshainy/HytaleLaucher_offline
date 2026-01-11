<script lang="ts" setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/authStore'

const authStore = useAuthStore()
const isChecking = ref(false)

async function goOnline() {
  isChecking.value = true
  try {
    await authStore.checkNetworkMode(true, 'user_request')
  } finally {
    isChecking.value = false
  }
}
</script>

<template>
  <div class="offline-footer">
    <div class="offline-footer__content">
      <div class="offline-footer__indicator"></div>
      <span class="offline-footer__text">{{ $t('offline.status') }}</span>
      <button
        class="offline-footer__button"
        @click="goOnline"
        :disabled="isChecking"
      >
        {{ isChecking ? $t('offline.checking') : $t('offline.go_online') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.offline-footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  width: 100%;
  height: 30px;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 0 8px;
  z-index: 10;
  background-color: rgba(18, 27, 37, 0.5);
}

.offline-footer__content {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: center;
}

.offline-footer__indicator {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: #f2486a;
  flex-shrink: 0;
}

.offline-footer__text {
  color: #d2d9e2;
  font-family: 'Nunito Sans', sans-serif;
  font-size: 13px;
  font-weight: 500;
}

.offline-footer__button {
  background-color: transparent;
  border: 1px solid #D2D9E2;
  border-radius: 2px;
  color: #d2d9e2;
  font-family: 'Nunito Sans', sans-serif;
  font-size: 13px;
  font-weight: 500;
  padding: 2px 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.offline-footer__button:hover:not(:disabled) {
  background-color: rgba(210, 217, 226, 0.1);
}

.offline-footer__button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
