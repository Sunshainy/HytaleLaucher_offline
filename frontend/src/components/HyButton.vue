<script lang="ts" setup>
defineProps<{
  type?: 'primary' | 'secondary' | 'tertiary' | 'destructive'
  small?: boolean
  disabled?: boolean
}>()
</script>

<template>
  <button
    class="hybutton"
    :class="[
      `hybutton--${type || 'primary'}`,
      { 'hybutton--small': small }
    ]"
    :disabled="disabled"
  >
    <slot />
  </button>
</template>

<style scoped>
.hybutton {
  height: 44px;
  text-transform: uppercase;
  font-weight: 800;
  font-size: 18px;
  letter-spacing: 0.4px;
  leading-trim: both;
  text-edge: cap;
  font-style: normal;
  line-height: normal;
  border-radius: 4px;
  border: 2px solid;
  cursor: pointer;
  background-size: cover;
  background-position: center;
  font-family: 'Nunito Sans', sans-serif;
  color: #d2d9e2;
  user-select: none;
  -webkit-user-select: none;
  -webkit-user-drag: none;
  box-sizing: border-box;
}

.hybutton--primary {
  background-image: url('@/assets/images/runes.png'), linear-gradient(to right, #465DA9, #4578AA);
  box-shadow: 0 0 rgba(0, 0, 0, 0.25) inset;
  background-size: 150%;
  border: 2px solid #598AC3;
}

.hybutton--primary:hover:not(:disabled) {
  background-image: url('@/assets/images/runes.png');
  background-color: #20578e;
  border: 2px solid #78A1FF;
  box-shadow: 0 0 0 2px #78a1ff;
}

.hybutton--primary:active:not(:disabled) {
  background-image: url('@/assets/images/runes.png');
  background-color: #2a74bd;
  border: 2px solid #78A1FF;
  box-shadow: 0 0 0 2px #78a1ff;
}

.hybutton--secondary {
  border-color: #47516b;
  background: #273349;
}

.hybutton--secondary:hover:not(:disabled) {
  border: 2px solid #576C95;
  box-shadow: 0 0 0 2px #576c95;
  background: #18283C;
}

.hybutton--secondary:active:not(:disabled) {
  border: 2px solid #576C95;
  box-shadow: 0 0 0 2px #576c95;
  background: #253F60;
}

.hybutton--tertiary {
  background: transparent;
  border-color: #47516b;
}

.hybutton--tertiary:hover:not(:disabled) {
  border: 2px solid #78A1FF;
  box-shadow: 0 0 0 2px #78a1ff;
}

.hybutton--tertiary:active:not(:disabled) {
  border: 2px solid #78A1FF;
  box-shadow: 0 0 0 2px #78a1ff;
  background: rgba(78, 89, 119, 0.3);
}

.hybutton--destructive {
  background: linear-gradient(to right, #8B2A2A, #A33030);
  border-color: #c94545;
}

.hybutton--destructive:hover:not(:disabled) {
  background: #6B1E1E;
  border: 2px solid #E74C3C;
  box-shadow: 0 0 0 2px #e74c3c;
}

.hybutton--destructive:active:not(:disabled) {
  background: #8B2A2A;
  border: 2px solid #E74C3C;
  box-shadow: 0 0 0 2px #e74c3c;
}

.hybutton:disabled {
  background-image: url('@/assets/images/runes.png');
  border: 2px solid #464E58;
  background: #39424D;
  cursor: not-allowed;
}

.hybutton--small {
  height: 32px;
  font-size: 16px;
  letter-spacing: 0.4px;
}
</style>
