<script lang="ts" setup>
import curseBg from '@/assets/images/curse-breaker-bg-2.jpg'
import curseGlow from '@/assets/images/curse-glow.png'
import cursePortal from '@/assets/images/curse-portal-front.png'
import curseTentacles from '@/assets/images/curse-tentacles.png'
import hypixelLogo from '@/assets/images/hypixel-studios-logo.png'
import runes from '@/assets/images/runes.png'

defineProps<{
  showGradient?: boolean
}>()

const bgStyle = `url(${curseBg})`
</script>

<template>
  <div class="curse-background" :class="{ 'curse-background--gradient': showGradient }">
    <img :src="curseGlow" class="curse-background__curse-glow" alt="" />
    <img :src="cursePortal" class="curse-background__curse-portal-front" alt="" />
    <img :src="hypixelLogo" class="curse-background__hypixel-logo" alt="Hypixel Studios" />
    <img :src="curseTentacles" class="curse-background__curse-tentacles" alt="" />
    <img :src="curseTentacles" class="curse-background__curse-tentacles-2" alt="" />
    <img :src="runes" class="curse-background__runes" alt="" />
  </div>
</template>

<style scoped>
.curse-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: v-bind(bgStyle);
  background-position: center;
  background-size: cover;
  z-index: 0;
}

.curse-background > img {
  position: absolute;
}

.curse-background__curse-glow {
  display: none;
  object-fit: cover;
  width: 50%;
  height: 50%;
  right: -150px;
  top: 100px;
  opacity: 0;
  transform-origin: center center;
}

.curse-background__curse-portal-front {
  height: 100%;
  width: auto;
  left: 0;
  top: 0;
  bottom: 0;
  opacity: 1;
  object-fit: contain;
}

.curse-background__hypixel-logo {
  position: absolute;
  left: 90px;
  top: 90px;
  max-width: 215px;
  width: 215px;
  height: auto;
  z-index: 2;
}

.curse-background__curse-tentacles,
.curse-background__curse-tentacles-2 {
  width: 100px;
  right: 53px;
  top: 138px;
  opacity: 1;
  animation: curse-tentacles-2-animation 4s ease-out infinite;
  transform-origin: right center;
}

.curse-background__curse-tentacles-2 {
  animation-delay: 2s;
}

.curse-background__runes {
  width: 240px;
  right: 30px;
  top: 75px;
  opacity: 1;
  animation: curse-tentacles-2-animation 4s ease-out infinite;
  position: absolute;
  z-index: 3;
}

.curse-background::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(48deg, #121B25 18.87%, rgba(18, 27, 37, 0.81) 50.18%, rgba(18, 27, 37, 0.47) 67.87%, rgba(18, 27, 37, 0) 84.38%);
  pointer-events: none;
  z-index: 0;
  opacity: 0;
  transition: opacity 0.7s ease;
}

.curse-background--gradient::after {
  opacity: 1;
}

@keyframes curse-glow-animation {
  0% { transform: translate(0); opacity: 0; }
  50% { transform: translate(-1px); opacity: 0.3; }
  100% { transform: translate(0); opacity: 0; }
}

@keyframes curse-portal-front-animation {
  0% { transform: translate(0) scale(1); opacity: 0; }
  50% { opacity: 0.5; }
  100% { transform: translate(-5px) scale(1.1); opacity: 0; }
}

@keyframes curse-tentacles-2-animation {
  0% { transform: scaleX(1); opacity: 0; }
  50% { opacity: 1; }
  100% { transform: scaleX(1.07); opacity: 0; }
}
</style>
