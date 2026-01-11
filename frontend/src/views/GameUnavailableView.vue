<script lang="ts" setup>
import Logo from '@/components/Logo.vue'
</script>

<template>
  <div class="game-unavailable">
    <div class="game-unavailable__container">
      <Logo class="game-unavailable__logo" />
      <div class="game-unavailable__content">
        <h1 class="game-unavailable__title">{{ $t('game_unavailable.title') }}</h1>
        <p class="game-unavailable__description">{{ $t('game_unavailable.description') }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.game-unavailable {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
  position: relative;
  overflow: hidden;
}

.game-unavailable__container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
  position: relative;
  z-index: 1;
}

.game-unavailable__logo {
  margin-bottom: 32px;
}

.game-unavailable__logo :deep(img) {
  height: 160px;
}

.game-unavailable__content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  max-width: 420px;
  padding: 0 24px;
}

.game-unavailable__title {
  font-size: 24px;
  font-weight: 800;
  color: #d2d9e2;
  margin: 0 0 14px;
}

.game-unavailable__description {
  font-size: 16px;
  font-weight: 500;
  color: #8b949f;
  line-height: 1.6;
  margin: 0;
}
</style>
