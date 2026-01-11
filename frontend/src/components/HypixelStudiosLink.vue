<script lang="ts" setup>
import { BrowserOpenURL } from '@wailsjs/runtime/runtime'
import hypixelLogo from '@/assets/images/hypixel-studios-logo.png'

function openHypixel() {
  BrowserOpenURL('https://hypixelstudios.com')
}
</script>

<template>
  <a class="hypixel-studios-link" href="#" @click.prevent="openHypixel">
    <img :src="hypixelLogo" alt="Hypixel Studios" class="hypixel-studios-link__logo" />
  </a>
</template>

<style scoped>
.hypixel-studios-link {
  display: flex;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  transition: opacity 0.2s ease;
}

.hypixel-studios-link:hover {
  opacity: 1;
}

.hypixel-studios-link__logo {
  height: auto;
  max-width: 90px;
  aspect-ratio: 45/23;
  opacity: 0.8;
  display: block;
}
</style>
