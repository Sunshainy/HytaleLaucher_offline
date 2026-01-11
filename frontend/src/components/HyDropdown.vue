<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface Option {
  label: string
  value: string | number
}

const props = withDefaults(defineProps<{
  modelValue?: string | number | null
  options: (Option | string | number)[]
  placeholder?: string
  disabled?: boolean
}>(), {
  modelValue: null,
  placeholder: 'Select an option',
  disabled: false
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
  'change': [value: string | number]
}>()

const isOpen = ref(false)
const isFocused = ref(false)

const selectedOption = computed(() => {
  return props.options.find(o => {
    const value = typeof o === 'object' ? o.value : o
    return value === props.modelValue
  }) || null
})

const displayText = computed(() => {
  if (selectedOption.value) {
    return typeof selectedOption.value === 'object'
      ? selectedOption.value.label
      : selectedOption.value
  }
  return props.placeholder
})

const chevronSrc = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAwAAAAMCAYAAABWdVznAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAEDSURBVHgBtY69isJAFIVHZWEhxZKwlWyRYotlZdln2HJfZh9hMxb6ANoIVtproxJFRMRCjBF/JlGiaCH4DzGlTeY6BgMRFBs9cGA457tzL0KPVoDn+Zd9KETfOe7JNE2bZX5mOPXOWxTFZ8t6DQiCwKFEuhDVJitoko2dKbQlD3g0yhbVf9ZRfbKjyVQ+4rS5Uge39C0o2hrkKnGHkFwhkkI2oGpbyJU7+Jj5Tj9RudrDweCbRNk1xqD3ZyM/fH58xX2sXS7m4d+fb+xudU9ApbqO9akJKlmAQuYwmO6gwjIv45UT1BpDPJpZMGauNw3pGnw2pHaNWKtvxG/B6AJwE76/DicmdRRP1qsgAAAAAElFTkSuQmCC'

function toggleDropdown() {
  if (!props.disabled) {
    isOpen.value = !isOpen.value
  }
}

function selectOption(option: Option | string | number) {
  const value = typeof option === 'object' ? option.value : option
  emit('update:modelValue', value)
  emit('change', value)
  isOpen.value = false
}

function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('.hy-dropdown')) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div
    class="hy-dropdown"
    :class="{
      'hy-dropdown--open': isOpen,
      'hy-dropdown--focused': isFocused,
      'hy-dropdown--disabled': disabled
    }"
  >
    <div
      class="hy-dropdown__header"
      @click="toggleDropdown"
      @mouseenter="isFocused = true"
      @mouseleave="isFocused = false"
      :tabindex="disabled ? -1 : 0"
      @focus="isFocused = true"
      @blur="isFocused = false"
    >
      <span class="hy-dropdown__text">{{ displayText }}</span>
      <img
        :src="chevronSrc"
        class="hy-dropdown__chevron"
        :class="{ 'hy-dropdown__chevron--open': isOpen }"
        draggable="false"
      />
    </div>
    <div v-if="isOpen" class="hy-dropdown__list">
      <div
        v-for="(option, index) in options"
        :key="index"
        class="hy-dropdown__item"
        :class="{ 'hy-dropdown__item--selected': (typeof option === 'object' ? option.value : option) === modelValue }"
        @click="selectOption(option)"
      >
        {{ typeof option === 'object' ? option.label : option }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.hy-dropdown {
  position: relative;
  user-select: none;
  -webkit-user-select: none;
}

.hy-dropdown__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background-color: #16212f;
  border: 2px solid #434E65;
  border-radius: 2px;
  cursor: pointer;
  height: 32px;
  user-select: none;
  -webkit-user-select: none;
}

.hy-dropdown__header:hover,
.hy-dropdown--focused .hy-dropdown__header {
  border-color: #6e82ac;
}

.hy-dropdown--open .hy-dropdown__header {
  border-color: #6e82ac;
  border-bottom: 2px solid #16212F;
  border-bottom-left-radius: 0;
  border-bottom-right-radius: 0;
}

.hy-dropdown--disabled .hy-dropdown__header {
  opacity: 0.5;
  cursor: not-allowed;
}

.hy-dropdown__text {
  color: #8b949f;
  font-size: 14px;
  font-family: 'Nunito Sans', sans-serif;
  font-weight: 500;
  flex: 1;
  text-align: left;
}

.hy-dropdown__chevron {
  transition: transform 0.2s ease;
  flex-shrink: 0;
  margin-left: 8px;
}

.hy-dropdown__chevron--open {
  transform: rotate(180deg);
}

.hy-dropdown__list {
  position: absolute;
  top: calc(100% - 2px);
  left: 0;
  right: 0;
  background-color: #16212f;
  border: 2px solid #6E82AC;
  border-top: none;
  border-radius: 0 0 2px 2px;
  max-height: 200px;
  overflow-y: auto;
  z-index: 1000;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
}

.hy-dropdown__item {
  padding: 6px 12px;
  color: #d2d9e2;
  font-size: 14px;
  font-family: 'Nunito Sans', sans-serif;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
  user-select: none;
  -webkit-user-select: none;
}

.hy-dropdown__item:hover {
  background-color: rgba(120, 161, 255, 0.1);
}

.hy-dropdown__item--selected {
  background-color: rgba(120, 161, 255, 0.15);
}
</style>
