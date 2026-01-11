<script lang="ts" setup>
import hytaleLogo from '@/assets/images/hytale.png'
</script>

<template>
  <div class="logo">
    <img class="logo__image" :src="hytaleLogo" alt="Hytale logo" draggable="false" />
  </div>
</template>

<style scoped>
.logo {
  display: flex;
  justify-content: center;
  align-items: center;
}

.logo__image {
  height: auto;
  max-height: 160px;
}
</style>
