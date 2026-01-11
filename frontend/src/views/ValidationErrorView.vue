<script lang="ts" setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Logo from '@/components/Logo.vue'
import HyButton from '@/components/HyButton.vue'
import ReportBug from '@/components/ReportBug.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const errorMessage = computed(() => {
  return route.query.error ?? t('error.unknown_error')
})

function goBack() {
  router.back()
}
</script>

<template>
  <div class="error-view">
    <div class="error-view__container">
      <Logo class="error-view__logo" />
      <div class="container text--center">
        <h2 class="error-view__title">{{ $t('validation_error.title') }}</h2>
        <p class="error-view__description" v-html="$t('validation_error.description')"></p>
        <div>
          <HyButton @click="goBack" class="error-view__continue-button" type="secondary">
            {{ $t('validation_error.continue') }}
          </HyButton>
        </div>
        <code class="error-view__message">{{ errorMessage }}</code>
      </div>
    </div>
    <ReportBug class="error-view__report-bug" />
  </div>
</template>

<style scoped>
.error-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
  position: relative;
}

.error-view__container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 32px;
}

.error-view__logo :deep(img) {
  height: 160px;
}

.error-view__title {
  font-size: 24px;
  font-weight: 800;
  color: #d2d9e2;
  margin: 0 0 16px;
  text-transform: uppercase;
}

.error-view__description {
  color: #8b949f;
  font-size: 16px;
  max-width: 400px;
  line-height: 1.6;
}

.error-view__message {
  display: block;
  background-color: rgba(22, 33, 47, 0.8);
  border: 1px solid #434E65;
  border-radius: 4px;
  padding: 16px;
  font-size: 14px;
  color: #f2486a;
  margin: 16px 0;
  max-width: 400px;
  word-break: break-word;
}

.error-view__report-bug {
  position: absolute;
  bottom: 20px;
  right: 20px;
}
</style>
