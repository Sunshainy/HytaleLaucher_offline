<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import HyButton from './HyButton.vue'

interface NewsArticle {
  id: string
  title: string
  description: string
  imageUrl?: string
  link?: string
}

const props = defineProps<{
  articles: NewsArticle[]
}>()

const emit = defineEmits<{
  (e: 'details', article: NewsArticle): void
}>()

const currentIndex = ref(0)
let autoRotateInterval: ReturnType<typeof setInterval> | null = null

const currentArticle = computed(() => {
  if (props.articles.length === 0) return null
  return props.articles[currentIndex.value]
})

const hasArticles = computed(() => props.articles.length > 0)

function goToSlide(index: number) {
  currentIndex.value = index
  resetAutoRotate()
}

function nextSlide() {
  if (props.articles.length > 0) {
    currentIndex.value = (currentIndex.value + 1) % props.articles.length
  }
}

function resetAutoRotate() {
  if (autoRotateInterval) {
    clearInterval(autoRotateInterval)
  }
  autoRotateInterval = setInterval(nextSlide, 8000)
}

function handleDetails() {
  if (currentArticle.value) {
    emit('details', currentArticle.value)
  }
}

function openLink() {
  if (currentArticle.value?.link) {
    // Use Wails runtime to open URL in browser
    if (window.runtime?.BrowserOpenURL) {
      window.runtime.BrowserOpenURL(currentArticle.value.link)
    }
  }
}

onMounted(() => {
  if (props.articles.length > 1) {
    autoRotateInterval = setInterval(nextSlide, 8000)
  }
})

onUnmounted(() => {
  if (autoRotateInterval) {
    clearInterval(autoRotateInterval)
  }
})
</script>

<template>
  <div class="news-carousel">
    <template v-if="hasArticles && currentArticle">
      <div class="news-carousel__item">
        <div
          class="news-carousel__image"
          :style="currentArticle.imageUrl ? { backgroundImage: `url(${currentArticle.imageUrl})` } : {}"
        />
        <div class="news-carousel__content">
          <h3 class="news-carousel__title">{{ currentArticle.title }}</h3>
          <p class="news-carousel__description">{{ currentArticle.description }}</p>
          <div class="news-carousel__details-button-container">
            <HyButton
              type="secondary"
              size="small"
              class="news-carousel__details-button"
              @click="openLink"
            >
              {{ $t('common.details') }}
            </HyButton>
          </div>
        </div>
      </div>
      <div v-if="articles.length > 1" class="news-carousel__controls">
        <div class="news-carousel__indicators">
          <button
            v-for="(article, index) in articles"
            :key="article.id"
            class="news-carousel__indicator"
            :class="{ 'news-carousel__indicator--active': index === currentIndex }"
            @click="goToSlide(index)"
          >
            <svg
              v-if="index === currentIndex"
              class="news-carousel__indicator-svg"
              viewBox="0 0 28 8"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <rect width="28" height="8" rx="2" fill="#78A1FF"/>
            </svg>
          </button>
        </div>
      </div>
    </template>
    <div v-else class="news-carousel__empty" />
  </div>
</template>

<style scoped>
.news-carousel {
  width: 512px;
  position: relative;
  display: flex;
  flex-direction: column;
  border-radius: 4px;
}

.news-carousel__item {
  display: flex;
  height: 240px;
  width: 100%;
  background-color: #121b25;
}

.news-carousel__image {
  width: 50%;
  height: 100%;
  background-size: cover;
  background-position: center;
}

.news-carousel__content {
  width: 50%;
  height: 100%;
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: flex-start;
  overflow: hidden;
}

.news-carousel__title {
  color: #d2d9e2;
  font-size: 16px;
  font-weight: 800;
  margin: 0;
  padding: 14px 14px 12px;
  font-family: 'Nunito Sans', sans-serif;
  line-height: 1.2;
}

.news-carousel__description {
  color: #8b949f;
  font-size: 14px;
  line-height: 1.5;
  margin: 0;
  padding: 0 14px;
  font-weight: 500;
  font-family: 'Nunito Sans', sans-serif;
  flex: 1;
  min-height: 84px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 4;
  text-overflow: ellipsis;
  word-break: break-word;
}

.news-carousel__details-button-container {
  display: flex;
  justify-content: flex-end;
  width: 100%;
  flex-grow: 1;
  align-self: flex-end;
}

.news-carousel__details-button {
  width: 186px;
  align-self: flex-end;
}

.news-carousel__controls {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 12px;
}

.news-carousel__indicators {
  display: flex;
  gap: 4px;
  align-items: center;
  height: 8px;
}

.news-carousel__indicator {
  width: 16px;
  height: 4px;
  border: none;
  background-color: #777f8d;
  cursor: pointer;
  transition: all 0.1s ease;
  transform-origin: center;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.news-carousel__indicator--active {
  background-color: transparent;
  width: 28px;
  height: 8px;
}

.news-carousel__indicator-svg {
  display: block;
  width: 100%;
  height: 100%;
}

.news-carousel__empty {
  height: 260px;
}
</style>
