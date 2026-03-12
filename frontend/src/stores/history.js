import { defineStore } from 'pinia'
import { GetHistory, ClearHistory, RemoveHistoryEntry } from '../../wailsjs/go/main/App'

export const useHistoryStore = defineStore('history', {
  state: () => ({
    entries: [],
    loaded: false,
  }),
  actions: {
    async load() {
      try {
        const entries = await GetHistory()
        this.entries = entries || []
        this.loaded = true
      } catch (e) {
        console.error('Failed to load history:', e)
      }
    },
    async clear() {
      try {
        await ClearHistory()
        this.entries = []
      } catch (e) {
        console.error('Failed to clear history:', e)
      }
    },
    async remove(id) {
      try {
        await RemoveHistoryEntry(id)
        this.entries = this.entries.filter(e => e.id !== id)
      } catch (e) {
        console.error('Failed to remove history entry:', e)
      }
    }
  }
})
