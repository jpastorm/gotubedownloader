<template>
  <div class="animate-slide-up space-y-2.5">
    <!-- Preview card (glass) -->
    <div class="glass-sm rounded-xl overflow-hidden">
      <div class="flex gap-3 p-3">
        <!-- Thumbnail with gradient overlay -->
        <div class="w-28 h-[72px] rounded-lg overflow-hidden bg-white/[0.03] flex-shrink-0 relative group">
          <img
            v-if="info.thumbnail"
            :src="info.thumbnail"
            :alt="info.title"
            class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            @error="$event.target.style.display='none'"
          />
          <div v-else class="w-full h-full flex items-center justify-center">
            <Play class="w-5 h-5 text-white/20" />
          </div>
          <!-- Duration badge (glass) -->
          <span v-if="info.durationStr" class="absolute bottom-1 right-1 text-[9px] text-white font-mono px-1.5 py-0.5 rounded-md" style="background: rgba(0,0,0,0.7); backdrop-filter: blur(4px);">
            {{ info.durationStr }}
          </span>
          <!-- Gradient overlay -->
          <div class="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent pointer-events-none"></div>
        </div>

        <!-- Info -->
        <div class="flex-1 min-w-0 flex flex-col justify-center py-0.5">
          <h3 class="text-[11px] font-medium text-white/90 truncate leading-snug">{{ info.title }}</h3>
          <p v-if="info.channel" class="text-[10px] text-white/35 mt-0.5 truncate">{{ info.channel }}</p>
          <div class="flex items-center gap-1.5 mt-1.5">
            <span v-if="info.type === 'playlist'" class="text-[9px] text-amber-300 font-mono glass-accent px-2 py-0.5 rounded-md border-amber-400/20">
              {{ info.videoCount || info.entries?.length || 0 }} videos
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Mode selector (glass toggle) -->
    <div class="flex gap-1 p-1 rounded-xl glass-sm">
      <button
        v-for="m in modes"
        :key="m.value"
        @click="selectedMode = m.value"
        :class="[
          'flex-1 flex items-center justify-center gap-1.5 py-2 rounded-lg text-[10px] font-medium transition-all duration-200 ease-spring',
          selectedMode === m.value
            ? 'glass text-white shadow-glass-sm'
            : 'text-white/35 hover:text-white/60 hover:bg-white/[0.03]'
        ]"
      >
        <component :is="m.icon" :class="['w-3.5 h-3.5', selectedMode === m.value ? 'text-accent-light' : '']" />
        {{ m.label }}
      </button>
    </div>

    <!-- Quality actions (glass cards) -->
    <div class="glass-sm rounded-xl overflow-hidden">
      <div class="p-1">
        <button
          v-for="(q, index) in qualities"
          :key="q.value"
          @click="$emit('download', selectedMode, q.value)"
          class="w-full flex items-center justify-between px-3 py-2.5 rounded-lg hover:bg-white/[0.05] transition-all duration-200 group"
          :style="{ animationDelay: (index * 50) + 'ms' }"
        >
          <div class="flex items-center gap-3">
            <div :class="['w-7 h-7 rounded-lg flex items-center justify-center transition-all duration-200', q.bgClass, 'group-hover:scale-110']">
              <component :is="q.icon" class="w-3.5 h-3.5" />
            </div>
            <div>
              <span class="text-[11px] font-medium text-white/70 group-hover:text-white transition-colors">{{ q.label }}</span>
              <span class="text-[9px] text-white/25 ml-2 font-mono">{{ selectedMode === 'video' ? q.videoDesc : q.audioDesc }}</span>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-[9px] text-white/15 font-mono">{{ selectedMode === 'video' ? 'MP4' : 'MP3' }}</span>
            <div class="w-5 h-5 rounded-md glass-button flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all duration-200 group-hover:translate-x-0 translate-x-1">
              <ChevronRight class="w-3 h-3 text-white/50" />
            </div>
          </div>
        </button>
      </div>
      <!-- Playlist hint -->
      <p v-if="info.type === 'playlist'" class="text-[9px] text-amber-400/40 text-center pb-2 px-3">
        All {{ info.videoCount || info.entries?.length || 0 }} videos will download at the chosen quality
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Video, Music, HardDrive, Disc, Disc3, Flame, Play, ChevronRight } from 'lucide-vue-next'

const props = defineProps({
  info: {
    type: Object,
    required: true,
  },
  defaultMode: {
    type: String,
    default: 'video',
  },
})

defineEmits(['download'])

const selectedMode = ref(props.defaultMode)

const modes = computed(() => {
  if (props.defaultMode === 'audio') {
    return [
      { value: 'audio', label: 'Audio', icon: Music },
      { value: 'video', label: 'Video', icon: Video },
    ]
  }
  return [
    { value: 'video', label: 'Video', icon: Video },
    { value: 'audio', label: 'Audio', icon: Music },
  ]
})

const qualities = [
  { value: 'low', label: 'Low', icon: HardDrive, videoDesc: '480p', audioDesc: '~96kbps', bgClass: 'bg-white/[0.04] text-white/40 group-hover:bg-white/[0.08]' },
  { value: 'med', label: 'Medium', icon: Disc, videoDesc: '720p', audioDesc: '~192kbps', bgClass: 'bg-blue-500/10 text-blue-400/70 group-hover:bg-blue-500/20' },
  { value: 'high', label: 'High', icon: Disc3, videoDesc: '1080p', audioDesc: '~256kbps', bgClass: 'bg-accent/10 text-accent-light/70 group-hover:bg-accent/20' },
  { value: 'source', label: 'Source', icon: Flame, videoDesc: 'Best available', audioDesc: 'Best available', bgClass: 'bg-amber-500/10 text-amber-400/70 group-hover:bg-amber-500/20' },
]
</script>
