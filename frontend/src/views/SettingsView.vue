<template>
  <div class="h-full overflow-y-auto">
    <div class="p-4 space-y-6 max-w-lg">

      <!-- Download Location -->
      <section class="animate-fade-in-up" style="animation-delay: 0ms;">
        <h3 class="text-[10px] text-white/20 uppercase tracking-widest font-medium mb-2.5 px-1">Download Location</h3>
        <button
          @click="settingsStore.selectFolder()"
          class="w-full flex items-center justify-between px-4 py-3 rounded-xl glass-sm hover:bg-white/[0.05] transition-all duration-200 group"
        >
          <div class="flex items-center gap-3 min-w-0">
            <div class="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center flex-shrink-0 group-hover:bg-accent/15 transition-colors">
              <FolderOpen class="w-4 h-4 text-accent-light/70" />
            </div>
            <span class="text-xs text-white/50 truncate group-hover:text-white/70 transition-colors">{{ settingsStore.settings.downloadDir || '~/Downloads' }}</span>
          </div>
          <ChevronRight class="w-4 h-4 text-white/15 group-hover:text-white/30 transition-colors flex-shrink-0" />
        </button>
      </section>

      <!-- Default Format -->
      <section class="animate-fade-in-up" style="animation-delay: 50ms;">
        <h3 class="text-[10px] text-white/20 uppercase tracking-widest font-medium mb-2.5 px-1">Default Format</h3>
        <div class="flex gap-2 p-1 rounded-xl glass-sm">
          <button
            @click="setMode('video')"
            :class="[
              'flex-1 py-2.5 rounded-lg text-xs font-medium transition-all duration-200 ease-spring flex items-center justify-center gap-2',
              settingsStore.settings.defaultMode === 'video'
                ? 'glass text-white shadow-glass-sm'
                : 'text-white/30 hover:text-white/50 hover:bg-white/[0.03]'
            ]"
          >
            <Video :class="['w-3.5 h-3.5', settingsStore.settings.defaultMode === 'video' ? 'text-accent-light' : '']" />
            Video
          </button>
          <button
            @click="setMode('audio')"
            :class="[
              'flex-1 py-2.5 rounded-lg text-xs font-medium transition-all duration-200 ease-spring flex items-center justify-center gap-2',
              settingsStore.settings.defaultMode === 'audio'
                ? 'glass text-white shadow-glass-sm'
                : 'text-white/30 hover:text-white/50 hover:bg-white/[0.03]'
            ]"
          >
            <Music :class="['w-3.5 h-3.5', settingsStore.settings.defaultMode === 'audio' ? 'text-accent-light' : '']" />
            Audio
          </button>
        </div>
      </section>

      <!-- Performance -->
      <section class="animate-fade-in-up" style="animation-delay: 100ms;">
        <h3 class="text-[10px] text-white/20 uppercase tracking-widest font-medium mb-2.5 px-1">Performance</h3>
        <div class="glass-sm rounded-xl overflow-hidden divide-y divide-white/[0.04]">
          <div class="flex items-center justify-between px-4 py-3">
            <span class="text-xs text-white/45">Concurrent Downloads</span>
            <div class="flex items-center gap-2">
              <button @click="adjustConcurrent(-1)" class="w-7 h-7 rounded-lg glass-button text-white/30 hover:text-white/60 text-xs flex items-center justify-center transition-all">−</button>
              <span class="text-xs text-white/60 font-mono w-5 text-center font-medium">{{ settingsStore.settings.maxConcurrent }}</span>
              <button @click="adjustConcurrent(1)" class="w-7 h-7 rounded-lg glass-button text-white/30 hover:text-white/60 text-xs flex items-center justify-center transition-all">+</button>
            </div>
          </div>
          <div class="flex items-center justify-between px-4 py-3">
            <span class="text-xs text-white/45">Fragment Threads</span>
            <div class="flex items-center gap-2">
              <button @click="adjustFragment(-1)" class="w-7 h-7 rounded-lg glass-button text-white/30 hover:text-white/60 text-xs flex items-center justify-center transition-all">−</button>
              <span class="text-xs text-white/60 font-mono w-5 text-center font-medium">{{ settingsStore.settings.fragmentConcurrent }}</span>
              <button @click="adjustFragment(1)" class="w-7 h-7 rounded-lg glass-button text-white/30 hover:text-white/60 text-xs flex items-center justify-center transition-all">+</button>
            </div>
          </div>
        </div>
      </section>

      <!-- Toggles -->
      <section class="animate-fade-in-up" style="animation-delay: 150ms;">
        <h3 class="text-[10px] text-white/20 uppercase tracking-widest font-medium mb-2.5 px-1">Options</h3>
        <div class="glass-sm rounded-xl overflow-hidden divide-y divide-white/[0.04]">
          <ToggleRow label="Skip Duplicates" v-model="settingsStore.settings.skipDuplicates" @change="save" />
          <ToggleRow label="Continue Partial Downloads" v-model="settingsStore.settings.continueDownloads" @change="save" />
          <ToggleRow label="Embed Metadata" v-model="settingsStore.settings.embedMetadata" @change="save" />
          <ToggleRow label="Download Subtitles" v-model="settingsStore.settings.downloadSubtitles" @change="save" />
          <ToggleRow label="Embed Subtitles" v-model="settingsStore.settings.embedSubtitles" @change="save" />
          <ToggleRow label="Notify on Complete" v-model="settingsStore.settings.notifyOnComplete" @change="save" />
        </div>
      </section>

      <!-- Advanced -->
      <section class="animate-fade-in-up" style="animation-delay: 200ms;">
        <h3 class="text-[10px] text-white/20 uppercase tracking-widest font-medium mb-2.5 px-1">Advanced</h3>
        <div class="glass-sm rounded-xl overflow-hidden divide-y divide-white/[0.04]">
          <TextRow label="Speed Limit" placeholder="e.g. 5M" v-model="settingsStore.settings.speedLimit" @change="save" />
          <TextRow label="Proxy URL" placeholder="socks5://127.0.0.1:1080" v-model="settingsStore.settings.proxy" @change="save" />
          <TextRow label="Subtitle Language" placeholder="en" v-model="settingsStore.settings.subtitleLang" @change="save" />
        </div>
      </section>

      <!-- Storage info -->
      <section class="pb-4 animate-fade-in-up" style="animation-delay: 250ms;">
        <h3 class="text-[10px] text-white/20 uppercase tracking-widest font-medium mb-2.5 px-1">Storage</h3>
        <div class="glass-sm rounded-xl overflow-hidden divide-y divide-white/[0.04]">
          <div class="flex items-center justify-between px-4 py-2.5">
            <span class="text-[10px] text-white/30">Config</span>
            <span class="text-[10px] text-white/15 font-mono">~/.gotubedownloader/config.json</span>
          </div>
          <div class="flex items-center justify-between px-4 py-2.5">
            <span class="text-[10px] text-white/30">History</span>
            <span class="text-[10px] text-white/15 font-mono">~/.gotubedownloader/history.json</span>
          </div>
          <div class="flex items-center justify-between px-4 py-2.5">
            <span class="text-[10px] text-white/30">Downloads</span>
            <span class="text-[10px] text-white/15 font-mono truncate max-w-[200px]">{{ settingsStore.settings.downloadDir }}</span>
          </div>
        </div>
      </section>

    </div>
  </div>
