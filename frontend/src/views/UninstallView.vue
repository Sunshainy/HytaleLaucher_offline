<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/appStore'
import { useNotificationStore } from '@/stores/notificationStore'
import { formatBytes } from '@/composables/formatBytes'
import PanelView from '@/components/PanelView.vue'
import HyButton from '@/components/HyButton.vue'

interface InstalledVersion {
  id: string
  channel: string
  version: string
  package: string
  installDir: string
  canValidateFiles: boolean
}

const router = useRouter()
const appStore = useAppStore()
const notificationStore = useNotificationStore()
const { t } = useI18n()

const isLoading = ref(true)
const versions = ref<InstalledVersion[]>([])
const fileSizes = ref<Record<string, number>>({})
const hasUserData = ref(false)
const showConfirmDialog = ref(false)

// Placeholder functions for backend calls
async function getInstalledVersions(): Promise<InstalledVersion[]> {
  return []
}

async function getFileSizes(dirs: string[]): Promise<Record<string, number>> {
  return {}
}

async function checkHasUserData(): Promise<boolean> {
  return true
}

async function resetGameSettings(): Promise<void> {
  // Would reset game settings
}

function close() {
  router.back()
}

function goBack() {
  router.go(-2)
}

function uninstall(version: InstalledVersion) {
  goBack()
  appStore.uninstall(version)
}

async function validate(version: InstalledVersion) {
  goBack()
  appStore.validateGameFiles(version)
}

async function resetSettings() {
  goBack()
  await resetGameSettings()
  notificationStore.showSuccess(t('uninstall.reset_game_settings_success'))
}

function showDeleteConfirm() {
  showConfirmDialog.value = true
}

function confirmDelete() {
  showConfirmDialog.value = false
  goBack()
  appStore.deleteUserData()
}

function cancelDelete() {
  showConfirmDialog.value = false
}

function handleEsc() {
  if (showConfirmDialog.value) {
    cancelDelete()
  } else {
    close()
  }
}

async function delay(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}

async function fetchFileSizes() {
  fileSizes.value = await getFileSizes(versions.value.map(v => v.installDir))
}

onMounted(async () => {
  hasUserData.value = await checkHasUserData()
  isLoading.value = true

  try {
    versions.value = await getInstalledVersions()
    await Promise.race([fetchFileSizes(), delay(1000)])
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <PanelView
    :title="$t('uninstall.title')"
    :show-close-button="true"
    :show-report-bug="true"
    :esc-handler="handleEsc"
    panel-width="480px"
    @close="close"
  >
    <div class="uninstall__section">
      <h2 class="uninstall__label">{{ $t('uninstall.installed_versions') }}</h2>

      <div v-if="isLoading">
        <span class="version-list__status-text">{{ $t('uninstall.scanning_versions') }}</span>
      </div>

      <div v-if="!isLoading && versions.length === 0">
        <span class="version-list__status-text">{{ $t('uninstall.no_versions_installed') }}</span>
      </div>

      <div v-if="!isLoading && versions.length > 0" class="version-list">
        <div v-for="version in versions" :key="version.id" class="version-list__item">
          <div class="version-list__info">
            <span class="version-list__channel">
              {{ version.channel }}
              <span v-if="version.package === 'lkg'" class="version-list__lkg">{{ $t('uninstall.lkg') }}</span>
            </span>
            <span class="version-list__version">{{ version.version }}</span>
            <span v-if="version.installDir in fileSizes" class="version-list__file-size">
              {{ formatBytes(fileSizes[version.installDir]) }}
            </span>
          </div>
          <div class="version-list__actions">
            <HyButton
              v-if="version.canValidateFiles"
              class="version-list__action-button"
              type="tertiary"
              @click="validate(version)"
              small
            >
              {{ $t('uninstall.validate') }}
            </HyButton>
            <HyButton
              class="version-list__action-button"
              type="tertiary"
              @click="uninstall(version)"
              small
            >
              {{ $t('uninstall.uninstall_version') }}
            </HyButton>
          </div>
        </div>
      </div>
    </div>

    <div v-if="!isLoading" class="uninstall__actions">
      <HyButton
        class="uninstall__action-button"
        type="tertiary"
        @click="resetSettings"
        small
      >
        {{ $t('uninstall.reset_game_settings') }}
      </HyButton>
      <HyButton
        class="uninstall__action-button"
        type="destructive"
        @click="showDeleteConfirm"
        :disabled="!hasUserData"
        small
      >
        {{ $t('uninstall.delete_user_data') }}
      </HyButton>
    </div>
  </PanelView>

  <Transition name="fade">
    <div v-if="showConfirmDialog" class="confirm-dialog-overlay" @click.self="cancelDelete">
      <div class="confirm-dialog">
        <p class="confirm-dialog__message">{{ $t('uninstall.confirm_delete_user_data') }}</p>
        <div class="confirm-dialog__actions">
          <HyButton type="tertiary" class="confirm-dialog__action-button" @click="cancelDelete" small>
            {{ $t('common.cancel') }}
          </HyButton>
          <HyButton type="destructive" class="confirm-dialog__action-button" @click="confirmDelete" small>
            {{ $t('common.delete') }}
          </HyButton>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.uninstall__section {
  margin-bottom: 20px;
}

.uninstall__label {
  font-size: 14px;
  font-weight: 600;
  color: #8b949f;
  margin: 0 0 12px;
  text-transform: uppercase;
}

.version-list__status-text {
  color: #8b949f;
  font-size: 14px;
}

.version-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.version-list__item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background-color: rgba(0, 0, 0, 0.2);
  border: 1px solid #434E65;
  border-radius: 4px;
}

.version-list__info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.version-list__channel {
  color: #d2d9e2;
  font-size: 14px;
  font-weight: 600;
}

.version-list__lkg {
  color: #78A1FF;
  font-size: 12px;
  margin-left: 8px;
}

.version-list__version {
  color: #8b949f;
  font-size: 12px;
}

.version-list__file-size {
  color: #8b949f;
  font-size: 12px;
}

.version-list__actions {
  display: flex;
  gap: 8px;
}

.uninstall__actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.confirm-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.confirm-dialog {
  background-color: rgba(22, 33, 47, 0.98);
  border: 2px solid #434E65;
  border-radius: 4px;
  padding: 24px;
  max-width: 360px;
}

.confirm-dialog__message {
  color: #d2d9e2;
  font-size: 14px;
  line-height: 1.6;
  margin: 0 0 20px;
}

.confirm-dialog__actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
