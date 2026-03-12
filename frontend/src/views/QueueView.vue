<template>
  <div class="h-full flex flex-col overflow-hidden">
    <!-- Header (glass) -->
    <div class="flex-shrink-0 flex items-center justify-between px-4 py-2.5 border-b border-white/[0.04]">
      <div class="flex items-center gap-2.5">
        <span class="text-xs text-white/60 font-medium">Queue</span>
        <span v-if="downloadStore.queue.length" class="text-[10px] text-white/25 font-mono glass-button px-2 py-0.5 rounded-md">
          {{ downloadStore.queue.length }}
        </span>
      </div>
      <div class="flex items-center gap-1.5">
        <button
          v-if="downloadStore.activeCount > 0"
          @click="downloadStore.cancelAll()"
          class="px-2.5 py-1 rounded-lg glass-button text-[10px] text-white/35 hover:text-red-400 hover:!bg-red-500/10 hover:!border-red-500/20 transition-all"
        >
          Cancel All
        </button>
        <button
          v-if="downloadStore.completedCount > 0"
          @click="downloadStore.clearCompleted()"
          class="px-2.5 py-1 rounded-lg glass-button text-[10px] text-white/35 hover:text-white/60 transition-all"
        >
          Clear Done
        </button>
      </div>
    </div>

    <!-- Queue list -->
    <div class="flex-1 overflow-y-auto p-1.5">
      <DownloadQueue
        :items="downloadStore.queue"
        @cancel="downloadStore.cancel"
        @retry="downloadStore.retry"
        @remove="downloadStore.remove"
        @open="handleOpen"
        @open-url="handleOpenUrl"
        @redownload="handleRedownload"
      />
    </div>
  </div>
</template>

<script setup>
import { useDownloadStore } from '../stores/downloads'
import { OpenFile, OpenFolder, OpenURL } from '../../wailsjs/go/main/App'
import DownloadQueue from '../components/DownloadQueue.vue'

const downloadStore = useDownloadStore()

function handleOpen(item) {
  if (item.filePath) {
    OpenFile(item.filePath)
  } else if (item.outputDir) {
    OpenFolder(item.outputDir)
  }
}

function handleOpenUrl(item) {
  if (item.url) {
    OpenURL(item.url)
  }
}

function handleRedownload(item) {
  downloadStore.download(item.url, item.mode, item.options?.quality || 'high')
}
</script>
