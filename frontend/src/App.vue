<template>
  <!-- Bootstrap Screen -->
  <BootstrapScreen v-if="!appStore.bootstrapReady" />

  <!-- Main App -->
  <div v-else :class="['h-screen w-screen flex flex-col overflow-hidden rounded-2xl relative', { 'macos-glass': isMac }]" :style="mainBgStyle">
    <!-- Subtle gradient orbs (ambient background) -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-20 -right-20 w-60 h-60 rounded-full" :class="isMac ? 'opacity-[0.06]' : 'opacity-[0.03]'" style="background: radial-gradient(circle, #FF3B3B, transparent 70%);"></div>
      <div class="absolute -bottom-32 -left-20 w-80 h-80 rounded-full" :class="isMac ? 'opacity-[0.04]' : 'opacity-[0.02]'" style="background: radial-gradient(circle, #D946EF, transparent 70%);"></div>
    </div>

    <!-- Glass outer border -->
    <div :class="['absolute inset-0 rounded-2xl pointer-events-none z-10', isMac ? 'border border-white/[0.12] shadow-[inset_0_1px_0_rgba(255,255,255,0.08)]' : 'border border-white/[0.06]']"></div>

    <!-- Drag bar (invisible, top area) -->
    <div class="wails-drag h-3 flex-shrink-0 cursor-grab relative z-20" />

    <!-- Navigation tabs -->
    <div class="flex-shrink-0 flex items-center gap-1 px-3 pb-2 relative z-20">
      <button
        v-for="tab in tabs"
        :key="tab.route"
        @click="$router.push(tab.route)"
        :class="[
          'wails-nodrag px-3 py-1.5 rounded-lg text-[11px] font-medium transition-all duration-200 ease-spring flex items-center gap-1.5',
          isActiveTab(tab.route)
            ? 'glass-sm text-white shadow-glass-sm'
            : 'text-white/40 hover:text-white/70 hover:bg-white/[0.04]'
        ]"
      >
        <component :is="tab.icon" :class="['w-3.5 h-3.5 transition-colors duration-200', isActiveTab(tab.route) ? 'text-accent-light' : '']" />
        {{ tab.label }}
      </button>

      <!-- Spacer -->
      <div class="flex-1" />

      <!-- Active downloads badge -->
      <div v-if="downloadStore.activeCount > 0" class="flex items-center gap-1.5 px-2.5 py-1 rounded-lg glass-accent animate-scale-in">
        <div class="w-1.5 h-1.5 rounded-full bg-accent-light animate-pulse-dot" />
        <span class="text-[10px] text-accent-light font-mono font-medium">{{ downloadStore.activeCount }}</span>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-hidden border-t border-white/[0.04] relative z-20">
      <router-view v-slot="{ Component }">
        <transition name="v" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </div>

    <!-- Status bar -->
    <div class="flex-shrink-0 h-7 flex items-center justify-between px-3 border-t border-white/[0.04] relative z-20">
      <div class="flex items-center gap-2">
        <div :class="['w-1.5 h-1.5 rounded-full transition-colors', appStore.ytdlpReady ? 'bg-emerald-400 shadow-[0_0_6px_rgba(52,211,153,0.4)]' : 'bg-red-400 shadow-[0_0_6px_rgba(248,113,113,0.4)]']" />
        <span class="text-[9px] text-white/20 font-mono">yt-dlp</span>
      </div>
      <div class="flex items-center gap-3">
        <span v-if="downloadStore.completedCount > 0" class="text-[9px] text-white/25 font-mono">
          {{ downloadStore.completedCount }} done
        </span>
        <span class="text-[9px] text-white/15 font-mono">ESC hide · ⌃⇧Y toggle</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from './stores/app'
import { useDownloadStore } from './stores/downloads'
import { useSettingsStore } from './stores/settings'
import BootstrapScreen from './components/BootstrapScreen.vue'
import { HideWindow } from '../wailsjs/go/main/App'
import { EventsOn } from '../wailsjs/runtime/runtime'
import { Download, Package, Settings } from 'lucide-vue-next'

const appStore = useAppStore()
const downloadStore = useDownloadStore()
const settingsStore = useSettingsStore()
const route = useRoute()

const isMac = ref(navigator.platform?.toUpperCase().includes('MAC') || navigator.userAgent?.includes('Macintosh'))

const mainBgStyle = computed(() => {
  if (isMac.value) {
    return { background: 'linear-gradient(145deg, rgba(15,23,42,0.82), rgba(11,17,32,0.88))' }
  }
  return { background: 'linear-gradient(145deg, #0F172A, #0B1120)' }
})

const tabs = [
  { icon: Download, label: 'Download', route: '/' },
  { icon: Package, label: 'Queue', route: '/queue' },
  { icon: Settings, label: 'Settings', route: '/settings' },
]

function isActiveTab(tabRoute) {
  return route.path === tabRoute
}

function handleKeydown(e) {
  if (e.key === 'Escape') {
    e.preventDefault()
    HideWindow()
  }
}

// Auto-hide overlay on blur (Raycast-style)
let blurTimer = null
function handleBlur() {
  // Don't hide during bootstrap (yt-dlp install)
  if (!appStore.bootstrapReady) return
  blurTimer = setTimeout(() => {
    HideWindow()
  }, 200)
}
function handleFocus() {
  if (blurTimer) {
    clearTimeout(blurTimer)
    blurTimer = null
  }
}

onMounted(async () => {
  window.addEventListener('keydown', handleKeydown)
  window.addEventListener('blur', handleBlur)
  window.addEventListener('focus', handleFocus)

  // Focus input when overlay is shown via hotkey/tray
  EventsOn('overlay:shown', () => {
    setTimeout(() => {
      const input = document.querySelector('input[type="text"]:not([disabled])')
      if (input) input.focus()
    }, 50)
  })

  await appStore.checkBootstrap()
  if (appStore.bootstrapReady) {
    downloadStore.init()
    await settingsStore.load()
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('blur', handleBlur)
  window.removeEventListener('focus', handleFocus)
})
</script>
