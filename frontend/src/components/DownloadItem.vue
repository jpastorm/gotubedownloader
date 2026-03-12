<template>
  <div
    class="px-2.5 py-2 rounded-xl transition-all duration-200 ease-spring group hover:bg-white/[0.03]"
    :class="{ 'opacity-40': item.status === 'cancelled', 'cursor-pointer': item.status === 'completed' }"
    @click="item.status === 'completed' && $emit('open', item)"
  >
    <div class="flex items-center gap-3">
      <!-- Status icon -->
      <div class="flex-shrink-0 w-6 h-6 flex items-center justify-center">
        <!-- Downloading: animated spinner with glow -->
        <div v-if="item.status === 'downloading'" class="relative">
          <div class="absolute inset-0 rounded-full bg-accent/15 blur-sm animate-pulse-glow"></div>
          <svg class="w-4.5 h-4.5 text-accent-light animate-spin-fast relative" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2.5" class="opacity-15" />
            <path d="M12 2a10 10 0 0 1 10 10" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" />
          </svg>
        </div>
        <!-- Queued -->
        <div v-else-if="item.status === 'queued'" class="w-2.5 h-2.5 rounded-full bg-white/15 animate-pulse-dot" />
        <!-- Completed -->
        <div v-else-if="item.status === 'completed'" class="w-5 h-5 rounded-full bg-emerald-500/15 flex items-center justify-center">
          <svg class="w-3 h-3 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <!-- Failed -->
        <div v-else-if="item.status === 'failed'" class="w-5 h-5 rounded-full bg-red-500/15 flex items-center justify-center">
          <svg class="w-3 h-3 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </div>
        <!-- Default -->
        <div v-else class="w-2.5 h-2.5 rounded-full bg-white/10" />
      </div>

      <!-- Info -->
      <div class="flex-1 min-w-0">
        <p class="text-[11px] text-white/70 truncate font-medium group-hover:text-white/90 transition-colors">{{ item.title || item.url || 'Unknown' }}</p>
        <div class="flex items-center gap-2 mt-0.5">
          <span :class="['text-[9px] font-mono font-medium', statusColor]">{{ statusText }}</span>
          <span v-if="item.speed" class="text-[9px] text-white/20 font-mono">{{ item.speed }}</span>
          <span v-if="item.eta && item.status === 'downloading'" class="text-[9px] text-white/15 font-mono">ETA {{ item.eta }}</span>
        </div>
      </div>

      <!-- Progress number (glowing) -->
      <div v-if="item.status === 'downloading'" class="flex-shrink-0 text-right">
        <span class="text-sm font-bold text-accent-light font-mono tabular-nums" style="text-shadow: 0 0 12px rgba(255,107,107,0.4);">{{ Math.round(item.progress || 0) }}%</span>
      </div>

      <!-- Mode badge (glass) -->
      <span class="flex-shrink-0 text-[9px] font-mono text-white/30 glass-button px-2 py-0.5 rounded-md">
        {{ item.mode === 'audio' ? 'MP3' : 'MP4' }}
      </span>

      <!-- Actions (glass buttons) -->
      <div class="flex-shrink-0 flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-all duration-200">
        <!-- Open in browser -->
        <button
          v-if="item.status === 'completed' && item.url"
          @click.stop="$emit('open-url', item)"
          class="p-1.5 rounded-lg glass-button text-white/30 hover:text-white/70 hover:!bg-white/[0.08] transition-all"
          title="Open in browser"
        >
          <ExternalLink class="w-3 h-3" />
        </button>
        <!-- Re-download -->
        <button
          v-if="item.status === 'completed'"
          @click.stop="$emit('redownload', item)"
          class="p-1.5 rounded-lg glass-button text-white/30 hover:text-accent-light hover:!bg-accent/10 hover:!border-accent/20 transition-all"
          title="Download again"
        >
          <RotateCcw class="w-3 h-3" />
        </button>
        <button
          v-if="item.status === 'downloading' || item.status === 'queued'"
          @click.stop="$emit('cancel', item.id)"
          class="p-1.5 rounded-lg glass-button hover:!bg-red-500/15 hover:!border-red-500/25 text-white/30 hover:text-red-400 transition-all"
          title="Cancel"
        >
          <X class="w-3 h-3" />
        </button>
        <button
          v-if="item.status === 'failed'"
          @click.stop="$emit('retry', item.id)"
          class="p-1.5 rounded-lg glass-button hover:!bg-amber-500/15 hover:!border-amber-500/25 text-white/30 hover:text-amber-400 transition-all"
          title="Retry"
        >
          <RotateCcw class="w-3 h-3" />
        </button>
        <button
          v-if="item.status === 'completed' || item.status === 'failed' || item.status === 'cancelled'"
          @click.stop="$emit('remove', item.id)"
          class="p-1.5 rounded-lg glass-button text-white/30 hover:text-white/60 transition-all"
          title="Remove"
        >
          <Trash2 class="w-3 h-3" />
        </button>
      </div>
    </div>

    <!-- Progress bar (gradient glow) -->
    <div v-if="item.status === 'downloading'" class="mt-2 ml-9 mr-1">
      <div class="h-[3px] rounded-full overflow-hidden bg-white/[0.04]">
        <div
          class="h-full rounded-full transition-all duration-500 ease-out relative progress-animated"
          :class="item.progress >= 100 ? 'bg-emerald-500' : ''"
          :style="{
            width: (item.progress || 0) + '%',
            background: item.progress < 100 ? 'linear-gradient(90deg, #FF3B3B, #D946EF, #FF6B6B)' : undefined,
            boxShadow: '0 0 8px rgba(255,59,59,0.3)',
          }"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { X, RotateCcw, ExternalLink, Trash2 } from 'lucide-vue-next'

const props = defineProps({
  item: { type: Object, required: true }
})

defineEmits(['cancel', 'retry', 'remove', 'open', 'open-url', 'redownload'])

const statusColor = computed(() => ({
  queued: 'text-white/30',
  downloading: 'text-accent-light',
  completed: 'text-emerald-400',
  failed: 'text-red-400',
  cancelled: 'text-white/20',
}[props.item.status] || 'text-white/30'))

const statusText = computed(() => ({
  queued: 'Queued',
  downloading: 'Downloading',
  completed: 'Done',
  failed: 'Failed',
  cancelled: 'Cancelled',
}[props.item.status] || props.item.status))
</script>
