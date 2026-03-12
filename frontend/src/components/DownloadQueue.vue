<template>
  <div>
    <!-- Empty state -->
    <div v-if="!items.length" class="flex flex-col items-center justify-center py-16 text-center animate-fade-in">
      <div class="w-14 h-14 rounded-2xl glass flex items-center justify-center mb-4">
        <Package class="w-6 h-6 text-white/15" />
      </div>
      <p class="text-xs text-white/25 font-light">No downloads in queue</p>
    </div>

    <!-- List -->
    <TransitionGroup v-else name="list" tag="div" class="space-y-0.5">
      <DownloadItem
        v-for="item in items"
        :key="item.id"
        :item="item"
        @cancel="$emit('cancel', item.id)"
        @retry="$emit('retry', item.id)"
        @remove="$emit('remove', item.id)"
        @open="$emit('open', item)"
        @open-url="$emit('open-url', item)"
        @redownload="$emit('redownload', item)"
      />
    </TransitionGroup>
  </div>
</template>

<script setup>
import DownloadItem from './DownloadItem.vue'
import { Package } from 'lucide-vue-next'

defineProps({
  items: { type: Array, default: () => [] }
})

defineEmits(['cancel', 'retry', 'remove', 'open', 'open-url', 'redownload'])
</script>
