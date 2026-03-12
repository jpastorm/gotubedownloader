<template>
  <div class="h-screen w-screen flex items-center justify-center relative" style="background: linear-gradient(145deg, #0F172A, #0B1120);">
    <!-- Background gradient orbs -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute top-1/4 left-1/3 w-64 h-64 rounded-full opacity-[0.04] animate-float" style="background: radial-gradient(circle, #FF3B3B, transparent 70%);"></div>
      <div class="absolute bottom-1/4 right-1/3 w-48 h-48 rounded-full opacity-[0.03] animate-float" style="background: radial-gradient(circle, #D946EF, transparent 70%); animation-delay: 1.5s;"></div>
    </div>

    <div class="wails-drag absolute inset-0" />

    <div class="wails-nodrag w-80 animate-scale-bounce relative z-10">
      <!-- Icon with gradient glow -->
      <div class="flex justify-center mb-6">
        <div class="relative">
          <div class="absolute inset-0 rounded-2xl bg-accent/20 blur-xl animate-pulse-glow"></div>
          <div class="relative w-16 h-16 rounded-2xl glass flex items-center justify-center">
            <ArrowDownToLine class="w-7 h-7 text-accent-light" />
          </div>
        </div>
      </div>

      <!-- Title with gradient -->
      <h1 class="text-lg font-semibold text-center text-gradient mb-1">GoTube Downloader</h1>
      <p class="text-xs text-white/30 text-center mb-8 font-light">Checking dependencies…</p>

      <!-- Status (glass cards) -->
      <div class="space-y-2 mb-8">
        <div class="flex items-center justify-between px-4 py-3 rounded-xl glass-sm">
          <span class="text-xs text-white/60 font-medium">yt-dlp</span>
          <div class="flex items-center gap-2">
            <div :class="['w-2 h-2 rounded-full transition-all', appStore.ytdlpReady ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.5)]' : 'bg-white/20']" />
            <span class="text-[10px] text-white/30 font-mono">{{ appStore.ytdlpReady ? 'Found' : 'Missing' }}</span>
          </div>
        </div>
        <div class="flex items-center justify-between px-4 py-3 rounded-xl glass-sm">
          <span class="text-xs text-white/60 font-medium">ffmpeg</span>
          <div class="flex items-center gap-2">
            <div :class="['w-2 h-2 rounded-full transition-all', appStore.ffmpegReady ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.5)]' : 'bg-amber-400 shadow-[0_0_8px_rgba(251,191,36,0.4)]']" />
            <span class="text-[10px] text-white/30 font-mono">{{ appStore.ffmpegReady ? 'Found' : 'Optional' }}</span>
          </div>
        </div>
      </div>

      <!-- Error -->
      <p v-if="appStore.bootstrapError" class="text-xs text-red-400/80 text-center mb-4 px-4 animate-fade-in-up">
        {{ appStore.bootstrapError }}
      </p>

      <!-- Install button (gradient accent) -->
      <button
        v-if="!appStore.ytdlpReady"
        @click="appStore.installYtDlp()"
        :disabled="appStore.installing"
        class="w-full py-3 rounded-xl text-xs font-semibold transition-all duration-300 ease-spring disabled:opacity-40 relative overflow-hidden group"
        style="background: linear-gradient(135deg, #FF3B3B, #D946EF); box-shadow: 0 4px 20px rgba(255,59,59,0.3);"
      >
        <div class="absolute inset-0 bg-white/10 opacity-0 group-hover:opacity-100 transition-opacity"></div>
        <span v-if="appStore.installing" class="flex items-center justify-center gap-2 relative">
          <svg class="w-3.5 h-3.5 animate-spin-fast" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" stroke-linecap="round" class="opacity-20" />
            <path d="M12 2a10 10 0 0 1 10 10" stroke="currentColor" stroke-width="3" stroke-linecap="round" />
          </svg>
          Installing…
        </span>
        <span v-else class="relative text-white">Install yt-dlp</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { useAppStore } from '../stores/app'
import { ArrowDownToLine } from 'lucide-vue-next'
const appStore = useAppStore()
</script>
