<template>
  <div class="h-full flex flex-col overflow-hidden">
    <!-- Header -->
    <div class="flex-shrink-0 flex items-center justify-between px-4 py-2.5 border-b border-white/[0.04]">
      <div class="flex items-center gap-2.5">
        <span class="text-xs text-white/60 font-medium">History</span>
        <span v-if="historyStore.entries.length" class="text-[10px] text-white/25 font-mono glass-button px-2 py-0.5 rounded-md">
          {{ historyStore.entries.length }}
        </span>
      </div>
      <button
        v-if="historyStore.entries.length"
        @click="historyStore.clear()"
        class="px-2.5 py-1 rounded-lg glass-button text-[10px] text-white/35 hover:text-red-400 hover:!bg-red-500/10 hover:!border-red-500/20 transition-all"
      >
        <Trash2 class="w-3 h-3 inline-block mr-1 -mt-px" />
        Clear All
      </button>
    </div>

    <!-- List -->
    <div class="flex-1 overflow-y-auto p-1.5">
      <!-- Empty -->
      <div v-if="!historyStore.entries.length" class="flex flex-col items-center justify-center py-16 text-center animate-fade-in">
        <div class="w-14 h-14 rounded-2xl glass flex items-center justify-center mb-4">
          <ClipboardList class="w-6 h-6 text-white/15" />
        </div>
        <p class="text-xs text-white/25 font-light">No download history yet</p>
        <p class="text-[10px] text-white/15 mt-1">Completed downloads will appear here</p>
      </div>

      <!-- Entries -->
      <div v-else class="space-y-0.5">
        <div
          v-for="entry in historyStore.entries"
          :key="entry.id"
          class="px-3 py-2.5 rounded-xl hover:bg-white/[0.03] transition-all duration-200 group"
        >
          <div class="flex items-center gap-3">
            <!-- Thumbnail with type badge -->
            <div class="w-10 h-10 rounded-lg overflow-hidden bg-white/[0.03] flex-shrink-0 border border-white/[0.04] relative">
              <img
                v-if="entry.thumbnail"
                :src="entry.thumbnail"
                class="w-full h-full object-cover"
                @error="$event.target.style.display='none'"
              />
              <div v-else class="w-full h-full flex items-center justify-center">
                <component :is="entry.mode === 'audio' ? Music : Video" class="w-3.5 h-3.5 text-white/15" />
              </div>
              <!-- Type badge overlay -->
              <div class="absolute bottom-0 right-0 px-1 py-px text-[7px] font-bold font-mono rounded-tl-md"
                :class="entry.mode === 'audio' ? 'bg-fuchsia-500/80 text-white' : 'bg-accent/80 text-white'">
                {{ entry.mode === 'audio' ? 'MP3' : 'MP4' }}
              </div>
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
              <p class="text-xs text-white/65 truncate font-medium leading-tight">{{ entry.title || truncateUrl(entry.url) }}</p>
              <div class="flex items-center gap-1.5 mt-0.5">
                <span v-if="entry.channel" class="text-[10px] text-white/25 truncate max-w-[120px]">{{ entry.channel }}</span>
                <span v-if="entry.channel" class="text-[10px] text-white/10">·</span>
                <span v-if="entry.duration" class="text-[10px] text-white/20 font-mono">{{ entry.duration }}</span>
                <span v-if="entry.duration" class="text-[10px] text-white/10">·</span>
                <Clock class="w-2.5 h-2.5 text-white/15 flex-shrink-0" />
                <span class="text-[10px] text-white/20 font-mono">{{ formatDate(entry.downloadedAt) }}</span>
              </div>
            </div>

            <!-- Action buttons -->
            <div class="flex items-center gap-1 flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
              <!-- Open in browser -->
              <button
                v-if="entry.url"
                @click.stop="openInBrowser(entry.url)"
                class="p-1.5 rounded-lg glass-button text-white/30 hover:text-white/70 hover:!bg-white/[0.08] transition-all"
                title="Open in browser"
              >
                <ExternalLink class="w-3.5 h-3.5" />
              </button>

              <!-- Re-download -->
              <button
                @click.stop="redownload(entry)"
                class="p-1.5 rounded-lg glass-button text-white/30 hover:text-accent-light hover:!bg-accent/10 hover:!border-accent/20 transition-all"
                title="Download again"
              >
                <RotateCcw class="w-3.5 h-3.5" />
              </button>

              <!-- Delete from history -->
              <button
                @click.stop="historyStore.remove(entry.id)"
                class="p-1.5 rounded-lg glass-button text-white/30 hover:text-red-400 hover:!bg-red-500/10 hover:!border-red-500/20 transition-all"
                title="Remove from history"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useHistoryStore } from '../stores/history'
import { useDownloadStore } from '../stores/downloads'
import { OpenURL } from '../../wailsjs/go/main/App'
import { ClipboardList, ExternalLink, RotateCcw, Trash2, Clock, Music, Video } from 'lucide-vue-next'

const historyStore = useHistoryStore()
const downloadStore = useDownloadStore()
const router = useRouter()

onMounted(() => {
  if (!historyStore.loaded) {
    historyStore.load()
  }
})

function openInBrowser(url) {
  OpenURL(url)
}

async function redownload(entry) {
  try {
    await downloadStore.download(entry.url, entry.mode, entry.quality)
    router.push('/queue')
  } catch (e) {
    console.error('Re-download failed:', e)
  }
}

function truncateUrl(url) {
  if (!url) return ''
  try {
    const u = new URL(url)
    const id = u.searchParams.get('v') || u.pathname.split('/').pop()
    return u.hostname + '/…/' + id
  } catch {
    return url.length > 40 ? url.slice(0, 40) + '…' : url
  }
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diff = now - d
  if (diff < 60000) return 'just now'
  if (diff < 3600000) return Math.floor(diff / 60000) + 'm ago'
  if (diff < 86400000) return Math.floor(diff / 3600000) + 'h ago'
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}
</script>