</template>

<script setup>
import { useSettingsStore } from '../stores/settings'
import { FolderOpen, Video, Music, ChevronRight } from 'lucide-vue-next'

const settingsStore = useSettingsStore()

function setMode(mode) {
  settingsStore.settings.defaultMode = mode
  save()
}

function adjustConcurrent(delta) {
  const v = settingsStore.settings.maxConcurrent + delta
  if (v >= 1 && v <= 8) {
    settingsStore.settings.maxConcurrent = v
    save()
  }
}

function adjustFragment(delta) {
  const v = settingsStore.settings.fragmentConcurrent + delta
  if (v >= 1 && v <= 16) {
    settingsStore.settings.fragmentConcurrent = v
    save()
  }
}

function save() {
  settingsStore.save()
}
</script>

<!-- Inline sub-components -->
<script>
import { defineComponent, h } from 'vue'

const ToggleRow = defineComponent({
  props: ['label', 'modelValue'],
  emits: ['update:modelValue', 'change'],
  setup(props, { emit }) {
    function toggle() {
      emit('update:modelValue', !props.modelValue)
      setTimeout(() => emit('change'), 10)
    }
    return () => h('div', {
      class: 'flex items-center justify-between px-4 py-3 hover:bg-white/[0.02] transition-colors cursor-pointer',
      onClick: toggle
    }, [
      h('span', { class: 'text-xs text-white/45' }, props.label),
      h('div', {
        class: [
          'w-9 h-5 rounded-full relative transition-all duration-300 cursor-pointer',
          props.modelValue
            ? 'shadow-[0_0_8px_rgba(255,59,59,0.3)]'
            : ''
        ].join(' '),
        style: props.modelValue
          ? 'background: linear-gradient(135deg, #FF3B3B, #D946EF)'
          : 'background: rgba(255,255,255,0.08)'
      }, [
        h('div', {
          class: [
            'absolute top-0.5 w-4 h-4 rounded-full shadow-sm transition-all duration-300',
            props.modelValue ? 'translate-x-[18px]' : 'translate-x-0.5'
          ].join(' '),
          style: 'background: white; box-shadow: 0 1px 3px rgba(0,0,0,0.3)'
        })
      ])
    ])
  }
})

const TextRow = defineComponent({
  props: ['label', 'placeholder', 'modelValue'],
  emits: ['update:modelValue', 'change'],
  setup(props, { emit }) {
    return () => h('div', {
      class: 'flex items-center justify-between px-4 py-3 hover:bg-white/[0.02] transition-colors'
    }, [
      h('span', { class: 'text-xs text-white/45 flex-shrink-0' }, props.label),
      h('input', {
        type: 'text',
        value: props.modelValue || '',
        placeholder: props.placeholder,
        class: 'w-40 text-right bg-transparent text-xs text-white/50 outline-none placeholder-white/15 font-mono focus:text-white/70 transition-colors',
        onInput: (e) => emit('update:modelValue', e.target.value),
        onBlur: () => emit('change')
      })
    ])
  }
})

export default {
  components: { ToggleRow, TextRow }
}
</script>
