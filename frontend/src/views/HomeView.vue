<template>
  <div class="h-full flex flex-col overflow-hidden">
    <!-- URL Input -->
    <div class="flex-shrink-0 p-3 pb-0">
      <URLInput
        ref="urlInputRef"
        :loading="downloadStore.analyzing"
        :error="downloadStore.analysisError"
        :default-mode="defaultMode"
        @submit="handleAnalyze"
        @instant-download="handleInstantDownload"
      />
    </div>

    <!-- Content area -->
    <div class="flex-1 overflow-y-auto p-3">
      <!-- Playlist choice prompt -->
      <div v-if="playlistChoice" class="animate-slide-up space-y-2.5">
        <div class="glass-sm rounded-xl p-4 text-center">
          <div class="flex items-center justify-center gap-2 mb-3">
            <ListMusic class="w-5 h-5 text-amber-400/70" />
            <p class="text-xs text-white/60 font-medium">Playlist Detected</p>
          </div>
          <p class="text-[10px] text-white/30 mb-4">What do you want to download?</p>
          <div class="flex gap-2">
            <button
              @click="choosePlaylistOption('single')"
              class="flex-1 py-2.5 rounded-xl glass-button text-[11px] font-medium text-white/60 hover:text-white hover:bg-white/[0.08] transition-all duration-200 flex items-center justify-center gap-2"
            >
              <Music2 class="w-3.5 h-3.5 text-accent-light/70" />
              This Song
            </button>
            <button
              @click="choosePlaylistOption('playlist')"
              class="flex-1 py-2.5 rounded-xl glass-button text-[11px] font-medium text-white/60 hover:text-white hover:bg-white/[0.08] transition-all duration-200 flex items-center justify-center gap-2"
            >
              <ListMusic class="w-3.5 h-3.5 text-amber-400/70" />
              Full Playlist
            </button>
          </div>
        </div>
      </div>

      <!-- Analysis result with quality picker -->
      <AnalysisResult
        v-else-if="downloadStore.analysisResult"
        :info="downloadStore.analysisResult"
        :default-mode="defaultMode"
        @download="handleDownload"
      />

      <!-- Quick recent list when idle -->
      <div v-else-if="!downloadStore.analyzing" class="animate-fade-in">
        <!-- Active downloads mini status -->
        <div v-if="activeItems.length > 0" class="mb-4">
          <p class="text-[9px] text-white/20 uppercase tracking-widest font-medium px-1 mb-2">Active</p>
          <div class="space-y-1.5">
            <div
              v-for="item in activeItems"
              :key="item.id"
              class="glass-sm rounded-xl px-3 py-2.5 animate-fade-in-up"
            >
              <div class="flex items-center gap-3">
                <div class="relative flex-shrink-0">
                  <div class="absolute inset-0 rounded-full bg-accent/15 blur-sm animate-pulse-glow"></div>
                  <svg class="w-4 h-4 text-accent-light animate-spin-fast relative" viewBox="0 0 24 24" fill="none">
                    <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2.5" class="opacity-15" />
                    <path d="M12 2a10 10 0 0 1 10 10" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" />
                  </svg>
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-[10px] text-white/60 truncate font-medium">{{ item.title || item.url }}</p>
                  <!-- Progress bar (gradient) -->
                  <div class="flex items-center gap-2 mt-1.5">
                    <div class="flex-1 h-1 rounded-full overflow-hidden bg-white/[0.04]">
                      <div
                        class="h-full rounded-full transition-all duration-500 progress-animated"
                        :style="{ width: (item.progress || 0) + '%', background: 'linear-gradient(90deg, #6366f1, #8b5cf6, #a78bfa)', boxShadow: '0 0 6px rgba(99,102,241,0.3)' }"
                      />
                    </div>
                    <span class="text-[10px] text-accent-light font-mono font-bold w-9 text-right tabular-nums" style="text-shadow: 0 0 10px rgba(129,140,248,0.3);">{{ Math.round(item.progress || 0) }}%</span>
                  </div>
                </div>
                <div class="flex-shrink-0 flex flex-col items-end gap-0.5">
                  <span v-if="item.speed" class="text-[9px] text-white/20 font-mono">{{ item.speed }}</span>
                  <span v-if="item.eta" class="text-[9px] text-white/15 font-mono">{{ item.eta }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent completed -->
        <div v-if="recentDone.length > 0">
          <p class="text-[9px] text-white/20 uppercase tracking-widest font-medium px-1 mb-2">Recent</p>
          <div class="space-y-0.5">
            <div
              v-for="item in recentDone"
              :key="item.id"
              class="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-white/[0.03] transition-all duration-200 cursor-pointer"
              @click="openItem(item)"
            >
              <div class="w-5 h-5 rounded-full bg-emerald-500/10 flex items-center justify-center flex-shrink-0">
                <svg class="w-2.5 h-2.5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <p class="flex-1 text-[10px] text-white/35 truncate font-medium">{{ item.title || item.url }}</p>
              <span class="text-[9px] text-white/15 font-mono flex-shrink-0">{{ item.mode === 'audio' ? 'MP3' : 'MP4' }}</span>
            </div>
          </div>
        </div>

        <!-- Empty state (premium) -->
        <div v-if="!activeItems.length && !recentDone.length" class="flex flex-col items-center justify-center py-14 text-center">
          <div class="relative mb-4">
            <div class="absolute inset-0 rounded-2xl bg-accent/10 blur-xl animate-pulse-glow"></div>
            <div class="relative w-14 h-14 rounded-2xl glass flex items-center justify-center">
              <ArrowDownToLine class="w-6 h-6 text-white/15" />
            </div>
          </div>
          <p class="text-[11px] text-white/40 font-medium">Paste a YouTube URL</p>
          <p class="text-[10px] text-white/20 mt-1 font-light">Videos, shorts, playlists</p>
        </div>
      </div>

      <!-- Loading skeleton (glass) + active downloads -->
      <div v-else-if="downloadStore.analyzing" class="mt-2 animate-fade-in">
        <div class="glass-sm rounded-xl p-3">
          <div class="flex gap-3">
            <div class="w-28 h-[72px] rounded-lg bg-white/[0.03] animate-pulse" />
            <div class="flex-1 space-y-2.5 py-1">
              <div class="h-2.5 bg-white/[0.04] rounded-full w-3/4 animate-pulse" />
              <div class="h-2 bg-white/[0.03] rounded-full w-1/2 animate-pulse" />
            </div>
          </div>
        </div>
        <div class="flex flex-col items-center gap-2 mt-4">
          <div class="flex items-center gap-2">
            <div class="w-1 h-1 rounded-full bg-accent/40 animate-pulse-dot" />
            <p class="text-[10px] text-white/20 font-light">{{ isPlaylistURL(lastUrl) ? 'Loading playlist…' : 'Analyzing…' }}</p>
          </div>
          <p v-if="!isPlaylistURL(lastUrl)" class="text-[9px] text-amber-400/40 font-light animate-fade-in">Press <kbd class="font-mono text-amber-400/60">↵ Enter</kbd> to skip &amp; download {{ defaultMode === 'audio' ? 'MP3' : 'MP4' }} Source</p>
        </div>

        <!-- Active downloads visible during analysis -->
        <div v-if="activeItems.length > 0" class="mt-4">
          <p class="text-[9px] text-white/20 uppercase tracking-widest font-medium px-1 mb-2">Active</p>
          <div class="space-y-1.5">
            <div
              v-for="item in activeItems"
              :key="item.id"
              class="glass-sm rounded-xl px-3 py-2.5 animate-fade-in-up"
            >
              <div class="flex items-center gap-3">
                <div class="relative flex-shrink-0">
                  <div class="absolute inset-0 rounded-full bg-accent/15 blur-sm animate-pulse-glow"></div>
                  <svg class="w-4 h-4 text-accent-light animate-spin-fast relative" viewBox="0 0 24 24" fill="none">
                    <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2.5" class="opacity-15" />
                    <path d="M12 2a10 10 0 0 1 10 10" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" />
                  </svg>
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-[10px] text-white/60 truncate font-medium">{{ item.title || item.url }}</p>
                  <div class="flex items-center gap-2 mt-1.5">
                    <div class="flex-1 h-1 rounded-full overflow-hidden bg-white/[0.04]">
                      <div
                        class="h-full rounded-full transition-all duration-500 progress-animated"
                        :style="{ width: (item.progress || 0) + '%', background: 'linear-gradient(90deg, #6366f1, #8b5cf6, #a78bfa)', boxShadow: '0 0 6px rgba(99,102,241,0.3)' }"
                      />
                    </div>
                    <span class="text-[10px] text-accent-light font-mono font-bold w-9 text-right tabular-nums" style="text-shadow: 0 0 10px rgba(129,140,248,0.3);">{{ Math.round(item.progress || 0) }}%</span>
                  </div>
                </div>
                <div class="flex-shrink-0 flex flex-col items-end gap-0.5">
                  <span v-if="item.speed" class="text-[9px] text-white/20 font-mono">{{ item.speed }}</span>
                  <span v-if="item.eta" class="text-[9px] text-white/15 font-mono">{{ item.eta }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useDownloadStore } from '../stores/downloads'
import { useSettingsStore } from '../stores/settings'
import { OpenFile, OpenFolder } from '../../wailsjs/go/main/App'
import URLInput from '../components/URLInput.vue'
import AnalysisResult from '../components/AnalysisResult.vue'
import { ArrowDownToLine, ListMusic, Music2 } from 'lucide-vue-next'

const downloadStore = useDownloadStore()
const settingsStore = useSettingsStore()
const urlInputRef = ref(null)
const lastUrl = ref('')
const playlistChoice = ref(false)
const playlistMode = ref(null) // 'single' or 'playlist'

const defaultMode = computed(() => settingsStore.settings.defaultMode || 'video')

const activeItems = computed(() => {
  return downloadStore.queue.filter(i => i.status === 'downloading' || i.status === 'queued').slice(0, 3)
})

const recentDone = computed(() => {
  return downloadStore.queue.filter(i => i.status === 'completed').slice(0, 4)
})

function isPlaylistURL(url) {
  return /[?&]list=/.test(url) || /\/playlist\b/.test(url)
}

function stripPlaylistParams(url) {
  try {
    const u = new URL(url)
    u.searchParams.delete('list')
    u.searchParams.delete('index')
    u.searchParams.delete('start_radio')
    return u.toString()
  } catch {
    return url.replace(/[&?]list=[^&]*/g, '').replace(/[&?]index=[^&]*/g, '')
  }
}

async function handleAnalyze(url) {
  lastUrl.value = url
  playlistChoice.value = false
  playlistMode.value = null

  if (isPlaylistURL(url)) {
    // Show choice: single song or full playlist
    playlistChoice.value = true
    return
  }

  await downloadStore.analyze(url)
}

async function choosePlaylistOption(choice) {
  playlistChoice.value = false
  playlistMode.value = choice
  const url = lastUrl.value

  if (choice === 'single') {
    // Strip playlist params so yt-dlp sees a clean single video URL
    const cleanUrl = stripPlaylistParams(url)
    await downloadStore.analyze(cleanUrl)
  } else {
    // Analyze as playlist (uses --flat-playlist)
    await downloadStore.analyzePlaylist(url)
  }
}

async function handleInstantDownload(url) {
  const targetUrl = url || lastUrl.value
  if (!targetUrl) return
  downloadStore.cancelAnalysis()
  playlistChoice.value = false
  urlInputRef.value?.clear()
  downloadStore.download(targetUrl, defaultMode.value, 'source').catch(e => {
    console.error('Instant download failed:', e)
  })
}

async function handleDownload(mode, quality) {
  const result = downloadStore.analysisResult
  const url = result?.url || lastUrl.value
  if (!url) return

  downloadStore.analysisResult = null
  downloadStore.analysisError = ''
  urlInputRef.value?.clear()

  if (result?.type === 'playlist' && result.entries?.length > 0) {
    // Batch download all playlist entries — build full URLs from IDs if needed
    const urls = result.entries.map(e => {
      if (e.url && e.url.startsWith('http')) return e.url
      if (e.id) return `https://www.youtube.com/watch?v=${e.id}`
      return null
    }).filter(Boolean)
    downloadStore.batchDownload(urls, mode, quality).catch(e => {
      console.error('Batch download failed:', e)
    })
  } else {
    downloadStore.download(url, mode, quality).catch(e => {
      console.error('Download failed:', e)
    })
  }
}

function openItem(item) {
  if (item.filePath) {
    OpenFile(item.filePath)
  } else if (item.outputDir) {
    OpenFolder(item.outputDir)
  }
}
</script>
