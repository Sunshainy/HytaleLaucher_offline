<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { BrowserOpenURL, ClipboardSetText } from '@wailsjs/runtime/runtime'
import HyButton from '@/components/HyButton.vue'
import PanelView from '@/components/PanelView.vue'

const router = useRouter()
const eulaText = ref('')
const isAccepting = ref(false)

// Placeholder for backend call
async function getEULAText(): Promise<string> {
  return `# End User License Agreement

This is a placeholder EULA text. The actual EULA would be fetched from the backend.

## Terms of Service

By using this software, you agree to the terms and conditions outlined in this agreement.

## Privacy Policy

Your privacy is important to us. Please review our privacy policy for more information.

For questions, contact: support@hytale.com`
}

async function acceptEULA(): Promise<void> {
  // Would call backend to accept EULA
}

async function declineEULA(): Promise<void> {
  // Would call backend to decline and quit
}

function escapeHtml(text: string): string {
  const map: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  }
  return text.replace(/[&<>"']/g, c => map[c])
}

function parseMarkdown(text: string): string {
  if (!text) return ''

  let html = escapeHtml(text)

  // Code blocks
  html = html.replace(/`([^`]+)`/g, '<code>$1</code>')

  // Bold
  html = html.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/__([^_]+)__/g, '<strong>$1</strong>')

  // Italic
  html = html.replace(/\*([^*]+)\*/g, '<em>$1</em>')

  // Headers
  html = html.replace(/^(#{1,6})\s+(.+)$/gm, (_, hashes, content) => {
    const level = hashes.length
    return `<h${level}>${content}</h${level}>`
  })

  // Links
  html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="#" data-href="$2">$1</a>')

  // Email links
  html = html.replace(/([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})/g, '<a href="#" data-href="mailto:$1">$1</a>')

  // Paragraphs
  html = html.split('\n\n').map(p => {
    if (p.trim() && !p.startsWith('<h') && !p.startsWith('<')) {
      return `<p>${p}</p>`
    }
    return p
  }).join('\n')

  return html
}

const parsedEula = computed(() => parseMarkdown(eulaText.value))

function handleClick(event: MouseEvent) {
  const target = event.target as HTMLElement
  const link = target.closest('a[data-href]')
  if (link) {
    event.preventDefault()
    const href = (link as HTMLElement).dataset.href
    if (href) {
      BrowserOpenURL(href)
    }
  }
}

async function accept() {
  isAccepting.value = true
  try {
    await acceptEULA()
    router.push({ name: 'init' })
  } catch (error) {
    router.push({ name: 'error', query: { error: String(error) } })
  } finally {
    isAccepting.value = false
  }
}

async function decline() {
  await declineEULA()
}

function copyText() {
  ClipboardSetText(eulaText.value)
}

onMounted(async () => {
  eulaText.value = await getEULAText()
})
</script>

<template>
  <PanelView
    :title="$t('eula.title')"
    :show-close-button="false"
    :show-report-bug="false"
    panel-width="640px"
    panel-max-height="80vh"
  >
    <div class="eula__section">
      <p class="eula__description">{{ $t('eula.description') }}</p>
    </div>
    <div class="eula__eula-container" @click="handleClick">
      <div v-if="eulaText" class="eula__eula-text" v-html="parsedEula"></div>
      <p v-else class="eula__eula-text">{{ $t('eula.loading') }}</p>
    </div>
    <div class="eula__actions">
      <HyButton
        class="eula__action-button"
        type="primary"
        @click="accept"
        :disabled="isAccepting"
      >
        {{ $t('eula.accept') }}
      </HyButton>
      <HyButton
        class="eula__action-button"
        type="tertiary"
        @click="decline"
      >
        {{ $t('eula.decline') }}
      </HyButton>
      <HyButton
        class="eula__action-button"
        type="tertiary"
        @click="copyText"
        small
      >
        {{ $t('eula.copy_text') }}
      </HyButton>
    </div>
  </PanelView>
</template>

<style scoped>
.eula__section {
  margin-bottom: 16px;
}

.eula__description {
  color: #8b949f;
  font-size: 14px;
  margin: 0;
}

.eula__eula-container {
  max-height: 300px;
  overflow-y: auto;
  background-color: rgba(0, 0, 0, 0.2);
  border: 1px solid #434E65;
  border-radius: 4px;
  padding: 16px;
  margin-bottom: 24px;
}

.eula__eula-text {
  color: #d2d9e2;
  font-size: 14px;
  line-height: 1.6;
}

.eula__eula-text :deep(h1),
.eula__eula-text :deep(h2),
.eula__eula-text :deep(h3) {
  color: #d2d9e2;
  margin: 16px 0 8px;
}

.eula__eula-text :deep(h1) {
  font-size: 18px;
}

.eula__eula-text :deep(h2) {
  font-size: 16px;
}

.eula__eula-text :deep(a) {
  color: #78A1FF;
  text-decoration: none;
}

.eula__eula-text :deep(a:hover) {
  text-decoration: underline;
}

.eula__eula-text :deep(code) {
  background-color: rgba(0, 0, 0, 0.3);
  padding: 2px 6px;
  border-radius: 2px;
}

.eula__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.eula__action-button {
  flex: 1;
  min-width: 120px;
}
</style>
