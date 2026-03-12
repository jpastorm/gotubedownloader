<template>
  <div class="relative group">
    <!-- Glow effect behind input on focus -->
    <div class="absolute -inset-1 rounded-2xl bg-gradient-to-r from-accent/20 via-fuchsia-500/15 to-accent/20 opacity-0 group-focus-within:opacity-100 blur-xl transition-opacity duration-500 pointer-events-none"></div>

    <!-- Glass input bar -->
    <div class="relative glass-input rounded-xl flex items-center gap-3 px-4 py-3">
      <!-- Search icon with gradient -->
      <div class="flex-shrink-0">
        <Search :class="['w-4 h-4 transition-colors duration-200', urlInput.trim() ? 'text-accent-light' : 'text-white/25']" />
      </div>

      <input
        ref="inputRef"
        v-model="urlInput"
        @keydown.enter="handleSubmit"
        @paste="handlePaste"
        type="text"
        placeholder="Paste YouTube URL…"
        class="flex-1 bg-transparent text-sm text-white placeholder-white/20 outline-none font-light tracking-wide"
        :disabled="disabled"
        spellcheck="false"
        autocomplete="off"
      />

      <!-- Loading spinner with instant download hint -->
      <div v-if="loading" class="flex-shrink-0 flex items-center gap-2">
        <div class="relative">
          <div class="absolute inset-0 rounded-full bg-accent/20 blur-md animate-pulse-glow"></div>
          <svg class="w-4 h-4 text-accent-light animate-spin-fast relative" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" class="opacity-20" />
            <path d="M12 2a10 10 0 0 1 10 10" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" />
          </svg>
        </div>
        <kbd v-if="!isPlaylistURL(urlInput)" class="px-1.5 py-0.5 rounded-md glass-button text-[9px] text-amber-400/70 font-mono animate-pulse-subtle">↵ {{ defaultMode === 'audio' ? 'MP3' : 'MP4' }}</kbd>
      </div>

      <!-- Enter hint (glass badge) -->
      <div v-else-if="urlInput.trim()" class="flex-shrink-0 animate-scale-in">
        <kbd class="px-1.5 py-0.5 rounded-md glass-button text-[10px] text-white/50 font-mono">↵</kbd>
      </div>
    </div>

    <!-- Error message -->
    <p v-if="error" class="mt-2 text-xs text-red-400/90 px-2 animate-fade-in-up flex items-center gap-1.5">
      <span class="w-1 h-1 rounded-full bg-red-400 flex-shrink-0"></span>
      {{ error }}
    </p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Search } from 'lucide-vue-next'

const props = defineProps({
  loading: Boolean,
  error: String,
  disabled: Boolean,
  defaultMode: {
    type: String,
    default: 'video',
  },
})

const emit = defineEmits(['submit', 'instant-download'])

const urlInput = ref('')
const inputRef = ref(null)

function isPlaylistURL(url) {
  // YouTube playlist patterns: list= param, /playlist path
  return /[?&]list=/.test(url) || /\/playlist\b/.test(url)
}

function handleSubmit() {
  const url = urlInput.value.trim()
  if (!url) return
  // Playlist URLs must go through analysis — no instant download
  if (isPlaylistURL(url)) {
    emit('submit', url)
    return
  }
  if (props.loading) {
    // During analysis, Enter = skip analysis and instant download
    emit('instant-download', url)
    return
  }
  // Enter = instant single-video download
  emit('instant-download', url)
}

function handlePaste(e) {
  setTimeout(() => {
    const text = urlInput.value.trim()
    if (text && (text.includes('youtube.com') || text.includes('youtu.be') || text.includes('http'))) {
      // Paste triggers analysis (shows video info + quality picker)
      emit('submit', text)
    }
  }, 50)
}

function clear() {
  urlInput.value = ''
}

onMounted(() => {
  inputRef.value?.focus()
})

defineExpose({ clear })
</script>
